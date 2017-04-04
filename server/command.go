package server

import (
	"errors"
	"log"
	"net/url"
	"runtime"
	"strings"

	"github.com/millken/go-ipset"
	"github.com/tidwall/redcon"
)

var (
	errInvalidCommand = errors.New("invalid command")
)

func (self *Server) serverCmd(conn redcon.Conn, cmd redcon.Command) {
	var args []string
	defer func() {
		if e := recover(); e != nil {
			buf := make([]byte, 4096)
			n := runtime.Stack(buf, false)
			buf = buf[0:n]
			log.Printf("handle cmd command panic: %s:%v", buf, e)
		}
	}()
	for _, arg := range cmd.Args {
		args = append(args, strings.ToLower(string(arg)))
	}
	cmdName := strings.ToLower(string(cmd.Args[0]))
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
	case "add":
		if len(cmd.Args) < 3 {
			conn.WriteString("Err wrong number of argument")
		} else {
			ipset.Add(args[1], args[2], args[3:]...)
			conn.WriteString("OK")
		}
	case "del":
		if len(cmd.Args) < 3 {
			conn.WriteString("Err wrong number of argument")
		} else {
			ipset.Del(args[1], args[2], args[3:]...)
			conn.WriteString("OK")
		}
	case "quit":
		conn.WriteString("OK")
		conn.Close()
	default:
		conn.WriteError("ERR handle command '" + string(cmd.Args[0]))
	}
}

func (self *Server) serverCmdAPI(addr string, stopC <-chan struct{}) {
	u, err := url.Parse(addr)
	if err != nil {
		log.Fatalf("cmd api parameter error : %s", err)
	}
	redisS := redcon.NewServerNetwork(
		u.Scheme,
		u.Host+u.Path,
		self.serverCmd,
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
	log.Printf("cmd api server exit\n")
}
