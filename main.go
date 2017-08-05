package main

import (
    "log"
    "flag"

    "GoHole/config"
    "GoHole/dnsserver"
)


func main(){

    // Command line options
    port := flag.String("p", "", "Set DNS server port")
    cfgFile := flag.String("c", "./config.json", "Config file")
    flag.Parse()

    log.Printf("Loading config..")
    config.CreateInstance(*cfgFile)
    if *port != ""{
        config.GetInstance().DNSPort = *port
    }

    dnsserver.ListenAndServe()

}
