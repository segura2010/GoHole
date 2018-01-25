package dnsserver

import (
    "fmt"
	"log"
	"net"
	"strings"
	"time"

    "github.com/miekg/dns"

    "GoHole/config"
    "GoHole/dnscache"
    "GoHole/logs"
    "GoHole/encryption"
)

func parseQuery(clientIp string, m *dns.Msg) {
	for _, q := range m.Question {

		log.Printf("Query for %s from %s", q.Name, clientIp)

		var err error = nil
		var ip string = ""
		cleanedName := q.Name[0:len(q.Name)-1] // remove the end "."
		qType := "A"
		cached := 0
		var cacheTTL int64 = 1
		isIpv4 := true

		if q.Qtype == dns.TypeA{
			ip, err = dnscache.GetDomainIPv4(cleanedName)
			cacheTTL, _ = dnscache.GetTTLDomainIPv4(cleanedName)
		}else if q.Qtype == dns.TypeAAAA{
			ip, err = dnscache.GetDomainIPv6(cleanedName)
			cacheTTL, _ = dnscache.GetTTLDomainIPv6(cleanedName)
			qType = "AAAA"
			isIpv4 = false
		}

		if ip != "" && err == nil {
			rr, err := dns.NewRR(fmt.Sprintf("%s %s %s", q.Name, qType, ip))
			if err == nil {
				m.Answer = append(m.Answer, rr)
			}
			cached = 1
		}else{
			// Request to a DNS server
			c := new(dns.Client)
			msg := new(dns.Msg)
			msg.SetQuestion(dns.Fqdn(q.Name), q.Qtype)
			msg.RecursionDesired = true

		    r, _, err := c.Exchange(msg, net.JoinHostPort(config.GetInstance().UpstreamDNSServer, "53"))
		    if r == nil {
		    	log.Printf("*** error: %s\n", err.Error())
		    	return
		    }

		    if r.Rcode != dns.RcodeSuccess {
		    	log.Printf(" *** invalid answer name %s after %s query for %s\n", q.Name, qType, q.Name)
		    	return
		    }
		    // Parse Answer
		    for _, a := range r.Answer {
		    	ans := strings.Split(a.String(), "\t")
		    	if len(ans) == 5 && ans[3] == qType{
		    		// Save on cache
		    		if q.Qtype == dns.TypeA{
		    			dnscache.AddDomainIPv4(cleanedName, ans[4], config.GetInstance().DomainCacheTime)
					}else if q.Qtype == dns.TypeAAAA{
						dnscache.AddDomainIPv6(cleanedName, ans[4], config.GetInstance().DomainCacheTime)
					}
		    	}
		    }
		    // Set answer for the client
		    m.Answer = r.Answer
		    cached = 0
		}

		// Add logs
		now := time.Now().Unix()
		logs.AddQuery(clientIp, cleanedName, cached, now)
		
		isBlocked := false
		isCached := true
		if cached == 0{
			isCached = false
		}else{
			// if it was cached and TTL is -1 (<0), then it is a blocked domain
			if cacheTTL < 0{
				// is a blocked domain
				isBlocked = true
			}
		}
		go logs.AddQueryToGraphite(isBlocked, isIpv4, isCached)
	}
}

func handleDnsRequest(w dns.ResponseWriter, r *dns.Msg) {
	m := new(dns.Msg)
	m.SetReply(r)
	m.Compress = false
	clientIp := w.RemoteAddr().String()
	clientIp = clientIp[0:strings.LastIndex(clientIp, ":")] // remove port

	switch r.Opcode {
	case dns.OpcodeQuery:
		parseQuery(clientIp, m)
	}

	w.WriteMsg(m)
}

func listenAndServeSecure(){
	serveraddr, err := net.ResolveUDPAddr("udp",":"+ config.GetInstance().SecureDNSPort)
	if err != nil {
		log.Fatal("Failed to start DNS Secure Server: %s\n", err)
	}
	conn, err := net.ListenUDP("udp", serveraddr)
	if err != nil {
		log.Fatal("Failed to start DNS Secure Server: %s\n", err)
	}
	defer conn.Close()

	log.Printf("Starting Secure DNS Server at %s\n", config.GetInstance().SecureDNSPort)

	//simple read
	for{
		buf := make([]byte, 2048)
		n, addr, err := conn.ReadFromUDP(buf)
        if err != nil {
            continue
        }

        query, err := encryption.Decrypt(buf[:n])
        if err != nil{
        	continue
        }

        m := new(dns.Msg)
        m.Unpack(query)
        clientIp := addr.String()
        clientIp = clientIp[0:strings.LastIndex(clientIp, ":")] // remove port
        parseQuery(clientIp, m)

        reply, err := m.Pack()
        if err != nil{
        	continue
        }
        eReply, err := encryption.Encrypt(reply)
        if err != nil{
        	continue
        }

        conn.WriteToUDP(eReply, addr)
	}
}

func ListenAndServe(){

	// add go.hole domain to our cache :)
	dnscache.AddDomainIPv4("go.hole", config.GetInstance().ServerIP, 0)

	// start the graphite statistics loop
	go logs.StartStatsLoop()

	dns.HandleFunc(".", handleDnsRequest)

	// Start DNS server
	port := config.GetInstance().DNSPort
	server := &dns.Server{Addr: ":" + port, Net: "udp"}

	log.Printf("Starting at %s\n", port)
	go listenAndServeSecure()

	err := server.ListenAndServe()
	defer server.Shutdown()
	if err != nil {
		log.Fatalf("Failed to start DNS Server: %s\n ", err.Error())
	}

}