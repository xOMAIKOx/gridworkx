use gridworks_sim::{advance, SimulationState};
use serde::Deserialize;
use serde_json::json;
use std::fs;
use std::path::PathBuf;

#[derive(Deserialize)]
struct Fixture {
    seed: u64,
    elapsed_ms: u64,
}

#[test]
fn wp010_fixture_generates_canonical_parity_result() {
    let fixture: Fixture = serde_json::from_str(include_str!(
        "../../../tests/deterministic/fixtures/wp010/fixture.json"
    ))
    .unwrap();
    let mut state = SimulationState::new(fixture.seed).unwrap();
    let transition = advance(&mut state, fixture.elapsed_ms).unwrap();
    let result = json!({
        "snapshot": transition.resulting_state.to_json().unwrap(),
        "digest": transition.resulting_digest,
        "events": transition.events,
    });
    let path =
        PathBuf::from(env!("CARGO_MANIFEST_DIR")).join("../../target/wp010-pure-result.json");
    fs::write(path, serde_json::to_vec_pretty(&result).unwrap()).unwrap();
    assert_eq!(
        transition.resulting_state.digest().unwrap(),
        result["digest"]
    );
}
