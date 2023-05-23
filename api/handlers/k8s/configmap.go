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

func ConfigMapHandler(app *S.App, client *kubernetes.Clientset, w http.ResponseWriter, r *http.Request) error {
	url := S.GetRequestParams(r, "/api/v1/cm/")
	log.Printf("ConfigMapHandler url: %v", url)

	cm, err := client.CoreV1().ConfigMaps(url.Namespace).Get(context.TODO(), url.Resource, metav1.GetOptions{})
	if err != nil {
		return S.RespondError(err)
	}
	cm.ObjectMeta.ManagedFields = nil

	S.RespondYAML(w, http.StatusOK, cm)
	return nil
}

func ConfigMapPodListHandler(app *S.App, client *kubernetes.Clientset, w http.ResponseWriter, r *http.Request) error {
	url := S.GetRequestParams(r, "/api/v1/graph/cm/")
	log.Printf("ConfigMapPodListHandler url: %v", url)

	g := S.Graph{Nodes: []S.Node{}, Edges: []S.Edge{}}
	cmnode := g.AddNode("cm", url.Resource, url.Resource, S.NodeOptions{Namespace: url.Namespace, Type: "cm"})

	podList, err := client.CoreV1().Pods(url.Namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return S.RespondError(err)
	}
	for _, pod := range podList.Items {
		for _, volume := range pod.Spec.Volumes {
			if volume.ConfigMap != nil {
				if volume.ConfigMap.Name == url.Resource {
					if !g.Includes(pod.ObjectMeta.Name) {
						podnode := g.AddNode("pod", string(pod.ObjectMeta.UID), pod.ObjectMeta.Name, S.NodeOptions{Namespace: url.Namespace, Type: "pod"})
						g.AddEdge(cmnode, podnode)
					}
				}
			}
			if volume.Projected != nil {
				for _, source := range volume.Projected.Sources {
					if source.ConfigMap != nil {
						if source.ConfigMap.Name == url.Resource {
							if !g.Includes(pod.ObjectMeta.Name) {
								podnode := g.AddNode("pod", string(pod.ObjectMeta.UID), pod.ObjectMeta.Name, S.NodeOptions{Namespace: url.Namespace, Type: "pod"})
								g.AddEdge(cmnode, podnode)
							}
						}
					}
				}
			}
		}
	}

	return S.RespondJSON(w, http.StatusOK, g)
}

func NamespaceConfigMapListHandler(app *S.App, client *kubernetes.Clientset, w http.ResponseWriter, r *http.Request) error {
	url := S.GetRequestParams(r, "/api/v1/namespace/cm/")
	log.Printf("NamespaceConfigMapListHandler url: %v", url)

	g := S.Graph{Nodes: []S.Node{}, Edges: []S.Edge{}}

	ns, err := client.CoreV1().Namespaces().Get(context.TODO(), url.Namespace, metav1.GetOptions{})
	if err != nil {
		return S.RespondError(err)
	}
	node := g.AddNode("ns", string(ns.ObjectMeta.UID), ns.ObjectMeta.Name, S.NodeOptions{Type: "namespace", Group: true})

	cmList, err := client.CoreV1().ConfigMaps(url.Namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return S.RespondError(err)
	}
	for _, cm := range cmList.Items {
		g.AddNode("cm", string(cm.ObjectMeta.UID), cm.ObjectMeta.Name, S.NodeOptions{Namespace: url.Namespace, Type: "cm", ParentNode: node, Extent: "parent"})
	}

	return S.RespondJSON(w, http.StatusOK, g)
}
