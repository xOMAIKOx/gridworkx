extends SceneTree

var controller

func _initialize() -> void:
	controller = load("res://scripts/opening_aggregate_controller.gd").new()
	root.add_child(controller)
	await _frames(3)
	await _capture_size(Vector2i(390, 844), "wp011-opening-portrait.png")
	controller.inspect_component("component.feed_conveyor")
	await _frames(2)
	await _capture("wp011-inspected-feed.png")
	controller.select_diagnosis("fault.motor_bearing_seizure")
	controller.repair_output_belt()
	controller.repair_feed_conveyor()
	await _frames(2)
	await _capture("wp011-partial-recovery.png")
	controller.produce_first_output()
	await _frames(2)
	await _capture("wp011-first-output.png")
	controller.restart_scenario()
	await _frames(2)
	await _capture_size(Vector2i(1440, 900), "wp011-opening-wide.png")
	quit()

func _frames(count: int) -> void:
	for _index in range(count):
		await process_frame

func _capture_size(size: Vector2i, name: String) -> void:
	root.get_viewport().size = size
	await _frames(2)
	await _capture(name)

func _capture(name: String) -> void:
	var image := root.get_viewport().get_texture().get_image()
	image.save_png("user://" + name)
