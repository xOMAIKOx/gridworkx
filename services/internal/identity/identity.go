package identity

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"strings"
	"sync"
	"time"
	"unicode"

	"golang.org/x/text/unicode/norm"
)

const HandleAlgorithmVersion = "unicode-15.1-skeleton-v1"

type AccountStatus string

const (
	AccountGuest              AccountStatus = "guest"
	AccountProtected          AccountStatus = "protected"
	AccountSuspended          AccountStatus = "suspended"
	AccountRecoveryRestricted AccountStatus = "recovery_restricted"
)

type DMPolicy string

const (
	DMEveryone DMPolicy = "everyone"
	DMNobody   DMPolicy = "nobody"
)

type ProfileVisibility string

const (
	VisibilityPublic  ProfileVisibility = "public"
	VisibilityPrivate ProfileVisibility = "private"
)

type Account struct {
	AccountID        string
	Status           AccountStatus
	CreatedAt        time.Time
	AuthorityVersion uint64
}
type Player struct {
	PlayerID  string
	AccountID string
	CreatedAt time.Time
}
type SessionRecord struct {
	SessionID   string
	AccountID   string
	TokenDigest string
	CreatedAt   time.Time
	LastUsedAt  time.Time
	ExpiresAt   time.Time
	RevokedAt   *time.Time
}
type ExternalIdentity struct {
	LinkID      string
	AccountID   string
	Provider    string
	Issuer      string
	Subject     string
	IdentityKey string
	LinkedAt    time.Time
	RevokedAt   *time.Time
}
type PlayerProfile struct {
	PlayerID                string
	AccountID               string
	Handle                  string
	HandleCanonical         string
	HandleSkeleton          string
	DisplayName             string
	AvatarAssetID           string
	Bio                     string
	Locale                  string
	Timezone                string
	Visibility              ProfileVisibility
	DMPolicy                DMPolicy
	Discoverable            bool
	NotificationPreferences map[string]bool
	CreatedAt               time.Time
	UpdatedAt               time.Time
}
type GuestIssuanceResult struct {
	Account         Account
	Player          Player
	Profile         PlayerProfile
	Session         SessionRecord
	RawSessionToken string
}
type ExternalIdentityAssertion struct {
	Provider       string
	Issuer         string
	Subject        string
	Verified       bool
	ProofReference string
}

type IdentityError string

func (e IdentityError) Error() string { return string(e) }

const (
	ErrInvalidIdempotency       IdentityError = "identity idempotency key is invalid"
	ErrIdempotencyConflict      IdentityError = "identity idempotency payload conflicts"
	ErrInvalidSession           IdentityError = "identity session is invalid"
	ErrSessionRevoked           IdentityError = "identity session is revoked or expired"
	ErrUnverifiedIdentity       IdentityError = "external identity proof is not verified"
	ErrExternalIdentityConflict IdentityError = "external identity is linked to another account"
	ErrAccountMergeRequired     IdentityError = "normal self-service account merge is not supported"
	ErrHandleInvalid            IdentityError = "handle is invalid"
	ErrHandleUnavailable        IdentityError = "handle is unavailable"
	ErrHandleReserved           IdentityError = "handle is reserved"
	ErrProfileInvalid           IdentityError = "profile value is invalid"
)

// IdentityStore is a repository/domain seam. Implementations can map these atomic operations to one DB transaction.
type IdentityStore interface {
	IssueGuest(idempotencyKey string, now time.Time) (GuestIssuanceResult, error)
	LinkExternal(sessionToken string, assertion ExternalIdentityAssertion, now time.Time) (ExternalIdentity, error)
	RevokeSession(sessionToken string, now time.Time) error
}

type InMemoryStore struct {
	mu        sync.Mutex
	accounts  map[string]Account
	players   map[string]Player
	profiles  map[string]PlayerProfile
	sessions  map[string]SessionRecord
	links     map[string]ExternalIdentity
	receipts  map[string]GuestIssuanceResult
	handles   map[string]string
	skeletons map[string]string
	reserved  map[string]struct{}
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{accounts: map[string]Account{}, players: map[string]Player{}, profiles: map[string]PlayerProfile{}, sessions: map[string]SessionRecord{}, links: map[string]ExternalIdentity{}, receipts: map[string]GuestIssuanceResult{}, handles: map[string]string{}, skeletons: map[string]string{}, reserved: map[string]struct{}{"gridworks": {}, "admin": {}, "administrator": {}, "moderator": {}, "support": {}, "system": {}, "official": {}}}
}

func (s *InMemoryStore) IssueGuest(idempotencyKey string, now time.Time) (GuestIssuanceResult, error) {
	if strings.TrimSpace(idempotencyKey) == "" {
		return GuestIssuanceResult{}, ErrInvalidIdempotency
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if prior, ok := s.receipts[idempotencyKey]; ok {
		return prior, nil
	}
	accountID, _ := opaqueID("account")
	playerID, _ := opaqueID("player")
	sessionID, _ := opaqueID("session")
	token, err := secureToken()
	if err != nil {
		return GuestIssuanceResult{}, err
	}
	handle, err := NormalizeHandle("guest-" + accountID[len(accountID)-8:])
	if err != nil {
		return GuestIssuanceResult{}, err
	}
	if err := s.reserveHandle(handle); err != nil {
		return GuestIssuanceResult{}, err
	}
	account := Account{AccountID: accountID, Status: AccountGuest, CreatedAt: now, AuthorityVersion: 1}
	player := Player{PlayerID: playerID, AccountID: accountID, CreatedAt: now}
	profile := PlayerProfile{PlayerID: playerID, AccountID: accountID, Handle: handle.Display, HandleCanonical: handle.Canonical, HandleSkeleton: handle.Skeleton, DisplayName: "Guest", Locale: "en", Timezone: "UTC", Visibility: VisibilityPublic, DMPolicy: DMEveryone, Discoverable: true, NotificationPreferences: map[string]bool{}, CreatedAt: now, UpdatedAt: now}
	session := SessionRecord{SessionID: sessionID, AccountID: accountID, TokenDigest: digestToken(token), CreatedAt: now, LastUsedAt: now, ExpiresAt: now.Add(30 * 24 * time.Hour)}
	result := GuestIssuanceResult{Account: account, Player: player, Profile: profile, Session: session, RawSessionToken: token}
	s.accounts[accountID], s.players[playerID], s.profiles[playerID], s.sessions[sessionID] = account, player, profile, session
	s.receipts[idempotencyKey] = result
	return result, nil
}
func (s *InMemoryStore) LinkExternal(sessionToken string, assertion ExternalIdentityAssertion, now time.Time) (ExternalIdentity, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	session, ok := s.findSession(sessionToken)
	if !ok || session.RevokedAt != nil || !now.Before(session.ExpiresAt) {
		return ExternalIdentity{}, ErrInvalidSession
	}
	if !assertion.Verified || strings.TrimSpace(assertion.Provider) == "" || strings.TrimSpace(assertion.Issuer) == "" || strings.TrimSpace(assertion.Subject) == "" {
		return ExternalIdentity{}, ErrUnverifiedIdentity
	}
	key := assertion.Issuer + "\x00" + assertion.Subject
	if existing, ok := s.links[key]; ok {
		if existing.AccountID == session.AccountID {
			return existing, nil
		}
		return ExternalIdentity{}, ErrExternalIdentityConflict
	}
	linkID, _ := opaqueID("identity-link")
	link := ExternalIdentity{LinkID: linkID, AccountID: session.AccountID, Provider: assertion.Provider, Issuer: assertion.Issuer, Subject: assertion.Subject, IdentityKey: key, LinkedAt: now, RevokedAt: nil}
	s.links[key] = link
	account := s.accounts[session.AccountID]
	account.Status = AccountProtected
	account.AuthorityVersion++
	s.accounts[account.AccountID] = account
	return link, nil
}
func (s *InMemoryStore) RevokeSession(token string, now time.Time) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	session, ok := s.findSession(token)
	if !ok {
		return ErrInvalidSession
	}
	session.RevokedAt = &now
	s.sessions[session.SessionID] = session
	return nil
}
func (s *InMemoryStore) findSession(token string) (SessionRecord, bool) {
	d := digestToken(token)
	for _, session := range s.sessions {
		if subtle.ConstantTimeCompare([]byte(session.TokenDigest), []byte(d)) == 1 {
			return session, true
		}
	}
	return SessionRecord{}, false
}
func (s *InMemoryStore) reserveHandle(handle NormalizedHandle) error {
	if _, ok := s.reserved[handle.Skeleton]; ok {
		return ErrHandleReserved
	}
	if _, ok := s.handles[handle.Canonical]; ok {
		return ErrHandleUnavailable
	}
	if _, ok := s.skeletons[handle.Skeleton]; ok {
		return ErrHandleUnavailable
	}
	s.handles[handle.Canonical] = handle.Canonical
	s.skeletons[handle.Skeleton] = handle.Canonical
	return nil
}

func opaqueID(prefix string) (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return prefix + "." + hex.EncodeToString(b), nil
}
func secureToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
func digestToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

type NormalizedHandle struct {
	Display   string
	Canonical string
	Skeleton  string
}

func NormalizeHandle(value string) (NormalizedHandle, error) {
	display := strings.TrimSpace(value)
	if display == "" || len([]rune(display)) < 3 || len([]rune(display)) > 32 {
		return NormalizedHandle{}, ErrHandleInvalid
	}
	normalized := norm.NFKC.String(display)
	var canonical strings.Builder
	var skeleton strings.Builder
	for _, r := range normalized {
		if unicode.IsControl(r) || unicode.Is(unicode.Cf, r) {
			return NormalizedHandle{}, ErrHandleInvalid
		}
		lower := unicode.ToLower(r)
		if unicode.IsSpace(lower) {
			return NormalizedHandle{}, ErrHandleInvalid
		}
		canonical.WriteRune(lower)
		skeleton.WriteRune(confusableRune(lower))
	}
	if canonical.String() == "" {
		return NormalizedHandle{}, ErrHandleInvalid
	}
	return NormalizedHandle{Display: display, Canonical: canonical.String(), Skeleton: skeleton.String()}, nil
}
func confusableRune(r rune) rune {
	switch r {
	case 'а', 'А':
		return 'a'
	case 'е', 'Е':
		return 'e'
	case 'о', 'О':
		return 'o'
	case 'р', 'Р':
		return 'p'
	case 'с', 'С':
		return 'c'
	case 'х', 'Х':
		return 'x'
	case 'у', 'У':
		return 'y'
	case 'і', 'І':
		return 'i'
	case 'ј', 'Ј':
		return 'j'
	case 'ѕ', 'Ѕ':
		return 's'
	default:
		return r
	}
}

func ValidateLocale(value string) bool {
	switch value {
	case "en", "es", "de", "fr", "bg", "ru", "it", "pt-PT", "pt-BR":
		return true
	default:
		return false
	}
}
func ValidateTimezone(value string) bool { _, err := time.LoadLocation(value); return err == nil }
func ValidateDisplayName(value string) bool {
	v := strings.TrimSpace(value)
	if v == "" || len([]rune(v)) > 80 {
		return false
	}
	for _, r := range v {
		if unicode.IsControl(r) || unicode.Is(unicode.Cf, r) {
			return false
		}
	}
	return true
}

// PublicPlayerProfile deliberately excludes account, session, provider and security state.
type PublicPlayerProfile struct {
	PlayerID      string
	Handle        string
	DisplayName   string
	AvatarAssetID string
	Bio           string
	Locale        string
	Visibility    ProfileVisibility
	Discoverable  bool
}

func PublicProfile(profile PlayerProfile) PublicPlayerProfile {
	return PublicPlayerProfile{PlayerID: profile.PlayerID, Handle: profile.Handle, DisplayName: profile.DisplayName, AvatarAssetID: profile.AvatarAssetID, Bio: profile.Bio, Locale: profile.Locale, Visibility: profile.Visibility, Discoverable: profile.Discoverable}
}
