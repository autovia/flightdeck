// Copyright (c) Autovia GmbH
// SPDX-License-Identifier: Apache-2.0

package k8s

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	S "github.com/autovia/flightdeck/api/structs"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/rest"
)

func PodHandler(app *S.App, c *S.Client, w http.ResponseWriter, r *http.Request) error {
	url := S.GetRequestParams(r, "/api/v1/pod/")
	log.Printf("PodHandler url: %v", url)

	pod, err := c.Clientset.CoreV1().Pods(url.Namespace).Get(context.TODO(), url.Resource, metav1.GetOptions{})
	if err != nil {
		return S.RespondError(err)
	}
	pod.ObjectMeta.ManagedFields = nil

	return S.RespondYAML(w, http.StatusOK, pod)
}

func PodVolumeHandler(app *S.App, c *S.Client, w http.ResponseWriter, r *http.Request) error {
	url := S.GetRequestParams(r, "/api/v1/vol/")
	log.Printf("PodVolumeHandler url: %v", url)

	pod, err := c.Clientset.CoreV1().Pods(url.Namespace).Get(context.TODO(), url.Resource, metav1.GetOptions{})
	if err != nil {
		return S.RespondError(err)
	}
	for _, volume := range pod.Spec.Volumes {
		if volume.Name == url.Subresource {
			S.RespondYAML(w, http.StatusOK, volume)
			break
		}
	}
	return nil
}

func PodFilesystemHandler(app *S.App, c *S.Client, w http.ResponseWriter, r *http.Request) error {
	url := S.GetRequestParams(r, "/api/v1/fs/")
	log.Printf("PodFilesystemHandler url: %v", url)

	container, err := c.GetContainer(url)
	if err != nil {
		return S.RespondError(err)
	}

	path := "/"
	if url.Path != "" {
		path = path + url.Path
	}

	// debian, ubuntu
	// cmd := []string{"find", path, "-maxdepth", "1", "-printf", "{\"size\": %s, \"name\":\"%p\", \"modified\":\"%TY-%Tm-%Td %TH:%TM:%.2TS\", \"user\":\"%u\", \"group\":\"%g\", \"permission\": \"%m\", \"type\":\"%y\"},"}
	// busybox, alpine
	cmd := []string{"find", path, "-maxdepth", "1", "-exec", "stat", "-c", "{\"name\": \"%n\", \"size\":%s, \"modified\":\"%.19y\", \"user\":\"%U\", \"group\":\"%G\", \"permission\": \"%A\", \"type\":\"%.1F\"},", "{}", "+"}

	buf, err := c.ExecCommand(url, container, cmd)
	if err != nil {
		return S.RespondError(err)

	}

	return S.RespondText(w, http.StatusOK, fmt.Sprintf("[%s]", strings.TrimSuffix(strings.TrimSpace(buf.String()), ",")))
}

func PodFileOpenHandler(app *S.App, c *S.Client, w http.ResponseWriter, r *http.Request) error {
	url := S.GetRequestParams(r, "/api/v1/file/")
	log.Printf("PodFileOpenHandler url: %v", url)

	container, err := c.GetContainer(url)
	if err != nil {
		return S.RespondError(err)
	}

	path := "/"
	if url.Path != "" {
		path = path + url.Path
	}

	cmd := []string{"cat", path}

	buf, err := c.ExecCommand(url, container, cmd)
	if err != nil {
		return S.RespondError(err)
	}

	return S.RespondText(w, http.StatusOK, buf.String())
}

func PodLogsHandler(app *S.App, c *S.Client, w http.ResponseWriter, r *http.Request) error {
	url := S.GetRequestParams(r, "/api/v1/logs/")
	log.Printf("PodLogsHandler url: %v", url)

	var request *rest.Request
	// container name
	if url.Subresource != "" {
		// open container on request
		request = c.Clientset.CoreV1().Pods(url.Namespace).GetLogs(url.Resource, &corev1.PodLogOptions{TailLines: app.PodLogTailLines, Container: url.Subresource})
	} else {
		pod, err := c.Clientset.CoreV1().Pods(url.Namespace).Get(context.TODO(), url.Resource, metav1.GetOptions{})
		if err != nil {
			return S.RespondError(err)
		}

		if len(pod.Spec.Containers) > 0 {
			request = c.Clientset.CoreV1().Pods(url.Namespace).GetLogs(url.Resource, &corev1.PodLogOptions{TailLines: app.PodLogTailLines, Container: pod.Spec.Containers[0].Name})
		} else {
			request = c.Clientset.CoreV1().Pods(url.Namespace).GetLogs(url.Resource, &corev1.PodLogOptions{TailLines: app.PodLogTailLines})
		}
	}

	podLogs, err := request.Stream(context.TODO())
	if err != nil {
		return S.RespondError(err)
	}
	defer podLogs.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, podLogs)
	if err != nil {
		return S.RespondError(err)
	}
	return S.RespondText(w, http.StatusOK, buf.String())
}

func PodGraphHandler(app *S.App, c *S.Client, w http.ResponseWriter, r *http.Request) error {
	// namespace/podname
	url := strings.Split(strings.Replace(r.URL.Path, "/api/v1/graph/pod/", "", -1), "/")
	log.Printf("PodHandler url: %v", url)
	namespace := url[0]
	podname := url[1]

	// S.Graph
	g := S.Graph{Nodes: []S.Node{}, Edges: []S.Edge{}}

	// pod
	pod, err := c.Clientset.CoreV1().Pods(namespace).Get(context.TODO(), podname, metav1.GetOptions{})
	if err != nil {
		return S.RespondError(err)
	}

	var containers []string
	for _, c := range pod.Spec.Containers {
		containers = append(containers, c.Name)
	}
	podnode := g.AddNode("pod", string(pod.ObjectMeta.UID), pod.ObjectMeta.Name, S.NodeOptions{Containers: containers, Namespace: namespace, Type: "podedge"})

	// services
	svcList, err := c.Clientset.CoreV1().Services(pod.Namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return S.RespondError(err)
	}
	for _, service := range svcList.Items {
		if service.Spec.Selector == nil {
			continue
		}
		selector := labels.Set(service.Spec.Selector).AsSelectorPreValidated()
		if selector.Matches(labels.Set(pod.Labels)) {
			svcnode := g.AddNode("svc", string(service.ObjectMeta.UID), service.ObjectMeta.Name, S.NodeOptions{Namespace: namespace, Type: "service"})
			g.AddEdge(podnode, svcnode)
		}
	}

	// network policy
	netpolList, err := c.Clientset.NetworkingV1().NetworkPolicies(pod.Namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return S.RespondError(err)
	}
	for _, netpol := range netpolList.Items {
		selector := labels.Set(netpol.Spec.PodSelector.MatchLabels).AsSelectorPreValidated()
		if selector.Matches(labels.Set(pod.Labels)) {
			netpolnode := g.AddNode("netpol", netpol.ObjectMeta.Name, netpol.ObjectMeta.Name, S.NodeOptions{Namespace: namespace, Type: "netpol"})
			g.AddEdge(podnode, netpolnode)
		}
	}

	// replicaset, deployment
	for _, podOwnerRefs := range pod.ObjectMeta.OwnerReferences {
		log.Printf("podOwnerRefs.Kind: %v", podOwnerRefs.Kind)
		if podOwnerRefs.Kind == "ReplicaSet" {
			replicaset, err := c.Clientset.AppsV1().ReplicaSets(pod.Namespace).Get(context.Background(), podOwnerRefs.Name, metav1.GetOptions{})
			if err != nil {
				return S.RespondError(err)
			}
			replicasetnode := g.AddNode("rs", string(replicaset.ObjectMeta.UID), replicaset.ObjectMeta.Name, S.NodeOptions{Namespace: namespace, Type: "rs"})
			g.AddEdge(replicasetnode, podnode)

			for _, replOwnerRefs := range replicaset.ObjectMeta.OwnerReferences {
				if replOwnerRefs.Kind == "Deployment" {
					replDeployment, err := c.Clientset.AppsV1().Deployments(pod.Namespace).Get(context.Background(), replOwnerRefs.Name, metav1.GetOptions{})
					if err != nil {
						return S.RespondError(err)
					}
					repldeploynode := g.AddNode("deploy", string(replDeployment.ObjectMeta.UID), replDeployment.ObjectMeta.Name, S.NodeOptions{Namespace: namespace, Type: "deploy"})
					g.AddEdge(repldeploynode, replicasetnode)
				}
			}
		}
		if podOwnerRefs.Kind == "Deployment" {
			deployment, err := c.Clientset.AppsV1().Deployments(pod.Namespace).Get(context.Background(), podOwnerRefs.Name, metav1.GetOptions{})
			if err != nil {
				return S.RespondError(err)
			}
			deploynode := g.AddNode("deploy", string(deployment.ObjectMeta.UID), deployment.ObjectMeta.Name, S.NodeOptions{Namespace: namespace, Type: "deploy"})
			g.AddEdge(deploynode, podnode)
		}
		if podOwnerRefs.Kind == "StatefulSet" {
			statefulset, err := c.Clientset.AppsV1().StatefulSets(pod.Namespace).Get(context.Background(), podOwnerRefs.Name, metav1.GetOptions{})
			if err != nil {
				return S.RespondError(err)
			}
			stsnode := g.AddNode("sts", string(statefulset.ObjectMeta.UID), statefulset.ObjectMeta.Name, S.NodeOptions{Namespace: namespace, Type: "sts"})
			g.AddEdge(stsnode, podnode)
		}
		if podOwnerRefs.Kind == "DaemonSet" {
			daemonset, err := c.Clientset.AppsV1().DaemonSets(pod.Namespace).Get(context.Background(), podOwnerRefs.Name, metav1.GetOptions{})
			if err != nil {
				return S.RespondError(err)
			}
			dsnode := g.AddNode("ds", string(daemonset.ObjectMeta.UID), daemonset.ObjectMeta.Name, S.NodeOptions{Namespace: namespace, Type: "ds"})
			g.AddEdge(dsnode, podnode)
		}
		if podOwnerRefs.Kind == "Job" {
			job, err := c.Clientset.BatchV1().Jobs(pod.Namespace).Get(context.Background(), podOwnerRefs.Name, metav1.GetOptions{})
			if err != nil {
				return S.RespondError(err)
			}
			jobnode := g.AddNode("job", string(job.ObjectMeta.UID), job.ObjectMeta.Name, S.NodeOptions{Namespace: namespace, Type: "job"})
			g.AddEdge(jobnode, podnode)

			for _, jobOwnerRefs := range job.ObjectMeta.OwnerReferences {
				if jobOwnerRefs.Kind == "CronJob" {
					cronjob, err := c.Clientset.BatchV1().CronJobs(pod.Namespace).Get(context.Background(), jobOwnerRefs.Name, metav1.GetOptions{})
					if err != nil {
						return S.RespondError(err)
					}
					cjnode := g.AddNode("cronjob", string(cronjob.ObjectMeta.UID), cronjob.ObjectMeta.Name, S.NodeOptions{Namespace: namespace, Type: "cronjob"})
					g.AddEdge(cjnode, jobnode)
				}
			}
		}
		log.Printf("podOwnerRefs %v \n", podOwnerRefs.Kind)
	}

	// volumes
	for _, volume := range pod.Spec.Volumes {
		switch {
		case volume.Secret != nil:
			if !g.Includes(volume.Secret.SecretName) {
				secretnode := g.AddNode("secret", volume.Secret.SecretName, volume.Secret.SecretName, S.NodeOptions{Namespace: namespace, Type: "secret"})
				g.AddEdge(secretnode, podnode)
			}
		case volume.ConfigMap != nil:
			if !g.Includes(volume.ConfigMap.Name) {
				cmnode := g.AddNode("cm", volume.ConfigMap.Name, volume.ConfigMap.Name, S.NodeOptions{Namespace: namespace, Type: "cm"})
				g.AddEdge(cmnode, podnode)
			}
		case volume.PersistentVolumeClaim != nil:
			if !g.Includes(volume.PersistentVolumeClaim.ClaimName) {
				pvc, err := c.Clientset.CoreV1().PersistentVolumeClaims(pod.Namespace).Get(context.Background(), volume.PersistentVolumeClaim.ClaimName, metav1.GetOptions{})
				if err != nil {
					return S.RespondError(err)
				}
				volnode := g.AddNode("pv", pvc.Spec.VolumeName, pvc.Spec.VolumeName, S.NodeOptions{Namespace: namespace, Type: "pv"})
				g.AddEdge(volnode, podnode)
				pvcnode := g.AddNode("pvc", volume.PersistentVolumeClaim.ClaimName, volume.PersistentVolumeClaim.ClaimName, S.NodeOptions{Namespace: namespace, Type: "pvc"})
				g.AddEdge(pvcnode, volnode)
				if pvc.Spec.StorageClassName != nil {
					scnode := g.AddNode("sc", *pvc.Spec.StorageClassName, *pvc.Spec.StorageClassName, S.NodeOptions{Type: "sc"})
					g.AddEdge(scnode, pvcnode)
				}
			}
		case volume.Projected != nil:
			for _, source := range volume.Projected.Sources {
				switch {
				// case source.DownwardAPI != nil:
				// case source.ServiceAccountToken != nil:
				case source.Secret != nil:
					if !g.Includes(source.Secret.Name) {
						volsecretnode := g.AddNode("secret", source.Secret.Name, source.Secret.Name, S.NodeOptions{Namespace: namespace, Type: "secret"})
						g.AddEdge(volsecretnode, podnode)
					}
				case source.ConfigMap != nil:
					if !g.Includes(source.ConfigMap.Name) {
						volcmnode := g.AddNode("cm", source.ConfigMap.Name, source.ConfigMap.Name, S.NodeOptions{Namespace: namespace, Type: "cm"})
						g.AddEdge(volcmnode, podnode)
					}
				}
			}
		default:
			if !g.Includes(volume.Name) {
				volnode := g.AddNode("vol", volume.Name, volume.Name, S.NodeOptions{Namespace: namespace, Type: "vol"})
				g.AddEdge(volnode, podnode)
			}
		}
	}

	// serviceaccount
	if pod.Spec.ServiceAccountName != "" {
		sanode := g.AddNode("sa", pod.Spec.ServiceAccountName, pod.Spec.ServiceAccountName, S.NodeOptions{Namespace: namespace, Type: "sa"})
		g.AddEdge(sanode, podnode)

		rbList, err := c.Clientset.RbacV1().RoleBindings(pod.Namespace).List(context.Background(), metav1.ListOptions{})
		if err != nil {
			return S.RespondError(err)
		}
		for _, rb := range rbList.Items {
			for _, subject := range rb.Subjects {
				if subject.Kind == "ServiceAccount" && subject.Name == pod.Spec.ServiceAccountName {
					rbnode := g.AddNode("rb", rb.ObjectMeta.Name, rb.ObjectMeta.Name, S.NodeOptions{Namespace: namespace, Type: "rb"})
					g.AddEdge(rbnode, sanode)

					rolenode := g.AddNode("role", rb.RoleRef.Name, rb.RoleRef.Name, S.NodeOptions{Namespace: namespace, Type: "role"})
					g.AddEdge(rolenode, rbnode)
					break
				}
			}
		}

		crbList, err := c.Clientset.RbacV1().ClusterRoleBindings().List(context.Background(), metav1.ListOptions{})
		if err != nil {
			return S.RespondError(err)
		}
		for _, crb := range crbList.Items {
			for _, subject := range crb.Subjects {
				if subject.Kind == "ServiceAccount" && subject.Name == pod.Spec.ServiceAccountName {
					crbnode := g.AddNode("crb", crb.ObjectMeta.Name, crb.ObjectMeta.Name, S.NodeOptions{Namespace: namespace, Type: "crb"})
					g.AddEdge(crbnode, sanode)

					rolenode := g.AddNode("c-role", crb.RoleRef.Name, crb.RoleRef.Name, S.NodeOptions{Namespace: namespace, Type: "cr"})
					g.AddEdge(rolenode, crbnode)
					break
				}
			}
		}
	}

	return S.RespondJSON(w, http.StatusOK, g)
}

func NamespacePodListHandler(app *S.App, c *S.Client, w http.ResponseWriter, r *http.Request) error {
	url := S.GetRequestParams(r, "/api/v1/namespace/pod/")
	log.Printf("NamespacePodListHandler url: %v", url)

	g := S.Graph{Nodes: []S.Node{}, Edges: []S.Edge{}}

	ns, err := c.Clientset.CoreV1().Namespaces().Get(context.TODO(), url.Namespace, metav1.GetOptions{})
	if err != nil {
		return S.RespondError(err)
	}
	node := g.AddNode("ns", string(ns.ObjectMeta.UID), ns.ObjectMeta.Name, S.NodeOptions{Type: "namespace", Group: true})

	podList, err := c.Clientset.CoreV1().Pods(url.Namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return S.RespondError(err)
	}
	for _, pod := range podList.Items {
		g.AddNode("pod", string(pod.ObjectMeta.UID), pod.ObjectMeta.Name, S.NodeOptions{Namespace: url.Namespace, Type: "pod", ParentNode: node, Extent: "parent"})
	}

	return S.RespondJSON(w, http.StatusOK, g)
}
