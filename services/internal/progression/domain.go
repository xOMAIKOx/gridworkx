package progression

import (
	"errors"
	"time"
)

var (
	ErrDuplicateSource    = errors.New("progression source event already applied")
	ErrUnknownSkill       = errors.New("progression skill is unknown")
	ErrInvalidEvent       = errors.New("progression event is invalid")
	ErrActiveEmployment   = errors.New("manager already has an active employer")
	ErrNoActiveEmployment = errors.New("manager has no active employment")
)

type ActivityClass string

const (
	ActivityTrivial    ActivityClass = "trivial"
	ActivityMeaningful ActivityClass = "meaningful"
)

type SkillEvent struct {
	SourceEventID  string
	PlayerID       string
	SkillID        string
	ActivityKind   string
	BaseXP         int64
	ActivityClass  ActivityClass
	RepetitionKey  string
	OccurrenceTime time.Time
	RulesVersion   string
}
type SkillState struct {
	PlayerID           string
	SkillID            string
	CumulativeXP       int64
	ProficiencyBPS     int
	ProgressionVersion string
	RepetitionKey      string
	RepetitionCount    int
}
type SkillAward struct {
	AwardedXP       int64
	CumulativeXP    int64
	ProficiencyBPS  int
	RepetitionCount int
	Replay          bool
}

var PlayerSkillIDs = map[string]struct{}{
	"skill.mechanical": {}, "skill.electrical": {}, "skill.process_engineering": {}, "skill.agriculture": {}, "skill.mining": {}, "skill.energy": {}, "skill.water": {}, "skill.logistics": {}, "skill.construction": {}, "skill.commerce": {}, "skill.finance": {}, "skill.management": {},
}

func ApplySkillEvent(state *SkillState, event SkillEvent, duplicate bool) (SkillAward, error) {
	if state.PlayerID != event.PlayerID || state.SkillID != event.SkillID || event.BaseXP < 0 || event.SourceEventID == "" || event.RepetitionKey == "" || event.RulesVersion != Version {
		return SkillAward{}, ErrInvalidEvent
	}
	if _, ok := PlayerSkillIDs[event.SkillID]; !ok {
		return SkillAward{}, ErrUnknownSkill
	}
	if duplicate {
		return SkillAward{CumulativeXP: state.CumulativeXP, ProficiencyBPS: state.ProficiencyBPS, RepetitionCount: state.RepetitionCount, Replay: true}, nil
	}
	count := 1
	if state.RepetitionKey == event.RepetitionKey {
		count = state.RepetitionCount + 1
	}
	multiplier := int64(10000)
	if event.ActivityClass == ActivityTrivial {
		multipliers := []int64{10000, 5000, 2500, 0}
		if count <= len(multipliers) {
			multiplier = multipliers[count-1]
		} else {
			multiplier = 0
		}
	}
	awarded := event.BaseXP * multiplier / 10000
	state.CumulativeXP += awarded
	state.ProficiencyBPS = int(minInt64(state.CumulativeXP*10, 10000))
	state.ProgressionVersion = Version
	state.RepetitionKey = event.RepetitionKey
	state.RepetitionCount = count
	return SkillAward{AwardedXP: awarded, CumulativeXP: state.CumulativeXP, ProficiencyBPS: state.ProficiencyBPS, RepetitionCount: count}, nil
}
func minInt64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}
