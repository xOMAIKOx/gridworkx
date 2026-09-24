// Code generated from packages/content/config/skills-managers.json; DO NOT EDIT.
package progression

const SharedProgressionVersion = "progression-0.1.0"

var SharedPlayerSkills = map[string][]string{
	"skill.mechanical": {"diagnosis", "repair", "production"}, "skill.electrical": {"diagnosis", "repair", "production"}, "skill.process_engineering": {"diagnosis", "repair", "production"}, "skill.agriculture": {"diagnosis", "repair", "production"}, "skill.mining": {"diagnosis", "repair", "production"}, "skill.energy": {"diagnosis", "repair", "production"}, "skill.water": {"diagnosis", "repair", "production"}, "skill.logistics": {"diagnosis", "repair", "production"}, "skill.construction": {"diagnosis", "repair", "production"}, "skill.commerce": {"diagnosis", "repair", "production"}, "skill.finance": {"diagnosis", "repair", "production"}, "skill.management": {"diagnosis", "repair", "production"},
}
var SharedManagerSkills = []string{"manager_skill.operations", "manager_skill.technical", "manager_skill.maintenance", "manager_skill.safety", "manager_skill.leadership", "manager_skill.logistics", "manager_skill.energy_efficiency", "manager_skill.crisis_response", "manager_skill.mentoring"}
var SharedRarityPotential = map[string]int{"Bronze": 5000, "Silver": 6500, "Gold": 8500, "Platinum": 10000}
var SharedRarityTraitCapacity = map[string]int{"Bronze": 1, "Silver": 2, "Gold": 3, "Platinum": 4}
var SharedTraitEffects = map[string]map[string]int{"trait.aggressive_operator": {"diagnostic_capability_bps": 100, "fatigue_contribution_bps": 300}, "trait.maintenance_first": {"diagnostic_capability_bps": 250, "fatigue_contribution_bps": -250}, "trait.cost_cutter": {"diagnostic_capability_bps": 0, "workload_contribution_bps": -300}, "trait.mentor": {"diagnostic_capability_bps": 150, "morale_contribution_bps": 250}, "trait.crisis_specialist": {"diagnostic_capability_bps": 500, "workload_contribution_bps": 400}}
var SharedTrivialMultipliers = []int64{10000, 5000, 2500, 0}

const SharedMeaningfulMultiplier int64 = 10000
