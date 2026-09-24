from pathlib import Path
import json

root = Path(__file__).resolve().parents[2]
required = [
    "crates/gridworks-sim/src/scenario.rs",
    "crates/gridworks-sim/tests/wp011_fixture.rs",
    "crates/gridworks-godot/src/lib.rs",
    "apps/game/scripts/opening_aggregate_controller.gd",
    "apps/game/scripts/wp011_headless_test.gd",
    "apps/game/scenes/opening_aggregate.tscn",
    "apps/game/scenes/wp011_test.tscn",
    "tests/deterministic/fixtures/wp011/fixture.json",
]
for path in required:
    if not (root / path).is_file():
        raise SystemExit(f"WP-011 missing required path: {path}")
scenario = (root / "crates/gridworks-sim/src/scenario.rs").read_text()
bridge = (root / "crates/gridworks-godot/src/lib.rs").read_text()
controller = (root / "apps/game/scripts/opening_aggregate_controller.gd").read_text()
workflow = (root / ".github/workflows/ci.yml").read_text()
fixture = json.loads((root / "tests/deterministic/fixtures/wp011/fixture.json").read_text())
if "opening_aggregate_scenario" not in scenario or "scenario.opening.aggregate" not in scenario:
    raise SystemExit("WP-011 canonical Rust scenario builder is missing")
if "create_scenario" not in bridge or "opening_aggregate_scenario" not in bridge:
    raise SystemExit("WP-011 scenario bridge operation is missing")
if "GridworksSimBridgeClient" not in controller or "evaluate_facility" not in controller:
    raise SystemExit("WP-011 controller does not use the accepted bridge wrapper")
import re
if any(re.search(pattern, controller.lower()) for pattern in [r"\bcredits\b", r"\brevenue\b", r"experience points", r"\bxp\b"]):
    raise SystemExit("WP-011 controller contains forbidden fake economy/progression scope")
required_fixture = {"scenario_id", "seed", "commands", "expected"}
if not required_fixture.issubset(fixture) or fixture["scenario_id"] != "scenario.opening.aggregate":
    raise SystemExit("WP-011 fixture is incomplete")
for key in ["initial_digest", "partial_recovery_digest", "final_digest", "effective_capacity", "accepted_runs", "finished_aggregate_quantity", "event_types"]:
    if key not in fixture["expected"]:
        raise SystemExit(f"WP-011 golden is missing {key}")
if "wp011_test.tscn" not in workflow or "PNG image data, 390 x 844" not in workflow or "PNG image data, 1440 x 900" not in workflow:
    raise SystemExit("WP-011 headless test is not wired into CI")
print("WP-011 integration shape is valid.")
