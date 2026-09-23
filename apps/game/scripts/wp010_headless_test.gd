extends Node

func fail(message: String) -> void:
	push_error(message)
	get_tree().quit(1)

func check(condition: bool, message: String) -> bool:
	if condition:
		return true
	fail(message)
	return false

func expect_error(response: Dictionary, code: String, context: String) -> bool:
	if not check(not response.get("ok", true), context + " unexpectedly succeeded"):
		return false
	if not check(response.get("error", {}).get("code", "") == code, context + " returned the wrong error code"):
		return false
	if not check(not response.has("result"), context + " exposed a partial result"):
		return false
	return true

func _ready() -> void:
	await get_tree().process_frame
	await get_tree().process_frame
	print("WP010 ready")
	_run()

func _run() -> void:
	var bridge = GridworksSimBridgeClient.new()
	if not check(bridge.is_available(), "GDExtension bridge wrapper is unavailable"):
		return
	var metadata = bridge.bridge_metadata()
	if not check(metadata.get("ok", false), "bridge metadata failed"):
		return
	if not check(metadata["result"]["bridge_version"] == "godot-rust-bridge-0.1.0", "bridge metadata version mismatch"):
		return
	print("WP010 wrapper and metadata present")

	var fixture = JSON.parse_string(FileAccess.get_file_as_string("res://bin/linux/debug/wp010-fixture.json"))
	var expected = JSON.parse_string(FileAccess.get_file_as_string("res://bin/linux/debug/wp010-pure-result.json"))
	var facility_expected = JSON.parse_string(FileAccess.get_file_as_string("res://bin/linux/debug/wp010-facility-result.json"))
	if not check(typeof(fixture) == TYPE_DICTIONARY and typeof(expected) == TYPE_DICTIONARY, "fixture artifacts are invalid"):
		return

	var created = bridge.create_snapshot(int(fixture["seed"]))
	if not check(created.get("ok", false), "snapshot creation failed"):
		return
	var advanced = bridge.advance_snapshot(created["result"]["snapshot"], int(fixture["elapsed_ms"]))
	if not check(advanced.get("ok", false), "elapsed snapshot advance failed"):
		return
	if not check(advanced["result"] == expected["after_elapsed"], "elapsed result diverged from committed Rust golden"):
		return
	var repeated = bridge.advance_snapshot(created["result"]["snapshot"], int(fixture["elapsed_ms"]))
	if not check(repeated == advanced, "repeated identical input produced a different result"):
		return
	print("WP010 elapsed golden and repeat compared")

	var batch_json = JSON.stringify(fixture["commands"])
	var batched = bridge.execute_command_batch(advanced["result"]["snapshot"], batch_json)
	if not check(batched.get("ok", false), "non-empty command batch failed: " + JSON.stringify(batched)):
		return
	if not check(batched["result"]["digest"] == fixture["expected"]["after_command_digest"], "command batch digest diverged from committed golden"):
		return
	if not check(batched["result"]["snapshot"] == expected["after_commands"]["snapshot"], "command batch snapshot diverged from committed Rust result"):
		return
	if not check(batched["result"]["events"] == expected["after_commands"]["events"], "command batch events diverged from committed Rust result"):
		return
	print("WP010 command batch golden compared")

	var duplicate = bridge.execute_command_batch(batched["result"]["snapshot"], batch_json)
	if not expect_error(duplicate, "bridge.duplicate_command", "duplicate command batch"):
		return
	var failing_commands = [fixture["commands"][0], {"command_id": "wp010.bad", "command_type": "unsupported.command", "schema_version": "schema-0.1.0", "rules_version": "rules-0.1.0", "effective_time_ms": 5000, "idempotency_key": "wp010.bad.v1", "payload": {"adjust_register": {"delta": 1}}}]
	var atomic_failure = bridge.execute_command_batch(advanced["result"]["snapshot"], JSON.stringify(failing_commands))
	if not expect_error(atomic_failure, "bridge.unsupported_command", "later command atomicity"):
		return
	print("WP010 duplicate and atomicity regressions compared")

	var final = bridge.advance_to_snapshot(batched["result"]["snapshot"], int(fixture["authoritative_time_ms"]))
	if not check(final.get("ok", false), "authoritative-time advance failed"):
		return
	if not check(final["result"] == expected["final"], "authoritative-time result diverged from committed Rust golden"):
		return
	var backwards = bridge.advance_to_snapshot(final["result"]["snapshot"], int(fixture["authoritative_time_ms"]) - 1)
	if not expect_error(backwards, "bridge.invalid_time", "backward authoritative time"):
		return
	print("WP010 authoritative-time regressions compared")

	var digest = bridge.digest_snapshot(created["result"]["snapshot"])
	if not check(digest.get("ok", false) and digest["result"]["digest"] == created["result"]["digest"], "digest calculation mismatch"):
		return
	var valid_digest = bridge.validate_digest(created["result"]["snapshot"], created["result"]["digest"])
	if not check(valid_digest.get("ok", false) and valid_digest["result"]["valid"], "correct digest validation failed"):
		return
	var wrong_digest = bridge.validate_digest(created["result"]["snapshot"], "0".repeat(64))
	if not expect_error(wrong_digest, "bridge.digest_mismatch", "wrong digest validation"):
		return
	print("WP010 digest regressions compared")

	var facility = bridge.evaluate_facility(facility_expected["snapshot"], fixture["facility_id"])
	if not check(facility.get("ok", false) and facility["result"] == facility_expected["evaluation"], "supplied-state facility evaluation diverged"):
		return
	var unknown_facility = bridge.evaluate_facility(facility_expected["snapshot"], "facility.unknown")
	if not expect_error(unknown_facility, "bridge.unknown_facility", "unknown facility"):
		return
	var oversized_facility_id = bridge.evaluate_facility(facility_expected["snapshot"], "f".repeat(129))
	if not expect_error(oversized_facility_id, "bridge.invalid_facility", "oversized facility id"):
		return
	print("WP010 supplied-state facility regressions compared")

	var unknown_field_snapshot = JSON.parse_string(created["result"]["snapshot"])
	unknown_field_snapshot["unknown_field"] = true
	if not expect_error(bridge.digest_snapshot(JSON.stringify(unknown_field_snapshot)), "bridge.invalid_json", "snapshot unknown field"):
		return
	var unsupported_snapshot = JSON.parse_string(created["result"]["snapshot"])
	unsupported_snapshot["schema_version"] = "schema-unsupported"
	unsupported_snapshot["state"]["schema_version"] = "schema-unsupported"
	if not expect_error(bridge.digest_snapshot(JSON.stringify(unsupported_snapshot)), "bridge.version_mismatch", "unsupported schema"):
		return
	var invalid_rng_snapshot = JSON.parse_string(created["result"]["snapshot"])
	invalid_rng_snapshot["state"]["rng"]["algorithm"] = "rng-unsupported"
	if not expect_error(bridge.digest_snapshot(JSON.stringify(invalid_rng_snapshot)), "bridge.invalid_rng", "invalid RNG state"):
		return
	var unsupported_rules_snapshot = JSON.parse_string(created["result"]["snapshot"])
	unsupported_rules_snapshot["rules_version"] = "rules-unsupported"
	unsupported_rules_snapshot["state"]["rules_version"] = "rules-unsupported"
	if not expect_error(bridge.digest_snapshot(JSON.stringify(unsupported_rules_snapshot)), "bridge.version_mismatch", "unsupported rules"):
		return
	var unsupported_kernel_snapshot = JSON.parse_string(created["result"]["snapshot"])
	unsupported_kernel_snapshot["kernel_revision"] = "kernel-unsupported"
	unsupported_kernel_snapshot["state"]["kernel_revision"] = "kernel-unsupported"
	if not expect_error(bridge.digest_snapshot(JSON.stringify(unsupported_kernel_snapshot)), "bridge.invalid_state", "unsupported kernel"):
		return
	if not expect_error(bridge.digest_snapshot("x".repeat(4 * 1024 * 1024 + 1)), "bridge.payload_too_large", "oversized snapshot"):
		return
	if not expect_error(bridge.execute_command_batch(created["result"]["snapshot"], "[" + " ".repeat(1024 * 1024) + "]"), "bridge.payload_too_large", "oversized command batch"):
		return
	if not expect_error(bridge.digest_snapshot("{"), "bridge.invalid_json", "malformed snapshot"):
		return
	print("WP010 negative and panic-safety regressions compared")
	get_tree().quit(0)
