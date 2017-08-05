package config

import (
    "encoding/json"
    "io/ioutil"
)

// MyConfig struct
// This is the struct that the config.json must have
type MyConfig struct {
    ServerIP string // the DNS server IP to redirect blocked ads
    DNSPort string // listen on port

    // RedisDB info
    RedisDB RedisConfig

    UpstreamDNSServer string
    DomainCacheTime int // time to save domains in cache (in seconds)
}

// DB Config
type RedisConfig struct {
    Host string
    Port string
    Pass string
}

var instance *MyConfig = nil

func CreateInstance(filename string) *MyConfig {
    var err error
    instance, err = loadConfig(filename)
    if err != nil {
        // use defaults
        instance = &MyConfig{
            ServerIP: "0.0.0.0",
            DNSPort: "5353",
            UpstreamDNSServer: "8.8.8.8",
            DomainCacheTime: 300,
            RedisDB: RedisConfig{
                Host: "localhost",
                Port: "6379",
                Pass: "",
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
