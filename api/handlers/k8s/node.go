// Copyright (c) Autovia GmbH
// SPDX-License-Identifier: Apache-2.0

package k8s

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	S "github.com/autovia/flightdeck/api/structs"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func NodeListHandler(app *S.App, c *S.Client, w http.ResponseWriter, r *http.Request) error {
	log.Print("NodeListHandler")

	g := S.Graph{Nodes: []S.Node{}, Edges: []S.Edge{}}

	nodeList, err := c.Clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return S.RespondError(err)
	}
	for _, node := range nodeList.Items {
		g.AddNode("node", string(node.ObjectMeta.UID), node.ObjectMeta.Name, S.NodeOptions{Type: "node", Draggable: false, Connectable: false})
	}
	return S.RespondJSON(w, http.StatusOK, g)
}

func NodeHandler(app *S.App, c *S.Client, w http.ResponseWriter, r *http.Request) error {
	url := strings.Split(strings.Replace(r.URL.Path, "/api/v1/node/", "", -1), "/")
	log.Printf("NodeHandler url: %v", url)
	resource := url[0]

	g := S.Graph{Nodes: []S.Node{}, Edges: []S.Edge{}}

	no, err := c.Clientset.CoreV1().Nodes().Get(context.TODO(), resource, metav1.GetOptions{})
	if err != nil {
		return S.RespondError(err)
	}
	node := g.AddNode("node", string(no.ObjectMeta.UID), no.ObjectMeta.Name, S.NodeOptions{Type: "node", Group: true})

	podList, err := c.Clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{
		FieldSelector: fmt.Sprintf("spec.nodeName=%s,status.phase=Running", no.Name),
	})
	if err != nil {
		return S.RespondError(err)
	}
	for _, pod := range podList.Items {
		g.AddNode("pod", string(pod.ObjectMeta.UID), pod.ObjectMeta.Name, S.NodeOptions{Namespace: pod.ObjectMeta.Namespace, Type: "pod", ParentNode: node, Extent: "parent"})
	}

	return S.RespondJSON(w, http.StatusOK, g)
}
