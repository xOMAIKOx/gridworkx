package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/xOMAIKOx/gridworkx/services/internal/progression"
	"strings"
	"testing"
	"time"
)

func TestPlayerProgressionRepositoryReplayConflict(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	repo := NewForTesting(db, nil)
	event := progression.SkillEvent{SourceEventID: "event.one", PlayerID: "player.one", SkillID: "skill.mechanical", ActivityKind: "diagnosis", BaseXP: 100, ActivityClass: progression.ActivityTrivial, RepetitionKey: "same", OccurrenceTime: time.Unix(0, 0), RulesVersion: progression.Version}
	mock.ExpectBegin()
	mock.ExpectExec("SELECT pg_advisory_xact_lock").WithArgs("progression\x00event.one").WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectQuery("SELECT awarded_xp,request_digest FROM gridworks.player_skill_events").WithArgs("event.one").WillReturnError(sql.ErrNoRows)
	mock.ExpectQuery("SELECT player_id,skill_id,cumulative_xp").WithArgs("player.one", "skill.mechanical").WillReturnError(sql.ErrNoRows)
	mock.ExpectExec("INSERT INTO gridworks.player_skills").WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec("INSERT INTO gridworks.player_skill_events").WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()
	award, err := repo.ApplyPlayerSkillEvent(context.Background(), event)
	if err != nil || award.AwardedXP != 100 {
		t.Fatalf("first award=%+v err=%v", award, err)
	}
	requestDigest := digest(strings.Join([]string{event.PlayerID, event.SkillID, event.ActivityKind, fmt.Sprint(event.BaseXP), string(event.ActivityClass), event.RepetitionKey, event.OccurrenceTime.UTC().Format(time.RFC3339Nano), event.RulesVersion}, "\x00"))
	mock.ExpectBegin()
	mock.ExpectExec("SELECT pg_advisory_xact_lock").WithArgs("progression\x00event.one").WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectQuery("SELECT awarded_xp,request_digest FROM gridworks.player_skill_events").WithArgs("event.one").WillReturnRows(sqlmock.NewRows([]string{"awarded_xp", "request_digest"}).AddRow(100, requestDigest))
	mock.ExpectQuery("SELECT cumulative_xp,proficiency_bps,repetition_count FROM gridworks.player_skills").WithArgs("player.one", "skill.mechanical").WillReturnRows(sqlmock.NewRows([]string{"cumulative_xp", "proficiency_bps", "repetition_count"}).AddRow(100, 1000, 1))
	mock.ExpectCommit()
	replay, err := repo.ApplyPlayerSkillEvent(context.Background(), event)
	if err != nil || !replay.Replay || replay.CumulativeXP != 100 || replay.ProficiencyBPS != 1000 {
		t.Fatalf("replay=%+v err=%v", replay, err)
	}
	mock.ExpectBegin()
	mock.ExpectExec("SELECT pg_advisory_xact_lock").WithArgs("progression\x00event.one").WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectQuery("SELECT awarded_xp,request_digest FROM gridworks.player_skill_events").WithArgs("event.one").WillReturnRows(sqlmock.NewRows([]string{"awarded_xp", "request_digest"}).AddRow(100, "0000000000000000000000000000000000000000000000000000000000000000"))
	mock.ExpectRollback()
	if _, err := repo.ApplyPlayerSkillEvent(context.Background(), event); err == nil {
		t.Fatal("changed digest should conflict")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestManagerProgressionRepositoryLocksAndPersistsSkillHistory(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	repo := NewForTesting(db, nil)
	now := time.Unix(0, 0)
	event := progression.ManagerProgressionEvent{SourceEventID: "manager.event", ManagerID: "manager.one", ActivityKind: "diagnosis", SkillID: "manager_skill.technical", AwardedXP: 100, OccurrenceTime: now, RulesVersion: progression.Version}
	mock.ExpectBegin()
	mock.ExpectExec("SELECT pg_advisory_xact_lock").WithArgs("progression\x00manager.progression.manager.event").WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectQuery("SELECT mutation_type,request_digest,result_ref FROM gridworks.manager_mutation_receipts").WithArgs("manager.progression.manager.event").WillReturnError(sql.ErrNoRows)
	mock.ExpectQuery("SELECT manager_id,rarity,total_xp,level FROM gridworks.managers").WithArgs("manager.one").WillReturnRows(sqlmock.NewRows([]string{"manager_id", "rarity", "total_xp", "level"}).AddRow("manager.one", "Gold", 0, 1))
	mock.ExpectQuery("SELECT skill_id,proficiency_bps,potential_bps FROM gridworks.manager_skills").WithArgs("manager.one").WillReturnRows(sqlmock.NewRows([]string{"skill_id", "proficiency_bps", "potential_bps"}).AddRow("manager_skill.technical", 0, 8500))
	mock.ExpectExec("UPDATE gridworks.managers").WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec("UPDATE gridworks.manager_skills").WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec("INSERT INTO gridworks.manager_progression_events").WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec("INSERT INTO gridworks.manager_mutation_receipts").WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()
	result, err := repo.ApplyManagerProgression(context.Background(), event)
	if err != nil || result.TotalXP != 100 {
		t.Fatalf("result=%+v err=%v", result, err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestManagerProgressionReplayAndOccurrenceConflict(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	repo := NewForTesting(db, nil)
	now := time.Unix(0, 0)
	event := progression.ManagerProgressionEvent{SourceEventID: "manager.replay", ManagerID: "manager.one", ActivityKind: "diagnosis", SkillID: "manager_skill.technical", AwardedXP: 100, OccurrenceTime: now, RulesVersion: progression.Version}
	digestValue := digest(strings.Join([]string{event.ManagerID, event.ActivityKind, event.SkillID, fmt.Sprint(event.AwardedXP), event.OccurrenceTime.UTC().Format(time.RFC3339Nano), event.RulesVersion}, "\x00"))
	mock.ExpectBegin()
	mock.ExpectExec("SELECT pg_advisory_xact_lock").WithArgs("progression\x00manager.progression.manager.replay").WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectQuery("SELECT mutation_type,request_digest,result_ref FROM gridworks.manager_mutation_receipts").WithArgs("manager.progression.manager.replay").WillReturnRows(sqlmock.NewRows([]string{"mutation_type", "request_digest", "result_ref"}).AddRow("manager.progression", digestValue, "manager.one"))
	mock.ExpectQuery("SELECT total_xp,level FROM gridworks.managers").WithArgs("manager.one").WillReturnRows(sqlmock.NewRows([]string{"total_xp", "level"}).AddRow(100, 1))
	mock.ExpectCommit()
	result, err := repo.ApplyManagerProgression(context.Background(), event)
	if err != nil || !result.Replay || result.TotalXP != 100 {
		t.Fatalf("replay=%+v err=%v", result, err)
	}
	event.OccurrenceTime = event.OccurrenceTime.Add(time.Hour)
	mock.ExpectBegin()
	mock.ExpectExec("SELECT pg_advisory_xact_lock").WithArgs("progression\x00manager.progression.manager.replay").WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectQuery("SELECT mutation_type,request_digest,result_ref FROM gridworks.manager_mutation_receipts").WithArgs("manager.progression.manager.replay").WillReturnRows(sqlmock.NewRows([]string{"mutation_type", "request_digest", "result_ref"}).AddRow("manager.progression", digestValue, "manager.one"))
	mock.ExpectRollback()
	if _, err := repo.ApplyManagerProgression(context.Background(), event); err == nil {
		t.Fatal("occurrence change should conflict")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestEmploymentCloseRejectsActiveAssignmentBeforeUpdate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	repo := NewForTesting(db, nil)
	when := time.Unix(0, 0)
	mock.ExpectBegin()
	mock.ExpectExec("SELECT pg_advisory_xact_lock").WithArgs("progression\x00employment.close.employment.one").WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectQuery("SELECT mutation_type,request_digest,result_ref FROM gridworks.manager_mutation_receipts").WithArgs("employment.close.employment.one").WillReturnError(sql.ErrNoRows)
	mock.ExpectQuery("SELECT manager_id FROM gridworks.manager_employment_history").WithArgs("employment.one").WillReturnRows(sqlmock.NewRows([]string{"manager_id"}).AddRow("manager.one"))
	mock.ExpectQuery("SELECT manager_id FROM gridworks.managers WHERE manager_id=\\$1 FOR UPDATE").WithArgs("manager.one").WillReturnRows(sqlmock.NewRows([]string{"manager_id"}).AddRow("manager.one"))
	mock.ExpectQuery("SELECT count\\(\\*\\) FROM gridworks.manager_facility_assignments").WithArgs("manager.one").WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
	mock.ExpectRollback()
	if err := repo.CloseManagerEmployment(context.Background(), "employment.one", "source.one", when); err != progression.ErrActiveAssignment {
		t.Fatalf("error=%v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestEmploymentMutationReplayConflictAndCloseLifecycle(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	repo := NewForTesting(db, nil)
	when := time.Unix(0, 0)
	openKey := "employment.open.employment.one"
	openDigest := digest(strings.Join([]string{"manager.one", "company.one", "operations", "source.one", when.UTC().Format(time.RFC3339Nano)}, "\x00"))
	mock.ExpectBegin()
	mock.ExpectExec("SELECT pg_advisory_xact_lock").WithArgs("progression\x00" + openKey).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectQuery("SELECT mutation_type,request_digest,result_ref FROM gridworks.manager_mutation_receipts").WithArgs(openKey).WillReturnError(sql.ErrNoRows)
	mock.ExpectQuery("SELECT manager_id FROM gridworks.managers").WithArgs("manager.one").WillReturnRows(sqlmock.NewRows([]string{"manager_id"}).AddRow("manager.one"))
	mock.ExpectQuery("SELECT count\\(\\*\\) FROM gridworks.manager_employment_history").WithArgs("manager.one").WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
	mock.ExpectExec("INSERT INTO gridworks.manager_employment_history").WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec("INSERT INTO gridworks.manager_mutation_receipts").WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()
	if err := repo.OpenManagerEmployment(context.Background(), "employment.one", "manager.one", "company.one", "operations", "source.one", when); err != nil {
		t.Fatal(err)
	}
	mock.ExpectBegin()
	mock.ExpectExec("SELECT pg_advisory_xact_lock").WithArgs("progression\x00" + openKey).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectQuery("SELECT mutation_type,request_digest,result_ref FROM gridworks.manager_mutation_receipts").WithArgs(openKey).WillReturnRows(sqlmock.NewRows([]string{"mutation_type", "request_digest", "result_ref"}).AddRow("employment.open", openDigest, "employment.one"))
	mock.ExpectCommit()
	if err := repo.OpenManagerEmployment(context.Background(), "employment.one", "manager.one", "company.one", "operations", "source.one", when); err != nil {
		t.Fatal(err)
	}
	closeKey := "employment.close.employment.one"
	closeDigest := digest(strings.Join([]string{"employment.one", "close.source", when.UTC().Format(time.RFC3339Nano)}, "\x00"))
	mock.ExpectBegin()
	mock.ExpectExec("SELECT pg_advisory_xact_lock").WithArgs("progression\x00" + closeKey).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectQuery("SELECT mutation_type,request_digest,result_ref FROM gridworks.manager_mutation_receipts").WithArgs(closeKey).WillReturnError(sql.ErrNoRows)
	mock.ExpectQuery("SELECT manager_id FROM gridworks.manager_employment_history").WithArgs("employment.one").WillReturnRows(sqlmock.NewRows([]string{"manager_id"}).AddRow("manager.one"))
	mock.ExpectQuery("SELECT manager_id FROM gridworks.managers").WithArgs("manager.one").WillReturnRows(sqlmock.NewRows([]string{"manager_id"}).AddRow("manager.one"))
	mock.ExpectQuery("SELECT count\\(\\*\\) FROM gridworks.manager_facility_assignments").WithArgs("manager.one").WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
	mock.ExpectExec("UPDATE gridworks.manager_employment_history").WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec("INSERT INTO gridworks.manager_mutation_receipts").WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()
	if err := repo.CloseManagerEmployment(context.Background(), "employment.one", "close.source", when); err != nil {
		t.Fatal(err)
	}
	mock.ExpectBegin()
	mock.ExpectExec("SELECT pg_advisory_xact_lock").WithArgs("progression\x00" + closeKey).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectQuery("SELECT mutation_type,request_digest,result_ref FROM gridworks.manager_mutation_receipts").WithArgs(closeKey).WillReturnRows(sqlmock.NewRows([]string{"mutation_type", "request_digest", "result_ref"}).AddRow("employment.close", closeDigest, "employment.one"))
	mock.ExpectCommit()
	if err := repo.CloseManagerEmployment(context.Background(), "employment.one", "close.source", when); err != nil {
		t.Fatal(err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}
