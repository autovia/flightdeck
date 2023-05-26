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

func IngressClassHandler(app *S.App, c *S.Client, w http.ResponseWriter, r *http.Request) error {
	url := S.GetRequestParams(r, "/api/v1/ic/")
	log.Printf("IngressClassHandler url: %v", url)

	netpol, err := c.Clientset.NetworkingV1().IngressClasses().Get(context.TODO(), url.Resource, metav1.GetOptions{})
	if err != nil {
		return S.RespondError(err)
	}
	netpol.ObjectMeta.ManagedFields = nil

	S.RespondYAML(w, http.StatusOK, netpol)
	return nil
}

func IngressClassListHandler(app *S.App, c *S.Client, w http.ResponseWriter, r *http.Request) error {
	log.Print("IngressClassListHandler")

	g := S.Graph{Nodes: []S.Node{}, Edges: []S.Edge{}}

	icList, err := c.Clientset.NetworkingV1().IngressClasses().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return S.RespondError(err)
	}
	for _, ic := range icList.Items {
		g.AddNode("ing", string(ic.ObjectMeta.UID), ic.ObjectMeta.Name, S.NodeOptions{Type: "ic", Draggable: false, Connectable: false})
	}
	return S.RespondJSON(w, http.StatusOK, g)
}
