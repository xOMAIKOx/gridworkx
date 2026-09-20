#!/usr/bin/env bash
set -euo pipefail

exec python3 tests/structural/check_systemd_templates.py
