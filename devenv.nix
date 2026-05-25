{
  pkgs,
  lib,
  config,
  ...
}:
{
  env = {
    AS212510_NET_MIKROTIK_ADDRESS = config.secretspec.secrets.AS212510_NET_MIKROTIK_ADDRESS or "";
    AS212510_NET_MIKROTIK_USERNAME = config.secretspec.secrets.AS212510_NET_MIKROTIK_USERNAME or "";
    AS212510_NET_MIKROTIK_PASSWORD = config.secretspec.secrets.AS212510_NET_MIKROTIK_PASSWORD or "";
  };

  packages = [
    pkgs.secretspec
  ];

  languages.go.enable = true;
}
