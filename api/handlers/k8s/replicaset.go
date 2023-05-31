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

func ReplicaSetHandler(app *S.App, c *S.Client, w http.ResponseWriter, r *http.Request) error {
	url := S.GetRequestParams(r, "/api/v1/rs/")
	log.Printf("ReplicaSetHandler url: %v", url)

	rs, err := c.Clientset.AppsV1().ReplicaSets(url.Namespace).Get(context.TODO(), url.Resource, metav1.GetOptions{})
	if err != nil {
		return S.RespondError(err)
	}
	rs.ObjectMeta.ManagedFields = nil

	return S.RespondFormat(r, w, http.StatusOK, rs)
}

func ReplicaSetPodListHandler(app *S.App, c *S.Client, w http.ResponseWriter, r *http.Request) error {
	url := S.GetRequestParams(r, "/api/v1/graph/rs/")
	log.Printf("ReplicaSetPodListHandler url: %v", url)

	g := S.Graph{Nodes: []S.Node{}, Edges: []S.Edge{}}
	rsnode := g.AddNode("rs", url.Resource, url.Resource, S.NodeOptions{Namespace: url.Namespace, Type: "rs"})

	podList, err := c.Clientset.CoreV1().Pods(url.Namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return S.RespondError(err)
	}
	for _, pod := range podList.Items {
		for _, podOwnerRefs := range pod.ObjectMeta.OwnerReferences {
			if podOwnerRefs.Kind == "ReplicaSet" {

				statefulset, err := c.Clientset.AppsV1().ReplicaSets(pod.Namespace).Get(context.Background(), podOwnerRefs.Name, metav1.GetOptions{})
				if err != nil {
					return S.RespondError(err)
				}

				if url.Resource == statefulset.ObjectMeta.Name {
					if !g.Includes(pod.ObjectMeta.Name) {
						podnode := g.AddNode("pod", string(pod.ObjectMeta.UID), pod.ObjectMeta.Name, S.NodeOptions{Namespace: url.Namespace, Type: "pod"})
						g.AddEdge(rsnode, podnode)
					}
				}
			}
		}
	}

	return S.RespondJSON(w, http.StatusOK, g)
}

func NamespaceReplicaSetListHandler(app *S.App, c *S.Client, w http.ResponseWriter, r *http.Request) error {
	url := S.GetRequestParams(r, "/api/v1/namespace/rs/")
	log.Printf("NamespaceReplicaSetListHandler url: %v", url)

	g := S.Graph{Nodes: []S.Node{}, Edges: []S.Edge{}}

	ns, err := c.Clientset.CoreV1().Namespaces().Get(context.TODO(), url.Namespace, metav1.GetOptions{})
	if err != nil {
		return S.RespondError(err)
	}
	node := g.AddNode("ns", string(ns.ObjectMeta.UID), ns.ObjectMeta.Name, S.NodeOptions{Type: "namespace", Group: true})

	rsList, err := c.Clientset.AppsV1().ReplicaSets(url.Namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return S.RespondError(err)
	}
	for _, rs := range rsList.Items {
		g.AddNode("rs", string(rs.ObjectMeta.UID), rs.ObjectMeta.Name, S.NodeOptions{Namespace: url.Namespace, Type: "rs", ParentNode: node, Extent: "parent"})
	}

	return S.RespondJSON(w, http.StatusOK, g)
}
