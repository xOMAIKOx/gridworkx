package postgres

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/xOMAIKOx/gridworkx/services/api/internal/api"
	"github.com/xOMAIKOx/gridworkx/services/internal/progression"
)

func TestProgressionRepositoryPlayerSkillReadUsesActualSQL(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	repo := NewForTesting(db, nil)
	mock.ExpectQuery("SELECT player_id,skill_id,cumulative_xp,proficiency_bps,progression_version FROM gridworks.player_skills").WithArgs("player.one").WillReturnRows(sqlmock.NewRows([]string{"player_id", "skill_id", "cumulative_xp", "proficiency_bps", "progression_version"}).AddRow("player.one", "skill.mechanical", 175, 1750, progression.Version))
	views, err := repo.GetPlayerSkills(context.Background(), api.Principal{PlayerID: "player.one"})
	if err != nil {
		t.Fatal(err)
	}
	if len(views) != 1 || views[0].CumulativeXP != 175 {
		t.Fatalf("unexpected views: %+v", views)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestProgressionRepositoryCompanyReadEnforcesOwner(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	repo := NewForTesting(db, nil)
	mock.ExpectQuery("SELECT EXISTS").WithArgs("company.one", "player.one").WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))
	if _, err := repo.GetCompanyManagers(context.Background(), api.Principal{PlayerID: "player.one"}, "company.one"); err != api.ErrForbidden {
		t.Fatalf("non-owner error=%v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestProgressionRepositoryCompanyRosterScansCanonicalViews(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	repo := NewForTesting(db, nil)
	mock.ExpectQuery("SELECT EXISTS").WithArgs("company.one", "player.one").WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))
	mock.ExpectQuery("SELECT m.manager_id,m.display_name,m.rarity").WithArgs("company.one").WillReturnRows(sqlmock.NewRows([]string{"manager_id", "display_name", "rarity", "specialization_id", "total_xp", "level", "workload_bps", "fatigue_bps", "morale_bps", "status", "jsonb_agg", "jsonb_agg", "facility_id"}).AddRow("manager.one", "Manager", "Gold", "specialization.aggregate_diagnostics", 0, 1, 0, 0, 10000, "Active", []byte(`[{"skill_id":"manager_skill.technical","proficiency_bps":8000,"potential_bps":8500}]`), []byte(`["trait.crisis_specialist"]`), ""))
	views, err := repo.GetCompanyManagers(context.Background(), api.Principal{PlayerID: "player.one"}, "company.one")
	if err != nil {
		t.Fatal(err)
	}
	if len(views) != 1 || views[0].Rarity != "Gold" || len(views[0].Skills) != 1 {
		t.Fatalf("unexpected manager views: %+v", views)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}
