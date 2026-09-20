use serde::{Deserialize, Serialize};
use std::collections::{BTreeMap, BTreeSet};
use std::fmt::{Display, Formatter};

pub const BASIS_POINTS_PER_WHOLE: u64 = 10_000;

#[derive(Debug, Clone, PartialEq, Eq)]
pub enum GraphError {
    InvalidId {
        kind: &'static str,
        id: String,
    },
    DuplicateFacilityId(String),
    DuplicateSystemId(String),
    DuplicateComponentId(String),
    DuplicateDependency(String),
    MissingReference {
        kind: &'static str,
        id: String,
    },
    ComponentAssignedToMultipleSystems(String),
    EmptyGraph(String),
    SelfDependency(String),
    DependencyCycle(Vec<String>),
    InvalidCapacity {
        kind: &'static str,
        id: String,
    },
    InvalidCapacityFactor {
        component_id: String,
        factor_bps: u16,
    },
    InvalidDependencyGroup(String),
    InvalidGraphState(String),
    ArithmeticOverflow,
}

impl Display for GraphError {
    fn fmt(&self, formatter: &mut Formatter<'_>) -> std::fmt::Result {
        write!(formatter, "{self:?}")
    }
}

impl std::error::Error for GraphError {}

#[derive(Debug, Clone, Copy, PartialEq, Eq, PartialOrd, Ord, Serialize, Deserialize)]
#[serde(deny_unknown_fields)]
pub enum DependencyType {
    #[serde(rename = "HARD")]
    Hard,
    #[serde(rename = "CAPACITY")]
    Capacity,
    #[serde(rename = "QUALITY")]
    Quality,
    #[serde(rename = "RELIABILITY")]
    Reliability,
    #[serde(rename = "COST")]
    Cost,
    #[serde(rename = "OPTIONAL")]
    Optional,
}

#[derive(Debug, Clone, Default, PartialEq, Eq, Serialize, Deserialize)]
#[serde(deny_unknown_fields)]
pub struct FacilityContext {
    pub world_id: Option<String>,
    pub region_id: Option<String>,
    pub site_id: Option<String>,
}

#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize)]
#[serde(deny_unknown_fields)]
pub struct Component {
    pub component_id: String,
    pub component_type: String,
    pub nominal_capacity: u64,
    pub available: bool,
    pub enabled: bool,
    pub capacity_factor_bps: u16,
}

impl Component {
    pub fn new(
        component_id: impl Into<String>,
        component_type: impl Into<String>,
        nominal_capacity: u64,
    ) -> Self {
        Self {
            component_id: component_id.into(),
            component_type: component_type.into(),
            nominal_capacity,
            available: true,
            enabled: true,
            capacity_factor_bps: BASIS_POINTS_PER_WHOLE as u16,
        }
    }
}

#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize)]
#[serde(deny_unknown_fields)]
pub struct System {
    pub system_id: String,
    pub system_type: String,
    pub nominal_capacity: u64,
    pub component_ids: Vec<String>,
    pub output_component_ids: Vec<String>,
}

impl System {
    pub fn new(
        system_id: impl Into<String>,
        system_type: impl Into<String>,
        nominal_capacity: u64,
        component_ids: Vec<String>,
    ) -> Self {
        Self {
            system_id: system_id.into(),
            system_type: system_type.into(),
            nominal_capacity,
            component_ids,
            output_component_ids: Vec::new(),
        }
    }

    pub fn with_outputs(mut self, output_component_ids: Vec<String>) -> Self {
        self.output_component_ids = output_component_ids;
        self
    }
}

#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize)]
#[serde(deny_unknown_fields)]
pub struct DependencyEdge {
    pub upstream_component_id: String,
    pub downstream_component_id: String,
    pub dependency_type: DependencyType,
    #[serde(default)]
    pub path_group: Option<String>,
}

impl DependencyEdge {
    pub fn new(
        upstream_component_id: impl Into<String>,
        downstream_component_id: impl Into<String>,
        dependency_type: DependencyType,
    ) -> Self {
        Self {
            upstream_component_id: upstream_component_id.into(),
            downstream_component_id: downstream_component_id.into(),
            dependency_type,
            path_group: None,
        }
    }

    pub fn with_path_group(mut self, path_group: impl Into<String>) -> Self {
        self.path_group = Some(path_group.into());
        self
    }
}

#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize)]
#[serde(deny_unknown_fields)]
pub struct Facility {
    pub facility_id: String,
    pub facility_type: String,
    pub context: FacilityContext,
    pub nominal_capacity: u64,
    pub systems: Vec<System>,
    pub components: Vec<Component>,
    pub dependencies: Vec<DependencyEdge>,
}

impl Facility {
    pub fn new(
        facility_id: impl Into<String>,
        facility_type: impl Into<String>,
        nominal_capacity: u64,
        systems: Vec<System>,
        components: Vec<Component>,
        dependencies: Vec<DependencyEdge>,
    ) -> Self {
        Self {
            facility_id: facility_id.into(),
            facility_type: facility_type.into(),
            context: FacilityContext::default(),
            nominal_capacity,
            systems,
            components,
            dependencies,
        }
    }

    pub fn validate(&self) -> Result<(), GraphError> {
        validate_id("facility", &self.facility_id)?;
        if self.facility_type.is_empty() {
            return Err(GraphError::InvalidGraphState(
                "facility type is required".to_owned(),
            ));
        }
        validate_capacity("facility", &self.facility_id, self.nominal_capacity)?;
        for id in [
            &self.context.world_id,
            &self.context.region_id,
            &self.context.site_id,
        ]
        .into_iter()
        .flatten()
        {
            validate_id("context", id)?;
        }
        if self.systems.is_empty() || self.components.is_empty() {
            return Err(GraphError::EmptyGraph(self.facility_id.clone()));
        }

        let mut component_map = BTreeMap::new();
        for component in &self.components {
            validate_id("component", &component.component_id)?;
            if component.component_type.is_empty() {
                return Err(GraphError::InvalidGraphState(
                    "component type is required".to_owned(),
                ));
            }
            validate_capacity(
                "component",
                &component.component_id,
                component.nominal_capacity,
            )?;
            if u64::from(component.capacity_factor_bps) > BASIS_POINTS_PER_WHOLE {
                return Err(GraphError::InvalidCapacityFactor {
                    component_id: component.component_id.clone(),
                    factor_bps: component.capacity_factor_bps,
                });
            }
            if component_map
                .insert(component.component_id.clone(), component)
                .is_some()
            {
                return Err(GraphError::DuplicateComponentId(
                    component.component_id.clone(),
                ));
            }
        }

        let mut system_map = BTreeMap::new();
        let mut assigned_components = BTreeMap::new();
        for system in &self.systems {
            validate_id("system", &system.system_id)?;
            if system.system_type.is_empty() {
                return Err(GraphError::InvalidGraphState(
                    "system type is required".to_owned(),
                ));
            }
            validate_capacity("system", &system.system_id, system.nominal_capacity)?;
            if system.component_ids.is_empty() {
                return Err(GraphError::EmptyGraph(system.system_id.clone()));
            }
            if system_map
                .insert(system.system_id.clone(), system)
                .is_some()
            {
                return Err(GraphError::DuplicateSystemId(system.system_id.clone()));
            }
            let mut system_components = BTreeSet::new();
            for component_id in &system.component_ids {
                if !system_components.insert(component_id) {
                    return Err(GraphError::InvalidGraphState(format!(
                        "component {component_id} is repeated in system {}",
                        system.system_id
                    )));
                }
                if !component_map.contains_key(component_id) {
                    return Err(GraphError::MissingReference {
                        kind: "component",
                        id: component_id.clone(),
                    });
                }
                if assigned_components
                    .insert(component_id.clone(), system.system_id.clone())
                    .is_some()
                {
                    return Err(GraphError::ComponentAssignedToMultipleSystems(
                        component_id.clone(),
                    ));
                }
            }
            let output_ids = if system.output_component_ids.is_empty() {
                Vec::new()
            } else {
                let mut outputs = BTreeSet::new();
                for component_id in &system.output_component_ids {
                    if !system_components.contains(component_id) {
                        return Err(GraphError::MissingReference {
                            kind: "system output component",
                            id: component_id.clone(),
                        });
                    }
                    if !outputs.insert(component_id) {
                        return Err(GraphError::InvalidGraphState(format!(
                            "system {} repeats output component {component_id}",
                            system.system_id
                        )));
                    }
                }
                system.output_component_ids.clone()
            };
            if output_ids.iter().any(|id| !system_components.contains(id)) {
                return Err(GraphError::InvalidGraphState(format!(
                    "system {} has an output outside its component set",
                    system.system_id
                )));
            }
        }
        if assigned_components.len() != component_map.len() {
            let missing = component_map
                .keys()
                .find(|component_id| !assigned_components.contains_key(*component_id))
                .cloned()
                .unwrap_or_default();
            return Err(GraphError::MissingReference {
                kind: "owning system for component",
                id: missing,
            });
        }

        let mut seen_edges = BTreeSet::new();
        let mut adjacency: BTreeMap<String, BTreeSet<String>> = BTreeMap::new();
        for edge in &self.dependencies {
            validate_id("upstream component", &edge.upstream_component_id)?;
            validate_id("downstream component", &edge.downstream_component_id)?;
            if edge.upstream_component_id == edge.downstream_component_id {
                return Err(GraphError::SelfDependency(
                    edge.upstream_component_id.clone(),
                ));
            }
            if !component_map.contains_key(&edge.upstream_component_id) {
                return Err(GraphError::MissingReference {
                    kind: "upstream component",
                    id: edge.upstream_component_id.clone(),
                });
            }
            if !component_map.contains_key(&edge.downstream_component_id) {
                return Err(GraphError::MissingReference {
                    kind: "downstream component",
                    id: edge.downstream_component_id.clone(),
                });
            }
            if edge.path_group.as_deref().is_some_and(str::is_empty) {
                return Err(GraphError::InvalidDependencyGroup(format!(
                    "{} -> {}",
                    edge.upstream_component_id, edge.downstream_component_id
                )));
            }
            let key = (
                edge.upstream_component_id.clone(),
                edge.downstream_component_id.clone(),
                edge.dependency_type,
                edge.path_group.clone(),
            );
            if !seen_edges.insert(key) {
                return Err(GraphError::DuplicateDependency(format!(
                    "{} -> {}",
                    edge.upstream_component_id, edge.downstream_component_id
                )));
            }
            adjacency
                .entry(edge.upstream_component_id.clone())
                .or_default()
                .insert(edge.downstream_component_id.clone());
        }
        detect_cycle(&component_map, &adjacency)?;
        Ok(())
    }

    pub fn evaluate(&self) -> Result<FacilityEvaluation, GraphError> {
        self.validate()?;
        let component_map: BTreeMap<_, _> = self
            .components
            .iter()
            .map(|component| (component.component_id.clone(), component))
            .collect();
        let topology = self.topological_order(&component_map)?;
        let incoming = incoming_edges(self);
        let mut component_evaluations = BTreeMap::new();

        for component_id in topology {
            let component =
                component_map
                    .get(&component_id)
                    .ok_or_else(|| GraphError::MissingReference {
                        kind: "component",
                        id: component_id.clone(),
                    })?;
            let local_capacity = local_capacity(component)?;
            let mut effective_capacity = local_capacity;
            let mut hard_blocked = false;
            let mut capacity_blocked = false;
            let edges = incoming.get(&component_id).cloned().unwrap_or_default();

            for (group, group_edges) in grouped_edges(&edges, DependencyType::Hard) {
                let maximum = group_edges
                    .iter()
                    .map(|edge| upstream_capacity(&component_evaluations, edge))
                    .max()
                    .unwrap_or(0);
                if group.is_some() {
                    hard_blocked |= maximum == 0;
                } else {
                    hard_blocked |= group_edges
                        .iter()
                        .any(|edge| upstream_capacity(&component_evaluations, edge) == 0);
                }
            }
            if hard_blocked {
                effective_capacity = 0;
            }

            for (group, group_edges) in grouped_edges(&edges, DependencyType::Capacity) {
                let constraint = if group.is_some() {
                    group_edges
                        .iter()
                        .map(|edge| upstream_capacity(&component_evaluations, edge))
                        .max()
                        .unwrap_or(0)
                } else {
                    group_edges
                        .iter()
                        .map(|edge| upstream_capacity(&component_evaluations, edge))
                        .min()
                        .unwrap_or(effective_capacity)
                };
                if constraint < effective_capacity {
                    capacity_blocked = true;
                    effective_capacity = constraint;
                }
            }

            component_evaluations.insert(
                component_id,
                ComponentEvaluation {
                    local_capacity,
                    effective_capacity,
                    hard_blocked,
                    capacity_blocked,
                },
            );
        }

        let mut system_evaluations = Vec::new();
        for system in &self.systems {
            let output_ids = output_components(system, self, &incoming)?;
            let output_capacity = output_ids
                .iter()
                .map(|id| component_evaluations[id].effective_capacity)
                .max()
                .unwrap_or(0);
            let effective_capacity = output_capacity.min(system.nominal_capacity);
            let mut bottlenecks = Vec::new();
            for component_id in &system.component_ids {
                let component = component_map[component_id];
                let evaluation = &component_evaluations[component_id];
                if evaluation.effective_capacity == effective_capacity
                    && evaluation.effective_capacity < component.nominal_capacity
                {
                    bottlenecks.push(BottleneckEvidence {
                        component_id: component_id.clone(),
                        reason: bottleneck_reason(component, evaluation),
                        effective_capacity: evaluation.effective_capacity,
                    });
                }
            }
            if bottlenecks.is_empty() && effective_capacity < system.nominal_capacity {
                if let Some(component_id) = output_ids.first() {
                    bottlenecks.push(BottleneckEvidence {
                        component_id: component_id.clone(),
                        reason: BottleneckReason::SystemCapacity,
                        effective_capacity,
                    });
                }
            }
            bottlenecks.sort_by(|left, right| left.component_id.cmp(&right.component_id));
            system_evaluations.push(SystemEvaluation {
                system_id: system.system_id.clone(),
                operational_state: operational_state(effective_capacity, system.nominal_capacity),
                nominal_capacity: system.nominal_capacity,
                effective_capacity,
                bottlenecks,
            });
        }

        let effective_capacity = system_evaluations
            .iter()
            .map(|evaluation| evaluation.effective_capacity)
            .min()
            .unwrap_or(0)
            .min(self.nominal_capacity);
        let mut bottlenecks = system_evaluations
            .iter()
            .filter(|evaluation| evaluation.effective_capacity == effective_capacity)
            .flat_map(|evaluation| evaluation.bottlenecks.clone())
            .collect::<Vec<_>>();
        bottlenecks.sort_by(|left, right| {
            left.component_id
                .cmp(&right.component_id)
                .then(left.reason.cmp(&right.reason))
        });
        bottlenecks.dedup();
        if bottlenecks.is_empty() && effective_capacity < self.nominal_capacity {
            bottlenecks.push(BottleneckEvidence {
                component_id: self
                    .systems
                    .first()
                    .and_then(|system| system.component_ids.first())
                    .cloned()
                    .unwrap_or_default(),
                reason: BottleneckReason::SystemCapacity,
                effective_capacity,
            });
        }
        Ok(FacilityEvaluation {
            facility_id: self.facility_id.clone(),
            operational_state: operational_state(effective_capacity, self.nominal_capacity),
            nominal_capacity: self.nominal_capacity,
            effective_capacity,
            system_evaluations,
            bottlenecks,
        })
    }

    fn topological_order(
        &self,
        component_map: &BTreeMap<String, &Component>,
    ) -> Result<Vec<String>, GraphError> {
        let mut indegree: BTreeMap<String, usize> =
            component_map.keys().map(|id| (id.clone(), 0)).collect();
        let mut adjacency: BTreeMap<String, BTreeSet<String>> = BTreeMap::new();
        for edge in &self.dependencies {
            adjacency
                .entry(edge.upstream_component_id.clone())
                .or_default()
                .insert(edge.downstream_component_id.clone());
        }
        for downstream_ids in adjacency.values() {
            for downstream_id in downstream_ids {
                *indegree.get_mut(downstream_id).ok_or_else(|| {
                    GraphError::MissingReference {
                        kind: "downstream component",
                        id: downstream_id.clone(),
                    }
                })? += 1;
            }
        }
        let mut ready = indegree
            .iter()
            .filter_map(|(id, degree)| (*degree == 0).then_some(id.clone()))
            .collect::<BTreeSet<_>>();
        let mut order = Vec::with_capacity(component_map.len());
        while let Some(id) = ready.pop_first() {
            order.push(id.clone());
            if let Some(downstream_ids) = adjacency.get(&id) {
                for downstream_id in downstream_ids {
                    let degree = indegree.get_mut(downstream_id).expect("validated node");
                    *degree -= 1;
                    if *degree == 0 {
                        ready.insert(downstream_id.clone());
                    }
                }
            }
        }
        if order.len() != component_map.len() {
            return Err(GraphError::DependencyCycle(order));
        }
        Ok(order)
    }
}

#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize)]
#[serde(deny_unknown_fields)]
pub enum OperationalState {
    Operational,
    Constrained,
    Unavailable,
}

#[derive(Debug, Clone, Copy, PartialEq, Eq, PartialOrd, Ord, Serialize, Deserialize)]
#[serde(deny_unknown_fields)]
pub enum BottleneckReason {
    ComponentUnavailable,
    ComponentCapacity,
    HardDependencyUnavailable,
    CapacityDependency,
    SystemCapacity,
}

#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize)]
#[serde(deny_unknown_fields)]
pub struct BottleneckEvidence {
    pub component_id: String,
    pub reason: BottleneckReason,
    pub effective_capacity: u64,
}

#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize)]
#[serde(deny_unknown_fields)]
pub struct SystemEvaluation {
    pub system_id: String,
    pub operational_state: OperationalState,
    pub nominal_capacity: u64,
    pub effective_capacity: u64,
    pub bottlenecks: Vec<BottleneckEvidence>,
}

#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize)]
#[serde(deny_unknown_fields)]
pub struct FacilityEvaluation {
    pub facility_id: String,
    pub operational_state: OperationalState,
    pub nominal_capacity: u64,
    pub effective_capacity: u64,
    pub system_evaluations: Vec<SystemEvaluation>,
    pub bottlenecks: Vec<BottleneckEvidence>,
}

#[derive(Debug, Clone, Copy)]
struct ComponentEvaluation {
    local_capacity: u64,
    effective_capacity: u64,
    hard_blocked: bool,
    capacity_blocked: bool,
}

fn validate_facility_ids(facilities: &[Facility]) -> Result<(), GraphError> {
    let mut ids = BTreeSet::new();
    for facility in facilities {
        if !ids.insert(&facility.facility_id) {
            return Err(GraphError::DuplicateFacilityId(
                facility.facility_id.clone(),
            ));
        }
        facility.validate()?;
    }
    Ok(())
}

pub fn validate_facilities(facilities: &[Facility]) -> Result<(), GraphError> {
    validate_facility_ids(facilities)
}

fn validate_id(kind: &'static str, id: &str) -> Result<(), GraphError> {
    let valid = !id.is_empty()
        && id
            .chars()
            .next()
            .is_some_and(|character| character.is_ascii_lowercase())
        && id.chars().all(|character| {
            character.is_ascii_lowercase()
                || character.is_ascii_digit()
                || "._-".contains(character)
        });
    if valid {
        Ok(())
    } else {
        Err(GraphError::InvalidId {
            kind,
            id: id.to_owned(),
        })
    }
}

fn validate_capacity(kind: &'static str, id: &str, capacity: u64) -> Result<(), GraphError> {
    if capacity == 0 {
        return Err(GraphError::InvalidCapacity {
            kind,
            id: id.to_owned(),
        });
    }
    Ok(())
}

fn local_capacity(component: &Component) -> Result<u64, GraphError> {
    if !component.available || !component.enabled {
        return Ok(0);
    }
    component
        .nominal_capacity
        .checked_mul(u64::from(component.capacity_factor_bps))
        .map(|capacity| capacity / BASIS_POINTS_PER_WHOLE)
        .ok_or(GraphError::ArithmeticOverflow)
}

fn incoming_edges(facility: &Facility) -> BTreeMap<String, Vec<&DependencyEdge>> {
    let mut incoming = BTreeMap::new();
    for edge in &facility.dependencies {
        incoming
            .entry(edge.downstream_component_id.clone())
            .or_insert_with(Vec::new)
            .push(edge);
    }
    incoming
}

fn grouped_edges<'a>(
    edges: &[&'a DependencyEdge],
    dependency_type: DependencyType,
) -> BTreeMap<Option<String>, Vec<&'a DependencyEdge>> {
    let mut grouped = BTreeMap::new();
    for edge in edges
        .iter()
        .copied()
        .filter(|edge| edge.dependency_type == dependency_type)
    {
        grouped
            .entry(edge.path_group.clone())
            .or_insert_with(Vec::new)
            .push(edge);
    }
    grouped
}

fn upstream_capacity(
    evaluations: &BTreeMap<String, ComponentEvaluation>,
    edge: &DependencyEdge,
) -> u64 {
    evaluations
        .get(&edge.upstream_component_id)
        .map_or(0, |evaluation| evaluation.effective_capacity)
}

fn output_components(
    system: &System,
    facility: &Facility,
    incoming: &BTreeMap<String, Vec<&DependencyEdge>>,
) -> Result<Vec<String>, GraphError> {
    if !system.output_component_ids.is_empty() {
        return Ok(system.output_component_ids.clone());
    }
    let system_components = system.component_ids.iter().collect::<BTreeSet<_>>();
    let mut outputs = system
        .component_ids
        .iter()
        .filter(|component_id| {
            !facility.dependencies.iter().any(|edge| {
                edge.upstream_component_id == **component_id
                    && system_components.contains(&edge.downstream_component_id)
            })
        })
        .cloned()
        .collect::<Vec<_>>();
    if outputs.is_empty() && !incoming.is_empty() {
        return Err(GraphError::InvalidGraphState(format!(
            "system {} has no terminal output component",
            system.system_id
        )));
    }
    outputs.sort();
    Ok(outputs)
}

fn bottleneck_reason(component: &Component, evaluation: &ComponentEvaluation) -> BottleneckReason {
    if !component.available || !component.enabled {
        BottleneckReason::ComponentUnavailable
    } else if evaluation.hard_blocked {
        BottleneckReason::HardDependencyUnavailable
    } else if evaluation.capacity_blocked {
        BottleneckReason::CapacityDependency
    } else if evaluation.local_capacity < component.nominal_capacity {
        BottleneckReason::ComponentCapacity
    } else {
        BottleneckReason::SystemCapacity
    }
}

fn operational_state(effective_capacity: u64, nominal_capacity: u64) -> OperationalState {
    if effective_capacity == 0 {
        OperationalState::Unavailable
    } else if effective_capacity < nominal_capacity {
        OperationalState::Constrained
    } else {
        OperationalState::Operational
    }
}

fn detect_cycle(
    components: &BTreeMap<String, &Component>,
    adjacency: &BTreeMap<String, BTreeSet<String>>,
) -> Result<(), GraphError> {
    let mut visiting = BTreeSet::new();
    let mut visited = BTreeSet::new();
    let mut stack = Vec::new();
    for component_id in components.keys() {
        visit_cycle(
            component_id,
            adjacency,
            &mut visiting,
            &mut visited,
            &mut stack,
        )?;
    }
    Ok(())
}

fn visit_cycle(
    component_id: &str,
    adjacency: &BTreeMap<String, BTreeSet<String>>,
    visiting: &mut BTreeSet<String>,
    visited: &mut BTreeSet<String>,
    stack: &mut Vec<String>,
) -> Result<(), GraphError> {
    if visited.contains(component_id) {
        return Ok(());
    }
    if visiting.contains(component_id) {
        let start = stack.iter().position(|id| id == component_id).unwrap_or(0);
        return Err(GraphError::DependencyCycle(stack[start..].to_vec()));
    }
    visiting.insert(component_id.to_owned());
    stack.push(component_id.to_owned());
    if let Some(children) = adjacency.get(component_id) {
        for child in children {
            visit_cycle(child, adjacency, visiting, visited, stack)?;
        }
    }
    stack.pop();
    visiting.remove(component_id);
    visited.insert(component_id.to_owned());
    Ok(())
}

pub fn aggregate_plant_fixture() -> Facility {
    let components = vec![
        Component::new("component.feed_stockpile", "stockpile", 100),
        Component::new("component.feed_conveyor", "conveyor", 100),
        Component::new("component.crusher", "crusher", 100),
        Component::new("component.screen", "screen", 100),
        Component::new("component.output_conveyor", "conveyor", 100),
        Component::new("component.finished_stockpile", "stockpile", 100),
    ];
    let system = System::new(
        "system.aggregate_line",
        "aggregate_processing",
        100,
        components
            .iter()
            .map(|component| component.component_id.clone())
            .collect(),
    )
    .with_outputs(vec!["component.finished_stockpile".to_owned()]);
    let edge_pairs = [
        ("component.feed_stockpile", "component.feed_conveyor"),
        ("component.feed_conveyor", "component.crusher"),
        ("component.crusher", "component.screen"),
        ("component.screen", "component.output_conveyor"),
        ("component.output_conveyor", "component.finished_stockpile"),
    ];
    let mut dependencies = Vec::new();
    for (upstream, downstream) in edge_pairs {
        dependencies.push(DependencyEdge::new(
            upstream,
            downstream,
            DependencyType::Hard,
        ));
        dependencies.push(DependencyEdge::new(
            upstream,
            downstream,
            DependencyType::Capacity,
        ));
    }
    Facility::new(
        "facility.aggregate_plant_fixture",
        "aggregate_processing_fixture",
        100,
        vec![system],
        components,
        dependencies,
    )
}

#[cfg(test)]
mod tests {
    use super::*;

    fn fixture() -> Facility {
        aggregate_plant_fixture()
    }

    #[test]
    fn aggregate_fixture_is_healthy_and_has_no_bottleneck() {
        let evaluation = fixture().evaluate().unwrap();
        assert_eq!(evaluation.operational_state, OperationalState::Operational);
        assert_eq!(evaluation.effective_capacity, 100);
        assert!(evaluation.bottlenecks.is_empty());
    }

    #[test]
    fn hard_unavailability_propagates_to_downstream_output() {
        let mut facility = fixture();
        facility.components[1].available = false;
        let evaluation = facility.evaluate().unwrap();
        assert_eq!(evaluation.operational_state, OperationalState::Unavailable);
        assert_eq!(evaluation.effective_capacity, 0);
        assert!(evaluation
            .bottlenecks
            .iter()
            .any(|evidence| evidence.reason == BottleneckReason::ComponentUnavailable));
    }

    #[test]
    fn partial_capacity_remains_operational_and_identifies_limiting_component() {
        let mut facility = fixture();
        facility.components[2].capacity_factor_bps = 300;
        let evaluation = facility.evaluate().unwrap();
        assert_eq!(evaluation.operational_state, OperationalState::Constrained);
        assert_eq!(evaluation.effective_capacity, 3);
        assert!(evaluation
            .bottlenecks
            .iter()
            .any(|evidence| evidence.component_id == "component.crusher"));
    }

    #[test]
    fn optional_branch_does_not_stop_the_primary_line() {
        let mut facility = fixture();
        facility
            .components
            .push(Component::new("component.optional_monitor", "monitor", 100));
        facility.systems[0]
            .component_ids
            .push("component.optional_monitor".to_owned());
        facility.dependencies.push(
            DependencyEdge::new(
                "component.optional_monitor",
                "component.finished_stockpile",
                DependencyType::Optional,
            )
            .with_path_group("optional.monitor"),
        );
        facility.components[6].available = false;
        let evaluation = facility.evaluate().unwrap();
        assert_eq!(evaluation.effective_capacity, 100);
        assert_eq!(evaluation.operational_state, OperationalState::Operational);
    }

    #[test]
    fn parallel_capacity_group_preserves_an_alternate_path() {
        let components = vec![
            Component::new("component.source_a", "source", 100),
            Component::new("component.source_b", "source", 100),
            Component::new("component.output", "output", 100),
        ];
        let system = System::new(
            "system.parallel",
            "parallel",
            100,
            components
                .iter()
                .map(|component| component.component_id.clone())
                .collect(),
        )
        .with_outputs(vec!["component.output".to_owned()]);
        let dependencies = vec![
            DependencyEdge::new(
                "component.source_a",
                "component.output",
                DependencyType::Hard,
            )
            .with_path_group("route"),
            DependencyEdge::new(
                "component.source_b",
                "component.output",
                DependencyType::Hard,
            )
            .with_path_group("route"),
            DependencyEdge::new(
                "component.source_a",
                "component.output",
                DependencyType::Capacity,
            )
            .with_path_group("route"),
            DependencyEdge::new(
                "component.source_b",
                "component.output",
                DependencyType::Capacity,
            )
            .with_path_group("route"),
        ];
        let mut facility = Facility::new(
            "facility.parallel_fixture",
            "parallel_fixture",
            100,
            vec![system],
            components,
            dependencies,
        );
        facility.components[0].available = false;
        let evaluation = facility.evaluate().unwrap();
        assert_eq!(evaluation.effective_capacity, 100);
        assert_eq!(evaluation.operational_state, OperationalState::Operational);
    }

    #[test]
    fn invalid_cycle_duplicate_and_missing_reference_are_typed_errors() {
        let mut cycle = fixture();
        cycle.dependencies.push(DependencyEdge::new(
            "component.finished_stockpile",
            "component.feed_stockpile",
            DependencyType::Hard,
        ));
        assert!(matches!(
            cycle.validate(),
            Err(GraphError::DependencyCycle(_))
        ));

        let mut duplicate = fixture();
        duplicate.components.push(duplicate.components[0].clone());
        assert_eq!(
            duplicate.validate(),
            Err(GraphError::DuplicateComponentId(
                "component.feed_stockpile".to_owned()
            ))
        );

        let mut missing = fixture();
        missing.dependencies.push(DependencyEdge::new(
            "component.unknown",
            "component.crusher",
            DependencyType::Hard,
        ));
        assert!(matches!(
            missing.validate(),
            Err(GraphError::MissingReference {
                kind: "upstream component",
                ..
            })
        ));
    }

    #[test]
    fn all_dependency_types_have_explicit_serialized_contracts() {
        let types = [
            DependencyType::Hard,
            DependencyType::Capacity,
            DependencyType::Quality,
            DependencyType::Reliability,
            DependencyType::Cost,
            DependencyType::Optional,
        ];
        let json = serde_json::to_string(&types).unwrap();
        for dependency_type in [
            "HARD",
            "CAPACITY",
            "QUALITY",
            "RELIABILITY",
            "COST",
            "OPTIONAL",
        ] {
            assert!(json.contains(dependency_type));
        }
    }

    #[test]
    fn graph_serialization_and_evaluation_are_deterministic() {
        let first = fixture();
        let mut second = fixture();
        second.components.reverse();
        second.systems[0].component_ids.reverse();
        second.dependencies.reverse();
        let first_json = serde_json::to_string(&first).unwrap();
        let second_json = serde_json::to_string(&second).unwrap();
        let first_evaluation = first.evaluate().unwrap();
        let second_evaluation = second.evaluate().unwrap();
        assert_eq!(first_evaluation, second_evaluation);
        assert_ne!(first_json, second_json);
    }
}
