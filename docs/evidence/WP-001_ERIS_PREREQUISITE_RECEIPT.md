# WP-001 ERIS Prerequisite Receipt

Authoritative receipt: [GitHub issue #1 comment 5751184899](https://github.com/xOMAIKOx/gridworkx/issues/1#issuecomment-5751184899)

- Host: ERIS (`eris`), Ubuntu 24.04.5 LTS, x86_64.
- Workspace: `/srv/gridworx/`.
- Audit: 2026-09-20T16:36:06Z.
- Exact parent: `0a96f2f05e3dd4e9e284486315422452bce51110`.
- Installed native packages: Rust/Cargo/rustfmt/Clippy 1.75.0, CMake 3.28.3, SQLite 3.45.1 plus development headers, Gitleaks 8.16.0, ShellCheck 0.9.0.
- Installed official user-local tool: Godot 4.7.2.stable.official.ed1daf0bf at `/home/michael/.local/bin/godot4`, SHA-512 verified against the official release checksum.
- Existing compatible tools retained: Git 2.43.0, GCC/G++ 13.3.0, Make 4.3, pkg-config 1.8.1, Go 1.26.5, Node.js 20.20.2/npm 10.8.2 and PostgreSQL client 16.15.
- No database, NATS, object-storage, Valkey or GRIDWORKS application daemon was installed, configured or started.
- No systemd application unit, deployment, nginx, DNS, firewall, WireGuard or externally reachable port work was performed.
- No Docker, Podman, Compose, Kubernetes, OCI runtime or container-dependent tooling was installed or invoked.
