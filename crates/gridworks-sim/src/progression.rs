use serde::{Deserialize, Serialize};
use std::collections::BTreeMap;
use std::fmt::{Display, Formatter};

pub const PROGRESSION_VERSION: &str = "progression-0.1.0";
pub const MAX_BPS: u16 = 10_000;
pub const PLAYER_SKILL_IDS: [&str; 12] = [
    "skill.mechanical",
    "skill.electrical",
    "skill.process_engineering",
    "skill.agriculture",
    "skill.mining",
    "skill.energy",
    "skill.water",
    "skill.logistics",
    "skill.construction",
    "skill.commerce",
    "skill.finance",
    "skill.management",
];
pub const MANAGER_SKILL_IDS: [&str; 9] = [
    "manager_skill.operations",
    "manager_skill.technical",
    "manager_skill.maintenance",
    "manager_skill.safety",
    "manager_skill.leadership",
    "manager_skill.logistics",
    "manager_skill.energy_efficiency",
    "manager_skill.crisis_response",
    "manager_skill.mentoring",
];

#[derive(Debug, Clone, PartialEq, Eq)]
pub enum ProgressionError {
    InvalidIdentifier(String),
    UnknownSkill(String),
    UnknownActivity(String),
    InvalidVersion,
    NegativeXP,
    XPOverflow,
    DuplicateSource(String),
    InvalidBps(u16),
    MissingManagerSkill(String),
}
impl Display for ProgressionError {
    fn fmt(&self, f: &mut Formatter<'_>) -> std::fmt::Result {
        write!(f, "{self:?}")
    }
}
impl std::error::Error for ProgressionError {}

#[derive(Debug, Clone, Copy, PartialEq, Eq, Serialize, Deserialize)]
#[serde(rename_all = "snake_case")]
pub enum ActivityClass {
    Trivial,
    Meaningful,
}
#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize)]
#[serde(deny_unknown_fields)]
pub struct SkillDefinition {
    pub skill_id: String,
    pub label_key: String,
    pub activity_classes: Vec<String>,
    pub progression_curve: String,
    pub rules_version: String,
}
pub fn player_skill_registry() -> Vec<SkillDefinition> {
    PLAYER_SKILL_IDS
        .iter()
        .map(|id| SkillDefinition {
            skill_id: (*id).to_owned(),
            label_key: format!("label.{id}"),
            activity_classes: vec![
                "diagnosis".to_owned(),
                "repair".to_owned(),
                "production".to_owned(),
            ],
            progression_curve: "curve.linear-capped-v1".to_owned(),
            rules_version: PROGRESSION_VERSION.to_owned(),
        })
        .collect()
}
#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize)]
#[serde(deny_unknown_fields)]
pub struct AntiGrindPolicy {
    pub policy_id: String,
    pub version: String,
    pub trivial_multipliers_bps: Vec<u16>,
    pub meaningful_multiplier_bps: u16,
}
pub fn anti_grind_policy() -> AntiGrindPolicy {
    AntiGrindPolicy {
        policy_id: "policy.progression.anti_grind".to_owned(),
        version: PROGRESSION_VERSION.to_owned(),
        trivial_multipliers_bps: vec![10_000, 5_000, 2_500, 0],
        meaningful_multiplier_bps: 10_000,
    }
}
#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize)]
#[serde(deny_unknown_fields)]
pub struct SkillActivityEvent {
    pub source_event_id: String,
    pub player_id: String,
    pub skill_id: String,
    pub activity_kind: String,
    pub base_xp: u64,
    pub activity_class: ActivityClass,
    pub repetition_key: String,
    pub occurrence_time_ms: u64,
    pub rules_version: String,
}
#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize, Default)]
#[serde(deny_unknown_fields)]
pub struct PlayerSkillState {
    pub player_id: String,
    pub skill_id: String,
    pub cumulative_xp: u64,
    pub proficiency_bps: u16,
    pub progression_version: String,
    pub last_source_event_id: Option<String>,
    pub last_event_time_ms: Option<u64>,
    pub repetition_key: Option<String>,
    pub repetition_count: u16,
    pub applied_source_events: Vec<String>,
}
#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize)]
pub struct SkillProgressionResult {
    pub awarded_xp: u64,
    pub cumulative_xp: u64,
    pub proficiency_bps: u16,
    pub replay: bool,
    pub repetition_count: u16,
}
pub fn proficiency_from_xp(xp: u64) -> u16 {
    xp.saturating_mul(10).min(u64::from(MAX_BPS)) as u16
}
pub fn apply_player_skill_event(
    state: &mut PlayerSkillState,
    event: &SkillActivityEvent,
) -> Result<SkillProgressionResult, ProgressionError> {
    if state.player_id != event.player_id {
        return Err(ProgressionError::InvalidIdentifier(event.player_id.clone()));
    }
    if state.skill_id != event.skill_id || !PLAYER_SKILL_IDS.contains(&event.skill_id.as_str()) {
        return Err(ProgressionError::UnknownSkill(event.skill_id.clone()));
    }
    if event.source_event_id.is_empty() || event.repetition_key.is_empty() {
        return Err(ProgressionError::InvalidIdentifier(
            "missing progression identity".to_owned(),
        ));
    }
    if event.rules_version != PROGRESSION_VERSION {
        return Err(ProgressionError::InvalidVersion);
    }
    if state
        .applied_source_events
        .iter()
        .any(|id| id == &event.source_event_id)
    {
        return Ok(SkillProgressionResult {
            awarded_xp: 0,
            cumulative_xp: state.cumulative_xp,
            proficiency_bps: state.proficiency_bps,
            replay: true,
            repetition_count: state.repetition_count,
        });
    }
    let policy = anti_grind_policy();
    let count = if state.repetition_key.as_deref() == Some(event.repetition_key.as_str()) {
        state.repetition_count.saturating_add(1)
    } else {
        1
    };
    let multiplier = match event.activity_class {
        ActivityClass::Meaningful => policy.meaningful_multiplier_bps,
        ActivityClass::Trivial => policy
            .trivial_multipliers_bps
            .get(usize::from(count - 1))
            .copied()
            .unwrap_or(0),
    };
    let awarded = event
        .base_xp
        .checked_mul(u64::from(multiplier))
        .ok_or(ProgressionError::XPOverflow)?
        .checked_div(u64::from(MAX_BPS))
        .ok_or(ProgressionError::XPOverflow)?;
    state.cumulative_xp = state
        .cumulative_xp
        .checked_add(awarded)
        .ok_or(ProgressionError::XPOverflow)?;
    state.proficiency_bps = proficiency_from_xp(state.cumulative_xp);
    state.progression_version = PROGRESSION_VERSION.to_owned();
    state.last_source_event_id = Some(event.source_event_id.clone());
    state.last_event_time_ms = Some(event.occurrence_time_ms);
    state.repetition_key = Some(event.repetition_key.clone());
    state.repetition_count = count;
    state
        .applied_source_events
        .push(event.source_event_id.clone());
    Ok(SkillProgressionResult {
        awarded_xp: awarded,
        cumulative_xp: state.cumulative_xp,
        proficiency_bps: state.proficiency_bps,
        replay: false,
        repetition_count: count,
    })
}

#[derive(Debug, Clone, Copy, PartialEq, Eq, Serialize, Deserialize, PartialOrd, Ord)]
pub enum ManagerRarity {
    Bronze,
    Silver,
    Gold,
    Platinum,
}
#[derive(Debug, Clone, Copy, PartialEq, Eq, Serialize, Deserialize)]
pub enum ManagerStatus {
    Active,
    Inactive,
}
#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize)]
#[serde(deny_unknown_fields)]
pub struct ManagerSkillState {
    pub skill_id: String,
    pub proficiency_bps: u16,
    pub potential_bps: u16,
}
#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize)]
#[serde(deny_unknown_fields)]
pub struct ManagerState {
    pub manager_id: String,
    pub display_name: String,
    pub rarity: ManagerRarity,
    pub specialization_id: String,
    pub total_xp: u64,
    pub level: u32,
    pub skills: Vec<ManagerSkillState>,
    pub traits: Vec<String>,
    pub workload_bps: u16,
    pub fatigue_bps: u16,
    pub morale_bps: u16,
    pub status: ManagerStatus,
    pub created_at_ms: u64,
    pub effective_from_ms: u64,
}
pub fn manager_potential_cap(
    rarity: ManagerRarity,
    skill_id: &str,
) -> Result<u16, ProgressionError> {
    if !MANAGER_SKILL_IDS.contains(&skill_id) {
        return Err(ProgressionError::MissingManagerSkill(skill_id.to_owned()));
    }
    Ok(match rarity {
        ManagerRarity::Bronze => 5_000,
        ManagerRarity::Silver => 6_500,
        ManagerRarity::Gold => 8_500,
        ManagerRarity::Platinum => 10_000,
    })
}
pub fn manager_level_from_xp(xp: u64) -> u32 {
    ((xp / 1000) as u32).saturating_add(1)
}
pub fn cap_manager_skill(
    manager: &mut ManagerState,
    skill_id: &str,
    requested_bps: u16,
) -> Result<u16, ProgressionError> {
    if requested_bps > MAX_BPS {
        return Err(ProgressionError::InvalidBps(requested_bps));
    }
    let cap = manager_potential_cap(manager.rarity, skill_id)?;
    let skill = manager
        .skills
        .iter_mut()
        .find(|s| s.skill_id == skill_id)
        .ok_or_else(|| ProgressionError::MissingManagerSkill(skill_id.to_owned()))?;
    skill.potential_bps = skill.potential_bps.min(cap);
    skill.proficiency_bps = requested_bps.min(skill.potential_bps);
    Ok(skill.proficiency_bps)
}
#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize)]
#[serde(deny_unknown_fields)]
pub struct TraitDefinition {
    pub trait_id: String,
    pub effects: BTreeMap<String, i16>,
    pub rules_version: String,
}
pub fn manager_trait_definitions() -> Vec<TraitDefinition> {
    [
        ("trait.aggressive_operator", 100i16),
        ("trait.maintenance_first", 250),
        ("trait.cost_cutter", 0),
        ("trait.mentor", 150),
        ("trait.crisis_specialist", 500),
    ]
    .into_iter()
    .map(|(id, effect)| {
        let mut effects = BTreeMap::new();
        effects.insert("diagnostic_capability_bps".to_owned(), effect);
        TraitDefinition {
            trait_id: id.to_owned(),
            effects,
            rules_version: PROGRESSION_VERSION.to_owned(),
        }
    })
    .collect()
}
pub fn bounded_fatigue_modifier(fatigue_bps: u16) -> Result<u16, ProgressionError> {
    if fatigue_bps > MAX_BPS {
        return Err(ProgressionError::InvalidBps(fatigue_bps));
    }
    Ok(MAX_BPS - fatigue_bps)
}
pub fn effective_diagnostic_capability(
    player_skill_bps: u16,
    manager: Option<&ManagerState>,
) -> Result<u16, ProgressionError> {
    if player_skill_bps > MAX_BPS {
        return Err(ProgressionError::InvalidBps(player_skill_bps));
    }
    let mut total = u32::from(player_skill_bps) * 6 / 10;
    if let Some(manager) = manager {
        if manager.workload_bps > MAX_BPS
            || manager.fatigue_bps > MAX_BPS
            || manager.morale_bps > MAX_BPS
        {
            return Err(ProgressionError::InvalidBps(MAX_BPS + 1));
        }
        let relevant = manager
            .skills
            .iter()
            .find(|s| s.skill_id == "manager_skill.technical")
            .map(|s| s.proficiency_bps)
            .unwrap_or(0);
        let crisis = manager
            .skills
            .iter()
            .find(|s| s.skill_id == "manager_skill.crisis_response")
            .map(|s| s.proficiency_bps)
            .unwrap_or(0);
        let base = (u32::from(relevant) * 6 + u32::from(crisis) * 4) / 10;
        let fatigue = u32::from(bounded_fatigue_modifier(manager.fatigue_bps)?);
        let morale = 8_000u32
            .saturating_add(u32::from(manager.morale_bps) / 5)
            .min(10_000);
        let contribution = base * fatigue / 10_000 * morale / 10_000;
        let trait_bonus: i32 = manager
            .traits
            .iter()
            .map(|id| {
                manager_trait_definitions()
                    .into_iter()
                    .find(|d| &d.trait_id == id)
                    .and_then(|d| d.effects.get("diagnostic_capability_bps").copied())
                    .map(i32::from)
                    .unwrap_or(0)
            })
            .sum();
        total = total
            .saturating_add(contribution)
            .saturating_add_signed(trait_bonus);
    }
    Ok(total.min(u32::from(MAX_BPS)) as u16)
}
fn fixture_manager(
    id: &str,
    name: &str,
    rarity: ManagerRarity,
    technical: u16,
    potential: u16,
    fatigue: u16,
    traits: Vec<&str>,
) -> ManagerState {
    ManagerState {
        manager_id: id.to_owned(),
        display_name: name.to_owned(),
        rarity,
        specialization_id: "specialization.aggregate_diagnostics".to_owned(),
        total_xp: 0,
        level: 1,
        skills: MANAGER_SKILL_IDS
            .iter()
            .map(|skill_id| ManagerSkillState {
                skill_id: (*skill_id).to_owned(),
                proficiency_bps: if *skill_id == "manager_skill.technical" {
                    technical
                } else {
                    0
                },
                potential_bps: if *skill_id == "manager_skill.technical" {
                    potential
                } else {
                    manager_potential_cap(rarity, skill_id).unwrap_or(0)
                },
            })
            .collect(),
        traits: traits.into_iter().map(str::to_owned).collect(),
        workload_bps: 0,
        fatigue_bps: fatigue,
        morale_bps: 10_000,
        status: ManagerStatus::Active,
        created_at_ms: 0,
        effective_from_ms: 0,
    }
}
pub fn veteran_gold_fixture() -> ManagerState {
    fixture_manager(
        "manager.fixture.gold",
        "Veteran Gold",
        ManagerRarity::Gold,
        8_000,
        8_500,
        500,
        vec!["trait.crisis_specialist"],
    )
}
pub fn fresh_platinum_fixture() -> ManagerState {
    fixture_manager(
        "manager.fixture.platinum",
        "Fresh Platinum",
        ManagerRarity::Platinum,
        1_000,
        10_000,
        0,
        Vec::new(),
    )
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::failure::{FailureCommand, FailureState, Observation};
    use crate::graph::aggregate_plant_fixture;
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
    fn player_progression_is_monotonic_and_diminishes_trivial_repetition() {
        let mut state = PlayerSkillState {
            player_id: "player.fixture".to_owned(),
            skill_id: "skill.mechanical".to_owned(),
            progression_version: PROGRESSION_VERSION.to_owned(),
            ..Default::default()
        };
        assert_eq!(
            apply_player_skill_event(&mut state, &event("event.1", 100, "same"))
                .unwrap()
                .awarded_xp,
            100
        );
        assert_eq!(
            apply_player_skill_event(&mut state, &event("event.2", 100, "same"))
                .unwrap()
                .awarded_xp,
            50
        );
        assert_eq!(
            apply_player_skill_event(&mut state, &event("event.3", 100, "same"))
                .unwrap()
                .awarded_xp,
            25
        );
        assert_eq!(
            apply_player_skill_event(&mut state, &event("event.4", 100, "same"))
                .unwrap()
                .awarded_xp,
            0
        );
        assert!(
            apply_player_skill_event(&mut state, &event("event.4", 100, "same"))
                .unwrap()
                .replay
        );
        assert_eq!(state.cumulative_xp, 175);
    }
    #[test]
    fn gold_can_outperform_platinum_without_rarity_multiplier() {
        let gold = veteran_gold_fixture();
        let platinum = fresh_platinum_fixture();
        assert!(
            gold.skills
                .iter()
                .find(|s| s.skill_id == "manager_skill.technical")
                .unwrap()
                .potential_bps
                < platinum
                    .skills
                    .iter()
                    .find(|s| s.skill_id == "manager_skill.technical")
                    .unwrap()
                    .potential_bps
        );
        assert!(
            effective_diagnostic_capability(4_000, Some(&gold)).unwrap()
                > effective_diagnostic_capability(4_000, Some(&platinum)).unwrap()
        );
    }
    #[test]
    fn capability_changes_confidence_not_observation_identity() {
        let facilities = vec![aggregate_plant_fixture()];
        let mut state = FailureState::with_definitions(crate::aggregate_fault_definitions());
        state
            .apply(
                &facilities,
                &FailureCommand::ActivateFault {
                    fault_instance_id: "fault.instance.fixture".to_owned(),
                    fault_type_id: "fault.motor_bearing_seizure".to_owned(),
                    component_id: "component.feed_conveyor".to_owned(),
                    severity_bps: 10_000,
                },
                0,
                crate::RULES_VERSION,
            )
            .unwrap();
        let low = effective_diagnostic_capability(1000, Some(&fresh_platinum_fixture())).unwrap();
        let high = effective_diagnostic_capability(4000, Some(&veteran_gold_fixture())).unwrap();
        let low_event = state
            .apply(
                &facilities,
                &FailureCommand::Inspect {
                    evidence_id: "evidence.low".to_owned(),
                    component_id: "component.feed_conveyor".to_owned(),
                    diagnostic_capability_bps: low,
                },
                0,
                crate::RULES_VERSION,
            )
            .unwrap();
        let identity = match low_event {
            crate::FailureEvent::EvidenceObserved {
                symptom_id,
                observation,
                ..
            } => {
                assert_eq!(observation, Observation::Unknown);
                symptom_id
            }
            _ => panic!("wrong event"),
        };
        let high_event = state
            .apply(
                &facilities,
                &FailureCommand::Inspect {
                    evidence_id: "evidence.high".to_owned(),
                    component_id: "component.feed_conveyor".to_owned(),
                    diagnostic_capability_bps: high,
                },
                0,
                crate::RULES_VERSION,
            )
            .unwrap();
        match high_event {
            crate::FailureEvent::EvidenceObserved {
                symptom_id,
                observation,
                ..
            } => {
                assert_eq!(symptom_id, identity);
                assert_eq!(observation, Observation::Observed);
            }
            _ => panic!("wrong event"),
        }
        assert!(high >= low);
    }
}
