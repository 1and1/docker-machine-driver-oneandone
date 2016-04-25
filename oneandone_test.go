package oneandone

import (
	"testing"

	"github.com/docker/machine/libmachine/drivers"
)

func TestSetConfigFromFlags(t *testing.T) {
	driver := new(Driver)

	createFlags := &drivers.CheckDriverOptions{
		FlagsValues: map[string]interface{}{
			"oneandone-api-key":     "856fb4d0188eaf8149eae9643bacf993",
			"oneandone-ssh-pass":    "6dW51RcH33DCPZw",
			"oneandone-datacenter":  "DE",
			"oneandone-size":        "XXL",
			"oneandone-os":          "centos7-64std",
			"oneandone-firewall-id": "8F9BC3547812E0736E0AFD671EC43A1A",
		},
		CreateFlags: driver.GetCreateFlags(),
	}

	err := driver.SetConfigFromFlags(createFlags)

	if err != nil {
		t.Errorf("Setting driver create flags failed. Error: " + err.Error())
	}
	if len(createFlags.InvalidFlags) > 0 {
		t.Errorf("Expected no invalid flags but found %d.", len(createFlags.InvalidFlags))
	}
}

func TestDefaultConfigFlags(t *testing.T) {
	driver := new(Driver)

	createFlags := &drivers.CheckDriverOptions{
		FlagsValues: map[string]interface{}{
			"oneandone-api-key": "token-key",
		},
		CreateFlags: driver.GetCreateFlags(),
	}

	if err := driver.SetConfigFromFlags(createFlags); err != nil {
		t.Errorf("Setting driver create flags failed. Error: " + err.Error())
	}

	user := driver.GetSSHUsername()

	if user != "root" {
		t.Errorf("Invalid default SSH user. expected 'root', found '%s'", user)
	}

	port, _ := driver.GetSSHPort()

	if port != 22 {
		t.Errorf("Invalid default SSH port. expected 22, found %d", port)
	}
	if driver.Appliance == "" {
		t.Errorf("oneandone default OS should not be empty")
	}
	if driver.Size == "" {
		t.Errorf("oneandone default size should not be empty")
	}
	if driver.Datacenter == "" {
		t.Errorf("oneandone default datacenter should not be empty")
	}
}
