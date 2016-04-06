# 1&amp;1 Cloud Driver for Docker Machine

1&amp;1 Cloud driver is a plugin for Docker Machine that allows you to automate the provisioning of Docker hosts on 1&amp;1 Server Cloud. The plugin is based on [1&amp;1 CloudServer Go SDK](https://github.com/1and1/oneandone-cloudserver-sdk-go) and [Cloud API](https://cloudpanel-api.1and1.com/documentation/1and1/). To acquire 1&amp;1 Cloud API credentials visit https://www.1and1.com.

## Requirements

  * [Docker Machine](https://docs.docker.com/engine/installation/) 0.5.1 or a newer version

## Installation

### From a Release

The latest version of `docker-machine-driver-oneandone` binary is available on the [GithHub Releases](https://github.com/1and1/docker-machine-driver-oneandone/releases) page.
Download the tar archive and extract it into a directory residing in your PATH.

```bash
sudo tar -C /usr/local/bin -xvzf docker-machine-driver-oneandone*.tar.gz
```

If required, modify the permissions and allow the plugin execution.

```bash
sudo chmod +x /usr/local/bin/docker-machine-driver-oneandone
```

### From Source

Make sure that you have installed [Go](http://www.golang.org) and configured [GOPATH](http://golang.org/doc/code.html#GOPATH) properly.

To download the repository and build the driver run the following:

```bash
go get -d -u github.com/1and1/docker-machine-driver-oneandone
make build
```

To use the driver run:

```bash
make install
```
The preceding command will install the driver into `/usr/local/bin`

Otherwise, set your PATH environment variable correctly, for instance as follows.

```bash
export PATH=$GOPATH/src/github.com/1and1/docker-machine-driver-oneandone/bin:$PATH
```

If you are running a Windows OS, you may need to install GNU Make, Bash shell and a few other bash utilities available with [Cygwin](https://www.cygwin.com).

## Usage

You might want to refer to Docker Machine [official documentation](https://docs.docker.com/machine/) before using the driver.

Verify that the Docker Machine can see 1&amp;1 driver.

```bash
docker-machine create -d oneandone --help
```
To create a docker host, provide your API access token and a firewall policy ID.
Make sure that the policy has opened [Docker required ports](https://docs.docker.com/swarm/plan-for-production/).
If you do not supply a firewall policy ID, the driver will create a new firewall policy with `Docker-Driver-Required-Policy_` prefix.
The policy will open the following TCP ports: 22, 80, 2375-2376 and 3375-3376. You may further customize the policy according to your needs and Docker requirements.
Before creating a new policy, the plugin will search for an existing policy with `Docker-Driver-Required-Policy_` prefix and the required ports.

```bash
docker-machine create -d oneandone \
--oneandone-api-key [API-TOKEN-KEY] \
--oneandone-firewall-id [FIREWALL-POLICY-ID] \
MyHostName
```

Available options:

  * `--oneandone-api-key`:  1&amp;1 Cloud API key.
  * `--oneandone-datacenter`: 1&amp;1 datacenter location.
  * `--oneandone-firewall-id`: 1&amp;1 firewall policy ID.
  * `--oneandone-flavor`: 1&amp;1 server flavor.
  * `--oneandone-ip-address`: Unassigned 1&amp;1 public IP address.
  * `--oneandone-loadbalancer-id`: 1&amp;1 load balancer ID.
  * `--oneandone-monitor-policy-id`: 1&amp;1 monitoring policy ID.
  * `--oneandone-os`: 1&amp;1 server appliance OS.
  * `--oneandone-server-description`: Server description.
  * `--oneandone-ssh-key`: SSH key.
  * `--oneandone-ssh-pass`: SSH password.

|          CLI Option              |   Default Value    | Environment Variable           | Required |
| -------------------------------- | ------------------ | ------------------------------ | -------- |
| `--oneandone-api-key`            |                    | `ONEANDONE_API_KEY`            | yes      |
| `--oneandone-datacenter`         | `US`               | `ONEANDONE_DATACENTER`         | yes      |
| `--oneandone-firewall-id`        |                    | `ONEANDONE_FIREWALL`           | no       |
| `--oneandone-flavor`             | `M`                | `ONEANDONE_FLAVOR`             | yes      |
| `--oneandone-ip-address`         |                    | `ONEANDONE_IP_ADDRESS`         | no       |
| `--oneandone-loadbalancer-id`    |                    | `ONEANDONE_LOADBALANCER`       | no       |
| `--oneandone-monitor-policy-id`  |                    | `ONEANDONE_MONITOR_POLICY`     | no       |
| `--oneandone-os`                 | `ubuntu1404-64std` | `ONEANDONE_OS`                 | yes      |
| `--oneandone-server-description` |                    | `ONEANDONE_SERVER_DESCRIPTION` | no       |
| `--oneandone-ssh-key`            | generated          | `ONEANDONE_SSH_KEY`            | yes      |
| `--oneandone-ssh-pass`           |                    | `ONEANDONE_SSH_PASSWORD`       | no       |

## License

This code is released under the Apache 2.0 License.

Copyright (c) 2016 1&amp;1 Internet SE
