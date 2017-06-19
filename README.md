# WildDNS

A tiny dockerized DNS server which always returns the same IP(s).

## Quickstart

```bash
$ docker run -d --name wilddns -p 53:53/udp mazzolino/wilddns --log 192.168.1.230 192.168.1.231 192.168.1.232
```

Testing:

```bash
$ dig +short my.example.host @localhost
192.168.1.230
192.168.1.231
192.168.1.232
$ dig +short another.example.host @localhost
192.168.1.231
192.168.1.232
192.168.1.230
```

## Details

WildDNS will answer requests for `A` records only. It will always return __all__ IP addresses given on the commandline, in random order.

When given the `--log` flag, every matching request will be logged to STDOUT. Example:

    2017-06-19_20:08:41     172.17.0.1      my.example.host.
    2017-06-19_20:08:53     172.17.0.1      another.example.host.

## Development

### Prerequisites

The following tools have to be installed in order to build this:

* Docker
* The [dapper](https://github.com/rancher/dapper) binary needs to be available in your `PATH`.

### Building

Run `make`. It builds a docker image `mazzolino/wilddns` by default.

Available build targets:

* `wilddns`: build static wilddns binary
* `image` (default): build Docker image. Optionally set image name via `make IMAGE=my-image-name`
* `push`: push Docker image to Docker registry hub

## Credits

The code structure was adopted from [microdns](https://github.com/fffaraz/microdns). Thanks!
