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

func RoleHandler(app *S.App, c *S.Client, w http.ResponseWriter, r *http.Request) error {
	url := S.GetRequestParams(r, "/api/v1/role/")
	log.Printf("RoleHandler url: %v", url)

	role, err := c.Clientset.RbacV1().Roles(url.Namespace).Get(context.TODO(), url.Resource, metav1.GetOptions{})
	if err != nil {
		return S.RespondError(err)
	}
	role.ObjectMeta.ManagedFields = nil

	return S.RespondYAML(w, http.StatusOK, role)
}

func NamespaceRoleListHandler(app *S.App, c *S.Client, w http.ResponseWriter, r *http.Request) error {
	url := S.GetRequestParams(r, "/api/v1/namespace/role/")
	log.Printf("NamespaceRoleListHandler url: %v", url)

	g := S.Graph{Nodes: []S.Node{}, Edges: []S.Edge{}}

	ns, err := c.Clientset.CoreV1().Namespaces().Get(context.TODO(), url.Namespace, metav1.GetOptions{})
	if err != nil {
		return S.RespondError(err)
	}
	node := g.AddNode("ns", string(ns.ObjectMeta.UID), ns.ObjectMeta.Name, S.NodeOptions{Type: "namespace", Group: true})

	roleList, err := c.Clientset.RbacV1().Roles(url.Namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return S.RespondError(err)
	}
	for _, role := range roleList.Items {
		g.AddNode("role", string(role.ObjectMeta.UID), role.ObjectMeta.Name, S.NodeOptions{Namespace: url.Namespace, Type: "role", ParentNode: node, Extent: "parent"})
	}

	return S.RespondJSON(w, http.StatusOK, g)
}
