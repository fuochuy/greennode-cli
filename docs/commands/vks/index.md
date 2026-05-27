# VKS Commands

VKS (VNG Kubernetes Service) commands for managing Kubernetes clusters and node groups.

```bash
grn vks <command> [options]
```

## Available commands

### Cluster

| Command | Description |
|---------|-------------|
| [list-clusters](list-clusters.md) | List all VKS clusters |
| [get-cluster](get-cluster.md) | Get cluster details |
| [create-cluster](create-cluster.md) | Create a new VKS cluster |
| [update-cluster](update-cluster.md) | Update a VKS cluster |
| [delete-cluster](delete-cluster.md) | Delete a VKS cluster |

### Node Group

| Command | Description |
|---------|-------------|
| [list-nodegroups](list-nodegroups.md) | List node groups for a cluster |
| [get-nodegroup](get-nodegroup.md) | Get node group details |
| [create-nodegroup](create-nodegroup.md) | Create a new node group |
| [update-nodegroup](update-nodegroup.md) | Update a node group |
| [delete-nodegroup](delete-nodegroup.md) | Delete a node group |

### Auto-Upgrade

| Command | Description |
|---------|-------------|
| [set-auto-upgrade-config](set-auto-upgrade-config.md) | Configure auto-upgrade schedule for a cluster |
| [delete-auto-upgrade-config](delete-auto-upgrade-config.md) | Delete auto-upgrade config for a cluster |

### Quota

| Command | Description |
|---------|-------------|
| [get-quota](get-quota.md) | Get VKS quota limits and current usage |

### Waiter

| Command | Description |
|---------|-------------|
| [wait-cluster-active](wait-cluster-active.md) | Wait until cluster reaches ACTIVE status |
