extends Node

func fail(message: String) -> void:
	push_error(message)
	get_tree().quit(1)

func parse_bridge(payload: String) -> Dictionary:
	var parsed = JSON.parse_string(payload)
	if typeof(parsed) != TYPE_DICTIONARY:
		fail("bridge returned non-object JSON")
	return parsed

func _ready() -> void:
	await get_tree().process_frame
	await get_tree().process_frame
	print("WP010 ready")
	_run()

func _run() -> void:
	if not ClassDB.class_exists("GridworksSimBridge"):
		var candidates: Array[String] = []
		for candidate in ClassDB.get_class_list():
			if "grid" in candidate.to_lower():
				candidates.append(candidate)
		fail("GDExtension class is absent from ClassDB: " + ",".join(candidates))
	print("WP010 class present")
	var bridge = ClassDB.instantiate("GridworksSimBridge")
	if bridge == null:
		fail("GDExtension class was not registered")
	for method in ["bridge_metadata", "create_snapshot", "advance_snapshot", "execute_command_batch", "digest_snapshot", "facility_summary"]:
		if not bridge.has_method(method):
			fail("missing bridge method: " + method)
	print("WP010 bridge instantiated")
	var metadata = parse_bridge(bridge.bridge_metadata())
	if not metadata.get("ok", false) or metadata["result"]["bridge_version"] != "godot-rust-bridge-0.1.0":
		fail("bridge metadata mismatch")
	print("WP010 metadata ok")
	var fixture = JSON.parse_string(FileAccess.get_file_as_string("res://bin/linux/debug/wp010-fixture.json"))
	var expected = JSON.parse_string(FileAccess.get_file_as_string("res://bin/linux/debug/wp010-pure-result.json"))
	print("WP010 fixture loaded")
	var created = parse_bridge(bridge.create_snapshot(int(fixture["seed"])))
	if not created.get("ok", false):
		fail("snapshot creation failed")
	print("WP010 snapshot created")
	var advanced = parse_bridge(bridge.advance_snapshot(created["result"]["snapshot"], int(fixture["elapsed_ms"])))
	if not advanced.get("ok", false):
		fail("snapshot advance failed")
	if advanced["result"] != expected:
		fail("Godot/native result diverged from pure Rust fixture result")
	print("WP010 advance compared")
	var repeated = parse_bridge(bridge.advance_snapshot(created["result"]["snapshot"], int(fixture["elapsed_ms"])))
	if repeated["result"] != advanced["result"]:
		fail("repeated stateless bridge call diverged")
	print("WP010 repeat compared")
	var batch = parse_bridge(bridge.execute_command_batch(created["result"]["snapshot"], "[]"))
	if not batch.get("ok", false) or batch["result"]["digest"] != created["result"]["digest"]:
		fail("empty command batch changed canonical state")
	print("WP010 batch compared")
	var digest = parse_bridge(bridge.digest_snapshot(created["result"]["snapshot"]))
	if not digest.get("ok", false) or digest["result"]["digest"] != created["result"]["digest"]:
		fail("digest operation mismatch")
	print("WP010 digest compared")
	var facility = parse_bridge(bridge.facility_summary())
	if not facility.get("ok", false):
		fail("facility summary failed")
	print("WP010 facility compared")
	var malformed = parse_bridge(bridge.advance_snapshot("{", 1))
	if malformed.get("ok", true) or malformed["error"]["code"] != "bridge.invalid_input":
		fail("malformed snapshot was not rejected")
	print("WP010 malformed compared")
	get_tree().quit(0)
