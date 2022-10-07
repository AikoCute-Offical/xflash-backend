package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/AikoCute-Offical/xflash-backend/api/panel"
	"github.com/AikoCute-Offical/xflash-backend/conf"
	"github.com/AikoCute-Offical/xflash-backend/core"
	"github.com/AikoCute-Offical/xflash-backend/node"
)

var (
	configFile   = flag.String("config", "/etc/Xflash/config.yml", "Config file for Xflash.")
	printVersion = flag.Bool("version", false, "show version")
)

var (
	codename = "xflash"
	intro    = "Xflashx backend based on Xray-core"
	version  = "v0.0.7_beta"
	codename = "Xflash"
	intro    = "A V2board backend based on Xray-core"
)

func showVersion() {
	fmt.Printf("%s %s (%s) \n", codename, version, intro)
}

func getConfig() *viper.Viper {
func startNodes(nodes []*conf.NodeConfig, core *core.Core) error {
	for i := range nodes {
		var apiClient = panel.New(nodes[i].ApiConfig)
		// Register controller service
		err := node.New(core, apiClient, nodes[i].ControllerConfig).Start()
		if err != nil {
			return fmt.Errorf("start node controller error: %v", err)
		}
	}
	return nil
}

func main() {
	flag.Parse()
	showVersion()
	if *printVersion {
		return
	}
	config := conf.New()
	err := config.LoadFromPath(*configFile)
	if err != nil {
		log.Panicf("can't unmarshal config file: %s \n", err)
	}
	x := core.New(config)
	x.Start()
	defer x.Close()
	err = startNodes(config.NodesConfig, x)
	if err != nil {
		log.Panicf("run nodes error: %v", err)
	}
	//Explicitly triggering GC to remove garbage from config loading.
	runtime.GC()
	// Running backend
	{
		osSignals := make(chan os.Signal, 1)
		signal.Notify(osSignals, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)
		<-osSignals
	}
}
