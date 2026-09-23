import fs from "node:fs";
import path from "node:path";
import Ajv from "ajv/dist/2020.js";

const root = path.resolve(new URL("../../..", import.meta.url).pathname);
const schemaRoot = path.join(root, "packages", "schemas");
const contentRoot = path.join(root, "packages", "content");
const semanticId = /^[a-z][a-z0-9]*(\.[a-z0-9][a-z0-9_-]*)+$/;

const readJson = (filePath) => JSON.parse(fs.readFileSync(filePath, "utf8"));
const ajv = new Ajv({ strict: true });
for (const fileName of [
  "semantic-id.schema.json",
  "versioned-envelope.schema.json",
  "domain-catalog.schema.json",
  "content-manifest.schema.json",
  "simulation-command.schema.json",
  "facility-graph.schema.json",
  "failure-state.schema.json",
  "failure-command.schema.json",
  "fault-definitions.schema.json",
  "material-state.schema.json",
  "material-command.schema.json",
  "ledger-state.schema.json",
  "ledger-command.schema.json",
  "identity-state.schema.json",
  "identity-command.schema.json",
  "identity-public-profile.schema.json",
  "company-state.schema.json",
  "company-command.schema.json",
  "company-public-profile.schema.json",
  "godot-bridge.schema.json",
  "simulation-state.schema.json",
]) {
  ajv.addSchema(readJson(path.join(schemaRoot, fileName)));
}

const schemaCatalog = readJson(path.join(schemaRoot, "catalog.json"));
for (const schema of schemaCatalog.schemas) {
  if (!fs.existsSync(path.join(schemaRoot, schema.path))) {
    throw new Error(`schema catalog: missing ${schema.path}`);
  }
}

const validate = (schemaId, document, label) => {
  const check = ajv.getSchema(schemaId);
  if (!check || !check(document)) {
    const details = check?.errors?.map((error) => `${error.instancePath || "/"} ${error.message}`).join("; ") ?? "schema unavailable";
    throw new Error(`${label}: ${details}`);
  }
};

const assertSemantic = (value, label) => {
  if (typeof value !== "string" || !semanticId.test(value)) {
    throw new Error(`${label}: invalid semantic identifier ${String(value)}`);
  }
};

const assertUnique = (values, label) => {
  const unique = new Set(values);
  if (unique.size !== values.length) {
    throw new Error(`${label}: duplicate identifiers`);
  }
};

const domainCatalog = readJson(path.join(contentRoot, "domain-catalog.json"));
const contentManifest = readJson(path.join(contentRoot, "content-manifest.json"));
const faultDefinitions = readJson(path.join(contentRoot, "config", "faults.json"));
const materialDefinitions = readJson(path.join(contentRoot, "config", "materials.json"));
validate("https://gridworks.example/schema/domain-catalog.schema.json", domainCatalog, "domain catalog");
validate("https://gridworks.example/schema/fault-definitions.schema.json", faultDefinitions, "fault definitions");
validate("https://gridworks.example/schema/material-state.schema.json", { resources: materialDefinitions.resources, recipes: materialDefinitions.recipes, inventories: [] }, "material definitions");
validate("https://gridworks.example/schema/content-manifest.schema.json", contentManifest, "content manifest");
assertUnique(domainCatalog.domains.map((domain) => domain.domain_id), "domain catalog");
assertUnique(contentManifest.bundles.map((bundle) => bundle.bundle_id), "content manifest");

const validFailureCommand = {
  failure: {
    action: "ActivateFault",
    data: {
      fault_instance_id: "fault.instance.fixture",
      fault_type_id: "fault.motor_bearing_seizure",
      component_id: "component.feed_conveyor",
      severity_bps: 10000,
    },
  },
};
validate("https://gridworks.example/schema/failure-command.schema.json", validFailureCommand, "failure command fixture");
const malformedFailureCommand = {
  failure: {action: "ActivateFault", data: {fault_instance_id: "fault.instance.fixture"}},
};
const failureCommandCheck = ajv.getSchema("https://gridworks.example/schema/failure-command.schema.json");
if (!failureCommandCheck || failureCommandCheck(malformedFailureCommand)) {
  throw new Error("failure command schema accepted a malformed ActivateFault payload");
}
const validLedgerCommand = {
  ledger: { Post: { transaction_id: "transaction.synthetic", transaction_type: "ledger.synthetic_grant", idempotency_key: "ledger.synthetic", source_ref: "event.synthetic", currency: "CRD", lines: [
    { line_id: "line.debit", sequence: 1, account_id: "account.system", side: "Debit", amount_minor: 1000 },
    { line_id: "line.credit", sequence: 2, account_id: "account.wallet", side: "Credit", amount_minor: 1000 },
  ] } },
};
validate("https://gridworks.example/schema/ledger-command.schema.json", validLedgerCommand, "ledger command fixture");
const validIdentityCommand = { mutation_type: "identity.guest_issue", idempotency_key: "identity.request.fixture" };
validate("https://gridworks.example/schema/identity-command.schema.json", validIdentityCommand, "identity command fixture");
validate("https://gridworks.example/schema/identity-public-profile.schema.json", { player_id: "player.fixture", handle: "fixture", display_name: "Fixture", locale: "en", visibility: "public", discoverable: true }, "public profile fixture");
validate("https://gridworks.example/schema/company-command.schema.json", { mutation_type: "company.create", idempotency_key: "company.fixture", company_type: "operating", name: "Fixture Co", owner: { type: "system", id: "principal.gridworks.system" } }, "company command fixture");
validate("https://gridworks.example/schema/company-command.schema.json", { mutation_type: "ownership.transfer", idempotency_key: "company.transfer.fixture", entity_type: "group", entity_id: "group.fixture", from: { type: "player", id: "player.fixture" }, to: { type: "system", id: "principal.gridworks.system" }, share_bps: 10000 }, "group ownership command fixture");
const malformedCompanyCommand = { mutation_type: "company.create", idempotency_key: "company.bad", company_type: "operating", name: "Fixture Co", owner: { type: "system", id: "principal.gridworks.system" }, company_id: "company.client-supplied" };
const companyCommandCheck = ajv.getSchema("https://gridworks.example/schema/company-command.schema.json");
if (!companyCommandCheck || companyCommandCheck(malformedCompanyCommand)) throw new Error("company command schema accepted caller-supplied authoritative company ID");
validate("https://gridworks.example/schema/company-public-profile.schema.json", { company_id: "company.fixture", name: "Fixture Co", visibility: "public" }, "company public profile fixture");
validate("https://gridworks.example/schema/godot-bridge.schema.json", { ok: true, operation: "bridge.metadata", result: {} }, "Godot bridge fixture");

const configFiles = fs.readdirSync(path.join(contentRoot, "config")).filter((fileName) => fileName.endsWith(".json"));
const semanticKeys = new Set(["config_id", "policy_id", "policy_version", "manifest_id", "principal_id", "system_principal", "layer_id", "mode_id", "asset_id", "catalogue_id"]);
for (const fileName of configFiles) {
  const document = readJson(path.join(contentRoot, "config", fileName));
  const walk = (value, location) => {
    if (!value || typeof value !== "object") return;
    for (const [key, child] of Object.entries(value)) {
      if (semanticKeys.has(key) && key !== "policy_version" && typeof child === "string") assertSemantic(child, `${fileName}${location}.${key}`);
      if (Array.isArray(child)) child.forEach((item, index) => walk(item, `${location}.${key}[${index}]`));
      else walk(child, `${location}.${key}`);
    }
  };
  walk(document, "");
}

const requiredNamespaces = ["resource.", "fault.", "transport.", "industry.", "event.", "contract.", "property.", "region.", "climate."];
for (const namespace of requiredNamespaces) {
  const present = domainCatalog.domains.some((domain) => domain.domain_id.startsWith(namespace));
  if (!present) throw new Error(`domain catalog: missing namespace ${namespace}`);
}

console.log(`Validated ${domainCatalog.domains.length} domain boundaries, ${contentManifest.bundles.length} content bundles and ${configFiles.length} configuration contracts.`);
