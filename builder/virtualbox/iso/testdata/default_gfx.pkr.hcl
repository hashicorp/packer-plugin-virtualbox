packer {
  required_plugins {
    virtualbox = {
      version = ">= 1.0.0"
      source  = "github.com/hashicorp/virtualbox"
    }
  }
}
source "virtualbox-iso" "ubuntu_2404" {
  # AMD64 image
  guest_os_type     = "Ubuntu_64"
  iso_url           = "https://releases.ubuntu.com/24.04/ubuntu-24.04.2-live-server-amd64.iso"
  iso_checksum      = "sha256:d6dab0c3a657988501b4bd76f1297c053df710e06e0c3aece60dead24f270b4d"
  # ARM64 image
  # guest_os_type     = "Ubuntu_arm64"
  # iso_url           = "https://cdimage.ubuntu.com/releases/24.04.2/release/ubuntu-24.04.2-live-server-arm64.iso"
  # iso_checksum      = "sha256:9fd122eedff09dc57d66e1c29acb8d7a207e2a877e762bdf30d2c913f95f03a4"
  iso_interface     = "virtio"
  ssh_username      = "packer"
  ssh_password      = "packer"
  shutdown_command  = "echo 'packer' | sudo -S shutdown -P now"
  vm_name           = "ubuntu_vm_test"
  memory            = 2048
  cpus              = 2
  # chipset           = "armv8virtual"
  hard_drive_interface = "virtio"
  # headless          = true
  http_directory    = "./testdata/http_ubuntu_2404"
  ssh_timeout       = "60m"
  boot_wait         = "20s"
  boot_command = [
    "c<wait><wait><wait>",
    "linux /casper/vmlinuz --- autoinstall ds=\"nocloud-net;seedfrom=http://{{.HTTPIP}}:{{.HTTPPort}}/\"<enter><wait>",
    "initrd /casper/initrd<enter><wait>",
    "boot<enter>"
  ]
  vboxmanage = [
        # Firmware 
        ["modifyvm", "{{.Name}}", "--firmware", "efi"],

        # Input devices
        ["modifyvm", "{{.Name}}", "--mouse", "ps2"],
        ["modifyvm", "{{.Name}}", "--keyboard", "ps2"],

        # Boot order
        ["modifyvm", "{{.Name}}", "--boot1", "disk"],
        ["modifyvm", "{{.Name}}", "--boot2", "dvd"],
        ["modifyvm", "{{.Name}}", "--boot3", "floppy"],
        ["modifyvm", "{{.Name}}", "--boot4", "none"],

        # Network
        ["modifyvm", "{{.Name}}", "--macaddress1", "080027F0F51D"],
        ["modifyvm", "{{.Name}}", "--nat-localhostreachable1", "on"],

        # Audio
        ["modifyvm", "{{.Name}}", "--audioin", "off"],
        ["modifyvm", "{{.Name}}", "--audioout", "off"],

        # Other settings
        ["modifyvm", "{{.Name}}", "--rtcuseutc", "on"],
        ["modifyvm", "{{.Name}}", "--usbxhci", "on"],
        ["modifyvm", "{{.Name}}", "--clipboard-mode", "disabled"]
      ]
}

build {
  sources = ["source.virtualbox-iso.ubuntu_2404"]
}

