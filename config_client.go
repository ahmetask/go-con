package gocon

import (
	"context"
	"encoding/json"
	"google.golang.org/grpc"
)

type AddConfigRequest struct {
	AppName   string
	Port      string
	Namespace string
	Content   interface{}
}

func GetConfig(address, appName string, configResponse interface{}) error {
	cc, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer cc.Close()

	client := NewConfigClient(cc)
	request := &GetConfigRequest{AppName: appName}

	resp, clientError := client.GetConfig(context.Background(), request)
	if clientError != nil {
		return clientError
	}

	marshallError := json.Unmarshal(resp.Value, &configResponse)
	if marshallError != nil {
		return marshallError
	}

	return nil
}

func SaveConfig(address string, request AddConfigRequest) error {
	cc, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer cc.Close()

	client := NewConfigClient(cc)

	bytes, bodyError := json.Marshal(&request.Content)

	if bodyError != nil {
		return bodyError
	}

	apiRequest := &SaveConfigRequest{
		AppName:   request.AppName,
		Port:      request.Port,
		Namespace: request.Namespace,
		Content:   bytes,
	}

	_, clientError := client.SaveConfig(context.Background(), apiRequest)
	if clientError != nil {
		return clientError
	}

	return nil
}

func RefreshConfig(address string, appName string) error {
	cc, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer cc.Close()

	client := NewConfigClient(cc)

	apiRequest := &RefreshCfgRequest{
		AppName: appName,
	}

	_, clientError := client.RefreshCfg(context.Background(), apiRequest)
	if clientError != nil {
		return clientError
	}

	return nil
}

func RefreshConfigs(address string) error {
	cc, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer cc.Close()

	client := NewConfigClient(cc)

	apiRequest := &RefreshCfgsRequest{}

	_, clientError := client.RefreshCfgs(context.Background(), apiRequest)
	if clientError != nil {
		return clientError
	}

	return nil
}
