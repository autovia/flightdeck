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

func NetworkPolicyHandler(app *S.App, c *S.Client, w http.ResponseWriter, r *http.Request) error {
	url := S.GetRequestParams(r, "/api/v1/netpol/")
	log.Printf("NetworkPolicyHandler url: %v", url)

	netpol, err := c.Clientset.NetworkingV1().NetworkPolicies(url.Namespace).Get(context.TODO(), url.Resource, metav1.GetOptions{})
	if err != nil {
		return S.RespondError(err)
	}
	netpol.ObjectMeta.ManagedFields = nil

	S.RespondYAML(w, http.StatusOK, netpol)
	return nil
}

func NamespaceNetworkPolicyListHandler(app *S.App, c *S.Client, w http.ResponseWriter, r *http.Request) error {
	url := S.GetRequestParams(r, "/api/v1/namespace/netpol/")
	log.Printf("NamespaceNetworkPolicyListHandler url: %v", url)

	g := S.Graph{Nodes: []S.Node{}, Edges: []S.Edge{}}

	ns, err := c.Clientset.CoreV1().Namespaces().Get(context.TODO(), url.Namespace, metav1.GetOptions{})
	if err != nil {
		return S.RespondError(err)
	}
	node := g.AddNode("ns", string(ns.ObjectMeta.UID), ns.ObjectMeta.Name, S.NodeOptions{Type: "namespace", Group: true})

	netpolList, err := c.Clientset.NetworkingV1().NetworkPolicies(url.Namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return S.RespondError(err)
	}
	for _, netpol := range netpolList.Items {
		g.AddNode("netpol", string(netpol.ObjectMeta.UID), netpol.ObjectMeta.Name, S.NodeOptions{Namespace: url.Namespace, Type: "netpol", ParentNode: node, Extent: "parent"})
	}

	return S.RespondJSON(w, http.StatusOK, g)
}
