use crate::{
    aggregate_fault_definitions, aggregate_material_fixture, aggregate_plant_fixture,
    ComponentConditionState, FailureCommand, FailureState, KernelError, SimulationState,
};

pub const OPENING_AGGREGATE_SCENARIO_ID: &str = "scenario.opening.aggregate";

pub fn opening_aggregate_scenario(seed: u64) -> Result<SimulationState, KernelError> {
    let facility = aggregate_plant_fixture();
    let mut state = SimulationState::with_facilities(seed, vec![facility])?;
    state.failure = FailureState::with_definitions(aggregate_fault_definitions());
    state
        .failure
        .component_conditions
        .push(ComponentConditionState {
            component_id: "component.crusher".to_owned(),
            condition_bps: 10_000,
            derate_bps: 300,
            deliberate_shutdown: false,
            bypasses: Vec::new(),
        });
    let facilities = state.facilities.clone();
    state
        .failure
        .apply(
            &facilities,
            &FailureCommand::ActivateFault {
                fault_instance_id: "fault.instance.opening.feed_motor".to_owned(),
                fault_type_id: "fault.motor_bearing_seizure".to_owned(),
                component_id: "component.feed_conveyor".to_owned(),
                severity_bps: 10_000,
            },
            0,
            &state.rules_version,
        )
        .map_err(KernelError::Failure)?;
    state
        .failure
        .apply(
            &facilities,
            &FailureCommand::ActivateFault {
                fault_instance_id: "fault.instance.opening.output_belt".to_owned(),
                fault_type_id: "fault.worn_belt".to_owned(),
                component_id: "component.output_conveyor".to_owned(),
                severity_bps: 4_000,
            },
            0,
            &state.rules_version,
        )
        .map_err(KernelError::Failure)?;
    state.material = aggregate_material_fixture();
    state.validate()?;
    Ok(state)
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::{execute, Command, MaterialCommand, ProofPayload, SCHEMA_VERSION};

    #[test]
    fn opening_scenario_starts_blocked_and_recovery_is_partial() {
        let mut state = opening_aggregate_scenario(1234).unwrap();
        assert_eq!(state.schema_version, SCHEMA_VERSION);
        assert_eq!(
            state
                .evaluate_facility("facility.aggregate_plant_fixture")
                .unwrap()
                .effective_capacity,
            0
        );
        let inspect = Command {
            command_id: "opening.inspect.feed".to_owned(),
            command_type: "failure.inspect".to_owned(),
            schema_version: crate::SCHEMA_VERSION.to_owned(),
            rules_version: crate::RULES_VERSION.to_owned(),
            effective_time_ms: 0,
            idempotency_key: "opening.inspect.feed.v1".to_owned(),
            payload: ProofPayload::Failure(FailureCommand::Inspect {
                evidence_id: "evidence.opening.feed".to_owned(),
                component_id: "component.feed_conveyor".to_owned(),
                diagnostic_capability_bps: 10_000,
            }),
        };
        execute(&mut state, &inspect).unwrap();
        let repair_belt = Command {
            command_id: "opening.repair.belt".to_owned(),
            command_type: "failure.intervene".to_owned(),
            schema_version: crate::SCHEMA_VERSION.to_owned(),
            rules_version: crate::RULES_VERSION.to_owned(),
            effective_time_ms: 0,
            idempotency_key: "opening.repair.belt.v1".to_owned(),
            payload: ProofPayload::Failure(FailureCommand::Intervene {
                intervention_id: "intervention.opening.belt".to_owned(),
                component_id: "component.output_conveyor".to_owned(),
                kind: crate::InterventionKind::Repair,
                fault_instance_id: Some("fault.instance.opening.output_belt".to_owned()),
                derate_bps: None,
                bypass: None,
            }),
        };
        execute(&mut state, &repair_belt).unwrap();
        assert_eq!(
            state
                .evaluate_facility("facility.aggregate_plant_fixture")
                .unwrap()
                .effective_capacity,
            0
        );
        let repair_motor = Command {
            command_id: "opening.repair.motor".to_owned(),
            command_type: "failure.intervene".to_owned(),
            schema_version: crate::SCHEMA_VERSION.to_owned(),
            rules_version: crate::RULES_VERSION.to_owned(),
            effective_time_ms: 0,
            idempotency_key: "opening.repair.motor.v1".to_owned(),
            payload: ProofPayload::Failure(FailureCommand::Intervene {
                intervention_id: "intervention.opening.motor".to_owned(),
                component_id: "component.feed_conveyor".to_owned(),
                kind: crate::InterventionKind::Repair,
                fault_instance_id: Some("fault.instance.opening.feed_motor".to_owned()),
                derate_bps: None,
                bypass: None,
            }),
        };
        execute(&mut state, &repair_motor).unwrap();
        assert_eq!(
            state
                .evaluate_facility("facility.aggregate_plant_fixture")
                .unwrap()
                .effective_capacity,
            3
        );
        let production = Command {
            command_id: "opening.production.first".to_owned(),
            command_type: "material.production".to_owned(),
            schema_version: crate::SCHEMA_VERSION.to_owned(),
            rules_version: crate::RULES_VERSION.to_owned(),
            effective_time_ms: 0,
            idempotency_key: "opening.production.first.v1".to_owned(),
            payload: ProofPayload::Material(MaterialCommand::ExecuteProduction {
                run_id: "run.opening.first".to_owned(),
                recipe_id: "recipe.aggregate_crush".to_owned(),
                facility_id: "facility.aggregate_plant_fixture".to_owned(),
                input_sources: vec![
                    crate::InputSource {
                        resource_id: "resource.raw_feed".to_owned(),
                        store_id: "inventory.aggregate_feed".to_owned(),
                    },
                    crate::InputSource {
                        resource_id: "resource.limestone".to_owned(),
                        store_id: "inventory.aggregate_feed".to_owned(),
                    },
                ],
                output_inventory_id: "inventory.aggregate_finished".to_owned(),
                requested_runs: 100,
            }),
        };
        execute(&mut state, &production).unwrap();
        assert_eq!(
            state
                .material
                .inventories
                .iter()
                .find(|i| i.inventory_id == "inventory.aggregate_finished")
                .unwrap()
                .quantity("resource.finished_aggregate", "grade.aggregate.standard")
                .unwrap(),
            24
        );
    }
}
