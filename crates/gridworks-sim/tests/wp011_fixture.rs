use gridworks_sim::{execute, opening_aggregate_scenario, Command, KernelEvent};
use serde::Deserialize;
use serde_json::json;
use std::fs;
use std::path::PathBuf;

#[derive(Deserialize)]
struct Fixture {
    scenario_id: String,
    seed: u64,
    commands: Vec<Command>,
    expected: Expected,
}

#[derive(Deserialize)]
struct Expected {
    initial_digest: String,
    partial_recovery_digest: String,
    final_digest: String,
    effective_capacity: u64,
    accepted_runs: u64,
    finished_aggregate_quantity: u64,
    event_types: Vec<String>,
}

fn event_type(event: &KernelEvent) -> &'static str {
    match event {
        KernelEvent::Failure(_) => "failure",
        KernelEvent::Material(_) => "material",
        KernelEvent::TimeAdvanced { .. } => "time.advanced",
        KernelEvent::RegisterAdjusted { .. } => "proof.register_adjusted",
        KernelEvent::SeededPulse { .. } => "proof.seeded_pulse",
    }
}

#[test]
fn wp011_fixture_asserts_committed_golden() {
    let fixture: Fixture = serde_json::from_str(include_str!(
        "../../../tests/deterministic/fixtures/wp011/fixture.json"
    ))
    .unwrap();
    assert_eq!(fixture.scenario_id, "scenario.opening.aggregate");
    let mut state = opening_aggregate_scenario(fixture.seed).unwrap();
    let initial_digest = state.digest().unwrap();
    let mut events = Vec::new();
    let mut partial_recovery_digest = String::new();
    let mut accepted_runs = 0;
    for (index, command) in fixture.commands.iter().enumerate() {
        let transition = execute(&mut state, command).unwrap();
        events.extend(transition.events.clone());
        if index == 4 {
            partial_recovery_digest = state.digest().unwrap();
            assert_eq!(
                state
                    .evaluate_facility("facility.aggregate_plant_fixture")
                    .unwrap()
                    .effective_capacity,
                fixture.expected.effective_capacity
            );
        }
        if index == 5 {
            if let Some(KernelEvent::Material(gridworks_sim::MaterialEvent::ProductionExecuted {
                accepted_runs: runs,
                ..
            })) = transition.events.first()
            {
                accepted_runs = *runs;
            }
        }
    }
    let final_digest = state.digest().unwrap();
    let finished = state
        .material
        .inventories
        .iter()
        .find(|inventory| inventory.inventory_id == "inventory.aggregate_finished")
        .unwrap()
        .quantity("resource.finished_aggregate", "grade.aggregate.standard")
        .unwrap();
    let event_types = events
        .iter()
        .map(event_type)
        .map(str::to_owned)
        .collect::<Vec<_>>();
    assert!(!fixture.expected.initial_digest.is_empty());
    assert!(!fixture.expected.partial_recovery_digest.is_empty());
    assert!(!fixture.expected.final_digest.is_empty());
    assert_eq!(initial_digest, fixture.expected.initial_digest);
    assert_eq!(
        partial_recovery_digest,
        fixture.expected.partial_recovery_digest
    );
    assert_eq!(final_digest, fixture.expected.final_digest);
    assert_eq!(accepted_runs, fixture.expected.accepted_runs);
    assert_eq!(finished, fixture.expected.finished_aggregate_quantity);
    assert_eq!(event_types, fixture.expected.event_types);

    let output = json!({"initial_digest": initial_digest, "partial_recovery_digest": partial_recovery_digest, "final_digest": final_digest, "effective_capacity": fixture.expected.effective_capacity, "accepted_runs": accepted_runs, "finished_aggregate_quantity": finished, "event_types": event_types});
    let target = PathBuf::from(env!("CARGO_MANIFEST_DIR")).join("../../target");
    fs::write(
        target.join("wp011-pure-result.json"),
        serde_json::to_vec_pretty(&output).unwrap(),
    )
    .unwrap();
}
