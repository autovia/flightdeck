// Copyright (c) Autovia GmbH
// SPDX-License-Identifier: Apache-2.0

package structs

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"gopkg.in/yaml.v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/types"
)

func RespondOK(w http.ResponseWriter, response string) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
	return nil
}

func RespondError(args ...interface{}) StatusError {
	code := 500
	err := errors.New("Internal Server Error")
	msg := ""

	for _, arg := range args {
		switch t := arg.(type) {
		case string:
			msg = t
		case int:
			code = t
		case error:
			err = t
		}
	}
	return StatusError{Code: code, Err: err, Msg: msg}
}

func RespondText(w http.ResponseWriter, status int, payload string) error {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(status)
	w.Write([]byte(payload))
	return nil
}

func RespondFormat[T any](r *http.Request, w http.ResponseWriter, status int, payload T) error {
	if r.URL.Query().Get("format") == "json" {
		return RespondJSON(w, http.StatusOK, payload)
	}
	return RespondYAML(w, http.StatusOK, payload)
}

func RespondYAML[T any](w http.ResponseWriter, status int, payload T) error {
	yaml, err := yaml.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return nil
	}
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(status)
	w.Write(yaml)
	return nil
}

type ObjectList struct {
	UID               types.UID   `json:"uid"`
	Name              string      `json:"name"`
	Namespace         string      `json:"namespace"`
	Group             string      `json:"group"`
	Version           string      `json:"version"`
	CreationTimestamp metav1.Time `json:"creationTimestamp"`
}

func RespondFilter(r *http.Request, w http.ResponseWriter, status int, list *unstructured.UnstructuredList) error {
	filter := r.URL.Query().Get("filter")

	var ret []ObjectList
	for _, item := range list.Items {
		if strings.Contains(item.GetName(), filter) {
			ret = append(ret, ObjectList{
				UID:               item.GetUID(),
				Name:              item.GetName(),
				Namespace:         item.GetNamespace(),
				Group:             item.GroupVersionKind().Group,
				Version:           item.GroupVersionKind().Version,
				CreationTimestamp: item.GetCreationTimestamp(),
			})
		}
	}

	return RespondJSON(w, http.StatusOK, ret)
}

func RespondJSON[T any](w http.ResponseWriter, status int, payload T) error {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return nil
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte(response))
	return nil
}
