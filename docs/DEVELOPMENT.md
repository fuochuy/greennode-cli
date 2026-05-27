# GreenNode CLI — Development Guide

## Developer Workflow: Feature/Bug → Release

### Phase 1: Development

```bash
# 1. Create a feature branch
git checkout main && git pull
git checkout -b feat/add-new-command

# 2. Code
# Add command in go/cmd/vks/<command>.go
# Register in go/cmd/vks/vks.go

# 3. Build and test
cd go
CGO_ENABLED=0 go build -o grn .
./grn vks <new-command> --help
./grn vks <new-command> --dry-run ...

# 4. Update GitHub Pages docs
# - Create docs/commands/vks/<command-name>.md
# - Add entry to docs/commands/vks/index.md table
# - Add nav entry to mkdocs.yml

# 5. Add changelog fragment
cd ..
./scripts/new-change -t feature -c vks -d "Add new command"

# 6. Commit + push
git add .
git commit -m "feat(vks): add new command"
git push -u origin feat/add-new-command
```

### Phase 2: Pull Request

```
5. Create PR on GitHub (feat/add-new-command → main)
6. CI runs tests
7. Review + merge PR to main
```

### Phase 3: Release

```bash
# 8. Checkout main
git checkout main
git pull

# 9. Bump version
./scripts/bump-version minor
# Updates go/cmd/root.go version, merges changelog, commits, tags

# 10. Push
git push && git push --tags
# → GitHub Actions: build binaries → GitHub Release → upload artifacts
```

### Phase 4: Users Install

```bash
# Download binary from GitHub Releases
curl -L -o grn https://github.com/vngcloud/greennode-cli/releases/latest/download/grn-darwin-arm64
chmod +x grn
sudo mv grn /usr/local/bin/

# Or build from source
git clone https://github.com/vngcloud/greennode-cli.git
cd greennode-cli/go && go build -o grn .
```

---

## Building

```bash
cd go

# Build for current platform
CGO_ENABLED=0 go build -o grn .

# Cross-compile for all platforms
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o grn-linux-amd64 .
GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o grn-linux-arm64 .
GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -o grn-darwin-amd64 .
GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -o grn-darwin-arm64 .
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -o grn-windows-amd64.exe .
```

---

## Hotfix Flow

```bash
git checkout main
cd go
# Fix bug
CGO_ENABLED=0 go build -o grn .
./grn vks <command-to-test>
cd ..
./scripts/new-change -t bugfix -c auth -d "Fix token refresh"
git commit -am "fix(auth): fix token refresh"
./scripts/bump-version patch
git push && git push --tags
```

---

## Changelog Management

```bash
# Interactive
./scripts/new-change

# CLI args
./scripts/new-change -t feature -c vks -d "Add new command"
./scripts/new-change -t bugfix -c auth -d "Fix token refresh"
```

Change types: `feature`, `bugfix`, `enhancement`, `api-change`

## Version Bumping

```bash
./scripts/bump-version patch   # 0.1.0 → 0.1.1 (bug fixes)
./scripts/bump-version minor   # 0.1.0 → 0.2.0 (new features)
./scripts/bump-version major   # 0.1.0 → 1.0.0 (breaking changes)
```

## CI/CD Workflows

| Workflow | Trigger | Purpose |
|----------|---------|---------|
| `run-tests.yml` | PR to main/develop | Build + test Go binary |
| `release.yml` | Tag push `v*`, manual dispatch | Build multi-platform binaries + GitHub Release |
| `deploy-docs.yml` | Push to main (docs/) | Deploy documentation to GitHub Pages |
| `stale.yml` | Daily schedule | Auto-close stale issues |
