// Copyright (c) Autovia GmbH
// SPDX-License-Identifier: Apache-2.0

package k8s

import (
	"context"
	"log"
	"net/http"

	S "github.com/autovia/flightdeck/api/structs"

	"k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func CustomResourceDefinitionHandler(app *S.App, apiclient *clientset.Clientset, w http.ResponseWriter, r *http.Request) error {
	url := S.GetRequestParams(r, "/api/v1/crd/")
	log.Printf("CustomResourceDefinitionHandler url: %v", url)

	// role, err := app.Client.ApiextensionsV1().RESTClient().CustomResourceDefinitions().Get(context.TODO(), url.Resource, metav1.GetOptions{})
	// clientset.ApiextensionsV1beta1().CustomResourceDefinitions().List(context.TODO(), metav1.ListOptions{})

	//crd, err := clientset.A ApiextensionsV1().CustomResourceDefinitions().Get(context.TODO(), url.Resource, metav1.GetOptions{})
	crd, err := apiclient.ApiextensionsV1().CustomResourceDefinitions().Get(context.TODO(), url.Resource, metav1.GetOptions{})
	if err != nil {
		return S.RespondError(err)
	}
	crd.ObjectMeta.ManagedFields = nil

	return S.RespondYAML(w, http.StatusOK, crd)
}

func CustomResourceDefinitionListHandler(app *S.App, apiclient *clientset.Clientset, w http.ResponseWriter, r *http.Request) error {
	log.Print("CustomResourceDefinitionListHandler")

	g := S.Graph{Nodes: []S.Node{}, Edges: []S.Edge{}}

	crdList, err := apiclient.ApiextensionsV1().CustomResourceDefinitions().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return S.RespondError(err)
	}
	for _, cr := range crdList.Items {
		g.AddNode("crd", string(cr.ObjectMeta.UID), cr.ObjectMeta.Name, S.NodeOptions{Type: "crd", Draggable: false, Connectable: false})
	}
	return S.RespondJSON(w, http.StatusOK, g)
}
