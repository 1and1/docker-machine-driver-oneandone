package main

import (
	"github.com/1and1/docker-machine-driver-oneandone"
	"github.com/docker/machine/libmachine/drivers/plugin"
)

func main() {
	plugin.RegisterDriver(new(oneandone.Driver))
}
