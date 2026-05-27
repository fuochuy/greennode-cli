# vServer CLI Setup Guide

Step-by-step guide to install the GreenNode CLI, configure credentials, and manage vServer instances — including optional Claude Desktop AI integration via MCP.

---

## Step 1 — Install the binary

=== "macOS (Apple Silicon)"

    ```bash
    curl -L -o grn https://github.com/vngcloud/greennode-cli/releases/latest/download/grn-darwin-arm64
    chmod +x grn
    sudo mv grn /usr/local/bin/
    ```

=== "macOS (Intel)"

    ```bash
    curl -L -o grn https://github.com/vngcloud/greennode-cli/releases/latest/download/grn-darwin-amd64
    chmod +x grn
    sudo mv grn /usr/local/bin/
    ```

=== "Linux (x86_64)"

    ```bash
    curl -L -o grn https://github.com/vngcloud/greennode-cli/releases/latest/download/grn-linux-amd64
    chmod +x grn
    sudo mv grn /usr/local/bin/
    ```

=== "Linux (ARM64)"

    ```bash
    curl -L -o grn https://github.com/vngcloud/greennode-cli/releases/latest/download/grn-linux-arm64
    chmod +x grn
    sudo mv grn /usr/local/bin/
    ```

=== "Build from source"

    Requires [Go 1.22+](https://go.dev/dl/).

    ```bash
    git clone https://github.com/vngcloud/greennode-cli.git
    cd greennode-cli/go
    go build -o grn .
    sudo mv grn /usr/local/bin/
    ```

Verify the installation:

```bash
grn --version
# grn-cli/1.3.1 Go/1.22.x darwin/arm64
```

---

## Step 2 — Get your credentials

1. Go to [VNG Cloud IAM Portal](https://hcm-3.console.vngcloud.vn/iam/)
2. Navigate to **Service Accounts** → create a new Service Account
3. Copy the **Client ID** and **Client Secret**

---

## Step 3 — Configure the CLI

Run the interactive setup wizard:

```bash
grn configure
```

You will be prompted for each value:

```
GRN Client ID [None]: <paste your Client ID>
GRN Client Secret [None]: <paste your Client Secret>
Default region name [HCM-3]:                        ← press Enter to keep HCM-3
Default output format [json]:                        ← press Enter to keep json
Project ID (leave blank to auto-detect) [None]:      ← press Enter to auto-detect
Fetching project_id from HCM-3...
Auto-detected project_id: pro-xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
```

Verify what was saved:

```bash
grn configure list
```

```
          Name                   Value            Type    Location
          ----                   -----            ----    --------
       profile               <not set>            None    None
     client_id    ****************bc6e     config-file    ~/.greenode/credentials
 client_secret    ****************c123     config-file    ~/.greenode/credentials
        region                   HCM-3     config-file    ~/.greenode/config
        output                    json     config-file    ~/.greenode/config
    project_id       pro-xxxxxxxx          config-file    ~/.greenode/config
```

---

## Step 4 — Explore available resources

Before creating a server, use these discovery commands to find valid IDs for each parameter:

```bash
# 1. List availability zones
grn vserver volume-type list

# 2. List OS images (use imageVersion to filter, e.g. "ubuntu-22")
grn vserver image list --type os
grn vserver image list --type os --image-version ubuntu-22

# 3. List GPU images
grn vserver image list --type gpu

# 4. List instance families and CPU platforms
grn vserver flavor list-families
grn vserver flavor list-codes

# 5. List flavors for a given family and CPU platform
grn vserver flavor list --family general-purpose --code code-g

# 6. List volume types for a zone
grn vserver volume-type list --zone-id <zone-id> --type SSD

# 7. List your VPCs and subnets
grn vserver vpc list
grn vserver subnet list --vpc-id <vpc-id>
```

---

## Step 5 — Create a server

```bash
grn vserver server create \
  --name my-server \
  --zone-id <zone-id> \
  --network-id <vpc-id> \
  --subnet-id <subnet-id> \
  --image-id <image-id> \
  --flavor-id <flavor-id> \
  --root-disk-type-id <volume-type-id> \
  --root-disk-size 40
```

!!! tip "Missing a value?"
    Run each command without the flag to see available options. For example, running `grn vserver server create` without `--zone-id` will list valid zone IDs.

---

## Step 6 — Manage servers

```bash
# List all servers
grn vserver server list

# Get server details
grn vserver server get --server-id <server-id>

# Start / stop / reboot
grn vserver server start  --server-id <server-id>
grn vserver server stop   --server-id <server-id>
grn vserver server reboot --server-id <server-id>

# Resize to a different flavor
grn vserver server resize --server-id <server-id> --flavor-id <new-flavor-id>

# Delete (shows preview and asks for confirmation)
grn vserver server delete --server-id <server-id>

# Delete without confirmation prompt
grn vserver server delete --server-id <server-id> --force
```

---

## Step 7 — Manage volumes, VPCs, and security groups

```bash
# Volumes
grn vserver volume list
grn vserver volume create --name my-vol --zone-id <zone-id> --volume-type-id <type-id> --size 100
grn vserver volume delete --volume-id <volume-id>

# VPCs and subnets
grn vserver vpc create --name my-vpc --cidr 10.0.0.0/16
grn vserver subnet create --vpc-id <vpc-id> --cidr 10.0.1.0/24 --zone-id <zone-id>

# Security groups and rules
grn vserver secgroup list
grn vserver secgroup create --name my-sg --description "Web servers"
grn vserver secgroup rule create \
  --secgroup-id <sg-id> \
  --direction ingress \
  --protocol tcp \
  --port-range-min 80 \
  --port-range-max 80 \
  --remote-ip-prefix 0.0.0.0/0 \
  --ether-type IPv4
```

---

## Step 8 — Claude Desktop MCP integration (optional)

This lets Claude Desktop manage your vServer resources in natural language — no commands needed.

### 8.1 Locate the Claude Desktop config file

=== "macOS"

    ```
    ~/Library/Application Support/Claude/claude_desktop_config.json
    ```

=== "Windows"

    ```
    %APPDATA%\Claude\claude_desktop_config.json
    ```

### 8.2 Add the MCP server

Open the config file and add the `mcpServers` block (keep any existing keys):

```json
{
  "mcpServers": {
    "greennode": {
      "command": "/usr/local/bin/grn",
      "args": ["mcp"]
    }
  }
}
```

!!! note
    If you installed `grn` to a custom path, replace `/usr/local/bin/grn` with the full path from `which grn`.

### 8.3 Restart Claude Desktop

Quit and reopen Claude Desktop. In a new conversation you can now say:

- *"List my vServer instances"*
- *"Create a server named test-01 in zone HCM03-1A using Ubuntu 22"*
- *"What flavors are available in the general-purpose family?"*
- *"Stop server srv-xxxxxxxx"*
- *"Show me all security groups"*

Claude will call the appropriate `grn` commands automatically and return the results.

---

## Quick reference

| Task | Command |
|------|---------|
| Configure credentials | `grn configure` |
| Show current config | `grn configure list` |
| List zones | `grn vserver volume-type list` |
| List images | `grn vserver image list --type os` |
| List flavors | `grn vserver flavor list --family <f> --code <c>` |
| List servers | `grn vserver server list` |
| Create server | `grn vserver server create --name ...` |
| Resize server | `grn vserver server resize --server-id ... --flavor-id ...` |
| Delete server | `grn vserver server delete --server-id ...` |
| List volumes | `grn vserver volume list` |
| List VPCs | `grn vserver vpc list` |
| List security groups | `grn vserver secgroup list` |
| Start MCP server | `grn mcp` |
