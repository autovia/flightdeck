// Copyright (c) Autovia GmbH
// SPDX-License-Identifier: Apache-2.0

package structs

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/autovia/flightdeck/api/utils"
	jwt "github.com/golang-jwt/jwt/v5"
	corev1 "k8s.io/api/core/v1"
	apierr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const (
	namespace   = "flightdeck"
	tokenSecret = "flightdeck-token"
	tokenKey    = "key"
)

type AuthManager interface {
	Authorize(r *http.Request) (bool, string, error)
	Secret() string
}

type authManager struct {
	secret string
	client *kubernetes.Clientset
}

type JwtCustomClaims struct {
	Username string `json:"username"`
	Token    string `json:"token"`
	jwt.RegisteredClaims
}

func (am *authManager) init() {
	secret, err := am.client.CoreV1().Secrets(namespace).Get(context.TODO(), tokenSecret, metav1.GetOptions{})
	if apierr.IsNotFound(err) {
		newToken, err := utils.RandomString()
		if err != nil {
			panic(err)
		}
		_, err = am.client.CoreV1().Secrets(namespace).Create(context.TODO(), &corev1.Secret{ObjectMeta: metav1.ObjectMeta{Name: tokenSecret}, StringData: map[string]string{tokenKey: newToken}}, metav1.CreateOptions{})
		if err != nil {
			panic(err)
		}
		am.secret = newToken
	} else {
		token := string(secret.Data[tokenKey])
		if len(token) == 0 {
			newToken, err := utils.RandomString()
			if err != nil {
				panic(err)
			}
			secret.StringData = map[string]string{tokenKey: newToken}
			_, err = am.client.CoreV1().Secrets(namespace).Update(context.TODO(), secret, metav1.UpdateOptions{})
			if err != nil {
				panic(err)
			}
			am.secret = newToken
		} else {
			am.secret = token
		}
	}
}

func (am *authManager) Authorize(r *http.Request) (bool, string, error) {
	cookie, err := r.Cookie("jwt")
	if err != nil {
		return false, "", errors.New("Authorization token missing")
	}
	token, err := jwt.ParseWithClaims(cookie.Value, &JwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(am.secret), nil
	})
	if claims, ok := token.Claims.(*JwtCustomClaims); ok && token.Valid {
		return true, claims.Token, nil
	}
	return false, "", err
}

func (am *authManager) Secret() string {
	return am.secret
}

func NewAuthManager(k8s *kubernetes.Clientset) AuthManager {
	am := &authManager{client: k8s}
	am.init()
	return am
}
