use gridworks_sim::{
    apply_player_skill_event, effective_diagnostic_capability, fresh_platinum_fixture,
    veteran_gold_fixture, ActivityClass, PlayerSkillState, SkillActivityEvent, PROGRESSION_VERSION,
};
use serde::Deserialize;

#[derive(Deserialize)]
struct Fixture {
    player: PlayerFixture,
    manager: ManagerFixture,
    observation: ObservationFixture,
}
#[derive(Deserialize)]
struct PlayerFixture {
    player_id: String,
    skill_id: String,
    trivial_awards: Vec<u64>,
    duplicate_award: u64,
    cumulative_xp: u64,
    proficiency_bps: u16,
}
#[derive(Deserialize)]
struct ManagerFixture {
    gold_current_technical_bps: u16,
    gold_potential_bps: u16,
    platinum_current_technical_bps: u16,
    platinum_potential_bps: u16,
    gold_cap_after_request_bps: u16,
    gold_diagnostic_capability_bps: u16,
    platinum_diagnostic_capability_bps: u16,
    fatigue_reduces_contribution: bool,
}
#[derive(Deserialize)]
struct ObservationFixture {
    symptom_id: String,
    low_capability_observation: String,
    high_capability_observation: String,
}

fn event(id: &str, base: u64, repetition: &str) -> SkillActivityEvent {
    SkillActivityEvent {
        source_event_id: id.to_owned(),
        player_id: "player.fixture".to_owned(),
        skill_id: "skill.mechanical".to_owned(),
        activity_kind: "diagnosis".to_owned(),
        base_xp: base,
        activity_class: ActivityClass::Trivial,
        repetition_key: repetition.to_owned(),
        occurrence_time_ms: base,
        rules_version: PROGRESSION_VERSION.to_owned(),
    }
}

#[test]
fn wp012_golden_asserts_progression_and_manager_capability() {
    let fixture: Fixture = serde_json::from_str(include_str!(
        "../../../tests/deterministic/fixtures/wp012/fixture.json"
    ))
    .unwrap();
    assert_eq!(fixture.player.player_id, "player.fixture");
    assert_eq!(fixture.observation.symptom_id, "symptom.abnormal_vibration");
    assert_eq!(fixture.observation.low_capability_observation, "Unknown");
    assert_eq!(fixture.observation.high_capability_observation, "Observed");
    assert_eq!(fixture.manager.gold_current_technical_bps, 8000);
    assert_eq!(fixture.manager.platinum_current_technical_bps, 1000);
    assert_eq!(fixture.manager.gold_potential_bps, 8500);
    assert_eq!(fixture.manager.platinum_potential_bps, 10000);
    assert!(fixture.manager.fatigue_reduces_contribution);
    let mut state = PlayerSkillState {
        player_id: fixture.player.player_id.clone(),
        skill_id: fixture.player.skill_id.clone(),
        progression_version: PROGRESSION_VERSION.to_owned(),
        ..Default::default()
    };
    let mut awards = Vec::new();
    for (index, base) in [100, 100, 100, 100].into_iter().enumerate() {
        awards.push(
            apply_player_skill_event(&mut state, &event(&format!("event.{index}"), base, "same"))
                .unwrap()
                .awarded_xp,
        );
    }
    let replay = apply_player_skill_event(&mut state, &event("event.3", 100, "same")).unwrap();
    assert_eq!(awards, fixture.player.trivial_awards);
    assert_eq!(replay.awarded_xp, fixture.player.duplicate_award);
    assert_eq!(state.cumulative_xp, fixture.player.cumulative_xp);
    assert_eq!(state.proficiency_bps, fixture.player.proficiency_bps);
    let mut gold = veteran_gold_fixture();
    assert_eq!(
        gridworks_sim::progression::cap_manager_skill(&mut gold, "manager_skill.technical", 10_000)
            .unwrap(),
        fixture.manager.gold_cap_after_request_bps
    );
    assert_eq!(
        effective_diagnostic_capability(4_000, Some(&veteran_gold_fixture())).unwrap(),
        fixture.manager.gold_diagnostic_capability_bps
    );
    assert_eq!(
        effective_diagnostic_capability(4_000, Some(&fresh_platinum_fixture())).unwrap(),
        fixture.manager.platinum_diagnostic_capability_bps
    );
    assert!(
        fixture.manager.gold_diagnostic_capability_bps
            > fixture.manager.platinum_diagnostic_capability_bps
    );
}
