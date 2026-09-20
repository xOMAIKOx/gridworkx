from pathlib import Path
import re

root = Path(__file__).resolve().parents[2]
forbidden_names = {"Dockerfile", "docker-compose.yml", "docker-compose.yaml", "compose.yml", "compose.yaml"}
forbidden_content = re.compile(r"docker|podman|docker-compose|podman-compose|containerd|nerdctl|kubectl|oci://", re.IGNORECASE)
scan_roots = [".github", "apps", "crates", "services", "packages", "db", "ops", "tests"]
violations = []
for scan_root in scan_roots:
    for path in (root / scan_root).rglob("*"):
        relative = path.relative_to(root)
        if not path.is_file() or ".git" in path.parts or "node_modules" in path.parts or "target" in path.parts:
            continue
        if relative.as_posix() in {"tests/security/scan_forbidden_runtime.py", "tests/structural/check_repository.py"}:
            continue
        if path.name in forbidden_names:
            violations.append(f"manifest: {path.relative_to(root)}")
            continue
        try:
            text = path.read_text()
        except UnicodeDecodeError:
            continue
        if forbidden_content.search(text):
            violations.append(f"content: {path.relative_to(root)}")
if violations:
    print("forbidden runtime references detected:")
    print("\n".join(violations))
    raise SystemExit(1)
print("No forbidden container runtime manifests or dependencies found in implementation paths.")
