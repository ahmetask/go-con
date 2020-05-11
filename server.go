package gocon

import (
	"context"
	"encoding/json"
	"fmt"
)

type Server struct {
	cm *ConfigManager
}

func (s *Server) GetConfig(_ context.Context, request *GetConfigRequest) (*GetConfigResponse, error) {
	s.cm.Logger.Info("Getting Config With AppName:" + request.AppName)
	cache, success := s.cm.cache.Get(request.AppName)

	if !success {
		return nil, NewError(EType.ConfigServer, nil, fmt.Sprintf("can not get cache with key:%s", request.AppName))
	}

	bytes, err := json.Marshal(cache)

	if err != nil {
		return nil, NewError(EType.Json, err)
	}

	response := &GetConfigResponse{
		Value:   bytes,
		Success: success,
	}

	return response, nil
}

func (s *Server) SaveConfig(_ context.Context, request *SaveConfigRequest) (*SuccessResponse, error) {
	s.cm.Logger.Info("Saving Config for:" + request.AppName)

	spec := Spec{
		AppName:   request.AppName,
		Namespace: request.Namespace,
		Port:      request.Port,
	}

	cfg := Config{
		Content: request.Content,
		Spec:    spec,
	}

	//saving
	err := s.cm.DbInstance.Save(cfg)

	if err != nil {
		return nil, err
	}

	//cached
	s.cm.cache.Set(request.AppName, cfg, -1)

	s.cm.triggerInstance <- request.AppName
	return &SuccessResponse{Success: true}, nil
}

func (s *Server) RefreshCfg(ctx context.Context, request *RefreshCfgRequest) (*SuccessResponse, error) {
	s.cm.Logger.Info(fmt.Sprintf("Refreshing Event Triggered For:%s", request.AppName))
	s.cm.triggerInstance <- request.AppName
	return &SuccessResponse{Success: true}, nil
}

func (s *Server) RefreshCfgs(ctx context.Context, in *RefreshCfgsRequest) (*SuccessResponse, error) {
	s.cm.Logger.Info("Refreshing Event Triggered")
	s.cm.triggerRefreshing <- true
	return &SuccessResponse{Success: true}, nil
}
