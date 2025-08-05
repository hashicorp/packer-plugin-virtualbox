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
