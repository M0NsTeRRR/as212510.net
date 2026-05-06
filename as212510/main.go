package as212510

import (
	"crypto/tls"
	"embed"
	"html/template"
	"log"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/caarlos0/env/v11"
	"github.com/go-routeros/routeros/v3"
)

var (
	cfg Config

	//go:embed all:templates/*.html
	tempFs embed.FS

	//go:embed static
	staticFiles embed.FS

	tmpl *template.Template

	version = "development"

	buildTime = "0"
)

type Config struct {
	HealthCheck struct {
		Address string `env:"ADDRESS" envDefault:":10240"`
	} `envPrefix:"HEALTHCHECK_"`
	Metric struct {
		Address string `env:"ADDRESS" envDefault:":10241"`
	} `envPrefix:"METRIC_"`
	Server struct {
		Address string `env:"ADDRESS" envDefault:":8080"`
	} `envPrefix:"SERVER_"`
	Asn      int `env:"ASN,required"`
	Mikrotik struct {
		Address                  string `env:"ADDRESS,required"`
		Tls                      bool   `env:"TLS" envDefault:"true"`
		SkipTLSVerify            bool   `env:"SKIP_TLS_VERIFY" envDefault:"true"`
		Username                 string `env:"USERNAME,required"`
		Password                 string `env:"PASSWORD,required"`
		BgpFirewallAddressListV6 string `env:"BGP_FIREWALL_ADDRESS_LIST_V6,required"`
	} `envPrefix:"MIKROTIK_"`
}

type peer struct {
	Name            string
	RemoteAs        string
	RemoteAddress   string
	AddressFamilies string
}

type bgp struct {
	As       int
	Prefixes []string
	Peers    []peer
}

type router struct {
	Name   string
	Bgp    bgp
	Client *routeros.Client
}

func runCommand(client *routeros.Client, command string) (routeros.Reply, error) {
	reply, err := client.Run(command)
	if err != nil {
		return routeros.Reply{}, err
	}
	return *reply, nil
}

func (r *router) identity() error {
	reply, err := runCommand(r.Client, "/system/identity/print")
	if err != nil {
		return err
	}

	r.Name = reply.Re[0].Map["name"]

	return nil
}

func (r *router) bgpInstance() error {
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

func (r *router) bgpNetworkv6() error {
	reply, err := runCommand(r.Client, "/ipv6/firewall/address-list/print")
	if err != nil {
		return err
	}

	for _, re := range reply.Re {
		if re.Map["list"] == cfg.Mikrotik.BgpFirewallAddressListV6 {
			r.Bgp.Prefixes = append(r.Bgp.Prefixes, re.Map["address"])
		}
	}

	return nil
}

func (r *router) bgpPeer() error {
	reply, err := runCommand(r.Client, "/routing/bgp/connection/print")
	if err != nil {
		return err
	}

	for _, re := range reply.Re {
		if re.Map["remote.as"] != strconv.Itoa(cfg.Asn) {
			r.Bgp.Peers = append(r.Bgp.Peers,
				peer{
					Name:            re.Map["name"],
					RemoteAs:        re.Map["remote.as"],
					RemoteAddress:   re.Map["remote.address"],
					AddressFamilies: re.Map["address-families"],
				},
			)
		}
	}

	return nil
}

func (r *router) information() error {
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

func viewHandler(w http.ResponseWriter, r *http.Request) {
	var client *routeros.Client
	var err error

	if cfg.Mikrotik.Tls {
		client, err = routeros.DialTLS(cfg.Mikrotik.Address, cfg.Mikrotik.Username, cfg.Mikrotik.Password, &tls.Config{
			InsecureSkipVerify: cfg.Mikrotik.SkipTLSVerify,
		})
	} else {
		client, err = routeros.Dial(cfg.Mikrotik.Address, cfg.Mikrotik.Username, cfg.Mikrotik.Password)
	}
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	defer client.Close()

	router := router{Client: client}

	if err := router.information(); err != nil {
		slog.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if err := tmpl.ExecuteTemplate(w, "index.html", router); err != nil {
		slog.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func Run() {
	slog.Info("Starting as212510", "version", version, "build_time", buildTime)

	err := env.ParseWithOptions(&cfg, env.Options{Prefix: "AS212510_NET_"})
	if err != nil {
		log.Fatal(err.Error())
	}

	staticFS := http.FS(staticFiles)
	fs := http.FileServer(staticFS)

	tmpl = template.Must(template.ParseFS(tempFs, "templates/*.html"))

	go startHealthcheck(cfg.HealthCheck.Address)
	go startMetric(cfg.Metric.Address)

	mux := http.NewServeMux()
	mux.HandleFunc("/", viewHandler)
	mux.Handle("/static/", fs)

	slog.Info("Server is starting", "address", cfg.Server.Address)
	if err := http.ListenAndServe(cfg.Server.Address, mux); err != nil {
		log.Fatal(err.Error())
	}
}
