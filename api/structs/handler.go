// Copyright (c) Autovia GmbH
// SPDX-License-Identifier: Apache-2.0

package structs

import (
	"context"
	"log"
	"net/http"
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
	H func(a *App, w http.ResponseWriter, r *http.Request) error
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
	H func(a *App, c *Client, w http.ResponseWriter, r *http.Request) error
}

func (kc KubeClient) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("KubeClient %v", r.URL)

	tokenIsValid, token, err := kc.App.AuthManager.Authorize(r)
	if tokenIsValid && len(token) > 0 {
		client, err := kc.App.NewKubeClient(token)
		if err != nil {
			log.Print(err)
			HandleError(w, err)
		}
		err = kc.H(kc.App, client, w, r)
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
	H func(a *App, c *Client, w http.ResponseWriter, r *http.Request) error
}

func (ac ApiClient) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("ApiClient %v", r.URL)

	tokenIsValid, token, err := ac.App.AuthManager.Authorize(r)
	if tokenIsValid && len(token) > 0 {
		client, err := ac.App.NewApiClient(token)
		if err != nil {
			log.Print(err)
			HandleError(w, err)
		}
		err = ac.H(ac.App, client, w, r)
		if err != nil {
			log.Print(err)
			HandleError(w, err)
		}
		return
	}
	HandleError(w, RespondError(http.StatusUnauthorized, err))
}
