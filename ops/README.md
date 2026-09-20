# Native runtime repository templates

The systemd units are examples for later deployment of the four accepted native processes. They are deliberately not installed, enabled or started by WP-001. Secrets and deployment-specific paths belong in external environment/configuration files.

`ops/scripts/validate-systemd.sh` performs static checks without mutating the host.
