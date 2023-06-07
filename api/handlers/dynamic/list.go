// Copyright (c) Autovia GmbH
// SPDX-License-Identifier: Apache-2.0

package dynamic

import (
	"context"
	"log"
	"net/http"

	S "github.com/autovia/flightdeck/api/structs"
	"github.com/autovia/flightdeck/api/utils"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func ListClusterResourcesHandler(app *S.App, c *S.Client, w http.ResponseWriter, r *http.Request) error {
	url := S.GetRequestParams(r, "/api/v1/list/")
	log.Printf("ListClusterResourcesHandler url: %v", url)

	list, err := c.Dynamic.Resource(utils.GVR[url.Scope]).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return S.RespondError(err)
	}

	return S.RespondFilter(r, w, http.StatusOK, list)
}

func ListNamespaceResourcesHandler(app *S.App, c *S.Client, w http.ResponseWriter, r *http.Request) error {
	url := S.GetRequestParams(r, "/api/v1/resources/")
	log.Printf("ListNamespaceResourcesHandler url: %v", url)

	var list *unstructured.UnstructuredList
	var err error

	if url.Scope == "all" {
		list, err = c.Dynamic.Resource(utils.GVR[url.Resource]).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			return S.RespondError(err)
		}
	} else {
		list, err = c.Dynamic.Resource(utils.GVR[url.Resource]).Namespace(url.Scope).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			return S.RespondError(err)
		}
	}

	return S.RespondGraph(r, w, http.StatusOK, list, &url)
}
