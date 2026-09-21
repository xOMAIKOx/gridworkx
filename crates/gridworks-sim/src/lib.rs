use serde::{Deserialize, Serialize};
use sha2::{Digest, Sha256};

pub mod failure;
pub mod graph;

pub use failure::{
    aggregate_fault_definitions, ComponentConditionState, Diagnosis, DiagnosisStatus, EvidenceItem,
    EvidenceSource, FailureCommand, FailureError, FailureEvent, FailureState, FaultDefinition,
    FaultInstance, FaultStatus, InterventionKind, Observation, SymptomObservation,
};
pub use graph::{
    aggregate_plant_fixture, BottleneckEvidence, BottleneckReason, Component, DependencyEdge,
    DependencyType, Facility, FacilityContext, FacilityEvaluation, GraphError, OperationalState,
    System, SystemEvaluation, BASIS_POINTS_PER_WHOLE,
};

pub const SCHEMA_VERSION: &str = "schema-0.1.0";
pub const RULES_VERSION: &str = "rules-0.1.0";
pub const KERNEL_REVISION: &str = "kernel-0.2.0";
pub const RNG_ALGORITHM: &str = "xorshift64star-v1";
pub const SNAPSHOT_TYPE: &str = "simulation.snapshot";

#[derive(Debug, Clone, PartialEq, Eq)]
pub enum KernelError {
    InvalidSeed,
    InvalidRngState,
    InvalidState(String),
    MalformedCommand(String),
    UnsupportedCommandType(String),
    VersionMismatch {
        expected_schema: String,
        actual_schema: String,
        expected_rules: String,
        actual_rules: String,
    },
    InvalidTimeMovement {
        current_time_ms: u64,
        requested_time_ms: u64,
    },
    DuplicateCommand {
        command_id: String,
        idempotency_key: String,
    },
    ArithmeticOverflow,
    EmptyChoiceSet,
    Serialization(String),
    Deserialization(String),
    Failure(FailureError),
}

#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize)]
#[serde(deny_unknown_fields)]
pub struct DeterministicRng {
    pub algorithm: String,
    pub seed: u64,
    pub state: u64,
    pub draws: u64,
}

impl DeterministicRng {
    pub fn new(seed: u64) -> Result<Self, KernelError> {
        if seed == 0 {
            return Err(KernelError::InvalidSeed);
        }
        Ok(Self {
            algorithm: RNG_ALGORITHM.to_owned(),
            seed,
            state: seed,
            draws: 0,
        })
    }

    fn validate(&self) -> Result<(), KernelError> {
        if self.algorithm != RNG_ALGORITHM || self.seed == 0 || self.state == 0 {
            return Err(KernelError::InvalidRngState);
        }
        Ok(())
    }

    fn next_u64(&mut self) -> Result<u64, KernelError> {
        self.validate()?;
        let mut value = self.state;
        value ^= value >> 12;
        value ^= value << 25;
        value ^= value >> 27;
        self.state = value;
        self.draws = self
            .draws
            .checked_add(1)
            .ok_or(KernelError::ArithmeticOverflow)?;
        Ok(value.wrapping_mul(0x2545_f491_4f6c_dd1d))
    }

    fn choose_index(&mut self, length: usize) -> Result<usize, KernelError> {
        if length == 0 {
            return Err(KernelError::EmptyChoiceSet);
        }
        let length = u64::try_from(length).map_err(|_| KernelError::ArithmeticOverflow)?;
        let threshold = u64::MAX - (u64::MAX % length);
        loop {
            let value = self.next_u64()?;
            if value < threshold {
                return usize::try_from(value % length)
                    .map_err(|_| KernelError::ArithmeticOverflow);
            }
        }
    }
}

#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize)]
#[serde(deny_unknown_fields)]
pub enum ProofPayload {
    #[serde(rename = "adjust_register")]
    AdjustRegister { delta: i64 },
    #[serde(rename = "seeded_pulse")]
    SeededPulse { options: Vec<i64> },
    #[serde(rename = "failure")]
    Failure(FailureCommand),
}

#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize)]
#[serde(deny_unknown_fields)]
pub struct Command {
    pub command_id: String,
    pub command_type: String,
    pub schema_version: String,
    pub rules_version: String,
    pub effective_time_ms: u64,
    pub idempotency_key: String,
    pub payload: ProofPayload,
}

impl Command {
    pub fn adjust_register(
        command_id: impl Into<String>,
        idempotency_key: impl Into<String>,
        effective_time_ms: u64,
        delta: i64,
    ) -> Self {
        Self {
            command_id: command_id.into(),
            command_type: "proof.adjust_register".to_owned(),
            schema_version: SCHEMA_VERSION.to_owned(),
            rules_version: RULES_VERSION.to_owned(),
            effective_time_ms,
            idempotency_key: idempotency_key.into(),
            payload: ProofPayload::AdjustRegister { delta },
        }
    }

    pub fn seeded_pulse(
        command_id: impl Into<String>,
        idempotency_key: impl Into<String>,
        effective_time_ms: u64,
        options: Vec<i64>,
    ) -> Self {
        Self {
            command_id: command_id.into(),
            command_type: "proof.seeded_pulse".to_owned(),
            schema_version: SCHEMA_VERSION.to_owned(),
            rules_version: RULES_VERSION.to_owned(),
            effective_time_ms,
            idempotency_key: idempotency_key.into(),
            payload: ProofPayload::SeededPulse { options },
        }
    }

    fn validate(&self) -> Result<(), KernelError> {
        if self.command_id.is_empty() {
            return Err(KernelError::MalformedCommand(
                "command_id is required".to_owned(),
            ));
        }
        if self.command_type.is_empty() {
            return Err(KernelError::MalformedCommand(
                "command_type is required".to_owned(),
            ));
        }
        if self.idempotency_key.is_empty() {
            return Err(KernelError::MalformedCommand(
                "idempotency_key is required".to_owned(),
            ));
        }
        if self.schema_version != SCHEMA_VERSION || self.rules_version != RULES_VERSION {
            return Err(KernelError::VersionMismatch {
                expected_schema: SCHEMA_VERSION.to_owned(),
                actual_schema: self.schema_version.clone(),
                expected_rules: RULES_VERSION.to_owned(),
                actual_rules: self.rules_version.clone(),
            });
        }
        Ok(())
    }
}

#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize)]
#[serde(deny_unknown_fields)]
pub struct ProofState {
    pub register: i64,
    pub pulse_count: u64,
    pub last_pulse: Option<i64>,
}

#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize)]
#[serde(deny_unknown_fields)]
pub struct ExecutedCommand {
    pub command_id: String,
    pub idempotency_key: String,
}

#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize)]
#[serde(deny_unknown_fields)]
pub struct SimulationState {
    pub schema_version: String,
    pub rules_version: String,
    pub kernel_revision: String,
    pub operational_time_ms: u64,
    pub rng: DeterministicRng,
    pub proof: ProofState,
    pub executed_commands: Vec<ExecutedCommand>,
    #[serde(default)]
    pub facilities: Vec<Facility>,
    #[serde(default)]
    pub failure: FailureState,
}

impl SimulationState {
    pub fn new(seed: u64) -> Result<Self, KernelError> {
        let state = Self {
            schema_version: SCHEMA_VERSION.to_owned(),
            rules_version: RULES_VERSION.to_owned(),
            kernel_revision: KERNEL_REVISION.to_owned(),
            operational_time_ms: 0,
            rng: DeterministicRng::new(seed)?,
            proof: ProofState {
                register: 0,
                pulse_count: 0,
                last_pulse: None,
            },
            executed_commands: Vec::new(),
            facilities: Vec::new(),
            failure: FailureState::default(),
        };
        state.validate()?;
        Ok(state)
    }

    pub fn with_facilities(seed: u64, facilities: Vec<Facility>) -> Result<Self, KernelError> {
        let mut state = Self::new(seed)?;
        graph::validate_facilities(&facilities)
            .map_err(|error| KernelError::InvalidState(format!("facility graph: {error}")))?;
        state.facilities = facilities;
        Ok(state)
    }

    pub fn evaluate_facility(&self, facility_id: &str) -> Result<FacilityEvaluation, KernelError> {
        let facility = self
            .facilities
            .iter()
            .find(|facility| facility.facility_id == facility_id)
            .ok_or_else(|| {
                KernelError::Failure(FailureError::UnknownComponent(facility_id.to_owned()))
            })?;
        self.failure
            .evaluate_facility(facility)
            .map_err(KernelError::Failure)
    }

    pub fn validate(&self) -> Result<(), KernelError> {
        if self.schema_version != SCHEMA_VERSION {
            return Err(KernelError::InvalidState(
                "unsupported schema version".to_owned(),
            ));
        }
        if self.rules_version != RULES_VERSION {
            return Err(KernelError::InvalidState(
                "unsupported rules version".to_owned(),
            ));
        }
        if self.kernel_revision != KERNEL_REVISION {
            return Err(KernelError::InvalidState(
                "unsupported kernel revision".to_owned(),
            ));
        }
        self.rng.validate()?;
        graph::validate_facilities(&self.facilities)
            .map_err(|error| KernelError::InvalidState(format!("facility graph: {error}")))?;
        self.failure
            .validate_against_facilities(&self.facilities)
            .map_err(KernelError::Failure)?;
        for (index, current) in self.executed_commands.iter().enumerate() {
            if current.command_id.is_empty() || current.idempotency_key.is_empty() {
                return Err(KernelError::InvalidState(
                    "executed command identity is incomplete".to_owned(),
                ));
            }
            if self.executed_commands[..index].iter().any(|previous| {
                previous.command_id == current.command_id
                    || previous.idempotency_key == current.idempotency_key
            }) {
                return Err(KernelError::InvalidState(
                    "executed command identity is duplicated".to_owned(),
                ));
            }
        }
        Ok(())
    }

    pub fn canonicalized(&self) -> Self {
        let mut canonical = self.clone();
        canonical.facilities = self
            .facilities
            .iter()
            .map(Facility::canonicalized)
            .collect();
        canonical
            .facilities
            .sort_by(|left, right| left.facility_id.cmp(&right.facility_id));
        canonical.failure = self.failure.canonicalized();
        canonical
    }

    pub fn snapshot(&self) -> SimulationSnapshot {
        let state = self.canonicalized();
        SimulationSnapshot {
            snapshot_type: SNAPSHOT_TYPE.to_owned(),
            schema_version: state.schema_version.clone(),
            rules_version: state.rules_version.clone(),
            kernel_revision: state.kernel_revision.clone(),
            state,
        }
    }

    pub fn to_json(&self) -> Result<String, KernelError> {
        serde_json::to_string(&self.snapshot())
            .map_err(|error| KernelError::Serialization(error.to_string()))
    }

    pub fn from_json(input: &str) -> Result<Self, KernelError> {
        let snapshot: SimulationSnapshot = serde_json::from_str(input)
            .map_err(|error| KernelError::Deserialization(error.to_string()))?;
        snapshot.validate()?;
        Ok(snapshot.state)
    }

    pub fn digest(&self) -> Result<String, KernelError> {
        let bytes = serde_json::to_vec(&self.snapshot())
            .map_err(|error| KernelError::Serialization(error.to_string()))?;
        Ok(sha256_hex(&bytes))
    }
}

#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize)]
#[serde(deny_unknown_fields)]
pub struct SimulationSnapshot {
    pub snapshot_type: String,
    pub schema_version: String,
    pub rules_version: String,
    pub kernel_revision: String,
    pub state: SimulationState,
}

impl SimulationSnapshot {
    pub fn validate(&self) -> Result<(), KernelError> {
        if self.snapshot_type != SNAPSHOT_TYPE {
            return Err(KernelError::Deserialization(
                "unsupported snapshot type".to_owned(),
            ));
        }
        if self.schema_version != SCHEMA_VERSION || self.rules_version != RULES_VERSION {
            return Err(KernelError::VersionMismatch {
                expected_schema: SCHEMA_VERSION.to_owned(),
                actual_schema: self.schema_version.clone(),
                expected_rules: RULES_VERSION.to_owned(),
                actual_rules: self.rules_version.clone(),
            });
        }
        if self.kernel_revision != KERNEL_REVISION {
            return Err(KernelError::InvalidState(
                "unsupported snapshot kernel revision".to_owned(),
            ));
        }
        self.state.validate()?;
        if self.state.schema_version != self.schema_version
            || self.state.rules_version != self.rules_version
            || self.state.kernel_revision != self.kernel_revision
        {
            return Err(KernelError::InvalidState(
                "snapshot metadata does not match state".to_owned(),
            ));
        }
        Ok(())
    }
}

#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize)]
#[serde(deny_unknown_fields)]
#[serde(tag = "event_type", content = "payload")]
pub enum KernelEvent {
    #[serde(rename = "time.advanced")]
    TimeAdvanced { from_time_ms: u64, to_time_ms: u64 },
    #[serde(rename = "proof.register_adjusted")]
    RegisterAdjusted {
        command_id: String,
        delta: i64,
        register: i64,
    },
    #[serde(rename = "proof.seeded_pulse")]
    SeededPulse {
        command_id: String,
        chosen_index: usize,
        chosen_value: i64,
        register: i64,
    },
    #[serde(rename = "failure")]
    Failure(FailureEvent),
}

#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize)]
#[serde(deny_unknown_fields)]
pub struct Transition {
    pub events: Vec<KernelEvent>,
    pub resulting_state: SimulationState,
    pub resulting_digest: String,
}

pub fn execute(state: &mut SimulationState, command: &Command) -> Result<Transition, KernelError> {
    state.validate()?;
    command.validate()?;
    if command.effective_time_ms != state.operational_time_ms {
        return Err(KernelError::InvalidTimeMovement {
            current_time_ms: state.operational_time_ms,
            requested_time_ms: command.effective_time_ms,
        });
    }
    if state.executed_commands.iter().any(|previous| {
        previous.command_id == command.command_id
            || previous.idempotency_key == command.idempotency_key
    }) {
        return Err(KernelError::DuplicateCommand {
            command_id: command.command_id.clone(),
            idempotency_key: command.idempotency_key.clone(),
        });
    }

    let mut next = state.clone();
    let event = match (&command.command_type[..], &command.payload) {
        ("proof.adjust_register", ProofPayload::AdjustRegister { delta }) => {
            next.proof.register = next
                .proof
                .register
                .checked_add(*delta)
                .ok_or(KernelError::ArithmeticOverflow)?;
            KernelEvent::RegisterAdjusted {
                command_id: command.command_id.clone(),
                delta: *delta,
                register: next.proof.register,
            }
        }
        ("proof.adjust_register", _) => {
            return Err(KernelError::MalformedCommand(
                "adjust_register payload does not match command type".to_owned(),
            ))
        }
        ("proof.seeded_pulse", ProofPayload::SeededPulse { options }) => {
            if options.is_empty() {
                return Err(KernelError::EmptyChoiceSet);
            }
            if options.len() > 64 {
                return Err(KernelError::MalformedCommand(
                    "seeded pulse has too many options".to_owned(),
                ));
            }
            let chosen_index = next.rng.choose_index(options.len())?;
            let chosen_value = options[chosen_index];
            next.proof.register = next
                .proof
                .register
                .checked_add(chosen_value)
                .ok_or(KernelError::ArithmeticOverflow)?;
            next.proof.pulse_count = next
                .proof
                .pulse_count
                .checked_add(1)
                .ok_or(KernelError::ArithmeticOverflow)?;
            next.proof.last_pulse = Some(chosen_value);
            KernelEvent::SeededPulse {
                command_id: command.command_id.clone(),
                chosen_index,
                chosen_value,
                register: next.proof.register,
            }
        }
        ("proof.seeded_pulse", _) => {
            return Err(KernelError::MalformedCommand(
                "seeded_pulse payload does not match command type".to_owned(),
            ))
        }
        (failure_type, ProofPayload::Failure(action)) if failure_type.starts_with("failure.") => {
            let facilities = next.facilities.clone();
            next.failure
                .apply(
                    &facilities,
                    action,
                    command.effective_time_ms,
                    &next.rules_version,
                )
                .map(KernelEvent::Failure)
                .map_err(KernelError::Failure)?
        }
        (failure_type, ProofPayload::Failure(_)) if failure_type.starts_with("failure.") => {
            return Err(KernelError::MalformedCommand(
                "failure payload does not match command type".to_owned(),
            ))
        }
        (unsupported, _) => {
            return Err(KernelError::UnsupportedCommandType(unsupported.to_owned()))
        }
    };
    next.executed_commands.push(ExecutedCommand {
        command_id: command.command_id.clone(),
        idempotency_key: command.idempotency_key.clone(),
    });
    next.validate()?;
    let resulting_digest = next.digest()?;
    *state = next.clone();
    Ok(Transition {
        events: vec![event],
        resulting_state: next,
        resulting_digest,
    })
}

pub fn advance(
    state: &mut SimulationState,
    authoritative_elapsed_ms: u64,
) -> Result<Transition, KernelError> {
    state.validate()?;
    let requested_time_ms = state
        .operational_time_ms
        .checked_add(authoritative_elapsed_ms)
        .ok_or(KernelError::ArithmeticOverflow)?;
    let mut next = state.clone();
    let from_time_ms = next.operational_time_ms;
    next.operational_time_ms = requested_time_ms;
    next.validate()?;
    let resulting_digest = next.digest()?;
    *state = next.clone();
    Ok(Transition {
        events: vec![KernelEvent::TimeAdvanced {
            from_time_ms,
            to_time_ms: requested_time_ms,
        }],
        resulting_state: next,
        resulting_digest,
    })
}

pub fn advance_to(
    state: &mut SimulationState,
    authoritative_time_ms: u64,
) -> Result<Transition, KernelError> {
    if authoritative_time_ms < state.operational_time_ms {
        return Err(KernelError::InvalidTimeMovement {
            current_time_ms: state.operational_time_ms,
            requested_time_ms: authoritative_time_ms,
        });
    }
    advance(state, authoritative_time_ms - state.operational_time_ms)
}

fn sha256_hex(bytes: &[u8]) -> String {
    let digest = Sha256::digest(bytes);
    let mut output = String::with_capacity(digest.len() * 2);
    for byte in digest {
        output.push_str(&format!("{byte:02x}"));
    }
    output
}

#[cfg(test)]
mod tests {
    use super::*;

    fn command_stream(time_ms: u64) -> [Command; 3] {
        [
            Command::adjust_register("command-1", "idempotency-1", time_ms, 4),
            Command::seeded_pulse("command-2", "idempotency-2", time_ms, vec![-3, 2, 7]),
            Command::adjust_register("command-3", "idempotency-3", time_ms, -1),
        ]
    }

    fn run_stream(
        mut state: SimulationState,
        commands: &[Command],
    ) -> (SimulationState, Vec<KernelEvent>) {
        let mut events = Vec::new();
        events.extend(advance(&mut state, 1_000).unwrap().events);
        for command in commands {
            events.extend(execute(&mut state, command).unwrap().events);
        }
        events.extend(advance(&mut state, 2_000).unwrap().events);
        (state, events)
    }

    #[test]
    fn initial_state_is_versioned_and_valid() {
        let state = SimulationState::new(7).unwrap();
        assert_eq!(state.operational_time_ms, 0);
        assert_eq!(state.schema_version, SCHEMA_VERSION);
        assert_eq!(state.rules_version, RULES_VERSION);
        assert_eq!(state.kernel_revision, KERNEL_REVISION);
        assert_eq!(state.rng.algorithm, RNG_ALGORITHM);
        assert_eq!(state.validate(), Ok(()));
    }

    #[test]
    fn valid_command_changes_only_proof_state() {
        let mut state = SimulationState::new(7).unwrap();
        advance(&mut state, 1_000).unwrap();
        let transition =
            execute(&mut state, &Command::adjust_register("id", "key", 1_000, 5)).unwrap();
        assert_eq!(state.proof.register, 5);
        assert_eq!(state.executed_commands.len(), 1);
        assert_eq!(transition.events.len(), 1);
    }

    #[test]
    fn malformed_and_unsupported_commands_are_rejected_without_mutation() {
        let mut state = SimulationState::new(7).unwrap();
        let before = state.clone();
        let malformed = Command {
            command_id: String::new(),
            command_type: "proof.adjust_register".to_owned(),
            schema_version: SCHEMA_VERSION.to_owned(),
            rules_version: RULES_VERSION.to_owned(),
            effective_time_ms: 0,
            idempotency_key: "key".to_owned(),
            payload: ProofPayload::AdjustRegister { delta: 1 },
        };
        assert!(matches!(
            execute(&mut state, &malformed),
            Err(KernelError::MalformedCommand(_))
        ));
        assert_eq!(state, before);
        let mismatched_payload = Command {
            command_id: "mismatch".to_owned(),
            command_type: "proof.adjust_register".to_owned(),
            schema_version: SCHEMA_VERSION.to_owned(),
            rules_version: RULES_VERSION.to_owned(),
            effective_time_ms: 0,
            idempotency_key: "mismatch-key".to_owned(),
            payload: ProofPayload::SeededPulse {
                options: vec![1, 2],
            },
        };
        assert!(matches!(
            execute(&mut state, &mismatched_payload),
            Err(KernelError::MalformedCommand(_))
        ));
        assert_eq!(state, before);
        let unsupported = Command {
            command_id: "id".to_owned(),
            command_type: "facility.component.advance".to_owned(),
            schema_version: SCHEMA_VERSION.to_owned(),
            rules_version: RULES_VERSION.to_owned(),
            effective_time_ms: 0,
            idempotency_key: "key".to_owned(),
            payload: ProofPayload::AdjustRegister { delta: 1 },
        };
        assert_eq!(
            execute(&mut state, &unsupported),
            Err(KernelError::UnsupportedCommandType(
                "facility.component.advance".to_owned()
            ))
        );
        assert_eq!(state, before);
    }

    #[test]
    fn version_mismatch_is_rejected() {
        let mut state = SimulationState::new(7).unwrap();
        let mut command = Command::adjust_register("id", "key", 0, 1);
        command.rules_version = "rules-other".to_owned();
        assert!(matches!(
            execute(&mut state, &command),
            Err(KernelError::VersionMismatch { .. })
        ));
    }

    #[test]
    fn backwards_time_is_rejected() {
        let mut state = SimulationState::new(7).unwrap();
        advance(&mut state, 1_000).unwrap();
        assert_eq!(
            advance_to(&mut state, 999),
            Err(KernelError::InvalidTimeMovement {
                current_time_ms: 1_000,
                requested_time_ms: 999,
            })
        );
    }

    #[test]
    fn duplicate_command_and_idempotency_key_are_rejected() {
        let mut state = SimulationState::new(7).unwrap();
        let command = Command::adjust_register("id", "key", 0, 1);
        execute(&mut state, &command).unwrap();
        assert_eq!(
            execute(&mut state, &command),
            Err(KernelError::DuplicateCommand {
                command_id: "id".to_owned(),
                idempotency_key: "key".to_owned(),
            })
        );
    }

    #[test]
    fn seeded_rng_is_reproducible_and_serializable() {
        let mut first = SimulationState::new(41).unwrap();
        let mut second = SimulationState::new(41).unwrap();
        let command = Command::seeded_pulse("pulse", "pulse-key", 0, vec![-4, 8, 12]);
        let first_result = execute(&mut first, &command).unwrap();
        let second_result = execute(&mut second, &command).unwrap();
        assert_eq!(first_result.events, second_result.events);
        assert_eq!(first.digest(), second.digest());
        assert_ne!(first.rng.draws, 0);
    }

    #[test]
    fn snapshot_round_trip_preserves_state_and_digest() {
        let mut state = SimulationState::new(41).unwrap();
        execute(&mut state, &Command::adjust_register("id", "key", 0, 3)).unwrap();
        let json = state.to_json().unwrap();
        let restored = SimulationState::from_json(&json).unwrap();
        assert_eq!(restored, state);
        assert_eq!(restored.digest(), state.digest());
    }

    #[test]
    fn facility_graph_round_trip_preserves_evaluation_inputs() {
        let state = SimulationState::with_facilities(41, vec![aggregate_plant_fixture()]).unwrap();
        let json = state.to_json().unwrap();
        let restored = SimulationState::from_json(&json).unwrap();
        assert_eq!(restored, state.canonicalized());
        assert_eq!(restored.digest(), state.digest());
        let evaluation = restored.facilities[0].evaluate().unwrap();
        assert_eq!(evaluation.effective_capacity, 100);
        assert_eq!(evaluation.operational_state, OperationalState::Operational);
    }

    #[test]
    fn failure_commands_use_kernel_replay_and_graph_projection() {
        let mut state =
            SimulationState::with_facilities(41, vec![aggregate_plant_fixture()]).unwrap();
        state.failure = FailureState::with_definitions(aggregate_fault_definitions());
        let activate = Command {
            command_id: "command.failure.activate".to_owned(),
            command_type: "failure.activate".to_owned(),
            schema_version: SCHEMA_VERSION.to_owned(),
            rules_version: RULES_VERSION.to_owned(),
            effective_time_ms: 0,
            idempotency_key: "idempotency.failure.activate".to_owned(),
            payload: ProofPayload::Failure(FailureCommand::ActivateFault {
                fault_instance_id: "fault.instance.feed_motor".to_owned(),
                fault_type_id: "fault.motor_bearing_seizure".to_owned(),
                component_id: "component.feed_conveyor".to_owned(),
                severity_bps: 10_000,
            }),
        };
        execute(&mut state, &activate).unwrap();
        assert_eq!(
            state
                .evaluate_facility("facility.aggregate_plant_fixture")
                .unwrap()
                .effective_capacity,
            0
        );
        let restored = SimulationState::from_json(&state.to_json().unwrap()).unwrap();
        assert_eq!(restored.digest(), state.digest());
    }

    #[test]
    fn failure_state_order_is_canonicalized_in_snapshot_digest() {
        let mut first =
            SimulationState::with_facilities(41, vec![aggregate_plant_fixture()]).unwrap();
        first.failure = FailureState::with_definitions(aggregate_fault_definitions());
        let mut second = first.clone();
        second.failure.fault_definitions.reverse();
        assert_eq!(first.to_json().unwrap(), second.to_json().unwrap());
        assert_eq!(first.digest(), second.digest());
    }

    #[test]
    fn snapshot_rejects_unknown_fields() {
        let input = SimulationState::new(41).unwrap().to_json().unwrap();
        let input = input.replacen(
            "{\"snapshot_type\"",
            "{\"unknown\":true,\"snapshot_type\"",
            1,
        );
        assert!(matches!(
            SimulationState::from_json(&input),
            Err(KernelError::Deserialization(_))
        ));
    }

    #[test]
    fn replay_produces_identical_events_state_and_digest() {
        let commands = command_stream(1_000);
        let (first_state, first_events) = run_stream(SimulationState::new(99).unwrap(), &commands);
        let (second_state, second_events) =
            run_stream(SimulationState::new(99).unwrap(), &commands);
        assert_eq!(first_events, second_events);
        assert_eq!(first_state, second_state);
        assert_eq!(first_state.digest(), second_state.digest());
    }

    #[test]
    fn altered_order_and_seed_change_replay_evidence() {
        let commands = command_stream(1_000);
        let (ordered_state, ordered_events) =
            run_stream(SimulationState::new(99).unwrap(), &commands);
        let reordered = [
            commands[1].clone(),
            commands[0].clone(),
            commands[2].clone(),
        ];
        let (reordered_state, reordered_events) =
            run_stream(SimulationState::new(99).unwrap(), &reordered);
        let (different_seed_state, different_seed_events) =
            run_stream(SimulationState::new(102).unwrap(), &commands);
        assert_ne!(ordered_events, reordered_events);
        assert_ne!(ordered_state.digest(), reordered_state.digest());
        assert_ne!(ordered_events, different_seed_events);
        assert_ne!(ordered_state.digest(), different_seed_state.digest());
    }

    #[test]
    fn generated_monotonic_time_sequence_stays_monotonic() {
        let mut state = SimulationState::new(7).unwrap();
        for elapsed in [0, 1, 7, 100, 1_000, 50_000] {
            let before = state.operational_time_ms;
            advance(&mut state, elapsed).unwrap();
            assert!(state.operational_time_ms >= before);
        }
    }
}
