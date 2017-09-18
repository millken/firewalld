package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"syscall"

	"github.com/judwhite/go-svc/svc"
	"github.com/millken/firewalld/server"
	"github.com/millken/firewalld/worker"
)

const Binary = "0.0.2"

var (
	flagSet        = flag.NewFlagSet("firewalld", flag.ExitOnError)
	configFilePath = flagSet.String("config", "", "the config file path to read")
	showVersion    = flagSet.Bool("version", false, "print version string and exit")
	showIpset      = flagSet.Bool("ipset", false, "print ipset&iptables string and exit")
)

type program struct {
	server *server.Server
}

func main() {
	defer log.Printf("main exit")
	prg := &program{}
	if err := svc.Run(prg, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGINT); err != nil {
		log.Fatal(err)
	}
}

func (p *program) Init(env svc.Environment) error {
	if env.IsWindowsService() {
		dir := filepath.Dir(os.Args[0])
		return os.Chdir(dir)
	}
	return nil
}

func (p *program) Start() error {
	flagSet.Parse(os.Args[1:])

	if *showVersion {
		fmt.Println(fmt.Sprintf("firewalld v%s (built w/%s)", Binary, runtime.Version()))
		os.Exit(0)
	}
	if *showIpset {
		cmd := `ipset destroy dbi
ipset destroy dbn
ipset destroy sbi
ipset destroy sbn
ipset destroy dwi
ipset destroy dwn
ipset destroy swi
ipset destroy swn
ipset create -exist dbi hash:ip hashsize 4096 maxelem 1048576 timeout 86400
ipset create -exist dbn hash:net timeout 86400
ipset create -exist sbi hash:ip hashsize 4096 maxelem 1048576
ipset create -exist sbn hash:net hashsize 2048 maxelem 524288
ipset create -exist dwi hash:ip hashsize 2048 maxelem 524288 timeout 86400
ipset create -exist dwn hash:net timeout 86400
ipset create -exist swi hash:ip hashsize 4096 maxelem 1048576
ipset create -exist swn hash:net hashsize 2048 maxelem 524288
iptables -A INPUT  -m set --match-set dbi src -p TCP -m multiport --dports 80,12377 -j DROP
iptables -A INPUT  -m set --match-set dbn src -p TCP -m multiport --dports 80,12377 -j DROP
iptables -A INPUT  -m set --match-set dbn sbi -p TCP -m multiport --dports 80,12377 -j REJECT
iptables -A INPUT  -m set --match-set dbn sbn -p TCP -m multiport --dports 80,12377 -j REJECT
iptables -I INPUT  -m set --match-set dwi src -p TCP -m multiport --dports 80,12377 -j ACCEPT
iptables -I INPUT  -m set --match-set dwn src -p TCP -m multiport --dports 80,12377 -j ACCEPT
iptables -I INPUT  -m set --match-set swi src -p TCP -m multiport --dports 80,12377 -j ACCEPT
iptables -I INPUT  -m set --match-set swn src -p TCP -m multiport --dports 80,12377 -j ACCEPT`
		fmt.Println(cmd)
		os.Exit(0)
	}
	var serverConf server.ServerConfig
	if *configFilePath != "" {
		d, err := ioutil.ReadFile(*configFilePath)
		if err != nil {
			panic(err)
		}
		err = json.Unmarshal(d, &serverConf)
		if err != nil {
			panic(err)
		}
	}

	loadConf, _ := json.MarshalIndent(serverConf, "", " ")
	fmt.Printf("loading with conf:%v\n", string(loadConf))
	app := server.NewServer(serverConf)
	go worker.Run()
	app.ServeAPI()
	p.server = app
	return nil
}

func (p *program) Stop() error {
	if p.server != nil {
		p.server.Stop()
	}
	return nil
}
