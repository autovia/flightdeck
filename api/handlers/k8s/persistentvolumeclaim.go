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

func PersistentVolumeClaim(app *S.App, c *S.Client, w http.ResponseWriter, r *http.Request) error {
	url := S.GetRequestParams(r, "/api/v1/pvc/")
	log.Printf("ReplicaSetHandler url: %v", url)

	pvc, err := c.Clientset.CoreV1().PersistentVolumeClaims(url.Namespace).Get(context.TODO(), url.Resource, metav1.GetOptions{})
	if err != nil {
		return S.RespondError(err)
	}
	pvc.ObjectMeta.ManagedFields = nil

	return S.RespondFormat(r, w, http.StatusOK, pvc)
}

func PersistentVolumeClaimPodListHandler(app *S.App, c *S.Client, w http.ResponseWriter, r *http.Request) error {
	url := S.GetRequestParams(r, "/api/v1/graph/pvc/")
	log.Printf("PersistentVolumeClaimPodListHandler url: %v", url)

	g := S.Graph{Nodes: []S.Node{}, Edges: []S.Edge{}}
	pvcnode := g.AddNode("pvc", url.Resource, url.Resource, S.NodeOptions{Namespace: url.Namespace, Type: "pvc"})

	podList, err := c.Clientset.CoreV1().Pods(url.Namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return S.RespondError(err)
	}
	for _, pod := range podList.Items {
		for _, volume := range pod.Spec.Volumes {
			if volume.PersistentVolumeClaim != nil {
				if volume.PersistentVolumeClaim.ClaimName == url.Resource {
					if !g.Includes(pod.ObjectMeta.Name) {
						podnode := g.AddNode("pod", string(pod.ObjectMeta.UID), pod.ObjectMeta.Name, S.NodeOptions{Namespace: url.Namespace, Type: "pod"})
						g.AddEdge(pvcnode, podnode)
					}
				}
			}
		}
	}

	return S.RespondJSON(w, http.StatusOK, g)
}

func NamespacePersistentVolumeClaimListHandler(app *S.App, c *S.Client, w http.ResponseWriter, r *http.Request) error {
	url := S.GetRequestParams(r, "/api/v1/namespace/pvc/")
	log.Printf("NamespacePersistentVolumeClaimListHandler url: %v", url)

	g := S.Graph{Nodes: []S.Node{}, Edges: []S.Edge{}}

	ns, err := c.Clientset.CoreV1().Namespaces().Get(context.TODO(), url.Namespace, metav1.GetOptions{})
	if err != nil {
		return S.RespondError(err)
	}
	node := g.AddNode("ns", string(ns.ObjectMeta.UID), ns.ObjectMeta.Name, S.NodeOptions{Type: "namespace", Group: true})

	pvcList, err := c.Clientset.CoreV1().PersistentVolumeClaims(url.Namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return S.RespondError(err)
	}
	for _, cm := range pvcList.Items {
		g.AddNode("pvc", string(cm.ObjectMeta.UID), cm.ObjectMeta.Name, S.NodeOptions{Namespace: url.Namespace, Type: "pvc", ParentNode: node, Extent: "parent"})
	}

	return S.RespondJSON(w, http.StatusOK, g)
}
