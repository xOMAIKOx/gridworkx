# Canonical schema catalogue

`packages/schemas` is the repository-owned source boundary for cross-language contracts. JSON is the human-readable external/API representation. Future generated Rust, Go, TypeScript and OpenAPI artifacts must derive from these definitions rather than creating parallel authoritative models.

Every externally meaningful payload is associated with schema, rules/content and release versions. Semantic IDs use namespaces and never contain translated prose. WP-002 owns simulation command/snapshot schemas, WP-003 adds the strict facility graph schema, and WP-004 adds strict failure-state/failure-command contracts, and WP-005 adds material-state/material-command contracts; Rust representations and these schemas must evolve together.
