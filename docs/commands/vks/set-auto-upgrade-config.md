# set-auto-upgrade-config

## Description

Configure auto-upgrade schedule for a cluster. Sets the days and time when automatic Kubernetes version upgrades will be performed.

## Synopsis

```
grn vks set-auto-upgrade-config
    --cluster-id <value>
    --weekdays <value>
    --time <value>
```

## Options

`--cluster-id` (required)
: The ID of the cluster.

`--weekdays` (required)
: Days of the week to perform auto upgrade, comma-separated. Valid values: `Mon`, `Tue`, `Wed`, `Thu`, `Fri`, `Sat`, `Sun`. Example: `Mon,Wed,Fri`

`--time` (required)
: Time of day to perform auto upgrade in 24-hour format `HH:mm`. Example: `03:00`

## Examples

Set auto-upgrade to run on weekdays at 3 AM:

```bash
grn vks set-auto-upgrade-config \
  --cluster-id k8s-xxxxx \
  --weekdays Mon,Tue,Wed,Thu,Fri \
  --time 03:00
```

Set auto-upgrade to run on weekends at midnight:

```bash
grn vks set-auto-upgrade-config \
  --cluster-id k8s-xxxxx \
  --weekdays Sat,Sun \
  --time 00:00
```
