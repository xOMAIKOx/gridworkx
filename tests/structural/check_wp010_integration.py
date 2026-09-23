from pathlib import Path
import json

root = Path(__file__).resolve().parents[2]
cargo = (root / "Cargo.toml").read_text()
bridge = root / "crates/gridworks-godot"
if '"crates/gridworks-godot"' not in cargo: raise SystemExit("WP-010 bridge crate is not in the workspace")
manifest = (bridge / "Cargo.toml").read_text()
if 'version = "=0.2.4"' not in manifest: raise SystemExit("WP-010 godot binding is not exactly pinned to 0.2.4")
if 'gridworks-sim' not in manifest or 'crate-type = ["cdylib", "rlib"]' not in manifest: raise SystemExit("WP-010 bridge crate contract is incomplete")
for path in ["apps/game/native/gridworks_sim.gdextension", "apps/game/scripts/wp010_headless_test.gd", "tests/deterministic/fixtures/wp010/fixture.json", "ops/scripts/build-godot-extension.sh", "crates/gridworks-sim/tests/wp010_fixture.rs"]:
    if not (root / path).is_file(): raise SystemExit(f"WP-010 missing required path: {path}")
descriptor=(root/"apps/game/native/gridworks_sim.gdextension").read_text()
if 'entry_symbol = "gdext_rust_init"' not in descriptor or 'linux.debug.x86_64' not in descriptor: raise SystemExit("WP-010 descriptor is incomplete")
fixture=json.loads((root/"tests/deterministic/fixtures/wp010/fixture.json").read_text())
if fixture.get("seed") != 1234 or fixture.get("elapsed_ms") != 5000: raise SystemExit("WP-010 fixture is not the pinned deterministic fixture")
print("WP-010 integration shape is valid.")
