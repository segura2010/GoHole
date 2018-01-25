package logs

import (
	"strconv"
	"time"

    "github.com/marpaia/graphite-golang"

    "GoHole/config"
)

type Statistics struct {
    Total int
    Blocked int
    NonBlocked int
    Cached int
    NonCached int
    Ipv4 int
    Ipv6 int
}

var statsInstance *Statistics = nil

func getGraphiteInstance() *graphite.Graphite {
	host := config.GetInstance().Graphite.Host
	port := config.GetInstance().Graphite.Port
	instance, err := graphite.NewGraphite(host, port)

	if err != nil {
		//instance = graphite.NewGraphiteNop(host, port)
		instance = nil
	}

    return instance
}

func getStatsInstance() *Statistics {
	if statsInstance == nil{
		statsInstance = &Statistics{
			Total:0,
			Blocked:0,
			NonBlocked:0,
			Cached:0,
			NonCached:0,
			Ipv4:0,
			Ipv6:0,
		}
	}

    return statsInstance
}

func AddQueryToGraphite(isBlocked, isIpv4, isCached bool){
	stats := getStatsInstance()
	stats.Total += 1
	// add query to blocked/non-blocked query metric
	if isBlocked {
		stats.Blocked += 1
	}else{
		stats.NonBlocked += 1
	}

	// add query to ipv4/ipv6 query metric
	if isIpv4 {
		stats.Ipv4 += 1
	}else{
		stats.Ipv6 += 1
	}

	// add query to cached/non-cached query metric
	if isCached {
		stats.Cached += 1
	}else{
		stats.NonCached += 1
	}
}

func resetStats(){
	stats := getStatsInstance()
	stats.Total = 0
	stats.Blocked = 0
	stats.NonBlocked = 0
	stats.Cached = 0
	stats.NonCached = 0
	stats.Ipv4 = 0
	stats.Ipv6 = 0
}

func sendQueriesToGraphite(){
	// The user should configure the graph to "summarize" (sum)
	// the metrics in order to see better graphs :)

	stats := getStatsInstance()
	Graphite := getGraphiteInstance()
	if Graphite == nil{
		return
	}

	// add query to total query metric
	Graphite.SimpleSend("gohole.queries.total", strconv.Itoa(stats.Total))

	// add query to blocked/non-blocked query metric
	Graphite.SimpleSend("gohole.queries.blocked", strconv.Itoa(stats.Blocked))
	Graphite.SimpleSend("gohole.queries.nonblocked", strconv.Itoa(stats.NonBlocked))

	// add query to ipv4/ipv6 query metrics
	Graphite.SimpleSend("gohole.queries.ipv4", strconv.Itoa(stats.Ipv4))
	Graphite.SimpleSend("gohole.queries.ipv6", strconv.Itoa(stats.Ipv6))

	// add query to cached/non-cached query metric
	Graphite.SimpleSend("gohole.queries.cached", strconv.Itoa(stats.Cached))
	Graphite.SimpleSend("gohole.queries.noncached", strconv.Itoa(stats.NonCached))

	resetStats()
	Graphite.Disconnect()
}

func StartStatsLoop(){
	// loop in which every 30s we send the stats to Graphite
	for{
		time.Sleep(time.Duration(30) * time.Second)
		sendQueriesToGraphite()
	}
}

