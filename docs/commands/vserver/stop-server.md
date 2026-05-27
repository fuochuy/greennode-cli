# stop-server

## Description

Stop a running vServer instance. The server must be in `ACTIVE` status for this operation to succeed.

## Synopsis

```
grn vserver stop-server
    --server-id <value>
```

## Options

`--server-id` (required)
: ID of the server to stop.

## Examples

Stop a server:

```bash
grn vserver stop-server --server-id srv-abc12345-6789-def0-1234-abcdef012345
```
