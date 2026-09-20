extends Node

const PRESENTATION_LAYERS := [
	"presentation.world_map",
	"presentation.portfolio_map",
	"presentation.site_view",
	"presentation.operational_detail",
]

func _ready() -> void:
	assert(PRESENTATION_LAYERS.size() == 4)
