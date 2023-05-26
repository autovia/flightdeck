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

func RoleBindingHandler(app *S.App, c *S.Client, w http.ResponseWriter, r *http.Request) error {
	url := S.GetRequestParams(r, "/api/v1/rb/")
	log.Printf("RoleBindingHandler url: %v", url)

	rb, err := c.Clientset.RbacV1().RoleBindings(url.Namespace).Get(context.TODO(), url.Resource, metav1.GetOptions{})
	if err != nil {
		return S.RespondError(err)
	}
	rb.ObjectMeta.ManagedFields = nil

	return S.RespondYAML(w, http.StatusOK, rb)
}

func NamespaceRoleBindingListHandler(app *S.App, c *S.Client, w http.ResponseWriter, r *http.Request) error {
	url := S.GetRequestParams(r, "/api/v1/namespace/rb/")
	log.Printf("NamespaceRoleBindingListHandler url: %v", url)

	g := S.Graph{Nodes: []S.Node{}, Edges: []S.Edge{}}

	ns, err := c.Clientset.CoreV1().Namespaces().Get(context.TODO(), url.Namespace, metav1.GetOptions{})
	if err != nil {
		return S.RespondError(err)
	}
	node := g.AddNode("ns", string(ns.ObjectMeta.UID), ns.ObjectMeta.Name, S.NodeOptions{Type: "namespace", Group: true})

	rbList, err := c.Clientset.RbacV1().RoleBindings(url.Namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return S.RespondError(err)
	}
	for _, rb := range rbList.Items {
		g.AddNode("rb", string(rb.ObjectMeta.UID), rb.ObjectMeta.Name, S.NodeOptions{Namespace: url.Namespace, Type: "rb", ParentNode: node, Extent: "parent"})
	}

	return S.RespondJSON(w, http.StatusOK, g)
}
