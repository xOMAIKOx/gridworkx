package company

import (
	"io"
	"testing"
	"time"

	"github.com/xOMAIKOx/gridworkx/services/internal/identity"
)

func TestCompanyCreationOwnershipAndIdempotency(t *testing.T) {
	store := NewStore()
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
	store := NewStore()
	now := time.Date(2026, 9, 23, 0, 0, 0, 0, time.UTC)
	company, err := store.CreateCompany(CompanyCreateRequest{CompanyType: Operating, Name: "Operating One", Owner: OwnerPrincipal{Type: PlayerPrincipal, ID: "player.one"}, IdempotencyKey: "company.one"}, now)
	if err != nil {
		t.Fatal(err)
	}
	if err := store.TransferOwnership(OwnershipTransferRequest{CompanyID: company.CompanyID, From: OwnerPrincipal{Type: PlayerPrincipal, ID: "player.one"}, To: OwnerPrincipal{Type: PlayerPrincipal, ID: "player.two"}, ShareBPS: 2500, IdempotencyKey: "transfer.one"}, now.Add(time.Minute)); err != nil {
		t.Fatal(err)
	}
	if err := store.TransferOwnership(OwnershipTransferRequest{CompanyID: company.CompanyID, From: OwnerPrincipal{Type: PlayerPrincipal, ID: "player.one"}, To: OwnerPrincipal{Type: PlayerPrincipal, ID: "player.two"}, ShareBPS: 2500, IdempotencyKey: "transfer.one"}, now.Add(time.Hour)); err != nil {
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
	if err := store.TransferOwnership(OwnershipTransferRequest{CompanyID: company.CompanyID, From: OwnerPrincipal{Type: PlayerPrincipal, ID: "player.one"}, To: OwnerPrincipal{Type: PlayerPrincipal, ID: "player.two"}, ShareBPS: 7501, IdempotencyKey: "transfer.two"}, now.Add(2*time.Hour)); err != ErrInsufficientShare {
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
	store := NewStore()
	now := time.Date(2026, 9, 23, 0, 0, 0, 0, time.UTC)
	company, err := store.CreateCompany(CompanyCreateRequest{CompanyType: Operating, Name: "Full Exit Co", Owner: OwnerPrincipal{Type: PlayerPrincipal, ID: "player.one"}, IdempotencyKey: "company.full"}, now)
	if err != nil {
		t.Fatal(err)
	}
	if err := store.TransferOwnership(OwnershipTransferRequest{CompanyID: company.CompanyID, From: OwnerPrincipal{Type: PlayerPrincipal, ID: "player.one"}, To: OwnerPrincipal{Type: PlayerPrincipal, ID: "player.two"}, ShareBPS: 10_000, IdempotencyKey: "transfer.full"}, now.Add(time.Minute)); err != nil {
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
	store := NewStore()
	now := time.Date(2026, 9, 23, 0, 0, 0, 0, time.UTC)
	company, err := store.CreateCompany(CompanyCreateRequest{CompanyType: System, Name: "Gridworks Treasury", Owner: OwnerPrincipal{Type: SystemPrincipal, ID: "principal.gridworks.system"}, IdempotencyKey: "system.one"}, now)
	if err != nil || len(store.ActiveOwnership(company.CompanyID)) != 1 {
		t.Fatal("system principal ownership failed")
	}
	_, err = store.CreateCompany(CompanyCreateRequest{CompanyType: Operating, Name: "GRIDWORKS", Owner: OwnerPrincipal{Type: PlayerPrincipal, ID: "player.one"}, IdempotencyKey: "reserved.one"}, now)
	if err != ErrNameReserved {
		t.Fatal("reserved company name was not rejected")
	}
	failing := NewStoreWithRandom(failingReader{})
	if _, err := failing.CreateCompany(CompanyCreateRequest{CompanyType: Operating, Name: "Failure Co", Owner: OwnerPrincipal{Type: PlayerPrincipal, ID: "player.one"}, IdempotencyKey: "failure.one"}, now); err == nil || len(failing.companies) != 0 || len(failing.ownerships) != 0 {
		t.Fatal("generation failure mutated company state")
	}
}

type failingReader struct{}

func (failingReader) Read([]byte) (int, error) { return 0, io.ErrUnexpectedEOF }

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
