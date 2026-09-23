package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"reflect"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/xOMAIKOx/gridworkx/services/api/internal/api"
)

type stepEntropy struct {
	call   int
	failAt int
}

func (s *stepEntropy) Read(p []byte) (int, error) {
	s.call++
	if s.failAt > 0 && s.call == s.failAt {
		return 0, errors.New("entropy failed")
	}
	for i := range p {
		p[i] = byte(s.call)
	}
	return len(p), nil
}

func newMockRepository(t *testing.T, entropy Entropy) (*Repository, sqlmock.Sqlmock, func()) {
	t.Helper()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	return NewForTesting(db, entropy), mock, func() { _ = db.Close() }
}

func TestIssueGuestSuccessAndReplayUseDurableReceipt(t *testing.T) {
	now := time.Date(2026, 9, 23, 10, 0, 0, 0, time.UTC)
	exp := now.Add(30 * 24 * time.Hour)
	entropy := &stepEntropy{}
	repo, mock, closeFn := newMockRepository(t, entropy)
	defer closeFn()

	accountID := "account." + strings.Repeat("01", 16)
	playerID := "player." + strings.Repeat("02", 16)
	sessionID := "session." + strings.Repeat("03", 16)
	raw := strings.Repeat("04", 32)
	ref := `{"account_id":"` + accountID + `","player_id":"` + playerID + `","session_id":"` + sessionID + `"}`

	mock.ExpectBegin()
	mock.ExpectExec("SELECT pg_advisory_xact_lock").WithArgs("identity\x00guest.success").WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectQuery("SELECT mutation_type,request_digest,result_ref FROM gridworks.identity_mutation_receipts").WithArgs("guest.success").WillReturnError(sql.ErrNoRows)
	mock.ExpectExec("INSERT INTO gridworks.accounts").WithArgs(accountID, now).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO gridworks.players").WithArgs(playerID, accountID, now).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO gridworks.player_profiles").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO gridworks.guest_sessions").WithArgs(sessionID, accountID, digest(raw), now, exp).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO gridworks.identity_mutation_receipts").WithArgs("guest.success", digest("identity.guest_issue\x00guest.success"), ref).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	first, err := repo.IssueGuest(context.Background(), "guest.success", now)
	if err != nil {
		t.Fatal(err)
	}
	if first.RawSessionToken != raw || first.Replay {
		t.Fatalf("unexpected first result: %+v", first)
	}
	var firstProfile map[string]any
	if err := json.Unmarshal(first.Profile, &firstProfile); err != nil {
		t.Fatal(err)
	}
	if strings.Contains(ref, raw) {
		t.Fatal("raw token leaked into receipt")
	}

	profile := []byte(`{"player_id":"` + playerID + `","handle":"guest-01010101","display_name":"Guest","avatar_asset_id":"","bio":"","locale":"en","timezone":"UTC","visibility":"public","dm_policy":"everyone","discoverable":true,"notification_preferences":{}}`)
	mock.ExpectBegin()
	mock.ExpectExec("SELECT pg_advisory_xact_lock").WithArgs("identity\x00guest.success").WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectQuery("SELECT mutation_type,request_digest,result_ref FROM gridworks.identity_mutation_receipts").WithArgs("guest.success").WillReturnRows(
		sqlmock.NewRows([]string{"mutation_type", "request_digest", "result_ref"}).AddRow("identity.guest_issue", digest("identity.guest_issue\x00guest.success"), ref),
	)
	mock.ExpectQuery("SELECT s.expires_at,json_build_object").WithArgs(accountID, sessionID, playerID).WillReturnRows(
		sqlmock.NewRows([]string{"expires_at", "profile"}).AddRow(exp, profile),
	)
	mock.ExpectCommit()

	replay, err := repo.IssueGuest(context.Background(), "guest.success", now)
	if err != nil {
		t.Fatal(err)
	}
	if !replay.Replay || replay.RawSessionToken != "" || replay.AccountID != accountID || replay.PlayerID != playerID || replay.SessionID != sessionID {
		t.Fatalf("unexpected replay: %+v", replay)
	}
	var replayProfile map[string]any
	if err := json.Unmarshal(replay.Profile, &replayProfile); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(firstProfile, replayProfile) {
		t.Fatalf("first/replay profile drift: %#v != %#v", firstProfile, replayProfile)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestUpdateProfileLocksBeforeMutationAndReplays(t *testing.T) {
	now := time.Date(2026, 9, 23, 10, 0, 0, 0, time.UTC)
	repo, mock, closeFn := newMockRepository(t, &stepEntropy{})
	defer closeFn()
	p := api.Principal{AccountID: "account.one", PlayerID: "player.one", SessionID: "session.one", Status: "guest"}
	display := "Updated"
	prefs := map[string]bool{"email": true}
	patch := api.ProfilePatch{DisplayName: &display, NotificationPreferences: prefs}
	payload, _ := jsonMarshal(patch)
	reqDigest := digest("profile.update\x00player.one\x00" + payload)
	meRows := func() *sqlmock.Rows {
		return sqlmock.NewRows([]string{"account_id", "player_id", "status", "profile"}).AddRow("account.one", "player.one", "guest", []byte(`{"display_name":"Updated"}`))
	}

	mock.ExpectBegin()
	mock.ExpectExec("SELECT pg_advisory_xact_lock").WithArgs("identity\x00profile.same").WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectQuery("SELECT mutation_type,request_digest,result_ref FROM gridworks.identity_mutation_receipts").WithArgs("profile.same").WillReturnError(sql.ErrNoRows)
	mock.ExpectExec(regexp.QuoteMeta("notification_preferences=COALESCE($8::jsonb,notification_preferences)")).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec("INSERT INTO gridworks.identity_mutation_receipts").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	mock.ExpectQuery("SELECT a.account_id,p.player_id,a.status,json_build_object").WithArgs("account.one", "player.one").WillReturnRows(meRows())

	if _, err := repo.UpdateProfile(context.Background(), p, "profile.same", patch, now); err != nil {
		t.Fatal(err)
	}

	mock.ExpectBegin()
	mock.ExpectExec("SELECT pg_advisory_xact_lock").WithArgs("identity\x00profile.same").WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectQuery("SELECT mutation_type,request_digest,result_ref FROM gridworks.identity_mutation_receipts").WithArgs("profile.same").WillReturnRows(
		sqlmock.NewRows([]string{"mutation_type", "request_digest", "result_ref"}).AddRow("profile.update", reqDigest, "player.one"),
	)
	mock.ExpectCommit()
	mock.ExpectQuery("SELECT a.account_id,p.player_id,a.status,json_build_object").WithArgs("account.one", "player.one").WillReturnRows(meRows())
	if _, err := repo.UpdateProfile(context.Background(), p, "profile.same", patch, now); err != nil {
		t.Fatal(err)
	}

	other := "Different"
	conflicting := api.ProfilePatch{DisplayName: &other}
	mock.ExpectBegin()
	mock.ExpectExec("SELECT pg_advisory_xact_lock").WithArgs("identity\x00profile.same").WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectQuery("SELECT mutation_type,request_digest,result_ref FROM gridworks.identity_mutation_receipts").WithArgs("profile.same").WillReturnRows(
		sqlmock.NewRows([]string{"mutation_type", "request_digest", "result_ref"}).AddRow("profile.update", reqDigest, "player.one"),
	)
	mock.ExpectRollback()
	if _, err := repo.UpdateProfile(context.Background(), p, "profile.same", conflicting, now); !errors.Is(err, api.ErrConflict) {
		t.Fatalf("conflict=%v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestCompanyAndGroupCreateLockOwnerReplayAndConflict(t *testing.T) {
	now := time.Date(2026, 9, 23, 10, 0, 0, 0, time.UTC)
	p := api.Principal{AccountID: "account.one", PlayerID: "player.one", SessionID: "session.one", Status: "guest"}

	t.Run("company", func(t *testing.T) {
		repo, mock, closeFn := newMockRepository(t, &stepEntropy{})
		defer closeFn()
		req := api.CompanyCreateRequest{CompanyType: "operating", Name: "  Acme  "}
		expectedDigest := digest("company.create\x00operating\x00  Acme  \x00player\x00player.one")
		mock.ExpectBegin()
		mock.ExpectExec("SELECT pg_advisory_xact_lock").WithArgs("company\x00company.same").WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectQuery("SELECT mutation_type,request_digest,result_ref FROM gridworks.company_mutation_receipts").WithArgs("company.same").WillReturnError(sql.ErrNoRows)
		mock.ExpectExec("INSERT INTO gridworks.companies").WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec("INSERT INTO gridworks.company_ownership").WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), "player.one", now).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec("INSERT INTO gridworks.ownership_history").WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec("INSERT INTO gridworks.company_mutation_receipts").WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()
		first, err := repo.CreateCompany(context.Background(), p, "company.same", req, now)
		if err != nil {
			t.Fatal(err)
		}
		if first.Name != "Acme" {
			t.Fatalf("trimmed result=%q", first.Name)
		}

		mock.ExpectBegin()
		mock.ExpectExec("SELECT pg_advisory_xact_lock").WithArgs("company\x00company.same").WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectQuery("SELECT mutation_type,request_digest,result_ref FROM gridworks.company_mutation_receipts").WithArgs("company.same").WillReturnRows(
			sqlmock.NewRows([]string{"mutation_type", "request_digest", "result_ref"}).AddRow("company.create", expectedDigest, first.ID),
		)
		mock.ExpectQuery("SELECT name_display FROM gridworks.companies").WithArgs(first.ID).WillReturnRows(sqlmock.NewRows([]string{"name_display"}).AddRow("Acme"))
		mock.ExpectCommit()
		replay, err := repo.CreateCompany(context.Background(), p, "company.same", req, now)
		if err != nil || replay != first {
			t.Fatalf("replay=%+v err=%v first=%+v", replay, err, first)
		}

		mock.ExpectBegin()
		mock.ExpectExec("SELECT pg_advisory_xact_lock").WithArgs("company\x00company.same").WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectQuery("SELECT mutation_type,request_digest,result_ref FROM gridworks.company_mutation_receipts").WithArgs("company.same").WillReturnRows(
			sqlmock.NewRows([]string{"mutation_type", "request_digest", "result_ref"}).AddRow("company.create", expectedDigest, first.ID),
		)
		mock.ExpectRollback()
		_, err = repo.CreateCompany(context.Background(), p, "company.same", api.CompanyCreateRequest{CompanyType: "operating", Name: "Different"}, now)
		if !errors.Is(err, api.ErrConflict) {
			t.Fatalf("contradictory company=%v", err)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("group", func(t *testing.T) {
		repo, mock, closeFn := newMockRepository(t, &stepEntropy{})
		defer closeFn()
		req := api.GroupCreateRequest{Name: "  Holding One  "}
		expectedDigest := digest("group.create\x00  Holding One  \x00player\x00player.one")
		mock.ExpectBegin()
		mock.ExpectExec("SELECT pg_advisory_xact_lock").WithArgs("company\x00group.same").WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectQuery("SELECT mutation_type,request_digest,result_ref FROM gridworks.company_mutation_receipts").WithArgs("group.same").WillReturnError(sql.ErrNoRows)
		mock.ExpectExec("INSERT INTO gridworks.company_groups").WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec("INSERT INTO gridworks.company_ownership").WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), "player.one", now).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec("INSERT INTO gridworks.ownership_history").WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec("INSERT INTO gridworks.company_mutation_receipts").WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()
		first, err := repo.CreateGroup(context.Background(), p, "group.same", req, now)
		if err != nil {
			t.Fatal(err)
		}
		if first.Name != "Holding One" {
			t.Fatalf("trimmed group=%q", first.Name)
		}

		mock.ExpectBegin()
		mock.ExpectExec("SELECT pg_advisory_xact_lock").WithArgs("company\x00group.same").WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectQuery("SELECT mutation_type,request_digest,result_ref FROM gridworks.company_mutation_receipts").WithArgs("group.same").WillReturnRows(
			sqlmock.NewRows([]string{"mutation_type", "request_digest", "result_ref"}).AddRow("group.create", expectedDigest, first.ID),
		)
		mock.ExpectQuery("SELECT name_display FROM gridworks.company_groups").WithArgs(first.ID).WillReturnRows(sqlmock.NewRows([]string{"name_display"}).AddRow("Holding One"))
		mock.ExpectCommit()
		replay, err := repo.CreateGroup(context.Background(), p, "group.same", req, now)
		if err != nil || replay != first {
			t.Fatalf("group replay=%+v err=%v first=%+v", replay, err, first)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatal(err)
		}
	})
}

func TestCompanyAndGroupEntropyFailuresRollback(t *testing.T) {
	now := time.Date(2026, 9, 23, 10, 0, 0, 0, time.UTC)
	p := api.Principal{PlayerID: "player.one"}

	t.Run("company entity id failure", func(t *testing.T) {
		repo, mock, closeFn := newMockRepository(t, &stepEntropy{failAt: 1})
		defer closeFn()
		mock.ExpectBegin()
		mock.ExpectExec("SELECT pg_advisory_xact_lock").WithArgs("company\x00company.entity.fail").WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectQuery("SELECT mutation_type,request_digest,result_ref FROM gridworks.company_mutation_receipts").WithArgs("company.entity.fail").WillReturnError(sql.ErrNoRows)
		mock.ExpectRollback()
		if _, err := repo.CreateCompany(context.Background(), p, "company.entity.fail", api.CompanyCreateRequest{CompanyType: "operating", Name: "Acme"}, now); err == nil {
			t.Fatal("expected company entity entropy failure")
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("company ownership id failure", func(t *testing.T) {
		repo, mock, closeFn := newMockRepository(t, &stepEntropy{failAt: 2})
		defer closeFn()
		mock.ExpectBegin()
		mock.ExpectExec("SELECT pg_advisory_xact_lock").WithArgs("company\x00company.fail").WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectQuery("SELECT mutation_type,request_digest,result_ref FROM gridworks.company_mutation_receipts").WithArgs("company.fail").WillReturnError(sql.ErrNoRows)
		mock.ExpectExec("INSERT INTO gridworks.companies").WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectRollback()
		if _, err := repo.CreateCompany(context.Background(), p, "company.fail", api.CompanyCreateRequest{CompanyType: "operating", Name: "Acme"}, now); err == nil {
			t.Fatal("expected company entropy failure")
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("group entity id failure", func(t *testing.T) {
		repo, mock, closeFn := newMockRepository(t, &stepEntropy{failAt: 1})
		defer closeFn()
		mock.ExpectBegin()
		mock.ExpectExec("SELECT pg_advisory_xact_lock").WithArgs("company\x00group.entity.fail").WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectQuery("SELECT mutation_type,request_digest,result_ref FROM gridworks.company_mutation_receipts").WithArgs("group.entity.fail").WillReturnError(sql.ErrNoRows)
		mock.ExpectRollback()
		if _, err := repo.CreateGroup(context.Background(), p, "group.entity.fail", api.GroupCreateRequest{Name: "Holding"}, now); err == nil {
			t.Fatal("expected group entity entropy failure")
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("group history id failure", func(t *testing.T) {
		repo, mock, closeFn := newMockRepository(t, &stepEntropy{failAt: 3})
		defer closeFn()
		mock.ExpectBegin()
		mock.ExpectExec("SELECT pg_advisory_xact_lock").WithArgs("company\x00group.fail").WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectQuery("SELECT mutation_type,request_digest,result_ref FROM gridworks.company_mutation_receipts").WithArgs("group.fail").WillReturnError(sql.ErrNoRows)
		mock.ExpectExec("INSERT INTO gridworks.company_groups").WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectRollback()
		if _, err := repo.CreateGroup(context.Background(), p, "group.fail", api.GroupCreateRequest{Name: "Holding"}, now); err == nil {
			t.Fatal("expected group entropy failure")
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatal(err)
		}
	})
}

func TestRepositoryHonorsParentContextCancellation(t *testing.T) {
	repo, mock, closeFn := newMockRepository(t, &stepEntropy{})
	defer closeFn()
	mock.ExpectQuery("SELECT player_id,handle_display,display_name").WithArgs("player.one").WillDelayFor(50 * time.Millisecond).WillReturnRows(
		sqlmock.NewRows([]string{"player_id", "handle_display", "display_name", "avatar_asset_id", "bio", "locale", "visibility", "discoverable"}).AddRow("player.one", "one", "One", "", "", "en", "public", true),
	)
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
	defer cancel()
	if _, err := repo.GetPublicPlayer(ctx, "player.one"); err == nil {
		t.Fatal("expected context cancellation")
	}
}

func jsonMarshal(v any) (string, error) {
	b, err := json.Marshal(v)
	return string(b), err
}

func TestGuestReplayRejectsCrossWiredReceipt(t *testing.T) {
	repo, mock, closeFn := newMockRepository(t, &stepEntropy{})
	defer closeFn()
	now := time.Date(2026, 9, 23, 10, 0, 0, 0, time.UTC)
	ref := `{"account_id":"account.one","player_id":"player.two","session_id":"session.one"}`
	mock.ExpectBegin()
	mock.ExpectExec("SELECT pg_advisory_xact_lock").WithArgs("identity\x00cross-wired").WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectQuery("SELECT mutation_type,request_digest,result_ref FROM gridworks.identity_mutation_receipts").WithArgs("cross-wired").WillReturnRows(sqlmock.NewRows([]string{"mutation_type", "request_digest", "result_ref"}).AddRow("identity.guest_issue", digest("identity.guest_issue\x00cross-wired"), ref))
	mock.ExpectQuery("SELECT s.expires_at,json_build_object").WithArgs("account.one", "session.one", "player.two").WillReturnError(sql.ErrNoRows)
	mock.ExpectRollback()
	if _, err := repo.IssueGuest(context.Background(), "cross-wired", now); err == nil {
		t.Fatal("cross-wired receipt was accepted")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}
