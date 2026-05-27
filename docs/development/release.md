# Release Process

## Adding changelog entries

```bash
./scripts/new-change                          # Interactive
./scripts/new-change -t feature -c vks -d "Add new command"  # CLI args
```

Change types: `feature`, `bugfix`, `enhancement`, `api-change`

## Creating a release

```bash
./scripts/bump-version patch   # 0.1.0 → 0.1.1
./scripts/bump-version minor   # 0.1.0 → 0.2.0
./scripts/bump-version major   # 0.1.0 → 1.0.0
git push && git push --tags    # Triggers GitHub Actions release
```

## Release flow

```
Developer workflow:
1. During development: ./scripts/new-change (add fragments per PR)
2. Ready to release:   ./scripts/bump-version minor
3. Push:               git push && git push --tags

GitHub Actions (automatic):
4. Build Go binaries for Linux/macOS/Windows (amd64 + arm64)
5. Create GitHub Release with binaries attached
```

## CI/CD workflows

| Workflow | Trigger | Purpose |
|----------|---------|---------|
| `run-tests.yml` | PR to main/develop | Build + test Go binary |
| `release.yml` | Tag push `v*`, manual dispatch | Build multi-platform binaries + GitHub Release |
| `deploy-docs.yml` | Push to main (docs/) | Deploy documentation to GitHub Pages |
| `stale.yml` | Daily schedule | Auto-close stale issues |
