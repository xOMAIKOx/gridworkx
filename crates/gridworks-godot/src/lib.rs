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

#[derive(GodotClass)]
#[class(base=Node)]
pub struct GridworksSimBridge {
    base: Base<Node>,
}

#[godot_api]
impl INode for GridworksSimBridge {
    fn init(base: Base<Node>) -> Self {
        Self { base }
    }
}

#[godot_api]
impl GridworksSimBridge {
    #[func]
    fn bridge_metadata(&self) -> GString {
        success(
            "bridge.metadata",
            json!({"bridge_version":BRIDGE_VERSION,"schema_version":SCHEMA_VERSION,"rules_version":RULES_VERSION,"kernel_revision":KERNEL_REVISION,"rng_algorithm":RNG_ALGORITHM}),
        )
    }

    #[func]
    fn create_snapshot(&self, seed: i64) -> GString {
        if seed <= 0 {
            return error("bridge.invalid_seed", "seed must be positive");
        }
        match SimulationState::new(seed as u64).and_then(|state| {
            state
                .to_json()
                .map(|snapshot| (snapshot, state.digest().unwrap_or_default()))
        }) {
            Ok((snapshot, digest)) => success(
                "snapshot.create",
                json!({"snapshot":snapshot,"digest":digest,"events":[]}),
            ),
            Err(err) => map_error(err),
        }
    }

    #[func]
    fn advance_snapshot(&self, snapshot: GString, elapsed_ms: i64) -> GString {
        if elapsed_ms < 0 {
            return error("bridge.invalid_time", "elapsed time must be non-negative");
        }
        let mut state = match state_from_json(&snapshot.to_string()) {
            Ok(v) => v,
            Err(e) => return e,
        };
        match advance(&mut state, elapsed_ms as u64).and_then(|t| {
            transition_result(t).map_err(|_| KernelError::Serialization("bridge result".to_owned()))
        }) {
            Ok(result) => success("snapshot.advance", result),
            Err(err) => map_error(err),
        }
    }

    #[func]
    fn execute_command_batch(&self, snapshot: GString, batch: GString) -> GString {
        if batch.to_string().len() > MAX_BATCH_BYTES {
            return error(
                "bridge.payload_too_large",
                "command batch exceeds the 1 MiB bridge limit",
            );
        }
        let mut state = match state_from_json(&snapshot.to_string()) {
            Ok(v) => v,
            Err(e) => return e,
        };
        let commands: Vec<Command> = match serde_json::from_str(&batch.to_string()) {
            Ok(v) => v,
            Err(_) => return error("bridge.invalid_input", "command batch JSON is invalid"),
        };
        let mut events = Vec::new();
        for command in commands {
            match execute(&mut state, &command) {
                Ok(t) => events.extend(t.events),
                Err(e) => return map_error(e),
            }
        }
        match state
            .digest()
            .and_then(|digest| state.to_json().map(|snapshot| (snapshot, digest)))
        {
            Ok((snapshot, digest)) => success(
                "command.batch",
                json!({"snapshot":snapshot,"digest":digest,"events":events}),
            ),
            Err(e) => map_error(e),
        }
    }

    #[func]
    fn digest_snapshot(&self, snapshot: GString) -> GString {
        let state = match state_from_json(&snapshot.to_string()) {
            Ok(v) => v,
            Err(e) => return e,
        };
        match state.digest() {
            Ok(d) => success("snapshot.digest", json!({"digest":d})),
            Err(e) => map_error(e),
        }
    }

    #[func]
    fn facility_summary(&self) -> GString {
        match aggregate_plant_fixture().evaluate() {
            Ok(v) => success(
                "facility.summary",
                serde_json::to_value(v).unwrap_or_else(|_| json!({})),
            ),
            Err(_) => error("bridge.facility_error", "facility evaluation failed"),
        }
    }
}

struct GridworksGodotExtension;
#[gdextension]
unsafe impl ExtensionLibrary for GridworksGodotExtension {}
