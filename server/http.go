package server

import (
	"log"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/millken/go-ipset"
)

//http://www.codegist.net/code/long2ip/
func Long2IP(long uint32) net.IP {
	return net.IPv4(byte(long>>24), byte(long>>16), byte(long>>8), byte(long))
}

func (self *Server) ipf(c *gin.Context) {
	v := c.DefaultPostForm("v", "1.0")
	i := c.DefaultPostForm("i", "0")
	s := c.DefaultPostForm("s", "")
	switch v {
	case "1.0":
		ipuint, err := strconv.ParseUint(i, 10, 32)
		if err != nil {
			c.JSON(200, gin.H{"errcode": 2000, "status": false})
			return
		}
		if s == "" {
			c.JSON(200, gin.H{"errcode": 2001, "status": false})
			return

		}
		ip := Long2IP(uint32(ipuint))
		ipset.Add(s, ip.String())
		c.JSON(200, gin.H{"errcode": -1, "status": true})

	}
	return
}

func (self *Server) serveHttpAPI(addr string, stopC <-chan struct{}) {
	u, err := url.Parse(addr)
	if err != nil {
		log.Fatalf("http api parameter error : %s", err)
	}

	httpS := gin.Default()
	httpS.POST("/ipf", self.ipf)

	httpSs := &http.Server{
		Addr:           u.Host,
		Handler:        httpS,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   5 * time.Second,
		MaxHeaderBytes: 1 << 10,
	}
	go func() {
		err := httpSs.ListenAndServe()
		if err != nil {
			log.Fatalf("failed to start the http server: %v", err)
		}
	}()
	<-stopC
	log.Printf("http api server exit\n")
}
