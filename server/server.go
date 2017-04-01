package server

import (
	"errors"
	"log"
	"net/http"
	"sync"

	"github.com/millken/firewalld/common"
	"github.com/tidwall/redcon"
)

var (
	errNamespaceNotFound = errors.New("namespace not found")
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
	// api server should disable the api request while starting until replay log finished and
	// also while we recovery we need to disable api.
	self.wg.Add(1)
	go func() {
		defer self.wg.Done()
		self.serveRedisAPI(self.conf.RedisAPIPort, self.stopC)
	}()
}

func (self *Server) GetHandler(cmdName string, cmd redcon.Command) (common.CommandFunc, redcon.Command, error) {
	if len(cmd.Args) < 2 {
		return nil, cmd, common.ErrInvalidArgs
	}
	rawKey := cmd.Args[1]

	_, _, err := common.ExtractNamesapce(rawKey)
	if err != nil {
		log.Printf("failed to get the namespace of the redis command:%v", rawKey)
		return nil, cmd, err
	}

	return nil, cmd, errNamespaceNotFound
}
