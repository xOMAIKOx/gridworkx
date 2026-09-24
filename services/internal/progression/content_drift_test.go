package progression

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestGeneratedProgressionContentMatchesCanonicalFile(t *testing.T) {
	path := filepath.Join("..", "..", "..", "packages", "content", "config", "skills-managers.json")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var canonical struct {
		ProgressionVersion string `json:"progression_version"`
		PlayerSkills       []struct {
			SkillID         string   `json:"skill_id"`
			ActivityClasses []string `json:"activity_classes"`
		} `json:"player_skills"`
		ManagerSkills []string `json:"manager_skills"`
		AntiGrind     struct {
			Trivial    []int64 `json:"trivial_multipliers_bps"`
			Meaningful int64   `json:"meaningful_multiplier_bps"`
		} `json:"anti_grind"`
	}
	if err := json.Unmarshal(data, &canonical); err != nil {
		t.Fatal(err)
	}
	if canonical.ProgressionVersion != SharedProgressionVersion {
		t.Fatalf("version drift: %s", canonical.ProgressionVersion)
	}
	if len(canonical.PlayerSkills) != len(SharedPlayerSkills) || len(canonical.ManagerSkills) != len(SharedManagerSkills) {
		t.Fatal("skill registry drift")
	}
	for _, skill := range canonical.PlayerSkills {
		classes, ok := SharedPlayerSkills[skill.SkillID]
		if !ok || len(classes) != len(skill.ActivityClasses) {
			t.Fatalf("player skill drift: %s", skill.SkillID)
		}
	}
	if canonical.AntiGrind.Meaningful != SharedMeaningfulMultiplier || len(canonical.AntiGrind.Trivial) != len(SharedTrivialMultipliers) {
		t.Fatal("anti-grind drift")
	}
	for i := range canonical.AntiGrind.Trivial {
		if canonical.AntiGrind.Trivial[i] != SharedTrivialMultipliers[i] {
			t.Fatal("anti-grind multiplier drift")
		}
	}
}
