// Copyright (c) Autovia GmbH
// SPDX-License-Identifier: Apache-2.0

package k8s

import (
	"context"
	"log"
	"net/http"

	S "github.com/autovia/flightdeck/api/structs"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func IngressHandler(app *S.App, client *kubernetes.Clientset, w http.ResponseWriter, r *http.Request) error {
	url := S.GetRequestParams(r, "/api/v1/ing/")
	log.Printf("IngressHandler url: %v", url)

	role, err := client.NetworkingV1().Ingresses(url.Namespace).Get(context.TODO(), url.Resource, metav1.GetOptions{})
	if err != nil {
		return S.RespondError(err)
	}
	role.ObjectMeta.ManagedFields = nil

	S.RespondYAML(w, http.StatusOK, role)
	return nil
}

func NamespaceIngressListHandler(app *S.App, client *kubernetes.Clientset, w http.ResponseWriter, r *http.Request) error {
	url := S.GetRequestParams(r, "/api/v1/namespace/ing/")
	log.Printf("NamespaceIngressListHandler url: %v", url)

	g := S.Graph{Nodes: []S.Node{}, Edges: []S.Edge{}}

	ns, err := client.CoreV1().Namespaces().Get(context.TODO(), url.Namespace, metav1.GetOptions{})
	if err != nil {
		return S.RespondError(err)
	}
	node := g.AddNode("ns", string(ns.ObjectMeta.UID), ns.ObjectMeta.Name, S.NodeOptions{Type: "namespace", Group: true})

	ingList, err := client.NetworkingV1().Ingresses(url.Namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return S.RespondError(err)
	}
	for _, ing := range ingList.Items {
		g.AddNode("ing", string(ing.ObjectMeta.UID), ing.ObjectMeta.Name, S.NodeOptions{Namespace: url.Namespace, Type: "ing", ParentNode: node, Extent: "parent"})
	}

	return S.RespondJSON(w, http.StatusOK, g)
}
