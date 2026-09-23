package postgres

import (
	"encoding/json"
	"strings"
	"testing"

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
