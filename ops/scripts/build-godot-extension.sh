#!/bin/sh
set -eu

cargo test -p gridworks-sim --test wp010_fixture
cargo build -p gridworks-godot
mkdir -p apps/game/bin/linux/debug
cp target/debug/libgridworks_godot.so apps/game/bin/linux/debug/libgridworks_godot.so
cp tests/deterministic/fixtures/wp010/fixture.json apps/game/bin/linux/debug/wp010-fixture.json
cp target/wp010-pure-result.json apps/game/bin/linux/debug/wp010-pure-result.json
test -s apps/game/bin/linux/debug/libgridworks_godot.so
test -s apps/game/bin/linux/debug/wp010-pure-result.json
mkdir -p apps/game/.godot
printf '%s\n' 'res://native/gridworks_sim.gdextension' > apps/game/.godot/extension_list.cfg
grep -Fx 'res://native/gridworks_sim.gdextension' apps/game/.godot/extension_list.cfg
if git ls-files --error-unmatch apps/game/bin/linux/debug/libgridworks_godot.so >/dev/null 2>&1; then
	echo "compiled Godot extension must not be tracked" >&2
	exit 1
fi
