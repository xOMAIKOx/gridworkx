use crate::{KernelError, SimulationState, RULES_VERSION, SCHEMA_VERSION};
use serde::{Deserialize, Serialize};

#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize)]
#[serde(deny_unknown_fields)]
pub struct SnapshotRecord {
    pub snapshot_id: String,
    pub owner_ref: Option<String>,
    pub schema_version: String,
    pub rules_version: String,
    pub kernel_revision: String,
    pub operational_time_ms: u64,
    pub state_digest: String,
    pub payload: String,
    pub previous_snapshot_id: Option<String>,
}

impl SnapshotRecord {
    pub fn from_state(
        snapshot_id: impl Into<String>,
        state: &SimulationState,
        previous_snapshot_id: Option<String>,
    ) -> Result<Self, KernelError> {
        let payload = state.to_json()?;
        Ok(Self {
            snapshot_id: snapshot_id.into(),
            owner_ref: None,
            schema_version: SCHEMA_VERSION.to_owned(),
            rules_version: RULES_VERSION.to_owned(),
            kernel_revision: state.kernel_revision.clone(),
            operational_time_ms: state.operational_time_ms,
            state_digest: state.digest()?,
            payload,
            previous_snapshot_id,
        })
    }

    pub fn validate(&self) -> Result<SimulationState, KernelError> {
        if self.schema_version != SCHEMA_VERSION || self.rules_version != RULES_VERSION {
            return Err(KernelError::InvalidState(
                "persistence record version is unsupported".to_owned(),
            ));
        }
        let state = SimulationState::from_json(&self.payload)?;
        if state.kernel_revision != self.kernel_revision
            || state.operational_time_ms != self.operational_time_ms
            || state.digest()? != self.state_digest
        {
            return Err(KernelError::InvalidState(
                "persistence snapshot integrity mismatch".to_owned(),
            ));
        }
        Ok(state)
    }
}

#[derive(Debug, Clone, Copy, PartialEq, Eq, Serialize, Deserialize)]
#[serde(deny_unknown_fields)]
pub enum ReceiptStatus {
    Accepted,
    Rejected,
}

#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize)]
#[serde(deny_unknown_fields)]
pub struct CommandReceipt {
    pub command_id: String,
    pub idempotency_key: String,
    pub command_type: String,
    pub status: ReceiptStatus,
    pub state_digest_before: Option<String>,
    pub state_digest_after: Option<String>,
    pub effective_time_ms: u64,
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::{Command, SimulationState};

    #[test]
    fn snapshot_record_round_trip_preserves_integrity() {
        let state = SimulationState::new(41).unwrap();
        let record = SnapshotRecord::from_state("snapshot.initial", &state, None).unwrap();
        assert_eq!(record.validate().unwrap(), state);
        let mut tampered = record.clone();
        tampered.payload = tampered.payload.replace("\"register\":0", "\"register\":1");
        assert!(matches!(
            tampered.validate(),
            Err(KernelError::InvalidState(_))
        ));
    }

    #[test]
    fn identical_state_digests_are_valid_for_distinct_snapshot_ids() {
        let state = SimulationState::new(41).unwrap();
        let first = SnapshotRecord::from_state("snapshot.first", &state, None).unwrap();
        let second = SnapshotRecord::from_state("snapshot.second", &state, None).unwrap();
        assert_eq!(first.state_digest, second.state_digest);
        assert_ne!(first.snapshot_id, second.snapshot_id);
    }

    #[test]
    fn receipt_carries_durable_command_identity_without_state_history() {
        let state = SimulationState::new(41).unwrap();
        let command = Command::adjust_register("command.receipt", "idempotency.receipt", 0, 1);
        let receipt = CommandReceipt {
            command_id: command.command_id,
            idempotency_key: command.idempotency_key,
            command_type: command.command_type,
            status: ReceiptStatus::Accepted,
            state_digest_before: Some(state.digest().unwrap()),
            state_digest_after: None,
            effective_time_ms: 0,
        };
        assert_eq!(receipt.status, ReceiptStatus::Accepted);
        assert!(state.executed_commands.is_empty());
    }
}
