## 1.1.6 (November 2, 2026)

### IMPROVEMENTS:

* Added `winrm_retry_interval` and `winrm_connect_timeout` communicator options to the `virtualbox-iso`, `virtualbox-ovf`, and `virtualbox-vm` builders (via `packer-plugin-sdk` update). [GH-195]

### DEPENDENCIES AND TOOLING:

* Updated `github.com/hashicorp/packer-plugin-sdk` from `v0.6.10` to `v0.6.13`. [GH-195] [GH-196] [GH-199]
* Updated Go version from `1.25.11` to `1.26.10`. [GH-199]
* Updated indirect dependencies, including `golang.org/x/crypto` to `v0.58.0` and `golang.org/x/net` to `v0.62.0`. [GH-196] [GH-199]
* Pinned GitHub Actions workflow refs to the latest commit SHAs. [GH-193]

## 1.1.5 (July 14, 2026)

### IMPROVEMENTS:

* Refreshed repository-wide copyright header coverage across builder, step, and version files.

### DEPENDENCIES AND TOOLING:

* Updated `github.com/hashicorp/packer-plugin-sdk` from `v0.6.9` to `v0.6.10`.

## 1.1.4 (June 16, 2026)

### IMPROVEMENTS:

* Added new VirtualBox hardware settings and configuration options.
* Added ARMv8 chipset support.
* Added missing floppy content handling when creating floppy media.
* Aligned storage controller names used by VBoxManage operations for better compatibility.

### BUG FIXES:

* Fixed shutdown message typo for ACPI.
* Updated plugin installation links in README.

### DEPENDENCIES AND TOOLING:

* Updated `github.com/hashicorp/packer-plugin-sdk` from `v0.6.1` to `v0.6.4`, and then to `v0.6.9`.
* Updated `github.com/ulikunitz/xz` to `v0.5.15`.
* Updated Go and linting/tooling configuration, including CI and release pipeline maintenance.
* Added backport workflow support.

## 1.1.3 (August 5, 2025)

### IMPROVEMENTS:

* Updated plugin release process: Plugin binaries are now published on the HashiCorp official [release site](https://releases.hashicorp.com/packer-plugin-virtualbox), ensuring a secure and standardized delivery pipeline.

### NOTES:
* **Binary Distribution Update**: To streamline our release process and align with other HashiCorp tools, all release binaries will now be published exclusively to the official HashiCorp [release](https://releases.hashicorp.com/packer-plugin-virtualbox) site. We will no longer attach release assets to GitHub Releases. Any scripts or automation that rely on the old location will need to be updated. For more information, see our post [here](https://discuss.hashicorp.com/t/important-update-official-packer-plugin-distribution-moving-to-releases-hashicorp-com/75972).

## 1.0.0 (June 14, 2021)

* Fix `Unknown option: --nested-hw-virt` bug [GH-26]
* Update packer-plugin-sdk to v0.2.3 [GH-27]

## 0.0.1 (April 16, 2021)

* VirtualBox Plugin break out from Packer core. Changes prior to break out can be found in [Packer's CHANGELOG](https://github.com/hashicorp/packer/blob/master/CHANGELOG.md).
