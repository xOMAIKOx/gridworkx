//! Versioned, deterministic simulation boundary foundation.
//!
//! WP-001 intentionally defines contracts and validation only. Gameplay rules belong to WP-002+
//! and must remain in this crate when implemented.

pub const SCHEMA_VERSION: &str = "schema-0.1.0";
pub const RULES_VERSION: &str = "rules-0.1.0";

#[derive(Debug, Clone, PartialEq, Eq)]
pub struct CommandEnvelope {
    pub command_id: String,
    pub command_type: String,
    pub payload_digest: String,
}

#[derive(Debug, Clone, PartialEq, Eq)]
pub struct CommandBatch {
    pub schema_version: String,
    pub rules_version: String,
    pub authoritative_timestamp_ms: i64,
    pub commands: Vec<CommandEnvelope>,
}

#[derive(Debug, Clone, PartialEq, Eq)]
pub struct SimulationSnapshot {
    pub schema_version: String,
    pub rules_version: String,
    pub operational_time_ms: i64,
    pub state_digest: String,
}

#[derive(Debug, Clone, PartialEq, Eq)]
pub enum BoundaryError {
    EmptyVersion(&'static str),
    NegativeTimestamp,
    EmptyCommandId,
    EmptyCommandType,
    EmptyPayloadDigest,
    EmptyStateDigest,
}

impl CommandBatch {
    pub fn validate(&self) -> Result<(), BoundaryError> {
        if self.schema_version.is_empty() {
            return Err(BoundaryError::EmptyVersion("schema_version"));
        }
        if self.rules_version.is_empty() {
            return Err(BoundaryError::EmptyVersion("rules_version"));
        }
        if self.authoritative_timestamp_ms < 0 {
            return Err(BoundaryError::NegativeTimestamp);
        }
        for command in &self.commands {
            if command.command_id.is_empty() {
                return Err(BoundaryError::EmptyCommandId);
            }
            if command.command_type.is_empty() {
                return Err(BoundaryError::EmptyCommandType);
            }
            if command.payload_digest.is_empty() {
                return Err(BoundaryError::EmptyPayloadDigest);
            }
        }
        Ok(())
    }
}

impl SimulationSnapshot {
    pub fn validate(&self) -> Result<(), BoundaryError> {
        if self.schema_version.is_empty() {
            return Err(BoundaryError::EmptyVersion("schema_version"));
        }
        if self.rules_version.is_empty() {
            return Err(BoundaryError::EmptyVersion("rules_version"));
        }
        if self.operational_time_ms < 0 {
            return Err(BoundaryError::NegativeTimestamp);
        }
        if self.state_digest.is_empty() {
            return Err(BoundaryError::EmptyStateDigest);
        }
        Ok(())
    }
}

/// Produces a stable non-cryptographic boundary digest without reading wall time or external state.
/// Cryptographic replay evidence is a later work-package concern.
pub fn deterministic_digest(snapshot: &SimulationSnapshot) -> String {
    let mut hash: u64 = 0xcbf29ce484222325;
    for byte in snapshot.schema_version.as_bytes() {
        hash = fnv1a_step(hash, *byte);
    }
    hash = fnv1a_step(hash, 0);
    for byte in snapshot.rules_version.as_bytes() {
        hash = fnv1a_step(hash, *byte);
    }
    hash = fnv1a_step(hash, 0);
    for byte in snapshot.operational_time_ms.to_le_bytes() {
        hash = fnv1a_step(hash, byte);
    }
    hash = fnv1a_step(hash, 0);
    for byte in snapshot.state_digest.as_bytes() {
        hash = fnv1a_step(hash, *byte);
    }
    format!("{hash:016x}")
}

fn fnv1a_step(hash: u64, byte: u8) -> u64 {
    (hash ^ u64::from(byte)).wrapping_mul(0x100000001b3)
}

#[cfg(test)]
mod tests {
    use super::*;

    fn snapshot() -> SimulationSnapshot {
        SimulationSnapshot {
            schema_version: SCHEMA_VERSION.to_owned(),
            rules_version: RULES_VERSION.to_owned(),
            operational_time_ms: 3_600_000,
            state_digest: "state-placeholder".to_owned(),
        }
    }

    #[test]
    fn boundary_snapshot_validates_without_a_clock() {
        assert_eq!(snapshot().validate(), Ok(()));
    }

    #[test]
    fn digest_is_repeatable_for_identical_inputs() {
        let first = deterministic_digest(&snapshot());
        let second = deterministic_digest(&snapshot());
        assert_eq!(first, second);
        assert_eq!(first.len(), 16);
    }

    #[test]
    fn invalid_command_batch_is_rejected() {
        let batch = CommandBatch {
            schema_version: SCHEMA_VERSION.to_owned(),
            rules_version: RULES_VERSION.to_owned(),
            authoritative_timestamp_ms: 0,
            commands: vec![CommandEnvelope {
                command_id: String::new(),
                command_type: "boundary.placeholder".to_owned(),
                payload_digest: "payload-placeholder".to_owned(),
            }],
        };
        assert_eq!(batch.validate(), Err(BoundaryError::EmptyCommandId));
    }
}
