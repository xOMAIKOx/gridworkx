import fs from "node:fs";
import path from "node:path";
import process from "node:process";
import Ajv from "ajv/dist/2020.js";

const root = path.resolve(new URL("../../..", import.meta.url).pathname);
const schemaRoot = path.join(root, "packages", "schemas");
const contentRoot = path.join(root, "packages", "content");
const semanticId = /^[a-z][a-z0-9]*(\.[a-z0-9][a-z0-9_-]*)+$/;

const readJson = (filePath) => JSON.parse(fs.readFileSync(filePath, "utf8"));
const ajv = new Ajv({ strict: true });
for (const fileName of ["semantic-id.schema.json", "versioned-envelope.schema.json", "domain-catalog.schema.json", "content-manifest.schema.json"]) {
  ajv.addSchema(readJson(path.join(schemaRoot, fileName)));
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
validate("https://gridworks.example/schema/domain-catalog.schema.json", domainCatalog, "domain catalog");
validate("https://gridworks.example/schema/content-manifest.schema.json", contentManifest, "content manifest");
assertUnique(domainCatalog.domains.map((domain) => domain.domain_id), "domain catalog");
assertUnique(contentManifest.bundles.map((bundle) => bundle.bundle_id), "content manifest");

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
