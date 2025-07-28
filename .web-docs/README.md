
The VirtualBox plugin is able to create
[VirtualBox](https://www.virtualbox.org) virtual machines and export them in
the OVA or OVF format.

### Installation
To install this plugin add this code into your Packer configuration and run [packer init](/packer/docs/commands/init)

```hcl
packer {
    required_plugins {
        virtualbox = {
          version = "~> 1"
          source  = "github.com/hashicorp/virtualbox"
        }
    }
}
```

Alternatively, you can use `packer plugins install` to manage installation of this plugin.

```sh
$ packer plugins install github.com/hashicorp/virtualbox
```

**Note: Update to Packer Plugin Installation**

With the new Packer release starting from version 1.14.0, the `packer init` command will automatically install official plugins from the [HashiCorp release site.](https://releases.hashicorp.com/)

Going forward, to use newer versions of official Packer plugins, you'll need to upgrade to Packer version 1.14.0 or later. If you're using an older version, you can still install plugins, but as a workaround, you'll need to [manually install them using the CLI.](https://developer.hashicorp.com/packer/docs/plugins/install#manually-install-plugins-using-the-cli)

There is no change to the syntax or commands for installing plugins.

### Components

The plugin comes with multiple builders able to create VirtualBox
machines, depending on the strategy you want to use to build the image. 
The following VirtualBox builders are supported:

#### Builders
- [virtualbox-iso](/packer/integrations/hashicorp/virtualbox/latest/components/builder/iso) - Starts from an ISO
  file, creates a brand new VirtualBox VM, installs an OS, provisions
  software within the OS, then exports that machine to create an image. This
  is best for people who want to start from scratch.

- [virtualbox-ovf](/packer/integrations/hashicorp/virtualbox/latest/components/builder/ovf) - This builder imports
  an existing OVF/OVA file, runs provisioners on top of that VM, and exports
  that machine to create an image. This is best if you have an existing
  VirtualBox VM export you want to use as the source. As an additional
  benefit, you can feed the artifact of this builder back into itself to
  iterate on a machine.

- [virtualbox-vm](/packer/integrations/hashicorp/virtualbox/latest/components/builder/vm) - This builder uses an
  existing VM to run defined provisioners on top of that VM, and optionally
  creates a snapshot to save the changes applied from the provisioners. In
  addition the builder is able to export that machine to create an image. The
  builder is able to attach to a defined snapshot as a starting point, which
  could be defined statically or dynamically via a variable.
