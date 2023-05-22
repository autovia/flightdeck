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

func ClusterRoleBindingHandler(client *kubernetes.Clientset, w http.ResponseWriter, r *http.Request) error {
	url := S.GetRequestParams(r, "/api/v1/crb/")
	log.Printf("ClusterRoleBindingHandler url: %v", url)

	role, err := client.RbacV1().ClusterRoleBindings().Get(context.TODO(), url.Resource, metav1.GetOptions{})
	if err != nil {
		return S.RespondError(err)
	}
	role.ObjectMeta.ManagedFields = nil

	return S.RespondYAML(w, http.StatusOK, role)
}

func ClusterRoleBindingListHandler(client *kubernetes.Clientset, w http.ResponseWriter, r *http.Request) error {
	log.Print("ClusterRoleBindingListHandler")

	g := S.Graph{Nodes: []S.Node{}, Edges: []S.Edge{}}

	crbList, err := client.RbacV1().ClusterRoleBindings().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return S.RespondError(err)
	}
	for _, crb := range crbList.Items {
		g.AddNode("crb", string(crb.ObjectMeta.UID), crb.ObjectMeta.Name, S.NodeOptions{Type: "crb", Draggable: false, Connectable: false})
	}
	return S.RespondJSON(w, http.StatusOK, g)
}
