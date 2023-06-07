// Copyright (c) Autovia GmbH
// SPDX-License-Identifier: Apache-2.0

package structs

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strings"

	v1 "k8s.io/api/core/v1"
	clientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
)

type Client struct {
	Clientset *kubernetes.Clientset
	ApiClient *clientset.Clientset
	Config    *rest.Config
	Dynamic   *dynamic.DynamicClient
	Discovery *discovery.DiscoveryClient
}

func (c *Client) GetContainer(url RequestUrl) (string, error) {
	var containerName string
	if url.Subresource != "" {
		containerName = url.Subresource
	} else {
		pod, err := c.Clientset.CoreV1().Pods(url.Scope).Get(context.TODO(), url.Resource, metav1.GetOptions{})
		if err != nil {
			return "", err
		}
		if len(pod.Spec.Containers) > 0 {
			containerName = pod.Spec.Containers[0].Name
		}
	}
	return containerName, nil
}

func (c *Client) ExecCommand(url RequestUrl, container string, cmd []string) (*bytes.Buffer, error) {
	req := c.Clientset.CoreV1().RESTClient().Post().Resource("pods").Name(url.Resource).Namespace(url.Scope).SubResource("exec")
	req.VersionedParams(&v1.PodExecOptions{
		Container: container,
		Command:   cmd,
		Stdin:     true,
		Stdout:    true,
		Stderr:    true,
	}, scheme.ParameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(c.Config, "POST", req.URL())
	if err != nil {
		return nil, err
	}
	outBuf := new(bytes.Buffer)
	errBuf := new(bytes.Buffer)
	err = exec.StreamWithContext(context.TODO(), remotecommand.StreamOptions{
		Stdin:  os.Stdin,
		Stdout: outBuf,
		Stderr: errBuf,
	})
	if err != nil {
		return nil, fmt.Errorf("can not exec cmd: %s", errBuf.String())
	}
	return outBuf, nil
}

func (c *Client) GetAPIResource(obj *unstructured.Unstructured) (*metav1.APIResource, error) {
	resourceList, err := c.Discovery.ServerResourcesForGroupVersion(obj.GetAPIVersion())
	if err != nil {
		return nil, err
	}
	resources := resourceList.APIResources
	var resource *metav1.APIResource
	for _, apiResource := range resources {
		if apiResource.Kind == obj.GetKind() && !strings.Contains(apiResource.Name, "/") {
			resource = &apiResource
			break
		}
	}
	if resource == nil {
		return nil, fmt.Errorf("unknown resource kind: %s", obj.GetKind())
	}
	return resource, nil
}
