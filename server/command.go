package server

import (
	"errors"
	"log"
	"runtime"
	"strconv"

	"github.com/tidwall/redcon"
)

var (
	errInvalidCommand = errors.New("invalid command")
)

func (self *Server) serverRedis(conn redcon.Conn, cmd redcon.Command) {
	defer func() {
		if e := recover(); e != nil {
			buf := make([]byte, 4096)
			n := runtime.Stack(buf, false)
			buf = buf[0:n]
			log.Printf("handle redis command panic: %s:%v", buf, e)
		}
	}()

	_, cmd, err := pipelineCommand(conn, cmd)
	if err != nil {
		conn.WriteError("pipeline error '" + err.Error() + "'")
		return
	}
	cmdName := qcmdlower(cmd.Args[0])
	switch cmdName {
	case "detach":
		hconn := conn.Detach()
		log.Printf("connection has been detached")
		go func() {
			defer hconn.Close()
			hconn.WriteString("OK")
			hconn.Flush()
		}()
	case "ping":
		conn.WriteString("PONG")
	case "quit":
		conn.WriteString("OK")
		conn.Close()
	default:
		h, cmd, err := self.GetHandler(cmdName, cmd)
		if err == nil {
			h(conn, cmd)
		} else {
			conn.WriteError("ERR handle command '" + string(cmd.Args[0]) + "' : " + err.Error())
		}
	}
}

func (self *Server) serveRedisAPI(port int, stopC <-chan struct{}) {
	redisS := redcon.NewServer(
		":"+strconv.Itoa(port),
		self.serverRedis,
		func(conn redcon.Conn) bool {
			//log.Printf("accept: %s", conn.RemoteAddr())
			return true
		},
		func(conn redcon.Conn, err error) {
			if err != nil {
				log.Printf("closed: %s, err: %v", conn.RemoteAddr(), err)
			}
		},
	)
	go func() {
		err := redisS.ListenAndServe()
		if err != nil {
			log.Fatalf("failed to start the redis server: %v", err)
		}
	}()
	<-stopC
	redisS.Close()
	log.Printf("redis api server exit\n")
}
