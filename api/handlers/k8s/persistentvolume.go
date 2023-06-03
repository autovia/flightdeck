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

func PersistentVolumeHandler(app *S.App, c *S.Client, w http.ResponseWriter, r *http.Request) error {
	url := S.GetRequestParams(r, "/api/v1/pv/")
	log.Printf("PersistentVolumeHandler url: %v", url)

	pv, err := c.Clientset.CoreV1().PersistentVolumes().Get(context.TODO(), url.Resource, metav1.GetOptions{})
	if err != nil {
		return S.RespondError(err)
	}
	pv.ObjectMeta.ManagedFields = nil

	return S.RespondFormat(r, w, http.StatusOK, pv)
}

func PersistentVolumeListHandler(app *S.App, c *S.Client, w http.ResponseWriter, r *http.Request) error {
	log.Print("PersistentVolumeListHandler")

	pvList, err := c.Clientset.CoreV1().PersistentVolumes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return S.RespondError(err)
	}

	return S.RespondJSON(w, http.StatusOK, pvList.Items)
}
