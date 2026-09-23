extends SceneTree

func fail(message: String) -> void:
	push_error(message)
	quit(1)

func parse_bridge(payload: String) -> Dictionary:
	var parsed = JSON.parse_string(payload)
	if typeof(parsed) != TYPE_DICTIONARY:
		fail("bridge returned non-object JSON")
	return parsed

func _init() -> void:
	call_deferred("_run")

func _run() -> void:
	var bridge = ClassDB.instantiate("GridworksSimBridge")
	if bridge == null:
		fail("GDExtension class was not registered")
	for method in ["bridge_metadata", "create_snapshot", "advance_snapshot", "execute_command_batch", "digest_snapshot", "facility_summary"]:
		if not bridge.has_method(method):
			fail("missing bridge method: " + method)
	var metadata = parse_bridge(bridge.bridge_metadata())
	if not metadata.get("ok", false) or metadata["result"]["bridge_version"] != "godot-rust-bridge-0.1.0":
		fail("bridge metadata mismatch")
	var fixture = JSON.parse_string(FileAccess.get_file_as_string("res://bin/linux/debug/wp010-fixture.json"))
	var expected = JSON.parse_string(FileAccess.get_file_as_string("res://bin/linux/debug/wp010-pure-result.json"))
	var created = parse_bridge(bridge.create_snapshot(int(fixture["seed"])))
	if not created.get("ok", false):
		fail("snapshot creation failed")
	var advanced = parse_bridge(bridge.advance_snapshot(created["result"]["snapshot"], int(fixture["elapsed_ms"])))
	if not advanced.get("ok", false):
		fail("snapshot advance failed")
	if advanced["result"] != expected:
		fail("Godot/native result diverged from pure Rust fixture result")
	var repeated = parse_bridge(bridge.advance_snapshot(created["result"]["snapshot"], int(fixture["elapsed_ms"])))
	if repeated["result"] != advanced["result"]:
		fail("repeated stateless bridge call diverged")
	var batch = parse_bridge(bridge.execute_command_batch(created["result"]["snapshot"], "[]"))
	if not batch.get("ok", false) or batch["result"]["digest"] != created["result"]["digest"]:
		fail("empty command batch changed canonical state")
	var digest = parse_bridge(bridge.digest_snapshot(created["result"]["snapshot"]))
	if not digest.get("ok", false) or digest["result"]["digest"] != created["result"]["digest"]:
		fail("digest operation mismatch")
	var facility = parse_bridge(bridge.facility_summary())
	if not facility.get("ok", false):
		fail("facility summary failed")
	var malformed = parse_bridge(bridge.advance_snapshot("{", 1))
	if malformed.get("ok", true) or malformed["error"]["code"] != "bridge.invalid_input":
		fail("malformed snapshot was not rejected")
	quit(0)
