from pathlib import Path

root = Path(__file__).resolve().parents[2]
unit_dir = root / "ops" / "systemd"
units = sorted(unit_dir.glob("gridworks-*.service"))
required = [
    "[Unit]",
    "[Service]",
    "[Install]",
    "User=gridworks",
    "Group=gridworks",
    "EnvironmentFile=-/etc/gridworks/gridworks.env",
    "WorkingDirectory=/srv/gridworx",
    "Restart=on-failure",
    "NoNewPrivileges=true",
    "ProtectSystem=strict",
    "ProtectHome=true",
    "RestrictAddressFamilies=AF_UNIX AF_INET AF_INET6",
]
if len(units) != 4:
    raise SystemExit(f"expected four systemd templates, found {len(units)}")
for unit in units:
    text = unit.read_text()
    missing = [line for line in required if line not in text]
    if missing:
        raise SystemExit(f"{unit}: missing {', '.join(missing)}")
    if "User=root" in text or "--bind 0.0.0.0" in text:
        raise SystemExit(f"{unit}: unsafe identity or bind")
    if "ExecStart=" not in text or "--bind $GRIDWORKS_" not in text:
        raise SystemExit(f"{unit}: external bind contract missing")
print(f"Native systemd templates are statically valid: {len(units)} checked.")
