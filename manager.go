package gocon

import (
	"encoding/json"
	"github.com/ahmetask/gocache/v2"
	"google.golang.org/grpc"
	"net"
	"time"
)

type ConfigManager struct {
	DbInstance            Db                    //Implement your own db logic
	Logger                ILogger               // Pass logger interface if you want
	ConfigServerPort      string                // It's crucial
	ConfigRefreshRate     time.Duration         // Default is 1 hour
	BroadcastModuleConfig BroadcastModuleConfig // This Flags are important
	cache                 *gocache.Cache
	broadcastModule       *broadcastModule
	specs                 []Spec
	triggerRefreshing     chan bool
	triggerInstance       chan string
}

func (cm *ConfigManager) Initialize() *Error {

	if cm.Logger == nil {
		cm.Logger = DefaultLogger()
	}

	if cm.DbInstance == nil {
		cm.Logger.Fatal("DB instance error")
		return NewError(EType.DB, nil, "You Should implement db interface")
	}

	cm.Logger.Info("Initializing Manager...")

	cm.triggerRefreshing = make(chan bool)
	cm.triggerInstance = make(chan string)

	cm.cache = gocache.NewCache(gocache.Eternal, 0)

	cm.Logger.Info("Getting Specs...")
	configListErr := cm.getSpecs()
	if configListErr != nil {
		cm.Logger.Error(configListErr)
		return configListErr
	}

	cm.Logger.Info("Getting Configs...")
	configErr := cm.getConfigs()
	if configErr != nil {
		cm.Logger.Error(configErr)
		return configErr
	}

	if cm.BroadcastModuleConfig.Available {
		broadcastModuleErr := cm.initBroadcastModule()
		if broadcastModuleErr != nil {
			cm.Logger.Error(broadcastModuleErr)
			return broadcastModuleErr
		}
	}

	go cm.refreshConfig()

	return cm.runConfigServer()
}

func (cm *ConfigManager) refreshAll() {
	cm.Logger.Info("Configs are refreshing...")

	err := cm.getSpecs()
	if err != nil {
		cm.Logger.Error(err)
	}

	err = cm.getConfigs()
	if err != nil {
		cm.Logger.Error(err)
	}

	if cm.BroadcastModuleConfig.Available {
		cm.runBroadcastModule()
	}
}

func (cm *ConfigManager) refreshConfig() {
	if cm.ConfigRefreshRate <= 0 {
		cm.ConfigRefreshRate = 1 * time.Hour
	}

	ticker := time.NewTicker(cm.ConfigRefreshRate)

	for {
		select {
		case <-ticker.C:
			cm.refreshAll()
		case <-cm.triggerRefreshing:
			cm.refreshAll()
		case appName := <-cm.triggerInstance:
			_ = cm.getSpecs()
			_ = cm.getConfig(appName)
			if cm.BroadcastModuleConfig.Available {
				cm.refreshListener(appName)
			}
		}
	}
}

func (cm *ConfigManager) getSpecs() *Error {
	specs, err := cm.DbInstance.ReadSpecs()

	if err != nil {
		cm.Logger.Error(err)
		return NewError(EType.DB, err)
	}

	cm.specs = specs

	bs, err := json.Marshal(specs)
	if err != nil {
		return NewError(EType.Json, err)
	}

	specConfig := Config{
		Spec: Spec{
			AppName: "go-con",
		},
		Content: bs,
	}

	cm.cache.Set("specs", specConfig, -1)

	return nil
}

func (cm *ConfigManager) getConfigs() *Error {
	for _, spec := range cm.specs {
		config, err := cm.DbInstance.Read(spec.AppName)
		if err != nil {
			cm.Logger.Error(err)
			return NewError(EType.DB, err)
		}
		cm.cache.Set(config.Spec.AppName, config, -1)
	}
	return nil
}

func (cm *ConfigManager) getConfig(appName string) *Error {
	config, err := cm.DbInstance.Read(appName)
	if err != nil {
		cm.Logger.Error(err)
		return NewError(EType.DB, err)
	}
	cm.cache.Set(config.Spec.AppName, config, -1)

	return nil
}

func (cm *ConfigManager) newCacheServer(config gocache.ServerConfig) *Error {

	valid, validationErr := validateServerConfig(config)

	if !valid {
		return validationErr
	}

	address := "0.0.0.0:" + config.Port
	listen, err := net.Listen("tcp", address)
	if err != nil {
		return NewError(EType.ConfigServer, err)
	}

	s := grpc.NewServer()
	RegisterConfigServer(s, &Server{cm: cm})

	serveError := s.Serve(listen)

	if serveError != nil {
		return NewError(EType.ConfigServer, serveError)
	}

	return nil
}

func validateServerConfig(config gocache.ServerConfig) (bool, *Error) {
	if config.Port == "" {
		return false, NewError(EType.Validation, nil, "port should not be empty")
	}

	return true, nil
}

func (cm *ConfigManager) runConfigServer() *Error {
	serverConfig := gocache.ServerConfig{
		CachePtr: cm.cache,
		Port:     cm.ConfigServerPort,
	}

	err := cm.newCacheServer(serverConfig)

	if err != nil {
		cm.Logger.Fatal(err)
		return NewError(EType.ConfigServer, err)
	}

	return nil
}

func (cm *ConfigManager) refreshListener(appName string) {
	cm.broadcastModule.broadcastToSingleListener(appName, cm.specs)
}

func (cm *ConfigManager) runBroadcastModule() {
	cm.broadcastModule.broadcastToAll(cm.specs)
}

func (cm *ConfigManager) initBroadcastModule() *Error {
	cm.Logger.Info("Initializing Broadcast module...")
	module, err := newBroadcastModule(cm.Logger, cm.cache, cm.BroadcastModuleConfig)
	if err != nil {
		return err
	}
	cm.broadcastModule = module
	return nil
}
