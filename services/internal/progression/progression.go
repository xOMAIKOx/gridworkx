package progression

import (
	"errors"
	"sort"
)

const Version = "progression-0.1.0"

var ErrInvalidView = errors.New("invalid progression view")

type PlayerSkillView struct {
	PlayerID           string `json:"player_id"`
	SkillID            string `json:"skill_id"`
	CumulativeXP       int64  `json:"cumulative_xp"`
	ProficiencyBPS     int    `json:"proficiency_bps"`
	ProgressionVersion string `json:"progression_version"`
}

type ManagerSkillView struct {
	SkillID        string `json:"skill_id"`
	ProficiencyBPS int    `json:"proficiency_bps"`
	PotentialBPS   int    `json:"potential_bps"`
}

type ManagerView struct {
	ManagerID        string             `json:"manager_id"`
	DisplayName      string             `json:"display_name"`
	Rarity           string             `json:"rarity"`
	SpecializationID string             `json:"specialization_id"`
	TotalXP          int64              `json:"total_xp"`
	Level            int                `json:"level"`
	Skills           []ManagerSkillView `json:"skills"`
	Traits           []string           `json:"traits"`
	WorkloadBPS      int                `json:"workload_bps"`
	FatigueBPS       int                `json:"fatigue_bps"`
	MoraleBPS        int                `json:"morale_bps"`
	Status           string             `json:"status"`
	ActiveEmployment bool               `json:"active_employment"`
	ActiveFacilityID string             `json:"active_facility_id,omitempty"`
	PotentialVisible bool               `json:"potential_visible"`
}

func ValidatePlayerSkillView(view PlayerSkillView) error {
	if view.PlayerID == "" || view.SkillID == "" || view.CumulativeXP < 0 || view.ProficiencyBPS < 0 || view.ProficiencyBPS > 10000 || view.ProgressionVersion == "" {
		return ErrInvalidView
	}
	return nil
}

func ValidateManagerView(view ManagerView) error {
	if view.ManagerID == "" || view.DisplayName == "" || view.TotalXP < 0 || view.Level < 1 || view.WorkloadBPS < 0 || view.WorkloadBPS > 10000 || view.FatigueBPS < 0 || view.FatigueBPS > 10000 || view.MoraleBPS < 0 || view.MoraleBPS > 10000 {
		return ErrInvalidView
	}
	for _, skill := range view.Skills {
		if skill.ProficiencyBPS < 0 || skill.ProficiencyBPS > skill.PotentialBPS || skill.PotentialBPS > 10000 {
			return ErrInvalidView
		}
	}
	return nil
}

func SortSkills(skills []PlayerSkillView) {
	sort.Slice(skills, func(i, j int) bool { return skills[i].SkillID < skills[j].SkillID })
}
func SortManagers(managers []ManagerView) {
	sort.Slice(managers, func(i, j int) bool { return managers[i].ManagerID < managers[j].ManagerID })
}
