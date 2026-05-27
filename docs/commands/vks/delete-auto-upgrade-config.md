# delete-auto-upgrade-config

## Description

Delete the auto-upgrade configuration for a cluster. This disables automatic Kubernetes version upgrades.

## Synopsis

```
grn vks delete-auto-upgrade-config
    --cluster-id <value>
    [--force]
```

## Options

`--cluster-id` (required)
: The ID of the cluster.

`--force` (optional)
: Skip the confirmation prompt.

## Examples

Delete auto-upgrade config with confirmation:

```bash
grn vks delete-auto-upgrade-config --cluster-id k8s-xxxxx
# Are you sure you want to delete the auto-upgrade config? (yes/no): yes
```

Delete without confirmation (for scripting):

```bash
grn vks delete-auto-upgrade-config --cluster-id k8s-xxxxx --force
```
