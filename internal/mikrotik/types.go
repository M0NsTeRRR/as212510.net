package mikrotik

import "github.com/go-routeros/routeros/v3"

type peer struct {
	Name          string
	RemoteAs      string
	RemoteAddress string
	Afi           string
}

type bgp struct {
	As       int
	Prefixes []string
	Peers    []peer
}

type Router struct {
	Name                     string
	Asn                      int
	Bgp                      bgp
	BgpFirewallAddressListV6 string
	Client                   *routeros.Client
}
