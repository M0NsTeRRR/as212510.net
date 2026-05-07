package config

import (
	"log"

	"github.com/caarlos0/env/v11"
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
		BgpFirewallAddressListV6 string `env:"BGP_FIREWALL_ADDRESSLIST_V6,required"`
	} `envPrefix:"MIKROTIK_"`
}

func Init() Config {
	cfg := Config{}
	if err := env.ParseWithOptions(&cfg, env.Options{Prefix: "AS212510_NET_"}); err != nil {
		log.Fatalf("error reading configuration from environment: %v", err)
	}
	return cfg
}
