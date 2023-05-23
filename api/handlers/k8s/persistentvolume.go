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

func PersistentVolumeHandler(app *S.App, client *kubernetes.Clientset, w http.ResponseWriter, r *http.Request) error {
	url := S.GetRequestParams(r, "/api/v1/pv/")
	log.Printf("PersistentVolumeHandler url: %v", url)

	pv, err := client.CoreV1().PersistentVolumes().Get(context.TODO(), url.Resource, metav1.GetOptions{})
	if err != nil {
		return S.RespondError(err)
	}
	pv.ObjectMeta.ManagedFields = nil

	return S.RespondYAML(w, http.StatusOK, pv)
}

func PersistentVolumeListHandler(app *S.App, client *kubernetes.Clientset, w http.ResponseWriter, r *http.Request) error {
	log.Print("PersistentVolumeListHandler")

	g := S.Graph{Nodes: []S.Node{}, Edges: []S.Edge{}}

	pvList, err := client.CoreV1().PersistentVolumes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return S.RespondError(err)
	}
	for _, pv := range pvList.Items {
		g.AddNode("pv", string(pv.ObjectMeta.UID), pv.ObjectMeta.Name, S.NodeOptions{Type: "pv", Draggable: false, Connectable: false})
	}

	return S.RespondJSON(w, http.StatusOK, g)
}
