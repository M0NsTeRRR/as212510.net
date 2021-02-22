package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/getsentry/sentry-go"
	sentryhttp "github.com/getsentry/sentry-go/http"
	"gopkg.in/routeros.v2"
	"gopkg.in/yaml.v2"
)

var (
	cfg = config{}

	configPath = flag.String("config", "", "Path to config")

	templates = template.Template{}
)

type config struct {
	Sentry struct {
		Dsn string `yaml:"dsn"`
	}
	Server struct {
		Address string `yaml:"address"`
		Cwd     string `yaml:"cwd"`
	} `yaml:"server"`
	Asn      int `yaml:"asn"`
	Mikrotik struct {
		Address  string `yaml:"address"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"mikrotik"`
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

func run(c *routeros.Client, command string) (routeros.Reply, error) {
	reply, err := c.Run(command)
	if err != nil {
		return routeros.Reply{}, err
	}
	return *reply, nil
}

func (r *router) identity(c *routeros.Client) error {
	reply, err := run(c, "/system/identity/print")
	if err != nil {
		return err
	}

	r.Name = reply.Re[0].Map["name"]

	return nil
}

func (r *router) bgpInstance(c *routeros.Client) error {
	reply, err := run(c, "/routing/bgp/instance/print")
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

func (r *router) bgpNetwork(c *routeros.Client) error {
	reply, err := run(c, "/routing/bgp/network/print")
	if err != nil {
		return err
	}

	for _, re := range reply.Re {
		r.Bgp.Prefixes = append(r.Bgp.Prefixes, re.Map["network"])
	}

	return nil
}

func (r *router) bgpPeer(c *routeros.Client) error {
	reply, err := run(c, "/routing/bgp/peer/print")
	if err != nil {
		return err
	}

	for _, re := range reply.Re {
		if re.Map["remote-as"] != strconv.Itoa(cfg.Asn) {
			r.Bgp.Peers = append(r.Bgp.Peers,
				peer{
					Name:            re.Map["name"],
					RemoteAs:        re.Map["remote-as"],
					RemoteAddress:   re.Map["remote-address"],
					AddressFamilies: re.Map["address-families"],
				},
			)
		}
	}

	return nil
}

func (r *router) information(c *routeros.Client) error {
	if err := r.identity(c); err != nil {
		return err
	}
	if err := r.bgpInstance(c); err != nil {
		return err
	}
	if err := r.bgpNetwork(c); err != nil {
		return err
	}
	if err := r.bgpPeer(c); err != nil {
		return err
	}

	return nil
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	router := router{}

	c, err := routeros.Dial(cfg.Mikrotik.Address, cfg.Mikrotik.Username, cfg.Mikrotik.Password)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if err := router.information(c); err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if err := templates.ExecuteTemplate(w, "index.html", router); err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func newConfig(path string, config *config) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	d := yaml.NewDecoder(f)
	if err := d.Decode(config); err != nil {
		return err
	}

	return nil
}

func main() {
	flag.Parse()

	err := newConfig(*configPath, &cfg)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("loaded config from file %s", *configPath)

	err = sentry.Init(sentry.ClientOptions{
		Dsn: cfg.Sentry.Dsn,
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("sentry initialized")

	templates = *template.Must(template.ParseFiles(
		filepath.Join(cfg.Server.Cwd, "./templates/index.html"),
		filepath.Join(cfg.Server.Cwd, "./templates/prompt.html"),
	))

	sentryHandler := sentryhttp.New(sentryhttp.Options{})

	mux := http.NewServeMux()
	mux.HandleFunc("/", sentryHandler.HandleFunc(viewHandler))
	mux.Handle("/css/", sentryHandler.Handle(http.StripPrefix("/css/", http.FileServer(http.Dir(filepath.Join(cfg.Server.Cwd, "./css"))))))

	log.Printf("server is starting on %s", cfg.Server.Address)
	if err := http.ListenAndServe(cfg.Server.Address, mux); err != nil {
		log.Fatal(err)
	}
}
