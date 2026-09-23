#![allow(elided_lifetimes_in_associated_constant)]

use godot::prelude::*;
use gridworks_sim::{
    advance, aggregate_plant_fixture, execute, Command, KernelError, SimulationState,
    KERNEL_REVISION, RNG_ALGORITHM, RULES_VERSION, SCHEMA_VERSION,
};
use serde_json::{json, Value};

const BRIDGE_VERSION: &str = "godot-rust-bridge-0.1.0";
const MAX_SNAPSHOT_BYTES: usize = 4 * 1024 * 1024;
const MAX_BATCH_BYTES: usize = 1024 * 1024;

fn error(code: &str, message: impl Into<String>) -> GString {
    serde_json::to_string(&json!({"ok": false, "error": {"code": code, "message": message.into()}}))
        .unwrap_or_else(|_| {
            "{\"ok\":false,\"error\":{\"code\":\"bridge.internal\",\"message\":\"bridge error\"}}"
                .to_owned()
        })
        .into()
}
fn success(operation: &str, result: Value) -> GString {
    serde_json::to_string(&json!({"ok": true, "operation": operation, "result": result}))
        .unwrap_or_else(|_| {
            "{\"ok\":false,\"error\":{\"code\":\"bridge.internal\",\"message\":\"bridge error\"}}"
                .to_owned()
        })
        .into()
}
fn map_error(err: KernelError) -> GString {
    match err {
        KernelError::VersionMismatch { .. } => error(
            "bridge.version_mismatch",
            "snapshot or command version is unsupported",
        ),
        KernelError::InvalidTimeMovement { .. } => error(
            "bridge.invalid_time",
            "authoritative time cannot move backwards",
        ),
        KernelError::DuplicateCommand { .. } => {
            error("bridge.duplicate_command", "command was already executed")
        }
        KernelError::MalformedCommand(_) | KernelError::Deserialization(_) => error(
            "bridge.invalid_input",
            "boundary JSON or command is invalid",
        ),
        KernelError::Serialization(_) => error(
            "bridge.serialization",
            "canonical result could not be serialized",
        ),
        _ => error(
            "bridge.simulation_error",
            "canonical simulation operation failed",
        ),
    }
}
fn state_from_json(snapshot: &str) -> Result<SimulationState, GString> {
    if snapshot.len() > MAX_SNAPSHOT_BYTES {
        return Err(error(
            "bridge.payload_too_large",
            "snapshot exceeds the 4 MiB bridge limit",
        ));
    }
    SimulationState::from_json(snapshot).map_err(map_error)
}
fn transition_result(transition: gridworks_sim::Transition) -> Result<Value, GString> {
    Ok(
        json!({"snapshot": transition.resulting_state.to_json().map_err(map_error)?, "digest": transition.resulting_digest, "events": transition.events}),
    )
}

struct GridworksGodotExtension;
#[gdextension]
unsafe impl ExtensionLibrary for GridworksGodotExtension {}

#[derive(GodotClass)]
#[class(init, base=Node)]
pub struct GridworksSimBridge {
    #[base]
    base: Base<Node>,
}
