package main

import (
	"flag"
	"fmt"
	"github.com/miekg/dns"
	"log"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var logflag bool
var ips []string

func main() {
	flag.BoolVar(&logflag, "log", false, "Log requests to stdout")
	flag.Parse()
	ips = flag.Args()
	if len(ips) == 0 {
		log.Fatalf("Please give some IPs %s\n", flag.Args())
	}

	dns.HandleFunc(".", handleRequest)

	go func() {
		srv := &dns.Server{Addr: ":53", Net: "udp"}
		err := srv.ListenAndServe()
		if err != nil {
			log.Fatalf("Failed to set udp listener %s\n", err.Error())
		}
	}()

	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	for {
		select {
		case s := <-sig:
			log.Fatalf("Signal (%d) received, stopping\n", s)
		}
	}
}

func handleRequest(w dns.ResponseWriter, r *dns.Msg) {
	// TODO parameterize
	ttl := 3600

	m := new(dns.Msg)
	m.SetReply(r)
	m.Authoritative = true

	if r.Question[0].Qtype == 1 { // Only answer questions for A records
		domain := r.Question[0].Name
		if logflag {
			t := time.Now()
			ip, _, _ := net.SplitHostPort(w.RemoteAddr().String())
			fmt.Printf("%d-%02d-%02d_%02d:%02d:%02d\t%s\t%s\n", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), ip, domain)
		}

		var rrs = make([]dns.RR, len(ips))

		for i, ip := range ips {
			rr := new(dns.A)
			rr.Hdr = dns.RR_Header{Name: domain, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: uint32(ttl)}
			rr.A = net.ParseIP(ip)
			rrs[i] = rr
		}

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
