# reboot-server

## Description

Reboot a vServer instance. Supports two reboot types:

- **SOFT** (default) — sends a graceful shutdown signal and waits for the OS to restart cleanly.
- **HARD** — equivalent to a power cycle; use when the OS is unresponsive.

## Synopsis

```
grn vserver reboot-server
    --server-id <value>
    [--type SOFT|HARD]
```

## Options

`--server-id` (required)
: ID of the server to reboot.

`--type` (optional, default: `SOFT`)
: Reboot type. `SOFT` for a graceful restart, `HARD` for a forced power cycle.

## Examples

Graceful reboot (default):

```bash
grn vserver reboot-server --server-id srv-abc12345-6789-def0-1234-abcdef012345
```

Hard reboot (force):

```bash
grn vserver reboot-server \
  --server-id srv-abc12345-6789-def0-1234-abcdef012345 \
  --type HARD
```
