# libvirt-aws

Partial AWS API emulation powered by libvirt.

## Installation

Currently the only supported (and tested) OS is ubuntu 24.04

### Prerequisites

You'll need to install and configure `libvirt` and `sqlite`. The start command can accept custom inputs for libvirt network name, image pool and uri and the sqlite database file.

### Building the binary

```bash
make
export $PATH=$PATH:$(git rev-parse --show-toplevel)/bin
libvirt-aws --help
```

## Running the server

The server can be started from the CLI using the `start` command. 

It provides customization options via command line flags:

```bash
$ libvirt-aws start --help
Starts the libvirt-aws emulator server

Usage:
  libvirt-aws start [flags]

Flags:
  -b, --bind-to string              Address to listen on
  -d, --database string             Path to db file (default "pool.db")
      --debug                       Enable debug logging
  -h, --help                        help for start
  -i, --libvirt-image-pool string   Name or UUID of libvirt image pool to use for EBS emulation. (default "default")
  -n, --libvirt-network string      Name or UUID of libvirt network to use for EIP emulation. (default "default")
  -l, --libvirt-uri string          Libvirtd URI (default "qemu:///system")
  -p, --port int                    TCP port to listen on (default 5100)
  -r, --region string               AWS region to pretend to be in (default "us-east-2")
```

## Supported Endpoints

- DescribeAvailabilityZones
