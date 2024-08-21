Type: `virtualbox-ovf`
Artifact BuilderId: `mitchellh.virtualbox`

This VirtualBox Packer builder is able to create
[VirtualBox](https://www.virtualbox.org/) virtual machines and export them in
the OVF format, starting from an existing OVF/OVA (exported virtual machine
image).

When exporting from VirtualBox make sure to choose OVF Version 2, since Version
1 is not compatible and will generate errors like this:

```shell-session
==> virtualbox-ovf: Progress state: VBOX_E_FILE_ERROR
==> virtualbox-ovf: VBoxManage: error: Appliance read failed
==> virtualbox-ovf: VBoxManage: error: Error reading "source.ova": element "Section" has no "type" attribute, line 21
==> virtualbox-ovf: VBoxManage: error: Details: code VBOX_E_FILE_ERROR (0x80bb0004), component Appliance, interface IAppliance
==> virtualbox-ovf: VBoxManage: error: Context: "int handleImportAppliance(HandlerArg*)" at line 304 of file VBoxManageAppliance.cpp
```

The builder builds a virtual machine by importing an existing OVF or OVA file.
It then boots this image, runs provisioners on this new VM, and exports that VM
to create the image. The imported machine is deleted prior to finishing the
build.

## Basic Example

Here is a basic example. This example is functional if you have an OVF matching
the settings here.

**JSON**

```json
{
  "type": "virtualbox-ovf",
  "source_path": "source.ovf",
  "ssh_username": "packer",
  "ssh_password": "packer",
  "shutdown_command": "echo 'packer' | sudo -S shutdown -P now"
}
```

**HCL2**

```hcl
source "virtualbox-ovf" "basic-example" {
  source_path = "source.ovf"
  ssh_username = "packer"
  ssh_password = "packer"
  shutdown_command = "echo 'packer' | sudo -S shutdown -P now"
}

build {
  sources = ["sources.virtualbox-ovf.basic-example"]
}
```


It is important to add a `shutdown_command`. By default Packer halts the virtual
machine and the file system may not be sync'd. Thus, changes made in a
provisioner might not be saved.

## Configuration Reference

There are many configuration options available for the builder. In addition to
the items listed here, you will want to look at the general configuration
references for [ISO](#iso-configuration),
[HTTP](#http-directory-configuration),
[Floppy](#floppy-configuration),
[Export](#export-configuration),
[Boot](#boot-configuration),
[Shutdown](#shutdown-configuration),
[Run](#run-configuration),
[Communicator](#communicator-configuration)
configuration references, which are
necessary for this build to succeed and can be found further down the page.

### Required:

<!-- Code generated from the comments of the Config struct in builder/virtualbox/ovf/config.go; DO NOT EDIT MANUALLY -->

- `checksum` (string) - The checksum for the source_path file. The type of the checksum is
  specified within the checksum field as a prefix, ex: "md5:{$checksum}".
  The type of the checksum can also be omitted and Packer will try to
  infer it based on string length. Valid values are "none", "{$checksum}",
  "md5:{$checksum}", "sha1:{$checksum}", "sha256:{$checksum}",
  "sha512:{$checksum}" or "file:{$path}". Here is a list of valid checksum
  values:
   * md5:090992ba9fd140077b0661cb75f7ce13
   * 090992ba9fd140077b0661cb75f7ce13
   * sha1:ebfb681885ddf1234c18094a45bbeafd91467911
   * ebfb681885ddf1234c18094a45bbeafd91467911
   * sha256:ed363350696a726b7932db864dda019bd2017365c9e299627830f06954643f93
   * ed363350696a726b7932db864dda019bd2017365c9e299627830f06954643f93
   * file:http://releases.ubuntu.com/20.04/SHA256SUMS
   * file:file://./local/path/file.sum
   * file:./local/path/file.sum
   * none
  Although the checksum will not be verified when it is set to "none",
  this is not recommended since these files can be very large and
  corruption does happen from time to time.

- `source_path` (string) - The filepath or URL to an OVF or OVA file that acts as the
  source of this build.

<!-- End of code generated from the comments of the Config struct in builder/virtualbox/ovf/config.go; -->


#### Optional:

<!-- Code generated from the comments of the Config struct in builder/virtualbox/ovf/config.go; DO NOT EDIT MANUALLY -->

- `import_flags` ([]string) - Additional flags to pass to
  VBoxManage import. This can be used to add additional command-line flags
  such as --eula-accept to accept a EULA in the OVF.

- `import_opts` (string) - Additional options to pass to the
  VBoxManage import. This can be useful for passing keepallmacs or
  keepnatmacs options for existing ovf images.

- `target_path` (string) - The path where the OVA should be saved
  after download. By default, it will go in the packer cache, with a hash of
  the original filename as its name.

- `vm_name` (string) - This is the name of the OVF file for the new virtual
  machine, without the file extension. By default this is packer-BUILDNAME,
  where "BUILDNAME" is the name of the build.

- `keep_registered` (bool) - Set this to true if you would like to keep
  the VM registered with virtualbox. Defaults to false.

- `skip_export` (bool) - Defaults to false. When enabled, Packer will
  not export the VM. Useful if the build output is not the resultant image,
  but created inside the VM.

<!-- End of code generated from the comments of the Config struct in builder/virtualbox/ovf/config.go; -->


<!-- Code generated from the comments of the VBoxVersionConfig struct in builder/virtualbox/common/vbox_version_config.go; DO NOT EDIT MANUALLY -->

- `virtualbox_version_file` (\*string) - The path within the virtual machine to
  upload a file that contains the VirtualBox version that was used to create
  the machine. This information can be useful for provisioning. By default
  this is .vbox_version, which will generally be upload it into the
  home directory. Set to an empty string to skip uploading this file, which
  can be useful when using the none communicator.

<!-- End of code generated from the comments of the VBoxVersionConfig struct in builder/virtualbox/common/vbox_version_config.go; -->


<!-- Code generated from the comments of the GuestAdditionsConfig struct in builder/virtualbox/common/guest_additions_config.go; DO NOT EDIT MANUALLY -->

- `guest_additions_mode` (string) - The method by which guest additions are
  made available to the guest for installation. Valid options are `upload`,
  `attach`, or `disable`. If the mode is `attach` the guest additions ISO will
  be attached as a CD device to the virtual machine. If the mode is `upload`
  the guest additions ISO will be uploaded to the path specified by
  `guest_additions_path`. The default value is `upload`. If `disable` is used,
  guest additions won't be downloaded, either.

- `guest_additions_interface` (string) - The interface type to use to mount guest additions when
  guest_additions_mode is set to attach. Will default to the value set in
  iso_interface, if iso_interface is set. Will default to "ide", if
  iso_interface is not set. Options are "ide" and "sata".

- `guest_additions_path` (string) - The path on the guest virtual machine
   where the VirtualBox guest additions ISO will be uploaded. By default this
   is `VBoxGuestAdditions.iso` which should upload into the login directory of
   the user. This is a [configuration
   template](/packer/docs/templates/legacy_json_templates/engine) where the `Version`
   variable is replaced with the VirtualBox version.

- `guest_additions_sha256` (string) - The SHA256 checksum of the guest
   additions ISO that will be uploaded to the guest VM. By default the
   checksums will be downloaded from the VirtualBox website, so this only needs
   to be set if you want to be explicit about the checksum.

- `guest_additions_url` (string) - The URL of the guest additions ISO
   to upload. This can also be a file URL if the ISO is at a local path. By
   default, the VirtualBox builder will attempt to find the guest additions ISO
   on the local file system. If it is not available locally, the builder will
   download the proper guest additions ISO from the internet.

<!-- End of code generated from the comments of the GuestAdditionsConfig struct in builder/virtualbox/common/guest_additions_config.go; -->


### VBoxManage configuration

<!-- Code generated from the comments of the VBoxManageConfig struct in builder/virtualbox/common/vboxmanage_config.go; DO NOT EDIT MANUALLY -->

- `vboxmanage` ([][]string) - Custom `VBoxManage` commands to execute in order to further customize
  the virtual machine being created. The example shown below sets the memory and number of CPUs
  within the virtual machine:
  
  In JSON:
  ```json
  "vboxmanage": [
     ["modifyvm", "{{.Name}}", "--memory", "1024"],
     ["modifyvm", "{{.Name}}", "--cpus", "2"]
  ]
  ```
  
  In HCL2:
  ```hcl
  vboxmanage = [
     ["modifyvm", "{{.Name}}", "--memory", "1024"],
     ["modifyvm", "{{.Name}}", "--cpus", "2"],
  ]
  ```
  
  The value of `vboxmanage` is an array of commands to execute. These commands are
  executed in the order defined. So in the above example, the memory will be set
  followed by the CPUs.
  Each command itself is an array of strings, where each string is an argument to
  `VBoxManage`. Each argument is treated as a [configuration
  template](/packer/docs/templates/legacy_json_templates/engine). The only available
  variable is `Name` which is replaced with the unique name of the VM, which is
  required for many VBoxManage calls.

- `vboxmanage_post` ([][]string) - Identical to vboxmanage,
  except that it is run after the virtual machine is shutdown, and before the
  virtual machine is exported.

<!-- End of code generated from the comments of the VBoxManageConfig struct in builder/virtualbox/common/vboxmanage_config.go; -->


### Http directory configuration

<!-- Code generated from the comments of the HTTPConfig struct in multistep/commonsteps/http_config.go; DO NOT EDIT MANUALLY -->

Packer will create an http server serving `http_directory` when it is set, a
random free port will be selected and the architecture of the directory
referenced will be available in your builder.

Example usage from a builder:

```
wget http://{{ .HTTPIP }}:{{ .HTTPPort }}/foo/bar/preseed.cfg
```

<!-- End of code generated from the comments of the HTTPConfig struct in multistep/commonsteps/http_config.go; -->


#### Optional:

<!-- Code generated from the comments of the HTTPConfig struct in multistep/commonsteps/http_config.go; DO NOT EDIT MANUALLY -->

- `http_directory` (string) - Path to a directory to serve using an HTTP server. The files in this
  directory will be available over HTTP that will be requestable from the
  virtual machine. This is useful for hosting kickstart files and so on.
  By default this is an empty string, which means no HTTP server will be
  started. The address and port of the HTTP server will be available as
  variables in `boot_command`. This is covered in more detail below.

- `http_content` (map[string]string) - Key/Values to serve using an HTTP server. `http_content` works like and
  conflicts with `http_directory`. The keys represent the paths and the
  values contents, the keys must start with a slash, ex: `/path/to/file`.
  `http_content` is useful for hosting kickstart files and so on. By
  default this is empty, which means no HTTP server will be started. The
  address and port of the HTTP server will be available as variables in
  `boot_command`. This is covered in more detail below.
  Example:
  ```hcl
    http_content = {
      "/a/b"     = file("http/b")
      "/foo/bar" = templatefile("${path.root}/preseed.cfg", { packages = ["nginx"] })
    }
  ```

- `http_port_min` (int) - These are the minimum and maximum port to use for the HTTP server
  started to serve the `http_directory`. Because Packer often runs in
  parallel, Packer will choose a randomly available port in this range to
  run the HTTP server. If you want to force the HTTP server to be on one
  port, make this minimum and maximum port the same. By default the values
  are `8000` and `9000`, respectively.

- `http_port_max` (int) - HTTP Port Max

- `http_bind_address` (string) - This is the bind address for the HTTP server. Defaults to 0.0.0.0 so that
  it will work with any network interface.

<!-- End of code generated from the comments of the HTTPConfig struct in multistep/commonsteps/http_config.go; -->


### Floppy configuration

<!-- Code generated from the comments of the FloppyConfig struct in multistep/commonsteps/floppy_config.go; DO NOT EDIT MANUALLY -->

A floppy can be made available for your build. This is most useful for
unattended Windows installs, which look for an Autounattend.xml file on
removable media. By default, no floppy will be attached. All files listed in
this setting get placed into the root directory of the floppy and the floppy
is attached as the first floppy device. The summary size of the listed files
must not exceed 1.44 MB. The supported ways to move large files into the OS
are using `http_directory` or [the file
provisioner](/packer/docs/provisioner/file).

<!-- End of code generated from the comments of the FloppyConfig struct in multistep/commonsteps/floppy_config.go; -->


#### Optional:

<!-- Code generated from the comments of the FloppyConfig struct in multistep/commonsteps/floppy_config.go; DO NOT EDIT MANUALLY -->

- `floppy_files` ([]string) - A list of files to place onto a floppy disk that is attached when the VM
  is booted. Currently, no support exists for creating sub-directories on
  the floppy. Wildcard characters (\\*, ?, and \[\]) are allowed. Directory
  names are also allowed, which will add all the files found in the
  directory to the floppy.

- `floppy_dirs` ([]string) - A list of directories to place onto the floppy disk recursively. This is
  similar to the `floppy_files` option except that the directory structure
  is preserved. This is useful for when your floppy disk includes drivers
  or if you just want to organize it's contents as a hierarchy. Wildcard
  characters (\\*, ?, and \[\]) are allowed. The maximum summary size of
  all files in the listed directories are the same as in `floppy_files`.

- `floppy_content` (map[string]string) - Key/Values to add to the floppy disk. The keys represent the paths, and
  the values contents. It can be used alongside `floppy_files` or
  `floppy_dirs`, which is useful to add large files without loading them
  into memory. If any paths are specified by both, the contents in
  `floppy_content` will take precedence.
  
  Usage example (HCL):
  
  ```hcl
  floppy_files = ["vendor-data"]
  floppy_content = {
    "meta-data" = jsonencode(local.instance_data)
    "user-data" = templatefile("user-data", { packages = ["nginx"] })
  }
  floppy_label = "cidata"
  ```

- `floppy_label` (string) - Floppy Label

<!-- End of code generated from the comments of the FloppyConfig struct in multistep/commonsteps/floppy_config.go; -->


### CD configuration

<!-- Code generated from the comments of the CDConfig struct in multistep/commonsteps/extra_iso_config.go; DO NOT EDIT MANUALLY -->

An iso (CD) containing custom files can be made available for your build.

By default, no extra CD will be attached. All files listed in this setting
get placed into the root directory of the CD and the CD is attached as the
second CD device.

This config exists to work around modern operating systems that have no
way to mount floppy disks, which was our previous go-to for adding files at
boot time.

<!-- End of code generated from the comments of the CDConfig struct in multistep/commonsteps/extra_iso_config.go; -->


#### Optional:

<!-- Code generated from the comments of the CDConfig struct in multistep/commonsteps/extra_iso_config.go; DO NOT EDIT MANUALLY -->

- `cd_files` ([]string) - A list of files to place onto a CD that is attached when the VM is
  booted. This can include either files or directories; any directories
  will be copied onto the CD recursively, preserving directory structure
  hierarchy. Symlinks will have the link's target copied into the directory
  tree on the CD where the symlink was. File globbing is allowed.
  
  Usage example (JSON):
  
  ```json
  "cd_files": ["./somedirectory/meta-data", "./somedirectory/user-data"],
  "cd_label": "cidata",
  ```
  
  Usage example (HCL):
  
  ```hcl
  cd_files = ["./somedirectory/meta-data", "./somedirectory/user-data"]
  cd_label = "cidata"
  ```
  
  The above will create a CD with two files, user-data and meta-data in the
  CD root. This specific example is how you would create a CD that can be
  used for an Ubuntu 20.04 autoinstall.
  
  Since globbing is also supported,
  
  ```hcl
  cd_files = ["./somedirectory/*"]
  cd_label = "cidata"
  ```
  
  Would also be an acceptable way to define the above cd. The difference
  between providing the directory with or without the glob is whether the
  directory itself or its contents will be at the CD root.
  
  Use of this option assumes that you have a command line tool installed
  that can handle the iso creation. Packer will use one of the following
  tools:
  
    * xorriso
    * mkisofs
    * hdiutil (normally found in macOS)
    * oscdimg (normally found in Windows as part of the Windows ADK)

- `cd_content` (map[string]string) - Key/Values to add to the CD. The keys represent the paths, and the values
  contents. It can be used alongside `cd_files`, which is useful to add large
  files without loading them into memory. If any paths are specified by both,
  the contents in `cd_content` will take precedence.
  
  Usage example (HCL):
  
  ```hcl
  cd_files = ["vendor-data"]
  cd_content = {
    "meta-data" = jsonencode(local.instance_data)
    "user-data" = templatefile("user-data", { packages = ["nginx"] })
  }
  cd_label = "cidata"
  ```

- `cd_label` (string) - CD Label

<!-- End of code generated from the comments of the CDConfig struct in multistep/commonsteps/extra_iso_config.go; -->


### Export configuration

#### Optional:

<!-- Code generated from the comments of the ExportConfig struct in builder/virtualbox/common/export_config.go; DO NOT EDIT MANUALLY -->

- `format` (string) - Either ovf or ova, this specifies the output format
  of the exported virtual machine. This defaults to ovf.

- `export_opts` ([]string) - Additional options to pass to the [VBoxManage
  export](https://www.virtualbox.org/manual/ch09.html#vboxmanage-export).
  This can be useful for passing product information to include in the
  resulting appliance file. Packer JSON configuration file example:
  
  In JSON:
  ```json
  {
    "type": "virtualbox-iso",
    "export_opts":
    [
      "--manifest",
      "--vsys", "0",
      "--description", "{{user `vm_description`}}",
      "--version", "{{user `vm_version`}}"
    ],
    "format": "ova",
  }
  ```
  
  In HCL2:
  ```hcl
  	source "virtualbox-iso" "basic-example" {
  		export_opts = [
  	          "--manifest",
  	          "--vsys", "0",
  	          "--description", "${var.vm_description}",
  	          "--version", "${var.vm_version}"
  	   	]
  		format = "ova"
   }
  ```
  
  A VirtualBox [VM
  description](https://www.virtualbox.org/manual/ch09.html#vboxmanage-export-ovf)
  may contain arbitrary strings; the GUI interprets HTML formatting. However,
  the JSON format does not allow arbitrary newlines within a value. Add a
  multi-line description by preparing the string in the shell before the
  packer call like this (shell `>` continuation character snipped for easier
  copy & paste):
  
  ```shell
  vm_description='some
  multiline
  description'
  
  vm_version='0.2.0'
  
  packer build \
      -var "vm_description=${vm_description}" \
      -var "vm_version=${vm_version}"         \
      "packer_conf.json"
  ```

<!-- End of code generated from the comments of the ExportConfig struct in builder/virtualbox/common/export_config.go; -->


### Output configuration

#### Optional:

<!-- Code generated from the comments of the OutputConfig struct in builder/virtualbox/common/output_config.go; DO NOT EDIT MANUALLY -->

- `output_directory` (string) - This is the path to the directory where the
  resulting virtual machine will be created. This may be relative or absolute.
  If relative, the path is relative to the working directory when packer
  is executed. This directory must not exist or be empty prior to running
  the builder. By default this is output-BUILDNAME where "BUILDNAME" is the
  name of the build.

- `output_filename` (string) - This is the base name of the file (excluding the file extension) where
  the resulting virtual machine will be created. By default this is the
  `vm_name`.

<!-- End of code generated from the comments of the OutputConfig struct in builder/virtualbox/common/output_config.go; -->


### Run configuration

#### Optional:

<!-- Code generated from the comments of the RunConfig struct in builder/virtualbox/common/run_config.go; DO NOT EDIT MANUALLY -->

- `headless` (bool) - Packer defaults to building VirtualBox virtual
  machines by launching a GUI that shows the console of the machine
  being built. When this value is set to true, the machine will start
  without a console.

- `vrdp_bind_address` (string) - The IP address that should be
  binded to for VRDP. By default packer will use 127.0.0.1 for this. If you
  wish to bind to all interfaces use 0.0.0.0.

- `vrdp_port_min` (int) - The minimum and maximum port
  to use for VRDP access to the virtual machine. Packer uses a randomly chosen
  port in this range that appears available. By default this is 5900 to
  6000. The minimum and maximum ports are inclusive.

- `vrdp_port_max` (int) - VRDP Port Max

<!-- End of code generated from the comments of the RunConfig struct in builder/virtualbox/common/run_config.go; -->


### Shutdown configuration

#### Optional:

<!-- Code generated from the comments of the ShutdownConfig struct in builder/virtualbox/common/shutdown_config.go; DO NOT EDIT MANUALLY -->

- `shutdown_command` (string) - The command to use to gracefully shut down the
  machine once all the provisioning is done. By default this is an empty
  string, which tells Packer to just forcefully shut down the machine unless a
  shutdown command takes place inside script so this may safely be omitted. If
  one or more scripts require a reboot it is suggested to leave this blank
  since reboots may fail and specify the final shutdown command in your
  last script.

- `shutdown_timeout` (duration string | ex: "1h5m2s") - The amount of time to wait after executing the
  shutdown_command for the virtual machine to actually shut down. If it
  doesn't shut down in this time, it is an error. By default, the timeout is
  5m or five minutes.

- `post_shutdown_delay` (duration string | ex: "1h5m2s") - The amount of time to wait after shutting
  down the virtual machine. If you get the error
  Error removing floppy controller, you might need to set this to 5m
  or so. By default, the delay is 0s or disabled.

- `disable_shutdown` (bool) - Packer normally halts the virtual machine after all provisioners have
  run when no `shutdown_command` is defined.  If this is set to `true`, Packer
  *will not* halt the virtual machine but will assume that you will send the stop
  signal yourself through the preseed.cfg or your final provisioner.
  Packer will wait for a default of 5 minutes until the virtual machine is shutdown.
  The timeout can be changed using `shutdown_timeout` option.

- `acpi_shutdown` (bool) - If it's set to true, it will shutdown the VM via power button. It could be a good option
  when keeping the machine state is necessary after shutting it down.

<!-- End of code generated from the comments of the ShutdownConfig struct in builder/virtualbox/common/shutdown_config.go; -->


### Communicator configuration

#### Optional common fields:

<!-- Code generated from the comments of the Config struct in communicator/config.go; DO NOT EDIT MANUALLY -->

- `communicator` (string) - Packer currently supports three kinds of communicators:
  
  -   `none` - No communicator will be used. If this is set, most
      provisioners also can't be used.
  
  -   `ssh` - An SSH connection will be established to the machine. This
      is usually the default.
  
  -   `winrm` - A WinRM connection will be established.
  
  In addition to the above, some builders have custom communicators they
  can use. For example, the Docker builder has a "docker" communicator
  that uses `docker exec` and `docker cp` to execute scripts and copy
  files.

- `pause_before_connecting` (duration string | ex: "1h5m2s") - We recommend that you enable SSH or WinRM as the very last step in your
  guest's bootstrap script, but sometimes you may have a race condition
  where you need Packer to wait before attempting to connect to your
  guest.
  
  If you end up in this situation, you can use the template option
  `pause_before_connecting`. By default, there is no pause. For example if
  you set `pause_before_connecting` to `10m` Packer will check whether it
  can connect, as normal. But once a connection attempt is successful, it
  will disconnect and then wait 10 minutes before connecting to the guest
  and beginning provisioning.

<!-- End of code generated from the comments of the Config struct in communicator/config.go; -->


<!-- Code generated from the comments of the CommConfig struct in builder/virtualbox/common/comm_config.go; DO NOT EDIT MANUALLY -->

- `host_port_min` (int) - The minimum port to use for the Communicator port on the host machine which is forwarded
  to the SSH or WinRM port on the guest machine. By default this is 2222.

- `host_port_max` (int) - The maximum port to use for the Communicator port on the host machine which is forwarded
  to the SSH or WinRM port on the guest machine. Because Packer often runs in parallel,
  Packer will choose a randomly available port in this range to use as the
  host port. By default this is 4444.

- `skip_nat_mapping` (bool) - Defaults to false. When enabled, Packer
  does not setup forwarded port mapping for communicator (SSH or WinRM) requests and uses ssh_port or winrm_port
  on the host to communicate to the virtual machine.

- `ssh_listen_address` (string) - The address where the SSH port forwarding will be set to listen on. This value defaults to `127.0.0.1`.

<!-- End of code generated from the comments of the CommConfig struct in builder/virtualbox/common/comm_config.go; -->


#### Optional SSH fields:

<!-- Code generated from the comments of the SSH struct in communicator/config.go; DO NOT EDIT MANUALLY -->

- `ssh_host` (string) - The address to SSH to. This usually is automatically configured by the
  builder.

- `ssh_port` (int) - The port to connect to SSH. This defaults to `22`.

- `ssh_username` (string) - The username to connect to SSH with. Required if using SSH.

- `ssh_password` (string) - A plaintext password to use to authenticate with SSH.

- `ssh_ciphers` ([]string) - This overrides the value of ciphers supported by default by Golang.
  The default value is [
    "aes128-gcm@openssh.com",
    "chacha20-poly1305@openssh.com",
    "aes128-ctr", "aes192-ctr", "aes256-ctr",
  ]
  
  Valid options for ciphers include:
  "aes128-ctr", "aes192-ctr", "aes256-ctr", "aes128-gcm@openssh.com",
  "chacha20-poly1305@openssh.com",
  "arcfour256", "arcfour128", "arcfour", "aes128-cbc", "3des-cbc",

- `ssh_clear_authorized_keys` (bool) - If true, Packer will attempt to remove its temporary key from
  `~/.ssh/authorized_keys` and `/root/.ssh/authorized_keys`. This is a
  mostly cosmetic option, since Packer will delete the temporary private
  key from the host system regardless of whether this is set to true
  (unless the user has set the `-debug` flag). Defaults to "false";
  currently only works on guests with `sed` installed.

- `ssh_key_exchange_algorithms` ([]string) - If set, Packer will override the value of key exchange (kex) algorithms
  supported by default by Golang. Acceptable values include:
  "curve25519-sha256@libssh.org", "ecdh-sha2-nistp256",
  "ecdh-sha2-nistp384", "ecdh-sha2-nistp521",
  "diffie-hellman-group14-sha1", and "diffie-hellman-group1-sha1".

- `ssh_certificate_file` (string) - Path to user certificate used to authenticate with SSH.
  The `~` can be used in path and will be expanded to the
  home directory of current user.

- `ssh_pty` (bool) - If `true`, a PTY will be requested for the SSH connection. This defaults
  to `false`.

- `ssh_timeout` (duration string | ex: "1h5m2s") - The time to wait for SSH to become available. Packer uses this to
  determine when the machine has booted so this is usually quite long.
  Example value: `10m`.
  This defaults to `5m`, unless `ssh_handshake_attempts` is set.

- `ssh_disable_agent_forwarding` (bool) - If true, SSH agent forwarding will be disabled. Defaults to `false`.

- `ssh_handshake_attempts` (int) - The number of handshakes to attempt with SSH once it can connect.
  This defaults to `10`, unless a `ssh_timeout` is set.

- `ssh_bastion_host` (string) - A bastion host to use for the actual SSH connection.

- `ssh_bastion_port` (int) - The port of the bastion host. Defaults to `22`.

- `ssh_bastion_agent_auth` (bool) - If `true`, the local SSH agent will be used to authenticate with the
  bastion host. Defaults to `false`.

- `ssh_bastion_username` (string) - The username to connect to the bastion host.

- `ssh_bastion_password` (string) - The password to use to authenticate with the bastion host.

- `ssh_bastion_interactive` (bool) - If `true`, the keyboard-interactive used to authenticate with bastion host.

- `ssh_bastion_private_key_file` (string) - Path to a PEM encoded private key file to use to authenticate with the
  bastion host. The `~` can be used in path and will be expanded to the
  home directory of current user.

- `ssh_bastion_certificate_file` (string) - Path to user certificate used to authenticate with bastion host.
  The `~` can be used in path and will be expanded to the
  home directory of current user.

- `ssh_file_transfer_method` (string) - `scp` or `sftp` - How to transfer files, Secure copy (default) or SSH
  File Transfer Protocol.
  
  **NOTE**: Guests using Windows with Win32-OpenSSH v9.1.0.0p1-Beta, scp
  (the default protocol for copying data) returns a a non-zero error code since the MOTW
  cannot be set, which cause any file transfer to fail. As a workaround you can override the transfer protocol
  with SFTP instead `ssh_file_transfer_method = "sftp"`.

- `ssh_proxy_host` (string) - A SOCKS proxy host to use for SSH connection

- `ssh_proxy_port` (int) - A port of the SOCKS proxy. Defaults to `1080`.

- `ssh_proxy_username` (string) - The optional username to authenticate with the proxy server.

- `ssh_proxy_password` (string) - The optional password to use to authenticate with the proxy server.

- `ssh_keep_alive_interval` (duration string | ex: "1h5m2s") - How often to send "keep alive" messages to the server. Set to a negative
  value (`-1s`) to disable. Example value: `10s`. Defaults to `5s`.

- `ssh_read_write_timeout` (duration string | ex: "1h5m2s") - The amount of time to wait for a remote command to end. This might be
  useful if, for example, packer hangs on a connection after a reboot.
  Example: `5m`. Disabled by default.

- `ssh_remote_tunnels` ([]string) - 

- `ssh_local_tunnels` ([]string) - 

<!-- End of code generated from the comments of the SSH struct in communicator/config.go; -->


<!-- Code generated from the comments of the SSHTemporaryKeyPair struct in communicator/config.go; DO NOT EDIT MANUALLY -->

- `temporary_key_pair_type` (string) - `dsa` | `ecdsa` | `ed25519` | `rsa` ( the default )
  
  Specifies the type of key to create. The possible values are 'dsa',
  'ecdsa', 'ed25519', or 'rsa'.
  
  NOTE: DSA is deprecated and no longer recognized as secure, please
  consider other alternatives like RSA or ED25519.

- `temporary_key_pair_bits` (int) - Specifies the number of bits in the key to create. For RSA keys, the
  minimum size is 1024 bits and the default is 4096 bits. Generally, 3072
  bits is considered sufficient. DSA keys must be exactly 1024 bits as
  specified by FIPS 186-2. For ECDSA keys, bits determines the key length
  by selecting from one of three elliptic curve sizes: 256, 384 or 521
  bits. Attempting to use bit lengths other than these three values for
  ECDSA keys will fail. Ed25519 keys have a fixed length and bits will be
  ignored.
  
  NOTE: DSA is deprecated and no longer recognized as secure as specified
  by FIPS 186-5, please consider other alternatives like RSA or ED25519.

<!-- End of code generated from the comments of the SSHTemporaryKeyPair struct in communicator/config.go; -->


- `ssh_keypair_name` (string) - If specified, this is the key that will be used for SSH with the
  machine. The key must match a key pair name loaded up into the remote.
  By default, this is blank, and Packer will generate a temporary keypair
  unless [`ssh_password`](#ssh_password) is used.
  [`ssh_private_key_file`](#ssh_private_key_file) or
  [`ssh_agent_auth`](#ssh_agent_auth) must be specified when
  [`ssh_keypair_name`](#ssh_keypair_name) is utilized.


- `ssh_private_key_file` (string) - Path to a PEM encoded private key file to use to authenticate with SSH.
  The `~` can be used in path and will be expanded to the home directory
  of current user.


- `ssh_agent_auth` (bool) - If true, the local SSH agent will be used to authenticate connections to
  the source instance. No temporary keypair will be created, and the
  values of [`ssh_password`](#ssh_password) and
  [`ssh_private_key_file`](#ssh_private_key_file) will be ignored. The
  environment variable `SSH_AUTH_SOCK` must be set for this option to work
  properly.


#### Optional WinRM fields:

<!-- Code generated from the comments of the WinRM struct in communicator/config.go; DO NOT EDIT MANUALLY -->

- `winrm_username` (string) - The username to use to connect to WinRM.

- `winrm_password` (string) - The password to use to connect to WinRM.

- `winrm_host` (string) - The address for WinRM to connect to.
  
  NOTE: If using an Amazon EBS builder, you can specify the interface
  WinRM connects to via
  [`ssh_interface`](/packer/integrations/hashicorp/amazon/latest/components/builder/ebs#ssh_interface)

- `winrm_no_proxy` (bool) - Setting this to `true` adds the remote
  `host:port` to the `NO_PROXY` environment variable. This has the effect of
  bypassing any configured proxies when connecting to the remote host.
  Default to `false`.

- `winrm_port` (int) - The WinRM port to connect to. This defaults to `5985` for plain
  unencrypted connection and `5986` for SSL when `winrm_use_ssl` is set to
  true.

- `winrm_timeout` (duration string | ex: "1h5m2s") - The amount of time to wait for WinRM to become available. This defaults
  to `30m` since setting up a Windows machine generally takes a long time.

- `winrm_use_ssl` (bool) - If `true`, use HTTPS for WinRM.

- `winrm_insecure` (bool) - If `true`, do not check server certificate chain and host name.

- `winrm_use_ntlm` (bool) - If `true`, NTLMv2 authentication (with session security) will be used
  for WinRM, rather than default (basic authentication), removing the
  requirement for basic authentication to be enabled within the target
  guest. Further reading for remote connection authentication can be found
  [here](https://msdn.microsoft.com/en-us/library/aa384295(v=vs.85).aspx).

<!-- End of code generated from the comments of the WinRM struct in communicator/config.go; -->


### Boot Configuration

<!-- Code generated from the comments of the BootConfig struct in bootcommand/config.go; DO NOT EDIT MANUALLY -->

The boot configuration is very important: `boot_command` specifies the keys
to type when the virtual machine is first booted in order to start the OS
installer. This command is typed after boot_wait, which gives the virtual
machine some time to actually load.

The boot_command is an array of strings. The strings are all typed in
sequence. It is an array only to improve readability within the template.

There are a set of special keys available. If these are in your boot
command, they will be replaced by the proper key:

-   `<bs>` - Backspace

-   `<del>` - Delete

-   `<enter> <return>` - Simulates an actual "enter" or "return" keypress.

-   `<esc>` - Simulates pressing the escape key.

-   `<tab>` - Simulates pressing the tab key.

-   `<f1> - <f12>` - Simulates pressing a function key.

-   `<up> <down> <left> <right>` - Simulates pressing an arrow key.

-   `<spacebar>` - Simulates pressing the spacebar.

-   `<insert>` - Simulates pressing the insert key.

-   `<home> <end>` - Simulates pressing the home and end keys.

  - `<pageUp> <pageDown>` - Simulates pressing the page up and page down
    keys.

-   `<menu>` - Simulates pressing the Menu key.

-   `<leftAlt> <rightAlt>` - Simulates pressing the alt key.

-   `<leftCtrl> <rightCtrl>` - Simulates pressing the ctrl key.

-   `<leftShift> <rightShift>` - Simulates pressing the shift key.

-   `<leftSuper> <rightSuper>` - Simulates pressing the ⌘ or Windows key.

  - `<wait> <wait5> <wait10>` - Adds a 1, 5 or 10 second pause before
    sending any additional keys. This is useful if you have to generally
    wait for the UI to update before typing more.

  - `<waitXX>` - Add an arbitrary pause before sending any additional keys.
    The format of `XX` is a sequence of positive decimal numbers, each with
    optional fraction and a unit suffix, such as `300ms`, `1.5h` or `2h45m`.
    Valid time units are `ns`, `us` (or `µs`), `ms`, `s`, `m`, `h`. For
    example `<wait10m>` or `<wait1m20s>`.

  - `<XXXOn> <XXXOff>` - Any printable keyboard character, and of these
    "special" expressions, with the exception of the `<wait>` types, can
    also be toggled on or off. For example, to simulate ctrl+c, use
    `<leftCtrlOn>c<leftCtrlOff>`. Be sure to release them, otherwise they
    will be held down until the machine reboots. To hold the `c` key down,
    you would use `<cOn>`. Likewise, `<cOff>` to release.

  - `{{ .HTTPIP }} {{ .HTTPPort }}` - The IP and port, respectively of an
    HTTP server that is started serving the directory specified by the
    `http_directory` configuration parameter. If `http_directory` isn't
    specified, these will be blank!

-   `{{ .Name }}` - The name of the VM.

Example boot command. This is actually a working boot command used to start an
CentOS 6.4 installer:

In JSON:

```json
"boot_command": [

	   "<tab><wait>",
	   " ks=http://{{ .HTTPIP }}:{{ .HTTPPort }}/centos6-ks.cfg<enter>"
	]

```

In HCL2:

```hcl
boot_command = [

	   "<tab><wait>",
	   " ks=http://{{ .HTTPIP }}:{{ .HTTPPort }}/centos6-ks.cfg<enter>"
	]

```

The example shown below is a working boot command used to start an Ubuntu
12.04 installer:

In JSON:

```json
"boot_command": [

	"<esc><esc><enter><wait>",
	"/install/vmlinuz noapic ",
	"preseed/url=http://{{ .HTTPIP }}:{{ .HTTPPort }}/preseed.cfg ",
	"debian-installer=en_US auto locale=en_US kbd-chooser/method=us ",
	"hostname={{ .Name }} ",
	"fb=false debconf/frontend=noninteractive ",
	"keyboard-configuration/modelcode=SKIP keyboard-configuration/layout=USA ",
	"keyboard-configuration/variant=USA console-setup/ask_detect=false ",
	"initrd=/install/initrd.gz -- <enter>"

]
```

In HCL2:

```hcl
boot_command = [

	"<esc><esc><enter><wait>",
	"/install/vmlinuz noapic ",
	"preseed/url=http://{{ .HTTPIP }}:{{ .HTTPPort }}/preseed.cfg ",
	"debian-installer=en_US auto locale=en_US kbd-chooser/method=us ",
	"hostname={{ .Name }} ",
	"fb=false debconf/frontend=noninteractive ",
	"keyboard-configuration/modelcode=SKIP keyboard-configuration/layout=USA ",
	"keyboard-configuration/variant=USA console-setup/ask_detect=false ",
	"initrd=/install/initrd.gz -- <enter>"

]
```

For more examples of various boot commands, see the sample projects from our
[community templates page](https://packer.io/community-tools#templates).

<!-- End of code generated from the comments of the BootConfig struct in bootcommand/config.go; -->


Please note that for the Virtuabox builder, the IP address of the HTTP server
Packer launches for you to access files like the preseed file in the example
above (`{{ .HTTPIP }}`) is hardcoded to 10.0.2.2. If you change the network
of your VM you must guarantee that you can still access this HTTP server.

The boot command is sent to the VM through the `VBoxManage` utility in as few
invocations as possible. We send each character in groups of 25, with a default
delay of 100ms between groups. The delay alleviates issues with latency and CPU
contention. If you notice missing keys, you can tune this delay by specifying
"boot_keygroup_interval" in your Packer template, for example:

**JSON**

```json
{
  "builders": [
    {
      "type": "virtualbox-ovf",
      "boot_keygroup_interval": "500ms"
      ...
    }
  ]
}
```

**HCL2**

```hcl
source "virtualbox-ovf" "basic-example" {
  boot_keygroup_interval = "500ms"
  # ...
}
```


#### Optional:

<!-- Code generated from the comments of the BootConfig struct in bootcommand/config.go; DO NOT EDIT MANUALLY -->

- `boot_keygroup_interval` (duration string | ex: "1h5m2s") - Time to wait after sending a group of key pressses. The value of this
  should be a duration. Examples are `5s` and `1m30s` which will cause
  Packer to wait five seconds and one minute 30 seconds, respectively. If
  this isn't specified, a sensible default value is picked depending on
  the builder type.

- `boot_wait` (duration string | ex: "1h5m2s") - The time to wait after booting the initial virtual machine before typing
  the `boot_command`. The value of this should be a duration. Examples are
  `5s` and `1m30s` which will cause Packer to wait five seconds and one
  minute 30 seconds, respectively. If this isn't specified, the default is
  `10s` or 10 seconds. To set boot_wait to 0s, use a negative number, such
  as "-1s"

- `boot_command` ([]string) - This is an array of commands to type when the virtual machine is first
  booted. The goal of these commands should be to type just enough to
  initialize the operating system installer. Special keys can be typed as
  well, and are covered in the section below on the boot command. If this
  is not specified, it is assumed the installer will start itself.

<!-- End of code generated from the comments of the BootConfig struct in bootcommand/config.go; -->


### SSH key pair automation

The VirtualBox builders can inject the current SSH key pair's public key into
the template using the `SSHPublicKey` template engine. This is the SSH public
key as a line in OpenSSH authorized_keys format.

When a private key is provided using `ssh_private_key_file`, the key's
corresponding public key can be accessed using the above engine.

- `ssh_private_key_file` (string) - Path to a PEM encoded private key file to use to authenticate with SSH.
  The `~` can be used in path and will be expanded to the home directory
  of current user.


If `ssh_password` and `ssh_private_key_file` are not specified, Packer will
automatically generate en ephemeral key pair. The key pair's public key can
be accessed using the template engine.

For example, the public key can be provided in the boot command as a URL
encoded string by appending `| urlquery` to the variable:

In JSON:

```json
"boot_command": [
  "<up><wait><tab> text ks=http://{{ .HTTPIP }}:{{ .HTTPPort }}/ks.cfg PACKER_USER={{ user `username` }} PACKER_AUTHORIZED_KEY={{ .SSHPublicKey | urlquery }}<enter>"
]
```

In HCL2:

```hcl
boot_command = [
  "<up><wait><tab> text ks=http://{{ .HTTPIP }}:{{ .HTTPPort }}/ks.cfg PACKER_USER={{ user `username` }} PACKER_AUTHORIZED_KEY={{ .SSHPublicKey | urlquery }}<enter>"
]
```

A kickstart could then leverage those fields from the kernel command line by
decoding the URL-encoded public key:

```shell
%post

# Newly created users need the file/folder framework for SSH key authentication.
umask 0077
mkdir /etc/skel/.ssh
touch /etc/skel/.ssh/authorized_keys

# Loop over the command line. Set interesting variables.
for x in $(cat /proc/cmdline)
do
  case $x in
    PACKER_USER=*)
      PACKER_USER="${x#*=}"
      ;;
    PACKER_AUTHORIZED_KEY=*)
      # URL decode $encoded into $PACKER_AUTHORIZED_KEY
      encoded=$(echo "${x#*=}" | tr '+' ' ')
      printf -v PACKER_AUTHORIZED_KEY '%b' "${encoded//%/\\x}"
      ;;
  esac
done

# Create/configure packer user, if any.
if [ -n "$PACKER_USER" ]
then
  useradd $PACKER_USER
  echo "%$PACKER_USER ALL=(ALL) NOPASSWD: ALL" >> /etc/sudoers.d/$PACKER_USER
  [ -n "$PACKER_AUTHORIZED_KEY" ] && echo $PACKER_AUTHORIZED_KEY >> $(eval echo ~"$PACKER_USER")/.ssh/authorized_keys
fi

%end
```

## Guest Additions

Packer will automatically download the proper guest additions for the version of
VirtualBox that is running and upload those guest additions into the virtual
machine so that provisioners can easily install them.

Packer downloads the guest additions from the official VirtualBox website, and
verifies the file with the official checksums released by VirtualBox.

After the virtual machine is up and the operating system is installed, Packer
uploads the guest additions into the virtual machine. The path where they are
uploaded is controllable by `guest_additions_path`, and defaults to
"VBoxGuestAdditions.iso". Without an absolute path, it is uploaded to the home
directory of the SSH user.

## VBoxManage Commands

<!-- Code generated from the comments of the VBoxManageConfig struct in builder/virtualbox/common/vboxmanage_config.go; DO NOT EDIT MANUALLY -->

In order to perform extra customization of the virtual machine, a template can
define extra calls to `VBoxManage` to perform.
[VBoxManage](https://www.virtualbox.org/manual/ch09.html) is the command-line
interface to VirtualBox where you can completely control VirtualBox. It can be
used to do things such as set RAM, CPUs, etc.

<!-- End of code generated from the comments of the VBoxManageConfig struct in builder/virtualbox/common/vboxmanage_config.go; -->


<!-- Code generated from the comments of the VBoxManageConfig struct in builder/virtualbox/common/vboxmanage_config.go; DO NOT EDIT MANUALLY -->

- `vboxmanage` ([][]string) - Custom `VBoxManage` commands to execute in order to further customize
  the virtual machine being created. The example shown below sets the memory and number of CPUs
  within the virtual machine:
  
  In JSON:
  ```json
  "vboxmanage": [
     ["modifyvm", "{{.Name}}", "--memory", "1024"],
     ["modifyvm", "{{.Name}}", "--cpus", "2"]
  ]
  ```
  
  In HCL2:
  ```hcl
  vboxmanage = [
     ["modifyvm", "{{.Name}}", "--memory", "1024"],
     ["modifyvm", "{{.Name}}", "--cpus", "2"],
  ]
  ```
  
  The value of `vboxmanage` is an array of commands to execute. These commands are
  executed in the order defined. So in the above example, the memory will be set
  followed by the CPUs.
  Each command itself is an array of strings, where each string is an argument to
  `VBoxManage`. Each argument is treated as a [configuration
  template](/packer/docs/templates/legacy_json_templates/engine). The only available
  variable is `Name` which is replaced with the unique name of the VM, which is
  required for many VBoxManage calls.

- `vboxmanage_post` ([][]string) - Identical to vboxmanage,
  except that it is run after the virtual machine is shutdown, and before the
  virtual machine is exported.

<!-- End of code generated from the comments of the VBoxManageConfig struct in builder/virtualbox/common/vboxmanage_config.go; -->
