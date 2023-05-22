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

//log.Printf("HERE IS YOUR USER: %s", r.Context().Value("user"))

func NamespaceListHandler(client *kubernetes.Clientset, w http.ResponseWriter, r *http.Request) error {
	log.Print("NamespaceListHandler")

	g := S.Graph{Nodes: []S.Node{}, Edges: []S.Edge{}}

	nsList, err := client.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return S.RespondError(err)
	}
	for _, ns := range nsList.Items {
		// SLOW
		// podList, err := app.Client.CoreV1().Pods(ns.ObjectMeta.Name).List(context.TODO(), metav1.ListOptions{})
		// if err != nil {
		// 	return S.RespondError(err)
		// }
		g.AddNode("ns", string(ns.ObjectMeta.UID), ns.ObjectMeta.Name, S.NodeOptions{Type: "namespace", Draggable: false, Connectable: false})
	}
	return S.RespondJSON(w, http.StatusOK, g)
}
