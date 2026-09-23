extends Node

func _ready() -> void:
	var controller = load("res://scripts/opening_aggregate_controller.gd").new()
	add_child(controller)
	await get_tree().process_frame
	await get_tree().process_frame
	var fixture = JSON.parse_string(FileAccess.get_file_as_string("res://bin/linux/debug/wp011-fixture.json"))
	if typeof(fixture) != TYPE_DICTIONARY:
		push_error("WP011 fixture is invalid")
		get_tree().quit(1)
		return
	var result = controller.run_headless_golden(fixture)
	if not result.get("ok", false):
		push_error("WP011 golden execution failed: " + JSON.stringify(result))
		get_tree().quit(1)
		return
	var expected = fixture["expected"]
	var checks := [
		[result["initial_digest"] == expected["initial_digest"], "initial digest"],
		[result["partial_recovery_digest"] == expected["partial_recovery_digest"], "partial recovery digest"],
		[result["final_digest"] == expected["final_digest"], "final digest"],
		[result["effective_capacity"] == int(expected["effective_capacity"]), "effective capacity"],
		[result["finished_aggregate_quantity"] == int(expected["finished_aggregate_quantity"]), "finished aggregate"],
		[result["event_types"] == expected["event_types"], "event semantics"],
	]
	for check in checks:
		if not check[0]:
			push_error("WP011 golden mismatch: " + check[1])
			get_tree().quit(1)
			return
	print("WP011 opening aggregate golden passed")
	get_tree().quit(0)
