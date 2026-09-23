package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

type scriptedRepo struct {
	pingFn          func(context.Context) error
	issueGuestFn    func(context.Context, string, time.Time) (GuestResult, error)
	authenticateFn  func(context.Context, string, time.Time) (Principal, error)
	revokeFn        func(context.Context, string, time.Time) error
	getMeFn         func(context.Context, Principal) (MeView, error)
	updateProfileFn func(context.Context, Principal, string, ProfilePatch, time.Time) (MeView, error)
	publicPlayerFn  func(context.Context, string) (PublicPlayerView, error)
	createCompanyFn func(context.Context, Principal, string, CompanyCreateRequest, time.Time) (CreatedEntity, error)
	publicCompanyFn func(context.Context, string) (PublicCompanyView, error)
	createGroupFn   func(context.Context, Principal, string, GroupCreateRequest, time.Time) (CreatedEntity, error)
}

func (s *scriptedRepo) Ping(ctx context.Context) error {
	if s.pingFn != nil {
		return s.pingFn(ctx)
	}
	return nil
}

func (s *scriptedRepo) IssueGuest(ctx context.Context, key string, now time.Time) (GuestResult, error) {
	if s.issueGuestFn != nil {
		return s.issueGuestFn(ctx, key, now)
	}
	return GuestResult{AccountID: "account.one", PlayerID: "player.one", SessionID: "session.one", ExpiresAt: now.Add(time.Hour), Profile: json.RawMessage(`{"display_name":"Guest"}`), RawSessionToken: "raw-once"}, nil
}

func (s *scriptedRepo) Authenticate(ctx context.Context, token string, now time.Time) (Principal, error) {
	if s.authenticateFn != nil {
		return s.authenticateFn(ctx, token, now)
	}
	if token != "valid-token" {
		return Principal{}, ErrUnauthorized
	}
	return Principal{AccountID: "account.one", PlayerID: "player.one", SessionID: "session.one", Status: "guest"}, nil
}

func (s *scriptedRepo) RevokeSession(context.Context, Principal, time.Time) error { return nil }

func (s *scriptedRepo) RevokePresentedSession(ctx context.Context, token string, now time.Time) error {
	if s.revokeFn != nil {
		return s.revokeFn(ctx, token, now)
	}
	return nil
}

func (s *scriptedRepo) GetMe(ctx context.Context, p Principal) (MeView, error) {
	if s.getMeFn != nil {
		return s.getMeFn(ctx, p)
	}
	return MeView{AccountID: p.AccountID, PlayerID: p.PlayerID, Status: p.Status, Profile: json.RawMessage(`{"display_name":"Guest"}`)}, nil
}

func (s *scriptedRepo) UpdateProfile(ctx context.Context, p Principal, key string, patch ProfilePatch, now time.Time) (MeView, error) {
	if s.updateProfileFn != nil {
		return s.updateProfileFn(ctx, p, key, patch, now)
	}
	return MeView{AccountID: p.AccountID, PlayerID: p.PlayerID, Status: p.Status, Profile: json.RawMessage(`{"display_name":"Updated"}`)}, nil
}

func (s *scriptedRepo) GetPublicPlayer(ctx context.Context, id string) (PublicPlayerView, error) {
	if s.publicPlayerFn != nil {
		return s.publicPlayerFn(ctx, id)
	}
	return PublicPlayerView{PlayerID: id, Handle: "public-player", DisplayName: "Public", Locale: "en", Visibility: "public", Discoverable: true}, nil
}

func (s *scriptedRepo) CreateCompany(ctx context.Context, p Principal, key string, req CompanyCreateRequest, now time.Time) (CreatedEntity, error) {
	if s.createCompanyFn != nil {
		return s.createCompanyFn(ctx, p, key, req, now)
	}
	return CreatedEntity{ID: "company.one", Name: strings.TrimSpace(req.Name)}, nil
}

func (s *scriptedRepo) GetPublicCompany(ctx context.Context, id string) (PublicCompanyView, error) {
	if s.publicCompanyFn != nil {
		return s.publicCompanyFn(ctx, id)
	}
	return PublicCompanyView{CompanyID: id, Name: "Company", Visibility: "public"}, nil
}

func (s *scriptedRepo) CreateGroup(ctx context.Context, p Principal, key string, req GroupCreateRequest, now time.Time) (CreatedEntity, error) {
	if s.createGroupFn != nil {
		return s.createGroupFn(ctx, p, key, req, now)
	}
	return CreatedEntity{ID: "group.one", Name: strings.TrimSpace(req.Name)}, nil
}

func perform(h http.Handler, method, path, body string, headers map[string]string) *httptest.ResponseRecorder {
	r := httptest.NewRequest(method, path, strings.NewReader(body))
	for k, v := range headers {
		r.Header.Set(k, v)
	}
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	return w
}

func authenticatedHeaders() map[string]string {
	return map[string]string{
		"Authorization":   "Bearer valid-token",
		"Content-Type":    "application/json",
		"Idempotency-Key": "test-key",
	}
}

func TestHTTPBoundaryRequiredCases(t *testing.T) {
	repo := &scriptedRepo{}
	h := NewServer(repo, "test").Mux()

	if w := perform(h, http.MethodPost, "/healthz", "", nil); w.Code != http.StatusMethodNotAllowed || w.Header().Get("Allow") != "GET" {
		t.Fatalf("static wrong method=%d allow=%q", w.Code, w.Header().Get("Allow"))
	}
	if w := perform(h, http.MethodPost, "/api/v1/players/player.one", "", nil); w.Code != http.StatusMethodNotAllowed || w.Header().Get("Allow") != "GET" {
		t.Fatalf("player wrong method=%d allow=%q", w.Code, w.Header().Get("Allow"))
	}
	if w := perform(h, http.MethodPost, "/api/v1/companies/company.one", "", nil); w.Code != http.StatusMethodNotAllowed || w.Header().Get("Allow") != "GET" {
		t.Fatalf("company wrong method=%d allow=%q", w.Code, w.Header().Get("Allow"))
	}
	if w := perform(h, http.MethodPost, "/api/v1/players/player.one/extra", "", nil); w.Code != http.StatusNotFound {
		t.Fatalf("extra parameterized segment=%d", w.Code)
	}
	if w := perform(h, http.MethodGet, "/unknown", "", nil); w.Code != http.StatusNotFound || !strings.Contains(w.Body.String(), "route.not_found") {
		t.Fatalf("unknown route=%d body=%s", w.Code, w.Body.String())
	}

	headers := authenticatedHeaders()
	headers["Content-Type"] = "application/jsonfoo"
	if w := perform(h, http.MethodPost, "/api/v1/companies", `{"company_type":"operating","name":"Acme"}`, headers); w.Code != http.StatusUnsupportedMediaType {
		t.Fatalf("invalid media type=%d", w.Code)
	}

	headers = authenticatedHeaders()
	if w := perform(h, http.MethodPost, "/api/v1/companies", "{", headers); w.Code != http.StatusBadRequest {
		t.Fatalf("malformed json=%d", w.Code)
	}
	if w := perform(h, http.MethodPost, "/api/v1/companies", `{"company_type":"operating","name":"Acme"} {"x":1}`, headers); w.Code != http.StatusBadRequest {
		t.Fatalf("trailing json=%d", w.Code)
	}
	tooLarge := strings.Repeat("x", MaxJSONBody+1)
	if w := perform(h, http.MethodPost, "/api/v1/companies", tooLarge, headers); w.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("oversized body=%d", w.Code)
	}

	if w := perform(h, http.MethodGet, "/healthz", "", nil); w.Header().Get("Access-Control-Allow-Origin") != "" {
		t.Fatalf("wildcard/unexpected CORS header=%q", w.Header().Get("Access-Control-Allow-Origin"))
	}
	if w := perform(h, http.MethodGet, "/api/v1/me", "", map[string]string{"Authorization": "Bearer valid-token"}); w.Code != http.StatusOK || w.Header().Get("Cache-Control") != "no-store" {
		t.Fatalf("auth cache contract=%d cache=%q", w.Code, w.Header().Get("Cache-Control"))
	}
}

func TestRequestIDGenerationPreservationAndPanicRecovery(t *testing.T) {
	repo := &scriptedRepo{}
	h := NewServer(repo, "test").Mux()

	w := perform(h, http.MethodGet, "/healthz", "", nil)
	if got := w.Header().Get("X-Request-ID"); got == "" {
		t.Fatal("missing generated request id")
	}

	w = perform(h, http.MethodGet, "/healthz", "", map[string]string{"X-Request-ID": "safe.ID_123:part-value"})
	if got := w.Header().Get("X-Request-ID"); got != "safe.ID_123:part-value" {
		t.Fatalf("safe request id not preserved: %q", got)
	}

	repo.getMeFn = func(context.Context, Principal) (MeView, error) {
		panic("test panic")
	}
	w = perform(h, http.MethodGet, "/api/v1/me", "", map[string]string{"Authorization": "Bearer valid-token"})
	if w.Code != http.StatusInternalServerError || strings.Contains(w.Body.String(), "test panic") {
		t.Fatalf("panic response=%d body=%s", w.Code, w.Body.String())
	}
}

func TestAuthenticationAndReplaySafeRevokeMatrix(t *testing.T) {
	status := "guest"
	revoked := false
	repo := &scriptedRepo{}
	repo.authenticateFn = func(_ context.Context, token string, _ time.Time) (Principal, error) {
		if token != "valid-token" || revoked {
			return Principal{}, ErrUnauthorized
		}
		return Principal{AccountID: "account.one", PlayerID: "player.one", SessionID: "session.one", Status: status}, nil
	}
	repo.revokeFn = func(_ context.Context, token string, _ time.Time) error {
		if token != "valid-token" {
			return ErrUnauthorized
		}
		revoked = true
		return nil
	}
	h := NewServer(repo, "test").Mux()

	if w := perform(h, http.MethodGet, "/api/v1/me", "", map[string]string{"Authorization": "Bearer valid-token"}); w.Code != http.StatusOK {
		t.Fatalf("valid auth=%d", w.Code)
	}
	for _, token := range []string{"wrong-token", "expired-token", "revoked-token"} {
		if w := perform(h, http.MethodGet, "/api/v1/me", "", map[string]string{"Authorization": "Bearer " + token}); w.Code != http.StatusUnauthorized {
			t.Fatalf("%s auth=%d", token, w.Code)
		}
	}
	for _, restricted := range []string{"suspended", "recovery_restricted"} {
		status = restricted
		if w := perform(h, http.MethodGet, "/api/v1/me", "", map[string]string{"Authorization": "Bearer valid-token"}); w.Code != http.StatusUnauthorized {
			t.Fatalf("%s auth=%d", restricted, w.Code)
		}
	}
	status = "guest"

	for i := 0; i < 2; i++ {
		if w := perform(h, http.MethodDelete, "/api/v1/auth/session", "", map[string]string{"Authorization": "Bearer valid-token"}); w.Code != http.StatusNoContent {
			t.Fatalf("revoke[%d]=%d", i, w.Code)
		}
	}
	if w := perform(h, http.MethodGet, "/api/v1/me", "", map[string]string{"Authorization": "Bearer valid-token"}); w.Code != http.StatusUnauthorized {
		t.Fatalf("revoked session still authenticates=%d", w.Code)
	}
}

func TestProfileContractAndIdempotency(t *testing.T) {
	fixed := time.Date(2026, 9, 23, 10, 0, 0, 0, time.UTC)
	receipts := map[string]string{}
	var updateTime time.Time
	repo := &scriptedRepo{}
	repo.publicPlayerFn = func(_ context.Context, id string) (PublicPlayerView, error) {
		if id == "private" {
			return PublicPlayerView{}, ErrNotFound
		}
		return PublicPlayerView{PlayerID: id, Handle: "visible", DisplayName: "Visible", Locale: "en", Visibility: "public", Discoverable: true}, nil
	}
	repo.updateProfileFn = func(_ context.Context, p Principal, key string, patch ProfilePatch, now time.Time) (MeView, error) {
		encoded, _ := json.Marshal(patch)
		if prior, ok := receipts[key]; ok && prior != string(encoded) {
			return MeView{}, ErrConflict
		}
		receipts[key] = string(encoded)
		updateTime = now
		return MeView{AccountID: p.AccountID, PlayerID: p.PlayerID, Status: p.Status, Profile: json.RawMessage(`{"display_name":"Updated"}`)}, nil
	}
	s := NewServer(repo, "test")
	s.Now = func() time.Time { return fixed }
	h := s.Mux()

	if w := perform(h, http.MethodGet, "/api/v1/me", "", map[string]string{"Authorization": "Bearer valid-token"}); w.Code != http.StatusOK {
		t.Fatalf("self read=%d", w.Code)
	}
	w := perform(h, http.MethodGet, "/api/v1/players/player.one", "", nil)
	if w.Code != http.StatusOK || strings.Contains(w.Body.String(), "account_id") || strings.Contains(w.Body.String(), "session") || strings.Contains(w.Body.String(), "provider") {
		t.Fatalf("public DTO leakage/status=%d body=%s", w.Code, w.Body.String())
	}
	if w := perform(h, http.MethodGet, "/api/v1/players/private", "", nil); w.Code != http.StatusNotFound {
		t.Fatalf("private player visibility=%d", w.Code)
	}

	headers := authenticatedHeaders()
	headers["Idempotency-Key"] = "profile.same"
	valid := `{"display_name":"Updated","bio":"safe","locale":"en","timezone":"UTC","visibility":"public","dm_policy":"everyone","discoverable":true,"notification_preferences":{"email":true}}`
	if w := perform(h, http.MethodPatch, "/api/v1/me/profile", valid, headers); w.Code != http.StatusOK {
		t.Fatalf("valid update=%d body=%s", w.Code, w.Body.String())
	}
	if !updateTime.Equal(fixed) {
		t.Fatalf("update time=%s want=%s", updateTime, fixed)
	}
	if w := perform(h, http.MethodPatch, "/api/v1/me/profile", valid, headers); w.Code != http.StatusOK {
		t.Fatalf("same-key replay=%d", w.Code)
	}
	if w := perform(h, http.MethodPatch, "/api/v1/me/profile", `{"display_name":"Different"}`, headers); w.Code != http.StatusConflict {
		t.Fatalf("contradictory idempotency=%d", w.Code)
	}
	if w := perform(h, http.MethodPatch, "/api/v1/me/profile", `{"handle":"forbidden"}`, headers); w.Code != http.StatusBadRequest {
		t.Fatalf("private/unknown field=%d", w.Code)
	}

	cases := []string{
		`{}`,
		`{"display_name":""}`,
		`{"bio":"bad\u0001"}`,
		`{"locale":"no_such_locale"}`,
		`{"timezone":"Not/AZone"}`,
		`{"visibility":"hidden"}`,
		`{"dm_policy":"maybe"}`,
	}
	for _, body := range cases {
		w := perform(h, http.MethodPatch, "/api/v1/me/profile", body, authenticatedHeaders())
		if w.Code != http.StatusBadRequest && w.Code != http.StatusUnprocessableEntity {
			t.Fatalf("invalid profile body=%s status=%d", body, w.Code)
		}
	}
}

func TestCompanyAndGroupHTTPContract(t *testing.T) {
	var owner Principal
	repo := &scriptedRepo{}
	repo.createCompanyFn = func(_ context.Context, p Principal, _ string, req CompanyCreateRequest, _ time.Time) (CreatedEntity, error) {
		owner = p
		switch req.Name {
		case "Taken", "TAKEN", "Раypal", "GRIDWORKS":
			return CreatedEntity{}, ErrConflict
		}
		return CreatedEntity{ID: "company.one", Name: strings.TrimSpace(req.Name)}, nil
	}
	repo.createGroupFn = func(_ context.Context, p Principal, _ string, req GroupCreateRequest, _ time.Time) (CreatedEntity, error) {
		owner = p
		if req.Name == "Taken" {
			return CreatedEntity{}, ErrConflict
		}
		return CreatedEntity{ID: "group.one", Name: strings.TrimSpace(req.Name)}, nil
	}
	repo.publicCompanyFn = func(_ context.Context, id string) (PublicCompanyView, error) {
		if id == "private" {
			return PublicCompanyView{}, ErrNotFound
		}
		return PublicCompanyView{CompanyID: id, Name: "PublicCo", Visibility: "public"}, nil
	}
	h := NewServer(repo, "test").Mux()

	for _, companyType := range []string{"operating", "holding"} {
		body := `{"company_type":"` + companyType + `","name":"Acme ` + companyType + `"}`
		if w := perform(h, http.MethodPost, "/api/v1/companies", body, authenticatedHeaders()); w.Code != http.StatusCreated || owner.PlayerID != "player.one" {
			t.Fatalf("%s create=%d owner=%+v", companyType, w.Code, owner)
		}
	}
	if w := perform(h, http.MethodPost, "/api/v1/companies", `{"company_type":"system","name":"System"}`, authenticatedHeaders()); w.Code != http.StatusForbidden {
		t.Fatalf("system company=%d", w.Code)
	}
	if w := perform(h, http.MethodPost, "/api/v1/companies", `{"company_type":"operating","name":"Acme","owner_id":"player.other"}`, authenticatedHeaders()); w.Code != http.StatusBadRequest {
		t.Fatalf("alternate owner field=%d", w.Code)
	}
	for _, name := range []string{"Taken", "TAKEN", "Раypal", "GRIDWORKS"} {
		body := `{"company_type":"operating","name":"` + name + `"}`
		if w := perform(h, http.MethodPost, "/api/v1/companies", body, authenticatedHeaders()); w.Code != http.StatusConflict {
			t.Fatalf("name conflict %q=%d", name, w.Code)
		}
	}
	if w := perform(h, http.MethodPost, "/api/v1/company-groups", `{"name":"Holding One"}`, authenticatedHeaders()); w.Code != http.StatusCreated || owner.PlayerID != "player.one" {
		t.Fatalf("group create=%d owner=%+v", w.Code, owner)
	}
	if w := perform(h, http.MethodPost, "/api/v1/company-groups", `{"name":"Taken"}`, authenticatedHeaders()); w.Code != http.StatusConflict {
		t.Fatalf("shared namespace group conflict=%d", w.Code)
	}
	if w := perform(h, http.MethodGet, "/api/v1/companies/private", "", nil); w.Code != http.StatusNotFound {
		t.Fatalf("private company=%d", w.Code)
	}
	if w := perform(h, http.MethodGet, "/api/v1/companies/company.one", "", nil); w.Code != http.StatusOK {
		t.Fatalf("public company=%d", w.Code)
	}
	if w := perform(h, http.MethodPost, "/api/v1/ownership", `{}`, authenticatedHeaders()); w.Code != http.StatusNotFound {
		t.Fatalf("raw ownership route exists=%d", w.Code)
	}
	if w := perform(h, http.MethodPost, "/api/v1/identity/link", `{}`, authenticatedHeaders()); w.Code != http.StatusNotFound {
		t.Fatalf("provider link route exists=%d", w.Code)
	}
}

func TestErrorEnvelopeNeverUsesRawInternalError(t *testing.T) {
	repo := &scriptedRepo{}
	repo.getMeFn = func(context.Context, Principal) (MeView, error) {
		return MeView{}, errors.New("postgres secret detail")
	}
	w := perform(NewServer(repo, "test").Mux(), http.MethodGet, "/api/v1/me", "", map[string]string{"Authorization": "Bearer valid-token"})
	if w.Code != http.StatusInternalServerError || strings.Contains(w.Body.String(), "postgres secret detail") {
		t.Fatalf("unsafe error response=%d body=%s", w.Code, w.Body.String())
	}
}
