use crate::graph::{is_valid_semantic_id, DependencyType, Facility, FacilityEvaluation};
use serde::{Deserialize, Serialize};
use std::collections::{BTreeMap, BTreeSet};
use std::fmt::{Display, Formatter};

pub const CONDITION_BPS_WHOLE: u16 = 10_000;

#[derive(Debug, Clone, PartialEq, Eq)]
pub enum FailureError {
    InvalidIdentifier(String),
    InvalidCondition(u16),
    InvalidSeverity(u16),
    InvalidConfidence(u16),
    InvalidDerate(u16),
    DuplicateFaultInstance(String),
    DuplicateEvidence(String),
    DuplicateDiagnosis(String),
    DuplicateCondition(String),
    DuplicateDefinition(String),
    UnknownComponent(String),
    UnknownFaultDefinition(String),
    UnknownFaultInstance(String),
    UnknownEvidence(String),
    IncompatibleFaultComponent {
        fault_type_id: String,
        component_type: String,
    },
    UnsupportedIntervention(String),
    InterventionPrecondition(String),
    AlreadyResolved(String),
    InvalidBypassTarget(String),
    InvalidShutdownRecovery(String),
    InvalidState(String),
    ArithmeticOverflow,
}

impl Display for FailureError {
    fn fmt(&self, formatter: &mut Formatter<'_>) -> std::fmt::Result {
        write!(formatter, "{self:?}")
    }
}

impl std::error::Error for FailureError {}

#[derive(Debug, Clone, Copy, PartialEq, Eq, PartialOrd, Ord, Serialize, Deserialize)]
#[serde(deny_unknown_fields)]
pub enum FaultStatus {
    Latent,
    Active,
    Resolved,
}

#[derive(Debug, Clone, Copy, PartialEq, Eq, PartialOrd, Ord, Serialize, Deserialize)]
#[serde(deny_unknown_fields)]
pub enum Observation {
    Observed,
    Absent,
    Unknown,
}

#[derive(Debug, Clone, Copy, PartialEq, Eq, PartialOrd, Ord, Serialize, Deserialize)]
#[serde(deny_unknown_fields)]
pub enum EvidenceSource {
    VisualInspection,
    Measurement,
    Telemetry,
    ComponentState,
}

#[derive(Debug, Clone, Copy, PartialEq, Eq, PartialOrd, Ord, Serialize, Deserialize)]
#[serde(deny_unknown_fields)]
pub enum DiagnosisStatus {
    Unresolved,
    Confirmed,
    Rejected,
}

#[derive(Debug, Clone, Copy, PartialEq, Eq, PartialOrd, Ord, Serialize, Deserialize)]
#[serde(deny_unknown_fields)]
pub enum InterventionKind {
    Inspect,
    Test,
    Clean,
    Calibrate,
    Lubricate,
    Patch,
    Repair,
    Rebuild,
    Replace,
    Bypass,
    Reroute,
    Derate,
    Shutdown,
    Restore,
    Outsource,
    SalvageDismantle,
}

#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize)]
#[serde(deny_unknown_fields)]
pub struct FaultDefinition {
    pub fault_type_id: String,
    pub component_types: Vec<String>,
    pub default_severity_bps: u16,
    pub capacity_factor_bps: u16,
    pub availability_effect: bool,
    pub symptom_ids: Vec<String>,
    pub permitted_interventions: Vec<InterventionKind>,
}

#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize)]
#[serde(deny_unknown_fields)]
pub struct FaultInstance {
    pub fault_instance_id: String,
    pub fault_type_id: String,
    pub component_id: String,
    pub severity_bps: u16,
    pub status: FaultStatus,
    pub onset_time_ms: u64,
    pub symptom_ids: Vec<String>,
    pub rules_version: String,
}

#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize)]
#[serde(deny_unknown_fields)]
pub struct ComponentConditionState {
    pub component_id: String,
    pub condition_bps: u16,
    pub derate_bps: u16,
    pub deliberate_shutdown: bool,
    pub bypasses: Vec<BypassRule>,
}

impl ComponentConditionState {
    fn default_for(component_id: &str) -> Self {
        Self {
            component_id: component_id.to_owned(),
            condition_bps: CONDITION_BPS_WHOLE,
            derate_bps: CONDITION_BPS_WHOLE,
            deliberate_shutdown: false,
            bypasses: Vec::new(),
        }
    }
}

#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize)]
#[serde(deny_unknown_fields)]
pub struct BypassRule {
    pub bypass_id: String,
    pub upstream_component_id: String,
    pub downstream_component_id: String,
    pub dependency_type: DependencyType,
    pub path_group: Option<String>,
    pub active: bool,
}

#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize)]
#[serde(deny_unknown_fields)]
pub struct BypassSpec {
    pub bypass_id: String,
    pub upstream_component_id: String,
    pub downstream_component_id: String,
    pub dependency_type: DependencyType,
    pub path_group: Option<String>,
}

#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize)]
#[serde(deny_unknown_fields)]
pub struct SymptomObservation {
    pub observation_id: String,
    pub symptom_id: String,
    pub component_id: String,
    pub observation: Observation,
    pub source: EvidenceSource,
    pub effective_time_ms: u64,
}

#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize)]
#[serde(deny_unknown_fields)]
pub struct EvidenceItem {
    pub evidence_id: String,
    pub symptom_id: String,
    pub component_id: String,
    pub observation_id: String,
    pub observation: Observation,
    pub source: EvidenceSource,
    pub effective_time_ms: u64,
    pub factual_value_bps: Option<u16>,
}

#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize)]
#[serde(deny_unknown_fields)]
pub struct Diagnosis {
    pub diagnosis_id: String,
    pub component_id: String,
    pub candidate_fault_type_id: String,
    pub fault_instance_id: Option<String>,
    pub evidence_ids: Vec<String>,
    pub confidence_bps: u16,
    pub status: DiagnosisStatus,
}

#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize)]
#[serde(deny_unknown_fields)]
#[serde(tag = "action", content = "data")]
pub enum FailureCommand {
    ActivateFault {
        fault_instance_id: String,
        fault_type_id: String,
        component_id: String,
        severity_bps: u16,
    },
    Inspect {
        evidence_id: String,
        component_id: String,
        diagnostic_capability_bps: u16,
    },
    Test {
        evidence_id: String,
        component_id: String,
        symptom_id: String,
        diagnostic_capability_bps: u16,
    },
    Diagnose {
        diagnosis_id: String,
        component_id: String,
        candidate_fault_type_id: String,
        fault_instance_id: Option<String>,
        evidence_ids: Vec<String>,
        diagnostic_capability_bps: u16,
    },
    Intervene {
        intervention_id: String,
        component_id: String,
        kind: InterventionKind,
        fault_instance_id: Option<String>,
        derate_bps: Option<u16>,
        bypass: Option<BypassSpec>,
    },
}

#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize)]
#[serde(deny_unknown_fields)]
pub enum FailureEvent {
    FaultActivated {
        fault_instance_id: String,
        fault_type_id: String,
        component_id: String,
    },
    EvidenceObserved {
        evidence_id: String,
        symptom_id: String,
        component_id: String,
        observation: Observation,
    },
    DiagnosisUpdated {
        diagnosis_id: String,
        component_id: String,
        confidence_bps: u16,
        status: DiagnosisStatus,
    },
    InterventionApplied {
        intervention_id: String,
        component_id: String,
        kind: InterventionKind,
    },
    FaultResolved {
        fault_instance_id: String,
        component_id: String,
    },
}

#[derive(Debug, Clone, Default, PartialEq, Eq, Serialize, Deserialize)]
#[serde(deny_unknown_fields)]
pub struct FailureState {
    pub fault_definitions: Vec<FaultDefinition>,
    pub component_conditions: Vec<ComponentConditionState>,
    pub faults: Vec<FaultInstance>,
    pub symptoms: Vec<SymptomObservation>,
    pub evidence: Vec<EvidenceItem>,
    pub diagnoses: Vec<Diagnosis>,
}

struct FailureContext<'a> {
    facilities: &'a [Facility],
    effective_time_ms: u64,
    rules_version: &'a str,
}

impl FailureState {
    pub fn with_definitions(fault_definitions: Vec<FaultDefinition>) -> Self {
        Self {
            fault_definitions,
            ..Self::default()
        }
    }

    pub fn canonicalized(&self) -> Self {
        let mut canonical = self.clone();
        canonical
            .fault_definitions
            .sort_by(|left, right| left.fault_type_id.cmp(&right.fault_type_id));
        canonical
            .component_conditions
            .sort_by(|left, right| left.component_id.cmp(&right.component_id));
        for condition in &mut canonical.component_conditions {
            condition
                .bypasses
                .sort_by(|left, right| left.bypass_id.cmp(&right.bypass_id));
        }
        canonical
            .faults
            .sort_by(|left, right| left.fault_instance_id.cmp(&right.fault_instance_id));
        canonical
            .symptoms
            .sort_by(|left, right| left.observation_id.cmp(&right.observation_id));
        canonical
            .evidence
            .sort_by(|left, right| left.evidence_id.cmp(&right.evidence_id));
        canonical
            .diagnoses
            .sort_by(|left, right| left.diagnosis_id.cmp(&right.diagnosis_id));
        canonical
    }

    pub fn validate(&self) -> Result<(), FailureError> {
        let mut definitions = BTreeSet::new();
        for definition in &self.fault_definitions {
            validate_identifier(&definition.fault_type_id)?;
            if !definitions.insert(&definition.fault_type_id) {
                return Err(FailureError::DuplicateDefinition(
                    definition.fault_type_id.clone(),
                ));
            }
            if definition.component_types.is_empty()
                || definition.symptom_ids.is_empty()
                || definition.permitted_interventions.is_empty()
            {
                return Err(FailureError::InvalidState(
                    "fault definition is incomplete".to_owned(),
                ));
            }
            validate_severity(definition.default_severity_bps)?;
            validate_bps(definition.capacity_factor_bps)?;
            for symptom_id in &definition.symptom_ids {
                validate_identifier(symptom_id)?;
            }
        }
        validate_unique_conditions(&self.component_conditions)?;
        let definitions = self
            .fault_definitions
            .iter()
            .map(|definition| (definition.fault_type_id.clone(), definition))
            .collect::<BTreeMap<_, _>>();
        let mut faults = BTreeSet::new();
        for fault in &self.faults {
            validate_identifier(&fault.fault_instance_id)?;
            validate_identifier(&fault.component_id)?;
            validate_severity(fault.severity_bps)?;
            if !faults.insert(&fault.fault_instance_id) {
                return Err(FailureError::DuplicateFaultInstance(
                    fault.fault_instance_id.clone(),
                ));
            }
            let definition = definitions
                .get(&fault.fault_type_id)
                .ok_or_else(|| FailureError::UnknownFaultDefinition(fault.fault_type_id.clone()))?;
            if fault.symptom_ids != definition.symptom_ids {
                return Err(FailureError::InvalidState(
                    "fault symptoms do not match definition".to_owned(),
                ));
            }
        }
        let mut evidence = BTreeSet::new();
        for item in &self.evidence {
            validate_identifier(&item.evidence_id)?;
            validate_identifier(&item.symptom_id)?;
            validate_identifier(&item.component_id)?;
            if !evidence.insert(&item.evidence_id) {
                return Err(FailureError::DuplicateEvidence(item.evidence_id.clone()));
            }
            if let Some(value) = item.factual_value_bps {
                validate_bps(value)?;
            }
        }
        let mut diagnoses = BTreeSet::new();
        for diagnosis in &self.diagnoses {
            validate_identifier(&diagnosis.diagnosis_id)?;
            validate_identifier(&diagnosis.component_id)?;
            validate_identifier(&diagnosis.candidate_fault_type_id)?;
            validate_confidence(diagnosis.confidence_bps)?;
            if let Some(fault_instance_id) = &diagnosis.fault_instance_id {
                let fault = self
                    .faults
                    .iter()
                    .find(|fault| &fault.fault_instance_id == fault_instance_id)
                    .ok_or_else(|| {
                        FailureError::InvalidState("diagnosis references unknown fault".to_owned())
                    })?;
                if fault.component_id != diagnosis.component_id
                    || fault.fault_type_id != diagnosis.candidate_fault_type_id
                    || fault.status == FaultStatus::Resolved
                {
                    return Err(FailureError::InvalidState(
                        "persisted diagnosis fault reference is inconsistent".to_owned(),
                    ));
                }
            }
            if !diagnoses.insert(&diagnosis.diagnosis_id) {
                return Err(FailureError::DuplicateDiagnosis(
                    diagnosis.diagnosis_id.clone(),
                ));
            }
        }
        Ok(())
    }

    pub fn validate_against_facilities(&self, facilities: &[Facility]) -> Result<(), FailureError> {
        self.validate()?;
        let components = facilities
            .iter()
            .flat_map(|facility| facility.components.iter())
            .map(|component| {
                (
                    component.component_id.as_str(),
                    component.component_type.as_str(),
                )
            })
            .collect::<BTreeMap<_, _>>();
        for condition in &self.component_conditions {
            ensure_component(&components, &condition.component_id)?;
            for bypass in &condition.bypasses {
                ensure_bypass_parts(
                    facilities,
                    &bypass.upstream_component_id,
                    &bypass.downstream_component_id,
                    bypass.dependency_type,
                    bypass.path_group.as_ref(),
                    &bypass.bypass_id,
                )?;
            }
        }
        for fault in &self.faults {
            let component_type = ensure_component(&components, &fault.component_id)?;
            let definition = self
                .fault_definitions
                .iter()
                .find(|definition| definition.fault_type_id == fault.fault_type_id)
                .ok_or_else(|| FailureError::UnknownFaultDefinition(fault.fault_type_id.clone()))?;
            if !definition
                .component_types
                .iter()
                .any(|kind| kind == component_type)
            {
                return Err(FailureError::IncompatibleFaultComponent {
                    fault_type_id: fault.fault_type_id.clone(),
                    component_type: component_type.to_owned(),
                });
            }
        }
        for evidence in &self.evidence {
            ensure_component(&components, &evidence.component_id)?;
        }
        for symptom in &self.symptoms {
            ensure_component(&components, &symptom.component_id)?;
        }
        for diagnosis in &self.diagnoses {
            ensure_component(&components, &diagnosis.component_id)?;
        }
        Ok(())
    }

    pub fn evaluate_facility(
        &self,
        facility: &Facility,
    ) -> Result<FacilityEvaluation, FailureError> {
        let projected = self.project_facility(facility)?;
        projected
            .evaluate()
            .map_err(|error| FailureError::InvalidState(format!("graph evaluation: {error}")))
    }

    pub fn project_facility(&self, facility: &Facility) -> Result<Facility, FailureError> {
        self.validate_against_facilities(std::slice::from_ref(facility))?;
        let mut projected = facility.clone();
        for condition in &self.component_conditions {
            let component = projected
                .components
                .iter_mut()
                .find(|component| component.component_id == condition.component_id)
                .ok_or_else(|| FailureError::UnknownComponent(condition.component_id.clone()))?;
            component.capacity_factor_bps = component
                .capacity_factor_bps
                .min(condition.condition_bps)
                .min(condition.derate_bps);
            if condition.deliberate_shutdown {
                component.available = false;
            }
        }
        for fault in self
            .faults
            .iter()
            .filter(|fault| fault.status == FaultStatus::Active)
        {
            let definition = self
                .fault_definitions
                .iter()
                .find(|definition| definition.fault_type_id == fault.fault_type_id)
                .ok_or_else(|| FailureError::UnknownFaultDefinition(fault.fault_type_id.clone()))?;
            let component = projected
                .components
                .iter_mut()
                .find(|component| component.component_id == fault.component_id)
                .ok_or_else(|| FailureError::UnknownComponent(fault.component_id.clone()))?;
            component.capacity_factor_bps = component
                .capacity_factor_bps
                .min(definition.capacity_factor_bps);
            if definition.availability_effect {
                component.available = false;
            }
        }
        for condition in &self.component_conditions {
            for bypass in condition.bypasses.iter().filter(|bypass| bypass.active) {
                projected.dependencies.retain(|edge| {
                    !(edge.upstream_component_id == bypass.upstream_component_id
                        && edge.downstream_component_id == bypass.downstream_component_id
                        && edge.dependency_type == bypass.dependency_type
                        && edge.path_group == bypass.path_group)
                });
            }
        }
        Ok(projected)
    }

    pub fn apply(
        &mut self,
        facilities: &[Facility],
        command: &FailureCommand,
        effective_time_ms: u64,
        rules_version: &str,
    ) -> Result<FailureEvent, FailureError> {
        self.validate_against_facilities(facilities)?;
        let context = FailureContext {
            facilities,
            effective_time_ms,
            rules_version,
        };
        match command {
            FailureCommand::ActivateFault {
                fault_instance_id,
                fault_type_id,
                component_id,
                severity_bps,
            } => self.activate_fault(
                fault_instance_id,
                fault_type_id,
                component_id,
                *severity_bps,
                &context,
            ),
            FailureCommand::Inspect {
                evidence_id,
                component_id,
                diagnostic_capability_bps,
            } => self.observe(
                evidence_id,
                component_id,
                None,
                EvidenceSource::VisualInspection,
                *diagnostic_capability_bps,
                &context,
            ),
            FailureCommand::Test {
                evidence_id,
                component_id,
                symptom_id,
                diagnostic_capability_bps,
            } => self.observe(
                evidence_id,
                component_id,
                Some(symptom_id),
                EvidenceSource::Measurement,
                *diagnostic_capability_bps,
                &context,
            ),
            FailureCommand::Diagnose {
                diagnosis_id,
                component_id,
                candidate_fault_type_id,
                fault_instance_id,
                evidence_ids,
                diagnostic_capability_bps,
            } => self.diagnose(
                diagnosis_id,
                component_id,
                candidate_fault_type_id,
                fault_instance_id.as_ref(),
                evidence_ids,
                *diagnostic_capability_bps,
                &context,
            ),
            FailureCommand::Intervene {
                intervention_id,
                component_id,
                kind,
                fault_instance_id,
                derate_bps,
                bypass,
            } => self.intervene(
                intervention_id,
                component_id,
                *kind,
                fault_instance_id.as_ref(),
                *derate_bps,
                bypass.as_ref(),
                &context,
            ),
        }
    }

    fn activate_fault(
        &mut self,
        fault_instance_id: &str,
        fault_type_id: &str,
        component_id: &str,
        severity_bps: u16,
        context: &FailureContext<'_>,
    ) -> Result<FailureEvent, FailureError> {
        validate_identifier(fault_instance_id)?;
        validate_identifier(fault_type_id)?;
        validate_severity(severity_bps)?;
        let component_type = find_component_type(context.facilities, component_id)?;
        if self
            .faults
            .iter()
            .any(|fault| fault.fault_instance_id == fault_instance_id)
        {
            return Err(FailureError::DuplicateFaultInstance(
                fault_instance_id.to_owned(),
            ));
        }
        let definition = self
            .fault_definitions
            .iter()
            .find(|definition| definition.fault_type_id == fault_type_id)
            .ok_or_else(|| FailureError::UnknownFaultDefinition(fault_type_id.to_owned()))?;
        if !definition
            .component_types
            .iter()
            .any(|kind| kind == component_type)
        {
            return Err(FailureError::IncompatibleFaultComponent {
                fault_type_id: fault_type_id.to_owned(),
                component_type: component_type.to_owned(),
            });
        }
        let instance = FaultInstance {
            fault_instance_id: fault_instance_id.to_owned(),
            fault_type_id: fault_type_id.to_owned(),
            component_id: component_id.to_owned(),
            severity_bps,
            status: FaultStatus::Active,
            onset_time_ms: context.effective_time_ms,
            symptom_ids: definition.symptom_ids.clone(),
            rules_version: context.rules_version.to_owned(),
        };
        self.faults.push(instance);
        Ok(FailureEvent::FaultActivated {
            fault_instance_id: fault_instance_id.to_owned(),
            fault_type_id: fault_type_id.to_owned(),
            component_id: component_id.to_owned(),
        })
    }

    fn observe(
        &mut self,
        evidence_id: &str,
        component_id: &str,
        requested_symptom_id: Option<&str>,
        source: EvidenceSource,
        diagnostic_capability_bps: u16,
        context: &FailureContext<'_>,
    ) -> Result<FailureEvent, FailureError> {
        validate_identifier(evidence_id)?;
        validate_confidence(diagnostic_capability_bps)?;
        find_component_type(context.facilities, component_id)?;
        if self
            .evidence
            .iter()
            .any(|evidence| evidence.evidence_id == evidence_id)
        {
            return Err(FailureError::DuplicateEvidence(evidence_id.to_owned()));
        }
        let symptom_id = match requested_symptom_id {
            Some(symptom_id) => {
                validate_identifier(symptom_id)?;
                symptom_id.to_owned()
            }
            None => self
                .faults
                .iter()
                .filter(|fault| {
                    fault.component_id == component_id && fault.status == FaultStatus::Active
                })
                .flat_map(|fault| fault.symptom_ids.iter())
                .min()
                .cloned()
                .unwrap_or_else(|| "symptom.component.operability".to_owned()),
        };
        let observation = if diagnostic_capability_bps < 5_000 {
            Observation::Unknown
        } else if self.faults.iter().any(|fault| {
            fault.component_id == component_id
                && fault.status == FaultStatus::Active
                && fault
                    .symptom_ids
                    .iter()
                    .any(|symptom| symptom == &symptom_id)
        }) {
            Observation::Observed
        } else {
            Observation::Absent
        };
        let observation_id = format!("observation.{evidence_id}");
        self.symptoms.push(SymptomObservation {
            observation_id: observation_id.clone(),
            symptom_id: symptom_id.clone(),
            component_id: component_id.to_owned(),
            observation,
            source,
            effective_time_ms: context.effective_time_ms,
        });
        self.evidence.push(EvidenceItem {
            evidence_id: evidence_id.to_owned(),
            symptom_id: symptom_id.clone(),
            component_id: component_id.to_owned(),
            observation_id,
            observation,
            source,
            effective_time_ms: context.effective_time_ms,
            factual_value_bps: (observation != Observation::Unknown)
                .then(|| {
                    self.condition_for(component_id)
                        .ok()
                        .map(|state| state.condition_bps)
                })
                .flatten(),
        });
        Ok(FailureEvent::EvidenceObserved {
            evidence_id: evidence_id.to_owned(),
            symptom_id,
            component_id: component_id.to_owned(),
            observation,
        })
    }

    #[allow(clippy::too_many_arguments)]
    fn diagnose(
        &mut self,
        diagnosis_id: &str,
        component_id: &str,
        candidate_fault_type_id: &str,
        fault_instance_id: Option<&String>,
        evidence_ids: &[String],
        diagnostic_capability_bps: u16,
        context: &FailureContext<'_>,
    ) -> Result<FailureEvent, FailureError> {
        validate_identifier(diagnosis_id)?;
        validate_identifier(candidate_fault_type_id)?;
        validate_confidence(diagnostic_capability_bps)?;
        find_component_type(context.facilities, component_id)?;
        if self
            .diagnoses
            .iter()
            .any(|diagnosis| diagnosis.diagnosis_id == diagnosis_id)
        {
            return Err(FailureError::DuplicateDiagnosis(diagnosis_id.to_owned()));
        }
        let definition = self
            .fault_definitions
            .iter()
            .find(|definition| definition.fault_type_id == candidate_fault_type_id)
            .ok_or_else(|| {
                FailureError::UnknownFaultDefinition(candidate_fault_type_id.to_owned())
            })?;
        let mut observations = Vec::new();
        for evidence_id in evidence_ids {
            let evidence = self
                .evidence
                .iter()
                .find(|evidence| &evidence.evidence_id == evidence_id)
                .ok_or_else(|| FailureError::UnknownEvidence(evidence_id.clone()))?;
            if evidence.component_id != component_id {
                return Err(FailureError::InterventionPrecondition(
                    "evidence targets another component".to_owned(),
                ));
            }
            observations.push(evidence);
        }
        let observed = observations
            .iter()
            .filter(|evidence| {
                evidence.observation == Observation::Observed
                    && definition.symptom_ids.contains(&evidence.symptom_id)
            })
            .count() as u16;
        let unknown = observations
            .iter()
            .filter(|evidence| evidence.observation == Observation::Unknown)
            .count() as u16;
        let absent = observations
            .iter()
            .filter(|evidence| {
                evidence.observation == Observation::Absent
                    && definition.symptom_ids.contains(&evidence.symptom_id)
            })
            .count() as u16;
        let raw_confidence = if observed > 0 && absent == 0 {
            7_000u16
                .saturating_add(observed.saturating_mul(1_000))
                .min(10_000)
        } else if unknown > 0 {
            3_000
        } else {
            1_000
        };
        let confidence_bps = (u32::from(raw_confidence) * u32::from(diagnostic_capability_bps)
            / u32::from(CONDITION_BPS_WHOLE)) as u16;
        let status = if confidence_bps >= 8_000 {
            DiagnosisStatus::Confirmed
        } else if absent > 0 && observed == 0 {
            DiagnosisStatus::Rejected
        } else {
            DiagnosisStatus::Unresolved
        };
        if let Some(fault_instance_id) = fault_instance_id {
            let fault = self
                .faults
                .iter()
                .find(|fault| &fault.fault_instance_id == fault_instance_id)
                .ok_or_else(|| FailureError::UnknownFaultInstance(fault_instance_id.clone()))?;
            if fault.component_id != component_id || fault.fault_type_id != candidate_fault_type_id
            {
                return Err(FailureError::InvalidState(
                    "diagnosis fault instance does not match component/type".to_owned(),
                ));
            }
            if fault.status == FaultStatus::Resolved {
                return Err(FailureError::InterventionPrecondition(
                    "diagnosis cannot bind a resolved fault".to_owned(),
                ));
            }
        }
        self.diagnoses.push(Diagnosis {
            diagnosis_id: diagnosis_id.to_owned(),
            component_id: component_id.to_owned(),
            candidate_fault_type_id: candidate_fault_type_id.to_owned(),
            fault_instance_id: fault_instance_id.cloned(),
            evidence_ids: evidence_ids.to_vec(),
            confidence_bps,
            status,
        });
        Ok(FailureEvent::DiagnosisUpdated {
            diagnosis_id: diagnosis_id.to_owned(),
            component_id: component_id.to_owned(),
            confidence_bps,
            status,
        })
    }

    #[allow(clippy::too_many_arguments)]
    fn intervene(
        &mut self,
        intervention_id: &str,
        component_id: &str,
        kind: InterventionKind,
        fault_instance_id: Option<&String>,
        derate_bps: Option<u16>,
        bypass: Option<&BypassSpec>,
        context: &FailureContext<'_>,
    ) -> Result<FailureEvent, FailureError> {
        validate_identifier(intervention_id)?;
        find_component_type(context.facilities, component_id)?;
        if kind == InterventionKind::Repair && fault_instance_id.is_none() {
            return Err(FailureError::InterventionPrecondition(
                "repair requires a fault instance".to_owned(),
            ));
        }
        if let Some(fault_instance_id) = fault_instance_id {
            let fault = self
                .faults
                .iter()
                .find(|fault| {
                    &fault.fault_instance_id == fault_instance_id
                        && fault.component_id == component_id
                })
                .ok_or_else(|| FailureError::UnknownFaultInstance(fault_instance_id.clone()))?;
            if fault.status == FaultStatus::Resolved {
                return Err(FailureError::AlreadyResolved(
                    fault.fault_instance_id.clone(),
                ));
            }
            if fault.status != FaultStatus::Active {
                return Err(FailureError::InterventionPrecondition(
                    "intervention requires an active fault".to_owned(),
                ));
            }
            let definition = self
                .fault_definitions
                .iter()
                .find(|definition| definition.fault_type_id == fault.fault_type_id)
                .ok_or_else(|| FailureError::UnknownFaultDefinition(fault.fault_type_id.clone()))?;
            if !definition.permitted_interventions.contains(&kind) {
                return Err(FailureError::UnsupportedIntervention(format!("{kind:?}")));
            }
        }
        match kind {
            InterventionKind::Repair
            | InterventionKind::Patch
            | InterventionKind::Lubricate
            | InterventionKind::Clean => {
                let condition = self.condition_for_mut(component_id);
                condition.condition_bps = condition
                    .condition_bps
                    .saturating_add(2_000)
                    .min(CONDITION_BPS_WHOLE);
                if let Some(fault_instance_id) = fault_instance_id {
                    resolve_fault(&mut self.faults, fault_instance_id, component_id)?;
                }
            }
            InterventionKind::Rebuild => {
                let condition = self.condition_for_mut(component_id);
                condition.condition_bps = CONDITION_BPS_WHOLE;
                if let Some(fault_instance_id) = fault_instance_id {
                    resolve_fault(&mut self.faults, fault_instance_id, component_id)?;
                }
            }
            InterventionKind::Replace => {
                let condition = self.condition_for_mut(component_id);
                condition.condition_bps = CONDITION_BPS_WHOLE;
                condition.derate_bps = CONDITION_BPS_WHOLE;
                condition.deliberate_shutdown = false;
                for fault in self.faults.iter_mut().filter(|fault| {
                    fault.component_id == component_id && fault.status == FaultStatus::Active
                }) {
                    fault.status = FaultStatus::Resolved;
                }
            }
            InterventionKind::Bypass | InterventionKind::Reroute => {
                let bypass = bypass.ok_or_else(|| {
                    FailureError::InvalidBypassTarget("bypass specification is required".to_owned())
                })?;
                ensure_bypass_parts(
                    context.facilities,
                    &bypass.upstream_component_id,
                    &bypass.downstream_component_id,
                    bypass.dependency_type,
                    bypass.path_group.as_ref(),
                    &bypass.bypass_id,
                )?;
                let condition = self.condition_for_mut(component_id);
                if condition
                    .bypasses
                    .iter()
                    .any(|existing| existing.bypass_id == bypass.bypass_id)
                {
                    return Err(FailureError::InterventionPrecondition(
                        "bypass already active".to_owned(),
                    ));
                }
                condition.bypasses.push(BypassRule {
                    bypass_id: bypass.bypass_id.clone(),
                    upstream_component_id: bypass.upstream_component_id.clone(),
                    downstream_component_id: bypass.downstream_component_id.clone(),
                    dependency_type: bypass.dependency_type,
                    path_group: bypass.path_group.clone(),
                    active: true,
                });
            }
            InterventionKind::Derate => {
                let derate_bps = derate_bps.ok_or(FailureError::InvalidDerate(0))?;
                validate_derate(derate_bps)?;
                let condition = self.condition_for_mut(component_id);
                condition.derate_bps = derate_bps;
            }
            InterventionKind::Shutdown => {
                self.condition_for_mut(component_id).deliberate_shutdown = true;
            }
            InterventionKind::Restore => {
                let condition = self.condition_for_mut(component_id);
                if !condition.deliberate_shutdown && condition.derate_bps == CONDITION_BPS_WHOLE {
                    return Err(FailureError::InvalidShutdownRecovery(
                        component_id.to_owned(),
                    ));
                }
                condition.deliberate_shutdown = false;
                condition.derate_bps = CONDITION_BPS_WHOLE;
            }
            InterventionKind::Inspect
            | InterventionKind::Test
            | InterventionKind::Calibrate
            | InterventionKind::Outsource
            | InterventionKind::SalvageDismantle => {}
        }
        Ok(FailureEvent::InterventionApplied {
            intervention_id: intervention_id.to_owned(),
            component_id: component_id.to_owned(),
            kind,
        })
    }

    fn condition_for(&self, component_id: &str) -> Result<ComponentConditionState, FailureError> {
        Ok(self
            .component_conditions
            .iter()
            .find(|condition| condition.component_id == component_id)
            .cloned()
            .unwrap_or_else(|| ComponentConditionState::default_for(component_id)))
    }

    fn condition_for_mut(&mut self, component_id: &str) -> &mut ComponentConditionState {
        if let Some(index) = self
            .component_conditions
            .iter()
            .position(|condition| condition.component_id == component_id)
        {
            return &mut self.component_conditions[index];
        }
        self.component_conditions
            .push(ComponentConditionState::default_for(component_id));
        self.component_conditions
            .last_mut()
            .expect("condition inserted")
    }
}

fn validate_identifier(identifier: &str) -> Result<(), FailureError> {
    if is_valid_semantic_id(identifier) {
        Ok(())
    } else {
        Err(FailureError::InvalidIdentifier(identifier.to_owned()))
    }
}

fn validate_bps(value: u16) -> Result<(), FailureError> {
    if u32::from(value) <= u32::from(CONDITION_BPS_WHOLE) {
        Ok(())
    } else {
        Err(FailureError::InvalidCondition(value))
    }
}

fn validate_severity(value: u16) -> Result<(), FailureError> {
    if u32::from(value) <= u32::from(CONDITION_BPS_WHOLE) {
        Ok(())
    } else {
        Err(FailureError::InvalidSeverity(value))
    }
}

fn validate_confidence(value: u16) -> Result<(), FailureError> {
    if u32::from(value) <= u32::from(CONDITION_BPS_WHOLE) {
        Ok(())
    } else {
        Err(FailureError::InvalidConfidence(value))
    }
}

fn validate_derate(value: u16) -> Result<(), FailureError> {
    if u32::from(value) <= u32::from(CONDITION_BPS_WHOLE) {
        Ok(())
    } else {
        Err(FailureError::InvalidDerate(value))
    }
}

fn validate_unique_conditions(conditions: &[ComponentConditionState]) -> Result<(), FailureError> {
    let mut ids = BTreeSet::new();
    for condition in conditions {
        validate_identifier(&condition.component_id)?;
        validate_bps(condition.condition_bps)?;
        validate_bps(condition.derate_bps)?;
        if !ids.insert(&condition.component_id) {
            return Err(FailureError::DuplicateCondition(
                condition.component_id.clone(),
            ));
        }
        let mut bypass_ids = BTreeSet::new();
        for bypass in &condition.bypasses {
            validate_identifier(&bypass.bypass_id)?;
            validate_identifier(&bypass.upstream_component_id)?;
            validate_identifier(&bypass.downstream_component_id)?;
            if !bypass_ids.insert(&bypass.bypass_id) {
                return Err(FailureError::InvalidState("duplicate bypass ID".to_owned()));
            }
        }
    }
    Ok(())
}

fn find_component_type<'a>(
    facilities: &'a [Facility],
    component_id: &str,
) -> Result<&'a str, FailureError> {
    facilities
        .iter()
        .flat_map(|facility| facility.components.iter())
        .find(|component| component.component_id == component_id)
        .map(|component| component.component_type.as_str())
        .ok_or_else(|| FailureError::UnknownComponent(component_id.to_owned()))
}

fn ensure_component<'a>(
    components: &BTreeMap<&'a str, &'a str>,
    component_id: &str,
) -> Result<&'a str, FailureError> {
    components
        .get(component_id)
        .copied()
        .ok_or_else(|| FailureError::UnknownComponent(component_id.to_owned()))
}

fn ensure_bypass_parts(
    facilities: &[Facility],
    upstream_component_id: &str,
    downstream_component_id: &str,
    dependency_type: DependencyType,
    path_group: Option<&String>,
    bypass_id: &str,
) -> Result<(), FailureError> {
    validate_identifier(bypass_id)?;
    validate_identifier(upstream_component_id)?;
    validate_identifier(downstream_component_id)?;
    let exists = facilities
        .iter()
        .flat_map(|facility| facility.dependencies.iter())
        .any(|edge| {
            edge.upstream_component_id == upstream_component_id
                && edge.downstream_component_id == downstream_component_id
                && edge.dependency_type == dependency_type
                && edge.path_group.as_ref() == path_group
        });
    if exists {
        Ok(())
    } else {
        Err(FailureError::InvalidBypassTarget(bypass_id.to_owned()))
    }
}

fn resolve_fault(
    faults: &mut [FaultInstance],
    fault_instance_id: &str,
    component_id: &str,
) -> Result<(), FailureError> {
    let fault = faults
        .iter_mut()
        .find(|fault| {
            fault.fault_instance_id == fault_instance_id && fault.component_id == component_id
        })
        .ok_or_else(|| FailureError::UnknownFaultInstance(fault_instance_id.to_owned()))?;
    if fault.status == FaultStatus::Resolved {
        return Err(FailureError::AlreadyResolved(fault_instance_id.to_owned()));
    }
    fault.status = FaultStatus::Resolved;
    Ok(())
}

pub fn aggregate_fault_definitions() -> Vec<FaultDefinition> {
    vec![
        FaultDefinition {
            fault_type_id: "fault.motor_bearing_seizure".to_owned(),
            component_types: vec!["conveyor".to_owned()],
            default_severity_bps: 10_000,
            capacity_factor_bps: 0,
            availability_effect: true,
            symptom_ids: vec![
                "symptom.abnormal_vibration".to_owned(),
                "symptom.no_rotation".to_owned(),
            ],
            permitted_interventions: vec![
                InterventionKind::Repair,
                InterventionKind::Replace,
                InterventionKind::Bypass,
            ],
        },
        FaultDefinition {
            fault_type_id: "fault.worn_belt".to_owned(),
            component_types: vec!["conveyor".to_owned()],
            default_severity_bps: 4_000,
            capacity_factor_bps: 7_000,
            availability_effect: false,
            symptom_ids: vec!["symptom.abnormal_vibration".to_owned()],
            permitted_interventions: vec![
                InterventionKind::Repair,
                InterventionKind::Replace,
                InterventionKind::Derate,
            ],
        },
        FaultDefinition {
            fault_type_id: "fault.screen_blockage".to_owned(),
            component_types: vec!["screen".to_owned()],
            default_severity_bps: 5_000,
            capacity_factor_bps: 3_000,
            availability_effect: false,
            symptom_ids: vec![
                "symptom.blockage_indication".to_owned(),
                "symptom.low_flow".to_owned(),
            ],
            permitted_interventions: vec![
                InterventionKind::Clean,
                InterventionKind::Repair,
                InterventionKind::Replace,
            ],
        },
    ]
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::graph::{aggregate_plant_fixture, OperationalState};
    use crate::RULES_VERSION;

    fn facilities() -> Vec<Facility> {
        vec![aggregate_plant_fixture()]
    }

    fn state() -> FailureState {
        FailureState::with_definitions(aggregate_fault_definitions())
    }

    #[test]
    fn seized_motor_has_causal_evidence_and_repair_recovery() {
        let facilities = facilities();
        let mut failure = state();
        failure
            .apply(
                &facilities,
                &FailureCommand::ActivateFault {
                    fault_instance_id: "fault.instance.feed_motor".to_owned(),
                    fault_type_id: "fault.motor_bearing_seizure".to_owned(),
                    component_id: "component.feed_conveyor".to_owned(),
                    severity_bps: 10_000,
                },
                0,
                RULES_VERSION,
            )
            .unwrap();
        let outage = failure.evaluate_facility(&facilities[0]).unwrap();
        assert_eq!(outage.effective_capacity, 0);
        failure
            .apply(
                &facilities,
                &FailureCommand::Inspect {
                    evidence_id: "evidence.motor_inspection".to_owned(),
                    component_id: "component.feed_conveyor".to_owned(),
                    diagnostic_capability_bps: 10_000,
                },
                0,
                RULES_VERSION,
            )
            .unwrap();
        assert_eq!(failure.evidence[0].observation, Observation::Observed);
        let before_invalid_repair = failure.clone();
        assert!(matches!(
            failure.apply(
                &facilities,
                &FailureCommand::Intervene {
                    intervention_id: "intervention.invalid_repair".to_owned(),
                    component_id: "component.feed_conveyor".to_owned(),
                    kind: InterventionKind::Repair,
                    fault_instance_id: None,
                    derate_bps: None,
                    bypass: None,
                },
                0,
                RULES_VERSION,
            ),
            Err(FailureError::InterventionPrecondition(_))
        ));
        assert_eq!(failure, before_invalid_repair);
        failure
            .apply(
                &facilities,
                &FailureCommand::Intervene {
                    intervention_id: "intervention.motor_repair".to_owned(),
                    component_id: "component.feed_conveyor".to_owned(),
                    kind: InterventionKind::Repair,
                    fault_instance_id: Some("fault.instance.feed_motor".to_owned()),
                    derate_bps: None,
                    bypass: None,
                },
                0,
                RULES_VERSION,
            )
            .unwrap();
        assert_eq!(failure.faults[0].status, FaultStatus::Resolved);
        assert_eq!(
            failure
                .evaluate_facility(&facilities[0])
                .unwrap()
                .effective_capacity,
            100
        );
    }

    #[test]
    fn worn_belt_and_blocked_screen_preserve_partial_operation() {
        let facilities = facilities();
        let mut failure = state();
        failure
            .apply(
                &facilities,
                &FailureCommand::ActivateFault {
                    fault_instance_id: "fault.instance.worn_belt".to_owned(),
                    fault_type_id: "fault.worn_belt".to_owned(),
                    component_id: "component.feed_conveyor".to_owned(),
                    severity_bps: 4_000,
                },
                0,
                RULES_VERSION,
            )
            .unwrap();
        let worn = failure.evaluate_facility(&facilities[0]).unwrap();
        assert_eq!(worn.effective_capacity, 70);
        assert_eq!(worn.operational_state, OperationalState::Constrained);
        let mut blocked = state();
        blocked
            .apply(
                &facilities,
                &FailureCommand::ActivateFault {
                    fault_instance_id: "fault.instance.screen_block".to_owned(),
                    fault_type_id: "fault.screen_blockage".to_owned(),
                    component_id: "component.screen".to_owned(),
                    severity_bps: 5_000,
                },
                0,
                RULES_VERSION,
            )
            .unwrap();
        let partial = blocked.evaluate_facility(&facilities[0]).unwrap();
        assert_eq!(partial.effective_capacity, 30);
        assert_eq!(partial.operational_state, OperationalState::Constrained);
    }

    #[test]
    fn unknown_evidence_does_not_fabricate_diagnosis() {
        let facilities = facilities();
        let mut failure = state();
        failure
            .apply(
                &facilities,
                &FailureCommand::Test {
                    evidence_id: "evidence.low_information".to_owned(),
                    component_id: "component.screen".to_owned(),
                    symptom_id: "symptom.blockage_indication".to_owned(),
                    diagnostic_capability_bps: 1_000,
                },
                0,
                RULES_VERSION,
            )
            .unwrap();
        assert_eq!(failure.evidence[0].observation, Observation::Unknown);
        assert_eq!(failure.evidence[0].factual_value_bps, None);
        failure
            .apply(
                &facilities,
                &FailureCommand::Diagnose {
                    diagnosis_id: "diagnosis.screen".to_owned(),
                    component_id: "component.screen".to_owned(),
                    candidate_fault_type_id: "fault.screen_blockage".to_owned(),
                    fault_instance_id: None,
                    evidence_ids: vec!["evidence.low_information".to_owned()],
                    diagnostic_capability_bps: 1_000,
                },
                0,
                RULES_VERSION,
            )
            .unwrap();
        assert_eq!(failure.diagnoses[0].status, DiagnosisStatus::Unresolved);
        assert!(failure.diagnoses[0].confidence_bps < 8_000);
    }

    #[test]
    fn repair_rejects_wrong_component_and_already_resolved_fault() {
        let facilities = facilities();
        let mut failure = state();
        failure
            .apply(
                &facilities,
                &FailureCommand::ActivateFault {
                    fault_instance_id: "fault.instance.repair_target".to_owned(),
                    fault_type_id: "fault.motor_bearing_seizure".to_owned(),
                    component_id: "component.feed_conveyor".to_owned(),
                    severity_bps: 10_000,
                },
                0,
                RULES_VERSION,
            )
            .unwrap();
        assert!(matches!(
            failure.apply(
                &facilities,
                &FailureCommand::Intervene {
                    intervention_id: "intervention.wrong_component".to_owned(),
                    component_id: "component.crusher".to_owned(),
                    kind: InterventionKind::Repair,
                    fault_instance_id: Some("fault.instance.repair_target".to_owned()),
                    derate_bps: None,
                    bypass: None,
                },
                0,
                RULES_VERSION,
            ),
            Err(FailureError::UnknownFaultInstance(_))
        ));
        failure
            .apply(
                &facilities,
                &FailureCommand::Intervene {
                    intervention_id: "intervention.valid_repair".to_owned(),
                    component_id: "component.feed_conveyor".to_owned(),
                    kind: InterventionKind::Repair,
                    fault_instance_id: Some("fault.instance.repair_target".to_owned()),
                    derate_bps: None,
                    bypass: None,
                },
                0,
                RULES_VERSION,
            )
            .unwrap();
        assert!(matches!(
            failure.apply(
                &facilities,
                &FailureCommand::Intervene {
                    intervention_id: "intervention.repeated_repair".to_owned(),
                    component_id: "component.feed_conveyor".to_owned(),
                    kind: InterventionKind::Repair,
                    fault_instance_id: Some("fault.instance.repair_target".to_owned()),
                    derate_bps: None,
                    bypass: None,
                },
                0,
                RULES_VERSION,
            ),
            Err(FailureError::AlreadyResolved(_))
        ));
    }

    #[test]
    fn diagnosis_rejects_mismatched_bound_fault_instance() {
        let facilities = facilities();
        let mut failure = state();
        failure
            .apply(
                &facilities,
                &FailureCommand::ActivateFault {
                    fault_instance_id: "fault.instance.motor_for_diagnosis".to_owned(),
                    fault_type_id: "fault.motor_bearing_seizure".to_owned(),
                    component_id: "component.feed_conveyor".to_owned(),
                    severity_bps: 10_000,
                },
                0,
                RULES_VERSION,
            )
            .unwrap();
        failure
            .apply(
                &facilities,
                &FailureCommand::Inspect {
                    evidence_id: "evidence.motor_for_diagnosis".to_owned(),
                    component_id: "component.feed_conveyor".to_owned(),
                    diagnostic_capability_bps: 10_000,
                },
                0,
                RULES_VERSION,
            )
            .unwrap();
        assert!(matches!(
            failure.apply(
                &facilities,
                &FailureCommand::Diagnose {
                    diagnosis_id: "diagnosis.mismatched_fault".to_owned(),
                    component_id: "component.feed_conveyor".to_owned(),
                    candidate_fault_type_id: "fault.worn_belt".to_owned(),
                    fault_instance_id: Some("fault.instance.motor_for_diagnosis".to_owned()),
                    evidence_ids: vec!["evidence.motor_for_diagnosis".to_owned()],
                    diagnostic_capability_bps: 10_000,
                },
                0,
                RULES_VERSION,
            ),
            Err(FailureError::InvalidState(_))
        ));
    }

    #[test]
    fn shutdown_is_deliberate_and_restore_is_not_a_fault() {
        let facilities = facilities();
        let mut failure = state();
        failure
            .apply(
                &facilities,
                &FailureCommand::Intervene {
                    intervention_id: "intervention.shutdown".to_owned(),
                    component_id: "component.crusher".to_owned(),
                    kind: InterventionKind::Shutdown,
                    fault_instance_id: None,
                    derate_bps: None,
                    bypass: None,
                },
                0,
                RULES_VERSION,
            )
            .unwrap();
        assert!(failure.faults.is_empty());
        assert_eq!(
            failure
                .evaluate_facility(&facilities[0])
                .unwrap()
                .effective_capacity,
            0
        );
        failure
            .apply(
                &facilities,
                &FailureCommand::Intervene {
                    intervention_id: "intervention.restore".to_owned(),
                    component_id: "component.crusher".to_owned(),
                    kind: InterventionKind::Restore,
                    fault_instance_id: None,
                    derate_bps: None,
                    bypass: None,
                },
                0,
                RULES_VERSION,
            )
            .unwrap();
        assert_eq!(
            failure
                .evaluate_facility(&facilities[0])
                .unwrap()
                .effective_capacity,
            100
        );
    }
}
