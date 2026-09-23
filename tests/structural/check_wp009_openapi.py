import json
from pathlib import Path

root = Path(__file__).resolve().parents[2]
openapi = root / "packages/schemas/openapi/gridworks-api-v1.json"
data = json.loads(openapi.read_text())
if data.get("openapi") != "3.1.0": raise SystemExit("WP-009 OpenAPI must be 3.1.0")
required = {
    ("GET", "/healthz"), ("GET", "/readyz"), ("GET", "/version"),
    ("POST", "/api/v1/auth/guest"), ("DELETE", "/api/v1/auth/session"),
    ("GET", "/api/v1/me"), ("PATCH", "/api/v1/me/profile"),
    ("GET", "/api/v1/players/{player_id}"), ("POST", "/api/v1/companies"),
    ("GET", "/api/v1/companies/{company_id}"), ("POST", "/api/v1/company-groups"),
}
actual = {(method.upper(), path) for path, item in data.get("paths", {}).items() for method in item}
if actual != required: raise SystemExit(f"WP-009 OpenAPI route inventory mismatch: {sorted(actual ^ required)}")
if "/api/v1/ownership" in str(data) or "/api/v1/identity" in str(data): raise SystemExit("WP-009 OpenAPI exposes forbidden raw identity/ownership route")
print("WP-009 OpenAPI contract is valid.")
