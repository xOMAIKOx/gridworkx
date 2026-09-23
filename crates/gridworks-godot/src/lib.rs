#![allow(elided_lifetimes_in_associated_constant)]

use godot::prelude::*;
use gridworks_sim::{
    advance, advance_to, execute, FailureError, KernelError, SimulationState, Transition,
    KERNEL_REVISION, RNG_ALGORITHM, RULES_VERSION, SCHEMA_VERSION,
};
use serde_json::{json, Value};
use std::panic::{catch_unwind, AssertUnwindSafe};

const BRIDGE_VERSION: &str = "godot-rust-bridge-0.1.0";
const MAX_SNAPSHOT_BYTES: usize = 4 * 1024 * 1024;
const MAX_BATCH_BYTES: usize = 1024 * 1024;
const MAX_FACILITY_ID_BYTES: usize = 128;

fn encode(value: Value) -> GString {
    serde_json::to_string(&value)
        .unwrap_or_else(|_| {
            "{\"ok\":false,\"bridge_version\":\"godot-rust-bridge-0.1.0\",\"schema_version\":\"schema-0.1.0\",\"rules_version\":\"rules-0.1.0\",\"operation\":\"bridge.error\",\"error\":{\"code\":\"bridge.internal\",\"message\":\"bridge response encoding failed\"}}".to_owned()
        })
        .into()
}

fn error(code: &str, message: &str) -> GString {
    encode(json!({
        "ok": false,
        "bridge_version": BRIDGE_VERSION,
        "schema_version": SCHEMA_VERSION,
        "rules_version": RULES_VERSION,
        "operation": "bridge.error",
        "error": {"code": code, "message": message},
    }))
}

fn success(operation: &str, result: Value) -> GString {
    encode(json!({
        "ok": true,
        "bridge_version": BRIDGE_VERSION,
        "schema_version": SCHEMA_VERSION,
        "rules_version": RULES_VERSION,
        "operation": operation,
        "result": result,
    }))
}

fn guarded<F>(operation: &str, function: F) -> GString
where
    F: FnOnce() -> Result<Value, GString>,
{
    match catch_unwind(AssertUnwindSafe(function)) {
        Ok(Ok(result)) => success(operation, result),
        Ok(Err(response)) => response,
        Err(_) => error("bridge.internal", "bridge operation failed safely"),
    }
}

fn map_error(error_value: KernelError) -> GString {
    match error_value {
        KernelError::InvalidSeed => error("bridge.invalid_seed", "seed must be positive"),
        KernelError::InvalidRngState => {
            error("bridge.invalid_rng", "deterministic RNG state is invalid")
        }
        KernelError::InvalidState(_) => {
            error("bridge.invalid_state", "simulation state is invalid")
        }
        KernelError::MalformedCommand(_) => {
            error("bridge.invalid_command", "command payload is malformed")
        }
        KernelError::UnsupportedCommandType(_) => {
            error("bridge.unsupported_command", "command type is unsupported")
        }
        KernelError::VersionMismatch { .. } => error(
            "bridge.version_mismatch",
            "schema or rules version is unsupported",
        ),
        KernelError::InvalidTimeMovement { .. } => error(
            "bridge.invalid_time",
            "authoritative time cannot move backwards or disagree with state",
        ),
        KernelError::DuplicateCommand { .. } => error(
            "bridge.duplicate_command",
            "command identity or idempotency key was already executed",
        ),
        KernelError::ArithmeticOverflow => error(
            "bridge.arithmetic",
            "simulation arithmetic overflowed safely",
        ),
        KernelError::EmptyChoiceSet => error(
            "bridge.invalid_command",
            "command choice set cannot be empty",
        ),
        KernelError::Serialization(_) => error(
            "bridge.serialization",
            "canonical result could not be serialized",
        ),
        KernelError::Deserialization(_) => error(
            "bridge.invalid_json",
            "boundary JSON is malformed or contains unknown fields",
        ),
        KernelError::Failure(FailureError::UnknownComponent(_)) => error(
            "bridge.unknown_facility",
            "facility is not present in the supplied state",
        ),
        KernelError::Failure(_) => error(
            "bridge.validation",
            "facility or failure validation rejected the operation",
        ),
        KernelError::Material(_) | KernelError::Ledger(_) => error(
            "bridge.validation",
            "canonical validation rejected the operation",
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

fn transition_result(transition: Transition) -> Result<Value, GString> {
    Ok(json!({
        "snapshot": transition.resulting_state.to_json().map_err(map_error)?,
        "digest": transition.resulting_digest,
        "events": transition.events,
    }))
}

fn validate_facility_id(facility_id: &str) -> Result<(), GString> {
    if facility_id.is_empty() || facility_id.len() > MAX_FACILITY_ID_BYTES {
        return Err(error(
            "bridge.invalid_facility",
            "facility_id must be between 1 and 128 UTF-8 bytes",
        ));
    }
    Ok(())
}

#[derive(GodotClass)]
#[class(init)]
pub struct GridworksSimBridge {}

#[godot_api]
impl GridworksSimBridge {
    #[func]
    fn bridge_metadata(&self) -> GString {
        guarded("bridge.metadata", || {
            Ok(json!({
                "bridge_version": BRIDGE_VERSION,
                "schema_version": SCHEMA_VERSION,
                "rules_version": RULES_VERSION,
                "kernel_revision": KERNEL_REVISION,
                "rng_algorithm": RNG_ALGORITHM,
            }))
        })
    }

    #[func]
    fn create_snapshot(&self, seed: i64) -> GString {
        guarded("snapshot.create", || {
            if seed <= 0 {
                return Err(error("bridge.invalid_seed", "seed must be positive"));
            }
            let state = SimulationState::new(seed as u64).map_err(map_error)?;
            let snapshot = state.to_json().map_err(map_error)?;
            let digest = state.digest().map_err(map_error)?;
            Ok(json!({"snapshot": snapshot, "digest": digest, "events": []}))
        })
    }

    #[func]
    fn advance_snapshot(&self, snapshot: GString, elapsed_ms: i64) -> GString {
        guarded("snapshot.advance", || {
            if elapsed_ms < 0 {
                return Err(error(
                    "bridge.invalid_time",
                    "elapsed time must be non-negative",
                ));
            }
            let mut state = state_from_json(&snapshot.to_string())?;
            advance(&mut state, elapsed_ms as u64)
                .map_err(map_error)
                .and_then(transition_result)
        })
    }

    #[func]
    fn advance_to_snapshot(&self, snapshot: GString, authoritative_time_ms: i64) -> GString {
        guarded("snapshot.advance_to", || {
            if authoritative_time_ms < 0 {
                return Err(error(
                    "bridge.invalid_time",
                    "authoritative time must be non-negative",
                ));
            }
            let mut state = state_from_json(&snapshot.to_string())?;
            advance_to(&mut state, authoritative_time_ms as u64)
                .map_err(map_error)
                .and_then(transition_result)
        })
    }

    #[func]
    fn execute_command_batch(&self, snapshot: GString, batch: GString) -> GString {
        guarded("command.batch", || {
            let batch_text = batch.to_string();
            if batch_text.len() > MAX_BATCH_BYTES {
                return Err(error(
                    "bridge.payload_too_large",
                    "command batch exceeds the 1 MiB bridge limit",
                ));
            }
            let mut state = state_from_json(&snapshot.to_string())?;
            let commands: Vec<gridworks_sim::Command> =
                serde_json::from_str(&batch_text).map_err(|_| {
                    error(
                        "bridge.invalid_command",
                        "command batch JSON is malformed or contains unknown fields",
                    )
                })?;
            let mut events = Vec::new();
            for command in commands {
                let transition = execute(&mut state, &command).map_err(map_error)?;
                events.extend(transition.events);
            }
            Ok(json!({
                "snapshot": state.to_json().map_err(map_error)?,
                "digest": state.digest().map_err(map_error)?,
                "events": events,
            }))
        })
    }

    #[func]
    fn digest_snapshot(&self, snapshot: GString) -> GString {
        guarded("snapshot.digest", || {
            let state = state_from_json(&snapshot.to_string())?;
            Ok(json!({"digest": state.digest().map_err(map_error)?}))
        })
    }

    #[func]
    fn validate_digest(&self, snapshot: GString, expected_digest: GString) -> GString {
        guarded("snapshot.digest.validate", || {
            let state = state_from_json(&snapshot.to_string())?;
            let digest = state.digest().map_err(map_error)?;
            if digest != expected_digest.to_string() {
                return Err(error(
                    "bridge.digest_mismatch",
                    "expected digest does not match canonical digest",
                ));
            }
            Ok(json!({"digest": digest, "valid": true}))
        })
    }

    #[func]
    fn evaluate_facility(&self, snapshot: GString, facility_id: GString) -> GString {
        guarded("facility.evaluate", || {
            let facility_id = facility_id.to_string();
            validate_facility_id(&facility_id)?;
            let state = state_from_json(&snapshot.to_string())?;
            let evaluation = state.evaluate_facility(&facility_id).map_err(map_error)?;
            serde_json::to_value(evaluation).map_err(|_| {
                error(
                    "bridge.internal",
                    "facility evaluation could not be encoded",
                )
            })
        })
    }
}

struct GridworksGodotExtension;
#[gdextension]
unsafe impl ExtensionLibrary for GridworksGodotExtension {}
