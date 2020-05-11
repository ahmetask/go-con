package gocon

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ahmetask/gocache/v2"
	"google.golang.org/grpc"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type BroadcastModuleConfig struct {
	Available bool //Activate Broadcast module
	IsLocal   bool //Run for local env
}

type broadcastModule struct {
	logger ILogger
	cache  *gocache.Cache
	k8s    *kubernetes.Clientset
	config BroadcastModuleConfig
}

func newBroadcastModule(logger ILogger, cache *gocache.Cache, broadcastConfig BroadcastModuleConfig) (*broadcastModule, *Error) {
	module := broadcastModule{
		logger: logger,
		cache:  cache,
		config: broadcastConfig,
	}

	//Kubernetes
	if !broadcastConfig.IsLocal {
		logger.Info("K8s configuration...")
		//it is available for k8s cluster
		config, err := getK8sInClusterConfig()
		if err != nil {
			return nil, err
		}

		// create kubernetes client
		k8s, clientErr := kubernetes.NewForConfig(config)
		if clientErr != nil {
			return nil, NewError(EType.BroadcastModule, clientErr)
		}

		module.k8s = k8s
	}

	return &module, nil
}

func getK8sInClusterConfig() (*rest.Config, *Error) {
	c, err := rest.InClusterConfig()
	if err != nil {
		return nil, NewError(EType.BroadcastModule, err)
	}
	return c, nil
}

func (bm *broadcastModule) broadcastToAll(specs []Spec) {
	for _, spec := range specs {
		if bm.config.IsLocal {
			bm.broadcastToLocalListener(spec)
		} else {
			bm.broadcastToK8sInstance(spec)
		}
	}
}

func (bm *broadcastModule) broadcastToSingleListener(appName string, specs []Spec) {
	for _, l := range specs {
		if l.AppName == appName {
			if bm.config.IsLocal {
				bm.broadcastToLocalListener(l)
			} else {
				bm.broadcastToK8sInstance(l)
			}
			break
		}
	}
}

func (bm *broadcastModule) broadcastToK8sInstance(spec Spec) {
	labelSelector := meta.LabelSelector{MatchLabels: map[string]string{"app": spec.AppName}}

	pods, err := bm.k8s.CoreV1().Pods(spec.Namespace).List(context.Background(),
		meta.ListOptions{LabelSelector: labels.Set(labelSelector.MatchLabels).String()})

	if err != nil {
		bm.logger.Error(err)
		return
	}
	bm.logger.Info(fmt.Sprintf("app:%s, %d Pods", spec.AppName, len(pods.Items)))
	for _, pod := range pods.Items {
		bm.callListener(pod.Status.PodIP, spec)
	}
}

func (bm *broadcastModule) broadcastToLocalListener(spec Spec) {
	bm.callListener("localhost", spec)
}

func (bm *broadcastModule) callListener(ipAddress string, spec Spec) {
	address := ipAddress + ":" + spec.Port

	bm.logger.Info(fmt.Sprintf("Activate listener For:%s|%s", spec.AppName, address))

	cc, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		bm.logger.Error(err)
		return
	}
	defer cc.Close()

	client := NewConfigCClient(cc)

	configFile, cacheError := bm.cache.Get(spec.AppName)

	if cacheError != true {
		bm.logger.Error(cacheError)
		return
	}

	bytes, jsonError := json.Marshal(configFile)
	if jsonError != nil {
		bm.logger.Error(cacheError)
		return
	}

	request := &RefreshConfigRequest{AppName: spec.AppName, Config: bytes}

	response, clientError := client.ChangeConfig(context.Background(), request)
	if clientError != nil {
		bm.logger.Error(clientError)
		return
	}

	bm.logger.Info(fmt.Sprintf(" For:%s|%s Result:%v", spec.AppName, address, response.Success))
}
