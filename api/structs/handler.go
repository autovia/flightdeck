// Copyright (c) Autovia GmbH
// SPDX-License-Identifier: Apache-2.0

package structs

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"gopkg.in/yaml.v2"
	clientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/client-go/kubernetes"
)

type Public struct {
	*App
	H func(e *App, w http.ResponseWriter, r *http.Request) error
}

func (p Public) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("PublicHandler %v", r.URL)
	type contextKey string
	u := contextKey("user")
	ctx := context.WithValue(r.Context(), u, "test")
	newReq := r.WithContext(ctx)

	err := p.H(p.App, w, newReq)
	if err != nil {
		log.Print(err)
		HandleError(w, err)
		return
	}
}

type Authorization struct {
	*App
	H func(e *App, w http.ResponseWriter, r *http.Request) error
}

func (auth Authorization) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("AuthorizationHandler %v", r.URL)

	tokenIsValid, token, err := auth.App.AuthManager.Authorize(r)
	if tokenIsValid && len(token) > 0 {
		ctx := context.WithValue(r.Context(), "token", token)
		err = auth.H(auth.App, w, r.WithContext(ctx))
		if err != nil {
			log.Print(err)
			HandleError(w, err)
		}
		return
	}
	HandleError(w, RespondError(http.StatusUnauthorized, err))
}

type KubeClient struct {
	*App
	H func(k *kubernetes.Clientset, w http.ResponseWriter, r *http.Request) error
}

func (kc KubeClient) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("KubeClient %v", r.URL)

	tokenIsValid, token, err := kc.App.AuthManager.Authorize(r)
	if tokenIsValid && len(token) > 0 {
		kubeclient, err := kc.App.NewKubeClient(token)
		if err != nil {
			log.Print(err)
			HandleError(w, err)
		}
		err = kc.H(kubeclient, w, r)
		if err != nil {
			log.Print(err)
			HandleError(w, err)
		}
		return
	}
	HandleError(w, RespondError(http.StatusUnauthorized, err))
}

type ApiClient struct {
	*App
	H func(k *clientset.Clientset, w http.ResponseWriter, r *http.Request) error
}

func (ac ApiClient) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("ApiClient %v", r.URL)

	tokenIsValid, token, err := ac.App.AuthManager.Authorize(r)
	if tokenIsValid && len(token) > 0 {
		apiclient, err := ac.App.NewApiClient(token)
		if err != nil {
			log.Print(err)
			HandleError(w, err)
		}
		err = ac.H(apiclient, w, r)
		if err != nil {
			log.Print(err)
			HandleError(w, err)
		}
		return
	}
	HandleError(w, RespondError(http.StatusUnauthorized, err))
}

func HandleError(w http.ResponseWriter, err error) {
	switch e := err.(type) {
	case StatusError:
		log.Printf("StatusError HTTP %d - %s - %s", e.Status(), e, e.Message())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(e.Status())
		if e.Message() != "" {
			w.Write([]byte(fmt.Sprintf("{\"error\":\"%s\"}", e.Message())))
		} else {
			//http.Error(w, e.Error(), e.Status())
			w.Write([]byte(fmt.Sprintf("{\"error\":\"%s\"}", e)))
		}
	case Error:
		log.Printf("Error HTTP %d - %s", e.Status(), e)
		http.Error(w, e.Error(), e.Status())
	default:
		log.Printf(fmt.Sprintf("default Error HTTP %v", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
	}
}

type Url struct {
	Namespace   string
	Resource    string
	Subresource string
}

func GetRequestParams(r *http.Request, path string) Url {
	request := strings.Split(strings.Replace(r.URL.Path, path, "", -1), "/")
	log.Printf("getRequestParams url: %v", request)
	u := Url{}
	switch {
	case len(request) == 1:
		u.Namespace = request[0]
	case len(request) == 2:
		u.Namespace = request[0]
		u.Resource = request[1]
	case len(request) == 3:
		u.Namespace = request[0]
		u.Resource = request[1]
		u.Subresource = request[2]
	}
	return u
}

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

func RespondJSON(w http.ResponseWriter, status int, payload interface{}) error {
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

func CheckIfError(err error) {
	if err == nil {
		return
	}

	log.Printf("\x1b[31;1m%s\x1b[0m\n", fmt.Sprintf("error: %s", err))
}

func Info(format string, args ...interface{}) {
	log.Printf("\x1b[34;1m%s\x1b[0m\n", fmt.Sprintf(format, args...))
}

func Warning(format string, args ...interface{}) {
	log.Printf("\x1b[36;1m%s\x1b[0m\n", fmt.Sprintf(format, args...))
}
