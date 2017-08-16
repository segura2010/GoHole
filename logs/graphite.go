package logs

import (
    "github.com/marpaia/graphite-golang"

    "GoHole/config"
)

func getGraphiteInstance() *graphite.Graphite {
	host := config.GetInstance().Graphite.Host
	port := config.GetInstance().Graphite.Port
	instance, err := graphite.NewGraphite(host, port)

	if err != nil {
		instance = graphite.NewGraphiteNop(host, port)
	}

    return instance
}

func AddQueryToGraphite(isBlocked, isIpv4, isCached bool){
	// We add just "1" as value, the user must configure the graph to "summarize" (sum)
	// the metrics in order to see the graphs :)

	Graphite := getGraphiteInstance()

	// add query to total query metric
	Graphite.SimpleSend("gohole.queries.total", "1")

	// add query to blocked/non-blocked query metric
	if isBlocked {
		Graphite.SimpleSend("gohole.queries.blocked", "1")
	}else{
		Graphite.SimpleSend("gohole.queries.nonblocked", "1")
	}

	// add query to ipv4/ipv6 query metric
	if isIpv4 {
		Graphite.SimpleSend("gohole.queries.ipv4", "1")
	}else{
		Graphite.SimpleSend("gohole.queries.ipv6", "1")
	}

	// add query to cached/non-cached query metric
	if isCached {
		Graphite.SimpleSend("gohole.queries.cached", "1")
	}else{
		Graphite.SimpleSend("gohole.queries.noncached", "1")
	}

	Graphite.Disconnect()
}

