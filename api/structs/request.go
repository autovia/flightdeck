// Copyright (c) Autovia GmbH
// SPDX-License-Identifier: Apache-2.0

package structs

import (
	"log"
	"net/http"
	"strings"
)

type RequestUrl struct {
	Scope       string
	Resource    string
	Subresource string
	Path        string
}

func GetRequestParams(r *http.Request, path string) RequestUrl {
	request := strings.Split(strings.Replace(r.URL.Path, path, "", -1), "/")
	log.Printf("getRequestParams url: %v", request)
	u := RequestUrl{}
	switch {
	case len(request) == 1:
		u.Scope = request[0]
	case len(request) == 2:
		u.Scope = request[0]
		u.Resource = request[1]
	case len(request) == 3:
		u.Scope = request[0]
		u.Resource = request[1]
		u.Subresource = request[2]
	case len(request) > 3:
		u.Scope = request[0]
		u.Resource = request[1]
		u.Subresource = request[2]
		u.Path = strings.Join(request[3:], "/")
	}
	return u
}
