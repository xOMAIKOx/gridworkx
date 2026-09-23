package postgres

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/xOMAIKOx/gridworkx/services/api/internal/api"
)

func TestGuestReceiptReferenceIsPostgresSafeAndSecretFree(t *testing.T) {
	ref, err := json.Marshal(guestReceiptRef{AccountID: "account.one", PlayerID: "player.one", SessionID: "session.one"})
	if err != nil {
		t.Fatal(err)
	}
	if strings.ContainsRune(string(ref), 0) || strings.Contains(string(ref), "raw-token") {
		t.Fatal("unsafe guest receipt reference")
	}
	var decoded guestReceiptRef
	if err := json.Unmarshal(ref, &decoded); err != nil || decoded.AccountID != "account.one" {
		t.Fatal("receipt reference did not round-trip")
	}
}

func TestDatabaseConflictMapping(t *testing.T) {
	if got := mapDB(&pgconn.PgError{Code: "23505"}); got != api.ErrConflict {
		t.Fatalf("duplicate mapped to %v", got)
	}
	if got := mapDB(&pgconn.PgError{Code: "23503"}); got == nil || got.Error() != "principal.not_found" {
		t.Fatalf("foreign-key mapping was %v", got)
	}
}

func TestRepositoryPublicReadUsesActualSQLMethod(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	repo := NewForTesting(db, rand.Reader)
	rows := sqlmock.NewRows([]string{"player_id", "handle_display", "display_name", "avatar_asset_id", "bio", "locale", "visibility", "discoverable"}).AddRow("player.one", "one", "One", "", "", "en", "public", true)
	mock.ExpectQuery("SELECT player_id,handle_display,display_name").WithArgs("player.one").WillReturnRows(rows)
	view, err := repo.GetPublicPlayer(context.Background(), "player.one")
	if err != nil {
		t.Fatal(err)
	}
	if view.PlayerID != "player.one" || view.Visibility != "public" {
		t.Fatal("unexpected public view")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

type failingEntropy struct{}

func (failingEntropy) Read([]byte) (int, error) { return 0, io.ErrUnexpectedEOF }
func TestGuestEntropyFailureRollsBackBeforeMutation(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	repo := NewForTesting(db, failingEntropy{})
	mock.ExpectBegin()
	mock.ExpectExec("SELECT pg_advisory_xact_lock").WithArgs("identity\x00guest.one").WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectQuery("SELECT mutation_type,request_digest,result_ref FROM gridworks.identity_mutation_receipts").WithArgs("guest.one").WillReturnError(sql.ErrNoRows)
	mock.ExpectRollback()
	_, err = repo.IssueGuest(context.Background(), "guest.one", time.Now().UTC())
	if err == nil {
		t.Fatal("expected entropy failure")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}
