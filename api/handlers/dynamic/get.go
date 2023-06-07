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
)

func GetClusterResourcesHandler(app *S.App, c *S.Client, w http.ResponseWriter, r *http.Request) error {
	url := S.GetRequestParams(r, "/api/v1/get/")
	log.Printf("GetClusterResourcesHandler url: %v", url)

	obj, err := c.Dynamic.Resource(utils.GVR[url.Scope]).Get(context.TODO(), url.Resource, metav1.GetOptions{})
	if err != nil {
		return S.RespondError(err)
	}

	return S.RespondJSON(w, http.StatusOK, obj)
}

func GetNamespaceResourcesHandler(app *S.App, c *S.Client, w http.ResponseWriter, r *http.Request) error {
	url := S.GetRequestParams(r, "/api/v1/resource/")
	log.Printf("GetNamespaceResourcesHandler url: %v", url)

	obj, err := c.Dynamic.Resource(utils.GVR[url.Resource]).Namespace(url.Scope).Get(context.TODO(), url.Resource, metav1.GetOptions{})
	if err != nil {
		return S.RespondError(err)
	}

	return S.RespondJSON(w, http.StatusOK, obj)
}
