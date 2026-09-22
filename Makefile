.PHONY: validate structure schemas systemd-check format-check test build security godot-check go-check rust-check admin-check

validate: structure schemas wp006-migration systemd-check rust-check go-check admin-check

structure:
	python3 tests/structural/check_repository.py

schemas:
	npm run validate --workspace @gridworks/tools

wp006-migration:
	python3 tests/structural/check_wp006_migration.py

systemd-check:
	./ops/scripts/validate-systemd.sh

format-check:
	cargo fmt --all -- --check
	python3 tests/structural/check_go_format.py
	rustfmt --version

rust-check:
	cargo fmt --all -- --check
	cargo test --workspace
	cargo clippy --workspace --all-targets --all-features -- -D warnings

go-check:
	go test ./services/...
	go vet ./services/...
	python3 tests/structural/check_go_format.py

admin-check:
	npm run typecheck --workspace @gridworks/admin
	npm run build --workspace @gridworks/admin
	npm run test --workspace @gridworks/admin

godot-check:
	godot --headless --path apps/game --editor --quit --audio-driver Dummy

build: rust-check go-check admin-check

security:
	gitleaks detect --source . --config .gitleaks.toml --no-banner --redact
	python3 tests/security/scan_forbidden_runtime.py
	shellcheck ops/scripts/*.sh

test: rust-check go-check admin-check schemas
