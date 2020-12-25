as212510.net website

# Requirements

- Go (version >= 1.15)

# Build
`$ go build`

# Usage
Use the -h flag to see full usage:  

```
$ as212510.net -h
Usage of as212510.net:
  -config string
        Path to config
```

Config is writed in yaml
```
server:
  address: ":8080"
  cwd: "/opt/as212510_net"
mikrotik:
  address: 192.168.0.1:8728
  username: admin
  password: password
```

`$ as212510.net -config /etc/as212510_net/as212510_net.yml`

# Licence

The code is under CeCILL license.

You can find all details here: https://cecill.info/licences/Licence_CeCILL_V2.1-en.html

# Credits

Copyright © Ludovic Ortega, 2021

Contributor(s):

-Ortega Ludovic - ludovic.ortega@adminafk.fr