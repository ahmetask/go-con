package gocon

import (
	"context"
	"encoding/json"
	"google.golang.org/grpc"
	"net"
)

var configFile interface{}

type configC struct {
	Port   string
	Logger ILogger
}

type configClientServer struct {
}

func (c *configClientServer) ChangeConfig(_ context.Context, in *RefreshConfigRequest) (*RefreshConfigResponse, error) {
	var resp Config
	//Unmarshal config file
	marshalError := json.Unmarshal(in.Config, &resp)
	if marshalError != nil {
		return &RefreshConfigResponse{Success: false}, marshalError
	}

	//Unmarshal config content
	marshalError = json.Unmarshal(resp.Content, &configFile)
	if marshalError != nil {
		return &RefreshConfigResponse{Success: false}, marshalError
	}

	return &RefreshConfigResponse{Success: true}, nil
}

//pass your config struct address
func ListenConfigChange(appPort string, configFileAddress interface{}, logger ILogger) error {
	client := configC{
		Port: appPort,
	}
	if logger == nil {
		client.Logger = DefaultLogger()
	}

	logger = client.Logger
	configFile = configFileAddress

	err := client.runConfigClient()
	return err
}

func (c *configC) runConfigClient() error {
	if c.Port == "" {
		c.Logger.Error("port should not be empty")
		return NewError(EType.ConfigClient, nil, "port should not be empty")
	}

	address := "0.0.0.0:" + c.Port
	listen, err := net.Listen("tcp", address)
	if err != nil {
		c.Logger.Error(err)
		return NewError(EType.ConfigClient, err)
	}

	s := grpc.NewServer()
	RegisterConfigCServer(s, &configClientServer{})

	serveError := s.Serve(listen)

	if serveError != nil {
		c.Logger.Error(err)
		return NewError(EType.ConfigClient, serveError)
	}
	return nil
}
