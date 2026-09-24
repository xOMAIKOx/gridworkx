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
	ErrUnknownActivity    = errors.New("progression activity is not applicable")
	ErrOverflow           = errors.New("progression arithmetic overflow")
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
type ManagerState struct {
	ManagerID string
	Rarity    string
	TotalXP   int64
	Level     int
	Skills    []ManagerSkillView
}

type SkillAward struct {
	AwardedXP       int64
	CumulativeXP    int64
	ProficiencyBPS  int
	RepetitionCount int
	Replay          bool
}

var PlayerSkillIDs = SharedPlayerSkills

func ApplySkillEvent(state *SkillState, event SkillEvent, duplicate bool) (SkillAward, error) {
	if state.PlayerID != event.PlayerID || state.SkillID != event.SkillID || event.BaseXP < 0 || event.SourceEventID == "" || event.RepetitionKey == "" || event.RulesVersion != Version {
		return SkillAward{}, ErrInvalidEvent
	}
	activityClasses, ok := PlayerSkillIDs[event.SkillID]
	if !ok {
		return SkillAward{}, ErrUnknownSkill
	}
	if event.ActivityClass != ActivityTrivial && event.ActivityClass != ActivityMeaningful {
		return SkillAward{}, ErrInvalidEvent
	}
	if !contains(activityClasses, event.ActivityKind) {
		return SkillAward{}, ErrUnknownActivity
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
	if event.BaseXP > (1<<63-1)/multiplier {
		return SkillAward{}, ErrOverflow
	}
	awarded := event.BaseXP * multiplier / 10000
	if awarded > (1<<63-1)-state.CumulativeXP {
		return SkillAward{}, ErrOverflow
	}
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

func contains(values []string, wanted string) bool {
	for _, value := range values {
		if value == wanted {
			return true
		}
	}
	return false
}

type ManagerProgressionEvent struct {
	SourceEventID  string
	ManagerID      string
	ActivityKind   string
	SkillID        string
	AwardedXP      int64
	OccurrenceTime time.Time
	RulesVersion   string
}
type ManagerProgressionResult struct {
	TotalXP int64
	Level   int
	Replay  bool
}

func ApplyManagerProgression(state *ManagerState, event ManagerProgressionEvent, duplicate bool) (ManagerProgressionResult, error) {
	if state.ManagerID != event.ManagerID || event.SourceEventID == "" || event.RulesVersion != Version || event.AwardedXP < 0 {
		return ManagerProgressionResult{}, ErrInvalidEvent
	}
	if duplicate {
		return ManagerProgressionResult{TotalXP: state.TotalXP, Level: state.Level, Replay: true}, nil
	}
	if event.AwardedXP > (1<<63-1)-state.TotalXP {
		return ManagerProgressionResult{}, ErrOverflow
	}
	state.TotalXP += event.AwardedXP
	state.Level = int(state.TotalXP/1000 + 1)
	if event.SkillID != "" {
		if !contains(SharedManagerSkills, event.SkillID) {
			return ManagerProgressionResult{}, ErrUnknownSkill
		}
		cap := SharedRarityPotential[string(state.Rarity)]
		for index := range state.Skills {
			if state.Skills[index].SkillID == event.SkillID {
				state.Skills[index].ProficiencyBPS = minInt(state.Skills[index].ProficiencyBPS+int(event.AwardedXP), minInt(state.Skills[index].PotentialBPS, cap))
			}
		}
	}
	return ManagerProgressionResult{TotalXP: state.TotalXP, Level: state.Level}, nil
}
func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
