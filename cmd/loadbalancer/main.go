package main

import (
	"math/rand"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
)

// easy loadbalancer for testing multiple leaderboard server
func main() {
	server := gin.Default()

	loadbalancer := &loadbalancer{
		servers: []string{"http://server1", "http://server2"},
	}

	server.POST("/api/score", loadbalancer.proxy)
	server.GET("/api/leaderboard", loadbalancer.proxy)

	server.Run(":8000")
}

type loadbalancer struct {
	servers []string
}

func (load loadbalancer) proxy(ctx *gin.Context) {
	p, err := url.Parse(load.servers[rand.Intn(len(load.servers))])
	if err != nil {
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(p)
	proxy.ServeHTTP(ctx.Writer, ctx.Request)
}
