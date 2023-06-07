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

func DaemonSetHandler(app *S.App, c *S.Client, w http.ResponseWriter, r *http.Request) error {
	url := S.GetRequestParams(r, "/api/v1/ds/")
	log.Printf("DaemonSetHandler url: %v", url)

	ds, err := c.Clientset.AppsV1().DaemonSets(url.Scope).Get(context.TODO(), url.Resource, metav1.GetOptions{})
	if err != nil {
		return S.RespondError(err)
	}
	ds.ObjectMeta.ManagedFields = nil

	S.RespondFormat(r, w, http.StatusOK, ds)
	return nil
}

func DaemonSetPodListHandler(app *S.App, c *S.Client, w http.ResponseWriter, r *http.Request) error {
	url := S.GetRequestParams(r, "/api/v1/graph/ds/")
	log.Printf("DaemonSetPodListHandler url: %v", url)

	g := S.Graph{Nodes: []S.Node{}, Edges: []S.Edge{}}
	dsnode := g.AddNode("ds", url.Resource, url.Resource, S.NodeOptions{Namespace: url.Scope, Type: "ds"})

	podList, err := c.Clientset.CoreV1().Pods(url.Scope).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return S.RespondError(err)
	}
	for _, pod := range podList.Items {
		for _, podOwnerRefs := range pod.ObjectMeta.OwnerReferences {
			if podOwnerRefs.Kind == "DaemonSet" {

				daemonset, err := c.Clientset.AppsV1().DaemonSets(pod.Namespace).Get(context.Background(), podOwnerRefs.Name, metav1.GetOptions{})
				if err != nil {
					return S.RespondError(err)
				}

				if url.Resource == daemonset.ObjectMeta.Name {
					if !g.Includes(pod.ObjectMeta.Name) {
						podnode := g.AddNode("pod", string(pod.ObjectMeta.UID), pod.ObjectMeta.Name, S.NodeOptions{Namespace: url.Scope, Type: "pod"})
						g.AddEdge(dsnode, podnode)
					}
				}
			}
		}
	}

	return S.RespondJSON(w, http.StatusOK, g)
}

func NamespaceDaemonSetListHandler(app *S.App, c *S.Client, w http.ResponseWriter, r *http.Request) error {
	url := S.GetRequestParams(r, "/api/v1/namespace/ds/")
	log.Printf("NamespaceDaemonSetListHandler url: %v", url)

	g := S.Graph{Nodes: []S.Node{}, Edges: []S.Edge{}}

	ns, err := c.Clientset.CoreV1().Namespaces().Get(context.TODO(), url.Scope, metav1.GetOptions{})
	if err != nil {
		return S.RespondError(err)
	}
	node := g.AddNode("ns", string(ns.ObjectMeta.UID), ns.ObjectMeta.Name, S.NodeOptions{Type: "namespace", Group: true})

	dsList, err := c.Clientset.AppsV1().DaemonSets(url.Scope).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return S.RespondError(err)
	}
	for _, ds := range dsList.Items {
		g.AddNode("ds", string(ds.ObjectMeta.UID), ds.ObjectMeta.Name, S.NodeOptions{Namespace: url.Scope, Type: "ds", ParentNode: node, Extent: "parent"})
	}

	return S.RespondJSON(w, http.StatusOK, g)
}
