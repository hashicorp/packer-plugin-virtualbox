# For full specification on the configuration of this file visit:
# https://github.com/hashicorp/integration-template#metadata-configuration
integration {
  name = "VirtualBox"
  description = "The VirtualBox plugin creates VirtualBox virtual machines and export them to an OVA or OVF format."
  identifier = "packer/hashicorp/virtualbox"
  component {
    type = "builder"
    name = "VirtualBox ISO"
    slug = "iso"
  }
  component {
    type = "builder"
    name = "VirtualBox OVF/OVA"
    slug = "ovf"
  }
  component {
    type = "builder"
    name = "VirtualBox VM"
    slug = "vm"
  }
}
