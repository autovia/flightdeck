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

func DeploymentHandler(app *S.App, c *S.Client, w http.ResponseWriter, r *http.Request) error {
	url := S.GetRequestParams(r, "/api/v1/deploy/")
	log.Printf("DeploymentHandler url: %v", url)

	deploy, err := c.Clientset.AppsV1().Deployments(url.Scope).Get(context.TODO(), url.Resource, metav1.GetOptions{})
	if err != nil {
		return S.RespondError(err)
	}
	deploy.ObjectMeta.ManagedFields = nil

	return S.RespondFormat(r, w, http.StatusOK, deploy)
}

func DeploymentPodListHandler(app *S.App, c *S.Client, w http.ResponseWriter, r *http.Request) error {
	url := S.GetRequestParams(r, "/api/v1/graph/deploy/")
	log.Printf("DeploymentPodListHandler url: %v", url)

	g := S.Graph{Nodes: []S.Node{}, Edges: []S.Edge{}}
	deploynode := g.AddNode("deploy", url.Resource, url.Resource, S.NodeOptions{Namespace: url.Scope, Type: "deploy"})

	podList, err := c.Clientset.CoreV1().Pods(url.Scope).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return S.RespondError(err)
	}
	for _, pod := range podList.Items {
		for _, podOwnerRefs := range pod.ObjectMeta.OwnerReferences {
			if podOwnerRefs.Kind == "ReplicaSet" {
				replicaset, err := c.Clientset.AppsV1().ReplicaSets(pod.Namespace).Get(context.Background(), podOwnerRefs.Name, metav1.GetOptions{})
				if err != nil {
					return S.RespondError(err)
				}
				for _, replOwnerRefs := range replicaset.ObjectMeta.OwnerReferences {
					if replOwnerRefs.Kind == "Deployment" {
						replDeployment, err := c.Clientset.AppsV1().Deployments(pod.Namespace).Get(context.Background(), replOwnerRefs.Name, metav1.GetOptions{})
						if err != nil {
							return S.RespondError(err)
						}
						if url.Resource == replDeployment.ObjectMeta.Name {
							if !g.Includes(pod.ObjectMeta.Name) {
								podnode := g.AddNode("pod", string(pod.ObjectMeta.UID), pod.ObjectMeta.Name, S.NodeOptions{Namespace: url.Scope, Type: "pod"})
								g.AddEdge(deploynode, podnode)
							}
						}
					}
				}
			}
			if podOwnerRefs.Kind == "Deployment" {
				deploy, err := c.Clientset.AppsV1().Deployments(pod.Namespace).Get(context.Background(), podOwnerRefs.Name, metav1.GetOptions{})
				if err != nil {
					return S.RespondError(err)
				}
				if url.Resource == deploy.ObjectMeta.Name {
					if !g.Includes(pod.ObjectMeta.Name) {
						podnode := g.AddNode("pod", string(pod.ObjectMeta.UID), pod.ObjectMeta.Name, S.NodeOptions{Namespace: url.Scope, Type: "pod"})
						g.AddEdge(deploynode, podnode)
					}
				}
			}
		}
	}

	return S.RespondJSON(w, http.StatusOK, g)
}

func NamespaceDeploymentListHandler(app *S.App, c *S.Client, w http.ResponseWriter, r *http.Request) error {
	url := S.GetRequestParams(r, "/api/v1/namespace/deploy/")
	log.Printf("NamespaceDeploymentListHandler url: %v", url)

	g := S.Graph{Nodes: []S.Node{}, Edges: []S.Edge{}}

	ns, err := c.Clientset.CoreV1().Namespaces().Get(context.TODO(), url.Scope, metav1.GetOptions{})
	if err != nil {
		return S.RespondError(err)
	}
	node := g.AddNode("ns", string(ns.ObjectMeta.UID), ns.ObjectMeta.Name, S.NodeOptions{Type: "namespace", Group: true})

	deployList, err := c.Clientset.AppsV1().Deployments(url.Scope).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return S.RespondError(err)
	}
	for _, deploy := range deployList.Items {
		g.AddNode("deploy", string(deploy.ObjectMeta.UID), deploy.ObjectMeta.Name, S.NodeOptions{Namespace: url.Scope, Type: "deploy", ParentNode: node, Extent: "parent"})
	}

	return S.RespondJSON(w, http.StatusOK, g)
}
