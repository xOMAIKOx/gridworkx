package company

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"strings"
	"sync"
	"time"

	"github.com/xOMAIKOx/gridworkx/services/internal/identity"
)

const OwnershipBPS = uint16(10000)

type CompanyType string

const (
	Holding   CompanyType = "holding"
	Operating CompanyType = "operating"
	System    CompanyType = "system"
)

type CompanyStatus string

const (
	Active    CompanyStatus = "active"
	Suspended CompanyStatus = "suspended"
	Archived  CompanyStatus = "archived"
)

type PrincipalType string

const (
	PlayerPrincipal PrincipalType = "player"
	SystemPrincipal PrincipalType = "system"
)

type OwnerPrincipal struct {
	Type PrincipalType
	ID   string
}
type Company struct {
	CompanyID            string
	Type                 CompanyType
	Name                 string
	NameCanonical        string
	NameSkeleton         string
	Status               CompanyStatus
	Description          string
	LogoAssetID          string
	IndustryID           string
	HeadquartersRegionID string
	Visibility           identity.ProfileVisibility
	CreatedAt            time.Time
	UpdatedAt            time.Time
}
type CompanyGroup struct {
	GroupID       string
	Name          string
	NameCanonical string
	NameSkeleton  string
	Status        CompanyStatus
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
type Ownership struct {
	OwnershipID   string
	EntityID      string
	Owner         OwnerPrincipal
	ShareBPS      uint16
	Active        bool
	EffectiveFrom time.Time
	EffectiveTo   *time.Time
}
type GroupMembership struct {
	MembershipID  string
	CompanyID     string
	GroupID       string
	Active        bool
	EffectiveFrom time.Time
	EffectiveTo   *time.Time
}
type OwnershipEvent struct {
	EventID       string
	EntityID      string
	MutationType  string
	FromOwner     *OwnerPrincipal
	ToOwner       *OwnerPrincipal
	ShareBPS      uint16
	EffectiveTime time.Time
	SourceRef     string
	ActorPlayerID string
	CreatedAt     time.Time
}
type CompanyCreateRequest struct {
	CompanyType    CompanyType
	Name           string
	Owner          OwnerPrincipal
	ActorPlayerID  string
	IdempotencyKey string
}
type OwnershipTransferRequest struct {
	CompanyID      string
	From           OwnerPrincipal
	To             OwnerPrincipal
	ShareBPS       uint16
	ActorPlayerID  string
	IdempotencyKey string
}
type GroupCreateRequest struct {
	Name           string
	Owner          OwnerPrincipal
	ActorPlayerID  string
	IdempotencyKey string
}
type GroupAssignRequest struct {
	CompanyID      string
	GroupID        string
	ActorPlayerID  string
	IdempotencyKey string
}
type MutationReceipt struct {
	MutationType  string
	RequestDigest string
	ResultRef     string
}

type CompanyError string

func (e CompanyError) Error() string { return string(e) }

const (
	ErrInvalidIdempotency  CompanyError = "company idempotency key is invalid"
	ErrIdempotencyConflict CompanyError = "company idempotency payload conflicts"
	ErrCompanyNotFound     CompanyError = "company not found"
	ErrGroupNotFound       CompanyError = "group not found"
	ErrOwnerInvalid        CompanyError = "owner principal is invalid"
	ErrNameInvalid         CompanyError = "company name is invalid"
	ErrNameUnavailable     CompanyError = "company name is unavailable"
	ErrNameReserved        CompanyError = "company name is reserved"
	ErrInvalidShare        CompanyError = "ownership share is invalid"
	ErrOwnershipTotal      CompanyError = "ownership total is invalid"
	ErrInsufficientShare   CompanyError = "ownership share is insufficient"
	ErrOwnershipConflict   CompanyError = "ownership relation conflicts"
	ErrGroupConflict       CompanyError = "group membership conflicts"
	ErrInvalidStatus       CompanyError = "company status is invalid"
	ErrCycle               CompanyError = "company group cycle is invalid"
)

type Store struct {
	mu          sync.Mutex
	random      io.Reader
	companies   map[string]Company
	groups      map[string]CompanyGroup
	ownerships  []Ownership
	memberships []GroupMembership
	history     []OwnershipEvent
	receipts    map[string]MutationReceipt
	names       map[string]string
	skeletons   map[string]string
	reserved    map[string]struct{}
}

func NewStore() *Store { return NewStoreWithRandom(rand.Reader) }
func NewStoreWithRandom(random io.Reader) *Store {
	return &Store{random: random, companies: map[string]Company{}, groups: map[string]CompanyGroup{}, receipts: map[string]MutationReceipt{}, names: map[string]string{}, skeletons: map[string]string{}, reserved: map[string]struct{}{"gridworks": {}, "official": {}, "system": {}, "admin": {}, "administrator": {}, "support": {}, "moderator": {}}}
}

func (s *Store) CreateCompany(req CompanyCreateRequest, now time.Time) (Company, error) {
	if strings.TrimSpace(req.IdempotencyKey) == "" {
		return Company{}, ErrInvalidIdempotency
	}
	if !validOwner(req.Owner) {
		return Company{}, ErrOwnerInvalid
	}
	if req.CompanyType != Holding && req.CompanyType != Operating && req.CompanyType != System {
		return Company{}, ErrInvalidStatus
	}
	digest := requestDigest("company.create", string(req.CompanyType), req.Name, string(req.Owner.Type), req.Owner.ID)
	s.mu.Lock()
	defer s.mu.Unlock()
	if receipt, ok := s.receipts[req.IdempotencyKey]; ok {
		if receipt.MutationType != "company.create" || receipt.RequestDigest != digest {
			return Company{}, ErrIdempotencyConflict
		}
		return s.companies[receipt.ResultRef], nil
	}
	name, err := normalizeCompanyName(req.Name)
	if err != nil {
		return Company{}, ErrNameInvalid
	}
	if _, ok := s.reserved[name.Skeleton]; ok {
		return Company{}, ErrNameReserved
	}
	if _, ok := s.names[name.Canonical]; ok {
		return Company{}, ErrNameUnavailable
	}
	if _, ok := s.skeletons[name.Skeleton]; ok {
		return Company{}, ErrNameUnavailable
	}
	id, err := opaqueID(s.random, "company")
	if err != nil {
		return Company{}, err
	}
	if _, ok := s.companies[id]; ok {
		return Company{}, ErrOwnershipConflict
	}
	company := Company{CompanyID: id, Type: req.CompanyType, Name: name.Display, NameCanonical: name.Canonical, NameSkeleton: name.Skeleton, Status: Active, Visibility: identity.VisibilityPublic, CreatedAt: now, UpdatedAt: now}
	ownershipID, err := opaqueID(s.random, "ownership")
	if err != nil {
		return Company{}, err
	}
	eventID, err := opaqueID(s.random, "ownership-event")
	if err != nil {
		return Company{}, err
	}
	s.companies[id] = company
	s.names[name.Canonical] = id
	s.skeletons[name.Skeleton] = id
	s.ownerships = append(s.ownerships, Ownership{OwnershipID: ownershipID, EntityID: id, Owner: req.Owner, ShareBPS: OwnershipBPS, Active: true, EffectiveFrom: now})
	s.history = append(s.history, OwnershipEvent{EventID: eventID, EntityID: id, MutationType: "company.create", ToOwner: &req.Owner, ShareBPS: OwnershipBPS, EffectiveTime: now, SourceRef: req.IdempotencyKey, ActorPlayerID: req.ActorPlayerID, CreatedAt: now})
	s.receipts[req.IdempotencyKey] = MutationReceipt{MutationType: "company.create", RequestDigest: digest, ResultRef: id}
	return company, nil
}

func (s *Store) CreateGroup(req GroupCreateRequest, now time.Time) (CompanyGroup, error) {
	if strings.TrimSpace(req.IdempotencyKey) == "" {
		return CompanyGroup{}, ErrInvalidIdempotency
	}
	if !validOwner(req.Owner) {
		return CompanyGroup{}, ErrOwnerInvalid
	}
	digest := requestDigest("group.create", req.Name, string(req.Owner.Type), req.Owner.ID)
	s.mu.Lock()
	defer s.mu.Unlock()
	if r, ok := s.receipts[req.IdempotencyKey]; ok {
		if r.MutationType != "group.create" || r.RequestDigest != digest {
			return CompanyGroup{}, ErrIdempotencyConflict
		}
		return s.groups[r.ResultRef], nil
	}
	name, err := normalizeCompanyName(req.Name)
	if err != nil {
		return CompanyGroup{}, ErrNameInvalid
	}
	if _, ok := s.reserved[name.Skeleton]; ok {
		return CompanyGroup{}, ErrNameReserved
	}
	if _, ok := s.names[name.Canonical]; ok {
		return CompanyGroup{}, ErrNameUnavailable
	}
	if _, ok := s.skeletons[name.Skeleton]; ok {
		return CompanyGroup{}, ErrNameUnavailable
	}
	id, err := opaqueID(s.random, "group")
	if err != nil {
		return CompanyGroup{}, err
	}
	if _, ok := s.groups[id]; ok {
		return CompanyGroup{}, ErrOwnershipConflict
	}
	group := CompanyGroup{GroupID: id, Name: name.Display, NameCanonical: name.Canonical, NameSkeleton: name.Skeleton, Status: Active, CreatedAt: now, UpdatedAt: now}
	oid, err := opaqueID(s.random, "ownership")
	if err != nil {
		return CompanyGroup{}, err
	}
	eventID, err := opaqueID(s.random, "ownership-event")
	if err != nil {
		return CompanyGroup{}, err
	}
	s.groups[id] = group
	s.names[name.Canonical] = id
	s.skeletons[name.Skeleton] = id
	s.ownerships = append(s.ownerships, Ownership{OwnershipID: oid, EntityID: id, Owner: req.Owner, ShareBPS: OwnershipBPS, Active: true, EffectiveFrom: now})
	s.history = append(s.history, OwnershipEvent{EventID: eventID, EntityID: id, MutationType: "group.create", ToOwner: &req.Owner, ShareBPS: OwnershipBPS, EffectiveTime: now, SourceRef: req.IdempotencyKey, ActorPlayerID: req.ActorPlayerID, CreatedAt: now})
	s.receipts[req.IdempotencyKey] = MutationReceipt{MutationType: "group.create", RequestDigest: digest, ResultRef: id}
	return group, nil
}

func (s *Store) TransferOwnership(req OwnershipTransferRequest, now time.Time) error {
	if strings.TrimSpace(req.IdempotencyKey) == "" {
		return ErrInvalidIdempotency
	}
	if req.From == req.To || req.ShareBPS == 0 || req.ShareBPS > OwnershipBPS {
		return ErrInvalidShare
	}
	if !validOwner(req.From) || !validOwner(req.To) {
		return ErrOwnerInvalid
	}
	digest := requestDigest("ownership.transfer", req.CompanyID, string(req.From.Type), req.From.ID, string(req.To.Type), req.To.ID, fmt.Sprint(req.ShareBPS))
	s.mu.Lock()
	defer s.mu.Unlock()
	if r, ok := s.receipts[req.IdempotencyKey]; ok {
		if r.MutationType != "ownership.transfer" || r.RequestDigest != digest {
			return ErrIdempotencyConflict
		}
		return nil
	}
	if _, ok := s.companies[req.CompanyID]; !ok {
		return ErrCompanyNotFound
	}
	var source *Ownership
	var target *Ownership
	total := uint32(0)
	for i := range s.ownerships {
		if s.ownerships[i].EntityID != req.CompanyID || !s.ownerships[i].Active {
			continue
		}
		total += uint32(s.ownerships[i].ShareBPS)
		if s.ownerships[i].Owner == req.From {
			source = &s.ownerships[i]
		}
		if s.ownerships[i].Owner == req.To {
			target = &s.ownerships[i]
		}
	}
	if total != uint32(OwnershipBPS) {
		return ErrOwnershipTotal
	}
	if source == nil || source.ShareBPS < req.ShareBPS {
		return ErrInsufficientShare
	}
	sourceShare := source.ShareBPS
	targetExists := target != nil
	targetShare := uint16(0)
	if targetExists {
		targetShare = target.ShareBPS
	}
	sourceID, err := opaqueID(s.random, "ownership")
	if err != nil {
		return err
	}
	targetID, err := opaqueID(s.random, "ownership")
	if err != nil {
		return err
	}
	eventID, err := opaqueID(s.random, "ownership-event")
	if err != nil {
		return err
	}
	for i := range s.ownerships {
		if s.ownerships[i].EntityID == req.CompanyID && s.ownerships[i].Active {
			s.ownerships[i].Active = false
			s.ownerships[i].EffectiveTo = &now
		}
	}
	if sourceShare > req.ShareBPS {
		s.ownerships = append(s.ownerships, Ownership{OwnershipID: sourceID, EntityID: req.CompanyID, Owner: req.From, ShareBPS: sourceShare - req.ShareBPS, Active: true, EffectiveFrom: now})
	}
	if targetExists {
		s.ownerships = append(s.ownerships, Ownership{OwnershipID: targetID, EntityID: req.CompanyID, Owner: req.To, ShareBPS: targetShare + req.ShareBPS, Active: true, EffectiveFrom: now})
	} else {
		s.ownerships = append(s.ownerships, Ownership{OwnershipID: targetID, EntityID: req.CompanyID, Owner: req.To, ShareBPS: req.ShareBPS, Active: true, EffectiveFrom: now})
	}
	s.history = append(s.history, OwnershipEvent{EventID: eventID, EntityID: req.CompanyID, MutationType: "ownership.transfer", FromOwner: &req.From, ToOwner: &req.To, ShareBPS: req.ShareBPS, EffectiveTime: now, SourceRef: req.IdempotencyKey, ActorPlayerID: req.ActorPlayerID, CreatedAt: now})
	s.receipts[req.IdempotencyKey] = MutationReceipt{MutationType: "ownership.transfer", RequestDigest: digest, ResultRef: req.CompanyID}
	return nil
}

func (s *Store) AssignCompany(req GroupAssignRequest, now time.Time) error {
	if strings.TrimSpace(req.IdempotencyKey) == "" {
		return ErrInvalidIdempotency
	}
	digest := requestDigest("group.assign", req.CompanyID, req.GroupID)
	s.mu.Lock()
	defer s.mu.Unlock()
	if r, ok := s.receipts[req.IdempotencyKey]; ok {
		if r.MutationType != "group.assign" || r.RequestDigest != digest {
			return ErrIdempotencyConflict
		}
		return nil
	}
	if _, ok := s.companies[req.CompanyID]; !ok {
		return ErrCompanyNotFound
	}
	if _, ok := s.groups[req.GroupID]; !ok {
		return ErrGroupNotFound
	}
	for i := range s.memberships {
		if s.memberships[i].CompanyID == req.CompanyID && s.memberships[i].Active {
			if s.memberships[i].GroupID == req.GroupID {
				s.receipts[req.IdempotencyKey] = MutationReceipt{MutationType: "group.assign", RequestDigest: digest, ResultRef: s.memberships[i].MembershipID}
				return nil
			}
			return ErrGroupConflict
		}
	}
	id, err := opaqueID(s.random, "membership")
	if err != nil {
		return err
	}
	s.memberships = append(s.memberships, GroupMembership{MembershipID: id, CompanyID: req.CompanyID, GroupID: req.GroupID, Active: true, EffectiveFrom: now})
	s.receipts[req.IdempotencyKey] = MutationReceipt{MutationType: "group.assign", RequestDigest: digest, ResultRef: id}
	return nil
}

func (s *Store) DetachCompany(req GroupAssignRequest, now time.Time) error {
	if strings.TrimSpace(req.IdempotencyKey) == "" {
		return ErrInvalidIdempotency
	}
	digest := requestDigest("group.detach", req.CompanyID, req.GroupID)
	s.mu.Lock()
	defer s.mu.Unlock()
	if receipt, ok := s.receipts[req.IdempotencyKey]; ok {
		if receipt.MutationType != "group.detach" || receipt.RequestDigest != digest {
			return ErrIdempotencyConflict
		}
		return nil
	}
	if _, ok := s.companies[req.CompanyID]; !ok {
		return ErrCompanyNotFound
	}
	if _, ok := s.groups[req.GroupID]; !ok {
		return ErrGroupNotFound
	}
	for i := range s.memberships {
		if s.memberships[i].CompanyID == req.CompanyID && s.memberships[i].Active {
			if s.memberships[i].GroupID != req.GroupID {
				return ErrGroupConflict
			}
			s.memberships[i].Active = false
			s.memberships[i].EffectiveTo = &now
			s.receipts[req.IdempotencyKey] = MutationReceipt{MutationType: "group.detach", RequestDigest: digest, ResultRef: s.memberships[i].MembershipID}
			return nil
		}
	}
	return ErrGroupConflict
}

func (s *Store) ReassignCompany(req GroupAssignRequest, now time.Time) error {
	if strings.TrimSpace(req.IdempotencyKey) == "" {
		return ErrInvalidIdempotency
	}
	digest := requestDigest("group.reassign", req.CompanyID, req.GroupID)
	s.mu.Lock()
	defer s.mu.Unlock()
	if receipt, ok := s.receipts[req.IdempotencyKey]; ok {
		if receipt.MutationType != "group.reassign" || receipt.RequestDigest != digest {
			return ErrIdempotencyConflict
		}
		return nil
	}
	if _, ok := s.companies[req.CompanyID]; !ok {
		return ErrCompanyNotFound
	}
	if _, ok := s.groups[req.GroupID]; !ok {
		return ErrGroupNotFound
	}
	for i := range s.memberships {
		if s.memberships[i].CompanyID == req.CompanyID && s.memberships[i].Active {
			if s.memberships[i].GroupID == req.GroupID {
				return ErrGroupConflict
			}
			s.memberships[i].Active = false
			s.memberships[i].EffectiveTo = &now
		}
	}
	id, err := opaqueID(s.random, "membership")
	if err != nil {
		return err
	}
	s.memberships = append(s.memberships, GroupMembership{MembershipID: id, CompanyID: req.CompanyID, GroupID: req.GroupID, Active: true, EffectiveFrom: now})
	s.receipts[req.IdempotencyKey] = MutationReceipt{MutationType: "group.reassign", RequestDigest: digest, ResultRef: id}
	return nil
}

func (s *Store) ActiveOwnership(entityID string) []Ownership {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := []Ownership{}
	for _, o := range s.ownerships {
		if o.EntityID == entityID && o.Active {
			out = append(out, o)
		}
	}
	return out
}
func (s *Store) History(entityID string) []OwnershipEvent {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := []OwnershipEvent{}
	for _, e := range s.history {
		if e.EntityID == entityID {
			out = append(out, e)
		}
	}
	return out
}
func validOwner(owner OwnerPrincipal) bool {
	if owner.ID == "" {
		return false
	}
	if owner.Type == SystemPrincipal {
		return owner.ID == "principal.gridworks.system"
	}
	return owner.Type == PlayerPrincipal && strings.HasPrefix(owner.ID, "player.")
}

func normalizeCompanyName(value string) (identity.NormalizedHandle, error) {
	display := strings.TrimSpace(value)
	if display == "" {
		return identity.NormalizedHandle{}, ErrNameInvalid
	}
	// Company display names may contain spaces; canonical policy shares WP-007 semantics by using its normalized handle key.
	key := strings.Join(strings.Fields(display), "-")
	normalized, err := identity.NormalizeHandle(key)
	if err != nil {
		return identity.NormalizedHandle{}, ErrNameInvalid
	}
	normalized.Display = display
	return normalized, nil
}

func requestDigest(parts ...string) string {
	sum := sha256.Sum256([]byte(strings.Join(parts, "\x00")))
	return hex.EncodeToString(sum[:])
}
func opaqueID(r io.Reader, prefix string) (string, error) {
	b := make([]byte, 16)
	if _, err := io.ReadFull(r, b); err != nil {
		return "", err
	}
	return prefix + "." + hex.EncodeToString(b), nil
}
