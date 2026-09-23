extends Control

const SCENARIO_ID := "scenario.opening.aggregate"
const FACILITY_ID := "facility.aggregate_plant_fixture"
const FINISHED_RESOURCE := "resource.finished_aggregate"
const FINISHED_GRADE := "grade.aggregate.standard"
const COMPONENTS := [
	{"id": "component.feed_conveyor", "label": "Feed conveyor"},
	{"id": "component.output_conveyor", "label": "Output conveyor"},
	{"id": "component.crusher", "label": "Crusher"},
	{"id": "component.screen", "label": "Screen"},
]
const DIAGNOSIS_CANDIDATES := [
	{"id": "fault.motor_bearing_seizure", "label": "Motor bearing seizure"},
	{"id": "fault.worn_belt", "label": "Worn belt"},
	{"id": "fault.screen_blockage", "label": "Screen blockage"},
]

var bridge: GridworksSimBridgeClient
var snapshot := ""
var initial_digest := ""
var selected_component_id := "component.feed_conveyor"
var presentation_state := "opening.intro"
var action_number := 0
var production_observed := false
var last_production_event: Dictionary = {}
var advisor_visible := true
var status_label: Label
var facility_label: Label
var selected_label: Label
var evidence_label: Label
var inventory_label: Label
var advisor_label: Label
var error_label: Label
var component_buttons: Array[Button] = []
var diagnosis_buttons: Array[Button] = []
var action_buttons: Array[Button] = []

func _ready() -> void:
	_build_ui()
	bridge = GridworksSimBridgeClient.new()
	if not bridge.is_available():
		_show_error("Native simulation bridge unavailable: " + JSON.stringify(bridge.availability_error()))
		return
	restart_scenario()

func _build_ui() -> void:
	set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	var scroll := ScrollContainer.new()
	scroll.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	add_child(scroll)
	var content := VBoxContainer.new()
	content.custom_minimum_size = Vector2(360, 860)
	content.add_theme_constant_override("separation", 10)
	scroll.add_child(content)
	var title := Label.new()
	title.text = "Opening Aggregate Plant"
	title.add_theme_font_size_override("font_size", 26)
	content.add_child(title)
	status_label = _section(content, "Onboarding state")
	facility_label = _section(content, "Process schematic")
	selected_label = _section(content, "Selected component")
	evidence_label = _section(content, "Evidence and diagnosis")
	inventory_label = _section(content, "Inventory and throughput")
	advisor_label = _section(content, "Advisor")
	error_label = _section(content, "Errors")
	var components_title := Label.new()
	components_title.text = "Select a component"
	components_title.add_theme_font_size_override("font_size", 17)
	content.add_child(components_title)
	var components := GridContainer.new()
	components.columns = 2
	content.add_child(components)
	for component in COMPONENTS:
		var button := Button.new()
		button.text = component["label"]
		button.custom_minimum_size = Vector2(170, 48)
		button.pressed.connect(_select_component.bind(component["id"]))
		components.add_child(button)
		component_buttons.append(button)
	var inspect_title := Label.new()
	inspect_title.text = "Canonical inspection"
	inspect_title.add_theme_font_size_override("font_size", 17)
	content.add_child(inspect_title)
	var inspect := Button.new()
	inspect.text = "Inspect selected component"
	inspect.custom_minimum_size = Vector2(350, 48)
	inspect.pressed.connect(inspect_selected_component)
	content.add_child(inspect)
	var diagnosis_title := Label.new()
	diagnosis_title.text = "Choose a diagnosis candidate"
	diagnosis_title.add_theme_font_size_override("font_size", 17)
	content.add_child(diagnosis_title)
	var diagnoses := GridContainer.new()
	diagnoses.columns = 1
	content.add_child(diagnoses)
	for candidate in DIAGNOSIS_CANDIDATES:
		var button := Button.new()
		button.text = candidate["label"]
		button.custom_minimum_size = Vector2(350, 48)
		button.pressed.connect(select_diagnosis.bind(candidate["id"]))
		diagnoses.add_child(button)
		diagnosis_buttons.append(button)
	var actions := GridContainer.new()
	actions.columns = 2
	content.add_child(actions)
	_add_button(actions, "Repair output belt", repair_output_belt)
	_add_button(actions, "Repair feed conveyor", repair_feed_conveyor)
	_add_button(actions, "Run aggregate production", produce_first_output)
	_add_button(actions, "Restart scenario", restart_scenario)
	var advisor_control := Button.new()
	advisor_control.text = "Dismiss / show advisor"
	advisor_control.custom_minimum_size = Vector2(350, 48)
	advisor_control.pressed.connect(_toggle_advisor)
	content.add_child(advisor_control)

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

func restart_scenario() -> Dictionary:
	if not bridge or not bridge.is_available():
		return {"ok": false, "error": {"code": "bridge.unavailable"}}
	var response: Dictionary = bridge.create_scenario(SCENARIO_ID, 1234)
	if not response.get("ok", false):
		_show_error(JSON.stringify(response))
		return response
	snapshot = response["result"]["snapshot"]
	initial_digest = response["result"]["digest"]
	action_number = 0
	production_observed = false
	last_production_event = {}
	presentation_state = "opening.intro"
	_clear_error()
	_render()
	return response

func _state() -> Dictionary:
	var parsed: Variant = JSON.parse_string(snapshot)
	return parsed.get("state", {}) if typeof(parsed) == TYPE_DICTIONARY else {}

func _component_label(component_id: String) -> String:
	for component in COMPONENTS:
		if component["id"] == component_id:
			return component["label"]
	return component_id

func _current_time() -> int:
	return int(_state().get("operational_time_ms", 0))

func _command(command_type: String, payload: Dictionary, label: String, command_id_override: String = "", idempotency_override: String = "") -> Dictionary:
	if snapshot.is_empty():
		_show_error("Scenario has not been created")
		return {"ok": false, "error": {"code": "scenario.unavailable"}}
	action_number += 1
	var command_id := command_id_override if not command_id_override.is_empty() else "opening." + label + "." + str(action_number)
	var idempotency_key := idempotency_override if not idempotency_override.is_empty() else command_id
	var command := {
		"command_id": command_id,
		"command_type": command_type,
		"schema_version": "schema-0.1.0",
		"rules_version": "rules-0.1.0",
		"effective_time_ms": _current_time(),
		"idempotency_key": idempotency_key,
		"payload": payload,
	}
	var response: Dictionary = bridge.execute_command_batch(snapshot, JSON.stringify([command]))
	if not response.get("ok", false):
		_show_error(JSON.stringify(response))
		return response
	var previous_snapshot := snapshot
	snapshot = response["result"]["snapshot"]
	for event in response["result"].get("events", []):
		if event.get("event_type", "") == "material" and event.get("payload", {}).has("ProductionExecuted"):
			last_production_event = event["payload"]["ProductionExecuted"]
			if int(last_production_event.get("accepted_runs", 0)) > 0:
				production_observed = true
	_clear_error()
	_render()
	if snapshot.is_empty():
		snapshot = previous_snapshot
	return response

func _failure_command(action: String, data: Dictionary, command_type: String, label: String, command_id_override: String = "", idempotency_override: String = "") -> Dictionary:
	return _command(command_type, {"failure": {"action": action, "data": data}}, label, command_id_override, idempotency_override)

func _select_component(component_id: String) -> void:
	selected_component_id = component_id
	_render()

func _component_key(component_id: String) -> String:
	match component_id:
		"component.feed_conveyor": return "feed"
		"component.output_conveyor": return "output"
		"component.crusher": return "crusher"
		"component.screen": return "screen"
	return component_id.trim_prefix("component.")

func inspect_selected_component() -> Dictionary:
	var evidence_id := "evidence.opening." + _component_key(selected_component_id)
	return _failure_command("Inspect", {"evidence_id": evidence_id, "component_id": selected_component_id, "diagnostic_capability_bps": 10_000}, "failure.inspect", "inspect." + selected_component_id.trim_prefix("component."))

func inspect_component(component_id: String, command_id_override: String = "", idempotency_override: String = "") -> Dictionary:
	_select_component(component_id)
	var evidence_id := "evidence.opening." + _component_key(component_id)
	return _failure_command("Inspect", {"evidence_id": evidence_id, "component_id": component_id, "diagnostic_capability_bps": 10_000}, "failure.inspect", "inspect." + component_id.trim_prefix("component."), command_id_override, idempotency_override)

func select_diagnosis(candidate_fault_type_id: String, command_id_override: String = "", idempotency_override: String = "", diagnosis_id_override: String = "") -> Dictionary:
	var evidence_ids: Array[String] = []
	for item in _state().get("failure", {}).get("evidence", []):
		if item.get("component_id", "") == selected_component_id:
			evidence_ids.append(item.get("evidence_id", ""))
	if evidence_ids.is_empty():
		return {"ok": false, "error": {"code": "opening.evidence_required"}}
	var diagnosis_id := diagnosis_id_override if not diagnosis_id_override.is_empty() else "diagnosis.opening." + _component_key(selected_component_id)
	return _failure_command("Diagnose", {"diagnosis_id": diagnosis_id, "component_id": selected_component_id, "candidate_fault_type_id": candidate_fault_type_id, "fault_instance_id": null, "evidence_ids": evidence_ids, "diagnostic_capability_bps": 10_000}, "failure.diagnose", "diagnose." + selected_component_id.trim_prefix("component."), command_id_override, idempotency_override)

func repair_output_belt(command_id_override: String = "", idempotency_override: String = "") -> Dictionary:
	return _failure_command("Intervene", {"intervention_id": "intervention.opening.belt", "component_id": "component.output_conveyor", "kind": "Repair", "fault_instance_id": "fault.instance.opening.output_belt", "derate_bps": null, "bypass": null}, "failure.intervene", "repair.belt", command_id_override, idempotency_override)

func repair_feed_conveyor(command_id_override: String = "", idempotency_override: String = "") -> Dictionary:
	return _failure_command("Intervene", {"intervention_id": "intervention.opening.motor", "component_id": "component.feed_conveyor", "kind": "Repair", "fault_instance_id": "fault.instance.opening.feed_motor", "derate_bps": null, "bypass": null}, "failure.intervene", "repair.motor", command_id_override, idempotency_override)

func produce_first_output(command_id_override: String = "", idempotency_override: String = "") -> Dictionary:
	var run_id := "run.opening.first." + str(action_number + 1)
	return _command("material.production", {"material": {"ExecuteProduction": {"run_id": run_id, "recipe_id": "recipe.aggregate_crush", "facility_id": FACILITY_ID, "input_sources": [{"resource_id": "resource.raw_feed", "store_id": "inventory.aggregate_feed"}, {"resource_id": "resource.limestone", "store_id": "inventory.aggregate_feed"}], "output_inventory_id": "inventory.aggregate_finished", "requested_runs": 100}}}, "produce", command_id_override, idempotency_override)

func _toggle_advisor() -> void:
	advisor_visible = not advisor_visible
	_render()

func _material_quantity(state: Dictionary, inventory_id: String, resource_id: String, grade_id: String) -> int:
	for inventory in state.get("material", {}).get("inventories", []):
		if inventory.get("inventory_id", "") != inventory_id:
			continue
		for lot in inventory.get("lots", []):
			if lot.get("resource_id", "") == resource_id and lot.get("grade_id", "") == grade_id:
				return int(lot.get("quantity", 0))
	return 0

func _derive_presentation_state(state: Dictionary, capacity: int, finished: int) -> String:
	var evidence_count: int = state.get("failure", {}).get("evidence", []).size()
	var diagnosis_count: int = state.get("failure", {}).get("diagnoses", []).size()
	if evidence_count == 0:
		return "opening.intro"
	if evidence_count < 4:
		return "opening.inspect"
	if diagnosis_count == 0:
		return "opening.diagnose"
	if capacity == 0:
		return "opening.intervene"
	if not production_observed:
		return "opening.produce"
	if _completion_predicate(state, capacity, finished):
		return "opening.complete"
	return "opening.produce"

func _completion_predicate(state: Dictionary, capacity: int, finished: int) -> bool:
	var feed_resolved := false
	for fault in state.get("failure", {}).get("faults", []):
		if fault.get("fault_instance_id", "") == "fault.instance.opening.feed_motor":
			feed_resolved = fault.get("status", "") == "Resolved"
	return feed_resolved and capacity > 0 and production_observed and finished > 0

func _render() -> void:
	var state := _state()
	var evaluation: Dictionary = bridge.evaluate_facility(snapshot, FACILITY_ID)
	if not evaluation.get("ok", false):
		_show_error(JSON.stringify(evaluation))
		return
	var capacity := int(evaluation["result"].get("effective_capacity", 0))
	var finished := _material_quantity(state, "inventory.aggregate_finished", FINISHED_RESOURCE, FINISHED_GRADE)
	presentation_state = _derive_presentation_state(state, capacity, finished)
	status_label.text = "State: %s | Effective capacity: %d%% | Completion: %s" % [presentation_state, capacity, "YES" if _completion_predicate(state, capacity, finished) else "NO"]
	facility_label.text = "Feed stockpile → Feed conveyor → Crusher → Screen → Output conveyor → Finished stockpile"
	selected_label.text = "%s (%s)\nCanonical facility state: %s" % [_component_label(selected_component_id), selected_component_id, str(evaluation["result"].get("operational_state", "Unknown"))]
	var evidence_lines: Array[String] = []
	for item in state.get("failure", {}).get("evidence", []):
		if item.get("component_id", "") == selected_component_id:
			evidence_lines.append(str(item.get("observation", "")) + " / " + str(item.get("symptom_id", "")))
	var diagnosis_lines: Array[String] = []
	for diagnosis in state.get("failure", {}).get("diagnoses", []):
		if diagnosis.get("component_id", "") == selected_component_id:
			diagnosis_lines.append(str(diagnosis.get("candidate_fault_type_id", "")) + " / " + str(diagnosis.get("status", "")) + " / " + str(diagnosis.get("confidence_bps", 0)) + " bps")
	evidence_label.text = "Observed evidence:\n" + ("\n".join(evidence_lines) if not evidence_lines.is_empty() else "Inspect this component to reveal evidence.") + "\nDiagnosis:\n" + ("\n".join(diagnosis_lines) if not diagnosis_lines.is_empty() else "No diagnosis recorded.")
	var raw := _material_quantity(state, "inventory.aggregate_feed", "resource.raw_feed", "grade.raw.standard")
	var limestone := _material_quantity(state, "inventory.aggregate_feed", "resource.limestone", "grade.limestone.standard")
	inventory_label.text = "Raw feed: %d\nLimestone: %d\nFinished aggregate: %d\nLast accepted runs: %d" % [raw, limestone, finished, int(last_production_event.get("accepted_runs", 0))]
	advisor_label.visible = advisor_visible
	advisor_label.text = _advisor_text()
	_refresh_action_buttons(state, capacity)

func _advisor_text() -> String:
	if not advisor_visible:
		return "Advisor dismissed."
	match presentation_state:
		"opening.intro": return "Inspect before spending. Select a component and inspect it."
		"opening.inspect": return "Compare evidence. Visible wear is not automatically the immediate blocker."
		"opening.diagnose": return "Choose the diagnosis that the observed evidence supports."
		"opening.intervene": return "Choose a recoverable intervention; imperfect repair can leave the plant constrained."
		"opening.produce": return "Useful partial throughput creates options. Run the canonical process to observe output."
		"opening.complete": return "The plant produced real aggregate. Use the new output to plan your next decision."
	return "Inspect before spending."

func _refresh_action_buttons(state: Dictionary, capacity: int) -> void:
	for button in diagnosis_buttons:
		button.disabled = not state.get("failure", {}).get("evidence", []).any(func(item): return item.get("component_id", "") == selected_component_id)
	var evidence_count: int = state.get("failure", {}).get("evidence", []).size()
	var diagnosis_count: int = state.get("failure", {}).get("diagnoses", []).size()
	for button in action_buttons:
		if button.text == "Repair output belt":
			button.disabled = evidence_count == 0
		elif button.text == "Repair feed conveyor":
			button.disabled = diagnosis_count == 0
		elif button.text == "Run aggregate production":
			button.disabled = capacity == 0

func _normalize_command(command: Dictionary) -> Dictionary:
	var normalized := command.duplicate(true)
	if normalized.has("effective_time_ms"):
		normalized["effective_time_ms"] = int(normalized["effective_time_ms"])
	var failure_data: Dictionary = normalized.get("payload", {}).get("failure", {}).get("data", {})
	for key in ["diagnostic_capability_bps", "severity_bps", "derate_bps"]:
		if failure_data.has(key) and failure_data[key] != null:
			failure_data[key] = int(failure_data[key])
	var material_data: Dictionary = normalized.get("payload", {}).get("material", {}).get("ExecuteProduction", {})
	if material_data.has("requested_runs"):
		material_data["requested_runs"] = int(material_data["requested_runs"])
	return normalized

func _event_types(response: Dictionary) -> Array[String]:
	var types: Array[String] = []
	for event in response.get("result", {}).get("events", []):
		types.append(str(event.get("event_type", "")))
	return types

func run_onboarding_contract_checks() -> Dictionary:
	var response := restart_scenario()
	if not response.get("ok", false):
		return response
	var initial_text := status_label.text + facility_label.text + selected_label.text + evidence_label.text
	if initial_text.contains("fault.") or initial_text.contains("seizure"):
		return {"ok": false, "error": {"code": "opening.answer_leak"}}
	for component_id in ["component.feed_conveyor", "component.output_conveyor", "component.crusher", "component.screen"]:
		var inspected := inspect_component(component_id)
		if not inspected.get("ok", false):
			return inspected
	if _state().get("failure", {}).get("evidence", []).size() < 4:
		return {"ok": false, "error": {"code": "opening.evidence_missing"}}
	_select_component("component.feed_conveyor")
	var diagnosis := select_diagnosis("fault.motor_bearing_seizure")
	if not diagnosis.get("ok", false):
		return diagnosis
	var belt := repair_output_belt()
	if not belt.get("ok", false) or int(bridge.evaluate_facility(snapshot, FACILITY_ID)["result"]["effective_capacity"]) != 0:
		return {"ok": false, "error": {"code": "opening.wrong_choice_not_recoverable"}}
	var motor := repair_feed_conveyor()
	if not motor.get("ok", false) or int(bridge.evaluate_facility(snapshot, FACILITY_ID)["result"]["effective_capacity"]) != 3:
		return {"ok": false, "error": {"code": "opening.partial_recovery_failed"}}
	var production := produce_first_output()
	if not production.get("ok", false) or int(last_production_event.get("accepted_runs", 0)) != 3 or not production_observed:
		return {"ok": false, "error": {"code": "opening.production_event_missing"}}
	var finished := _material_quantity(_state(), "inventory.aggregate_finished", FINISHED_RESOURCE, FINISHED_GRADE)
	if finished != 24 or not _completion_predicate(_state(), 3, finished):
		return {"ok": false, "error": {"code": "opening.completion_invalid"}}
	var before_bad := snapshot
	var bad := bridge.digest_snapshot("{")
	if bad.get("ok", true) or snapshot != before_bad:
		return {"ok": false, "error": {"code": "opening.error_corrupted_snapshot"}}
	var restarted := restart_scenario()
	if not restarted.get("ok", false) or restarted["result"]["digest"] != initial_digest:
		return {"ok": false, "error": {"code": "opening.restart_not_deterministic"}}
	return {"ok": true}

func run_headless_golden(fixture: Dictionary) -> Dictionary:
	var response: Dictionary = restart_scenario()
	if not response.get("ok", false):
		return response
	var event_types: Array[String] = []
	var initial: String = response["result"]["digest"]
	var responses: Array[Dictionary] = []
	responses.append(inspect_component("component.feed_conveyor", "opening.inspect.feed", "opening.inspect.feed.v1"))
	responses.append(inspect_component("component.output_conveyor", "opening.inspect.output", "opening.inspect.output.v1"))
	_select_component("component.feed_conveyor")
	responses.append(select_diagnosis("fault.motor_bearing_seizure", "opening.diagnose.feed", "opening.diagnose.feed.v1", "diagnosis.opening.feed"))
	responses.append(repair_output_belt("opening.repair.belt", "opening.repair.belt.v1"))
	responses.append(repair_feed_conveyor("opening.repair.motor", "opening.repair.motor.v1"))
	var partial: String = bridge.digest_snapshot(snapshot).get("result", {}).get("digest", "")
	responses.append(produce_first_output("opening.production.first", "opening.production.first.v1"))
	for response_item in responses:
		if not response_item.get("ok", false):
			return response_item
		event_types.append_array(_event_types(response_item))
	var state := _state()
	var evaluation: Dictionary = bridge.evaluate_facility(snapshot, FACILITY_ID)
	var finished: int = _material_quantity(state, "inventory.aggregate_finished", FINISHED_RESOURCE, FINISHED_GRADE)
	return {"ok": true, "initial_digest": initial, "partial_recovery_digest": partial, "final_digest": bridge.digest_snapshot(snapshot)["result"]["digest"], "effective_capacity": int(evaluation["result"]["effective_capacity"]), "accepted_runs": int(last_production_event.get("accepted_runs", 0)), "finished_aggregate_quantity": finished, "event_types": event_types, "completed": _completion_predicate(state, int(evaluation["result"]["effective_capacity"]), finished)}
