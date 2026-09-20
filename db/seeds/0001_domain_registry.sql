INSERT INTO gridworks.domain_contracts (domain_id, owner_boundary, contract_state, schema_version)
VALUES
    ('identity.account', 'services/api', 'foundation', 'schema-0.1.0'),
    ('player.profile', 'services/api', 'foundation', 'schema-0.1.0'),
    ('company.ownership', 'services/api', 'foundation', 'schema-0.1.0'),
    ('world.region', 'packages/content', 'foundation', 'schema-0.1.0'),
    ('world.site', 'packages/content', 'foundation', 'schema-0.1.0'),
    ('facility.component', 'crates/gridworks-sim', 'foundation', 'schema-0.1.0'),
    ('resource.inventory', 'services/api', 'foundation', 'schema-0.1.0'),
    ('fault.catalog', 'packages/content', 'foundation', 'schema-0.1.0'),
    ('industry.catalog', 'packages/content', 'foundation', 'schema-0.1.0'),
    ('event.envelope', 'packages/schemas', 'foundation', 'schema-0.1.0'),
    ('transport.modes', 'packages/content', 'foundation', 'schema-0.1.0'),
    ('contract.work_exchange', 'services/api', 'foundation', 'schema-0.1.0'),
    ('property.real_estate', 'packages/content', 'foundation', 'schema-0.1.0'),
    ('region.catalog', 'packages/content', 'foundation', 'schema-0.1.0'),
    ('climate.profile', 'packages/content', 'foundation', 'schema-0.1.0'),
    ('construction.maintenance', 'services/api', 'foundation', 'schema-0.1.0'),
    ('skill.progression', 'services/api', 'foundation', 'schema-0.1.0'),
    ('manager.operations', 'services/api', 'foundation', 'schema-0.1.0'),
    ('trade.marketplace', 'services/api', 'foundation', 'schema-0.1.0'),
    ('consortium.jv', 'services/api', 'foundation', 'schema-0.1.0'),
    ('accounting.ledger', 'services/api', 'foundation', 'schema-0.1.0'),
    ('business.lifecycle', 'services/api', 'foundation', 'schema-0.1.0'),
    ('social.messaging', 'services/realtime', 'foundation', 'schema-0.1.0'),
    ('notification.preferences', 'services/worker', 'foundation', 'schema-0.1.0'),
    ('localization.catalogue', 'packages/localization', 'foundation', 'schema-0.1.0'),
    ('ownership.gridworks', 'packages/content', 'foundation', 'schema-0.1.0'),
    ('treasury.authority', 'services/api', 'foundation', 'schema-0.1.0'),
    ('onboarding.starter_path', 'apps/game', 'foundation', 'schema-0.1.0'),
    ('business.operating_state', 'services/api', 'foundation', 'schema-0.1.0'),
    ('world.presentation', 'apps/game', 'foundation', 'schema-0.1.0'),
    ('time.configuration', 'packages/content', 'foundation', 'schema-0.1.0'),
    ('bootstrap.economy', 'packages/content', 'foundation', 'schema-0.1.0'),
    ('moderation.policy', 'services/realtime', 'foundation', 'schema-0.1.0')
ON CONFLICT (domain_id) DO NOTHING;

INSERT INTO gridworks.system_principals (principal_id, display_label, principal_kind)
VALUES ('ownership.gridworks', 'GRIDWORKS', 'non_player_system')
ON CONFLICT (principal_id) DO NOTHING;
