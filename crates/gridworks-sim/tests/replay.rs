use gridworks_sim::{
    advance, execute, Command, KernelError, KernelEvent, SimulationState, RULES_VERSION,
};

fn apply_stream(state: &mut SimulationState) -> Vec<KernelEvent> {
    let mut events = Vec::new();
    events.extend(advance(state, 5_000).expect("time transition").events);
    events.extend(
        execute(
            state,
            &Command::adjust_register("replay-1", "replay-key-1", 5_000, 11),
        )
        .expect("register transition")
        .events,
    );
    events.extend(
        execute(
            state,
            &Command::seeded_pulse("replay-2", "replay-key-2", 5_000, vec![-2, 3, 9]),
        )
        .expect("seeded transition")
        .events,
    );
    events.extend(advance(state, 7_000).expect("time transition").events);
    events
}

#[test]
fn serialized_replay_is_identical() {
    let initial = SimulationState::new(1234).unwrap();
    let snapshot = initial.to_json().unwrap();
    let mut first_state = SimulationState::from_json(&snapshot).unwrap();
    let mut second_state = SimulationState::from_json(&snapshot).unwrap();
    let first_events = apply_stream(&mut first_state);
    let second_events = apply_stream(&mut second_state);

    assert_eq!(first_events, second_events);
    assert_eq!(first_state, second_state);
    assert_eq!(first_state.digest(), second_state.digest());
    assert_eq!(first_events.len(), 4);
}

#[test]
fn incompatible_rules_are_rejected_before_execution() {
    let mut state = SimulationState::new(1234).unwrap();
    let mut command = Command::adjust_register("bad-version", "bad-version-key", 0, 1);
    command.rules_version = "rules-incompatible".to_owned();

    assert!(matches!(
        execute(&mut state, &command),
        Err(KernelError::VersionMismatch {
            expected_rules,
            actual_rules,
            ..
        }) if expected_rules == RULES_VERSION && actual_rules == "rules-incompatible"
    ));
    assert_eq!(state.proof.register, 0);
}
