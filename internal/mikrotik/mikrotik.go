package mikrotik

import (
	"crypto/tls"
	"strconv"

	"github.com/go-routeros/routeros/v3"
	"github.com/m0nsterrr/as212510.net/v3/internal/config"
)

func NewRouter(config config.Config) (*Router, error) {
	var client *routeros.Client
	var err error

	if config.Mikrotik.Tls {
		client, err = routeros.DialTLS(config.Mikrotik.Address, config.Mikrotik.Username, config.Mikrotik.Password, &tls.Config{
			InsecureSkipVerify: config.Mikrotik.SkipTLSVerify,
		})
	} else {
		client, err = routeros.Dial(config.Mikrotik.Address, config.Mikrotik.Username, config.Mikrotik.Password)
	}

	if err != nil {
		return nil, err
	}

	return &Router{Asn: config.Asn, BgpFirewallAddressListV6: config.Mikrotik.BgpFirewallAddressListV6, Client: client}, nil
}

func runCommand(client *routeros.Client, command string) (routeros.Reply, error) {
	reply, err := client.Run(command)
	if err != nil {
		return routeros.Reply{}, err
	}
	return *reply, nil
}

func (r *Router) identity() error {
	reply, err := runCommand(r.Client, "/system/identity/print")
	if err != nil {
		return err
	}

	r.Name = reply.Re[0].Map["name"]

	return nil
}

func (r *Router) bgpInstance() error {
	reply, err := runCommand(r.Client, "/routing/bgp/template/print")
	if err != nil {
		return err
	}

	as, err := strconv.Atoi(reply.Re[0].Map["as"])
	if err != nil {
		return err
	}
	r.Bgp.As = as

	return nil
}

func (r *Router) bgpNetworkv6() error {
	reply, err := runCommand(r.Client, "/ipv6/firewall/address-list/print")
	if err != nil {
		return err
	}

	for _, re := range reply.Re {
		if re.Map["list"] == r.BgpFirewallAddressListV6 {
			r.Bgp.Prefixes = append(r.Bgp.Prefixes, re.Map["address"])
		}
	}

	return nil
}

func (r *Router) bgpPeer() error {
	reply, err := runCommand(r.Client, "/routing/bgp/connection/print")
	if err != nil {
		return err
	}

	for _, re := range reply.Re {
		if re.Map["remote.as"] != strconv.Itoa(r.Asn) {
			r.Bgp.Peers = append(r.Bgp.Peers,
				peer{
					Name:          re.Map["name"],
					RemoteAs:      re.Map["remote.as"],
					RemoteAddress: re.Map["remote.address"],
					Afi:           re.Map["afi"],
				},
			)
		}
	}

	return nil
}

func (r *Router) Information() error {
	if err := r.identity(); err != nil {
		return err
	}
	if err := r.bgpInstance(); err != nil {
		return err
	}
	if err := r.bgpNetworkv6(); err != nil {
		return err
	}
	if err := r.bgpPeer(); err != nil {
		return err
	}

	return nil
}

func (r *Router) Close() {
	r.Client.Close()
}
