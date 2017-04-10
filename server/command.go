package server

import (
	"errors"
	"log"
	"net/url"
	"runtime"
	"strings"

	"github.com/millken/go-ipset"
	"github.com/millken/raphanus"
	"github.com/tidwall/redcon"
)

const (
	DynamicWhitelistIp  = "dwi"
	DynamicWhitelistNet = "dwn"
	DynamicBlacklistIp  = "dbi"
	DynamicBlacklistNet = "dbn"
	StaticWhitelistIp   = "swi"
	StaticWhitelistNet  = "swn"
	StaticBlacklistIp   = "sbi"
	StaticBlacklistNet  = "sbn"
)

var (
	errInvalidCommand = errors.New("invalid command")
	db                = raphanus.New(256)
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
	case "add" + DynamicWhitelistIp:
		if len(cmd.Args) != 3 {
			conn.WriteString("Err wrong number of argument")
		} else {
			ipset.Add(DynamicWhitelistIp, args[1], "timeout", args[2])
			conn.WriteString("OK")
		}
	case "del" + DynamicWhitelistIp:
		if len(cmd.Args) != 2 {
			conn.WriteString("Err wrong number of argument")
		} else {
			ipset.Del(DynamicWhitelistIp, args[1])
			conn.WriteString("OK")
		}
	case "add" + DynamicWhitelistNet:
		if len(cmd.Args) != 3 {
			conn.WriteString("Err wrong number of argument")
		} else {
			ipset.Add(DynamicWhitelistNet, args[1], "timeout", args[2])
			conn.WriteString("OK")
		}
	case "del" + DynamicWhitelistNet:
		if len(cmd.Args) != 2 {
			conn.WriteString("Err wrong number of argument")
		} else {
			ipset.Del(DynamicWhitelistNet, args[1])
			conn.WriteString("OK")
		}
	case "add" + DynamicBlacklistIp:
		if len(cmd.Args) != 3 {
			conn.WriteString("Err wrong number of argument")
		} else {
			ipset.Add(DynamicBlacklistIp, args[1], "timeout", args[2])
			conn.WriteString("OK")
		}
	case "del" + DynamicBlacklistIp:
		if len(cmd.Args) != 2 {
			conn.WriteString("Err wrong number of argument")
		} else {
			ipset.Del(DynamicBlacklistIp, args[1])
			conn.WriteString("OK")
		}
	case "add" + DynamicBlacklistNet:
		if len(cmd.Args) != 3 {
			conn.WriteString("Err wrong number of argument")
		} else {
			ipset.Add(DynamicBlacklistNet, args[1], "timeout", args[2])
			conn.WriteString("OK")
		}
	case "del" + DynamicBlacklistNet:
		if len(cmd.Args) != 2 {
			conn.WriteString("Err wrong number of argument")
		} else {
			ipset.Del(DynamicBlacklistNet, args[1])
			conn.WriteString("OK")
		}
	case "add" + StaticWhitelistIp:
		if len(cmd.Args) != 2 {
			conn.WriteString("Err wrong number of argument")
		} else {
			ipset.Add(StaticWhitelistIp, args[1])
			conn.WriteString("OK")
		}
	case "del" + StaticWhitelistIp:
		if len(cmd.Args) != 2 {
			conn.WriteString("Err wrong number of argument")
		} else {
			ipset.Del(StaticWhitelistIp, args[1])
			conn.WriteString("OK")
		}
	case "add" + StaticWhitelistNet:
		if len(cmd.Args) != 2 {
			conn.WriteString("Err wrong number of argument")
		} else {
			ipset.Add(StaticWhitelistNet, args[1])
			conn.WriteString("OK")
		}
	case "del" + StaticWhitelistNet:
		if len(cmd.Args) != 2 {
			conn.WriteString("Err wrong number of argument")
		} else {
			ipset.Del(StaticWhitelistNet, args[1])
			conn.WriteString("OK")
		}
	case "add" + StaticBlacklistIp:
		if len(cmd.Args) != 2 {
			conn.WriteString("Err wrong number of argument")
		} else {
			ipset.Add(StaticBlacklistIp, args[1])
			conn.WriteString("OK")
		}
	case "del" + StaticBlacklistIp:
		if len(cmd.Args) != 2 {
			conn.WriteString("Err wrong number of argument")
		} else {
			ipset.Del(StaticBlacklistIp, args[1])
			conn.WriteString("OK")
		}
	case "add" + StaticBlacklistNet:
		if len(cmd.Args) != 2 {
			conn.WriteString("Err wrong number of argument")
		} else {
			ipset.Add(StaticBlacklistNet, args[1])
			conn.WriteString("OK")
		}
	case "del" + StaticBlacklistNet:
		if len(cmd.Args) != 2 {
			conn.WriteString("Err wrong number of argument")
		} else {
			ipset.Del(StaticBlacklistNet, args[1])
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
