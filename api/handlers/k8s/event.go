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

func EventHandler(app *S.App, client *kubernetes.Clientset, w http.ResponseWriter, r *http.Request) error {
	url := S.GetRequestParams(r, "/api/v1/ev/")
	log.Printf("EventHandler url: %v", url)

	role, err := client.CoreV1().Events(url.Namespace).Get(context.TODO(), url.Resource, metav1.GetOptions{})
	if err != nil {
		return S.RespondError(err)
	}
	role.ObjectMeta.ManagedFields = nil

	S.RespondYAML(w, http.StatusOK, role)
	return nil
}

func NamespaceEventListHandler(app *S.App, client *kubernetes.Clientset, w http.ResponseWriter, r *http.Request) error {
	url := S.GetRequestParams(r, "/api/v1/namespace/ev/")
	log.Printf("NamespaceEventListHandler url: %v", url)

	g := S.Graph{Nodes: []S.Node{}, Edges: []S.Edge{}}

	ns, err := client.CoreV1().Namespaces().Get(context.TODO(), url.Namespace, metav1.GetOptions{})
	if err != nil {
		return S.RespondError(err)
	}
	node := g.AddNode("ns", string(ns.ObjectMeta.UID), ns.ObjectMeta.Name, S.NodeOptions{Type: "namespace", Group: true})

	evList, err := client.CoreV1().Events(url.Namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return S.RespondError(err)
	}
	for _, ev := range evList.Items {
		g.AddNode("ev", string(ev.ObjectMeta.UID), ev.ObjectMeta.Name, S.NodeOptions{Namespace: url.Namespace, Type: "ev", ParentNode: node, Extent: "parent"})
	}

	return S.RespondJSON(w, http.StatusOK, g)
}
