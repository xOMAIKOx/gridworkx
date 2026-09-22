package identity

import (
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func TestGuestIssuanceIsAtomicAndIdempotent(t *testing.T) {
	store := NewInMemoryStore()
	now := time.Date(2026, 9, 22, 0, 0, 0, 0, time.UTC)
	first, err := store.IssueGuest("request-1", now)
	if err != nil {
		t.Fatal(err)
	}
	second, err := store.IssueGuest("request-1", now.Add(time.Minute))
	if err != nil {
		t.Fatal(err)
	}
	if first.Account.AccountID != second.Account.AccountID || first.Player.PlayerID != second.Player.PlayerID || first.RawSessionToken != second.RawSessionToken {
		t.Fatal("duplicate issuance was not idempotent")
	}
	encoded, _ := json.Marshal(first.Session)
	if strings.Contains(string(encoded), first.RawSessionToken) {
		t.Fatal("raw session token persisted")
	}
	if first.Account.AccountID == first.Player.PlayerID || first.Account.Status != AccountGuest {
		t.Fatal("identity hierarchy is invalid")
	}
}

func TestLinkRequiresProofAndPreservesGuestIdentity(t *testing.T) {
	store := NewInMemoryStore()
	now := time.Date(2026, 9, 22, 0, 0, 0, 0, time.UTC)
	guest, _ := store.IssueGuest("request-1", now)
	_, err := store.LinkExternal(guest.RawSessionToken, ExternalIdentityAssertion{Provider: "google", Issuer: "https://accounts.google.example", Subject: "sub-1"}, now)
	if err != ErrUnverifiedIdentity {
		t.Fatalf("expected unverified identity rejection, got %v", err)
	}
	link, err := store.LinkExternal(guest.RawSessionToken, ExternalIdentityAssertion{Provider: "google", Issuer: "https://accounts.google.example", Subject: "sub-1", Verified: true, ProofReference: "proof-1"}, now)
	if err != nil {
		t.Fatal(err)
	}
	if link.AccountID != guest.Account.AccountID || store.accounts[guest.Account.AccountID].Status != AccountProtected {
		t.Fatal("link did not preserve guest account")
	}
	replay, err := store.LinkExternal(guest.RawSessionToken, ExternalIdentityAssertion{Provider: "google", Issuer: "https://accounts.google.example", Subject: "sub-1", Verified: true}, now)
	if err != nil || replay.LinkID != link.LinkID {
		t.Fatal("same identity link was not idempotent")
	}
	other, _ := store.IssueGuest("request-2", now)
	_, err = store.LinkExternal(other.RawSessionToken, ExternalIdentityAssertion{Provider: "google", Issuer: "https://accounts.google.example", Subject: "sub-1", Verified: true}, now)
	if err != ErrExternalIdentityConflict {
		t.Fatalf("expected cross-account conflict, got %v", err)
	}
}

func TestSessionsRevokeExpireAndHandlesNormalize(t *testing.T) {
	store := NewInMemoryStore()
	now := time.Date(2026, 9, 22, 0, 0, 0, 0, time.UTC)
	guest, _ := store.IssueGuest("request-1", now)
	if err := store.RevokeSession(guest.RawSessionToken, now); err != nil {
		t.Fatal(err)
	}
	if _, err := store.LinkExternal(guest.RawSessionToken, ExternalIdentityAssertion{Provider: "apple", Issuer: "https://apple.example", Subject: "sub", Verified: true}, now); err != ErrInvalidSession {
		t.Fatalf("revoked session accepted: %v", err)
	}
	if !ValidateLocale("pt-BR") || ValidateLocale("xx") || !ValidateTimezone("UTC") || ValidateTimezone("not/a-zone") {
		t.Fatal("locale/timezone validation failed")
	}
	if !ValidateDisplayName("Player") || ValidateDisplayName("bad\u0000name") {
		t.Fatal("display-name validation failed")
	}
	first, err := NormalizeHandle("Cafe\u0301")
	if err != nil {
		t.Fatal(err)
	}
	second, err := NormalizeHandle("Café")
	if err != nil || first.Canonical != second.Canonical || first.Skeleton != second.Skeleton {
		t.Fatal("NFKC handle normalization failed")
	}
	reserved, err := NormalizeHandle("аdmin")
	if err != nil || store.reserveHandle(reserved) != ErrHandleReserved {
		t.Fatalf("expected confusable reserved rejection")
	}
	if _, err := NormalizeHandle("name\u200B"); err != ErrHandleInvalid {
		t.Fatalf("expected invisible-control rejection, got %v", err)
	}
}

func TestPublicProfileExcludesPrivateIdentityState(t *testing.T) {
	profile := PlayerProfile{PlayerID: "player.private", AccountID: "account.private", Handle: "handle", DisplayName: "Name", Locale: "en", Visibility: VisibilityPublic, Discoverable: true}
	public := PublicProfile(profile)
	encoded, _ := json.Marshal(public)
	if strings.Contains(string(encoded), "AccountID") || strings.Contains(string(encoded), "Session") || strings.Contains(string(encoded), "provider") {
		t.Fatal("public profile exposed private identity state")
	}
}
