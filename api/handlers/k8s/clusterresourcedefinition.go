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

func CustomResourceDefinitionHandler(app *S.App, c *S.Client, w http.ResponseWriter, r *http.Request) error {
	url := S.GetRequestParams(r, "/api/v1/crd/")
	log.Printf("CustomResourceDefinitionHandler url: %v", url)

	crd, err := c.ApiClient.ApiextensionsV1().CustomResourceDefinitions().Get(context.TODO(), url.Resource, metav1.GetOptions{})
	if err != nil {
		return S.RespondError(err)
	}
	crd.ObjectMeta.ManagedFields = nil

	return S.RespondFormat(r, w, http.StatusOK, crd)
}

func CustomResourceDefinitionListHandler(app *S.App, c *S.Client, w http.ResponseWriter, r *http.Request) error {
	log.Print("CustomResourceDefinitionListHandler")

	crdList, err := c.ApiClient.ApiextensionsV1().CustomResourceDefinitions().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return S.RespondError(err)
	}

	return S.RespondJSON(w, http.StatusOK, crdList.Items)
}
