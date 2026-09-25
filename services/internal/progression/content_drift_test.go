package progression

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
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
			SkillID          string   `json:"skill_id"`
			LabelKey         string   `json:"label_key"`
			ActivityClasses  []string `json:"activity_classes"`
			ProgressionCurve string   `json:"progression_curve"`
			RulesVersion     string   `json:"rules_version"`
		} `json:"player_skills"`
		ManagerSkills []string `json:"manager_skills"`
		RarityPolicy  map[string]struct {
			PotentialBPS  int `json:"potential_bps"`
			TraitCapacity int `json:"trait_capacity"`
		} `json:"rarity_policy"`
		Traits []struct {
			TraitID string         `json:"trait_id"`
			Effects map[string]int `json:"effects"`
		} `json:"traits"`
		AntiGrind struct {
			Version    string  `json:"version"`
			Trivial    []int64 `json:"trivial_multipliers_bps"`
			Meaningful int64   `json:"meaningful_multiplier_bps"`
		} `json:"anti_grind"`
		ManagerProgression struct {
			LevelXP       int64 `json:"level_xp_per_level"`
			SkillBPSPerXP int64 `json:"skill_bps_per_xp"`
		} `json:"manager_progression"`
	}
	if err := json.Unmarshal(data, &canonical); err != nil {
		t.Fatal(err)
	}
	if canonical.ProgressionVersion != SharedProgressionVersion || canonical.AntiGrind.Version != SharedProgressionVersion {
		t.Fatalf("version drift: %+v", canonical)
	}
	if len(canonical.PlayerSkills) != len(SharedPlayerSkills) || len(canonical.ManagerSkills) != len(SharedManagerSkills) {
		t.Fatal("skill registry drift")
	}
	for _, skill := range canonical.PlayerSkills {
		metadata, ok := SharedPlayerSkillMetadataByID[skill.SkillID]
		if !ok || metadata.LabelKey != skill.LabelKey || !reflect.DeepEqual(metadata.ActivityClasses, skill.ActivityClasses) || metadata.ProgressionCurve != skill.ProgressionCurve || metadata.RulesVersion != skill.RulesVersion {
			t.Fatalf("player skill metadata drift: %s", skill.SkillID)
		}
	}
	if !reflect.DeepEqual(canonical.ManagerSkills, SharedManagerSkills) || canonical.ManagerProgression.LevelXP != SharedManagerProgression.LevelXP || canonical.ManagerProgression.SkillBPSPerXP != SharedManagerProgression.SkillBPSPerXP {
		t.Fatal("manager content drift")
	}
	if len(canonical.RarityPolicy) != len(SharedRarityPotential) || len(SharedRarityPolicyKeys()) != len(canonical.RarityPolicy) {
		t.Fatal("rarity key-set drift")
	}
	for rarity, policy := range canonical.RarityPolicy {
		if SharedRarityPotential[rarity] != policy.PotentialBPS || SharedRarityTraitCapacity[rarity] != policy.TraitCapacity {
			t.Fatalf("rarity drift: %s", rarity)
		}
	}
	if len(canonical.Traits) != len(SharedTraitEffects) {
		t.Fatal("trait key-set drift")
	}
	for _, trait := range canonical.Traits {
		if !reflect.DeepEqual(SharedTraitEffects[trait.TraitID], trait.Effects) {
			t.Fatalf("trait drift: %s", trait.TraitID)
		}
	}
	if canonical.AntiGrind.Meaningful != SharedMeaningfulMultiplier || !reflect.DeepEqual(canonical.AntiGrind.Trivial, SharedTrivialMultipliers) {
		t.Fatal("anti-grind drift")
	}
}

func SharedRarityPolicyKeys() []string {
	keys := make([]string, 0, len(SharedRarityPotential))
	for key := range SharedRarityPotential {
		keys = append(keys, key)
	}
	return keys
}
