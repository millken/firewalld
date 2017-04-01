package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"syscall"
	"time"

	"github.com/judwhite/go-svc/svc"
	"github.com/millken/firewalld/common"
	"github.com/millken/firewalld/server"
)

var (
	flagSet        = flag.NewFlagSet("zanredisdb", flag.ExitOnError)
	configFilePath = flagSet.String("config", "", "the config file path to read")
	showVersion    = flagSet.Bool("version", false, "print version string and exit")
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
		fmt.Println(common.String("firewalld"))
		os.Exit(0)
	}
	var configFile server.ConfigFile
	if *configFilePath != "" {
		d, err := ioutil.ReadFile(*configFilePath)
		if err != nil {
			panic(err)
		}
		err = json.Unmarshal(d, &configFile)
		if err != nil {
			panic(err)
		}
	}
	if configFile.ServerConf.DataDir == "" {
		tmpDir, err := ioutil.TempDir("", fmt.Sprintf("rocksdb-test-%d", time.Now().UnixNano()))
		if err != nil {
			panic(err)
		}
		configFile.ServerConf.DataDir = tmpDir
	}

	serverConf := configFile.ServerConf

	loadConf, _ := json.MarshalIndent(configFile, "", " ")
	fmt.Printf("loading with conf:%v\n", string(loadConf))
	bip := server.GetIPv4ForInterfaceName(serverConf.BroadcastInterface)
	if bip == "" || bip == "0.0.0.0" {
		panic("broadcast ip can not be found")
	} else {
		serverConf.BroadcastAddr = bip
	}
	fmt.Printf("broadcast ip is :%v\n", bip)
	app := server.NewServer(serverConf)
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
