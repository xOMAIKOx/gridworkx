extends Control

const SCENARIO_ID := "scenario.opening.aggregate"
const FACILITY_ID := "facility.aggregate_plant_fixture"
const FEED_COMPONENT := "component.feed_conveyor"
const OUTPUT_COMPONENT := "component.output_conveyor"
const FEED_FAULT := "fault.instance.opening.feed_motor"
const OUTPUT_FAULT := "fault.instance.opening.output_belt"
const FINISHED_RESOURCE := "resource.finished_aggregate"
const FINISHED_GRADE := "grade.aggregate.standard"

var bridge: GridworksSimBridgeClient
var snapshot := ""
var initial_digest := ""
var action_number := 0
var status_label: Label
var facility_label: Label
var evidence_label: Label
var inventory_label: Label
var advisor_label: Label
var error_label: Label
var action_buttons: Array[Button] = []

func _ready() -> void:
	_build_ui()
	bridge = GridworksSimBridgeClient.new()
	if not bridge.is_available():
		_show_error("Native simulation bridge unavailable: " + JSON.stringify(bridge.availability_error()))
		return
	_restart_scenario()

func _build_ui() -> void:
	set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	var scroll := ScrollContainer.new()
	scroll.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	add_child(scroll)
	var content := VBoxContainer.new()
	content.custom_minimum_size = Vector2(360, 760)
	content.add_theme_constant_override("separation", 10)
	scroll.add_child(content)
	var title := Label.new()
	title.text = "Opening Aggregate Plant"
	title.add_theme_font_size_override("font_size", 26)
	content.add_child(title)
	status_label = _section(content, "Plant status")
	facility_label = _section(content, "Process schematic")
	evidence_label = _section(content, "Evidence and diagnosis")
	inventory_label = _section(content, "Inventory and throughput")
	advisor_label = _section(content, "Advisor")
	error_label = _section(content, "Errors")
	var actions := GridContainer.new()
	actions.columns = 2
	content.add_child(actions)
	_add_button(actions, "Inspect feed conveyor", _inspect_feed)
	_add_button(actions, "Inspect output belt", _inspect_output)
	_add_button(actions, "Diagnose feed seizure", _diagnose_feed)
	_add_button(actions, "Repair output belt", _repair_output)
	_add_button(actions, "Repair feed conveyor", _repair_feed)
	_add_button(actions, "Run aggregate production", _produce)
	_add_button(actions, "Restart scenario", _restart_scenario)

func _section(parent: Control, heading: String) -> Label:
	var panel := VBoxContainer.new()
	panel.add_theme_constant_override("separation", 3)
	parent.add_child(panel)
	var title := Label.new()
	title.text = heading
	title.add_theme_font_size_override("font_size", 17)
	panel.add_child(title)
	var value := Label.new()
	value.autowrap_mode = TextServer.AUTOWRAP_WORD_SMART
	panel.add_child(value)
	return value

func _add_button(parent: Control, label: String, callback: Callable) -> void:
	var button := Button.new()
	button.text = label
	button.custom_minimum_size = Vector2(170, 48)
	button.pressed.connect(callback)
	parent.add_child(button)
	action_buttons.append(button)

func _show_error(message: String) -> void:
	if error_label:
		error_label.text = message
	if status_label:
		status_label.text = "NON-AUTHORITATIVE ERROR — retry or restart"

func _clear_error() -> void:
	if error_label:
		error_label.text = "None"

func _restart_scenario() -> void:
	if not bridge or not bridge.is_available():
		return
	var response := bridge.create_scenario(SCENARIO_ID, 1234)
	if not response.get("ok", false):
		_show_error(JSON.stringify(response))
		return
	snapshot = response["result"]["snapshot"]
	initial_digest = response["result"]["digest"]
	action_number = 0
	_clear_error()
	_render()

func _command(command_type: String, payload: Dictionary, label: String) -> void:
	if snapshot.is_empty():
		_show_error("Scenario has not been created")
		return
	action_number += 1
	var command := {
		"command_id": "opening." + label + "." + str(action_number),
		"command_type": command_type,
		"schema_version": "schema-0.1.0",
		"rules_version": "rules-0.1.0",
		"effective_time_ms": int(_state().get("operational_time_ms", 0)),
		"idempotency_key": "opening." + label + "." + str(action_number),
		"payload": payload,
	}
	var response := bridge.execute_command_batch(snapshot, JSON.stringify([command]))
	if not response.get("ok", false):
		_show_error(JSON.stringify(response))
		return
	snapshot = response["result"]["snapshot"]
	_clear_error()
	_render()

func _failure_command(action: String, data: Dictionary, command_type: String, label: String) -> void:
	_command(command_type, {"failure": {"action": action, "data": data}}, label)

func _inspect_feed() -> void:
	_failure_command("Inspect", {"evidence_id": "evidence.opening.feed", "component_id": FEED_COMPONENT, "diagnostic_capability_bps": 10_000}, "failure.inspect", "inspect.feed")

func _inspect_output() -> void:
	_failure_command("Inspect", {"evidence_id": "evidence.opening.output", "component_id": OUTPUT_COMPONENT, "diagnostic_capability_bps": 10_000}, "failure.inspect", "inspect.output")

func _diagnose_feed() -> void:
	_failure_command("Diagnose", {"diagnosis_id": "diagnosis.opening.feed", "component_id": FEED_COMPONENT, "candidate_fault_type_id": "fault.motor_bearing_seizure", "fault_instance_id": null, "evidence_ids": ["evidence.opening.feed"], "diagnostic_capability_bps": 10_000}, "failure.diagnose", "diagnose.feed")

func _repair_output() -> void:
	_failure_command("Intervene", {"intervention_id": "intervention.opening.belt", "component_id": OUTPUT_COMPONENT, "kind": "Repair", "fault_instance_id": OUTPUT_FAULT, "derate_bps": null, "bypass": null}, "failure.intervene", "repair.belt")

func _repair_feed() -> void:
	_failure_command("Intervene", {"intervention_id": "intervention.opening.motor", "component_id": FEED_COMPONENT, "kind": "Repair", "fault_instance_id": FEED_FAULT, "derate_bps": null, "bypass": null}, "failure.intervene", "repair.motor")

func _produce() -> void:
	_command("material.production", {"material": {"ExecuteProduction": {"run_id": "run.opening.first." + str(action_number + 1), "recipe_id": "recipe.aggregate_crush", "facility_id": FACILITY_ID, "input_sources": [{"resource_id": "resource.raw_feed", "store_id": "inventory.aggregate_feed"}, {"resource_id": "resource.limestone", "store_id": "inventory.aggregate_feed"}], "output_inventory_id": "inventory.aggregate_finished", "requested_runs": 100}}}, "produce")

func _state() -> Dictionary:
	var parsed = JSON.parse_string(snapshot)
	return parsed.get("state", {}) if typeof(parsed) == TYPE_DICTIONARY else {}

func _material_quantity(state: Dictionary, inventory_id: String, resource_id: String, grade_id: String) -> int:
	for inventory in state.get("material", {}).get("inventories", []):
		if inventory.get("inventory_id", "") != inventory_id:
			continue
		for lot in inventory.get("lots", []):
			if lot.get("resource_id", "") == resource_id and lot.get("grade_id", "") == grade_id:
				return int(lot.get("quantity", 0))
	return 0

func _render() -> void:
	var state := _state()
	var evaluation := bridge.evaluate_facility(snapshot, FACILITY_ID)
	if not evaluation.get("ok", false):
		_show_error(JSON.stringify(evaluation))
		return
	var capacity := int(evaluation["result"].get("effective_capacity", 0))
	var state_name := str(evaluation["result"].get("operational_state", "Unknown"))
	var faults := []
	for fault in state.get("failure", {}).get("faults", []):
		faults.append(str(fault.get("fault_type_id", "")) + " / " + str(fault.get("status", "")))
	var evidence := []
	for item in state.get("failure", {}).get("evidence", []):
		evidence.append(str(item.get("component_id", "")) + ": " + str(item.get("observation", "")))
	var finished := _material_quantity(state, "inventory.aggregate_finished", FINISHED_RESOURCE, FINISHED_GRADE)
	status_label.text = "State: %s | Effective capacity: %d%% | Digest: %s" % [state_name, capacity, str(bridge.digest_snapshot(snapshot).get("result", {}).get("digest", "" )).left(12)]
	facility_label.text = "Feed stockpile → Feed conveyor → Crusher → Screen → Output conveyor → Finished stockpile\nCritical feed seizure and non-critical belt wear are canonical faults."
	evidence_label.text = "Faults:\n" + ("\n".join(faults) if faults.size() > 0 else "Inspect components to reveal evidence") + "\nEvidence:\n" + ("\n".join(evidence) if evidence.size() > 0 else "None")
	inventory_label.text = "Raw feed and limestone are canonical inputs.\nFinished aggregate: %d\nFirst-output completion: %s" % [finished, "YES" if capacity > 0 and finished > 0 else "NO"]
	advisor_label.text = "Inspect before spending. Visible wear is not necessarily the immediate blocker. Restore a small useful throughput, then observe real material output."

func _normalize_command(command: Dictionary) -> Dictionary:
	var normalized := command.duplicate(true)
	for key in ["effective_time_ms"]:
		if normalized.has(key):
			normalized[key] = int(normalized[key])
	var failure_data: Dictionary = normalized.get("payload", {}).get("failure", {}).get("data", {})
	for key in ["diagnostic_capability_bps", "severity_bps", "derate_bps"]:
		if failure_data.has(key) and failure_data[key] != null:
			failure_data[key] = int(failure_data[key])
	var material_data: Dictionary = normalized.get("payload", {}).get("material", {}).get("ExecuteProduction", {})
	if material_data.has("requested_runs"):
		material_data["requested_runs"] = int(material_data["requested_runs"])
	return normalized

func run_headless_golden(fixture: Dictionary) -> Dictionary:
	var response := bridge.create_scenario(str(fixture["scenario_id"]), int(fixture["seed"]))
	if not response.get("ok", false):
		return response
	var current: String = response["result"]["snapshot"]
	var event_types: Array[String] = []
	var partial_recovery_digest := ""
	for index in range(fixture["commands"].size()):
		var command = _normalize_command(fixture["commands"][index])
		var result := bridge.execute_command_batch(current, JSON.stringify([command]))
		if not result.get("ok", false):
			return result
		current = result["result"]["snapshot"]
		if index == 4:
			partial_recovery_digest = bridge.digest_snapshot(current)["result"]["digest"]
		for event in result["result"].get("events", []):
			event_types.append(str(event.get("event_type", "")))
	var evaluation: Dictionary = bridge.evaluate_facility(current, FACILITY_ID)
	var state: Dictionary = JSON.parse_string(current)
	var finished: int = _material_quantity(state.get("state", {}), "inventory.aggregate_finished", FINISHED_RESOURCE, FINISHED_GRADE)
	return {"ok": true, "initial_digest": response["result"]["digest"], "partial_recovery_digest": partial_recovery_digest, "final_digest": bridge.digest_snapshot(current)["result"]["digest"], "effective_capacity": int(evaluation["result"]["effective_capacity"]), "finished_aggregate_quantity": finished, "event_types": event_types}
