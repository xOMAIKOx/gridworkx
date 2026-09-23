class_name GridworksSimBridgeClient
extends RefCounted

const EXPECTED_BRIDGE_VERSION := "godot-rust-bridge-0.1.0"
const EXPECTED_SCHEMA_VERSION := "schema-0.1.0"
const EXPECTED_RULES_VERSION := "rules-0.1.0"

var _native: Object = null
var _availability_error: Dictionary = {}

func _init() -> void:
	if not ClassDB.class_exists("GridworksSimBridge"):
		_availability_error = _local_error("bridge.unavailable", "GridworksSimBridge is not registered")
		return
	_native = ClassDB.instantiate("GridworksSimBridge")
	if _native == null:
		_availability_error = _local_error("bridge.unavailable", "GridworksSimBridge could not be instantiated")

func is_available() -> bool:
	return _native != null and _availability_error.is_empty()

func availability_error() -> Dictionary:
	return _availability_error.duplicate(true)

func _local_error(code: String, message: String) -> Dictionary:
	return {
		"ok": false,
		"bridge_version": EXPECTED_BRIDGE_VERSION,
		"schema_version": EXPECTED_SCHEMA_VERSION,
		"rules_version": EXPECTED_RULES_VERSION,
		"operation": "bridge.wrapper",
		"error": {"code": code, "message": message},
	}

func _invoke(method: StringName, arguments: Array) -> Dictionary:
	if not is_available():
		return availability_error()
	if not _native.has_method(method):
		return _local_error("bridge.method_unavailable", "required bridge method is unavailable")
	var parsed = JSON.parse_string(_native.callv(method, arguments))
	if typeof(parsed) != TYPE_DICTIONARY:
		return _local_error("bridge.invalid_envelope", "bridge returned a non-object envelope")
	for field in ["bridge_version", "schema_version", "rules_version", "ok"]:
		if not parsed.has(field):
			return _local_error("bridge.invalid_envelope", "bridge envelope is missing a version field")
	if parsed["bridge_version"] != EXPECTED_BRIDGE_VERSION:
		return _local_error("bridge.version_mismatch", "bridge version is unsupported")
	if parsed["schema_version"] != EXPECTED_SCHEMA_VERSION or parsed["rules_version"] != EXPECTED_RULES_VERSION:
		return _local_error("bridge.version_mismatch", "schema or rules version is unsupported")
	if parsed["ok"] and (not parsed.has("operation") or not parsed.has("result")):
		return _local_error("bridge.invalid_envelope", "success envelope is incomplete")
	if not parsed["ok"] and (not parsed.has("operation") or not parsed.has("error")):
		return _local_error("bridge.invalid_envelope", "error envelope is incomplete")
	return parsed

func bridge_metadata() -> Dictionary:
	return _invoke(&"bridge_metadata", [])

func create_snapshot(seed: int) -> Dictionary:
	return _invoke(&"create_snapshot", [seed])

func advance_snapshot(snapshot: String, elapsed_ms: int) -> Dictionary:
	return _invoke(&"advance_snapshot", [snapshot, elapsed_ms])

func advance_to_snapshot(snapshot: String, authoritative_time_ms: int) -> Dictionary:
	return _invoke(&"advance_to_snapshot", [snapshot, authoritative_time_ms])

func execute_command_batch(snapshot: String, batch: String) -> Dictionary:
	return _invoke(&"execute_command_batch", [snapshot, batch])

func digest_snapshot(snapshot: String) -> Dictionary:
	return _invoke(&"digest_snapshot", [snapshot])

func validate_digest(snapshot: String, expected_digest: String) -> Dictionary:
	return _invoke(&"validate_digest", [snapshot, expected_digest])

func evaluate_facility(snapshot: String, facility_id: String) -> Dictionary:
	return _invoke(&"evaluate_facility", [snapshot, facility_id])
