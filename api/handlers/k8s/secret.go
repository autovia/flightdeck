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

func SecretHandler(app *S.App, c *S.Client, w http.ResponseWriter, r *http.Request) error {
	url := S.GetRequestParams(r, "/api/v1/secret/")
	log.Printf("SecretHandler url: %v", url)

	secret, err := c.Clientset.CoreV1().Secrets(url.Namespace).Get(context.TODO(), url.Resource, metav1.GetOptions{})
	if err != nil {
		return S.RespondError(err)
	}
	secret.ObjectMeta.ManagedFields = nil

	return S.RespondYAML(w, http.StatusOK, secret)
}

func SecretPodListHandler(app *S.App, c *S.Client, w http.ResponseWriter, r *http.Request) error {
	url := S.GetRequestParams(r, "/api/v1/graph/secret/")
	log.Printf("SecretHandler url: %v", url)

	g := S.Graph{Nodes: []S.Node{}, Edges: []S.Edge{}}
	secnode := g.AddNode("secret", url.Resource, url.Resource, S.NodeOptions{Namespace: url.Namespace, Type: "secret"})

	podList, err := c.Clientset.CoreV1().Pods(url.Namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return S.RespondError(err)
	}
	for _, pod := range podList.Items {
		for _, volume := range pod.Spec.Volumes {
			if volume.Secret != nil {
				if volume.Secret.SecretName == url.Resource {
					if !g.Includes(pod.ObjectMeta.Name) {
						podnode := g.AddNode("pod", string(pod.ObjectMeta.UID), pod.ObjectMeta.Name, S.NodeOptions{Namespace: url.Namespace, Type: "pod"})
						g.AddEdge(secnode, podnode)
					}
				}
			}
			if volume.Projected != nil {
				for _, source := range volume.Projected.Sources {
					if source.Secret != nil {
						if source.Secret.Name == url.Resource {
							if !g.Includes(pod.ObjectMeta.Name) {
								podnode := g.AddNode("pod", string(pod.ObjectMeta.UID), pod.ObjectMeta.Name, S.NodeOptions{Namespace: url.Namespace, Type: "pod"})
								g.AddEdge(secnode, podnode)
							}
						}
					}
				}
			}
		}
	}

	return S.RespondJSON(w, http.StatusOK, g)
}

func NamespaceSecretListHandler(app *S.App, c *S.Client, w http.ResponseWriter, r *http.Request) error {
	url := S.GetRequestParams(r, "/api/v1/namespace/secret/")
	log.Printf("NamespaceSecretListHandler url: %v", url)

	g := S.Graph{Nodes: []S.Node{}, Edges: []S.Edge{}}

	ns, err := c.Clientset.CoreV1().Namespaces().Get(context.TODO(), url.Namespace, metav1.GetOptions{})
	if err != nil {
		return S.RespondError(err)
	}
	node := g.AddNode("ns", string(ns.ObjectMeta.UID), ns.ObjectMeta.Name, S.NodeOptions{Type: "namespace", Group: true})

	secList, err := c.Clientset.CoreV1().Secrets(url.Namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return S.RespondError(err)
	}
	for _, cm := range secList.Items {
		g.AddNode("secret", string(cm.ObjectMeta.UID), cm.ObjectMeta.Name, S.NodeOptions{Namespace: url.Namespace, Type: "secret", ParentNode: node, Extent: "parent"})
	}

	return S.RespondJSON(w, http.StatusOK, g)
}
