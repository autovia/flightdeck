// Copyright (c) Autovia GmbH
// SPDX-License-Identifier: Apache-2.0

package auth

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	S "github.com/autovia/flightdeck/api/structs"

	authv1 "k8s.io/api/authentication/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	jwt "github.com/golang-jwt/jwt/v5"
)

func LoginTokenHandler(app *S.App, w http.ResponseWriter, r *http.Request) error {
	authHeader := r.Header.Get("Authorization")
	token := strings.TrimPrefix(authHeader, "Bearer ")

	if len(token) == 0 {
		return S.RespondError("Can not login without a token")
	}

	kubeclient, err := app.NewKubeClient(token)
	if err != nil {
		log.Print(err)
		return S.RespondError("Can not load kube client")
	}

	tokenReview, err := kubeclient.AuthenticationV1().TokenReviews().Create(context.TODO(), &authv1.TokenReview{
		Spec: authv1.TokenReviewSpec{
			Token: token,
		},
	}, metav1.CreateOptions{})
	if err != nil {
		log.Print("Can not request token review")
		return S.RespondError(err)
	}

	expiration := time.Now().Add(2 * time.Hour)

	if tokenReview.Status.Authenticated {
		claims := &S.JwtCustomClaims{
			Username: tokenReview.Status.User.Username,
			Token:    token,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(expiration),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		}
		newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		accessToken, err := newToken.SignedString([]byte(app.AuthManager.Secret()))
		if err != nil {
			log.Print("Can not sign token")
			return S.RespondError(err)
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "jwt",
			Value:    accessToken,
			Expires:  expiration,
			Path:     "/api/v1",
			HttpOnly: true,
		})

		return S.RespondJSON(w, http.StatusOK, "")
	}
	return S.RespondError(err)
}

func ResetTokenHandler(app *S.App, w http.ResponseWriter, r *http.Request) error {
	cookie := &http.Cookie{
		Name:     "jwt",
		Value:    "",
		MaxAge:   -1,
		Path:     "/api/v1",
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)

	return S.RespondJSON(w, http.StatusOK, "valid")
}

func StatusTokenHandler(app *S.App, w http.ResponseWriter, r *http.Request) error {
	token := r.Context().Value("token").(string)

	kubeclient, err := app.NewKubeClient(token)
	if err != nil {
		log.Print("Can not create kube client")
		return S.RespondError(err)
	}

	tokenReview, err := kubeclient.AuthenticationV1().TokenReviews().Create(context.TODO(), &authv1.TokenReview{
		Spec: authv1.TokenReviewSpec{
			Token: token,
		},
	}, metav1.CreateOptions{})
	if err != nil {
		log.Print("Can not request token review")
		return S.RespondError(err)
	}

	if tokenReview.Status.Authenticated {
		return S.RespondJSON(w, http.StatusOK, "valid")
	}

	return S.RespondError("token not valid")
}
