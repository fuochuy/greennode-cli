# Contributing

See [CONTRIBUTING.md](https://github.com/vngcloud/greennode-cli/blob/main/CONTRIBUTING.md) for the full contributing guide.

## Quick start

```bash
git clone https://github.com/vngcloud/greennode-cli.git
cd greennode-cli/go
go build -o grn .
./grn --version
```

## Adding a new command

1. Create `cmd/vks/<command_name>.go`
2. Define `cobra.Command` with Use, Short, RunE
3. Register in `cmd/vks/vks.go`
4. Add `validator.ValidateID()` for ID args
5. Add `--dry-run` for create/update/delete
6. Create `docs/commands/vks/<command-name>.md` — command reference page
7. Add entry to `docs/commands/vks/index.md` table
8. Add nav entry to `mkdocs.yml`

## Adding a new service

1. Create `cmd/<service>/` directory
2. Create parent command with `cobra.Command`
3. Register in `cmd/root.go`
