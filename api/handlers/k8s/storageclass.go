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

func StorageClassHandler(app *S.App, c *S.Client, w http.ResponseWriter, r *http.Request) error {
	url := S.GetRequestParams(r, "/api/v1/sc/")
	log.Printf("StorageClassHandler url: %v", url)

	sc, err := c.Clientset.StorageV1().StorageClasses().Get(context.TODO(), url.Resource, metav1.GetOptions{})
	if err != nil {
		return S.RespondError(err)
	}
	sc.ObjectMeta.ManagedFields = nil

	return S.RespondFormat(r, w, http.StatusOK, sc)
}

func StorageClassListHandler(app *S.App, c *S.Client, w http.ResponseWriter, r *http.Request) error {
	log.Print("StorageClassListHandler")

	g := S.Graph{Nodes: []S.Node{}, Edges: []S.Edge{}}

	scList, err := c.Clientset.StorageV1().StorageClasses().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return S.RespondError(err)
	}
	for _, sc := range scList.Items {
		g.AddNode("sc", string(sc.ObjectMeta.UID), sc.ObjectMeta.Name, S.NodeOptions{Type: "sc", Draggable: false, Connectable: false})
	}
	return S.RespondJSON(w, http.StatusOK, g)
}
