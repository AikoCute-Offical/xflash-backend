package core

import (
	"encoding/json"
	io "io/ioutil"
	"log"
	"sync"

	"github.com/AikoCute-Offical/xflash-backend/conf"
	"github.com/AikoCute-Offical/xflash-backend/core/app/dispatcher"
	_ "github.com/AikoCute-Offical/xflash-backend/core/distro/all"
	"github.com/xtls/xray-core/app/proxyman"
	"github.com/xtls/xray-core/app/stats"
	"github.com/xtls/xray-core/common/serial"
	"github.com/xtls/xray-core/core"
	"github.com/xtls/xray-core/features/inbound"
	"github.com/xtls/xray-core/features/outbound"
	"github.com/xtls/xray-core/features/routing"
	coreConf "github.com/xtls/xray-core/infra/conf"
)

// Core Structure
type Core struct {
	access     sync.Mutex
	Server     *core.Instance
	ihm        inbound.Manager
	ohm        outbound.Manager
	dispatcher *dispatcher.DefaultDispatcher
}

func New(c *conf.Conf) *Core {
	return &Core{Server: getCore(c)}
}

func parseConnectionConfig(c *conf.ConnetionConfig) (policy *coreConf.Policy) {
	policy = &coreConf.Policy{
		StatsUserUplink:   true,
		StatsUserDownlink: true,
		Handshake:         &c.Handshake,
		ConnectionIdle:    &c.ConnIdle,
		UplinkOnly:        &c.UplinkOnly,
		DownlinkOnly:      &c.DownlinkOnly,
		BufferSize:        &c.BufferSize,
	}
	return
}

func getCore(AikoCute-OfficalConfig *conf.Conf) *core.Instance {
	// Log Config
	coreLogConfig := &coreConf.LogConfig{}
	coreLogConfig.LogLevel = AikoCute-OfficalConfig.LogConfig.Level
	coreLogConfig.AccessLog = AikoCute-OfficalConfig.LogConfig.AccessPath
	coreLogConfig.ErrorLog = AikoCute-OfficalConfig.LogConfig.ErrorPath
	// DNS config
	coreDnsConfig := &coreConf.DNSConfig{}
	if AikoCute-OfficalConfig.DnsConfigPath != "" {
		if data, err := io.ReadFile(AikoCute-OfficalConfig.DnsConfigPath); err != nil {
			log.Panicf("Failed to read DNS config file at: %s", AikoCute-OfficalConfig.DnsConfigPath)
		} else {
			if err = json.Unmarshal(data, coreDnsConfig); err != nil {
				log.Panicf("Failed to unmarshal DNS config: %s", AikoCute-OfficalConfig.DnsConfigPath)
			}
		}
	}
	dnsConfig, err := coreDnsConfig.Build()
	if err != nil {
		log.Panicf("Failed to understand DNS config, Please check: https://xtls.github.io/config/dns.html for help: %s", err)
	}
	// Routing config
	coreRouterConfig := &coreConf.RouterConfig{}
	if AikoCute-OfficalConfig.RouteConfigPath != "" {
		if data, err := io.ReadFile(AikoCute-OfficalConfig.RouteConfigPath); err != nil {
			log.Panicf("Failed to read Routing config file at: %s", AikoCute-OfficalConfig.RouteConfigPath)
		} else {
			if err = json.Unmarshal(data, coreRouterConfig); err != nil {
				log.Panicf("Failed to unmarshal Routing config: %s", AikoCute-OfficalConfig.RouteConfigPath)
			}
		}
	}
	routeConfig, err := coreRouterConfig.Build()
	if err != nil {
		log.Panicf("Failed to understand Routing config  Please check: https://xtls.github.io/config/routing.html for help: %s", err)
	}
	// Custom Inbound config
	var coreCustomInboundConfig []coreConf.InboundDetourConfig
	if AikoCute-OfficalConfig.InboundConfigPath != "" {
		if data, err := io.ReadFile(AikoCute-OfficalConfig.InboundConfigPath); err != nil {
			log.Panicf("Failed to read Custom Inbound config file at: %s", AikoCute-OfficalConfig.OutboundConfigPath)
		} else {
			if err = json.Unmarshal(data, &coreCustomInboundConfig); err != nil {
				log.Panicf("Failed to unmarshal Custom Inbound config: %s", AikoCute-OfficalConfig.OutboundConfigPath)
			}
		}
	}
	var inBoundConfig []*core.InboundHandlerConfig
	for _, config := range coreCustomInboundConfig {
		oc, err := config.Build()
		if err != nil {
			log.Panicf("Failed to understand Inbound config, Please check: https://xtls.github.io/config/inbound.html for help: %s", err)
		}
		inBoundConfig = append(inBoundConfig, oc)
	}
	// Custom Outbound config
	var coreCustomOutboundConfig []coreConf.OutboundDetourConfig
	if AikoCute-OfficalConfig.OutboundConfigPath != "" {
		if data, err := io.ReadFile(AikoCute-OfficalConfig.OutboundConfigPath); err != nil {
			log.Panicf("Failed to read Custom Outbound config file at: %s", AikoCute-OfficalConfig.OutboundConfigPath)
		} else {
			if err = json.Unmarshal(data, &coreCustomOutboundConfig); err != nil {
				log.Panicf("Failed to unmarshal Custom Outbound config: %s", AikoCute-OfficalConfig.OutboundConfigPath)
			}
		}
	}
	var outBoundConfig []*core.OutboundHandlerConfig
	for _, config := range coreCustomOutboundConfig {
		oc, err := config.Build()
		if err != nil {
			log.Panicf("Failed to understand Outbound config, Please check: https://xtls.github.io/config/outbound.html for help: %s", err)
		}
		outBoundConfig = append(outBoundConfig, oc)
	}
	// Policy config
	levelPolicyConfig := parseConnectionConfig(AikoCute-OfficalConfig.ConnectionConfig)
	corePolicyConfig := &coreConf.PolicyConfig{}
	corePolicyConfig.Levels = map[uint32]*coreConf.Policy{0: levelPolicyConfig}
	policyConfig, _ := corePolicyConfig.Build()
	// Build Core conf
	config := &core.Config{
		App: []*serial.TypedMessage{
			serial.ToTypedMessage(coreLogConfig.Build()),
			serial.ToTypedMessage(&dispatcher.Config{}),
			serial.ToTypedMessage(&stats.Config{}),
			serial.ToTypedMessage(&proxyman.InboundConfig{}),
			serial.ToTypedMessage(&proxyman.OutboundConfig{}),
			serial.ToTypedMessage(policyConfig),
			serial.ToTypedMessage(dnsConfig),
			serial.ToTypedMessage(routeConfig),
		},
		Inbound:  inBoundConfig,
		Outbound: outBoundConfig,
	}
	server, err := core.New(config)
	if err != nil {
		log.Panicf("failed to create instance: %s", err)
	}
	log.Printf("Core Version: %s", core.Version())

	return server
}

// Start the Core
func (p *Core) Start() {
	p.access.Lock()
	defer p.access.Unlock()
	log.Print("Start the xflash..")
	if err := p.Server.Start(); err != nil {
		log.Panicf("Failed to start instance: %s", err)
	}
	p.ihm = p.Server.GetFeature(inbound.ManagerType()).(inbound.Manager)
	p.ohm = p.Server.GetFeature(outbound.ManagerType()).(outbound.Manager)
	p.dispatcher = p.Server.GetFeature(routing.DispatcherType()).(*dispatcher.DefaultDispatcher)
	return
}

// Close  the core
func (p *Core) Close() {
	p.access.Lock()
	defer p.access.Unlock()
	p.Server.Close()
	return
}
