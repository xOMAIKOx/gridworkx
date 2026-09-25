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
	ErrVersionMismatch    = errors.New("progression version mismatch")
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
	ManagerID          string
	Rarity             string
	TotalXP            int64
	Level              int
	Skills             []ManagerSkillView
	ProgressionVersion string
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
	if state.PlayerID != event.PlayerID || state.SkillID != event.SkillID || event.BaseXP < 0 || event.SourceEventID == "" || event.RepetitionKey == "" {
		return SkillAward{}, ErrInvalidEvent
	}
	if event.RulesVersion != Version || (state.ProgressionVersion != "" && state.ProgressionVersion != Version) {
		return SkillAward{}, ErrVersionMismatch
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
	multiplier := SharedMeaningfulMultiplier
	if event.ActivityClass == ActivityTrivial {
		if count <= len(SharedTrivialMultipliers) {
			multiplier = SharedTrivialMultipliers[count-1]
		} else {
			multiplier = 0
		}
	}
	if multiplier == 0 {
		state.ProgressionVersion = Version
		state.RepetitionKey = event.RepetitionKey
		state.RepetitionCount = count
		return SkillAward{AwardedXP: 0, CumulativeXP: state.CumulativeXP, ProficiencyBPS: state.ProficiencyBPS, RepetitionCount: count}, nil
	}
	if event.BaseXP > (1<<63-1)/multiplier {
		return SkillAward{}, ErrOverflow
	}
	awarded := event.BaseXP * multiplier / 10000
	if awarded > (1<<63-1)-state.CumulativeXP {
		return SkillAward{}, ErrOverflow
	}
	state.CumulativeXP += awarded
	if state.CumulativeXP >= 1000 {
		state.ProficiencyBPS = 10000
	} else {
		state.ProficiencyBPS = int(state.CumulativeXP * 10)
	}
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
	if state.ManagerID != event.ManagerID || event.SourceEventID == "" || event.AwardedXP < 0 {
		return ManagerProgressionResult{}, ErrInvalidEvent
	}
	if event.RulesVersion != Version || (state.ProgressionVersion != "" && state.ProgressionVersion != Version) {
		return ManagerProgressionResult{}, ErrVersionMismatch
	}
	if _, ok := SharedRarityPotential[state.Rarity]; !ok {
		return ManagerProgressionResult{}, ErrInvalidEvent
	}
	if duplicate {
		return ManagerProgressionResult{TotalXP: state.TotalXP, Level: state.Level, Replay: true}, nil
	}
	var target *ManagerSkillView
	if event.SkillID != "" {
		if !contains(SharedManagerSkills, event.SkillID) {
			return ManagerProgressionResult{}, ErrUnknownSkill
		}
		for index := range state.Skills {
			if state.Skills[index].SkillID == event.SkillID {
				target = &state.Skills[index]
				break
			}
		}
		if target == nil {
			return ManagerProgressionResult{}, ErrInvalidEvent
		}
	}
	if event.AwardedXP > (1<<63-1)-state.TotalXP {
		return ManagerProgressionResult{}, ErrOverflow
	}
	newTotal := state.TotalXP + event.AwardedXP
	newLevel := int(newTotal/SharedManagerProgression.LevelXP + 1)
	var newSkill int
	if target != nil {
		if event.AwardedXP > (1<<63-1)/SharedManagerProgression.SkillBPSPerXP {
			return ManagerProgressionResult{}, ErrOverflow
		}
		newSkill = target.ProficiencyBPS + int(event.AwardedXP*SharedManagerProgression.SkillBPSPerXP)
		if newSkill > target.PotentialBPS {
			newSkill = target.PotentialBPS
		}
	}
	state.TotalXP = newTotal
	state.Level = newLevel
	if target != nil {
		target.ProficiencyBPS = newSkill
	}
	return ManagerProgressionResult{TotalXP: state.TotalXP, Level: state.Level}, nil
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
