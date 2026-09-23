use gridworks_sim::{
    advance, advance_to, aggregate_plant_fixture, execute, Command, SimulationState,
};
use serde::Deserialize;
use serde_json::json;
use std::fs;
use std::path::PathBuf;

#[derive(Deserialize)]
struct Fixture {
    seed: u64,
    elapsed_ms: u64,
    authoritative_time_ms: u64,
    facility_id: String,
    commands: Vec<Command>,
    expected: Expected,
}

#[derive(Deserialize)]
struct Expected {
    after_elapsed_digest: String,
    after_command_digest: String,
    final_digest: String,
    event_types: Vec<String>,
    command_event_types: Vec<String>,
}

fn event_types(events: &[gridworks_sim::KernelEvent]) -> Vec<String> {
    events
        .iter()
        .map(|event| match event {
            gridworks_sim::KernelEvent::TimeAdvanced { .. } => "time.advanced",
            gridworks_sim::KernelEvent::RegisterAdjusted { .. } => "proof.register_adjusted",
            gridworks_sim::KernelEvent::SeededPulse { .. } => "proof.seeded_pulse",
            gridworks_sim::KernelEvent::Failure(_) => "failure",
            gridworks_sim::KernelEvent::Material(_) => "material",
        })
        .map(str::to_owned)
        .collect()
}

#[test]
fn wp010_fixture_asserts_committed_golden_and_writes_parity_artifacts() {
    let fixture: Fixture = serde_json::from_str(include_str!(
        "../../../tests/deterministic/fixtures/wp010/fixture.json"
    ))
    .unwrap();
    let mut state = SimulationState::new(fixture.seed).unwrap();
    let first = advance(&mut state, fixture.elapsed_ms).unwrap();
    let after_elapsed_digest = first.resulting_digest.clone();
    let mut events = first.events.clone();
    let mut command_events = Vec::new();
    for command in &fixture.commands {
        let transition = execute(&mut state, command).unwrap();
        command_events.extend(transition.events.clone());
        events.extend(transition.events);
    }
    let after_command_digest = state.digest().unwrap();
    let after_commands_snapshot = state.to_json().unwrap();
    let final_transition = advance_to(&mut state, fixture.authoritative_time_ms).unwrap();
    events.extend(final_transition.events.clone());
    let final_digest = final_transition.resulting_digest.clone();

    assert!(
        !fixture.expected.after_elapsed_digest.is_empty(),
        "committed golden is missing after_elapsed_digest"
    );
    assert!(
        !fixture.expected.after_command_digest.is_empty(),
        "committed golden is missing after_command_digest"
    );
    assert!(
        !fixture.expected.final_digest.is_empty(),
        "committed golden is missing final_digest"
    );
    assert_eq!(after_elapsed_digest, fixture.expected.after_elapsed_digest);
    assert_eq!(after_command_digest, fixture.expected.after_command_digest);
    assert_eq!(final_digest, fixture.expected.final_digest);
    assert_eq!(event_types(&events), fixture.expected.event_types);
    assert_eq!(
        event_types(&command_events),
        fixture.expected.command_event_types
    );

    let result = json!({
        "after_elapsed": {"snapshot": first.resulting_state.to_json().unwrap(), "digest": after_elapsed_digest, "events": first.events},
        "after_commands": {"snapshot": after_commands_snapshot, "digest": after_command_digest, "events": command_events},
        "final": {"snapshot": final_transition.resulting_state.to_json().unwrap(), "digest": final_digest, "events": final_transition.events},
        "event_types": event_types(&events),
    });
    let target = PathBuf::from(env!("CARGO_MANIFEST_DIR")).join("../../target");
    fs::write(
        target.join("wp010-pure-result.json"),
        serde_json::to_vec_pretty(&result).unwrap(),
    )
    .unwrap();

    let facility_state =
        SimulationState::with_facilities(fixture.seed, vec![aggregate_plant_fixture()]).unwrap();
    let facility = facility_state
        .evaluate_facility(&fixture.facility_id)
        .unwrap();
    fs::write(
        target.join("wp010-facility-result.json"),
        serde_json::to_vec_pretty(&json!({
            "snapshot": facility_state.to_json().unwrap(),
            "facility_id": fixture.facility_id,
            "evaluation": facility,
        }))
        .unwrap(),
    )
    .unwrap();
}
