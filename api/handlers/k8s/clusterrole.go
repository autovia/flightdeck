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

func ClusterRoleHandler(client *kubernetes.Clientset, w http.ResponseWriter, r *http.Request) error {
	url := S.GetRequestParams(r, "/api/v1/c-role/")
	log.Printf("ClusterRoleHandler url: %v", url)

	role, err := client.RbacV1().ClusterRoles().Get(context.TODO(), url.Resource, metav1.GetOptions{})
	if err != nil {
		return S.RespondError(err)
	}
	role.ObjectMeta.ManagedFields = nil

	return S.RespondYAML(w, http.StatusOK, role)
}

func ClusterRoleListHandler(client *kubernetes.Clientset, w http.ResponseWriter, r *http.Request) error {
	log.Print("ClusterRoleListHandler")

	g := S.Graph{Nodes: []S.Node{}, Edges: []S.Edge{}}

	crList, err := client.RbacV1().ClusterRoles().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return S.RespondError(err)
	}
	for _, cr := range crList.Items {
		g.AddNode("c-role", string(cr.ObjectMeta.UID), cr.ObjectMeta.Name, S.NodeOptions{Type: "cr", Draggable: false, Connectable: false})
	}
	return S.RespondJSON(w, http.StatusOK, g)
}
