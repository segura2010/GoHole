package dnsserver

import (
    "fmt"
	"log"
	"net"
	"strings"

    "github.com/miekg/dns"

    "GoHole/config"
    "GoHole/dnscache"
)

func parseQuery(m *dns.Msg) {
	for _, q := range m.Question {
		switch q.Qtype {
		case dns.TypeA:
			// IPv4
			log.Printf("Query for %s\n", q.Name)
			cleanedName := q.Name[0:len(q.Name)-1] // remove the end "."
			ip, err := dnscache.GetDomainIPv4(cleanedName)
			log.Printf("IP: %s", ip)
			if ip != "" && err == nil {
				rr, err := dns.NewRR(fmt.Sprintf("%s A %s", q.Name, ip))
				if err == nil {
					m.Answer = append(m.Answer, rr)
				}
			}else{
				// Request to a DNS server
				c := new(dns.Client)
				msg := new(dns.Msg)
				msg.SetQuestion(dns.Fqdn(q.Name), dns.TypeA)
				msg.RecursionDesired = true

			    r, _, err := c.Exchange(msg, net.JoinHostPort(config.GetInstance().UpstreamDNSServer, "53"))
			    if r == nil {
			    	log.Printf("*** error: %s\n", err.Error())
			    	return
			    }

			    if r.Rcode != dns.RcodeSuccess {
			    	log.Printf(" *** invalid answer name %s after A query for %s\n", q.Name, q.Name)
			    	return
			    }
			    // Parse Answer
			    for _, a := range r.Answer {
			    	ans := strings.Split(a.String(), "\t")
			    	if len(ans) == 5 && ans[3] == "A"{
			    		// Save on cache
			    		dnscache.AddDomainIPv4(cleanedName, ans[4], config.GetInstance().DomainCacheTime)
			    	}
			    }
			    // Set answer for the client
			    m.Answer = r.Answer
			}

		case dns.TypeAAAA:
			// IPv6
			log.Printf("Query for %s\n", q.Name)
			cleanedName := q.Name[0:len(q.Name)-1] // remove the end "."
			ip, err := dnscache.GetDomainIPv6(cleanedName)
			log.Printf("IP: %s", ip)
			if ip != "" && err == nil {
				rr, err := dns.NewRR(fmt.Sprintf("%s AAAA %s", q.Name, ip))
				if err == nil {
					m.Answer = append(m.Answer, rr)
				}
			}else{
				// Request to a DNS server
				c := new(dns.Client)
				msg := new(dns.Msg)
				msg.SetQuestion(dns.Fqdn(q.Name), dns.TypeAAAA)
				msg.RecursionDesired = true

			    r, _, err := c.Exchange(msg, net.JoinHostPort(config.GetInstance().UpstreamDNSServer, "53"))
			    if r == nil {
			    	log.Printf("*** error: %s\n", err.Error())
			    	return
			    }

			    if r.Rcode != dns.RcodeSuccess {
			    	log.Printf(" *** invalid answer name %s after AAAA query for %s\n", q.Name, q.Name)
			    	return
			    }
			    // Parse Answer
			    for _, a := range r.Answer {
			    	ans := strings.Split(a.String(), "\t")
			    	if len(ans) == 5 && ans[3] == "AAAA"{
			    		// Save on cache
			    		dnscache.AddDomainIPv6(cleanedName, ans[4], config.GetInstance().DomainCacheTime)
			    	}
			    }
			    // Set answer for the client
			    m.Answer = r.Answer
			}
		}
	}
}

func handleDnsRequest(w dns.ResponseWriter, r *dns.Msg) {
	m := new(dns.Msg)
	m.SetReply(r)
	m.Compress = false

	switch r.Opcode {
	case dns.OpcodeQuery:
		parseQuery(m)
	}

	w.WriteMsg(m)
}

func ListenAndServe(){

	dns.HandleFunc(".", handleDnsRequest)

	// Start DNS server
	port := config.GetInstance().DNSPort
	server := &dns.Server{Addr: ":" + port, Net: "udp"}

	log.Printf("Starting at %s\n", port)

	err := server.ListenAndServe()
	defer server.Shutdown()
	if err != nil {
		log.Fatalf("Failed to start DNS Server: %s\n ", err.Error())
	}

}