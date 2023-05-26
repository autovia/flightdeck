// Copyright (c) Autovia GmbH
// SPDX-License-Identifier: Apache-2.0

package structs

import (
	"bytes"
	"context"
	"os"

	v1 "k8s.io/api/core/v1"
	clientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
)

type Client struct {
	Clientset *kubernetes.Clientset
	ApiClient *clientset.Clientset
	Config    *rest.Config
}

func (c *Client) GetContainer(url Url) (string, error) {
	var containerName string
	if url.Subresource != "" {
		containerName = url.Subresource
	} else {
		pod, err := c.Clientset.CoreV1().Pods(url.Namespace).Get(context.TODO(), url.Resource, metav1.GetOptions{})
		if err != nil {
			return "", err
		}
		if len(pod.Spec.Containers) > 0 {
			containerName = pod.Spec.Containers[0].Name
		}
	}
	return containerName, nil
}

func (c *Client) ExecCommand(url Url, container string, cmd []string) (*bytes.Buffer, error) {
	req := c.Clientset.CoreV1().RESTClient().Post().Resource("pods").Name(url.Resource).Namespace(url.Namespace).SubResource("exec")
	req.VersionedParams(&v1.PodExecOptions{
		Container: container,
		Command:   cmd,
		Stdin:     true,
		Stdout:    true,
	}, scheme.ParameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(c.Config, "POST", req.URL())
	if err != nil {
		return nil, err
	}
	buf := new(bytes.Buffer)
	err = exec.StreamWithContext(context.TODO(), remotecommand.StreamOptions{
		Stdin:  os.Stdin,
		Stdout: buf,
	})
	if err != nil {
		return nil, err
	}
	return buf, nil
}
