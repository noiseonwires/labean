# Labean

[![Go Report Card](https://goreportcard.com/badge/github.com/noiseonwires/labean)](https://goreportcard.com/report/github.com/noiseonwires/labean)
[![Go Version](https://img.shields.io/github/go-mod/go-version/noiseonwires/labean)](https://github.com/noiseonwires/labean/blob/master/go.mod)
[![License](https://img.shields.io/github/license/noiseonwires/labean)](https://github.com/noiseonwires/labean/blob/master/LICENSE)
[![Latest Release](https://img.shields.io/github/v/release/noiseonwires/labean)](https://github.com/noiseonwires/labean/releases)
[![Last Commit](https://img.shields.io/github/last-commit/noiseonwires/labean)](https://github.com/noiseonwires/labean/commits)
[![Stars](https://img.shields.io/github/stars/noiseonwires/labean?style=social)](https://github.com/noiseonwires/labean/stargazers)

Labean is a simple HTTP/HTTPS web (port) knocker for GNU/Linux.

### What is port knocking?
Port knocking is a method of externally opening ports on a firewall by doing some actions, e.g. generating a connection attempt on a set of prespecified closed ports.  [[Wikipedia](https://en.wikipedia.org/wiki/Port_knocking)]

The primary purpose of port knocking is to prevent an attacker from scanning a system for potentially exploitable services by doing a port scan, because unless the attacker sends the correct knock sequence, the protected ports will appear closed.
Another purpose is to disguise some services (e.g. VPN or proxy) running on the server. For example, the server receives connections on 443 port and acts as a usual HTTPS web server, but after knocking it will modify firewall rules to allow specified IP to connect to another service on the same port - for third-party observers (like corporate or ISP Deep-Packet-Inspection systems) it will look exactly like ordinary HTTPS.

### Why Labean?
Classic implementations of port knockers (like [knockd](http://www.zeroflux.org/projects/knock)) allow a client to open ports or start services by generating a connection attempt on a set of prespecified closed ports. This is simple and usually reliable, but there are some tricky cases. First of all, this requires using special clients or scripts on the client device (including mobile gadgets or network routers). A more significant problem is that all ports and protocols except standard 80 (HTTP) and 443 (HTTPS) can be banned on a corporate or ISP firewall, so you can't use 'classic' knocking in this case. That's why Labean was created.
SStarted back in 2018, Labean was among the earliest web-based port knockers of its kind, and it remains a long-living project to this day.

### How does it work?
Briefly: there is a front-end web server (like [nginx](http://nginx.org/) or [caddy](https://caddyserver.com/)) running on your VDS/VPS/etc., and it serves some ordinary web content like cute kittens' videos, Linux distros' ISOs, or a Wikipedia mirror. But when you want to connect to the hidden service (VPN, proxy, ssh daemon, etc.), you perform a GET request (using cURL or a web browser) like:

`https://someserver.org/secret/vpn/on`,

process basic authentication, then your front-end web server reverse-proxies the request to Labean, and Labean starts the service or applies firewall rules to open access exclusively for your IP address. When you finish working with the hidden service, you just make a GET request to

`https://someserver.org/secret/vpn/off`

and the access becomes closed.

Another way is to use *automatic timeout control*.

For example, if you only need to open a firewall port temporarily to accept a connection, you can set a timeout of 30 or 60 seconds. After making the GET request, you start your client and initiate a connection to the hidden service; then the timeout expires and the rule is deleted. Your established connection is still active, but no one else, even from your IP, can establish a new one without knocking again.

### Building and installation
Labean is written in the Go language, so you'll need the Go compiler installed on your system. Please refer to the [official Go website](https://golang.org/doc/install) and your distro's manpages to get and install the Go compiler. I tried to build Labean with Go 1.6 and 1.10 and everything was fine; perhaps it will work with even older versions.

``` # clone the latest source code:
git clone https://github.com/uprt/labean.git
cd labean
go build 
```

If everything is okay, you'll have a `labean` binary executable in your current directory.

I recommend putting the `labean` binary in the `/usr/sbin` directory and `labean.conf` in `/etc/` (see the next chapter about this configuration file).

Labean can be started simply:

`./labean <path_to_config_file>`

After starting, `labean` will send its logs to the syslog service, so you can check them in the /var/log/syslog or /var/log/messages file depending on your distro.

A better way is to run it as a service. If your distro uses [LSB](https://wiki.debian.org/LSBInitScripts)-compatible init system or [Systemd](https://www.freedesktop.org/wiki/Software/systemd/), feel free to use sample service files that you can find in `examples` dir in the source tree:
```
cp ./examples/labean.init.ex /etc/init.d/labean
/etc/init.d/labean start
update-rc.d labean defaults
# or...
cp ./examples/labean.service.ex /etc/systemd/system/labean.service # this path can be different in your distro
systemctl daemon-reload
systemctl start labean
systemctl enable labean
```

Another important step is to decide how Labean should be exposed. You have two options.

The first one is to put Labean behind a frontend web server (reverse-proxy) like Caddy or Nginx, which performs the following things:

- serve HTTPS (SSL) connections with valid certificates
- require basic HTTP authorization for the 'secret' URL
- add an 'X-Real-IP' or similar header to the HTTP request
- pass the request to Labean (running on another TCP port like 8080) for the 'secret' URL

You can find the simplest variant of such config in the `examples` directory of the source tree, but you can use any other web server/reverse-proxy if you want.

The second option is to let Labean handle everything itself: it can serve HTTPS (just point it to your `tls_cert` and `tls_key`), require a token or basic authentication, serve a static website from a directory, and read the real client IP directly from the connection. This way you don't need any reverse-proxy at all. See the config fields below for details.

### Config file
Labean config file is a simple JSON document. 
```
{
  "listen": "127.0.0.1:8080",
  "url_prefix": "secret",
  "external_ip": "192.30.253.113",
  "external_ipv6": "9c49:226f:fd31:bc73:4c84:29ed:13aa:db07",
  "real_ip_header": "X-Real-IP",
  "allow_explicit_ips": false,
  "tasks": [ 
    {
      "name": "vpn",
      "timeout": 30,
      "on_command": "iptables -t nat -A PREROUTING -p tcp -s {clientIP} --dport 443 -j REDIRECT --to-port 4443",
      "off_command": "iptables -t nat -D PREROUTING -p tcp -s {clientIP} --dport 443 -j REDIRECT --to-port 4443",
      "on_command_v6": "ip6tables -t nat -A PREROUTING -p tcp -s {clientIP} --dport 443 -j REDIRECT --to-port 4443",
      "off_command_v6": "ip6tables -t nat -D PREROUTING -p tcp -s {clientIP} --dport 443 -j REDIRECT --to-port 4443"
    },
    {
      "name": "sshd",
      "timeout": 0,
      "on_command": "/etc/init.d/sshd start",
      "off_command": "/etc/init.d/sshd stop"
    }
  ]
}
```

Here is the description of its fields:
- `"listen"`: IP address and port to listen on for incoming connections from the reverse proxy; 127.0.0.1 for IPv4 localhost, ::1 for IPv6 localhost, and so on;
- `"url_prefix"`: if your reverse-proxy can't rewrite URLs, you can set the prefix to trim; if your reverse-proxy performs trimming, leave it empty;
- `"external_ip"`: IP of your server. This is not a 'must-have' option; it's just syntactic sugar to replace {serverIP} in the command strings if you need it;
- `"external_ipv6"`: the same, but for IPv6 (if you use it);
- `"real_ip_header"`: the name of the HTTP header with the real client IP added by the reverse-proxy (usually it is "X-Real-IP" or "X-Forwarded-For"); leave it empty if you don't use a reverse-proxy;
- `"allow_explicit_ips"`: if true, you can explicitly specify the client IP address in the GET request, like `https://someserver.org/secret/service/off/?ip=123.56.78.9`. This can be helpful when you've established a VPN connection and later want to manually deactivate the secret service using Labean's internal (local) IP inside the tunnel. This feature allows you to stop services started by other users, so use it very carefully; it is disabled by default.
- `"tls_cert"` and `"tls_key"`: paths to the certificate and private key files. If both are set, Labean serves HTTPS directly instead of plain HTTP (so you can run it without a reverse-proxy). Leave them empty to serve plain HTTP;
- `"www_root"`: path to a directory to serve as a static website on the root URL. If set, Labean acts as a simple static file server for any non-task request; leave it empty to keep the default behavior;
- `"auth_token"`: if set, every task request must provide this token via the `X-Auth-Token` header or the `token` query parameter (e.g. `.../secret/vpn/on?token=...`); leave it empty to disable token auth;
- `"username"` and `"password"`: if either is set, every task request must pass HTTP Basic authentication with these credentials; leave them empty to disable basic auth;
- `"log"`: where to send logs. Use `"syslog"` (the default) to log to the system syslog service, or `"stdout"` to log to standard output. On platforms without syslog (e.g. Windows) Labean always logs to stdout regardless of this setting;
- `"tasks"`: the array of 'tasks' to start or stop hidden services;
- `"name"`: the unique name of the hidden service. You will use it in your HTTP GET queries: `https://someserver.org/secret/<name>/{on|off}`;
- `"timeout"`: timeout to automatically switch your hidden service or firewall rule "off" after "on". If it is set to 0, the timeout feature will be disabled and you'll need to switch your service "off" manually;
- `"on_command"`: command line to start your service or activate the firewall rule. You can use the `{serverIP}` and `{clientIP}` macros in it; they will be automatically replaced by the corresponding values;
- `"off_command"`: the same as `"on_command"`, but does exactly the opposite thing :)
- `"on_command_v6"` and `"off_command_v6"`: same as above, but for IPv6 (if you use it)

### Bugs and ideas?
Feel free to open issues on GitHub or make pull requests if you want to improve something. I will be grateful for any help and feedback.

### License
Labean is released under the BSD 2-Clause "Simplified" License.

Copyright (c) 2018-2026, Kirill aka Noiseonwires
