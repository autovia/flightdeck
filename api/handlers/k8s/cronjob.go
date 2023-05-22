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

func CronJobHandler(client *kubernetes.Clientset, w http.ResponseWriter, r *http.Request) error {
	url := S.GetRequestParams(r, "/api/v1/cronjob/")
	log.Printf("CronJobHandler url: %v", url)

	cronjob, err := client.BatchV1().CronJobs(url.Namespace).Get(context.TODO(), url.Resource, metav1.GetOptions{})
	if err != nil {
		return S.RespondError(err)
	}
	cronjob.ObjectMeta.ManagedFields = nil

	S.RespondYAML(w, http.StatusOK, cronjob)
	return nil
}

func CronJobPodListHandler(client *kubernetes.Clientset, w http.ResponseWriter, r *http.Request) error {
	url := S.GetRequestParams(r, "/api/v1/graph/cronjob/")
	log.Printf("CronJobPodListHandler url: %v", url)

	g := S.Graph{Nodes: []S.Node{}, Edges: []S.Edge{}}
	cjnode := g.AddNode("cronjob", url.Resource, url.Resource, S.NodeOptions{Namespace: url.Namespace, Type: "cronjob"})

	podList, err := client.CoreV1().Pods(url.Namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return S.RespondError(err)
	}

	for _, pod := range podList.Items {
		for _, podOwnerRefs := range pod.ObjectMeta.OwnerReferences {
			if podOwnerRefs.Kind == "Job" {
				job, err := client.BatchV1().Jobs(pod.Namespace).Get(context.Background(), podOwnerRefs.Name, metav1.GetOptions{})
				if err != nil {
					return S.RespondError(err)
				}
				for _, jobOwnerRefs := range job.ObjectMeta.OwnerReferences {
					if jobOwnerRefs.Kind == "CronJob" {
						cronjob, err := client.BatchV1().CronJobs(pod.Namespace).Get(context.Background(), jobOwnerRefs.Name, metav1.GetOptions{})
						if err != nil {
							return S.RespondError(err)
						}
						if url.Resource == cronjob.ObjectMeta.Name {
							if !g.Includes(pod.ObjectMeta.Name) {
								podnode := g.AddNode("pod", string(pod.ObjectMeta.UID), pod.ObjectMeta.Name, S.NodeOptions{Namespace: url.Namespace, Type: "pod"})
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

func NamespaceCronJobListHandler(client *kubernetes.Clientset, w http.ResponseWriter, r *http.Request) error {
	url := S.GetRequestParams(r, "/api/v1/namespace/cronjob/")
	log.Printf("NamespaceCronJobListHandler url: %v", url)

	g := S.Graph{Nodes: []S.Node{}, Edges: []S.Edge{}}

	ns, err := client.CoreV1().Namespaces().Get(context.TODO(), url.Namespace, metav1.GetOptions{})
	if err != nil {
		return S.RespondError(err)
	}
	node := g.AddNode("ns", string(ns.ObjectMeta.UID), ns.ObjectMeta.Name, S.NodeOptions{Type: "namespace", Group: true})

	cronjobList, err := client.BatchV1().CronJobs(url.Namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return S.RespondError(err)
	}
	for _, cronjob := range cronjobList.Items {
		g.AddNode("cronjob", string(cronjob.ObjectMeta.UID), cronjob.ObjectMeta.Name, S.NodeOptions{Namespace: url.Namespace, Type: "cronjob", ParentNode: node, Extent: "parent"})
	}

	return S.RespondJSON(w, http.StatusOK, g)
}