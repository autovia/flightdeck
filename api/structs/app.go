// Copyright (c) Autovia GmbH
// SPDX-License-Identifier: Apache-2.0

package structs

import (
	"net/http"

	clientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

type App struct {
	Addr           *string
	Router         *http.ServeMux
	KubeConfigPath *string
	ApiServerHost  *string
	ProxyUrl       *string
	InCluster      *bool
	FileServer     *bool
	FileServerPath *string
	Client         *Client
	AuthManager    AuthManager
	// config
	PodLogTailLines   *int64
	DefaultConfigName *string
}

func (app *App) LoadKubeContext(context string) error {
	var config *rest.Config
	var err error

	if *app.InCluster {
		config, err = rest.InClusterConfig()
		if err != nil {
			panic(err.Error())
		}
	} else {
		if context == "" {
			config, err = clientcmd.BuildConfigFromFlags("", *app.KubeConfigPath)
			if err != nil {
				return err
			}
		} else {
			config, err = clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
				&clientcmd.ClientConfigLoadingRules{ExplicitPath: *app.KubeConfigPath},
				&clientcmd.ConfigOverrides{
					CurrentContext: context,
				}).ClientConfig()
			if err != nil {
				return err
			}
		}
	}

	k8sclientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	apiclientset, err := clientset.NewForConfig(config)
	if err != nil {
		return err
	}

	app.Client = &Client{
		Clientset: k8sclientset,
		ApiClient: apiclientset,
		Config:    config,
	}

	return nil
}

func (app *App) NewKubeClient(token string) (*Client, error) {
	config, err := app.buildConfigFromToken(token)
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return &Client{Clientset: clientset, Config: config}, nil
}

func (app *App) NewApiClient(token string) (*Client, error) {
	config, err := app.buildConfigFromToken(token)
	if err != nil {
		return nil, err
	}

	apiclientset, err := clientset.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return &Client{ApiClient: apiclientset, Config: config}, nil
}

func (app *App) buildConfigFromToken(token string) (*rest.Config, error) {
	var clientCfg *rest.Config
	var err error

	if *app.InCluster {
		cfg, err := rest.InClusterConfig()
		if err != nil {
			return nil, err
		}
		apiConfig := api.NewConfig()
		apiConfig.Clusters[*app.DefaultConfigName] = &api.Cluster{
			Server:                   cfg.Host,
			CertificateAuthority:     cfg.TLSClientConfig.CAFile,
			CertificateAuthorityData: cfg.TLSClientConfig.CAData,
			InsecureSkipTLSVerify:    cfg.TLSClientConfig.Insecure,
		}
		apiConfig.AuthInfos[*app.DefaultConfigName] = &api.AuthInfo{Token: token}
		apiConfig.Contexts[*app.DefaultConfigName] = &api.Context{
			Cluster:  *app.DefaultConfigName,
			AuthInfo: *app.DefaultConfigName,
		}
		apiConfig.CurrentContext = *app.DefaultConfigName

		clientConfig := clientcmd.NewDefaultClientConfig(
			*apiConfig,
			&clientcmd.ConfigOverrides{},
		)

		clientCfg, err = clientConfig.ClientConfig()
		if err != nil {
			return nil, err
		}
	} else {
		clientCfg, err = clientcmd.BuildConfigFromFlags("", *app.KubeConfigPath)
		if err != nil {
			return nil, err
		}
	}
	return clientCfg, nil
}

// Use configmap
func (app *App) LoadConfig() {
	// PodLogOptions, the number of lines from the end of the logs to show
	var lines int64 = 1024
	app.PodLogTailLines = &lines

	var configname string = "kubernetes"
	app.DefaultConfigName = &configname
}
