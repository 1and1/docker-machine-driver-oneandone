package oneandone

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"strconv"
	"strings"
	"time"

	oneandone "github.com/1and1/oneandone-cloudserver-sdk-go"
	"github.com/docker/machine/libmachine/drivers"
	"github.com/docker/machine/libmachine/log"
	"github.com/docker/machine/libmachine/mcnflag"
	"github.com/docker/machine/libmachine/ssh"
	"github.com/docker/machine/libmachine/state"
)

type Driver struct {
	drivers.BaseDriver
	APIKey           string
	MachineID        string
	Description      string
	Appliance        string
	Datacenter       string
	Size             string
	SSHKey           string
	FirewallPolicy   string
	LoadBalancer     string
	MonitoringPolicy string

	applianceID  string
	datacenterID string
	sizeID       string
	password     string
	publicIpID   string
	client       *oneandone.API
}

const (
	defaultOS            = "ubuntu1404-64std"
	defaultSize          = "M"
	defaultDatacenter    = "US"
	firewallPolicyPrefix = "Docker-Driver-Required-Policy_"
)

// GetCreateFlags registers the flags this driver adds to
// "docker hosts create"
func (d *Driver) GetCreateFlags() []mcnflag.Flag {
	return []mcnflag.Flag{
		mcnflag.StringFlag{
			EnvVar: "ONEANDONE_API_KEY",
			Name:   "oneandone-api-key",
			Usage:  "1&1 API Key",
		},
		mcnflag.StringFlag{
			EnvVar: "ONEANDONE_SSH_PASSWORD",
			Name:   "oneandone-ssh-pass",
			Usage:  "1&1 SSH Password",
		},
		mcnflag.StringFlag{
			EnvVar: "ONEANDONE_DATACENTER",
			Name:   "oneandone-datacenter",
			Usage:  "1&1 Data Center. Default: " + defaultDatacenter,
			Value:  defaultDatacenter,
		},
		mcnflag.StringFlag{
			EnvVar: "ONEANDONE_SIZE",
			Name:   "oneandone-size",
			Usage:  "1&1 Cloud Server Size. Default: " + defaultSize,
			Value:  defaultSize,
		},
		mcnflag.StringFlag{
			EnvVar: "ONEANDONE_OS",
			Name:   "oneandone-os",
			Usage:  "1&1 Appliance OS. Default: " + defaultOS,
			Value:  defaultOS,
		},
		mcnflag.StringFlag{
			EnvVar: "ONEANDONE_FIREWALL",
			Name:   "oneandone-firewall-id",
			Usage:  "1&1 Firewall Policy ID",
		},
		mcnflag.StringFlag{
			EnvVar: "ONEANDONE_IP_ADDRESS",
			Name:   "oneandone-ip-address",
			Usage:  "Unassigned 1&1 Public IP Address",
		},
		mcnflag.StringFlag{
			EnvVar: "ONEANDONE_LOADBALANCER",
			Name:   "oneandone-loadbalancer-id",
			Usage:  "1&1 Load Balancer ID",
		},
		mcnflag.StringFlag{
			EnvVar: "ONEANDONE_MONITOR_POLICY",
			Name:   "oneandone-monitor-policy-id",
			Usage:  "1&1 Monitoring Policy ID",
		},
		mcnflag.StringFlag{
			EnvVar: "ONEANDONE_SERVER_DESCRIPTION",
			Name:   "oneandone-server-description",
			Usage:  "1&1 Cloud Server Description",
		},
	}
}

func (d *Driver) GetSSHHostname() (string, error) {
	return d.GetIP()
}

// DriverName returns the name of the driver
func (d *Driver) DriverName() string {
	return "oneandone"
}

func (d *Driver) SetConfigFromFlags(flags drivers.DriverOptions) error {
	d.APIKey = flags.String("oneandone-api-key")
	d.IPAddress = flags.String("oneandone-ip-address")
	d.Appliance = flags.String("oneandone-os")
	d.Datacenter = flags.String("oneandone-datacenter")
	d.Size = flags.String("oneandone-size")
	d.FirewallPolicy = flags.String("oneandone-firewall-id")
	d.LoadBalancer = flags.String("oneandone-loadbalancer-id")
	d.MonitoringPolicy = flags.String("oneandone-monitor-policy-id")
	d.Description = flags.String("oneandone-server-description")
	d.SSHUser = drivers.DefaultSSHUser
	d.SSHPort = drivers.DefaultSSHPort
	d.password = flags.String("oneandone-ssh-pass")
	d.SetSwarmConfigFromFlags(flags)

	return nil
}

func (d *Driver) PreCreateCheck() error {

	log.Info("Validating 1&1 Cloud Server parameters...")

	if d.APIKey == "" {
		return fmt.Errorf("oneandone driver requires the --oneandone-api-key option")
	}

	if err := d.validateFirewallPolicy(); err != nil {
		return err
	}

	if err := d.validateDatacenter(); err != nil {
		return err
	}

	if err := d.validateFlavorSize(); err != nil {
		return err
	}

	if err := d.validateOS(); err != nil {
		return err
	}

	if err := d.validatePublicIP(); err != nil {
		return err
	}

	if err := d.validateLoadBalancer(); err != nil {
		return err
	}

	if err := d.validateMonitorPolicy(); err != nil {
		return err
	}

	return nil
}

func (d *Driver) Create() error {
	client := d.getClient()

	if d.FirewallPolicy == "" {
		log.Infof("Creating a firewall policy '%s'...", firewallPolicyPrefix)
		rand.Seed(time.Now().UnixNano())

		fpr := oneandone.FirewallPolicyRequest{
			Name: firewallPolicyPrefix + strconv.Itoa(rand.Intn(999999)),
			Rules: []oneandone.FirewallPolicyRule{
				{
					Protocol: "TCP",
					PortFrom: oneandone.Int2Pointer(22),
					PortTo:   oneandone.Int2Pointer(22),
				},
				{
					Protocol: "TCP",
					PortFrom: oneandone.Int2Pointer(80),
					PortTo:   oneandone.Int2Pointer(80),
				},
				{
					Protocol: "TCP",
					PortFrom: oneandone.Int2Pointer(2375),
					PortTo:   oneandone.Int2Pointer(2376),
				},
				{
					Protocol: "TCP",
					PortFrom: oneandone.Int2Pointer(3375),
					PortTo:   oneandone.Int2Pointer(3376),
				},
			},
		}
		_, fp, err := client.CreateFirewallPolicy(&fpr)
		if err == nil {
			log.Debug("Waiting for firewall policy to be available...")
			client.WaitForState(fp, "ACTIVE", 10, 30)
		} else {
			return err
		}
		d.FirewallPolicy = fp.Id
	}

	log.Debug("Creating SSH key...")
	key, err := d.createSSHKey()
	if err != nil {
		return err
	}
	d.SSHKey = key

	log.Info("Creating 1&1 Cloud Server...")

	request := oneandone.ServerRequest{
		Name:        d.MachineName,
		Description: d.Description,
		ApplianceId: d.applianceID,
		PowerOn:     true,
		Password:    d.password,
		Hardware: oneandone.Hardware{
			FixedInsSizeId: d.sizeID,
		},
		DatacenterId:       d.datacenterID,
		FirewallPolicyId:   d.FirewallPolicy,
		IpId:               d.publicIpID,
		LoadBalancerId:     d.LoadBalancer,
		MonitoringPolicyId: d.MonitoringPolicy,
		SSHKey:             d.SSHKey,
	}
	_, machine, err := client.CreateServer(&request)

	if err != nil {
		return err
	}
	d.MachineID = machine.Id

	if d.password == "" && machine.FirstPassword != "" {
		d.password = machine.FirstPassword
	}

	log.Info("Waiting for IP address to become available...")

	for {
		machine, err = client.GetServer(d.MachineID)
		if err != nil {
			return err
		}
		if len(machine.Ips) > 0 {
			if machine.Ips[0].Ip != "" {
				d.IPAddress = machine.Ips[0].Ip
				break
			}
		}

		log.Debug("IP address not yet available")
		time.Sleep(5 * time.Second)
	}

	log.Info("Finishing configuration and starting machine...")
	err = client.WaitForState(machine, "POWERED_ON", 10, 400)
	if err != nil {
		return err
	}
	log.Infof("Created 1&1 Cloud Server, ID: %s, Public IP: %s", d.MachineID, d.IPAddress)

	return nil
}

func (d *Driver) createSSHKey() (string, error) {
	if err := ssh.GenerateSSHKey(d.GetSSHKeyPath()); err != nil {
		return "", err
	}

	publicKey, err := ioutil.ReadFile(d.publicSSHKeyPath())
	if err != nil {
		return "", err
	}

	return string(publicKey), nil
}

func (d *Driver) GetURL() (string, error) {
	if err := drivers.MustBeRunning(d); err != nil {
		return "", err
	}

	ip, err := d.GetIP()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("tcp://%s:2376", ip), nil
}

func (d *Driver) GetState() (state.State, error) {
	status, err := d.getClient().GetServerStatus(d.MachineID)
	if err != nil {
		return state.Error, err
	}
	switch status.State {
	case "POWERING_ON":
		return state.Starting, nil
	case "POWERED_ON":
		return state.Running, nil
	case "POWERING_OFF":
		return state.Stopping, nil
	case "POWERED_OFF":
		return state.Stopped, nil
	}
	return state.None, nil
}

func (d *Driver) Start() error {
	vmState, err := d.GetState()
	if err != nil {
		return err
	}
	if vmState == state.Running || vmState == state.Starting {
		log.Infof("Host is already running or starting")
		return nil
	}
	log.Debugf("starting %s", d.MachineName)
	_, err = d.getClient().StartServer(d.MachineID)
	return err
}

func (d *Driver) Stop() error {
	vmState, err := d.GetState()
	if err != nil {
		return err
	}
	if vmState == state.Stopped {
		log.Infof("Host is already stopped")
		return nil
	}
	log.Debugf("stopping %s", d.MachineName)
	_, err = d.getClient().ShutdownServer(d.MachineID, false)
	return err
}

func (d *Driver) Remove() error {
	client := d.getClient()

	// Delete the virtual machine
	log.Debugf("removing %s", d.MachineName)
	_, err := client.DeleteServer(d.MachineID, false)
	if err != nil {
		if strings.Contains(err.Error(), "NOT_FOUND") {
			log.Infof("1&1 server '%s' doesn't exist, assuming it is already deleted", d.MachineName)
		} else {
			return err
		}
	}

	return nil
}

func (d *Driver) Restart() error {
	vmState, err := d.GetState()
	if err != nil {
		return err
	}
	if vmState == state.Stopped {
		log.Infof("Host is already stopped, use start command to run it")
		return nil
	}
	log.Debugf("restarting %s", d.MachineName)
	_, err = d.getClient().RebootServer(d.MachineID, false)
	return err
}

func (d *Driver) Kill() error {
	vmState, err := d.GetState()
	if err != nil {
		return err
	}
	if vmState == state.Stopped {
		log.Infof("Host is already stopped")
		return nil
	}
	log.Debugf("killing %s", d.MachineName)
	_, err = d.getClient().ShutdownServer(d.MachineID, true) // hardware shutdown
	return err
}

func (d *Driver) getClient() *oneandone.API {
	log.Debug("getting client")
	if d.client == nil {
		d.client = oneandone.New(d.APIKey, oneandone.BaseUrl)
	}
	return d.client
}

func (d *Driver) publicSSHKeyPath() string {
	return d.GetSSHKeyPath() + ".pub"
}

func (d *Driver) validateOS() error {
	log.Debugf("validating  '%s' server appliance", d.Appliance)
	appliances, err := d.getClient().ListServerAppliances(0, 0, "", d.Appliance, "")
	if err != nil {
		return err
	}

	for _, sa := range appliances {
		if strings.ToLower(d.Appliance) == strings.ToLower(sa.Name) {
			if strings.ToLower(sa.OsFamily) != "linux" {
				return fmt.Errorf("1&1 docker driver supports only linux OS")
			}
			d.applianceID = sa.Id
			return nil
		}
	}

	return fmt.Errorf("1&1 server appliance '%s' not found", d.Appliance)
}

func (d *Driver) validateDatacenter() error {
	if d.Datacenter == "" || strings.ToUpper(d.Datacenter) == defaultDatacenter {
		return nil
	}
	log.Debugf("validating  '%s' datacenter", d.Datacenter)
	dcs, err := d.getClient().ListDatacenters()
	if err != nil {
		return err
	}
	for _, dc := range dcs {
		if strings.ToUpper(d.Datacenter) == strings.ToUpper(dc.CountryCode) {
			d.datacenterID = dc.Id
			return nil
		}
	}
	return fmt.Errorf("1&1 datacenter '%s' could not be found", d.Datacenter)
}

func (d *Driver) validatePublicIP() error {
	if d.IPAddress == "" {
		return nil
	}
	log.Debugf("validating 1&1 IP address '%s'", d.IPAddress)

	ips, err := d.getClient().ListPublicIps(0, 0, "", d.IPAddress, "")
	if err != nil {
		if strings.Contains(err.Error(), "NOT_FOUND") || strings.Contains(err.Error(), "not found") {
			return fmt.Errorf("1&1 public IP '%s' could not be found", d.IPAddress)
		}
		return err
	}
	for _, pa := range ips {
		if d.IPAddress == pa.IpAddress {
			if pa.AssignedTo != nil {
				return fmt.Errorf("public IP '%s' is already assigned to server", d.IPAddress)
			}
			d.publicIpID = pa.Id
			return nil
		}
	}
	return fmt.Errorf("1&1 public IP '%s' could not be found", d.IPAddress)
}

func (d *Driver) validateLoadBalancer() error {
	if d.LoadBalancer == "" {
		return nil
	}
	log.Debugf("validating 1&1 load balancer ID '%s'", d.LoadBalancer)

	_, err := d.getClient().GetLoadBalancer(d.LoadBalancer)
	if err != nil {
		if strings.Contains(err.Error(), "NOT_FOUND") {
			return fmt.Errorf("load balancer with ID '%s' could not be found", d.LoadBalancer)
		}
		return err
	}
	return nil
}

func (d *Driver) validateMonitorPolicy() error {
	if d.MonitoringPolicy == "" {
		return nil
	}
	log.Debugf("validating 1&1 monitoring policy ID '%s'", d.MonitoringPolicy)

	_, err := d.getClient().GetMonitoringPolicy(d.MonitoringPolicy)
	if err != nil {
		if strings.Contains(err.Error(), "NOT_FOUND") {
			return fmt.Errorf("monitoring policy with ID '%s' could not be found", d.MonitoringPolicy)
		}
		return err
	}
	return nil
}

func (d *Driver) validateFlavorSize() error {
	log.Debugf("validating  '%s' flavor", d.Size)
	if strings.ToUpper(d.Size) == "S" {
		return fmt.Errorf("1&1 Cloud Server size should be 'M' or larger")
	}

	fixedSizes, err := d.getClient().ListFixedInstanceSizes()
	if err != nil {
		return err
	}

	for _, f := range fixedSizes {
		if strings.ToUpper(d.Size) == strings.ToUpper(f.Name) {
			d.sizeID = f.Id
			return nil
		}
	}

	return fmt.Errorf("Cloud Server size '%s' not available", d.Size)
}

func (d *Driver) validateFirewallPolicy() error {
	if d.FirewallPolicy == "" {
		policies, _ := d.getClient().ListFirewallPolicies(0, 0, "", firewallPolicyPrefix, "")

		if len(policies) > 0 {
			var p22, p80, p2375, p2376, p3375, p3376 bool
			for _, pol := range policies {
				if strings.HasPrefix(pol.Name, firewallPolicyPrefix) {
					for _, rule := range pol.Rules {
						if *(rule.PortTo) == 22 || *(rule.PortFrom) == 22 {
							p22 = true
							continue
						}
						if *(rule.PortTo) == 80 || *(rule.PortFrom) == 80 {
							p80 = true
							continue
						}
						if *(rule.PortTo) == 2375 || *(rule.PortFrom) == 2375 {
							p2375 = true
							continue
						}
						if *(rule.PortTo) == 2376 || *(rule.PortFrom) == 2376 {
							p2376 = true
							continue
						}
						if *(rule.PortTo) == 3375 || *(rule.PortFrom) == 3375 {
							p3375 = true
							continue
						}
						if *(rule.PortTo) == 3376 || *(rule.PortFrom) == 3376 {
							p3376 = true
							continue
						}
					}
					if p22 && p80 && (p2375 || p2376) && (p3375 || p3376) {
						d.FirewallPolicy = pol.Id
						return nil
					}
				}
			}
		}
		return nil
	}
	_, err := d.getClient().GetFirewallPolicy(d.FirewallPolicy)
	if err != nil {
		if strings.Contains(err.Error(), "NOT_FOUND") {
			return fmt.Errorf("firewall policy with ID '%s' could not be found", d.FirewallPolicy)
		}
		return fmt.Errorf("Failed to validate firewall policy. Error: %s", err.Error())
	}
	return nil
}
