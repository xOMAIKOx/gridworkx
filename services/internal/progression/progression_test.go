package progression

import "testing"

func TestViewsRejectUnboundedState(t *testing.T) {
	if ValidatePlayerSkillView(PlayerSkillView{PlayerID: "p", SkillID: "skill.mechanical", CumulativeXP: -1, ProficiencyBPS: 0, ProgressionVersion: Version}) == nil {
		t.Fatal("negative xp accepted")
	}
	if ValidateManagerView(ManagerView{ManagerID: "m", DisplayName: "M", Level: 1, Skills: []ManagerSkillView{{SkillID: "manager_skill.technical", ProficiencyBPS: 9000, PotentialBPS: 8000}}}) == nil {
		t.Fatal("skill above potential accepted")
	}
}
