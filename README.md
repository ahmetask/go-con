# go-con

go-con is a helper for building a config-server.

## Installation

You need [golang 1.13](https://golang.org/dl/) to use go-con.


```go
 go get github.com/ahmetask/go-con
```
## Example 
[Tutorial](https://medium.com/@ahmet9417/golang-config-server-as-a-service-3655bf832c2d)

[Server](https://github.com/ahmetask/go-con-server-example)

[Client](https://github.com/ahmetask/go-con-client-example)

[Client-2](https://github.com/ahmetask/go-con-client-example-2)

[UI Helper(Connects UI with Server)](https://github.com/ahmetask/go-con-manager)

[UI (Simple React Project that control Server)](https://github.com/ahmetask/go-con-ui)


## Usage
- First, implement your DB interface

```go
package main

import (
	gocon "github.com/ahmetask/go-con"
)

type CustomDB struct {
}

func (db *CustomDB) Save(config gocon.Config) error {

	return nil
}

func (db *CustomDB) Read(appName string) (*gocon.Config, error) {
	var config gocon.Config

	return &config, nil
}

func (db *CustomDB) ReadSpecs() ([]gocon.Spec, error) {

	return []gocon.Spec, nil
}


```
- Then run config server

```go
configManager := gocon.ConfigManager{
	DbInstance:       &customDB,
	ConfigServerPort: "8080",
	BroadcastModuleConfig: gocon.BroadcastModuleConfig{
		Available: true,
		IsLocal:   false,
	},
}

err := configManager.Initialize()

panic(err)

```

- You can also add your own logger instance if you want. It's not necessary. Go-con uses default logger

- It also provides helper functions for Golang projects to control the server and client.

- if you use another programing languages use proto files to interract go-con server

- How to get the updated config? (See [Client](https://github.com/ahmetask/go-con-client-example/blob/master/main.go#L35))

```go
type Config struct {
	
}

var CustomConfig Config

go func() {
	err := gocon.ListenConfigChange("port", &CustomConfig, nil)
	fmt.Println(err)
}()
```
- It will listen config server and change your CustomConfig instance whenever update arrives.

- So that's all. Don't forget to look at the examples. You can change it according to your own wishes

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.
