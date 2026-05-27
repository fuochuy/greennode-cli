# list-flavors

## Description

List available vServer flavors. Results can be filtered by availability zone and instance family.

## Synopsis

```
grn vserver list-flavors
    [--zone <value>]
    [--instance-family <value>]
    [--page <value>]
    [--page-size <value>]
    [--no-paginate]
```

## Options

`--zone` (optional)
: Filter by availability zone (e.g. `HCM-3a`, `HCM-3b`).

`--instance-family` (optional)
: Filter by instance family (e.g. `General-Purpose`, `Compute-Optimized`).

`--page` (optional, default: `1`)
: Page number to retrieve.

`--page-size` (optional, default: `50`)
: Number of flavors per page.

`--no-paginate` (optional)
: Return a single page instead of paginating automatically.

## Examples

List all flavors:

```bash
grn vserver list-flavors
```

Filter by zone:

```bash
grn vserver list-flavors --zone HCM-3a
```

Filter by instance family:

```bash
grn vserver list-flavors --instance-family General-Purpose
```

Filter by both zone and instance family:

```bash
grn vserver list-flavors --zone HCM-3a --instance-family Compute-Optimized
```
