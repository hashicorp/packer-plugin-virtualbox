
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

### Components

The plugin comes with multiple builders able to create VirtualBox
machines, depending on the strategy you want to use to build the image. 
The following VirtualBox builders are supported:

#### Builders
- [virtualbox-iso](../.web-docs/components/builder/iso/README.md) - Starts from an ISO
  file, creates a brand new VirtualBox VM, installs an OS, provisions
  software within the OS, then exports that machine to create an image. This
  is best for people who want to start from scratch.

- [virtualbox-ovf](../.web-docs/components/builder/ovf/README.md) - This builder imports
  an existing OVF/OVA file, runs provisioners on top of that VM, and exports
  that machine to create an image. This is best if you have an existing
  VirtualBox VM export you want to use as the source. As an additional
  benefit, you can feed the artifact of this builder back into itself to
  iterate on a machine.

- [virtualbox-vm](../.web-docs/components/builder/vm/README.md) - This builder uses an
  existing VM to run defined provisioners on top of that VM, and optionally
  creates a snapshot to save the changes applied from the provisioners. In
  addition the builder is able to export that machine to create an image. The
  builder is able to attach to a defined snapshot as a starting point, which
  could be defined statically or dynamically via a variable.
