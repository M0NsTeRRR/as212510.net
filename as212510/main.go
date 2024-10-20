package as212510

import (
	"embed"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/kelseyhightower/envconfig"
	"gopkg.in/routeros.v2"
)

var (
	cfg Config

	//go:embed all:templates/*.html
	tempFs embed.FS

	//go:embed static
	staticFiles embed.FS

	tmpl *template.Template

	version = "development"
)

type Config struct {
	HealthCheck struct {
		Address string `default:":10240"`
	}
	Metric struct {
		Address string `default:":10241"`
	}
	Server struct {
		Address string `default:":8080"`
	}
	Asn      int `required:"true"`
	Mikrotik struct {
		Address                  string `required:"true"`
		Username                 string `required:"true"`
		Password                 string `required:"true"`
		BgpFirewallAddressListV6 string `required:"true"`
	}
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
	Name string
	Bgp  bgp
}

func runCommand(client *routeros.Client, command string) (routeros.Reply, error) {
	reply, err := client.Run(command)
	if err != nil {
		return routeros.Reply{}, err
	}
	return *reply, nil
}

func (r *router) identity(client *routeros.Client) error {
	reply, err := runCommand(client, "/system/identity/print")
	if err != nil {
		return err
	}

	r.Name = reply.Re[0].Map["name"]

	return nil
}

func (r *router) bgpInstance(client *routeros.Client) error {
	reply, err := runCommand(client, "/routing/bgp/template/print")
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

func (r *router) bgpNetworkv6(c *routeros.Client) error {
	reply, err := runCommand(c, "/ipv6/firewall/address-list/print")
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

func (r *router) bgpPeer(client *routeros.Client) error {
	reply, err := runCommand(client, "/routing/bgp/connection/print")
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

func (r *router) information(client *routeros.Client) error {
	if err := r.identity(client); err != nil {
		return err
	}
	if err := r.bgpInstance(client); err != nil {
		return err
	}
	if err := r.bgpNetworkv6(client); err != nil {
		return err
	}
	if err := r.bgpPeer(client); err != nil {
		return err
	}

	return nil
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	router := router{}

	client, err := routeros.Dial(cfg.Mikrotik.Address, cfg.Mikrotik.Username, cfg.Mikrotik.Password)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if err := router.information(client); err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if err := tmpl.ExecuteTemplate(w, "index.html", router); err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func Run() {
	log.Printf("Starting %s %s", os.Args[0], version)

	err := envconfig.Process("as212510_net", &cfg)
	if err != nil {
		log.Fatal(err.Error())
	}

	var staticFS = http.FS(staticFiles)
	fs := http.FileServer(staticFS)

	tmpl = template.Must(template.ParseFS(tempFs, "templates/*.html"))

	go startHealthcheck(cfg.HealthCheck.Address)
	go startMetric(cfg.Metric.Address)

	mux := http.NewServeMux()
	mux.HandleFunc("/", viewHandler)
	mux.Handle("/static/", fs)

	log.Printf("Server is starting on %s", cfg.Server.Address)
	if err := http.ListenAndServe(cfg.Server.Address, mux); err != nil {
		log.Fatal(err)
	}
}
