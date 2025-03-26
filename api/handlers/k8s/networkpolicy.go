// Copyright (c) Autovia GmbH
// SPDX-License-Identifier: Apache-2.0

package k8s

import (
	"context"
	"log"
	"net/http"
	"strings"

	S "github.com/autovia/flightdeck/api/structs"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

func NetworkPolicyHandler(app *S.App, c *S.Client, w http.ResponseWriter, r *http.Request) error {
	url := S.GetRequestParams(r, "/api/v1/netpol/")
	log.Printf("NetworkPolicyHandler url: %v", url)

	netpol, err := c.Clientset.NetworkingV1().NetworkPolicies(url.Scope).Get(context.TODO(), url.Resource, metav1.GetOptions{})
	if err != nil {
		return S.RespondError(err)
	}
	netpol.ObjectMeta.ManagedFields = nil

	return S.RespondFormat(r, w, http.StatusOK, netpol)
}

func NetworkPolicyGraphHandler(app *S.App, c *S.Client, w http.ResponseWriter, r *http.Request) error {
	// namespace/netpol
	url := strings.Split(strings.Replace(r.URL.Path, "/api/v1/graph/netpol/", "", -1), "/")
	log.Printf("NetworkPolicyGraphHandler url: %v", url)
	namespace := url[0]
	np := url[1]

	// S.Graph
	g := S.Graph{Nodes: []S.Node{}, Edges: []S.Edge{}, Direction: "RL"}
	var node S.Node

	// netpol
	netpol, err := c.Clientset.NetworkingV1().NetworkPolicies(namespace).Get(context.TODO(), np, metav1.GetOptions{})
	if err != nil {
		return S.RespondError(err)
	}
	netpol.ObjectMeta.ManagedFields = nil

	// PodSelector
	podList, err := c.Clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{
		LabelSelector: labels.Set(netpol.Spec.PodSelector.MatchLabels).String(),
	})
	if err != nil {
		return S.RespondError(err)
	}
	for _, pod := range podList.Items {
		node = g.AddNode("pod", string(pod.ObjectMeta.UID), pod.ObjectMeta.Name, S.NodeOptions{Namespace: pod.Namespace, Type: "podnetpol", Position: "right"})
	}

	for i := range netpol.Spec.Egress {
		for j := range netpol.Spec.Egress[i].To {
			egressNamespaceSelector, err := metav1.LabelSelectorAsSelector(netpol.Spec.Egress[i].To[j].NamespaceSelector)
			if err != nil {
				return S.RespondError("egressNamespaceSelector", err)
			}

			// get list of namespaces
			egressNamespaceList, err := c.Clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{
				LabelSelector: egressNamespaceSelector.String(),
			})
			if err != nil {
				return S.RespondError("egressNamespaceList", err)
			}
			// get pods for each namespace
			for _, egressNamespace := range egressNamespaceList.Items {
				egressPodSelector, err := metav1.LabelSelectorAsSelector(netpol.Spec.Egress[i].To[j].PodSelector)
				if err != nil {
					return S.RespondError("egressPodSelector", err)
				}

				egressPodList, err := c.Clientset.CoreV1().Pods(egressNamespace.Name).List(context.TODO(), metav1.ListOptions{
					LabelSelector: egressPodSelector.String(),
				})
				if err != nil {
					return S.RespondError("egressNamespace", err)
				}
				for _, egressPod := range egressPodList.Items {
					g.AddEdge(node, g.AddNode("pod", string(egressPod.ObjectMeta.UID), egressPod.ObjectMeta.Name, S.NodeOptions{Namespace: egressPod.Namespace, Type: "podnetpol", Position: "left"}), S.EdgeOptions{
						Type:         "netpolegress",
						MarkerStart:  "start",
						MarkerEnd:    "end",
						Animated:     true,
						Data:         "egress-data",
						SourceHandle: "egress",
						TargetHandle: "ingress",
					})
					g.Direction = "LR"
				}
			}
		}
	}

	for i := range netpol.Spec.Ingress {
		for j := range netpol.Spec.Ingress[i].From {
			ingressNamespaceSelector, err := metav1.LabelSelectorAsSelector(netpol.Spec.Ingress[i].From[j].NamespaceSelector)
			if err != nil {
				return S.RespondError("ingressNamespaceSelector", err)
			}

			// get list of namespaces
			ingressNamespaceList, err := c.Clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{
				LabelSelector: ingressNamespaceSelector.String(),
			})
			if err != nil {
				return S.RespondError("ingressNamespaceList", err)
			}
			// get pods for each namespace
			for _, ingressNamespace := range ingressNamespaceList.Items {
				ingressPodSelector, err := metav1.LabelSelectorAsSelector(netpol.Spec.Ingress[i].From[j].PodSelector)
				if err != nil {
					return S.RespondError("ingressPodSelector", err)
				}

				ingressPodList, err := c.Clientset.CoreV1().Pods(ingressNamespace.Name).List(context.TODO(), metav1.ListOptions{
					LabelSelector: ingressPodSelector.String(),
				})
				if err != nil {
					return S.RespondError("ingressNamespace", err)
				}
				for _, ingressPod := range ingressPodList.Items {
					g.AddEdge(
						g.AddNode("pod", string(ingressPod.ObjectMeta.UID), ingressPod.ObjectMeta.Name, S.NodeOptions{Namespace: ingressPod.Namespace, Type: "podnetpol", Position: "left"}), node,
						S.EdgeOptions{
							Type:         "netpolingress",
							MarkerStart:  "start",
							MarkerEnd:    "end",
							Animated:     true,
							Data:         "ingress-data",
							SourceHandle: "egress",
							TargetHandle: "ingress",
						})
				}
			}
		}
	}

	return S.RespondJSON(w, http.StatusOK, g)
}

func NamespaceNetworkPolicyListHandler(app *S.App, c *S.Client, w http.ResponseWriter, r *http.Request) error {
	url := S.GetRequestParams(r, "/api/v1/namespace/netpol/")
	log.Printf("NamespaceNetworkPolicyListHandler url: %v", url)

	g := S.Graph{Nodes: []S.Node{}, Edges: []S.Edge{}}

	ns, err := c.Clientset.CoreV1().Namespaces().Get(context.TODO(), url.Scope, metav1.GetOptions{})
	if err != nil {
		return S.RespondError(err)
	}
	node := g.AddNode("ns", string(ns.ObjectMeta.UID), ns.ObjectMeta.Name, S.NodeOptions{Type: "namespace", Group: true})

	netpolList, err := c.Clientset.NetworkingV1().NetworkPolicies(url.Scope).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return S.RespondError(err)
	}
	for _, netpol := range netpolList.Items {
		g.AddNode("netpol", string(netpol.ObjectMeta.UID), netpol.ObjectMeta.Name, S.NodeOptions{Namespace: url.Scope, Type: "netpol", ParentNode: node, Extent: "parent"})
	}

	return S.RespondJSON(w, http.StatusOK, g)
}
