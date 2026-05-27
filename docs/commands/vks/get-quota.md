# get-quota

## Description

Get VKS quota limits and current usage for the current user, including maximum number of clusters, node groups per cluster, nodes per node group, and current cluster count.

## Synopsis

```
grn vks get-quota
```

## Options

No command-specific options. See [Global Options](../../usage/global-options.md) for flags available on all commands.

## Output fields

| Field | Description |
|-------|-------------|
| `maxClusters` | Maximum number of clusters allowed |
| `maxNodeGroupsPerCluster` | Maximum number of node groups per cluster |
| `maxNodesPerNodeGroup` | Maximum number of nodes per node group |
| `numClusters` | Current number of clusters in use |

## Examples

Get quota:

```bash
grn vks get-quota
```

Output:

```json
{
    "maxClusters": 200,
    "maxNodeGroupsPerCluster": 20,
    "maxNodesPerNodeGroup": 10,
    "numClusters": 4
}
```

Get only the maximum cluster limit:

```bash
grn vks get-quota --query maxClusters
```

Check remaining cluster capacity:

```bash
grn vks get-quota --output json | jq '.maxClusters - .numClusters'
```
