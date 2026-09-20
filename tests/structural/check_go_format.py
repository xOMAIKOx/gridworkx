from pathlib import Path
import subprocess

root = Path(__file__).resolve().parents[2]
files = sorted((root / "services").rglob("*.go"))
if not files:
    raise SystemExit("no Go files found")
result = subprocess.run(["gofmt", "-l", *map(str, files)], cwd=root, text=True, capture_output=True, check=False)
if result.returncode != 0:
    raise SystemExit(result.stderr)
if result.stdout.strip():
    print("Go files requiring gofmt:")
    print(result.stdout)
    raise SystemExit(1)
print(f"gofmt clean: {len(files)} files checked.")
