# SwarmDNS

A tiny dockerized DNS service for Docker swarm mode. It always returns the IP(s) of all active swarm nodes.

That makes it easy to host an arbitrary number of swarm services on a subdomain. Just add an `NS` record for the chosen subdomain for _every manager node_ in the swarm.

As the service works on manager nodes only, you should have more than one manager node for fail-safe operation.

## Quickstart

```bash
$ docker service create --name swarmdns \
                        -p 53:53/udp \
                        --mount type=bind,source=/var/run/docker.sock,target=/var/run/docker.sock,readonly \
                        --constraint "node.role == manager" \
                        mazzolino/swarmdns
                        --domain swarm.example.com
```

Alternatively, deploy the service stack definition supplied in this repository:

```bash
docker stack deploy --compose-file docker-compose.yml swarmdns
```

### Testing

With a cluster of 3 nodes:

```bash
$ docker node ls
ID                            HOSTNAME            STATUS              AVAILABILITY        MANAGER STATUS
4mqk9wohilllRkj7zppwie18h     swarm3              Ready               Active              Reachable
hhv80nx8r2jadchRohk4h3pfx *   swarm2              Ready               Active              Reachable
xx4zcnjnr80yletg4pnx00b4n     swarm1              Ready               Active              Leader
```

Here's the output:

```bash
$ dig +short foo.swarm.example.com @<IP OF ANY SWARM NODE>
192.168.1.230
192.168.1.231
192.168.1.232
$ dig +short bar.swarm.example.com @<IP OF ANY SWARM NODE>
192.168.1.231
192.168.1.232
192.168.1.230
```

## How it works

SwarmDNS will answer requests for `A` records only, and only for names in the domains specified at the commandline. It will always return the IP addresses of __all active nodes__ in the swarm, in random order. (The `AVAILABILITY` column in `docker node ls` shows which nodes are currently `Active`.)

The list of active nodes is refreshed once a minute. The TTL of the returned records is also set to 60 seconds.

### Options

When given the `--log` flag, every matching request will be logged to STDOUT. Example:

    Request:   172.17.0.1      foo.swarm.example.com.
    Request:   172.17.0.1      bar.swarm.example.com.

## Development

### Prerequisites

* Docker
* Docker-Compose

### Building

Just run `docker-compose build`. It builds a docker image `mazzolino/swarmdns` by default.

### Testing

(Only works if your host is a swarm manager node.)

```bash
$ docker-compose up -d
$ dig foo.swarm.example.com @localhost
```

## Credits

This is a fork of [WildDNS](https://github.com/djmaze/wilddns). The code structure was originally adopted from [microdns](https://github.com/fffaraz/microdns). Thanks!
