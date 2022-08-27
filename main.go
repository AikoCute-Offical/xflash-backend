package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path"
	"runtime"
	"strings"
	"syscall"

	"github.com/AikoCute-Offical/xflash-backend/api/panel"
	"github.com/AikoCute-Offical/xflash-backend/conf"
	"github.com/AikoCute-Offical/xflash-backend/core"
	"github.com/AikoCute-Offical/xflash-backend/node"
	"github.com/spf13/viper"
)

var (
	configFile   = flag.String("config", "", "Config file for XrayR.")
	printVersion = flag.Bool("version", false, "show version")
)

var (
	version  = "0.0.2"
	codename = "xflash"
	intro    = "Xflashx backend based on Xray-core"
)

func showVersion() {
	fmt.Printf("%s %s (%s) \n", codename, version, intro)
}

func getConfig() *viper.Viper {
	config := viper.New()
	// Set custom path and name
	if *configFile != "" {
		configName := path.Base(*configFile)
		configFileExt := path.Ext(*configFile)
		configNameOnly := strings.TrimSuffix(configName, configFileExt)
		configPath := path.Dir(*configFile)
		config.SetConfigName(configNameOnly)
		config.SetConfigType(strings.TrimPrefix(configFileExt, "."))
		config.AddConfigPath(configPath)
		// Set ASSET Path and Config Path for XrayR
		os.Setenv("XRAY_LOCATION_ASSET", configPath)
		os.Setenv("XRAY_LOCATION_CONFIG", configPath)
	} else {
		// Set default config path
		config.SetConfigName("config")
		config.SetConfigType("yml")
		config.AddConfigPath(".")
	}
	if err := config.ReadInConfig(); err != nil {
		log.Panicf("Fatal error config file: %s \n", err)
	}
	return config
}

func startNodes(nodes []*conf.NodeConfig, core *core.Core) error {
	for i, _ := range nodes {
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
	config := getConfig()
	c := conf.New()
	err := config.Unmarshal(c)
	if err != nil {
		log.Panicf("can't unmarshal config file: %s \n", err)
	}
	x := core.New(c)
	x.Start()
	defer x.Close()
	err = startNodes(c.NodesConfig, x)
	if err != nil {
		log.Panicf("run nodes error: %v", err)
	}
	//Explicitly triggering GC to remove garbage from config loading.
	runtime.GC()
	// Running backend
	{
		osSignals := make(chan os.Signal, 1)
		signal.Notify(osSignals, os.Interrupt, os.Kill, syscall.SIGTERM)
		<-osSignals
	}
}
