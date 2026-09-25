package postgres

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/xOMAIKOx/gridworkx/services/api/internal/api"
	"github.com/xOMAIKOx/gridworkx/services/internal/company"
	"github.com/xOMAIKOx/gridworkx/services/internal/identity"
	"github.com/xOMAIKOx/gridworkx/services/internal/progression"
)

type Entropy interface{ Read([]byte) (int, error) }
type Repository struct {
	db      *sql.DB
	entropy Entropy
}

func NewForTesting(db *sql.DB, entropy Entropy) *Repository {
	return &Repository{db: db, entropy: entropy}
}

func Open(ctx context.Context, dsn string) (*Repository, error) {
	if strings.TrimSpace(dsn) == "" {
		return nil, errors.New("GRIDWORKS_DATABASE_URL is required")
	}
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(30 * time.Minute)
	db.SetConnMaxIdleTime(5 * time.Minute)
	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, err
	}
	return &Repository{db: db, entropy: rand.Reader}, nil
}
func (r *Repository) Close() error                   { return r.db.Close() }
func (r *Repository) Ping(ctx context.Context) error { return r.db.PingContext(ctx) }
func opaqueIDFrom(source Entropy, prefix string) (string, error) {
	b := make([]byte, 16)
	if _, err := io.ReadFull(source, b); err != nil {
		return "", err
	}
	return prefix + "." + hex.EncodeToString(b), nil
}
func tokenFrom(source Entropy) (string, error) {
	b := make([]byte, 32)
	if _, err := io.ReadFull(source, b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
func digest(v string) string { sum := sha256.Sum256([]byte(v)); return hex.EncodeToString(sum[:]) }
func mapDB(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return api.ErrNotFound
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505":
			return api.ErrConflict
		case "23503":
			return &api.APIError{Status: 422, Code: "principal.not_found", Message: "referenced resource is not available"}
		case "23514", "23502":
			return &api.APIError{Status: 422, Code: "mutation.invalid_state", Message: "mutation violates a domain rule"}
		case "22P02":
			return &api.APIError{Status: 400, Code: "request.invalid_value", Message: "request value is invalid"}
		case "P0001":
			return api.ErrConflict
		}
	}
	return err
}

type guestReceiptRef struct {
	AccountID string `json:"account_id"`
	PlayerID  string `json:"player_id"`
	SessionID string `json:"session_id"`
}

func decodeGuestReceiptRef(value string) (guestReceiptRef, error) {
	var ref guestReceiptRef
	decoder := json.NewDecoder(bytes.NewReader([]byte(value)))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&ref); err != nil {
		return guestReceiptRef{}, err
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		return guestReceiptRef{}, errors.New("trailing guest receipt data")
	}
	if ref.AccountID == "" || ref.PlayerID == "" || ref.SessionID == "" || strings.ContainsRune(value, 0) {
		return guestReceiptRef{}, errors.New("invalid guest receipt")
	}
	return ref, nil
}

func lockIdempotency(ctx context.Context, tx *sql.Tx, namespace, key string) error {
	_, err := tx.ExecContext(ctx, `SELECT pg_advisory_xact_lock(hashtextextended($1,0))`, namespace+"\x00"+key)
	return err
}

func (r *Repository) IssueGuest(ctx context.Context, key string, now time.Time) (api.GuestResult, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	d := digest("identity.guest_issue\x00" + key)
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return api.GuestResult{}, err
	}
	defer tx.Rollback()
	if err := lockIdempotency(ctx, tx, "identity", key); err != nil {
		return api.GuestResult{}, err
	}
	var typ, storedDigest, resultRef string
	err = tx.QueryRowContext(ctx, `SELECT mutation_type,request_digest,result_ref FROM gridworks.identity_mutation_receipts WHERE idempotency_key=$1 FOR UPDATE`, key).Scan(&typ, &storedDigest, &resultRef)
	if err == nil {
		if typ != "identity.guest_issue" || storedDigest != d {
			return api.GuestResult{}, api.ErrConflict
		}
		ref, decodeErr := decodeGuestReceiptRef(resultRef)
		if decodeErr != nil {
			return api.GuestResult{}, decodeErr
		}
		var exp time.Time
		var profile []byte
		if err := tx.QueryRowContext(ctx, `SELECT s.expires_at,json_build_object('player_id',pr.player_id,'handle',pr.handle_display,'display_name',pr.display_name,'avatar_asset_id',COALESCE(pr.avatar_asset_id,''),'bio',COALESCE(pr.bio,''),'locale',pr.locale,'timezone',pr.timezone,'visibility',pr.visibility,'dm_policy',pr.dm_policy,'discoverable',pr.discoverable,'notification_preferences',pr.notification_preferences) FROM gridworks.guest_sessions s JOIN gridworks.accounts a ON a.account_id=s.account_id JOIN gridworks.players p ON p.account_id=a.account_id JOIN gridworks.player_profiles pr ON pr.player_id=p.player_id WHERE s.session_id=$2 AND s.account_id=$1 AND p.player_id=$3`, ref.AccountID, ref.SessionID, ref.PlayerID).Scan(&exp, &profile); err != nil {
			return api.GuestResult{}, mapDB(err)
		}
		if err := tx.Commit(); err != nil {
			return api.GuestResult{}, err
		}
		return api.GuestResult{AccountID: ref.AccountID, PlayerID: ref.PlayerID, SessionID: ref.SessionID, ExpiresAt: exp, Profile: profile, Replay: true}, nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return api.GuestResult{}, err
	}
	accountID, err := opaqueIDFrom(r.entropy, "account")
	if err != nil {
		return api.GuestResult{}, err
	}
	playerID, err := opaqueIDFrom(r.entropy, "player")
	if err != nil {
		return api.GuestResult{}, err
	}
	sessionID, err := opaqueIDFrom(r.entropy, "session")
	if err != nil {
		return api.GuestResult{}, err
	}
	raw, err := tokenFrom(r.entropy)
	if err != nil {
		return api.GuestResult{}, err
	}
	exp := now.Add(30 * 24 * time.Hour)
	h, err := identity.NormalizeHandle("guest-" + accountID[len(accountID)-8:])
	if err != nil {
		return api.GuestResult{}, err
	}
	if _, err = tx.ExecContext(ctx, `INSERT INTO gridworks.accounts(account_id,status,authority_version,created_at) VALUES($1,'guest',1,$2)`, accountID, now); err != nil {
		return api.GuestResult{}, err
	}
	if _, err = tx.ExecContext(ctx, `INSERT INTO gridworks.players(player_id,account_id,created_at) VALUES($1,$2,$3)`, playerID, accountID, now); err != nil {
		return api.GuestResult{}, err
	}
	if _, err = tx.ExecContext(ctx, `INSERT INTO gridworks.player_profiles(player_id,account_id,handle_display,handle_canonical,handle_skeleton,display_name,locale,timezone,visibility,dm_policy,discoverable,notification_preferences,created_at,updated_at) VALUES($1,$2,$3,$4,$5,'Guest','en','UTC','public','everyone',true,'{}'::jsonb,$6,$6)`, playerID, accountID, h.Display, h.Canonical, h.Skeleton, now); err != nil {
		return api.GuestResult{}, err
	}
	if _, err = tx.ExecContext(ctx, `INSERT INTO gridworks.guest_sessions(session_id,account_id,token_digest,created_at,last_used_at,expires_at) VALUES($1,$2,$3,$4,$4,$5)`, sessionID, accountID, digest(raw), now, exp); err != nil {
		return api.GuestResult{}, err
	}
	refJSON, marshalErr := json.Marshal(guestReceiptRef{AccountID: accountID, PlayerID: playerID, SessionID: sessionID})
	if marshalErr != nil {
		return api.GuestResult{}, marshalErr
	}
	profile, marshalErr := json.Marshal(map[string]any{"player_id": playerID, "handle": h.Display, "display_name": "Guest", "avatar_asset_id": "", "bio": "", "locale": "en", "timezone": "UTC", "visibility": "public", "dm_policy": "everyone", "discoverable": true, "notification_preferences": map[string]bool{}})
	if marshalErr != nil {
		return api.GuestResult{}, marshalErr
	}
	resultRef = string(refJSON)
	if _, err = tx.ExecContext(ctx, `INSERT INTO gridworks.identity_mutation_receipts(idempotency_key,mutation_type,request_digest,result_ref) VALUES($1,'identity.guest_issue',$2,$3)`, key, d, resultRef); err != nil {
		return api.GuestResult{}, err
	}
	if err = tx.Commit(); err != nil {
		return api.GuestResult{}, err
	}
	return api.GuestResult{AccountID: accountID, PlayerID: playerID, SessionID: sessionID, ExpiresAt: exp, Profile: profile, RawSessionToken: raw}, nil
}
func (r *Repository) Authenticate(ctx context.Context, raw string, now time.Time) (api.Principal, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	var p api.Principal
	err := r.db.QueryRowContext(ctx, `SELECT a.account_id,p.player_id,s.session_id,a.status FROM gridworks.guest_sessions s JOIN gridworks.accounts a ON a.account_id=s.account_id JOIN gridworks.players p ON p.account_id=a.account_id WHERE s.token_digest=$1 AND s.revoked_at IS NULL AND s.expires_at>$2 AND a.status IN ('guest','protected')`, digest(raw), now).Scan(&p.AccountID, &p.PlayerID, &p.SessionID, &p.Status)
	if err != nil {
		return api.Principal{}, mapDB(err)
	}
	return p, nil
}
func (r *Repository) RevokePresentedSession(ctx context.Context, raw string, now time.Time) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	result, err := r.db.ExecContext(ctx, `UPDATE gridworks.guest_sessions SET revoked_at=COALESCE(revoked_at,$1) WHERE token_digest=$2`, now, digest(raw))
	if err != nil {
		return mapDB(err)
	}
	if n, _ := result.RowsAffected(); n == 0 {
		return api.ErrUnauthorized
	}
	return nil
}
func (r *Repository) RevokeSession(ctx context.Context, p api.Principal, now time.Time) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	_, err := r.db.ExecContext(ctx, `UPDATE gridworks.guest_sessions SET revoked_at=$1 WHERE session_id=$2 AND revoked_at IS NULL`, now, p.SessionID)
	return err
}
func (r *Repository) GetMe(ctx context.Context, p api.Principal) (api.MeView, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	var v api.MeView
	var payload []byte
	err := r.db.QueryRowContext(ctx, `SELECT a.account_id,p.player_id,a.status,json_build_object('player_id',pr.player_id,'handle',pr.handle_display,'display_name',pr.display_name,'avatar_asset_id',COALESCE(pr.avatar_asset_id,''),'bio',COALESCE(pr.bio,''),'locale',pr.locale,'timezone',pr.timezone,'visibility',pr.visibility,'dm_policy',pr.dm_policy,'discoverable',pr.discoverable,'notification_preferences',pr.notification_preferences) FROM gridworks.accounts a JOIN gridworks.players p ON p.account_id=a.account_id JOIN gridworks.player_profiles pr ON pr.player_id=p.player_id WHERE a.account_id=$1 AND p.player_id=$2`, p.AccountID, p.PlayerID).Scan(&v.AccountID, &v.PlayerID, &v.Status, &payload)
	if err != nil {
		return api.MeView{}, mapDB(err)
	}
	v.Profile = json.RawMessage(payload)
	return v, nil
}
func (r *Repository) UpdateProfile(ctx context.Context, p api.Principal, key string, patch api.ProfilePatch, now time.Time) (api.MeView, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	payload, _ := json.Marshal(patch)
	d := digest("profile.update\x00" + p.PlayerID + "\x00" + string(payload))
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return api.MeView{}, err
	}
	defer tx.Rollback()
	if err := lockIdempotency(ctx, tx, "identity", key); err != nil {
		return api.MeView{}, err
	}
	var typ, sd string
	var ref string
	err = tx.QueryRowContext(ctx, `SELECT mutation_type,request_digest,result_ref FROM gridworks.identity_mutation_receipts WHERE idempotency_key=$1 FOR UPDATE`, key).Scan(&typ, &sd, &ref)
	if err == nil {
		if typ != "profile.update" || sd != d {
			return api.MeView{}, api.ErrConflict
		}
		if err = tx.Commit(); err != nil {
			return api.MeView{}, err
		}
		return r.GetMe(ctx, p)
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return api.MeView{}, err
	}
	if patch.DisplayName != nil && !identity.ValidateDisplayName(*patch.DisplayName) {
		return api.MeView{}, errors.New("invalid display name")
	}
	if patch.Locale != nil && !identity.ValidateLocale(*patch.Locale) {
		return api.MeView{}, errors.New("invalid locale")
	}
	if patch.Timezone != nil && !identity.ValidateTimezone(*patch.Timezone) {
		return api.MeView{}, errors.New("invalid timezone")
	}
	var prefs any
	if patch.NotificationPreferences != nil {
		encoded, marshalErr := json.Marshal(patch.NotificationPreferences)
		if marshalErr != nil {
			return api.MeView{}, marshalErr
		}
		prefs = encoded
	}
	_, err = tx.ExecContext(ctx, `UPDATE gridworks.player_profiles SET display_name=COALESCE($1,display_name),bio=COALESCE($2,bio),locale=COALESCE($3,locale),timezone=COALESCE($4,timezone),visibility=COALESCE($5,visibility),dm_policy=COALESCE($6,dm_policy),discoverable=COALESCE($7,discoverable),notification_preferences=COALESCE($8::jsonb,notification_preferences),updated_at=$9 WHERE player_id=$10 AND account_id=$11`, patch.DisplayName, patch.Bio, patch.Locale, patch.Timezone, patch.Visibility, patch.DMPolicy, patch.Discoverable, prefs, now, p.PlayerID, p.AccountID)
	if err != nil {
		return api.MeView{}, err
	}
	if _, err = tx.ExecContext(ctx, `INSERT INTO gridworks.identity_mutation_receipts(idempotency_key,mutation_type,request_digest,result_ref) VALUES($1,'profile.update',$2,$3)`, key, d, p.PlayerID); err != nil {
		return api.MeView{}, err
	}
	if err = tx.Commit(); err != nil {
		return api.MeView{}, err
	}
	return r.GetMe(ctx, p)
}
func (r *Repository) GetPublicPlayer(ctx context.Context, id string) (api.PublicPlayerView, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	var v api.PublicPlayerView
	err := r.db.QueryRowContext(ctx, `SELECT player_id,handle_display,display_name,COALESCE(avatar_asset_id,''),COALESCE(bio,''),locale,visibility,discoverable FROM gridworks.player_profiles WHERE player_id=$1 AND visibility='public' AND discoverable=true`, id).Scan(&v.PlayerID, &v.Handle, &v.DisplayName, &v.AvatarAssetID, &v.Bio, &v.Locale, &v.Visibility, &v.Discoverable)
	if err != nil {
		return api.PublicPlayerView{}, mapDB(err)
	}
	return v, nil
}
func (r *Repository) CreateCompany(ctx context.Context, p api.Principal, key string, req api.CompanyCreateRequest, now time.Time) (api.CreatedEntity, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if req.CompanyType == "system" {
		return api.CreatedEntity{}, api.ErrForbidden
	}
	displayName := strings.TrimSpace(req.Name)
	name, err := company.NormalizeBusinessName(displayName)
	if err != nil {
		return api.CreatedEntity{}, err
	}
	d := digest("company.create\x00" + req.CompanyType + "\x00" + req.Name + "\x00player\x00" + p.PlayerID)
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return api.CreatedEntity{}, err
	}
	defer tx.Rollback()
	if err := lockIdempotency(ctx, tx, "company", key); err != nil {
		return api.CreatedEntity{}, err
	}
	var typ, storedDigest, existingID string
	err = tx.QueryRowContext(ctx, `SELECT mutation_type,request_digest,result_ref FROM gridworks.company_mutation_receipts WHERE idempotency_key=$1 FOR UPDATE`, key).Scan(&typ, &storedDigest, &existingID)
	if err == nil {
		if typ != "company.create" || storedDigest != d {
			return api.CreatedEntity{}, api.ErrConflict
		}
		var existingName string
		if err := tx.QueryRowContext(ctx, `SELECT name_display FROM gridworks.companies WHERE company_id=$1`, existingID).Scan(&existingName); err != nil {
			return api.CreatedEntity{}, mapDB(err)
		}
		if err := tx.Commit(); err != nil {
			return api.CreatedEntity{}, err
		}
		return api.CreatedEntity{ID: existingID, Name: existingName}, nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return api.CreatedEntity{}, err
	}
	id, err := opaqueIDFrom(r.entropy, "company")
	if err != nil {
		return api.CreatedEntity{}, err
	}
	if _, err = tx.ExecContext(ctx, `INSERT INTO gridworks.companies(company_id,company_type,name_display,name_canonical,name_skeleton,status,visibility,created_at,updated_at) VALUES($1,$2,$3,$4,$5,'active','public',$6,$6)`, id, req.CompanyType, displayName, name.Canonical, name.Skeleton, now); err != nil {
		return api.CreatedEntity{}, mapDB(err)
	}
	oid, err := opaqueIDFrom(r.entropy, "ownership")
	if err != nil {
		return api.CreatedEntity{}, err
	}
	eid, err := opaqueIDFrom(r.entropy, "ownership-event")
	if err != nil {
		return api.CreatedEntity{}, err
	}
	if _, err = tx.ExecContext(ctx, `INSERT INTO gridworks.company_ownership(ownership_id,entity_id,entity_type,owner_type,owner_id,share_bps,effective_from) VALUES($1,$2,'company','player',$3,10000,$4)`, oid, id, p.PlayerID, now); err != nil {
		return api.CreatedEntity{}, mapDB(err)
	}
	if _, err = tx.ExecContext(ctx, `INSERT INTO gridworks.ownership_history(event_id,entity_id,mutation_type,to_owner_type,to_owner_id,share_bps,effective_time,source_ref,actor_player_id) VALUES($1,$2,'company.create','player',$3,10000,$4,$5,$3)`, eid, id, p.PlayerID, now, key); err != nil {
		return api.CreatedEntity{}, err
	}
	if _, err = tx.ExecContext(ctx, `INSERT INTO gridworks.company_mutation_receipts(idempotency_key,mutation_type,request_digest,result_ref) VALUES($1,'company.create',$2,$3)`, key, d, id); err != nil {
		return api.CreatedEntity{}, err
	}
	if err = tx.Commit(); err != nil {
		return api.CreatedEntity{}, err
	}
	return api.CreatedEntity{ID: id, Name: displayName}, nil
}
func (r *Repository) GetPublicCompany(ctx context.Context, id string) (api.PublicCompanyView, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	var v api.PublicCompanyView
	err := r.db.QueryRowContext(ctx, `SELECT company_id,name_display,COALESCE(description,''),COALESCE(logo_asset_id,''),COALESCE(industry_id,''),COALESCE(headquarters_region_id,''),visibility FROM gridworks.companies WHERE company_id=$1 AND visibility='public'`, id).Scan(&v.CompanyID, &v.Name, &v.Description, &v.LogoAssetID, &v.IndustryID, &v.HeadquartersRegionID, &v.Visibility)
	if err != nil {
		return api.PublicCompanyView{}, mapDB(err)
	}
	return v, nil
}
func (r *Repository) CreateGroup(ctx context.Context, p api.Principal, key string, req api.GroupCreateRequest, now time.Time) (api.CreatedEntity, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	displayName := strings.TrimSpace(req.Name)
	name, err := company.NormalizeBusinessName(displayName)
	if err != nil {
		return api.CreatedEntity{}, err
	}
	d := digest("group.create\x00" + req.Name + "\x00player\x00" + p.PlayerID)
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return api.CreatedEntity{}, err
	}
	defer tx.Rollback()
	if err := lockIdempotency(ctx, tx, "company", key); err != nil {
		return api.CreatedEntity{}, err
	}
	var typ, storedDigest, existingID string
	err = tx.QueryRowContext(ctx, `SELECT mutation_type,request_digest,result_ref FROM gridworks.company_mutation_receipts WHERE idempotency_key=$1 FOR UPDATE`, key).Scan(&typ, &storedDigest, &existingID)
	if err == nil {
		if typ != "group.create" || storedDigest != d {
			return api.CreatedEntity{}, api.ErrConflict
		}
		var existingName string
		if err := tx.QueryRowContext(ctx, `SELECT name_display FROM gridworks.company_groups WHERE group_id=$1`, existingID).Scan(&existingName); err != nil {
			return api.CreatedEntity{}, mapDB(err)
		}
		if err := tx.Commit(); err != nil {
			return api.CreatedEntity{}, err
		}
		return api.CreatedEntity{ID: existingID, Name: existingName}, nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return api.CreatedEntity{}, err
	}
	id, err := opaqueIDFrom(r.entropy, "group")
	if err != nil {
		return api.CreatedEntity{}, err
	}
	if _, err = tx.ExecContext(ctx, `INSERT INTO gridworks.company_groups(group_id,name_display,name_canonical,name_skeleton,status,created_at,updated_at) VALUES($1,$2,$3,$4,'active',$5,$5)`, id, displayName, name.Canonical, name.Skeleton, now); err != nil {
		return api.CreatedEntity{}, mapDB(err)
	}
	oid, err := opaqueIDFrom(r.entropy, "ownership")
	if err != nil {
		return api.CreatedEntity{}, err
	}
	eid, err := opaqueIDFrom(r.entropy, "ownership-event")
	if err != nil {
		return api.CreatedEntity{}, err
	}
	if _, err = tx.ExecContext(ctx, `INSERT INTO gridworks.company_ownership(ownership_id,entity_id,entity_type,owner_type,owner_id,share_bps,effective_from) VALUES($1,$2,'group','player',$3,10000,$4)`, oid, id, p.PlayerID, now); err != nil {
		return api.CreatedEntity{}, mapDB(err)
	}
	if _, err = tx.ExecContext(ctx, `INSERT INTO gridworks.ownership_history(event_id,entity_id,mutation_type,to_owner_type,to_owner_id,share_bps,effective_time,source_ref,actor_player_id) VALUES($1,$2,'group.create','player',$3,10000,$4,$5,$3)`, eid, id, p.PlayerID, now, key); err != nil {
		return api.CreatedEntity{}, err
	}
	if _, err = tx.ExecContext(ctx, `INSERT INTO gridworks.company_mutation_receipts(idempotency_key,mutation_type,request_digest,result_ref) VALUES($1,'group.create',$2,$3)`, key, d, id); err != nil {
		return api.CreatedEntity{}, err
	}
	if err = tx.Commit(); err != nil {
		return api.CreatedEntity{}, err
	}
	return api.CreatedEntity{ID: id, Name: displayName}, nil
}

func (r *Repository) GetPlayerSkills(ctx context.Context, p api.Principal) ([]progression.PlayerSkillView, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	rows, err := r.db.QueryContext(ctx, `SELECT player_id,skill_id,cumulative_xp,proficiency_bps,progression_version FROM gridworks.player_skills WHERE player_id=$1 ORDER BY skill_id`, p.PlayerID)
	if err != nil {
		return nil, mapDB(err)
	}
	defer rows.Close()
	views := make([]progression.PlayerSkillView, 0)
	for rows.Next() {
		var view progression.PlayerSkillView
		if err := rows.Scan(&view.PlayerID, &view.SkillID, &view.CumulativeXP, &view.ProficiencyBPS, &view.ProgressionVersion); err != nil {
			return nil, err
		}
		if err := progression.ValidatePlayerSkillView(view); err != nil {
			return nil, err
		}
		views = append(views, view)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return views, nil
}

func (r *Repository) authorizeCompany(ctx context.Context, p api.Principal, companyID string) error {
	var owner bool
	err := r.db.QueryRowContext(ctx, `SELECT EXISTS (SELECT 1 FROM gridworks.company_ownership WHERE entity_id=$1 AND entity_type='company' AND owner_type='player' AND owner_id=$2 AND active)`, companyID, p.PlayerID).Scan(&owner)
	if err != nil {
		return mapDB(err)
	}
	if !owner {
		return api.ErrForbidden
	}
	return nil
}

func (r *Repository) GetCompanyManagers(ctx context.Context, p api.Principal, companyID string) ([]progression.ManagerView, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := r.authorizeCompany(ctx, p, companyID); err != nil {
		return nil, err
	}
	rows, err := r.db.QueryContext(ctx, `SELECT m.manager_id,m.display_name,m.rarity,m.specialization_id,m.total_xp,m.level,m.workload_bps,m.fatigue_bps,m.morale_bps,m.status,COALESCE((SELECT jsonb_agg(jsonb_build_object('skill_id',ms.skill_id,'proficiency_bps',ms.proficiency_bps,'potential_bps',ms.potential_bps) ORDER BY ms.skill_id) FROM gridworks.manager_skills ms WHERE ms.manager_id=m.manager_id),'[]'::jsonb),COALESCE((SELECT jsonb_agg(mt.trait_id ORDER BY mt.trait_id) FROM gridworks.manager_traits mt WHERE mt.manager_id=m.manager_id),'[]'::jsonb),COALESCE((SELECT mfa.facility_id FROM gridworks.manager_facility_assignments mfa WHERE mfa.manager_id=m.manager_id AND mfa.active AND mfa.assignment_role='primary' ORDER BY mfa.effective_from DESC LIMIT 1),'') FROM gridworks.managers m JOIN gridworks.manager_employment_history eh ON eh.manager_id=m.manager_id AND eh.company_id=$1 AND eh.state='active' ORDER BY m.manager_id`, companyID)
	if err != nil {
		return nil, mapDB(err)
	}
	defer rows.Close()
	views := make([]progression.ManagerView, 0)
	for rows.Next() {
		var view progression.ManagerView
		var skillsJSON, traitsJSON []byte
		if err := rows.Scan(&view.ManagerID, &view.DisplayName, &view.Rarity, &view.SpecializationID, &view.TotalXP, &view.Level, &view.WorkloadBPS, &view.FatigueBPS, &view.MoraleBPS, &view.Status, &skillsJSON, &traitsJSON, &view.ActiveFacilityID); err != nil {
			return nil, err
		}
		if err := json.Unmarshal(skillsJSON, &view.Skills); err != nil {
			return nil, err
		}
		if err := json.Unmarshal(traitsJSON, &view.Traits); err != nil {
			return nil, err
		}
		view.ActiveEmployment = true
		view.PotentialVisible = true
		if err := progression.ValidateManagerView(view); err != nil {
			return nil, err
		}
		views = append(views, view)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return views, nil
}

func (r *Repository) GetCompanyManager(ctx context.Context, p api.Principal, companyID, managerID string) (progression.ManagerView, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := r.authorizeCompany(ctx, p, companyID); err != nil {
		return progression.ManagerView{}, err
	}
	rows, err := r.GetCompanyManagers(ctx, p, companyID)
	if err != nil {
		return progression.ManagerView{}, err
	}
	for _, view := range rows {
		if view.ManagerID == managerID {
			return view, nil
		}
	}
	return progression.ManagerView{}, api.ErrNotFound
}

func (r *Repository) ApplyPlayerSkillEvent(ctx context.Context, event progression.SkillEvent) (progression.SkillAward, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return progression.SkillAward{}, err
	}
	defer tx.Rollback()
	if err := lockIdempotency(ctx, tx, "progression", event.SourceEventID); err != nil {
		return progression.SkillAward{}, err
	}
	requestDigest := digest(strings.Join([]string{event.PlayerID, event.SkillID, event.ActivityKind, fmt.Sprint(event.BaseXP), string(event.ActivityClass), event.RepetitionKey, event.OccurrenceTime.UTC().Format(time.RFC3339Nano), event.RulesVersion}, "\x00"))
	var existingAward int64
	var storedDigest string
	err = tx.QueryRowContext(ctx, `SELECT awarded_xp,request_digest FROM gridworks.player_skill_events WHERE source_event_id=$1 FOR UPDATE`, event.SourceEventID).Scan(&existingAward, &storedDigest)
	if err == nil {
		if storedDigest != requestDigest {
			return progression.SkillAward{}, api.ErrConflict
		}
		var cumulative int64
		var proficiency, repetitions int
		if err := tx.QueryRowContext(ctx, `SELECT cumulative_xp,proficiency_bps,repetition_count FROM gridworks.player_skills WHERE player_id=$1 AND skill_id=$2`, event.PlayerID, event.SkillID).Scan(&cumulative, &proficiency, &repetitions); err != nil {
			return progression.SkillAward{}, mapDB(err)
		}
		if err := tx.Commit(); err != nil {
			return progression.SkillAward{}, err
		}
		return progression.SkillAward{AwardedXP: 0, CumulativeXP: cumulative, ProficiencyBPS: proficiency, RepetitionCount: repetitions, Replay: true}, nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return progression.SkillAward{}, mapDB(err)
	}
	var state progression.SkillState
	err = tx.QueryRowContext(ctx, `SELECT player_id,skill_id,cumulative_xp,proficiency_bps,progression_version,COALESCE(repetition_key,''),repetition_count FROM gridworks.player_skills WHERE player_id=$1 AND skill_id=$2 FOR UPDATE`, event.PlayerID, event.SkillID).Scan(&state.PlayerID, &state.SkillID, &state.CumulativeXP, &state.ProficiencyBPS, &state.ProgressionVersion, &state.RepetitionKey, &state.RepetitionCount)
	if errors.Is(err, sql.ErrNoRows) {
		state = progression.SkillState{PlayerID: event.PlayerID, SkillID: event.SkillID, ProgressionVersion: progression.Version}
	} else if err != nil {
		return progression.SkillAward{}, mapDB(err)
	}
	award, err := progression.ApplySkillEvent(&state, event, false)
	if err != nil {
		return progression.SkillAward{}, err
	}
	_, err = tx.ExecContext(ctx, `INSERT INTO gridworks.player_skills(player_id,skill_id,cumulative_xp,proficiency_bps,progression_version,last_source_event_id,last_event_time,repetition_key,repetition_count) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9) ON CONFLICT (player_id,skill_id) DO UPDATE SET cumulative_xp=EXCLUDED.cumulative_xp,proficiency_bps=EXCLUDED.proficiency_bps,progression_version=EXCLUDED.progression_version,last_source_event_id=EXCLUDED.last_source_event_id,last_event_time=EXCLUDED.last_event_time,repetition_key=EXCLUDED.repetition_key,repetition_count=EXCLUDED.repetition_count`, state.PlayerID, state.SkillID, state.CumulativeXP, state.ProficiencyBPS, state.ProgressionVersion, event.SourceEventID, event.OccurrenceTime, state.RepetitionKey, state.RepetitionCount)
	if err != nil {
		return progression.SkillAward{}, mapDB(err)
	}
	_, err = tx.ExecContext(ctx, `INSERT INTO gridworks.player_skill_events(source_event_id,player_id,skill_id,activity_kind,base_xp,awarded_xp,request_digest,activity_class,repetition_key,occurrence_time,rules_version) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`, event.SourceEventID, event.PlayerID, event.SkillID, event.ActivityKind, event.BaseXP, award.AwardedXP, requestDigest, event.ActivityClass, event.RepetitionKey, event.OccurrenceTime, event.RulesVersion)
	if err != nil {
		return progression.SkillAward{}, mapDB(err)
	}
	if err := tx.Commit(); err != nil {
		return progression.SkillAward{}, err
	}
	return award, nil
}

func (r *Repository) OpenManagerEmployment(ctx context.Context, employmentID, managerID, companyID, role, sourceRef string, effectiveFrom time.Time) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	key := "employment.open." + employmentID
	requestDigest := digest(strings.Join([]string{managerID, companyID, role, sourceRef, effectiveFrom.UTC().Format(time.RFC3339Nano)}, "\x00"))
	if err := lockIdempotency(ctx, tx, "progression", key); err != nil {
		return err
	}
	var mutationType, storedDigest, resultRef string
	err = tx.QueryRowContext(ctx, `SELECT mutation_type,request_digest,result_ref FROM gridworks.manager_mutation_receipts WHERE idempotency_key=$1 FOR UPDATE`, key).Scan(&mutationType, &storedDigest, &resultRef)
	if err == nil {
		if mutationType != "employment.open" || storedDigest != requestDigest {
			return api.ErrConflict
		}
		if err := tx.Commit(); err != nil {
			return err
		}
		return nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return mapDB(err)
	}
	if err := tx.QueryRowContext(ctx, `SELECT manager_id FROM gridworks.managers WHERE manager_id=$1 FOR UPDATE`, managerID).Scan(new(string)); err != nil {
		return mapDB(err)
	}
	var active int
	if err := tx.QueryRowContext(ctx, `SELECT count(*) FROM gridworks.manager_employment_history WHERE manager_id=$1 AND state='active'`, managerID).Scan(&active); err != nil {
		return mapDB(err)
	}
	if active > 0 {
		return progression.ErrActiveEmployment
	}
	if _, err = tx.ExecContext(ctx, `INSERT INTO gridworks.manager_employment_history(employment_id,manager_id,company_id,role,state,effective_from,source_ref) VALUES($1,$2,$3,$4,'active',$5,$6)`, employmentID, managerID, companyID, role, effectiveFrom, sourceRef); err != nil {
		return mapDB(err)
	}
	if _, err = tx.ExecContext(ctx, `INSERT INTO gridworks.manager_mutation_receipts(idempotency_key,mutation_type,request_digest,result_ref) VALUES($1,'employment.open',$2,$3)`, key, requestDigest, employmentID); err != nil {
		return mapDB(err)
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (r *Repository) CloseManagerEmployment(ctx context.Context, employmentID, sourceRef string, effectiveTo time.Time) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	key := "employment.close." + employmentID
	requestDigest := digest(strings.Join([]string{employmentID, sourceRef, effectiveTo.UTC().Format(time.RFC3339Nano)}, "\x00"))
	if err := lockIdempotency(ctx, tx, "progression", key); err != nil {
		return err
	}
	var mutationType, storedDigest, resultRef string
	err = tx.QueryRowContext(ctx, `SELECT mutation_type,request_digest,result_ref FROM gridworks.manager_mutation_receipts WHERE idempotency_key=$1 FOR UPDATE`, key).Scan(&mutationType, &storedDigest, &resultRef)
	if err == nil {
		if mutationType != "employment.close" || storedDigest != requestDigest {
			return api.ErrConflict
		}
		if err := tx.Commit(); err != nil {
			return err
		}
		return nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return mapDB(err)
	}
	result, err := tx.ExecContext(ctx, `UPDATE gridworks.manager_employment_history SET state='closed',effective_to=$1 WHERE employment_id=$2 AND state='active'`, effectiveTo, employmentID)
	if err != nil {
		return mapDB(err)
	}
	n, _ := result.RowsAffected()
	if n == 0 {
		return progression.ErrNoActiveEmployment
	}
	if _, err = tx.ExecContext(ctx, `INSERT INTO gridworks.manager_mutation_receipts(idempotency_key,mutation_type,request_digest,result_ref) VALUES($1,'employment.close',$2,$3)`, key, requestDigest, employmentID); err != nil {
		return mapDB(err)
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (r *Repository) ApplyManagerProgression(ctx context.Context, event progression.ManagerProgressionEvent) (progression.ManagerProgressionResult, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return progression.ManagerProgressionResult{}, err
	}
	defer tx.Rollback()
	key := "manager.progression." + event.SourceEventID
	requestDigest := digest(strings.Join([]string{event.ManagerID, event.ActivityKind, event.SkillID, fmt.Sprint(event.AwardedXP), event.OccurrenceTime.UTC().Format(time.RFC3339Nano), event.RulesVersion}, "\x00"))
	if err := lockIdempotency(ctx, tx, "progression", key); err != nil {
		return progression.ManagerProgressionResult{}, err
	}
	var mutationType, storedDigest, resultRef string
	err = tx.QueryRowContext(ctx, `SELECT mutation_type,request_digest,result_ref FROM gridworks.manager_mutation_receipts WHERE idempotency_key=$1 FOR UPDATE`, key).Scan(&mutationType, &storedDigest, &resultRef)
	if err == nil {
		if mutationType != "manager.progression" || storedDigest != requestDigest {
			return progression.ManagerProgressionResult{}, api.ErrConflict
		}
		var totalXP int64
		var level int
		if err := tx.QueryRowContext(ctx, `SELECT total_xp,level FROM gridworks.managers WHERE manager_id=$1`, event.ManagerID).Scan(&totalXP, &level); err != nil {
			return progression.ManagerProgressionResult{}, mapDB(err)
		}
		if err := tx.Commit(); err != nil {
			return progression.ManagerProgressionResult{}, err
		}
		return progression.ManagerProgressionResult{TotalXP: totalXP, Level: level, Replay: true}, nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return progression.ManagerProgressionResult{}, mapDB(err)
	}
	var state progression.ManagerState
	if err := tx.QueryRowContext(ctx, `SELECT manager_id,rarity,total_xp,level FROM gridworks.managers WHERE manager_id=$1 FOR UPDATE`, event.ManagerID).Scan(&state.ManagerID, &state.Rarity, &state.TotalXP, &state.Level); err != nil {
		return progression.ManagerProgressionResult{}, mapDB(err)
	}
	rows, err := tx.QueryContext(ctx, `SELECT skill_id,proficiency_bps,potential_bps FROM gridworks.manager_skills WHERE manager_id=$1 FOR UPDATE`, event.ManagerID)
	if err != nil {
		return progression.ManagerProgressionResult{}, err
	}
	for rows.Next() {
		var skill progression.ManagerSkillView
		if err := rows.Scan(&skill.SkillID, &skill.ProficiencyBPS, &skill.PotentialBPS); err != nil {
			rows.Close()
			return progression.ManagerProgressionResult{}, err
		}
		state.Skills = append(state.Skills, skill)
	}
	rows.Close()
	result, err := progression.ApplyManagerProgression(&state, event, false)
	if err != nil {
		return progression.ManagerProgressionResult{}, err
	}
	if _, err = tx.ExecContext(ctx, `UPDATE gridworks.managers SET total_xp=$1,level=$2 WHERE manager_id=$3`, state.TotalXP, state.Level, event.ManagerID); err != nil {
		return progression.ManagerProgressionResult{}, mapDB(err)
	}
	if event.SkillID != "" {
		for _, skill := range state.Skills {
			if skill.SkillID == event.SkillID {
				if _, err = tx.ExecContext(ctx, `UPDATE gridworks.manager_skills SET proficiency_bps=$1 WHERE manager_id=$2 AND skill_id=$3`, skill.ProficiencyBPS, event.ManagerID, event.SkillID); err != nil {
					return progression.ManagerProgressionResult{}, mapDB(err)
				}
			}
		}
	}
	if _, err = tx.ExecContext(ctx, `INSERT INTO gridworks.manager_progression_events(source_event_id,manager_id,activity_kind,awarded_xp,skill_id,occurrence_time,rules_version) VALUES($1,$2,$3,$4,$5,$6,$7)`, event.SourceEventID, event.ManagerID, event.ActivityKind, event.AwardedXP, event.SkillID, event.OccurrenceTime, event.RulesVersion); err != nil {
		return progression.ManagerProgressionResult{}, mapDB(err)
	}
	if _, err = tx.ExecContext(ctx, `INSERT INTO gridworks.manager_mutation_receipts(idempotency_key,mutation_type,request_digest,result_ref) VALUES($1,'manager.progression',$2,$3)`, key, requestDigest, event.ManagerID); err != nil {
		return progression.ManagerProgressionResult{}, mapDB(err)
	}
	if err := tx.Commit(); err != nil {
		return progression.ManagerProgressionResult{}, err
	}
	return result, nil
}

func (r *Repository) OpenManagerFacilityAssignment(ctx context.Context, assignmentID, managerID, companyID, facilityID, role, sourceRef string, effectiveFrom time.Time) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	key := "assignment.open." + assignmentID
	requestDigest := digest(strings.Join([]string{managerID, companyID, facilityID, role, sourceRef, effectiveFrom.UTC().Format(time.RFC3339Nano)}, "\x00"))
	if err := lockIdempotency(ctx, tx, "progression", key); err != nil {
		return err
	}
	var mutationType, storedDigest, resultRef string
	err = tx.QueryRowContext(ctx, `SELECT mutation_type,request_digest,result_ref FROM gridworks.manager_mutation_receipts WHERE idempotency_key=$1 FOR UPDATE`, key).Scan(&mutationType, &storedDigest, &resultRef)
	if err == nil {
		if mutationType != "assignment.open" || storedDigest != requestDigest {
			return api.ErrConflict
		}
		if err := tx.Commit(); err != nil {
			return err
		}
		return nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return mapDB(err)
	}
	if err := tx.QueryRowContext(ctx, `SELECT manager_id FROM gridworks.managers WHERE manager_id=$1 FOR UPDATE`, managerID).Scan(new(string)); err != nil {
		return mapDB(err)
	}
	var activeEmployer string
	if err := tx.QueryRowContext(ctx, `SELECT company_id FROM gridworks.manager_employment_history WHERE manager_id=$1 AND state='active'`, managerID).Scan(&activeEmployer); err != nil {
		return mapDB(err)
	}
	if activeEmployer != companyID {
		return api.ErrForbidden
	}
	if role == "primary" {
		var count int
		if err := tx.QueryRowContext(ctx, `SELECT count(*) FROM gridworks.manager_facility_assignments WHERE manager_id=$1 AND assignment_role='primary' AND active`, managerID).Scan(&count); err != nil {
			return mapDB(err)
		}
		if count > 0 {
			return api.ErrConflict
		}
	}
	if _, err = tx.ExecContext(ctx, `INSERT INTO gridworks.manager_facility_assignments(assignment_id,manager_id,company_id,facility_id,assignment_role,active,effective_from,source_ref) VALUES($1,$2,$3,$4,$5,true,$6,$7)`, assignmentID, managerID, companyID, facilityID, role, effectiveFrom, sourceRef); err != nil {
		return mapDB(err)
	}
	if _, err = tx.ExecContext(ctx, `INSERT INTO gridworks.manager_mutation_receipts(idempotency_key,mutation_type,request_digest,result_ref) VALUES($1,'assignment.open',$2,$3)`, key, requestDigest, assignmentID); err != nil {
		return mapDB(err)
	}
	return tx.Commit()
}

func (r *Repository) CloseManagerFacilityAssignment(ctx context.Context, assignmentID, sourceRef string, effectiveTo time.Time) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	key := "assignment.close." + assignmentID
	requestDigest := digest(strings.Join([]string{assignmentID, sourceRef, effectiveTo.UTC().Format(time.RFC3339Nano)}, "\x00"))
	if err := lockIdempotency(ctx, tx, "progression", key); err != nil {
		return err
	}
	var mutationType, storedDigest, resultRef string
	err = tx.QueryRowContext(ctx, `SELECT mutation_type,request_digest,result_ref FROM gridworks.manager_mutation_receipts WHERE idempotency_key=$1 FOR UPDATE`, key).Scan(&mutationType, &storedDigest, &resultRef)
	if err == nil {
		if mutationType != "assignment.close" || storedDigest != requestDigest {
			return api.ErrConflict
		}
		if err := tx.Commit(); err != nil {
			return err
		}
		return nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return mapDB(err)
	}
	result, err := tx.ExecContext(ctx, `UPDATE gridworks.manager_facility_assignments SET active=false,effective_to=$1 WHERE assignment_id=$2 AND active`, effectiveTo, assignmentID)
	if err != nil {
		return mapDB(err)
	}
	n, _ := result.RowsAffected()
	if n == 0 {
		return api.ErrNotFound
	}
	if _, err = tx.ExecContext(ctx, `INSERT INTO gridworks.manager_mutation_receipts(idempotency_key,mutation_type,request_digest,result_ref) VALUES($1,'assignment.close',$2,$3)`, key, requestDigest, assignmentID); err != nil {
		return mapDB(err)
	}
	return tx.Commit()
}
