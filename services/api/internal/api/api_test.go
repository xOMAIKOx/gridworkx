package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

type fakeRepo struct {
	principal    Principal
	companyOwner Principal
	lastKey      string
	revoked      bool
}

func (f *fakeRepo) Ping(context.Context) error { return nil }
func (f *fakeRepo) IssueGuest(_ context.Context, _ string, _ time.Time) (GuestResult, error) {
	return GuestResult{AccountID: "account.one", PlayerID: "player.one", SessionID: "session.one", ExpiresAt: time.Now().Add(time.Hour), RawSessionToken: "raw-once"}, nil
}
func (f *fakeRepo) Authenticate(_ context.Context, token string, _ time.Time) (Principal, error) {
	if token != "valid-token" {
		return Principal{}, ErrUnauthorized
	}
	return f.principal, nil
}
func (f *fakeRepo) RevokeSession(_ context.Context, _ Principal, _ time.Time) error {
	f.revoked = true
	return nil
}
func (f *fakeRepo) RevokePresentedSession(context.Context, string, time.Time) error {
	f.revoked = true
	return nil
}
func (f *fakeRepo) GetMe(_ context.Context, _ Principal) (MeView, error) {
	return MeView{AccountID: "account.one", PlayerID: "player.one", Status: "guest", Profile: json.RawMessage(`{"display_name":"Guest"}`)}, nil
}
func (f *fakeRepo) UpdateProfile(_ context.Context, _ Principal, key string, _ ProfilePatch, _ time.Time) (MeView, error) {
	f.lastKey = key
	return MeView{AccountID: "account.one", PlayerID: "player.one", Status: "guest", Profile: json.RawMessage(`{"display_name":"Updated"}`)}, nil
}
func (f *fakeRepo) GetPublicPlayer(_ context.Context, _ string) (PublicPlayerView, error) {
	return PublicPlayerView{PlayerID: "player.one", Handle: "guest", DisplayName: "Guest", Locale: "en", Visibility: "public", Discoverable: true}, nil
}
func (f *fakeRepo) CreateCompany(_ context.Context, p Principal, key string, req CompanyCreateRequest, _ time.Time) (CreatedEntity, error) {
	f.companyOwner = p
	f.lastKey = key
	return CreatedEntity{ID: "company.one", Name: req.Name}, nil
}
func (f *fakeRepo) GetPublicCompany(_ context.Context, _ string) (PublicCompanyView, error) {
	return PublicCompanyView{CompanyID: "company.one", Name: "Company", Visibility: "public"}, nil
}
func (f *fakeRepo) CreateGroup(_ context.Context, p Principal, key string, req GroupCreateRequest, _ time.Time) (CreatedEntity, error) {
	f.companyOwner = p
	f.lastKey = key
	return CreatedEntity{ID: "group.one", Name: req.Name}, nil
}

func TestGuestRequiresIdempotencyAndReturnsNoStore(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/api/v1/auth/guest", strings.NewReader(`{}`))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	repo := &fakeRepo{}
	NewServer(repo, "test").Mux().ServeHTTP(w, r)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status=%d", w.Code)
	}
	r = httptest.NewRequest(http.MethodPost, "/api/v1/auth/guest", nil)
	r.Header.Set("Idempotency-Key", "guest.one")
	w = httptest.NewRecorder()
	NewServer(repo, "test").Mux().ServeHTTP(w, r)
	if w.Code != http.StatusOK || w.Header().Get("Cache-Control") != "no-store" {
		t.Fatalf("guest response status/cache=%d/%s", w.Code, w.Header().Get("Cache-Control"))
	}
}
func TestAuthenticatedCompanyDerivesOwnerAndRejectsUnknownPatch(t *testing.T) {
	repo := &fakeRepo{principal: Principal{AccountID: "account.one", PlayerID: "player.one", SessionID: "session.one", Status: "guest"}}
	h := NewServer(repo, "test").Mux()
	r := httptest.NewRequest(http.MethodPost, "/api/v1/companies", strings.NewReader(`{"company_type":"operating","name":"Company"}`))
	r.Header.Set("Authorization", "Bearer valid-token")
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Idempotency-Key", "company.one")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if w.Code != http.StatusCreated || repo.companyOwner.PlayerID != "player.one" {
		t.Fatalf("company status/owner=%d/%+v", w.Code, repo.companyOwner)
	}
	r = httptest.NewRequest(http.MethodPatch, "/api/v1/me/profile", strings.NewReader(`{"handle":"forbidden"}`))
	r.Header.Set("Authorization", "Bearer valid-token")
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Idempotency-Key", "profile.one")
	w = httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("unknown patch status=%d", w.Code)
	}
}
func TestUnauthorizedAndHealthRoutes(t *testing.T) {
	h := NewServer(&fakeRepo{}, "test").Mux()
	r := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if w.Code != 200 {
		t.Fatal(w.Code)
	}
	r = httptest.NewRequest(http.MethodGet, "/api/v1/me", nil)
	w = httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if w.Code != 401 {
		t.Fatal(w.Code)
	}
}

func TestRestrictedStatusAndRevokeReplay(t *testing.T) {
	repo := &fakeRepo{principal: Principal{AccountID: "account.one", PlayerID: "player.one", SessionID: "session.one", Status: "suspended"}}
	h := NewServer(repo, "test").Mux()
	r := httptest.NewRequest(http.MethodGet, "/api/v1/me", nil)
	r.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("restricted status=%d", w.Code)
	}
	repo.principal.Status = "guest"
	for i := 0; i < 2; i++ {
		r = httptest.NewRequest(http.MethodDelete, "/api/v1/auth/session", nil)
		r.Header.Set("Authorization", "Bearer valid-token")
		w = httptest.NewRecorder()
		h.ServeHTTP(w, r)
		if w.Code != http.StatusNoContent {
			t.Fatalf("revoke replay status=%d", w.Code)
		}
	}
}

func TestProfilePatchValidation(t *testing.T) {
	repo := &fakeRepo{principal: Principal{AccountID: "account.one", PlayerID: "player.one", SessionID: "session.one", Status: "guest"}}
	h := NewServer(repo, "test").Mux()
	cases := []string{`{}`, `{"bio":"` + strings.Repeat("x", 1001) + `"}`, `{"visibility":"hidden"}`, `{"dm_policy":"maybe"}`}
	for _, body := range cases {
		r := httptest.NewRequest(http.MethodPatch, "/api/v1/me/profile", strings.NewReader(body))
		r.Header.Set("Authorization", "Bearer valid-token")
		r.Header.Set("Content-Type", "application/json")
		r.Header.Set("Idempotency-Key", "profile.validation")
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)
		if w.Code != http.StatusBadRequest && w.Code != http.StatusUnprocessableEntity {
			t.Fatalf("profile validation status=%d body=%s", w.Code, w.Body.String())
		}
	}
}

func TestRoutingReturnsStable404And405(t *testing.T) {
	h := NewServer(&fakeRepo{}, "test").Mux()
	r := httptest.NewRequest(http.MethodGet, "/unknown", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if w.Code != http.StatusNotFound {
		t.Fatalf("unknown route status=%d", w.Code)
	}
	r = httptest.NewRequest(http.MethodPost, "/healthz", nil)
	w = httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if w.Code != http.StatusMethodNotAllowed || w.Header().Get("Allow") == "" {
		t.Fatalf("wrong method status/allow=%d/%s", w.Code, w.Header().Get("Allow"))
	}
}

func TestRequestIDReflectionUsesNarrowSafeAlphabet(t *testing.T) {
	h := NewServer(&fakeRepo{}, "test").Mux()

	r := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	r.Header.Set("X-Request-ID", "safe.ID_123:part-value")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if got := w.Header().Get("X-Request-ID"); got != "safe.ID_123:part-value" {
		t.Fatalf("safe request ID was not preserved: %q", got)
	}

	for _, unsafeID := range []string{"bad=value", "bad value", "bad=value\nfield"} {
		r = httptest.NewRequest(http.MethodGet, "/healthz", nil)
		r.Header.Set("X-Request-ID", unsafeID)
		w = httptest.NewRecorder()
		h.ServeHTTP(w, r)
		if got := w.Header().Get("X-Request-ID"); got == unsafeID || got == "" {
			t.Fatalf("unsafe request ID %q was reflected as %q", unsafeID, got)
		}
	}
}
