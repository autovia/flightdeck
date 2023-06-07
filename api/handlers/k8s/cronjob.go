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

func CronJobHandler(app *S.App, c *S.Client, w http.ResponseWriter, r *http.Request) error {
	url := S.GetRequestParams(r, "/api/v1/cronjob/")
	log.Printf("CronJobHandler url: %v", url)

	cronjob, err := c.Clientset.BatchV1().CronJobs(url.Scope).Get(context.TODO(), url.Resource, metav1.GetOptions{})
	if err != nil {
		return S.RespondError(err)
	}
	cronjob.ObjectMeta.ManagedFields = nil

	return S.RespondFormat(r, w, http.StatusOK, cronjob)
}

func CronJobPodListHandler(app *S.App, c *S.Client, w http.ResponseWriter, r *http.Request) error {
	url := S.GetRequestParams(r, "/api/v1/graph/cronjob/")
	log.Printf("CronJobPodListHandler url: %v", url)

	g := S.Graph{Nodes: []S.Node{}, Edges: []S.Edge{}}
	cjnode := g.AddNode("cronjob", url.Resource, url.Resource, S.NodeOptions{Namespace: url.Scope, Type: "cronjob"})

	podList, err := c.Clientset.CoreV1().Pods(url.Scope).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return S.RespondError(err)
	}

	for _, pod := range podList.Items {
		for _, podOwnerRefs := range pod.ObjectMeta.OwnerReferences {
			if podOwnerRefs.Kind == "Job" {
				job, err := c.Clientset.BatchV1().Jobs(pod.Namespace).Get(context.Background(), podOwnerRefs.Name, metav1.GetOptions{})
				if err != nil {
					return S.RespondError(err)
				}
				for _, jobOwnerRefs := range job.ObjectMeta.OwnerReferences {
					if jobOwnerRefs.Kind == "CronJob" {
						cronjob, err := c.Clientset.BatchV1().CronJobs(pod.Namespace).Get(context.Background(), jobOwnerRefs.Name, metav1.GetOptions{})
						if err != nil {
							return S.RespondError(err)
						}
						if url.Resource == cronjob.ObjectMeta.Name {
							if !g.Includes(pod.ObjectMeta.Name) {
								podnode := g.AddNode("pod", string(pod.ObjectMeta.UID), pod.ObjectMeta.Name, S.NodeOptions{Namespace: url.Scope, Type: "pod"})
								g.AddEdge(cjnode, podnode)
							}
						}
					}
				}
			}
		}
	}

	return S.RespondJSON(w, http.StatusOK, g)
}

func NamespaceCronJobListHandler(app *S.App, c *S.Client, w http.ResponseWriter, r *http.Request) error {
	url := S.GetRequestParams(r, "/api/v1/namespace/cronjob/")
	log.Printf("NamespaceCronJobListHandler url: %v", url)

	g := S.Graph{Nodes: []S.Node{}, Edges: []S.Edge{}}

	ns, err := c.Clientset.CoreV1().Namespaces().Get(context.TODO(), url.Scope, metav1.GetOptions{})
	if err != nil {
		return S.RespondError(err)
	}
	node := g.AddNode("ns", string(ns.ObjectMeta.UID), ns.ObjectMeta.Name, S.NodeOptions{Type: "namespace", Group: true})

	cronjobList, err := c.Clientset.BatchV1().CronJobs(url.Scope).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return S.RespondError(err)
	}
	for _, cronjob := range cronjobList.Items {
		g.AddNode("cronjob", string(cronjob.ObjectMeta.UID), cronjob.ObjectMeta.Name, S.NodeOptions{Namespace: url.Scope, Type: "cronjob", ParentNode: node, Extent: "parent"})
	}

	return S.RespondJSON(w, http.StatusOK, g)
}
