# create-server

## Description

Create a new vServer instance. Use `--dry-run` to validate parameters without sending the request.

## Synopsis

```
grn vserver create-server
    --name <value>
    --flavor-id <value>
    --image-id <value>
    --network-id <value>
    --subnet-id <value>
    --root-disk-type-id <value>
    --zone-id <value>
    [--root-disk-size <value>]
    [--root-disk-encryption-type <value>]
    [--encryption-volume]
    [--data-disk-type-id <value>]
    [--data-disk-size <value>]
    [--data-disk-encryption-type <value>]
    [--data-disk-name <value>]
    [--attach-floating]
    [--security-group <value>]
    [--ssh-key-id <value>]
    [--user-name <value>]
    [--user-password <value>]
    [--expire-password]
    [--server-group-id <value>]
    [--host-group-id <value>]
    [--enable-backup]
    [--backup-instance-point-id <value>]
    [--snapshot-instance-point-id <value>]
    [--period <value>]
    [--is-poc]
    [--is-enable-auto-renew]
    [--os-licence]
    [--user-data <value>]
    [--user-data-base64-encoded]
    [--dry-run]
```

## Options

### Required

`--name`
: Server name. 2–65 characters, alphanumeric, hyphens, and underscores; must start and end with an alphanumeric character.

`--flavor-id`
: Flavor ID that determines the CPU/RAM size. Use `grn vserver list-flavors` to find available flavor IDs.

`--image-id`
: OS image ID.

`--network-id`
: VPC network ID.

`--subnet-id`
: Subnet ID within the network.

`--root-disk-type-id`
: Volume type ID for the root disk. Use `grn vserver list-volume-types` to find available IDs.

`--zone-id`
: Availability zone ID (e.g. `HCM03-1A`).

### Root disk

`--root-disk-size` (default: `20`)
: Root disk size in GiB. Minimum 20 GiB.

`--root-disk-encryption-type` (optional)
: Encryption type for the root disk.

`--encryption-volume` (default: `false`)
: Encrypt the root volume.

### Data disk (optional)

`--data-disk-type-id`
: Volume type ID for the data disk.

`--data-disk-size`
: Data disk size in GiB. If `0` (default), no data disk is created.

`--data-disk-encryption-type`
: Encryption type for the data disk.

`--data-disk-name`
: Name for the data disk.

### Network

`--attach-floating` (default: `false`)
: Attach a public floating IP to the server.

`--security-group`
: Comma-separated list of security group IDs to attach.

### Authentication

`--ssh-key-id`
: SSH key ID to inject. The key must already exist in the project.

`--user-name`
: OS login username to create.

`--user-password`
: Password for the OS user.

`--expire-password` (default: `true`)
: Force a password change on first login.

### Placement

`--server-group-id`
: Server group ID for anti-affinity / affinity placement.

`--host-group-id`
: Dedicated host group ID.

### Backup / restore

`--enable-backup` (default: `false`)
: Enable automatic backup for the server.

`--backup-instance-point-id`
: Restore from this backup instance point.

`--snapshot-instance-point-id`
: Restore from this snapshot instance point.

### Billing

`--period` (default: `1`)
: Billing period in months.

`--is-poc` (default: `false`)
: Mark the instance as a PoC (proof-of-concept).

`--is-enable-auto-renew` (default: `false`)
: Enable automatic renewal at the end of the billing period.

`--os-licence` (default: `false`)
: Include OS licence cost in billing.

### User data

`--user-data`
: User data script passed to cloud-init on first boot.

`--user-data-base64-encoded` (default: `false`)
: Set when the `--user-data` value is already base64-encoded.

### Other

`--dry-run`
: Validate all parameters without creating the server.

## Examples

Minimal creation (SSH key auth):

```bash
grn vserver create-server \
  --name my-server \
  --flavor-id flav-dca36eca-2018-4607-bde4-ac9d9cd31655 \
  --image-id img-5948b4a0-1a88-4254-bb36-cc0faa0c38e3 \
  --network-id net-83f3e9a9-71e6-438d-89e8-bde44722d987 \
  --subnet-id sub-0305c2cc-a89f-4cd3-bdd0-49cd312049a6 \
  --root-disk-type-id vtype-7a7a8610-34f5-11ee-be56-0242ac120002 \
  --zone-id HCM03-1A \
  --ssh-key-id key-abc12345
```

Full creation with floating IP, backup, and data disk:

```bash
grn vserver create-server \
  --name prod-server-01 \
  --flavor-id flav-dca36eca-2018-4607-bde4-ac9d9cd31655 \
  --image-id img-5948b4a0-1a88-4254-bb36-cc0faa0c38e3 \
  --network-id net-83f3e9a9-71e6-438d-89e8-bde44722d987 \
  --subnet-id sub-0305c2cc-a89f-4cd3-bdd0-49cd312049a6 \
  --root-disk-type-id vtype-7a7a8610-34f5-11ee-be56-0242ac120002 \
  --zone-id HCM03-1A \
  --root-disk-size 40 \
  --data-disk-type-id vtype-7a7a8610-34f5-11ee-be56-0242ac120002 \
  --data-disk-size 100 \
  --attach-floating \
  --security-group secg-52042c19-2706-44db-b38c-2310bb853357 \
  --user-name admin \
  --user-password "Aa123456@" \
  --enable-backup \
  --period 1
```

Validate without creating:

```bash
grn vserver create-server \
  --name my-server \
  --flavor-id flav-dca36eca-2018-4607-bde4-ac9d9cd31655 \
  --image-id img-5948b4a0-1a88-4254-bb36-cc0faa0c38e3 \
  --network-id net-83f3e9a9-71e6-438d-89e8-bde44722d987 \
  --subnet-id sub-0305c2cc-a89f-4cd3-bdd0-49cd312049a6 \
  --root-disk-type-id vtype-7a7a8610-34f5-11ee-be56-0242ac120002 \
  --zone-id HCM03-1A \
  --dry-run
```
