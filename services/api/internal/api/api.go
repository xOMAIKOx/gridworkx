package api

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"log"
	"mime"
	"net/http"
	"strings"
	"time"
	"unicode"

	"github.com/xOMAIKOx/gridworkx/services/internal/company"
	"github.com/xOMAIKOx/gridworkx/services/internal/identity"
	"github.com/xOMAIKOx/gridworkx/services/internal/progression"
)

const MaxJSONBody = 1 << 20

type Principal struct {
	AccountID string
	PlayerID  string
	SessionID string
	Status    string
}

type GuestResult struct {
	AccountID       string
	PlayerID        string
	SessionID       string
	ExpiresAt       time.Time
	Profile         json.RawMessage
	RawSessionToken string
	Replay          bool
}

type MeView struct {
	AccountID string          `json:"account_id"`
	PlayerID  string          `json:"player_id"`
	Status    string          `json:"status"`
	Profile   json.RawMessage `json:"profile"`
}

type PublicPlayerView struct {
	PlayerID      string `json:"player_id"`
	Handle        string `json:"handle"`
	DisplayName   string `json:"display_name"`
	AvatarAssetID string `json:"avatar_asset_id,omitempty"`
	Bio           string `json:"bio,omitempty"`
	Locale        string `json:"locale"`
	Visibility    string `json:"visibility"`
	Discoverable  bool   `json:"discoverable"`
}

type PublicCompanyView struct {
	CompanyID            string `json:"company_id"`
	Name                 string `json:"name"`
	Description          string `json:"description,omitempty"`
	LogoAssetID          string `json:"logo_asset_id,omitempty"`
	IndustryID           string `json:"industry_id,omitempty"`
	HeadquartersRegionID string `json:"headquarters_region_id,omitempty"`
	Visibility           string `json:"visibility"`
}

type PublicGroupView struct {
	GroupID    string `json:"group_id"`
	Name       string `json:"name"`
	Visibility string `json:"visibility"`
}

type ProfilePatch struct {
	DisplayName             *string         `json:"display_name,omitempty"`
	Bio                     *string         `json:"bio,omitempty"`
	Locale                  *string         `json:"locale,omitempty"`
	Timezone                *string         `json:"timezone,omitempty"`
	Visibility              *string         `json:"visibility,omitempty"`
	DMPolicy                *string         `json:"dm_policy,omitempty"`
	Discoverable            *bool           `json:"discoverable,omitempty"`
	NotificationPreferences map[string]bool `json:"notification_preferences,omitempty"`
}

type CompanyCreateRequest struct {
	CompanyType string `json:"company_type"`
	Name        string `json:"name"`
}

type GroupCreateRequest struct {
	Name string `json:"name"`
}

type CreatedEntity struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type errorEnvelope struct {
	Error errorBody `json:"error"`
}

type errorBody struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	RequestID string `json:"request_id"`
	Timestamp string `json:"timestamp"`
}

type RevocationRepository interface {
	RevokePresentedSession(context.Context, string, time.Time) error
}

type IdentityRepository interface {
	IssueGuest(context.Context, string, time.Time) (GuestResult, error)
	Authenticate(context.Context, string, time.Time) (Principal, error)
	RevokeSession(context.Context, Principal, time.Time) error
	GetMe(context.Context, Principal) (MeView, error)
	UpdateProfile(context.Context, Principal, string, ProfilePatch, time.Time) (MeView, error)
	GetPublicPlayer(context.Context, string) (PublicPlayerView, error)
}

type CompanyRepository interface {
	CreateCompany(context.Context, Principal, string, CompanyCreateRequest, time.Time) (CreatedEntity, error)
	GetPublicCompany(context.Context, string) (PublicCompanyView, error)
	CreateGroup(context.Context, Principal, string, GroupCreateRequest, time.Time) (CreatedEntity, error)
}

type Repository interface {
	IdentityRepository
	CompanyRepository
	Ping(context.Context) error
}

type ProgressionReadRepository interface {
	GetPlayerSkills(context.Context, Principal) ([]progression.PlayerSkillView, error)
	GetCompanyManagers(context.Context, Principal, string) ([]progression.ManagerView, error)
	GetCompanyManager(context.Context, Principal, string, string) (progression.ManagerView, error)
}

type APIError struct {
	Status  int
	Code    string
	Message string
}

func (e *APIError) Error() string { return e.Code }

var (
	ErrUnauthorized = &APIError{Status: http.StatusUnauthorized, Code: "auth.invalid", Message: "authentication required"}
	ErrForbidden    = &APIError{Status: http.StatusForbidden, Code: "auth.forbidden", Message: "operation is not permitted"}
	ErrNotFound     = &APIError{Status: http.StatusNotFound, Code: "resource.not_found", Message: "resource not found"}
	ErrConflict     = &APIError{Status: http.StatusConflict, Code: "mutation.conflict", Message: "mutation conflicts with existing state"}
)

type Server struct {
	Repo    Repository
	Version string
	Now     func() time.Time
	MaxBody int64
}

func NewServer(repo Repository, version string) *Server {
	return &Server{Repo: repo, Version: version, Now: time.Now, MaxBody: MaxJSONBody}
}

type routeSpec struct {
	Method  string
	Pattern string
	Handler func(*Server) http.HandlerFunc
}

var routeRegistry = []routeSpec{
	{Method: http.MethodGet, Pattern: "/healthz", Handler: func(s *Server) http.HandlerFunc { return s.health }},
	{Method: http.MethodGet, Pattern: "/readyz", Handler: func(s *Server) http.HandlerFunc { return s.ready }},
	{Method: http.MethodGet, Pattern: "/version", Handler: func(s *Server) http.HandlerFunc { return s.version }},
	{Method: http.MethodPost, Pattern: "/api/v1/auth/guest", Handler: func(s *Server) http.HandlerFunc { return s.guest }},
	{Method: http.MethodDelete, Pattern: "/api/v1/auth/session", Handler: func(s *Server) http.HandlerFunc { return s.revoke }},
	{Method: http.MethodGet, Pattern: "/api/v1/me", Handler: func(s *Server) http.HandlerFunc { return s.me }},
	{Method: http.MethodGet, Pattern: "/api/v1/me/skills", Handler: func(s *Server) http.HandlerFunc { return s.skills }},
	{Method: http.MethodPatch, Pattern: "/api/v1/me/profile", Handler: func(s *Server) http.HandlerFunc { return s.profile }},
	{Method: http.MethodGet, Pattern: "/api/v1/players/{player_id}", Handler: func(s *Server) http.HandlerFunc { return s.player }},
	{Method: http.MethodPost, Pattern: "/api/v1/companies", Handler: func(s *Server) http.HandlerFunc { return s.createCompany }},
	{Method: http.MethodGet, Pattern: "/api/v1/companies/{company_id}", Handler: func(s *Server) http.HandlerFunc { return s.company }},
	{Method: http.MethodGet, Pattern: "/api/v1/companies/{company_id}/managers", Handler: func(s *Server) http.HandlerFunc { return s.managers }},
	{Method: http.MethodGet, Pattern: "/api/v1/companies/{company_id}/managers/{manager_id}", Handler: func(s *Server) http.HandlerFunc { return s.manager }},
	{Method: http.MethodPost, Pattern: "/api/v1/company-groups", Handler: func(s *Server) http.HandlerFunc { return s.createGroup }},
}

func RouteInventory() map[string][]string {
	inventory := make(map[string][]string, len(routeRegistry))
	for _, route := range routeRegistry {
		inventory[route.Pattern] = append(inventory[route.Pattern], route.Method)
	}
	return inventory
}

func (s *Server) Mux() http.Handler {
	mux := http.NewServeMux()
	for _, route := range routeRegistry {
		mux.HandleFunc(route.Method+" "+route.Pattern, route.Handler(s))
	}
	mux.HandleFunc("/", s.notFound)
	return s.middleware(s.methodGuard(mux))
}

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusWriter) WriteHeader(status int) {
	if w.status == 0 {
		w.status = status
	}
	w.ResponseWriter.WriteHeader(status)
}

func (w *statusWriter) Write(p []byte) (int, error) {
	if w.status == 0 {
		w.status = 200
	}
	return w.ResponseWriter.Write(p)
}

func (s *Server) middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get("X-Request-ID")
		if !validRequestID(requestID) {
			requestID = newRequestID()
		}
		ctx := context.WithValue(r.Context(), requestIDKey{}, requestID)
		r = r.WithContext(ctx)
		sw := &statusWriter{ResponseWriter: w}
		sw.Header().Set("X-Request-ID", requestID)
		sw.Header().Set("X-Content-Type-Options", "nosniff")
		start := time.Now()
		defer func() {
			if recover() != nil {
				writeError(sw, r, &APIError{Status: 500, Code: "internal.panic", Message: "internal server error"})
			}
			log.Printf("api_request request_id=%s method=%s path=%s status=%d duration_ms=%d", requestID, r.Method, r.URL.Path, sw.status, time.Since(start).Milliseconds())
		}()
		next.ServeHTTP(sw, r)
	})
}

type requestIDKey struct{}

func requestID(ctx context.Context) string { v, _ := ctx.Value(requestIDKey{}).(string); return v }

func validRequestID(v string) bool {
	if len(v) < 1 || len(v) > 96 {
		return false
	}
	for _, r := range v {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			continue
		}
		if r != '.' && r != '_' && r != ':' && r != '-' {
			return false
		}
	}
	return true
}

func newRequestID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "request-unavailable"
	}
	return hex.EncodeToString(b)
}

func (s *Server) now() time.Time {
	if s.Now != nil {
		return s.Now().UTC()
	}
	return time.Now().UTC()
}

func routePathMatches(pattern, path string) bool {
	patternParts := strings.Split(strings.Trim(pattern, "/"), "/")
	pathParts := strings.Split(strings.Trim(path, "/"), "/")
	if len(patternParts) != len(pathParts) {
		return false
	}
	for i, part := range patternParts {
		if strings.HasPrefix(part, "{") && strings.HasSuffix(part, "}") {
			if pathParts[i] == "" {
				return false
			}
			continue
		}
		if part != pathParts[i] {
			return false
		}
	}
	return true
}

func (s *Server) methodGuard(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		allowed := make([]string, 0, 2)
		for _, route := range routeRegistry {
			if !routePathMatches(route.Pattern, r.URL.Path) {
				continue
			}
			if r.Method == route.Method {
				next.ServeHTTP(w, r)
				return
			}
			allowed = append(allowed, route.Method)
		}
		if len(allowed) > 0 {
			w.Header().Set("Allow", strings.Join(allowed, ", "))
			writeError(w, r, &APIError{Status: 405, Code: "route.method_not_allowed", Message: "method not allowed"})
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (s *Server) notFound(w http.ResponseWriter, r *http.Request) {
	writeError(w, r, &APIError{Status: http.StatusNotFound, Code: "route.not_found", Message: "route not found"})
}

func (s *Server) health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, 200, map[string]string{"status": "ok", "role": "gridworks-api"}, false)
}

func (s *Server) ready(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()
	if err := s.Repo.Ping(ctx); err != nil {
		writeError(w, r, &APIError{Status: 503, Code: "service.not_ready", Message: "persistence is unavailable"})
		return
	}
	writeJSON(w, 200, map[string]string{"status": "ready"}, false)
}

func (s *Server) version(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, 200, map[string]string{"api_version": "v1", "version": s.Version}, false)
}

func requireJSON(r *http.Request) error {
	media, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil || strings.ToLower(media) != "application/json" {
		return &APIError{Status: 415, Code: "request.unsupported_media_type", Message: "application/json is required"}
	}
	return nil
}

func decodeJSON(w http.ResponseWriter, r *http.Request, dst any, max int64) error {
	if err := requireJSON(r); err != nil {
		return err
	}
	body := io.LimitReader(r.Body, max+1)
	data, err := io.ReadAll(body)
	if err != nil {
		return &APIError{Status: 400, Code: "request.read_failed", Message: "request body could not be read"}
	}
	if int64(len(data)) > max {
		return &APIError{Status: 413, Code: "request.body_too_large", Message: "request body is too large"}
	}
	dec := json.NewDecoder(strings.NewReader(string(data)))
	dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil {
		return &APIError{Status: 400, Code: "request.invalid_json", Message: "request body is invalid"}
	}
	var extra any
	if err := dec.Decode(&extra); err != io.EOF {
		return &APIError{Status: 400, Code: "request.multiple_values", Message: "request body must contain one JSON value"}
	}
	return nil
}

func writeJSON(w http.ResponseWriter, status int, v any, noStore bool) {
	w.Header().Set("Content-Type", "application/json")
	if noStore {
		w.Header().Set("Cache-Control", "no-store")
	}
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, r *http.Request, err error) {
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		apiErr = &APIError{Status: 500, Code: "internal.error", Message: "internal server error"}
	}
	writeJSON(w, apiErr.Status, errorEnvelope{Error: errorBody{Code: apiErr.Code, Message: apiErr.Message, RequestID: requestID(r.Context()), Timestamp: time.Now().UTC().Format(time.RFC3339Nano)}}, true)
}

func idempotency(r *http.Request) (string, error) {
	v := strings.TrimSpace(r.Header.Get("Idempotency-Key"))
	if len(v) < 1 || len(v) > 160 {
		return "", &APIError{Status: 400, Code: "mutation.invalid_idempotency", Message: "Idempotency-Key is required"}
	}
	for _, c := range v {
		if c < 0x21 || c > 0x7e {
			return "", &APIError{Status: 400, Code: "mutation.invalid_idempotency", Message: "Idempotency-Key is invalid"}
		}
	}
	return v, nil
}

func bearer(r *http.Request) (string, error) {
	parts := strings.Fields(r.Header.Get("Authorization"))
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return "", ErrUnauthorized
	}
	return parts[1], nil
}

func (s *Server) principal(r *http.Request) (Principal, string, error) {
	token, err := bearer(r)
	if err != nil {
		return Principal{}, "", err
	}
	p, err := s.Repo.Authenticate(r.Context(), token, s.now())
	if err != nil || (p.Status != "guest" && p.Status != "protected") {
		return Principal{}, "", ErrUnauthorized
	}
	return p, token, nil
}

func (s *Server) guest(w http.ResponseWriter, r *http.Request) {
	key, err := idempotency(r)
	if err != nil {
		writeError(w, r, err)
		return
	}
	if r.Body != nil && r.ContentLength != 0 {
		if err := decodeJSON(w, r, &struct{}{}, s.MaxBody); err != nil {
			writeError(w, r, err)
			return
		}
	}
	result, err := s.Repo.IssueGuest(r.Context(), key, s.now())
	if err != nil {
		writeError(w, r, mapError(err))
		return
	}
	response := map[string]any{"account_id": result.AccountID, "player_id": result.PlayerID, "session_id": result.SessionID, "expires_at": result.ExpiresAt, "profile": json.RawMessage(result.Profile), "replay": result.Replay}
	if result.Replay {
		response["credential_unavailable"] = true
	} else {
		response["session_token"] = result.RawSessionToken
	}
	writeJSON(w, 200, response, true)
}

func (s *Server) revoke(w http.ResponseWriter, r *http.Request) {
	token, err := bearer(r)
	if err != nil {
		writeError(w, r, err)
		return
	}
	if revoker, ok := s.Repo.(RevocationRepository); ok {
		err = revoker.RevokePresentedSession(r.Context(), token, s.now())
	} else {
		p, _, authErr := s.principal(r)
		if authErr != nil {
			writeError(w, r, authErr)
			return
		}
		err = s.Repo.RevokeSession(r.Context(), p, s.now())
	}
	if err != nil {
		writeError(w, r, mapError(err))
		return
	}
	_ = token
	writeJSON(w, 204, nil, true)
}

func (s *Server) me(w http.ResponseWriter, r *http.Request) {
	p, _, err := s.principal(r)
	if err != nil {
		writeError(w, r, err)
		return
	}
	v, err := s.Repo.GetMe(r.Context(), p)
	if err != nil {
		writeError(w, r, mapError(err))
		return
	}
	writeJSON(w, 200, v, true)
}

func (s *Server) profile(w http.ResponseWriter, r *http.Request) {
	p, _, err := s.principal(r)
	if err != nil {
		writeError(w, r, err)
		return
	}
	key, err := idempotency(r)
	if err != nil {
		writeError(w, r, err)
		return
	}
	var patch ProfilePatch
	if err := decodeJSON(w, r, &patch, s.MaxBody); err != nil {
		writeError(w, r, err)
		return
	}
	if err := validateProfilePatch(patch); err != nil {
		writeError(w, r, err)
		return
	}
	v, err := s.Repo.UpdateProfile(r.Context(), p, key, patch, s.now())
	if err != nil {
		writeError(w, r, mapError(err))
		return
	}
	writeJSON(w, 200, v, true)
}

func (s *Server) player(w http.ResponseWriter, r *http.Request) {
	v, err := s.Repo.GetPublicPlayer(r.Context(), r.PathValue("player_id"))
	if err != nil {
		writeError(w, r, mapError(err))
		return
	}
	writeJSON(w, 200, v, false)
}

func (s *Server) createCompany(w http.ResponseWriter, r *http.Request) {
	p, _, err := s.principal(r)
	if err != nil {
		writeError(w, r, err)
		return
	}
	key, err := idempotency(r)
	if err != nil {
		writeError(w, r, err)
		return
	}
	var req CompanyCreateRequest
	if err := decodeJSON(w, r, &req, s.MaxBody); err != nil {
		writeError(w, r, err)
		return
	}
	if req.CompanyType == "system" {
		writeError(w, r, ErrForbidden)
		return
	}
	v, err := s.Repo.CreateCompany(r.Context(), p, key, req, s.now())
	if err != nil {
		writeError(w, r, mapError(err))
		return
	}
	writeJSON(w, 201, v, true)
}

func (s *Server) skills(w http.ResponseWriter, r *http.Request) {
	p, _, err := s.principal(r)
	if err != nil {
		writeError(w, r, err)
		return
	}
	repo, ok := s.Repo.(ProgressionReadRepository)
	if !ok {
		writeError(w, r, &APIError{Status: 500, Code: "progression.unavailable", Message: "progression read service is unavailable"})
		return
	}
	views, err := repo.GetPlayerSkills(r.Context(), p)
	if err != nil {
		writeError(w, r, mapError(err))
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"skills": views}, true)
}

func (s *Server) managers(w http.ResponseWriter, r *http.Request) {
	p, _, err := s.principal(r)
	if err != nil {
		writeError(w, r, err)
		return
	}
	repo, ok := s.Repo.(ProgressionReadRepository)
	if !ok {
		writeError(w, r, &APIError{Status: 500, Code: "progression.unavailable", Message: "progression read service is unavailable"})
		return
	}
	views, err := repo.GetCompanyManagers(r.Context(), p, r.PathValue("company_id"))
	if err != nil {
		writeError(w, r, mapError(err))
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"managers": views}, true)
}

func (s *Server) manager(w http.ResponseWriter, r *http.Request) {
	p, _, err := s.principal(r)
	if err != nil {
		writeError(w, r, err)
		return
	}
	repo, ok := s.Repo.(ProgressionReadRepository)
	if !ok {
		writeError(w, r, &APIError{Status: 500, Code: "progression.unavailable", Message: "progression read service is unavailable"})
		return
	}
	view, err := repo.GetCompanyManager(r.Context(), p, r.PathValue("company_id"), r.PathValue("manager_id"))
	if err != nil {
		writeError(w, r, mapError(err))
		return
	}
	writeJSON(w, http.StatusOK, view, true)
}

func (s *Server) company(w http.ResponseWriter, r *http.Request) {
	v, err := s.Repo.GetPublicCompany(r.Context(), r.PathValue("company_id"))
	if err != nil {
		writeError(w, r, mapError(err))
		return
	}
	writeJSON(w, 200, v, false)
}

func (s *Server) createGroup(w http.ResponseWriter, r *http.Request) {
	p, _, err := s.principal(r)
	if err != nil {
		writeError(w, r, err)
		return
	}
	key, err := idempotency(r)
	if err != nil {
		writeError(w, r, err)
		return
	}
	var req GroupCreateRequest
	if err := decodeJSON(w, r, &req, s.MaxBody); err != nil {
		writeError(w, r, err)
		return
	}
	v, err := s.Repo.CreateGroup(r.Context(), p, key, req, s.now())
	if err != nil {
		writeError(w, r, mapError(err))
		return
	}
	writeJSON(w, 201, v, true)
}

func validateProfilePatch(p ProfilePatch) error {
	if p.DisplayName == nil && p.Bio == nil && p.Locale == nil && p.Timezone == nil && p.Visibility == nil && p.DMPolicy == nil && p.Discoverable == nil && p.NotificationPreferences == nil {
		return &APIError{Status: 400, Code: "profile.empty_patch", Message: "profile patch is empty"}
	}
	if p.DisplayName != nil && !identity.ValidateDisplayName(*p.DisplayName) {
		return &APIError{Status: 422, Code: "profile.invalid_display_name", Message: "display name is invalid"}
	}
	if p.Bio != nil {
		if len([]rune(*p.Bio)) > 1000 {
			return &APIError{Status: 422, Code: "profile.invalid_bio", Message: "bio is too long"}
		}
		for _, r := range *p.Bio {
			if unicode.IsControl(r) || unicode.Is(unicode.Cf, r) {
				return &APIError{Status: 422, Code: "profile.invalid_bio", Message: "bio contains forbidden characters"}
			}
		}
	}
	if p.Locale != nil && !identity.ValidateLocale(*p.Locale) {
		return &APIError{Status: 422, Code: "profile.invalid_locale", Message: "locale is invalid"}
	}
	if p.Timezone != nil && !identity.ValidateTimezone(*p.Timezone) {
		return &APIError{Status: 422, Code: "profile.invalid_timezone", Message: "timezone is invalid"}
	}
	if p.Visibility != nil && *p.Visibility != "public" && *p.Visibility != "private" {
		return &APIError{Status: 422, Code: "profile.invalid_visibility", Message: "visibility is invalid"}
	}
	if p.DMPolicy != nil && *p.DMPolicy != "everyone" && *p.DMPolicy != "nobody" {
		return &APIError{Status: 422, Code: "profile.invalid_dm_policy", Message: "DM policy is invalid"}
	}
	if p.NotificationPreferences != nil && len(p.NotificationPreferences) > 32 {
		return &APIError{Status: 422, Code: "profile.invalid_notifications", Message: "notification preferences are invalid"}
	}
	return nil
}

func mapError(err error) error {
	if err == nil {
		return nil
	}
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr
	}
	var companyErr company.CompanyError
	if errors.As(err, &companyErr) {
		switch companyErr {
		case company.ErrNameInvalid:
			return &APIError{Status: 422, Code: "company.invalid_name", Message: "company name is invalid"}
		case company.ErrNameReserved, company.ErrNameUnavailable, company.ErrIdempotencyConflict, company.ErrOwnershipConflict, company.ErrGroupConflict:
			return &APIError{Status: 409, Code: "mutation.conflict", Message: "mutation conflicts with existing state"}
		case company.ErrOwnerInvalid, company.ErrInvalidShare, company.ErrOwnershipTotal:
			return &APIError{Status: 422, Code: "mutation.invalid_state", Message: "mutation is invalid"}
		case company.ErrCompanyNotFound, company.ErrGroupNotFound:
			return ErrNotFound
		}
	}
	var identityErr identity.IdentityError
	if errors.As(err, &identityErr) {
		switch identityErr {
		case identity.ErrInvalidSession, identity.ErrSessionRevoked:
			return ErrUnauthorized
		case identity.ErrIdempotencyConflict, identity.ErrExternalIdentityConflict:
			return ErrConflict
		default:
			return &APIError{Status: 422, Code: "identity.invalid", Message: "identity mutation is invalid"}
		}
	}
	return err
}

func DigestToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}
