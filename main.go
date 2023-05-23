// Copyright (c) Autovia GmbH
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"flag"
	"log"
	"net/http"
	"path/filepath"

	"k8s.io/client-go/util/homedir"

	S "github.com/autovia/flightdeck/api/structs"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	app := &S.App{}

	app.Addr = flag.String("addr", ":3000", "TCP address for the server to listen on, in the form host:port")
	app.InCluster = flag.Bool("incluster", false, "If true; use in-cluster k8s client else .kube/config file")
	app.ApiServerHost = flag.String("apiserver", "https://192.168.122.8:6443", "TCP address for the server to listen on, in the form host:port")
	app.ProxyUrl = flag.String("proxy", "socks5://host.docker.internal:2022", "TCP address for the server to listen on, in the form host:port")
	app.FileServer = flag.Bool("fileserver", false, "If true; serve frontend files from dist/ folder else use server from package.json")
	app.FileServerPath = flag.String("fileserverpath", "./dist", "Folder containing the frontend files")
	if home := homedir.HomeDir(); home != "" {
		app.KubeConfigPath = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		app.KubeConfigPath = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// Config
	app.LoadConfig()

	// Kube client
	app.LoadKubeContext("")

	// Router
	app.Router = http.NewServeMux()
	InitRoutes(app)

	// JWT manager
	app.AuthManager = S.NewAuthManager(app.Client)

	// Server
	srv := &http.Server{
		Addr:    *app.Addr,
		Handler: app.Router,
		//TLSConfig:    cfg,
		//TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
	}
	log.Printf("Listen on %s", *app.Addr)
	log.Fatal(srv.ListenAndServe())
}
