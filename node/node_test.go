package node_test

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"testing"

	"github.com/AikoCute-Offical/xflash-backend/api/xflash"
	. "github.com/AikoCute-Offical/xflash-backend/node"

	"github.com/AikoCute-Offical/xflash-backend/api"
	_ "github.com/AikoCute-Offical/xflash-backend/example/distro/all"
	"github.com/xtls/xray-core/core"
	"github.com/xtls/xray-core/infra/conf"
)

func TestController(t *testing.T) {
	serverConfig := &conf.Config{
		Stats:     &conf.StatsConfig{},
		LogConfig: &conf.LogConfig{LogLevel: "debug"},
	}
	policyConfig := &conf.PolicyConfig{}
	policyConfig.Levels = map[uint32]*conf.Policy{0: &conf.Policy{
		StatsUserUplink:   true,
		StatsUserDownlink: true,
	}}
	serverConfig.Policy = policyConfig
	config, _ := serverConfig.Build()

	server, err := core.New(config)
	defer server.Close()
	if err != nil {
		t.Errorf("failed to create instance: %s", err)
	}
	if err = server.Start(); err != nil {
		t.Errorf("Failed to start instance: %s", err)
	}
	certConfig := &CertConfig{
		CertMode:   "http",
		CertDomain: "test.ss.tk",
		Provider:   "alidns",
		Email:      "ss@ss.com",
	}
	controlerconfig := &Config{
		UpdatePeriodic: 5,
		CertConfig:     certConfig,
	}
	apiConfig := &api.Config{
		APIHost:  "http://127.0.0.1:667",
		Key:      "123",
		NodeID:   41,
		NodeType: "V2ray",
	}
	apiclient := xflash.New(apiConfig)
	c := New(server, apiclient, controlerconfig)
	fmt.Println("Sleep 1s")
	err = c.Start()
	if err != nil {
		t.Error(err)
	}
	//Explicitly triggering GC to remove garbage from config loading.
	runtime.GC()

	{
		osSignals := make(chan os.Signal, 1)
		signal.Notify(osSignals, os.Interrupt, os.Kill, syscall.SIGTERM)
		<-osSignals
	}
}
