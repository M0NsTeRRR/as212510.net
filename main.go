package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"gopkg.in/routeros.v2"
)

const (
	templatesDir = "templates/"
)

var (
	listenAddr       = flag.String("listenAddr", "127.0.0.1:8080", "HTTP server address")
	mikrotikAddress  = flag.String("mikrotik.addr", "", "MikroTik address")
	mikrotikUsername = flag.String("mikrotik.username", "", "MikroTik username")
	mikrotikPassword = flag.String("mikrotik.password", "", "MikroTik password")

	templates = template.Must(template.ParseFiles(
		fmt.Sprintf("%s/index.html", templatesDir),
		fmt.Sprintf("%s/prompt.html", templatesDir),
	))
)

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

func identity(c *routeros.Client, r *router) error {
	reply, err := run(c, "/system/identity/print")
	if err != nil {
		return err
	}

	r.Name = reply.Re[0].Map["name"]

	return nil
}

func bgpInstance(c *routeros.Client, r *router) error {
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

func bgpNetwork(c *routeros.Client, r *router) error {
	reply, err := run(c, "/routing/bgp/network/print")
	if err != nil {
		return err
	}

	for _, re := range reply.Re {
		r.Bgp.Prefixes = append(r.Bgp.Prefixes, re.Map["network"])
	}

	return nil
}

func bgpPeer(c *routeros.Client, r *router) error {
	reply, err := run(c, "/routing/bgp/peer/print")
	if err != nil {
		return err
	}

	for _, re := range reply.Re {
		r.Bgp.Peers = append(r.Bgp.Peers,
			peer{
				Name:            re.Map["name"],
				RemoteAs:        re.Map["remote-as"],
				RemoteAddress:   re.Map["remote-address"],
				AddressFamilies: re.Map["address-families"],
			},
		)
	}

	return nil
}

func routerInfo(r *router) error {
	c, err := routeros.Dial(*mikrotikAddress, *mikrotikUsername, *mikrotikPassword)

	if err != nil {
		return err
	}
	if identity(c, r) != nil {
		return err
	}
	if bgpInstance(c, r) != nil {
		return err
	}
	if bgpNetwork(c, r) != nil {
		return err
	}
	if bgpPeer(c, r) != nil {
		return err
	}

	return nil
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	router := router{}

	if err := routerInfo(&router); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Println(err)
		return
	}

	if err := templates.ExecuteTemplate(w, "index.html", router); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Println(err)
		return
	}
}

func main() {
	flag.Parse()

	mux := http.NewServeMux()
	mux.HandleFunc("/", viewHandler)
	mux.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./css"))))

	if err := http.ListenAndServe(*listenAddr, mux); err != nil {
		log.Fatal(err)
	}
}
