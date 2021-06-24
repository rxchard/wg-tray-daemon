# wg-tray-daemon

Tested on Ubuntu 20.04

Note: this is incomplete and doesn't allow configuration at the moment.
The wireguard interface is hard-coded [here](https://github.com/rxchard/wg-tray-daemon/blob/main/pkg/wireguard/wireguard.go#L29).

Works with the [tray companion](https://github.com/rxchard/wg-tray)

The daemon has to run under a user that has access to wireguard. You could create a service for the daemon so it always runs.

Example service file:

```toml
[Unit]
Description=Wireguard Tray Daemon

[Service]
ExecStart=/opt/wg-tray-daemon
StandardOutput=null

[Install]
WantedBy=multi-user.target

```
