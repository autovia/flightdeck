// Copyright (c) Autovia GmbH
// SPDX-License-Identifier: Apache-2.0

package k8s

import (
	"log"
	"net/http"

	S "github.com/autovia/flightdeck/api/structs"
)

func ClusterVersion(app *S.App, c *S.Client, w http.ResponseWriter, r *http.Request) error {
	log.Print("ClusterVersion")

	versionInfo, err := c.Discovery.ServerVersion()
	if err != nil {
		S.RespondError(err)
	}

	return S.RespondJSON(w, http.StatusOK, versionInfo)
}
