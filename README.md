# 1&amp;1 Cloud Driver for Docker Machine

## Table of Contents
* [Overview](#overview)
* [Requirements](#requirements)
* [Installation](#installation)
  * [From a Release](#from-a-release)
  * [From Source](#from-source)
* [Usage](#usage)
  * [Available Options](#available-options)
  * [Example](#example)
* [License](#license)

## Overview

The 1&amp;1 Cloud Driver is a plugin for Docker Machine which allows you to automate the provisioning of Docker hosts on 1&amp;1 Cloud Servers. The plugin is based on the [1&amp;1 CloudServer Go SDK](https://github.com/1and1/oneandone-cloudserver-sdk-go) and [Cloud API](https://cloudpanel-api.1and1.com/documentation/1and1/). 

To acquire 1&amp;1 Cloud API credentials visit https://www.1and1.com.

## Requirements

  * [Docker Machine](https://docs.docker.com/machine/install-machine/) 0.5.1 or a newer version

Windows and Mac OS X users may install [Docker Toolbox](https://www.docker.com/products/docker-toolbox) package that includes the latest version of the Docker Machine.

## Installation

### From a Release

The latest version of the `docker-machine-driver-oneandone` binary is available on the [GithHub Releases](https://github.com/1and1/docker-machine-driver-oneandone/releases) page.
Download the `tar` archive and extract it into a directory residing in your PATH. Select the binary that corresponds to your OS and according to the file name prefix:

* Linux: docker-machine-driver-oneandone-linux
* Mac OS X: docker-machine-driver-oneandone-darwin
* Windows: docker-machine-driver-oneandone-windows

To extract and install the binary, Linux and Mac users can use the Terminal and the following commands:

```bash
sudo tar -C /usr/local/bin -xvzf docker-machine-driver-oneandone*.tar.gz
```

If required, modify the permissions to make the plugin executable:

```bash
sudo chmod +x /usr/local/bin/docker-machine-driver-oneandone
```

Windows users may run the above commands without `sudo` in Docker Quickstart Terminal that is installed with [Docker Toolbox](https://www.docker.com/products/docker-toolbox).

### From Source

Make sure you have installed [Go](http://www.golang.org) and configured [GOPATH](http://golang.org/doc/code.html#GOPATH) properly.

To download the repository and build the driver run the following:

```bash
go get -d -u github.com/1and1/docker-machine-driver-oneandone
cd $GOPATH/src/github.com/1and1/docker-machine-driver-oneandone
make build
```

To use the driver run:

```bash
make install
```

This command will install the driver into `/usr/local/bin`. 

Otherwise, set your PATH environment variable correctly. For example:

```bash
export PATH=$GOPATH/src/github.com/1and1/docker-machine-driver-oneandone/bin:$PATH
```

If you are running Windows, you may also need to install GNU Make, Bash shell and a few other Bash utilities available with [Cygwin](https://www.cygwin.com).

## Usage

You may want to refer to the Docker Machine [official documentation](https://docs.docker.com/machine/) before using the driver.

Verify that Docker Machine can see the 1&amp;1 driver:

```bash
docker-machine create -d oneandone --help
```

To create a Docker host, provide your API access token and a firewall policy ID.
Make sure that the policy has opened [the ports required by Docker](https://docs.docker.com/swarm/plan-for-production/).

If you do not supply a firewall policy ID, the driver will create a new firewall policy with a prefix of `Docker-Driver-Required-Policy_`. The policy will open the following TCP ports: 

* 22
* 80
* 2375-2376
* 3375-3376

You may further customize the policy according to your needs and Docker's requirements.
Before creating a new policy, the plugin will search for an existing policy with the `Docker-Driver-Required-Policy_` prefix and the required ports.

```bash
docker-machine create -d oneandone \
--oneandone-api-key [API-TOKEN-KEY] \
--oneandone-firewall-id [FIREWALL-POLICY-ID] \
MyHostName
```

### Available Options

  * `--oneandone-api-key`:  1&amp;1 Cloud API key.
  * `--oneandone-datacenter`: 1&amp;1 data center location.
  * `--oneandone-firewall-id`: 1&amp;1 firewall policy ID.
  * `--oneandone-size`: 1&amp;1 Cloud Server size.
  * `--oneandone-ip-address`: Unassigned 1&amp;1 public IP address.
  * `--oneandone-loadbalancer-id`: 1&amp;1 load balancer ID.
  * `--oneandone-monitor-policy-id`: 1&amp;1 monitoring policy ID.
  * `--oneandone-os`: 1&amp;1 server appliance OS.
  * `--oneandone-server-description`: 1&amp;1 Cloud Server description.
  * `--oneandone-ssh-pass`: SSH password.

|          CLI Option              |   Default Value    | Environment Variable           | Required |
| -------------------------------- | ------------------ | ------------------------------ | -------- |
| `--oneandone-api-key`            |                    | `ONEANDONE_API_KEY`            | yes      |
| `--oneandone-datacenter`         | `US`               | `ONEANDONE_DATACENTER`         | yes      |
| `--oneandone-firewall-id`        |                    | `ONEANDONE_FIREWALL`           | no       |
| `--oneandone-size`               | `M`                | `ONEANDONE_SIZE`               | yes      |
| `--oneandone-ip-address`         |                    | `ONEANDONE_IP_ADDRESS`         | no       |
| `--oneandone-loadbalancer-id`    |                    | `ONEANDONE_LOADBALANCER`       | no       |
| `--oneandone-monitor-policy-id`  |                    | `ONEANDONE_MONITOR_POLICY`     | no       |
| `--oneandone-os`                 | `ubuntu1404-64std` | `ONEANDONE_OS`                 | yes      |
| `--oneandone-server-description` |                    | `ONEANDONE_SERVER_DESCRIPTION` | no       |
| `--oneandone-ssh-pass`           |                    | `ONEANDONE_SSH_PASSWORD`       | no       |

Valid values for `--oneandone-size` are `M`, `L`, `XL`, `XXL`, `3XL`, `4XL` and `5XL`.

Available parameters for `--oneandone-datacenter` are shown in the next table.

| Parameter |                 Data Center Location                 |
|-----------|------------------------------------------------------|
| `DE`      | Germany                                              |
| `ES`      | Spain                                                |
| `GB`      | United Kingdom of Great Britain and Northern Ireland |
| `US`      | United States of America                             |

Supported values for `--oneandone-os` are listed below.

|              Parameter                |
|---------------------------------------|
| `centos7-64min`                       |
| `centos7-64std`                       |
| `centos7-64std+cpanel`                |
| `centos7-64std+Plesk12unlimited`      |
| `centos7-64std+Plesk12.5unlimited`    |
| `ubuntu1204-64min`                    |
| `ubuntu1204-64std`                    |
| `ubuntu1204-64std+Plesk12.5unlimited` |
| `ubuntu1204-64std+Plesk12unlimited`   |
| `ubuntu1404-64std`                    |
| `ubuntu1404-64min`                    |
| `ubuntu1404-64std+Plesk12unlimited`   |
| `ubuntu1404-64std+Plesk12.5unlimited` |

 
### Example

```
docker-machine --debug create --driver oneandone \
 --oneandone-api-key              b92bd5bb3dc909cfd48b6370d3bf560c \
 --oneandone-datacenter           GB \
 --oneandone-firewall-id          D8D6964A24D9A709941064AFD5FA83BC \
 --oneandone-size                 XXL \
 --oneandone-ip-address           109.228.51.239 \
 --oneandone-loadbalancer-id      F78677C1364BE21973E530FB7E7D805E \
 --oneandone-monitor-policy-id    6027B730256C9585B269DAA8B1788DEC \
 --oneandone-os                   centos7-64std \
 --oneandone-server-description   My-Docker-host-description \
 --oneandone-ssh-pass             MyPassSecret.170 \   
MyDockerHostName
```

**Note:** When creating a new machine, if you provide an IP address and a load balancer ID make sure that they belong to the same data center as the machine being created.
Also, all OS appliances may not be available in all data centers.

## License

This code is released under the Apache 2.0 License.

Copyright (c) 2016 1&amp;1 Internet SE
