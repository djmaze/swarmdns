package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/djmaze/swarmdns/swarm"
	"github.com/miekg/dns"
)

type arrayFlags []string

func (i *arrayFlags) String() string {
	return "my string representation"
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

const NodeRefreshInterval = 60
const TTL = NodeRefreshInterval

var logger *log.Logger
var client swarm.Client
var logflag bool
var ips []string
var mutex = &sync.Mutex{}
var swarmDomains arrayFlags

func main() {
	var err error

	flag.BoolVar(&logflag, "log", false, "Log requests to stdout")
	flag.Var(&swarmDomains, "domain", "[required] Domain to resolve addresses for (can be specified multiple times)")
	flag.Parse()

	if len(swarmDomains) == 0 {
		flag.Usage()
		fmt.Fprintf(os.Stderr, "No domains given. Aborting.")
		os.Exit(1)
	}

	logger = log.New(os.Stderr, "", 0)

	logger.Printf("Using domains: %v", swarmDomains)

	client, err = swarm.NewClient()
	if err != nil {
		panic(err)
	}

	refreshNodeIPs()

	// Get IPs on every interval
	ticker := time.NewTicker(time.Second * NodeRefreshInterval)
	go func() {
		for range ticker.C {
			refreshNodeIPs()
		}
	}()

	dns.HandleFunc(".", handleRequest)

	go func() {
		srv := &dns.Server{Addr: ":53", Net: "udp"}
		err := srv.ListenAndServe()
		if err != nil {
			logger.Fatalf("Failed to set udp listener %s\n", err.Error())
		}
	}()

	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	for {
		select {
		case s := <-sig:
			logger.Fatalf("Signal (%d) received, stopping\n", s)
			ticker.Stop()
		}
	}
}

func matchingDomain(domain string) *string {
	normalizedDomain := strings.ToLower(domain)

	for _, name := range swarmDomains {
		if (normalizedDomain == name+".") || strings.HasSuffix(normalizedDomain, "."+name+".") {
			return &name
		}
	}
	return nil
}

func handleRequest(w dns.ResponseWriter, r *dns.Msg) {
	m := new(dns.Msg)
	m.SetReply(r)
	m.Authoritative = true

	domain := r.Question[0].Name
	swarmDomain := matchingDomain(domain)

	if swarmDomain != nil && r.Question[0].Qtype == 1 { // Only answer questions for A records on supported domains
		if logflag {
			ip, _, _ := net.SplitHostPort(w.RemoteAddr().String())
			logger.Printf("Request: %15s %s", ip, domain)
		}

		mutex.Lock()
		var rrs = make([]dns.RR, len(ips))
		for i, ip := range ips {
			rr := new(dns.A)
			rr.Hdr = dns.RR_Header{Name: domain, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: uint32(TTL)}
			rr.A = net.ParseIP(ip)
			rrs[i] = rr
		}
		mutex.Unlock()

		m.Answer = shuffleRRs(rrs)
	} else {
		m.Answer = []dns.RR{}
	}

	w.WriteMsg(m)
}

func shuffleRRs(src []dns.RR) []dns.RR {
	dest := make([]dns.RR, len(src))
	perm := rand.Perm(len(src))
	for i, v := range perm {
		dest[v] = src[i]
	}
	return dest
}

func refreshNodeIPs() {
	var err error

	mutex.Lock()
	ips, err = client.ListActiveNodeIPs()
	logger.Printf("Refreshed node IPs: %v\n", ips)
	mutex.Unlock()
	if err != nil {
		panic(err)
	}
}
