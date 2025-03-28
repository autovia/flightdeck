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

func StatefulSetHandler(app *S.App, c *S.Client, w http.ResponseWriter, r *http.Request) error {
	url := S.GetRequestParams(r, "/api/v1/sts/")
	log.Printf("StatefulSetHandler url: %v", url)

	sts, err := c.Clientset.AppsV1().StatefulSets(url.Scope).Get(context.TODO(), url.Resource, metav1.GetOptions{})
	if err != nil {
		return S.RespondError(err)
	}
	sts.ObjectMeta.ManagedFields = nil

	return S.RespondFormat(r, w, http.StatusOK, sts)
}

func StatefulSetPodListHandler(app *S.App, c *S.Client, w http.ResponseWriter, r *http.Request) error {
	url := S.GetRequestParams(r, "/api/v1/graph/sts/")
	log.Printf("StatefulSetPodListHandler url: %v", url)

	g := S.Graph{Nodes: []S.Node{}, Edges: []S.Edge{}}
	stsnode := g.AddNode("sts", url.Resource, url.Resource, S.NodeOptions{Namespace: url.Scope, Type: "sts"})

	podList, err := c.Clientset.CoreV1().Pods(url.Scope).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return S.RespondError(err)
	}
	for _, pod := range podList.Items {
		for _, podOwnerRefs := range pod.ObjectMeta.OwnerReferences {
			if podOwnerRefs.Kind == "StatefulSet" {

				statefulset, err := c.Clientset.AppsV1().StatefulSets(pod.Namespace).Get(context.Background(), podOwnerRefs.Name, metav1.GetOptions{})
				if err != nil {
					return S.RespondError(err)
				}

				if url.Resource == statefulset.ObjectMeta.Name {
					if !g.Includes(pod.ObjectMeta.Name) {
						podnode := g.AddNode("pod", string(pod.ObjectMeta.UID), pod.ObjectMeta.Name, S.NodeOptions{Namespace: url.Scope, Type: "pod"})
						g.AddEdge(stsnode, podnode, S.EdgeOptions{})
					}
				}
			}
		}
	}

	return S.RespondJSON(w, http.StatusOK, g)
}

func NamespaceStatefulSetListHandler(app *S.App, c *S.Client, w http.ResponseWriter, r *http.Request) error {
	url := S.GetRequestParams(r, "/api/v1/namespace/sts/")
	log.Printf("NamespaceStatefulSetListHandler url: %v", url)

	g := S.Graph{Nodes: []S.Node{}, Edges: []S.Edge{}}

	ns, err := c.Clientset.CoreV1().Namespaces().Get(context.TODO(), url.Scope, metav1.GetOptions{})
	if err != nil {
		return S.RespondError(err)
	}
	node := g.AddNode("ns", string(ns.ObjectMeta.UID), ns.ObjectMeta.Name, S.NodeOptions{Type: "namespace", Group: true})

	stsList, err := c.Clientset.AppsV1().StatefulSets(url.Scope).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return S.RespondError(err)
	}
	for _, sts := range stsList.Items {
		g.AddNode("sts", string(sts.ObjectMeta.UID), sts.ObjectMeta.Name, S.NodeOptions{Namespace: url.Scope, Type: "sts", ParentNode: node, Extent: "parent"})
	}

	return S.RespondJSON(w, http.StatusOK, g)
}
