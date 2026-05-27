# Changelog

## 1.3.1

### Enhancements
* **configure**: Auto-detect project_id during grn configure by calling vServer /v1/projects with the given credentials and region

## 1.3.0

### Features
* **configure**: Add project_id field to configure wizard, configure list/get/set, and GRN_DEFAULT_PROJECT_ID env var override

## 1.2.0

### Features
* **vks**: Add get-quota command to retrieve VKS quota limits and current usage
* **vks**: Add get-quota command reference page to GitHub Pages docs

## 1.1.0

### Features
* **core**: Add GRN_ACCESS_KEY_ID/GRN_SECRET_ACCESS_KEY env var credential support

## 1.0.0

### Features
* **core**: Rewrite CLI in Go - single binary, zero dependencies, cross-platform

## 0.2.1

### Bug Fixes
* **security**: Remove GRN_CLIENT_ID/SECRET env var support for credential security

## 0.2.0

### Features
* **vks**: Add set-auto-upgrade-config and delete-auto-upgrade-config commands

## 0.1.0

### Features
* **core**: Initial release of Greenode CLI with VKS service support
