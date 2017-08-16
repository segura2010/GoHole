package config

import (
    "encoding/json"
    "io/ioutil"
    "log"
)

// MyConfig struct
// This is the struct that the config.json must have
type MyConfig struct {
    ServerIP string // the DNS server IP to redirect blocked ads
    DNSPort string // listen on port

    // RedisDB info
    RedisDB RedisConfig

    // Graphite info
    Graphite GraphiteConfig

    UpstreamDNSServer string
    DomainCacheTime int // time to save domains in cache (in seconds)
}

// DB Config
type RedisConfig struct {
    Host string
    Port string
    Pass string
}

// Graphite Config
type GraphiteConfig struct {
    Host string
    Port int
}

var instance *MyConfig = nil

func CreateInstance(filename string) *MyConfig {
    var err error
    instance, err = loadConfig(filename)
    if err != nil {
        log.Printf("Error loading config file: %s\nUsing default config.", err)
        // use defaults
        instance = &MyConfig{
            ServerIP: "0.0.0.0",
            DNSPort: "53",
            UpstreamDNSServer: "8.8.8.8",
            DomainCacheTime: 1800,
            RedisDB: RedisConfig{
                Host: "localhost",
                Port: "6379",
                Pass: "",
            },
            Graphite: GraphiteConfig{
                Host: "localhost",
                Port: 2003,
            },
        }
    }

    return instance
}

func GetInstance() *MyConfig {
    return instance
}

func loadConfig(filename string) (*MyConfig, error){
    var s *MyConfig

    bytes, err := ioutil.ReadFile(filename)
    if err != nil {
        return s, err
    }
    // Unmarshal json
    err = json.Unmarshal(bytes, &s)
    return s, err
}
