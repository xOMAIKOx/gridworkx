package progression

import (
	"testing"
	"time"
)

func TestTrustedSkillEventFailsClosedAndReplaysZero(t *testing.T) {
	state := SkillState{PlayerID: "player.one", SkillID: "skill.mechanical", ProgressionVersion: Version}
	event := SkillEvent{SourceEventID: "event.one", PlayerID: "player.one", SkillID: "skill.mechanical", ActivityKind: "diagnosis", BaseXP: 100, ActivityClass: ActivityTrivial, RepetitionKey: "same", OccurrenceTime: time.Unix(0, 0), RulesVersion: Version}
	award, err := ApplySkillEvent(&state, event, false)
	if err != nil || award.AwardedXP != 100 {
		t.Fatalf("award=%+v err=%v", award, err)
	}
	replay, err := ApplySkillEvent(&state, event, true)
	if err != nil || replay.AwardedXP != 0 || !replay.Replay {
		t.Fatalf("replay=%+v err=%v", replay, err)
	}
	event.ActivityKind = "unknown"
	if _, err := ApplySkillEvent(&state, event, false); err != ErrUnknownActivity {
		t.Fatalf("unknown activity err=%v", err)
	}
	event.ActivityKind = "diagnosis"
	event.BaseXP = 1<<63 - 1
	if _, err := ApplySkillEvent(&state, event, false); err != ErrOverflow {
		t.Fatalf("overflow err=%v", err)
	}
}

func TestManagerProgressionIsCappedAndReplaySafe(t *testing.T) {
	state := ManagerState{ManagerID: "manager.one", Rarity: "Gold", TotalXP: 0, Level: 1, Skills: []ManagerSkillView{{SkillID: "manager_skill.technical", PotentialBPS: 8500}}}
	event := ManagerProgressionEvent{SourceEventID: "manager.event.one", ManagerID: "manager.one", SkillID: "manager_skill.technical", ActivityKind: "diagnosis", AwardedXP: 9000, OccurrenceTime: time.Unix(0, 0), RulesVersion: Version}
	result, err := ApplyManagerProgression(&state, event, false)
	if err != nil || result.TotalXP != 9000 || state.Skills[0].ProficiencyBPS != 8500 {
		t.Fatalf("result=%+v state=%+v err=%v", result, state, err)
	}
	replay, err := ApplyManagerProgression(&state, event, true)
	if err != nil || !replay.Replay {
		t.Fatalf("replay=%+v err=%v", replay, err)
	}
}

func TestTrustedProgressionZeroMultiplierVersionAndNoMutation(t *testing.T) {
	state := SkillState{PlayerID: "player.one", SkillID: "skill.mechanical", ProgressionVersion: Version, CumulativeXP: 100, ProficiencyBPS: 1000, RepetitionKey: "same", RepetitionCount: 3}
	event := SkillEvent{SourceEventID: "event.zero", PlayerID: "player.one", SkillID: "skill.mechanical", ActivityKind: "diagnosis", BaseXP: 100, ActivityClass: ActivityTrivial, RepetitionKey: "same", OccurrenceTime: time.Unix(0, 0), RulesVersion: Version}
	before := state
	award, err := ApplySkillEvent(&state, event, false)
	if err != nil || award.AwardedXP != 0 || state.CumulativeXP != before.CumulativeXP {
		t.Fatalf("zero multiplier award=%+v state=%+v err=%v", award, state, err)
	}
	state = before
	event.RulesVersion = "progression-old"
	if _, err := ApplySkillEvent(&state, event, false); err != ErrVersionMismatch || state != before {
		t.Fatalf("version failure state=%+v err=%v", state, err)
	}
}

func TestManagerProgressionRejectsInvalidSkillWithoutMutation(t *testing.T) {
	state := ManagerState{ManagerID: "manager.one", Rarity: "Gold", TotalXP: 10, Level: 1, Skills: []ManagerSkillView{{SkillID: "manager_skill.technical", PotentialBPS: 8500}}}
	before := state
	_, err := ApplyManagerProgression(&state, ManagerProgressionEvent{SourceEventID: "manager.bad", ManagerID: "manager.one", SkillID: "manager_skill.missing", AwardedXP: 10, RulesVersion: Version}, false)
	if err != ErrUnknownSkill || state.TotalXP != before.TotalXP || state.Level != before.Level {
		t.Fatalf("invalid manager skill mutated state=%+v err=%v", state, err)
	}
}
