# Canonical schema catalogue

`packages/schemas` is the repository-owned source boundary for cross-language contracts. JSON is the human-readable external/API representation. Future generated Rust, Go, TypeScript and OpenAPI artifacts must derive from these definitions rather than creating parallel authoritative models.

Every externally meaningful payload is associated with schema, rules/content and release versions. Semantic IDs use namespaces and never contain translated prose. WP-002 adds the simulation command and snapshot schemas; the Rust representation and these schemas must evolve together.
