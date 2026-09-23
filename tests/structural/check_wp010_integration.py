from pathlib import Path
import json

root = Path(__file__).resolve().parents[2]
cargo = (root / "Cargo.toml").read_text()
bridge = root / "crates/gridworks-godot"
if '"crates/gridworks-godot"' not in cargo:
    raise SystemExit("WP-010 bridge crate is not in the workspace")
manifest = (bridge / "Cargo.toml").read_text()
if 'godot = "=0.2.4"' not in manifest and 'version = "=0.2.4"' not in manifest:
    raise SystemExit("WP-010 godot binding is not exactly pinned to 0.2.4")
if 'gridworks-sim' not in manifest or 'crate-type = ["cdylib"]' not in manifest:
    raise SystemExit("WP-010 bridge crate contract is incomplete")
for path in [
    "apps/game/native/gridworks_sim.gdextension",
    "apps/game/scripts/sim_bridge.gd",
    "apps/game/scripts/wp010_headless_test.gd",
    "apps/game/scenes/wp010_test.tscn",
    "tests/deterministic/fixtures/wp010/fixture.json",
    "ops/scripts/build-godot-extension.sh",
    "crates/gridworks-sim/tests/wp010_fixture.rs",
]:
    if not (root / path).is_file():
        raise SystemExit(f"WP-010 missing required path: {path}")
descriptor = (root / "apps/game/native/gridworks_sim.gdextension").read_text()
if 'entry_symbol = "gdext_rust_init"' not in descriptor or 'linux.debug.x86_64' not in descriptor:
    raise SystemExit("WP-010 descriptor is incomplete")
fixture = json.loads((root / "tests/deterministic/fixtures/wp010/fixture.json").read_text())
required = {"seed", "elapsed_ms", "authoritative_time_ms", "facility_id", "commands", "expected"}
if not required.issubset(fixture) or len(fixture["commands"]) == 0:
    raise SystemExit("WP-010 fixture does not pin the required command golden")
if not all(fixture["expected"].get(key) for key in ["after_elapsed_digest", "after_command_digest", "final_digest"]):
    raise SystemExit("WP-010 fixture is missing committed golden digests")
build = (root / "ops/scripts/build-godot-extension.sh").read_text()
if ".godot/extension_list.cfg" in build or "printf" in build:
    raise SystemExit("WP-010 build helper must not manufacture Godot extension-list project data")
workflow = (root / ".github/workflows/ci.yml").read_text()
if "--headless --path apps/game --import --rendering-method gl_compatibility --rendering-driver opengl3 --audio-driver Dummy" not in workflow:
    raise SystemExit("WP-010 CI is missing the Godot import discovery scan")
if "--editor" in workflow or "--quit-after" in workflow or "-eq 134" in workflow or "discovery_status" in workflow:
    raise SystemExit("WP-010 CI must use strict import discovery without forced editor exit or abort whitelist")
bridge_source = (bridge / "src/lib.rs").read_text()
for method in ["advance_to_snapshot", "evaluate_facility", "validate_digest"]:
    if method not in bridge_source:
        raise SystemExit(f"WP-010 bridge is missing required method: {method}")
if "catch_unwind" not in bridge_source:
    raise SystemExit("WP-010 bridge is missing panic-safety guard")
print("WP-010 integration shape is valid.")
