package postgres

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"strings"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/xOMAIKOx/gridworkx/services/api/internal/api"
	"github.com/xOMAIKOx/gridworkx/services/internal/identity"
)

type Repository struct{ db *sql.DB }

func Open(ctx context.Context, dsn string) (*Repository, error) {
	if strings.TrimSpace(dsn) == "" {
		return nil, errors.New("GRIDWORKS_DATABASE_URL is required")
	}
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, err
	}
	return &Repository{db: db}, nil
}
func (r *Repository) Close() error                   { return r.db.Close() }
func (r *Repository) Ping(ctx context.Context) error { return r.db.PingContext(ctx) }
func opaqueID(prefix string) (string, error) {
	b := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return "", err
	}
	return prefix + "." + hex.EncodeToString(b), nil
}
func token() (string, error) {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
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
	return err
}

func (r *Repository) IssueGuest(ctx context.Context, key string, now time.Time) (api.GuestResult, error) {
	d := digest("identity.guest_issue\x00" + key)
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return api.GuestResult{}, err
	}
	defer tx.Rollback()
	var typ, storedDigest, resultRef string
	err = tx.QueryRowContext(ctx, `SELECT mutation_type,request_digest,result_ref FROM gridworks.identity_mutation_receipts WHERE idempotency_key=$1 FOR UPDATE`, key).Scan(&typ, &storedDigest, &resultRef)
	if err == nil {
		if typ != "identity.guest_issue" || storedDigest != d {
			return api.GuestResult{}, api.ErrConflict
		}
		parts := strings.Split(resultRef, "\x00")
		if len(parts) != 3 {
			return api.GuestResult{}, errors.New("invalid guest receipt")
		}
		var exp time.Time
		if err := tx.QueryRowContext(ctx, `SELECT expires_at FROM gridworks.guest_sessions WHERE session_id=$1`, parts[2]).Scan(&exp); err != nil {
			return api.GuestResult{}, mapDB(err)
		}
		if err := tx.Commit(); err != nil {
			return api.GuestResult{}, err
		}
		return api.GuestResult{AccountID: parts[0], PlayerID: parts[1], SessionID: parts[2], ExpiresAt: exp, Replay: true}, nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return api.GuestResult{}, err
	}
	accountID, err := opaqueID("account")
	if err != nil {
		return api.GuestResult{}, err
	}
	playerID, err := opaqueID("player")
	if err != nil {
		return api.GuestResult{}, err
	}
	sessionID, err := opaqueID("session")
	if err != nil {
		return api.GuestResult{}, err
	}
	raw, err := token()
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
	resultRef = strings.Join([]string{accountID, playerID, sessionID}, "\x00")
	if _, err = tx.ExecContext(ctx, `INSERT INTO gridworks.identity_mutation_receipts(idempotency_key,mutation_type,request_digest,result_ref) VALUES($1,'identity.guest_issue',$2,$3)`, key, d, resultRef); err != nil {
		return api.GuestResult{}, err
	}
	if err = tx.Commit(); err != nil {
		return api.GuestResult{}, err
	}
	return api.GuestResult{AccountID: accountID, PlayerID: playerID, SessionID: sessionID, ExpiresAt: exp, RawSessionToken: raw}, nil
}
func (r *Repository) Authenticate(ctx context.Context, raw string, now time.Time) (api.Principal, error) {
	var p api.Principal
	err := r.db.QueryRowContext(ctx, `SELECT a.account_id,p.player_id,s.session_id,a.status FROM gridworks.guest_sessions s JOIN gridworks.accounts a ON a.account_id=s.account_id JOIN gridworks.players p ON p.account_id=a.account_id WHERE s.token_digest=$1 AND s.revoked_at IS NULL AND s.expires_at>$2`, digest(raw), now).Scan(&p.AccountID, &p.PlayerID, &p.SessionID, &p.Status)
	if err != nil {
		return api.Principal{}, mapDB(err)
	}
	return p, nil
}
func (r *Repository) RevokeSession(ctx context.Context, p api.Principal, now time.Time) error {
	_, err := r.db.ExecContext(ctx, `UPDATE gridworks.guest_sessions SET revoked_at=$1 WHERE session_id=$2 AND revoked_at IS NULL`, now, p.SessionID)
	return err
}
func (r *Repository) GetMe(ctx context.Context, p api.Principal) (api.MeView, error) {
	var v api.MeView
	var payload []byte
	err := r.db.QueryRowContext(ctx, `SELECT a.account_id,p.player_id,a.status,json_build_object('player_id',pr.player_id,'handle',pr.handle_display,'display_name',pr.display_name,'avatar_asset_id',pr.avatar_asset_id,'bio',pr.bio,'locale',pr.locale,'timezone',pr.timezone,'visibility',pr.visibility,'dm_policy',pr.dm_policy,'discoverable',pr.discoverable,'notification_preferences',pr.notification_preferences) FROM gridworks.accounts a JOIN gridworks.players p ON p.account_id=a.account_id JOIN gridworks.player_profiles pr ON pr.player_id=p.player_id WHERE a.account_id=$1 AND p.player_id=$2`, p.AccountID, p.PlayerID).Scan(&v.AccountID, &v.PlayerID, &v.Status, &payload)
	if err != nil {
		return api.MeView{}, mapDB(err)
	}
	v.Profile = json.RawMessage(payload)
	return v, nil
}
func (r *Repository) UpdateProfile(ctx context.Context, p api.Principal, key string, patch api.ProfilePatch, now time.Time) (api.MeView, error) {
	payload, _ := json.Marshal(patch)
	d := digest("profile.update\x00" + p.PlayerID + "\x00" + string(payload))
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return api.MeView{}, err
	}
	defer tx.Rollback()
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
	_, err = tx.ExecContext(ctx, `UPDATE gridworks.player_profiles SET display_name=COALESCE($1,display_name),bio=COALESCE($2,bio),locale=COALESCE($3,locale),timezone=COALESCE($4,timezone),visibility=COALESCE($5,visibility),dm_policy=COALESCE($6,dm_policy),discoverable=COALESCE($7,discoverable),notification_preferences=COALESCE($8,notification_preferences),updated_at=$9 WHERE player_id=$10 AND account_id=$11`, patch.DisplayName, patch.Bio, patch.Locale, patch.Timezone, patch.Visibility, patch.DMPolicy, patch.Discoverable, prefs, now, p.PlayerID, p.AccountID)
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
	var v api.PublicPlayerView
	err := r.db.QueryRowContext(ctx, `SELECT player_id,handle_display,display_name,avatar_asset_id,bio,locale,visibility,discoverable FROM gridworks.player_profiles WHERE player_id=$1 AND visibility='public' AND discoverable=true`, id).Scan(&v.PlayerID, &v.Handle, &v.DisplayName, &v.AvatarAssetID, &v.Bio, &v.Locale, &v.Visibility, &v.Discoverable)
	if err != nil {
		return api.PublicPlayerView{}, mapDB(err)
	}
	return v, nil
}
func (r *Repository) CreateCompany(ctx context.Context, p api.Principal, key string, req api.CompanyCreateRequest, now time.Time) (api.CreatedEntity, error) {
	if req.CompanyType == "system" {
		return api.CreatedEntity{}, api.ErrForbidden
	}
	displayName := strings.TrimSpace(req.Name)
	name, err := identity.NormalizeHandleWithMax(strings.Join(strings.Fields(displayName), "-"), 80)
	if err != nil {
		return api.CreatedEntity{}, err
	}
	id, err := opaqueID("company")
	if err != nil {
		return api.CreatedEntity{}, err
	}
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return api.CreatedEntity{}, err
	}
	defer tx.Rollback()
	if _, err = tx.ExecContext(ctx, `INSERT INTO gridworks.companies(company_id,company_type,name_display,name_canonical,name_skeleton,status,visibility,created_at,updated_at) VALUES($1,$2,$3,$4,$5,'active','public',$6,$6)`, id, req.CompanyType, displayName, name.Canonical, name.Skeleton, now); err != nil {
		return api.CreatedEntity{}, mapDB(err)
	}
	oid, _ := opaqueID("ownership")
	eid, _ := opaqueID("ownership-event")
	if _, err = tx.ExecContext(ctx, `INSERT INTO gridworks.company_ownership(ownership_id,entity_id,entity_type,owner_type,owner_id,share_bps,effective_from) VALUES($1,$2,'company','player',$3,10000,$4)`, oid, id, p.PlayerID, now); err != nil {
		return api.CreatedEntity{}, mapDB(err)
	}
	if _, err = tx.ExecContext(ctx, `INSERT INTO gridworks.ownership_history(event_id,entity_id,mutation_type,to_owner_type,to_owner_id,share_bps,effective_time,source_ref,actor_player_id) VALUES($1,$2,'company.create','player',$3,10000,$4,$5,$3)`, eid, id, p.PlayerID, now, key); err != nil {
		return api.CreatedEntity{}, err
	}
	d := digest("company.create\x00" + req.CompanyType + "\x00" + req.Name + "\x00" + p.PlayerID)
	if _, err = tx.ExecContext(ctx, `INSERT INTO gridworks.company_mutation_receipts(idempotency_key,mutation_type,request_digest,result_ref) VALUES($1,'company.create',$2,$3)`, key, d, id); err != nil {
		return api.CreatedEntity{}, err
	}
	if err = tx.Commit(); err != nil {
		return api.CreatedEntity{}, err
	}
	return api.CreatedEntity{ID: id, Name: req.Name}, nil
}
func (r *Repository) GetPublicCompany(ctx context.Context, id string) (api.PublicCompanyView, error) {
	var v api.PublicCompanyView
	err := r.db.QueryRowContext(ctx, `SELECT company_id,name_display,description,logo_asset_id,industry_id,headquarters_region_id,visibility FROM gridworks.companies WHERE company_id=$1 AND visibility='public'`, id).Scan(&v.CompanyID, &v.Name, &v.Description, &v.LogoAssetID, &v.IndustryID, &v.HeadquartersRegionID, &v.Visibility)
	if err != nil {
		return api.PublicCompanyView{}, mapDB(err)
	}
	return v, nil
}
func (r *Repository) CreateGroup(ctx context.Context, p api.Principal, key string, req api.GroupCreateRequest, now time.Time) (api.CreatedEntity, error) {
	displayName := strings.TrimSpace(req.Name)
	name, err := identity.NormalizeHandleWithMax(strings.Join(strings.Fields(displayName), "-"), 80)
	if err != nil {
		return api.CreatedEntity{}, err
	}
	id, err := opaqueID("group")
	if err != nil {
		return api.CreatedEntity{}, err
	}
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return api.CreatedEntity{}, err
	}
	defer tx.Rollback()
	if _, err = tx.ExecContext(ctx, `INSERT INTO gridworks.company_groups(group_id,name_display,name_canonical,name_skeleton,status,created_at,updated_at) VALUES($1,$2,$3,$4,'active',$5,$5)`, id, displayName, name.Canonical, name.Skeleton, now); err != nil {
		return api.CreatedEntity{}, mapDB(err)
	}
	oid, _ := opaqueID("ownership")
	eid, _ := opaqueID("ownership-event")
	if _, err = tx.ExecContext(ctx, `INSERT INTO gridworks.company_ownership(ownership_id,entity_id,entity_type,owner_type,owner_id,share_bps,effective_from) VALUES($1,$2,'group','player',$3,10000,$4)`, oid, id, p.PlayerID, now); err != nil {
		return api.CreatedEntity{}, mapDB(err)
	}
	if _, err = tx.ExecContext(ctx, `INSERT INTO gridworks.ownership_history(event_id,entity_id,mutation_type,to_owner_type,to_owner_id,share_bps,effective_time,source_ref,actor_player_id) VALUES($1,$2,'group.create','player',$3,10000,$4,$5,$3)`, eid, id, p.PlayerID, now, key); err != nil {
		return api.CreatedEntity{}, err
	}
	d := digest("group.create\x00" + req.Name + "\x00" + p.PlayerID)
	if _, err = tx.ExecContext(ctx, `INSERT INTO gridworks.company_mutation_receipts(idempotency_key,mutation_type,request_digest,result_ref) VALUES($1,'group.create',$2,$3)`, key, d, id); err != nil {
		return api.CreatedEntity{}, err
	}
	if err = tx.Commit(); err != nil {
		return api.CreatedEntity{}, err
	}
	return api.CreatedEntity{ID: id, Name: req.Name}, nil
}
