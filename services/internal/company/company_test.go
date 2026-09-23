package company

import (
	"crypto/rand"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/xOMAIKOx/gridworkx/services/internal/identity"
)

func fixtureResolver() PlayerResolver {
	known := map[string]struct{}{"player.one": {}, "player.two": {}}
	return PlayerResolverFunc(func(playerID string) bool { _, ok := known[playerID]; return ok })
}
func testStore() *Store { return NewStoreWithResolver(rand.Reader, fixtureResolver()) }

func TestCompanyCreationOwnershipAndIdempotency(t *testing.T) {
	store := testStore()
	now := time.Date(2026, 9, 23, 0, 0, 0, 0, time.UTC)
	req := CompanyCreateRequest{CompanyType: Operating, Name: "Acme Works", Owner: OwnerPrincipal{Type: PlayerPrincipal, ID: "player.one"}, IdempotencyKey: "company.create.one"}
	company, err := store.CreateCompany(req, now)
	if err != nil {
		t.Fatal(err)
	}
	replay, err := store.CreateCompany(req, now.Add(time.Hour))
	if err != nil || replay.CompanyID != company.CompanyID {
		t.Fatal("company creation was not idempotent")
	}
	active := store.ActiveOwnership(company.CompanyID)
	if len(active) != 1 || active[0].ShareBPS != OwnershipBPS || active[0].Owner.ID != "player.one" {
		t.Fatal("initial ownership is not exact")
	}
	if len(store.History(company.CompanyID)) != 1 {
		t.Fatal("initial ownership history missing")
	}
	if _, err := store.CreateCompany(CompanyCreateRequest{CompanyType: Operating, Name: "ACME works", Owner: req.Owner, IdempotencyKey: "company.create.two"}, now); err != ErrNameUnavailable {
		t.Fatal("case-insensitive company name collision was not rejected")
	}
}

func TestOwnershipTransferAndGroupHistory(t *testing.T) {
	store := testStore()
	now := time.Date(2026, 9, 23, 0, 0, 0, 0, time.UTC)
	company, err := store.CreateCompany(CompanyCreateRequest{CompanyType: Operating, Name: "Operating One", Owner: OwnerPrincipal{Type: PlayerPrincipal, ID: "player.one"}, IdempotencyKey: "company.one"}, now)
	if err != nil {
		t.Fatal(err)
	}
	if err := store.TransferOwnership(OwnershipTransferRequest{EntityType: "company", EntityID: company.CompanyID, From: OwnerPrincipal{Type: PlayerPrincipal, ID: "player.one"}, To: OwnerPrincipal{Type: PlayerPrincipal, ID: "player.two"}, ShareBPS: 2500, IdempotencyKey: "transfer.one"}, now.Add(time.Minute)); err != nil {
		t.Fatal(err)
	}
	if err := store.TransferOwnership(OwnershipTransferRequest{EntityType: "company", EntityID: company.CompanyID, From: OwnerPrincipal{Type: PlayerPrincipal, ID: "player.one"}, To: OwnerPrincipal{Type: PlayerPrincipal, ID: "player.two"}, ShareBPS: 2500, IdempotencyKey: "transfer.one"}, now.Add(time.Hour)); err != nil {
		t.Fatal(err)
	}
	active := store.ActiveOwnership(company.CompanyID)
	total := uint16(0)
	for _, ownership := range active {
		total += ownership.ShareBPS
	}
	if total != OwnershipBPS || len(store.History(company.CompanyID)) != 2 {
		t.Fatal("ownership transfer did not preserve exact total/history")
	}
	if err := store.TransferOwnership(OwnershipTransferRequest{EntityType: "company", EntityID: company.CompanyID, From: OwnerPrincipal{Type: PlayerPrincipal, ID: "player.one"}, To: OwnerPrincipal{Type: PlayerPrincipal, ID: "player.two"}, ShareBPS: 7501, IdempotencyKey: "transfer.two"}, now.Add(2*time.Hour)); err != ErrInsufficientShare {
		t.Fatal("insufficient ownership was not rejected")
	}
	group, err := store.CreateGroup(GroupCreateRequest{Name: "Acme Group", Owner: OwnerPrincipal{Type: PlayerPrincipal, ID: "player.one"}, IdempotencyKey: "group.one"}, now)
	if err != nil {
		t.Fatal(err)
	}
	if err := store.AssignCompany(GroupAssignRequest{CompanyID: company.CompanyID, GroupID: group.GroupID, IdempotencyKey: "assign.one"}, now); err != nil {
		t.Fatal(err)
	}
	if err := store.AssignCompany(GroupAssignRequest{CompanyID: company.CompanyID, GroupID: group.GroupID, IdempotencyKey: "assign.one"}, now); err != nil {
		t.Fatal(err)
	}
}

func TestFullExitAndGroupDetachReassignHistory(t *testing.T) {
	store := testStore()
	now := time.Date(2026, 9, 23, 0, 0, 0, 0, time.UTC)
	company, err := store.CreateCompany(CompanyCreateRequest{CompanyType: Operating, Name: "Full Exit Co", Owner: OwnerPrincipal{Type: PlayerPrincipal, ID: "player.one"}, IdempotencyKey: "company.full"}, now)
	if err != nil {
		t.Fatal(err)
	}
	if err := store.TransferOwnership(OwnershipTransferRequest{EntityType: "company", EntityID: company.CompanyID, From: OwnerPrincipal{Type: PlayerPrincipal, ID: "player.one"}, To: OwnerPrincipal{Type: PlayerPrincipal, ID: "player.two"}, ShareBPS: 10_000, IdempotencyKey: "transfer.full"}, now.Add(time.Minute)); err != nil {
		t.Fatal(err)
	}
	active := store.ActiveOwnership(company.CompanyID)
	if len(active) != 1 || active[0].ShareBPS != 10_000 || active[0].Owner.ID != "player.two" {
		t.Fatal("full exit created invalid active ownership")
	}
	g1, _ := store.CreateGroup(GroupCreateRequest{Name: "Group One", Owner: OwnerPrincipal{Type: PlayerPrincipal, ID: "player.one"}, IdempotencyKey: "group.one"}, now)
	g2, _ := store.CreateGroup(GroupCreateRequest{Name: "Group Two", Owner: OwnerPrincipal{Type: PlayerPrincipal, ID: "player.one"}, IdempotencyKey: "group.two"}, now)
	if err := store.AssignCompany(GroupAssignRequest{CompanyID: company.CompanyID, GroupID: g1.GroupID, IdempotencyKey: "assign.one"}, now); err != nil {
		t.Fatal(err)
	}
	if err := store.ReassignCompany(GroupAssignRequest{CompanyID: company.CompanyID, GroupID: g2.GroupID, IdempotencyKey: "reassign.one"}, now.Add(time.Hour)); err != nil {
		t.Fatal(err)
	}
	if err := store.DetachCompany(GroupAssignRequest{CompanyID: company.CompanyID, GroupID: g2.GroupID, IdempotencyKey: "detach.one"}, now.Add(2*time.Hour)); err != nil {
		t.Fatal(err)
	}
	if err := store.DetachCompany(GroupAssignRequest{CompanyID: company.CompanyID, GroupID: g2.GroupID, IdempotencyKey: "detach.one"}, now.Add(3*time.Hour)); err != nil {
		t.Fatal("duplicate detach replay should be idempotent")
	}
}

func TestSystemOwnershipReservedNamesAndGenerationFailure(t *testing.T) {
	store := testStore()
	now := time.Date(2026, 9, 23, 0, 0, 0, 0, time.UTC)
	company, err := store.CreateCompany(CompanyCreateRequest{CompanyType: System, Name: "Gridworks Treasury", Owner: OwnerPrincipal{Type: SystemPrincipal, ID: "principal.gridworks.system"}, IdempotencyKey: "system.one"}, now)
	if err != nil || len(store.ActiveOwnership(company.CompanyID)) != 1 {
		t.Fatal("system principal ownership failed")
	}
	_, err = store.CreateCompany(CompanyCreateRequest{CompanyType: Operating, Name: "GRIDWORKS", Owner: OwnerPrincipal{Type: PlayerPrincipal, ID: "player.one"}, IdempotencyKey: "reserved.one"}, now)
	if err != ErrNameReserved {
		t.Fatal("reserved company name was not rejected")
	}
	failing := NewStoreWithResolver(failingReader{}, fixtureResolver())
	if _, err := failing.CreateCompany(CompanyCreateRequest{CompanyType: Operating, Name: "Failure Co", Owner: OwnerPrincipal{Type: PlayerPrincipal, ID: "player.one"}, IdempotencyKey: "failure.one"}, now); err == nil || len(failing.companies) != 0 || len(failing.ownerships) != 0 {
		t.Fatal("generation failure mutated company state")
	}
}

type failingReader struct{}

func (failingReader) Read([]byte) (int, error) { return 0, io.ErrUnexpectedEOF }

func TestCompanyNameDisplayBoundsAreIndependentFromCanonicalKey(t *testing.T) {
	store := testStore()
	now := time.Date(2026, 9, 23, 0, 0, 0, 0, time.UTC)
	longName := "ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	if _, err := store.CreateCompany(CompanyCreateRequest{CompanyType: Operating, Name: longName, Owner: OwnerPrincipal{Type: PlayerPrincipal, ID: "player.one"}, IdempotencyKey: "long-name"}, now); err != nil {
		t.Fatal("valid >32-rune company name rejected", err)
	}
	overlong := strings.Repeat("a ", 41)
	if _, err := store.CreateCompany(CompanyCreateRequest{CompanyType: Operating, Name: overlong, Owner: OwnerPrincipal{Type: PlayerPrincipal, ID: "player.one"}, IdempotencyKey: "overlong-name"}, now); err != ErrNameInvalid {
		t.Fatal("over-80 display name was accepted after whitespace collapse")
	}
}

func TestCompanyProfileUsesSharedIdentityNormalization(t *testing.T) {
	first, err := identity.NormalizeHandle("Café")
	if err != nil {
		t.Fatal(err)
	}
	second, err := identity.NormalizeHandle("Cafe\u0301")
	if err != nil || first.Canonical != second.Canonical {
		t.Fatal("company/player normalization contract diverged")
	}
}

func TestGroupOwnershipTransferAndPlayerResolution(t *testing.T) {
	store := testStore()
	now := time.Date(2026, 9, 23, 0, 0, 0, 0, time.UTC)
	group, err := store.CreateGroup(GroupCreateRequest{Name: "Transfer Group", Owner: OwnerPrincipal{Type: PlayerPrincipal, ID: "player.one"}, IdempotencyKey: "group.transfer"}, now)
	if err != nil {
		t.Fatal(err)
	}
	if err := store.TransferOwnership(OwnershipTransferRequest{EntityType: "group", EntityID: group.GroupID, From: OwnerPrincipal{Type: PlayerPrincipal, ID: "player.one"}, To: OwnerPrincipal{Type: PlayerPrincipal, ID: "player.two"}, ShareBPS: 2_000, IdempotencyKey: "group.transfer.ownership"}, now.Add(time.Minute)); err != nil {
		t.Fatal(err)
	}
	if err := store.TransferOwnership(OwnershipTransferRequest{EntityType: "group", EntityID: group.GroupID, From: OwnerPrincipal{Type: PlayerPrincipal, ID: "player.one"}, To: OwnerPrincipal{Type: PlayerPrincipal, ID: "player.two"}, ShareBPS: 2_000, IdempotencyKey: "group.transfer.ownership"}, now.Add(2*time.Minute)); err != nil {
		t.Fatal("group transfer replay failed")
	}
	if err := store.TransferOwnership(OwnershipTransferRequest{EntityType: "group", EntityID: group.GroupID, From: OwnerPrincipal{Type: PlayerPrincipal, ID: "player.one"}, To: OwnerPrincipal{Type: SystemPrincipal, ID: "principal.gridworks.system"}, ShareBPS: 2_000, IdempotencyKey: "group.transfer.ownership"}, now.Add(3*time.Minute)); err != ErrIdempotencyConflict {
		t.Fatal("group contradictory idempotency was accepted")
	}
	if err := store.TransferOwnership(OwnershipTransferRequest{EntityType: "group", EntityID: group.GroupID, From: OwnerPrincipal{Type: PlayerPrincipal, ID: "player.one"}, To: OwnerPrincipal{Type: SystemPrincipal, ID: "principal.gridworks.system"}, ShareBPS: 8_000, IdempotencyKey: "group.transfer.full"}, now.Add(4*time.Minute)); err != nil {
		t.Fatal("group full exit failed")
	}
	if err := store.TransferOwnership(OwnershipTransferRequest{EntityType: "group", EntityID: "group.missing", From: OwnerPrincipal{Type: PlayerPrincipal, ID: "player.one"}, To: OwnerPrincipal{Type: SystemPrincipal, ID: "principal.gridworks.system"}, ShareBPS: 1, IdempotencyKey: "group.invalid"}, now.Add(5*time.Minute)); err != ErrGroupNotFound {
		t.Fatal("invalid group entity was accepted")
	}
	active := store.ActiveOwnership(group.GroupID)
	total := uint16(0)
	for _, ownership := range active {
		total += ownership.ShareBPS
	}
	if total != OwnershipBPS {
		t.Fatal("group transfer did not preserve exact ownership total")
	}
	if _, err := store.CreateCompany(CompanyCreateRequest{CompanyType: Operating, Name: "Unknown Owner Co", Owner: OwnerPrincipal{Type: PlayerPrincipal, ID: "player.does-not-exist"}, IdempotencyKey: "unknown-owner"}, now); err != ErrOwnerInvalid {
		t.Fatal("unknown player principal was accepted")
	}
}

func TestReassignFailureDoesNotCloseExistingMembership(t *testing.T) {
	store := testStore()
	now := time.Date(2026, 9, 23, 0, 0, 0, 0, time.UTC)
	company, _ := store.CreateCompany(CompanyCreateRequest{CompanyType: Operating, Name: "Atomic Reassign Co", Owner: OwnerPrincipal{Type: PlayerPrincipal, ID: "player.one"}, IdempotencyKey: "atomic.company"}, now)
	groupA, _ := store.CreateGroup(GroupCreateRequest{Name: "Atomic Group A", Owner: OwnerPrincipal{Type: PlayerPrincipal, ID: "player.one"}, IdempotencyKey: "atomic.group.a"}, now)
	groupB, _ := store.CreateGroup(GroupCreateRequest{Name: "Atomic Group B", Owner: OwnerPrincipal{Type: PlayerPrincipal, ID: "player.one"}, IdempotencyKey: "atomic.group.b"}, now)
	if err := store.AssignCompany(GroupAssignRequest{CompanyID: company.CompanyID, GroupID: groupA.GroupID, IdempotencyKey: "atomic.assign"}, now); err != nil {
		t.Fatal(err)
	}
	store.random = failingReader{}
	if err := store.ReassignCompany(GroupAssignRequest{CompanyID: company.CompanyID, GroupID: groupB.GroupID, IdempotencyKey: "atomic.reassign"}, now.Add(time.Minute)); err == nil {
		t.Fatal("expected reassign generation failure")
	}
	if len(store.memberships) != 1 || !store.memberships[0].Active || store.memberships[0].GroupID != groupA.GroupID {
		t.Fatal("failed reassign closed existing membership")
	}
}
