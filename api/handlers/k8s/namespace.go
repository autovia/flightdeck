// Copyright (c) Autovia GmbH
// SPDX-License-Identifier: Apache-2.0

package k8s

import (
	"context"
	"log"
	"net/http"

	S "github.com/autovia/flightdeck/api/structs"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func NamespaceListHandler(app *S.App, c *S.Client, w http.ResponseWriter, r *http.Request) error {
	log.Print("NamespaceListHandler")

	g := S.Graph{Nodes: []S.Node{}, Edges: []S.Edge{}}

	nsList, err := c.Clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return S.RespondError(err)
	}
	g.AddNode("ns", "all", "all", S.NodeOptions{Type: "namespace", Draggable: false, Connectable: false})
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
