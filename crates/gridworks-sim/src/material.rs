use crate::failure::FailureState;
use crate::graph::is_valid_semantic_id;
use crate::graph::Facility;
use serde::{Deserialize, Serialize};
use std::collections::{BTreeMap, BTreeSet};
use std::fmt::{Display, Formatter};

pub const QUANTITY_SCALE: u64 = 1;

#[derive(Debug, Clone, PartialEq, Eq)]
pub enum MaterialError {
    InvalidIdentifier(String),
    DuplicateResource(String),
    DuplicateRecipe(String),
    DuplicateInventory(String),
    DuplicateLot(String),
    UnknownResource(String),
    UnknownRecipe(String),
    UnknownInventory(String),
    GradeMismatch {
        resource_id: String,
        grade_id: String,
    },
    InvalidQuantity,
    QuantityOverflow,
    QuantityUnderflow,
    CapacityExceeded(String),
    InvalidRecipe(String),
    InvalidTransfer(String),
    InvalidProductionTarget(String),
    FacilityUnavailable,
    InvalidState(String),
}
impl Display for MaterialError {
    fn fmt(&self, f: &mut Formatter<'_>) -> std::fmt::Result {
        write!(f, "{self:?}")
    }
}
impl std::error::Error for MaterialError {}

#[derive(Debug, Clone, Copy, PartialEq, Eq, PartialOrd, Ord, Serialize, Deserialize)]
#[serde(deny_unknown_fields)]
pub enum ResourceCategory {
    RawMaterial,
    Intermediate,
    FinishedGood,
    Waste,
    Byproduct,
}
#[derive(Debug, Clone, Copy, PartialEq, Eq, PartialOrd, Ord, Serialize, Deserialize)]
#[serde(deny_unknown_fields)]
pub enum QuantityUnit {
    BaseQuantity,
}
#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize)]
#[serde(deny_unknown_fields)]
pub struct ResourceDefinition {
    pub resource_id: String,
    pub category: ResourceCategory,
    pub unit: QuantityUnit,
    pub grade_ids: Vec<String>,
    pub rules_version: String,
}
#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize)]
#[serde(deny_unknown_fields)]
pub struct MaterialLot {
    pub resource_id: String,
    pub grade_id: String,
    pub quantity: u64,
}
#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize)]
#[serde(deny_unknown_fields)]
pub struct InventoryStore {
    pub inventory_id: String,
    pub facility_id: Option<String>,
    pub capacity: u64,
    pub permitted_resource_ids: Vec<String>,
    pub lots: Vec<MaterialLot>,
}
impl InventoryStore {
    pub fn quantity(&self, resource: &str, grade: &str) -> u64 {
        self.lots
            .iter()
            .filter(|l| l.resource_id == resource && l.grade_id == grade)
            .map(|l| l.quantity)
            .sum()
    }
    pub fn total_quantity(&self) -> u64 {
        self.lots.iter().map(|l| l.quantity).sum()
    }
    pub fn free_capacity(&self) -> u64 {
        self.capacity.saturating_sub(self.total_quantity())
    }
    fn add(&mut self, resource: &str, grade: &str, quantity: u64) -> Result<(), MaterialError> {
        if quantity == 0 {
            return Err(MaterialError::InvalidQuantity);
        }
        if let Some(lot) = self
            .lots
            .iter_mut()
            .find(|l| l.resource_id == resource && l.grade_id == grade)
        {
            lot.quantity = lot
                .quantity
                .checked_add(quantity)
                .ok_or(MaterialError::QuantityOverflow)?;
        } else {
            self.lots.push(MaterialLot {
                resource_id: resource.to_owned(),
                grade_id: grade.to_owned(),
                quantity,
            });
        }
        Ok(())
    }
    fn remove(&mut self, resource: &str, grade: &str, quantity: u64) -> Result<(), MaterialError> {
        if quantity == 0 {
            return Err(MaterialError::InvalidQuantity);
        }
        let lot = self
            .lots
            .iter_mut()
            .find(|l| l.resource_id == resource && l.grade_id == grade)
            .ok_or(MaterialError::QuantityUnderflow)?;
        if lot.quantity < quantity {
            return Err(MaterialError::QuantityUnderflow);
        }
        lot.quantity -= quantity;
        if lot.quantity == 0 {
            self.lots
                .retain(|l| !(l.resource_id == resource && l.grade_id == grade));
        }
        Ok(())
    }
}
#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize)]
#[serde(deny_unknown_fields)]
pub struct RecipeInput {
    pub resource_id: String,
    pub grade_id: String,
    pub quantity_per_run: u64,
}
#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize)]
#[serde(deny_unknown_fields)]
pub struct RecipeOutput {
    pub resource_id: String,
    pub grade_id: String,
    pub quantity_per_run: u64,
    pub yield_bps: u16,
}
#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize)]
#[serde(deny_unknown_fields)]
pub struct RecipeDefinition {
    pub recipe_id: String,
    pub facility_type_ids: Vec<String>,
    pub inputs: Vec<RecipeInput>,
    pub outputs: Vec<RecipeOutput>,
    pub rules_version: String,
}
#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize)]
#[serde(deny_unknown_fields)]
pub struct InputSource {
    pub resource_id: String,
    pub store_id: String,
}
#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize)]
#[serde(deny_unknown_fields)]
pub enum MaterialCommand {
    Transfer {
        transfer_id: String,
        source_inventory_id: String,
        destination_inventory_id: String,
        resource_id: String,
        grade_id: String,
        quantity: u64,
    },
    ExecuteProduction {
        run_id: String,
        recipe_id: String,
        facility_id: String,
        input_sources: Vec<InputSource>,
        output_inventory_id: String,
        requested_runs: u64,
    },
}
#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize)]
#[serde(deny_unknown_fields)]
pub enum MaterialLimit {
    None,
    FacilityCapacity,
    MissingInput,
    InsufficientInput,
    DestinationCapacity,
    FacilityUnavailable,
}
#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize)]
#[serde(deny_unknown_fields)]
pub enum MaterialEvent {
    InventoryTransferred {
        transfer_id: String,
        resource_id: String,
        grade_id: String,
        accepted_quantity: u64,
        limit: MaterialLimit,
    },
    ProductionExecuted {
        run_id: String,
        recipe_id: String,
        accepted_runs: u64,
        limit: MaterialLimit,
        input_quantity: u64,
        output_quantity: u64,
    },
}
#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize, Default)]
#[serde(deny_unknown_fields)]
pub struct MaterialState {
    pub resources: Vec<ResourceDefinition>,
    pub recipes: Vec<RecipeDefinition>,
    pub inventories: Vec<InventoryStore>,
}
impl MaterialState {
    pub fn canonicalized(&self) -> Self {
        let mut s = self.clone();
        s.resources
            .sort_by(|a, b| a.resource_id.cmp(&b.resource_id));
        s.recipes.sort_by(|a, b| a.recipe_id.cmp(&b.recipe_id));
        for r in &mut s.recipes {
            r.inputs.sort_by(|a, b| {
                a.resource_id
                    .cmp(&b.resource_id)
                    .then(a.grade_id.cmp(&b.grade_id))
            });
            r.outputs.sort_by(|a, b| {
                a.resource_id
                    .cmp(&b.resource_id)
                    .then(a.grade_id.cmp(&b.grade_id))
            });
        }
        s.inventories
            .sort_by(|a, b| a.inventory_id.cmp(&b.inventory_id));
        for i in &mut s.inventories {
            i.permitted_resource_ids.sort();
            i.lots.sort_by(|a, b| {
                a.resource_id
                    .cmp(&b.resource_id)
                    .then(a.grade_id.cmp(&b.grade_id))
            });
        }
        s
    }
    pub fn validate(&self, facilities: &[Facility]) -> Result<(), MaterialError> {
        let resources = self.resource_map()?;
        let mut recipe_ids = BTreeSet::new();
        for r in &self.recipes {
            validate_id(&r.recipe_id)?;
            if !recipe_ids.insert(&r.recipe_id)
                || r.inputs.is_empty()
                || r.outputs.is_empty()
                || r.facility_type_ids.is_empty()
            {
                return Err(MaterialError::InvalidRecipe(r.recipe_id.clone()));
            }
            for i in &r.inputs {
                validate_id(&i.resource_id)?;
                validate_id(&i.grade_id)?;
                if i.quantity_per_run == 0 {
                    return Err(MaterialError::InvalidQuantity);
                }
                grade(&resources, &i.resource_id, &i.grade_id)?;
            }
            for o in &r.outputs {
                validate_id(&o.resource_id)?;
                validate_id(&o.grade_id)?;
                if o.quantity_per_run == 0 || o.yield_bps > 10_000 {
                    return Err(MaterialError::InvalidRecipe(r.recipe_id.clone()));
                }
                grade(&resources, &o.resource_id, &o.grade_id)?;
            }
        }
        let facility_ids = facilities
            .iter()
            .map(|f| f.facility_id.as_str())
            .collect::<BTreeSet<_>>();
        let mut inventory_ids = BTreeSet::new();
        for i in &self.inventories {
            validate_id(&i.inventory_id)?;
            if i.capacity == 0 || !inventory_ids.insert(&i.inventory_id) {
                return Err(MaterialError::DuplicateInventory(i.inventory_id.clone()));
            }
            if let Some(f) = &i.facility_id {
                if !facility_ids.contains(f.as_str()) {
                    return Err(MaterialError::InvalidState(i.inventory_id.clone()));
                }
            }
            for r in &i.permitted_resource_ids {
                if !resources.contains_key(r.as_str()) {
                    return Err(MaterialError::UnknownResource(r.clone()));
                }
            }
            let mut lots = BTreeSet::new();
            for l in &i.lots {
                if l.quantity == 0 {
                    return Err(MaterialError::InvalidQuantity);
                }
                grade(&resources, &l.resource_id, &l.grade_id)?;
                if !lots.insert((l.resource_id.as_str(), l.grade_id.as_str())) {
                    return Err(MaterialError::DuplicateLot(l.resource_id.clone()));
                }
            }
            if i.total_quantity() > i.capacity {
                return Err(MaterialError::CapacityExceeded(i.inventory_id.clone()));
            }
        }
        Ok(())
    }
    pub fn execute(
        &mut self,
        facilities: &[Facility],
        failure: &FailureState,
        command: &MaterialCommand,
    ) -> Result<MaterialEvent, MaterialError> {
        match command {
            MaterialCommand::Transfer {
                transfer_id,
                source_inventory_id,
                destination_inventory_id,
                resource_id,
                grade_id,
                quantity,
            } => self.transfer(
                transfer_id,
                source_inventory_id,
                destination_inventory_id,
                resource_id,
                grade_id,
                *quantity,
            ),
            MaterialCommand::ExecuteProduction {
                run_id,
                recipe_id,
                facility_id,
                input_sources,
                output_inventory_id,
                requested_runs,
            } => self.produce(
                facilities,
                failure,
                run_id,
                recipe_id,
                facility_id,
                input_sources,
                output_inventory_id,
                *requested_runs,
            ),
        }
    }
    fn resource_map(&self) -> Result<BTreeMap<&str, &ResourceDefinition>, MaterialError> {
        let mut m = BTreeMap::new();
        for r in &self.resources {
            validate_id(&r.resource_id)?;
            if r.grade_ids.is_empty() || m.insert(r.resource_id.as_str(), r).is_some() {
                return Err(MaterialError::DuplicateResource(r.resource_id.clone()));
            }
            for g in &r.grade_ids {
                validate_id(g)?;
            }
        }
        Ok(m)
    }
    fn index(&self, id: &str) -> Result<usize, MaterialError> {
        self.inventories
            .iter()
            .position(|i| i.inventory_id == id)
            .ok_or_else(|| MaterialError::UnknownInventory(id.to_owned()))
    }
    fn transfer(
        &mut self,
        id: &str,
        source: &str,
        destination: &str,
        resource: &str,
        grade_id: &str,
        requested: u64,
    ) -> Result<MaterialEvent, MaterialError> {
        validate_id(id)?;
        if requested == 0 {
            return Err(MaterialError::InvalidQuantity);
        };
        let m = self.resource_map()?;
        grade(&m, resource, grade_id)?;
        let si = self.index(source)?;
        let di = self.index(destination)?;
        if si == di {
            return Err(MaterialError::InvalidTransfer(
                "source equals destination".to_owned(),
            ));
        }
        let available = self.inventories[si].quantity(resource, grade_id);
        let space = self.inventories[di].free_capacity();
        let accepted = requested.min(available).min(space);
        let limit = if accepted == requested {
            MaterialLimit::None
        } else if available == 0 {
            MaterialLimit::MissingInput
        } else if space == 0 {
            MaterialLimit::DestinationCapacity
        } else {
            MaterialLimit::InsufficientInput
        };
        if accepted > 0 {
            self.inventories[si].remove(resource, grade_id, accepted)?;
            self.inventories[di].add(resource, grade_id, accepted)?;
        }
        Ok(MaterialEvent::InventoryTransferred {
            transfer_id: id.to_owned(),
            resource_id: resource.to_owned(),
            grade_id: grade_id.to_owned(),
            accepted_quantity: accepted,
            limit,
        })
    }
    #[allow(clippy::too_many_arguments)]
    fn produce(
        &mut self,
        facilities: &[Facility],
        failure: &FailureState,
        run_id: &str,
        recipe_id: &str,
        facility_id: &str,
        sources: &[InputSource],
        output_id: &str,
        requested: u64,
    ) -> Result<MaterialEvent, MaterialError> {
        validate_id(run_id)?;
        if requested == 0 {
            return Err(MaterialError::InvalidQuantity);
        };
        let recipe = self
            .recipes
            .iter()
            .find(|r| r.recipe_id == recipe_id)
            .cloned()
            .ok_or_else(|| MaterialError::UnknownRecipe(recipe_id.to_owned()))?;
        let facility = facilities
            .iter()
            .find(|f| f.facility_id == facility_id)
            .ok_or_else(|| MaterialError::InvalidProductionTarget(facility_id.to_owned()))?;
        if !recipe.facility_type_ids.contains(&facility.facility_type) {
            return Err(MaterialError::InvalidProductionTarget(
                facility_id.to_owned(),
            ));
        }
        let eval = failure
            .evaluate_facility(facility)
            .map_err(|_| MaterialError::FacilityUnavailable)?;
        let mut accepted = requested
            .checked_mul(eval.effective_capacity)
            .ok_or(MaterialError::QuantityOverflow)?
            .checked_div(eval.nominal_capacity)
            .ok_or(MaterialError::QuantityUnderflow)?;
        let mut limit = if eval.effective_capacity == 0 {
            MaterialLimit::FacilityUnavailable
        } else if accepted < requested {
            MaterialLimit::FacilityCapacity
        } else {
            MaterialLimit::None
        };
        let source_map = sources
            .iter()
            .map(|s| (s.resource_id.as_str(), s.store_id.as_str()))
            .collect::<BTreeMap<_, _>>();
        for input in &recipe.inputs {
            let store_id = source_map.get(input.resource_id.as_str()).ok_or_else(|| {
                MaterialError::InvalidTransfer("missing recipe input source".to_owned())
            })?;
            let store = self
                .inventories
                .iter()
                .find(|s| s.inventory_id == *store_id)
                .ok_or_else(|| MaterialError::UnknownInventory((*store_id).to_owned()))?;
            let possible =
                store.quantity(&input.resource_id, &input.grade_id) / input.quantity_per_run;
            if possible < accepted {
                accepted = possible;
                limit = MaterialLimit::InsufficientInput;
            }
        }
        let oi = self.index(output_id)?;
        let output_per_run = recipe
            .outputs
            .iter()
            .map(|o| o.quantity_per_run.saturating_mul(u64::from(o.yield_bps)) / 10_000)
            .sum::<u64>();
        let possible_output = if output_per_run == 0 {
            0
        } else {
            self.inventories[oi].free_capacity() / output_per_run
        };
        if possible_output < accepted {
            accepted = possible_output;
            limit = MaterialLimit::DestinationCapacity;
        }
        if accepted == 0 {
            return Ok(MaterialEvent::ProductionExecuted {
                run_id: run_id.to_owned(),
                recipe_id: recipe_id.to_owned(),
                accepted_runs: 0,
                limit,
                input_quantity: 0,
                output_quantity: 0,
            });
        }
        for input in &recipe.inputs {
            let store_id = source_map.get(input.resource_id.as_str()).ok_or_else(|| {
                MaterialError::InvalidTransfer("missing recipe input source".to_owned())
            })?;
            let store = self
                .inventories
                .iter_mut()
                .find(|s| s.inventory_id == *store_id)
                .ok_or_else(|| MaterialError::UnknownInventory((*store_id).to_owned()))?;
            store.remove(
                &input.resource_id,
                &input.grade_id,
                input
                    .quantity_per_run
                    .checked_mul(accepted)
                    .ok_or(MaterialError::QuantityOverflow)?,
            )?;
        }
        let mut total: u64 = 0;
        for output in &recipe.outputs {
            let q = output
                .quantity_per_run
                .checked_mul(accepted)
                .ok_or(MaterialError::QuantityOverflow)?
                .checked_mul(u64::from(output.yield_bps))
                .ok_or(MaterialError::QuantityOverflow)?
                / 10_000;
            self.inventories[oi].add(&output.resource_id, &output.grade_id, q)?;
            total = total
                .checked_add(q)
                .ok_or(MaterialError::QuantityOverflow)?;
        }
        let input_total = recipe
            .inputs
            .iter()
            .map(|i| i.quantity_per_run.saturating_mul(accepted))
            .sum();
        Ok(MaterialEvent::ProductionExecuted {
            run_id: run_id.to_owned(),
            recipe_id: recipe_id.to_owned(),
            accepted_runs: accepted,
            limit,
            input_quantity: input_total,
            output_quantity: total,
        })
    }
}
fn grade(
    m: &BTreeMap<&str, &ResourceDefinition>,
    resource: &str,
    grade_id: &str,
) -> Result<(), MaterialError> {
    let r = m
        .get(resource)
        .ok_or_else(|| MaterialError::UnknownResource(resource.to_owned()))?;
    if r.grade_ids.iter().any(|g| g == grade_id) {
        Ok(())
    } else {
        Err(MaterialError::GradeMismatch {
            resource_id: resource.to_owned(),
            grade_id: grade_id.to_owned(),
        })
    }
}
fn validate_id(id: &str) -> Result<(), MaterialError> {
    if is_valid_semantic_id(id) {
        Ok(())
    } else {
        Err(MaterialError::InvalidIdentifier(id.to_owned()))
    }
}

pub fn aggregate_material_fixture() -> MaterialState {
    MaterialState {
        resources: vec![
            ResourceDefinition {
                resource_id: "resource.raw_feed".to_owned(),
                category: ResourceCategory::RawMaterial,
                unit: QuantityUnit::BaseQuantity,
                grade_ids: vec!["grade.raw.standard".to_owned()],
                rules_version: "rules-0.1.0".to_owned(),
            },
            ResourceDefinition {
                resource_id: "resource.limestone".to_owned(),
                category: ResourceCategory::RawMaterial,
                unit: QuantityUnit::BaseQuantity,
                grade_ids: vec!["grade.limestone.standard".to_owned()],
                rules_version: "rules-0.1.0".to_owned(),
            },
            ResourceDefinition {
                resource_id: "resource.finished_aggregate".to_owned(),
                category: ResourceCategory::FinishedGood,
                unit: QuantityUnit::BaseQuantity,
                grade_ids: vec!["grade.aggregate.standard".to_owned()],
                rules_version: "rules-0.1.0".to_owned(),
            },
            ResourceDefinition {
                resource_id: "resource.aggregate_waste".to_owned(),
                category: ResourceCategory::Waste,
                unit: QuantityUnit::BaseQuantity,
                grade_ids: vec!["grade.waste.standard".to_owned()],
                rules_version: "rules-0.1.0".to_owned(),
            },
        ],
        recipes: vec![RecipeDefinition {
            recipe_id: "recipe.aggregate_crush".to_owned(),
            facility_type_ids: vec!["aggregate_processing_fixture".to_owned()],
            inputs: vec![
                RecipeInput {
                    resource_id: "resource.raw_feed".to_owned(),
                    grade_id: "grade.raw.standard".to_owned(),
                    quantity_per_run: 10,
                },
                RecipeInput {
                    resource_id: "resource.limestone".to_owned(),
                    grade_id: "grade.limestone.standard".to_owned(),
                    quantity_per_run: 2,
                },
            ],
            outputs: vec![
                RecipeOutput {
                    resource_id: "resource.finished_aggregate".to_owned(),
                    grade_id: "grade.aggregate.standard".to_owned(),
                    quantity_per_run: 8,
                    yield_bps: 10_000,
                },
                RecipeOutput {
                    resource_id: "resource.aggregate_waste".to_owned(),
                    grade_id: "grade.waste.standard".to_owned(),
                    quantity_per_run: 4,
                    yield_bps: 10_000,
                },
            ],
            rules_version: "rules-0.1.0".to_owned(),
        }],
        inventories: vec![
            InventoryStore {
                inventory_id: "inventory.aggregate_feed".to_owned(),
                facility_id: Some("facility.aggregate_plant_fixture".to_owned()),
                capacity: 1_000,
                permitted_resource_ids: vec![
                    "resource.raw_feed".to_owned(),
                    "resource.limestone".to_owned(),
                ],
                lots: vec![
                    MaterialLot {
                        resource_id: "resource.raw_feed".to_owned(),
                        grade_id: "grade.raw.standard".to_owned(),
                        quantity: 100,
                    },
                    MaterialLot {
                        resource_id: "resource.limestone".to_owned(),
                        grade_id: "grade.limestone.standard".to_owned(),
                        quantity: 20,
                    },
                ],
            },
            InventoryStore {
                inventory_id: "inventory.aggregate_finished".to_owned(),
                facility_id: Some("facility.aggregate_plant_fixture".to_owned()),
                capacity: 100,
                permitted_resource_ids: vec![
                    "resource.finished_aggregate".to_owned(),
                    "resource.aggregate_waste".to_owned(),
                ],
                lots: Vec::new(),
            },
        ],
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::{aggregate_fault_definitions, aggregate_plant_fixture, FailureCommand};

    fn facilities() -> Vec<Facility> {
        vec![aggregate_plant_fixture()]
    }
    fn failure() -> FailureState {
        FailureState::with_definitions(aggregate_fault_definitions())
    }
    fn production_command(runs: u64) -> MaterialCommand {
        MaterialCommand::ExecuteProduction {
            run_id: "run.aggregate".to_owned(),
            recipe_id: "recipe.aggregate_crush".to_owned(),
            facility_id: "facility.aggregate_plant_fixture".to_owned(),
            input_sources: vec![
                InputSource {
                    resource_id: "resource.raw_feed".to_owned(),
                    store_id: "inventory.aggregate_feed".to_owned(),
                },
                InputSource {
                    resource_id: "resource.limestone".to_owned(),
                    store_id: "inventory.aggregate_feed".to_owned(),
                },
            ],
            output_inventory_id: "inventory.aggregate_finished".to_owned(),
            requested_runs: runs,
        }
    }

    #[test]
    fn healthy_aggregate_run_consumes_inputs_and_produces_outputs() {
        let facilities = facilities();
        let mut materials = aggregate_material_fixture();
        let event = materials
            .execute(&facilities, &failure(), &production_command(5))
            .unwrap();
        assert_eq!(
            event,
            MaterialEvent::ProductionExecuted {
                run_id: "run.aggregate".to_owned(),
                recipe_id: "recipe.aggregate_crush".to_owned(),
                accepted_runs: 5,
                limit: MaterialLimit::None,
                input_quantity: 60,
                output_quantity: 60
            }
        );
        assert_eq!(
            materials.inventories[0].quantity("resource.raw_feed", "grade.raw.standard"),
            50
        );
        assert_eq!(
            materials.inventories[1]
                .quantity("resource.finished_aggregate", "grade.aggregate.standard"),
            40
        );
    }

    #[test]
    fn shortage_and_destination_capacity_are_bounded_without_negative_or_overfill() {
        let facilities = facilities();
        let mut shortage = aggregate_material_fixture();
        shortage.inventories[0].lots[0].quantity = 15;
        let event = shortage
            .execute(&facilities, &failure(), &production_command(5))
            .unwrap();
        assert!(matches!(
            event,
            MaterialEvent::ProductionExecuted {
                accepted_runs: 1,
                limit: MaterialLimit::InsufficientInput,
                ..
            }
        ));
        assert_eq!(
            shortage.inventories[0].quantity("resource.raw_feed", "grade.raw.standard"),
            5
        );
        let mut full = aggregate_material_fixture();
        full.inventories[1].lots.push(MaterialLot {
            resource_id: "resource.finished_aggregate".to_owned(),
            grade_id: "grade.aggregate.standard".to_owned(),
            quantity: 95,
        });
        let event = full
            .execute(&facilities, &failure(), &production_command(1))
            .unwrap();
        assert!(matches!(
            event,
            MaterialEvent::ProductionExecuted {
                accepted_runs: 0,
                limit: MaterialLimit::DestinationCapacity,
                ..
            }
        ));
        assert_eq!(
            full.inventories[0].quantity("resource.raw_feed", "grade.raw.standard"),
            100
        );
        assert_eq!(full.inventories[1].total_quantity(), 95);
    }

    #[test]
    fn grade_mismatch_and_transfer_conservation_are_enforced() {
        let mut materials = aggregate_material_fixture();
        let facilities = facilities();
        let before =
            materials.inventories[0].total_quantity() + materials.inventories[1].total_quantity();
        let transfer = MaterialCommand::Transfer {
            transfer_id: "transfer.raw".to_owned(),
            source_inventory_id: "inventory.aggregate_feed".to_owned(),
            destination_inventory_id: "inventory.aggregate_finished".to_owned(),
            resource_id: "resource.raw_feed".to_owned(),
            grade_id: "grade.wrong".to_owned(),
            quantity: 10,
        };
        assert!(matches!(
            materials.execute(&facilities, &failure(), &transfer),
            Err(MaterialError::GradeMismatch { .. })
        ));
        let valid = MaterialCommand::Transfer {
            transfer_id: "transfer.raw".to_owned(),
            source_inventory_id: "inventory.aggregate_feed".to_owned(),
            destination_inventory_id: "inventory.aggregate_finished".to_owned(),
            resource_id: "resource.raw_feed".to_owned(),
            grade_id: "grade.raw.standard".to_owned(),
            quantity: 10,
        };
        let event = materials.execute(&facilities, &failure(), &valid).unwrap();
        assert!(matches!(
            event,
            MaterialEvent::InventoryTransferred {
                accepted_quantity: 10,
                ..
            }
        ));
        let after =
            materials.inventories[0].total_quantity() + materials.inventories[1].total_quantity();
        assert_eq!(before, after);
    }

    #[test]
    fn failure_projected_capacity_limits_production_without_duplicate_math() {
        let facilities = facilities();
        let mut failure_state = failure();
        failure_state
            .apply(
                &facilities,
                &FailureCommand::ActivateFault {
                    fault_instance_id: "fault.instance.screen".to_owned(),
                    fault_type_id: "fault.screen_blockage".to_owned(),
                    component_id: "component.screen".to_owned(),
                    severity_bps: 5_000,
                },
                0,
                crate::RULES_VERSION,
            )
            .unwrap();
        let mut materials = aggregate_material_fixture();
        let event = materials
            .execute(&facilities, &failure_state, &production_command(10))
            .unwrap();
        assert!(matches!(
            event,
            MaterialEvent::ProductionExecuted {
                accepted_runs: 3,
                limit: MaterialLimit::FacilityCapacity,
                ..
            }
        ));
    }

    #[test]
    fn material_state_canonicalizes_insertion_order() {
        let first = aggregate_material_fixture();
        let mut second = first.clone();
        second.resources.reverse();
        second.recipes.reverse();
        second.inventories.reverse();
        assert_eq!(
            serde_json::to_string(&first.canonicalized()).unwrap(),
            serde_json::to_string(&second.canonicalized()).unwrap()
        );
    }
}
