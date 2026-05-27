# delete-server

## Description

Delete a vServer instance. Before executing, the command fetches and displays a preview of the server that will be removed.

Use `--deleteAllVolumes` to also delete all volumes attached to the server. By default, volumes are retained.

**This action is irreversible.**

## Synopsis

```
grn vserver delete-server
    --server-id <value>
    [--deleteAllVolumes]
```

## Options

`--server-id` (required)
: ID of the server to delete.

`--deleteAllVolumes` (optional, default: `false`)
: Delete all volumes associated with the server. If omitted, attached volumes are preserved.

## Examples

Delete a server (retain volumes):

```bash
grn vserver delete-server --server-id srv-abc12345-6789-def0-1234-abcdef012345
```

Delete a server and all its volumes:

```bash
grn vserver delete-server \
  --server-id srv-abc12345-6789-def0-1234-abcdef012345 \
  --deleteAllVolumes
```
