# start-server

## Description

Start a stopped vServer instance. The server must be in `STOPPED` status for this operation to succeed.

## Synopsis

```
grn vserver start-server
    --server-id <value>
```

## Options

`--server-id` (required)
: ID of the server to start.

## Examples

Start a server:

```bash
grn vserver start-server --server-id srv-abc12345-6789-def0-1234-abcdef012345
```
