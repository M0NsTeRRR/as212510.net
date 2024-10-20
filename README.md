as212510.net website

[![Publish releases](https://github.com/M0NsTeRRR/as212510.net/actions/workflows/releases.yml/badge.svg)](https://github.com/M0NsTeRRR/as212510.net/actions/workflows/releases.yml)

# Dev
## Requirements
- MikroTik RouterOS v7
- [Devcontainer](https://code.visualstudio.com/docs/devcontainers/containers)

## Build
`$ go run .`

# Usage
## Helm chart
See [helm-charts](https://github.com/M0NsTeRRR/helm-charts)

## Docker
`docker pull ghcr.io/m0nsterrr/as212510.net:latest`

## Non docker
Download the binary  

Set environment variables :
```bash
# Mandatory
AS212510_NET_ASN="212510"
AS212510_NET_MIKROTIK_ADDRESS="192.168.0.1:8728"
AS212510_NET_MIKROTIK_USERNAME="as212510.net"
AS212510_NET_MIKROTIK_PASSWORD="password"
AS212510_NET_MIKROTIK_BGPFIREWALLADDRESSLISTV6="bgp-networks"
# Optional
AS212510_NET_HEALTHCHECK_ADDRESS=":10240" # default to :10240
AS212510_NET_EXPORTER_ADDRESS=":10241" # default to :10241
AS212510_NET_SERVER_ADDRESS=":8080" # default to :8080
```

`$ as212510.net`

# Contributing

We welcome and encourage contributions to this project! Please read the [Contributing guide](CONTRIBUTING.md). Also make sure to check the [Code of Conduct](CODE_OF_CONDUCT.md) and adhere to its guidelines

# Security

See [SECURITY.md](SECURITY.md) file for details.

# Licence

The code is under CeCILL license.

You can find all details here: https://cecill.info/licences/Licence_CeCILL_V2.1-en.html

# Credits

Copyright © Ludovic Ortega, 2023

Contributor(s):

-Ortega Ludovic - ludovic.ortega@adminafk.fr
