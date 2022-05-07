package main

import (
	"math/rand"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
)

// easy dispatcher for testing multiple leaderboard server
func main() {
	server := gin.Default()

	dispatcher := &dispatcher{
		servers: []string{"http://server1", "http://server2"},
	}

	server.POST("/api/v1/score", dispatcher.proxy)
	server.GET("/api/v1/leaderboard", dispatcher.proxy)

	server.Run(":8000")
}

type dispatcher struct {
	servers []string
}

func (load dispatcher) proxy(ctx *gin.Context) {
	p, err := url.Parse(load.servers[rand.Intn(len(load.servers))])
	if err != nil {
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(p)
	proxy.ServeHTTP(ctx.Writer, ctx.Request)
}
