package api

import (
	"context"
	"github.com/xOMAIKOx/gridworkx/services/internal/progression"
	"net/http"
	"net/http/httptest"
	"testing"
)

func (f *fakeRepo) GetPlayerSkills(context.Context, Principal) ([]progression.PlayerSkillView, error) {
	return []progression.PlayerSkillView{{PlayerID: "player.one", SkillID: "skill.mechanical", CumulativeXP: 175, ProficiencyBPS: 1750, ProgressionVersion: progression.Version}}, nil
}
func (f *fakeRepo) GetCompanyManagers(context.Context, Principal, string) ([]progression.ManagerView, error) {
	return []progression.ManagerView{{ManagerID: "manager.one", DisplayName: "Manager", Rarity: "Gold", SpecializationID: "specialization.aggregate_diagnostics", Level: 1, Skills: []progression.ManagerSkillView{{SkillID: "manager_skill.technical", ProficiencyBPS: 8000, PotentialBPS: 8500}}, Traits: []string{"trait.crisis_specialist"}, MoraleBPS: 10000, ActiveEmployment: true, PotentialVisible: true}}, nil
}
func (f *fakeRepo) GetCompanyManager(ctx context.Context, p Principal, companyID, managerID string) (progression.ManagerView, error) {
	views, err := f.GetCompanyManagers(ctx, p, companyID)
	if err != nil {
		return progression.ManagerView{}, err
	}
	for _, v := range views {
		if v.ManagerID == managerID {
			return v, nil
		}
	}
	return progression.ManagerView{}, ErrNotFound
}

func TestProgressionReadRoutesRequireAuthAndReturnViews(t *testing.T) {
	repo := &fakeRepo{principal: Principal{AccountID: "account.one", PlayerID: "player.one", SessionID: "session.one", Status: "guest"}}
	h := NewServer(repo, "test").Mux()
	for _, path := range []string{"/api/v1/me/skills", "/api/v1/companies/company.one/managers", "/api/v1/companies/company.one/managers/manager.one"} {
		r := httptest.NewRequest(http.MethodGet, path, nil)
		r.Header.Set("Authorization", "Bearer valid-token")
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)
		if w.Code != http.StatusOK {
			t.Fatalf("%s status=%d body=%s", path, w.Code, w.Body.String())
		}
	}
	r := httptest.NewRequest(http.MethodGet, "/api/v1/me/skills", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("unauthenticated skills status=%d", w.Code)
	}
}
