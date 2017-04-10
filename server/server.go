package server

import (
	"log"
	"net/http"
	"sync"
)

type Server struct {
	mutex  sync.Mutex
	conf   ServerConfig
	stopC  chan struct{}
	wg     sync.WaitGroup
	router http.Handler
}

func NewServer(conf ServerConfig) *Server {
	s := &Server{
		conf:  conf,
		stopC: make(chan struct{}),
	}
	return s
}

func (self *Server) Stop() {
	close(self.stopC)
	self.wg.Wait()
	log.Printf("server stopped")
}

func (self *Server) ServeAPI() {
	self.wg.Add(2)
	go func() {
		defer self.wg.Done()
		self.serverCmdAPI(self.conf.CmdApiAddr, self.stopC)
	}()
	go func() {
		defer self.wg.Done()
		self.serveHttpAPI(self.conf.HttpApiAddr, self.stopC)
	}()
}
