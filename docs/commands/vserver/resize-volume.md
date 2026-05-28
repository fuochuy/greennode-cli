# resize-volume

Resize a volume's storage size, change its volume type, or both in a single operation.

## Synopsis

```bash
greennode vserver volume resize --volume-id <volume-id> [--size <gib>] [--volume-type-id <id>] [flags]
```

At least one of `--size` or `--volume-type-id` must be provided.

## Flags

| Flag | Required | Description |
|------|----------|-------------|
| `--volume-id` | Yes | Volume ID to resize |
| `--size` | No* | New size in GiB (must be ≥ current size) |
| `--volume-type-id` | No* | New volume type ID |
| `--dry-run` | No | Validate parameters without sending the request |

\* At least one of `--size` or `--volume-type-id` is required.

## Examples

Expand a volume to 100 GiB:

```bash
greennode vserver volume resize --volume-id vol-abc123 --size 100
```

Change the volume type:

```bash
greennode vserver volume resize --volume-id vol-abc123 --volume-type-id vtype-xyz789
```

Resize and change the type simultaneously:

```bash
greennode vserver volume resize --volume-id vol-abc123 --size 200 --volume-type-id vtype-xyz789
```

Dry-run to validate before applying:

```bash
greennode vserver volume resize --volume-id vol-abc123 --size 100 --dry-run
```

## Notes

- Volume size can only be increased, not decreased.
- To see available volume types for a zone, run: `greennode vserver volume-type list --zone-id <zone-id> --type <type>`
- The volume status may change to `RESIZING` during the operation and return to `AVAILABLE` when done.
