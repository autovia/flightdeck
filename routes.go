// Copyright (c) Autovia GmbH
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"net/http"

	"github.com/autovia/flightdeck/api/handlers/auth"
	"github.com/autovia/flightdeck/api/handlers/dynamic"
	"github.com/autovia/flightdeck/api/handlers/k8s"
	S "github.com/autovia/flightdeck/api/structs"
)

func InitRoutes(app *S.App) {
	// Public
	app.Router.Handle("/api/v1/auth/login", S.Public{app, auth.LoginTokenHandler})

	// App
	app.Router.Handle("/api/v1/auth/reset", S.Authorization{app, auth.ResetTokenHandler})
	app.Router.Handle("/api/v1/auth/status", S.Authorization{app, auth.StatusTokenHandler})

	// Kubernetes

	// type CoreV1Interface interface
	// 	RESTClient() rest.Interface
	// 	ComponentStatusesGetter
	// 	ConfigMapsGetter
	app.Router.Handle("/api/v1/cm/", S.KubeClient{app, k8s.ConfigMapHandler})
	app.Router.Handle("/api/v1/graph/cm/", S.KubeClient{app, k8s.ConfigMapPodListHandler})
	app.Router.Handle("/api/v1/namespace/cm/", S.KubeClient{app, k8s.NamespaceConfigMapListHandler})

	// 	EndpointsGetter
	// 	EventsGetter
	app.Router.Handle("/api/v1/ev/", S.KubeClient{app, k8s.EventHandler})
	app.Router.Handle("/api/v1/namespace/ev/", S.KubeClient{app, k8s.NamespaceEventListHandler})

	// 	LimitRangesGetter
	// 	NamespacesGetter
	app.Router.Handle("/api/v1/namespaces", S.KubeClient{app, k8s.NamespaceListHandler})

	// 	NodesGetter
	app.Router.Handle("/api/v1/node/", S.KubeClient{app, k8s.NodeHandler})
	app.Router.Handle("/api/v1/no", S.KubeClient{app, k8s.NodeListHandler})

	// 	PersistentVolumesGetter
	app.Router.Handle("/api/v1/pv/", S.KubeClient{app, k8s.PersistentVolumeHandler})
	app.Router.Handle("/api/v1/pv", S.KubeClient{app, k8s.PersistentVolumeListHandler})

	// 	PersistentVolumeClaimsGetter
	app.Router.Handle("/api/v1/pvc/", S.KubeClient{app, k8s.PersistentVolumeClaim})
	app.Router.Handle("/api/v1/graph/pvc/", S.KubeClient{app, k8s.PersistentVolumeClaimPodListHandler})
	app.Router.Handle("/api/v1/namespace/pvc/", S.KubeClient{app, k8s.NamespacePersistentVolumeClaimListHandler})

	// 	PodsGetter
	app.Router.Handle("/api/v1/pod/", S.KubeClient{app, k8s.PodHandler})
	app.Router.Handle("/api/v1/graph/pod/", S.KubeClient{app, k8s.PodGraphHandler})
	app.Router.Handle("/api/v1/pod", S.KubeClient{app, k8s.NamespacePodListHandler})

	// 	PodTemplatesGetter
	// 	ReplicationControllersGetter
	// 	ResourceQuotasGetter
	// 	SecretsGetter
	app.Router.Handle("/api/v1/secret/", S.KubeClient{app, k8s.SecretHandler})
	app.Router.Handle("/api/v1/graph/secret/", S.KubeClient{app, k8s.SecretPodListHandler})
	app.Router.Handle("/api/v1/namespace/secret/", S.KubeClient{app, k8s.NamespaceSecretListHandler})

	// 	ServicesGetter
	app.Router.Handle("/api/v1/svc/", S.KubeClient{app, k8s.ServiceHandler})
	app.Router.Handle("/api/v1/graph/svc/", S.KubeClient{app, k8s.ServicePodListHandler})
	app.Router.Handle("/api/v1/namespace/svc/", S.KubeClient{app, k8s.NamespaceServiceListHandler})

	// 	ServiceAccountsGetter
	app.Router.Handle("/api/v1/sa/", S.KubeClient{app, k8s.ServiceAccountHandler})
	app.Router.Handle("/api/v1/graph/sa/", S.KubeClient{app, k8s.ServiceAccountPodListHandler})
	app.Router.Handle("/api/v1/namespace/sa/", S.KubeClient{app, k8s.NamespaceServiceAccountListHandler})

	// type AppsV1Interface interface
	// 	RESTClient() rest.Interface
	// 	ControllerRevisionsGetter
	// 	DaemonSetsGetter
	app.Router.Handle("/api/v1/ds/", S.KubeClient{app, k8s.DaemonSetHandler})
	app.Router.Handle("/api/v1/graph/ds/", S.KubeClient{app, k8s.DaemonSetPodListHandler})
	app.Router.Handle("/api/v1/namespace/ds/", S.KubeClient{app, k8s.NamespaceDaemonSetListHandler})

	// 	DeploymentsGetter
	app.Router.Handle("/api/v1/deploy/", S.KubeClient{app, k8s.DeploymentHandler})
	app.Router.Handle("/api/v1/graph/deploy/", S.KubeClient{app, k8s.DeploymentPodListHandler})
	app.Router.Handle("/api/v1/namespace/deploy/", S.KubeClient{app, k8s.NamespaceDeploymentListHandler})

	// 	ReplicaSetsGetter
	app.Router.Handle("/api/v1/rs/", S.KubeClient{app, k8s.ReplicaSetHandler})
	app.Router.Handle("/api/v1/graph/rs/", S.KubeClient{app, k8s.ReplicaSetPodListHandler})
	app.Router.Handle("/api/v1/namespace/rs/", S.KubeClient{app, k8s.NamespaceReplicaSetListHandler})

	// 	StatefulSetsGetter
	app.Router.Handle("/api/v1/sts/", S.KubeClient{app, k8s.StatefulSetHandler})
	app.Router.Handle("/api/v1/graph/sts/", S.KubeClient{app, k8s.StatefulSetPodListHandler})
	app.Router.Handle("/api/v1/namespace/sts/", S.KubeClient{app, k8s.NamespaceStatefulSetListHandler})

	// type RbacV1Interface interface
	// 	RESTClient() rest.Interface
	// 	ClusterRolesGetter
	app.Router.Handle("/api/v1/c-role/", S.KubeClient{app, k8s.ClusterRoleHandler})
	app.Router.Handle("/api/v1/c-role", S.KubeClient{app, k8s.ClusterRoleListHandler})

	// 	ClusterRoleBindingsGetter
	app.Router.Handle("/api/v1/crb/", S.KubeClient{app, k8s.ClusterRoleBindingHandler})
	app.Router.Handle("/api/v1/crb", S.KubeClient{app, k8s.ClusterRoleBindingListHandler})

	// 	RolesGetter
	app.Router.Handle("/api/v1/role/", S.KubeClient{app, k8s.RoleHandler})
	app.Router.Handle("/api/v1/namespace/role/", S.KubeClient{app, k8s.NamespaceRoleListHandler})

	// 	RoleBindingsGetter
	app.Router.Handle("/api/v1/rb/", S.KubeClient{app, k8s.RoleBindingHandler})
	app.Router.Handle("/api/v1/namespace/rb/", S.KubeClient{app, k8s.NamespaceRoleBindingListHandler})

	// type NetworkingV1Interface interface
	// 	RESTClient() rest.Interface
	// 	IngressesGetter
	app.Router.Handle("/api/v1/ing/", S.KubeClient{app, k8s.IngressHandler})
	app.Router.Handle("/api/v1/namespace/ing/", S.KubeClient{app, k8s.NamespaceIngressListHandler})

	// 	IngressClassesGetter
	app.Router.Handle("/api/v1/ic/", S.KubeClient{app, k8s.IngressClassHandler})
	app.Router.Handle("/api/v1/ic", S.KubeClient{app, k8s.IngressClassListHandler})

	// 	NetworkPoliciesGetter
	app.Router.Handle("/api/v1/netpol/", S.KubeClient{app, k8s.NetworkPolicyHandler})
	app.Router.Handle("/api/v1/namespace/netpol/", S.KubeClient{app, k8s.NamespaceNetworkPolicyListHandler})

	// type BatchV1Interface interface
	// 	RESTClient() rest.Interface
	// 	CronJobsGetter
	app.Router.Handle("/api/v1/cronjob/", S.KubeClient{app, k8s.CronJobHandler})
	app.Router.Handle("/api/v1/graph/cronjob/", S.KubeClient{app, k8s.CronJobPodListHandler})
	app.Router.Handle("/api/v1/namespace/cronjob/", S.KubeClient{app, k8s.NamespaceCronJobListHandler})

	// 	JobsGetter
	app.Router.Handle("/api/v1/job/", S.KubeClient{app, k8s.JobHandler})
	app.Router.Handle("/api/v1/graph/job/", S.KubeClient{app, k8s.JobPodListHandler})
	app.Router.Handle("/api/v1/namespace/job/", S.KubeClient{app, k8s.NamespaceJobListHandler})

	// type StorageV1Interface interface {
	// 	RESTClient() rest.Interface
	// 	CSIDriversGetter
	// 	CSINodesGetter
	// 	CSIStorageCapacitiesGetter
	// 	StorageClassesGetter
	app.Router.Handle("/api/v1/sc/", S.KubeClient{app, k8s.StorageClassHandler})
	app.Router.Handle("/api/v1/sc", S.KubeClient{app, k8s.StorageClassListHandler})

	// type ApiextensionsV1Interface interface {
	// 	RESTClient() rest.Interface
	// 	CustomResourceDefinitionsGetter
	app.Router.Handle("/api/v1/crd/", S.ApiClient{app, k8s.CustomResourceDefinitionHandler})
	app.Router.Handle("/api/v1/crd", S.ApiClient{app, k8s.CustomResourceDefinitionListHandler})

	// Pod scoped
	app.Router.Handle("/api/v1/pod/vol/", S.KubeClient{app, k8s.PodVolumeHandler})
	app.Router.Handle("/api/v1/pod/logs/", S.KubeClient{app, k8s.PodLogsHandler})
	app.Router.Handle("/api/v1/pod/fs/", S.KubeClient{app, k8s.PodFilesystemHandler})
	app.Router.Handle("/api/v1/pod/file/", S.KubeClient{app, k8s.PodFileOpenHandler})

	// Dynamic Client
	app.Router.Handle("/api/v1/list/", S.ApiClient{app, dynamic.ListClusterResourcesHandler})
	app.Router.Handle("/api/v1/resources/", S.ApiClient{app, dynamic.ListNamespaceResourcesHandler})

	// Main
	if *app.FileServer {
		app.Router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, *app.FileServerPath+"/index.html")
		})
		//app.Router.Handle("/k8s/", http.StripPrefix("/k8s", http.FileServer(http.Dir(*app.FileServerPath+"/k8s/"))))
		app.Router.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir(*app.FileServerPath+"/static/"))))
	}
}

//https://pkg.go.dev/k8s.io/client-go@v0.27.0/kubernetes

// type Interface interface {
// 	Discovery() discovery.DiscoveryInterface
// 	AdmissionregistrationV1() admissionregistrationv1.AdmissionregistrationV1Interface
// 	AdmissionregistrationV1alpha1() admissionregistrationv1alpha1.AdmissionregistrationV1alpha1Interface
// 	AdmissionregistrationV1beta1() admissionregistrationv1beta1.AdmissionregistrationV1beta1Interface
// 	InternalV1alpha1() internalv1alpha1.InternalV1alpha1Interface
// 	AppsV1() appsv1.AppsV1Interface
// 	AppsV1beta1() appsv1beta1.AppsV1beta1Interface
// 	AppsV1beta2() appsv1beta2.AppsV1beta2Interface
// 	AuthenticationV1() authenticationv1.AuthenticationV1Interface
// 	AuthenticationV1alpha1() authenticationv1alpha1.AuthenticationV1alpha1Interface
// 	AuthenticationV1beta1() authenticationv1beta1.AuthenticationV1beta1Interface
// 	AuthorizationV1() authorizationv1.AuthorizationV1Interface
// 	AuthorizationV1beta1() authorizationv1beta1.AuthorizationV1beta1Interface
// 	AutoscalingV1() autoscalingv1.AutoscalingV1Interface
// 	AutoscalingV2() autoscalingv2.AutoscalingV2Interface
// 	AutoscalingV2beta1() autoscalingv2beta1.AutoscalingV2beta1Interface
// 	AutoscalingV2beta2() autoscalingv2beta2.AutoscalingV2beta2Interface
// 	BatchV1() batchv1.BatchV1Interface
// 	BatchV1beta1() batchv1beta1.BatchV1beta1Interface
// 	CertificatesV1() certificatesv1.CertificatesV1Interface
// 	CertificatesV1beta1() certificatesv1beta1.CertificatesV1beta1Interface
// 	CertificatesV1alpha1() certificatesv1alpha1.CertificatesV1alpha1Interface
// 	CoordinationV1beta1() coordinationv1beta1.CoordinationV1beta1Interface
// 	CoordinationV1() coordinationv1.CoordinationV1Interface
// 	CoreV1() corev1.CoreV1Interface
// 	DiscoveryV1() discoveryv1.DiscoveryV1Interface
// 	DiscoveryV1beta1() discoveryv1beta1.DiscoveryV1beta1Interface
// 	EventsV1() eventsv1.EventsV1Interface
// 	EventsV1beta1() eventsv1beta1.EventsV1beta1Interface
// 	ExtensionsV1beta1() extensionsv1beta1.ExtensionsV1beta1Interface
// 	FlowcontrolV1alpha1() flowcontrolv1alpha1.FlowcontrolV1alpha1Interface
// 	FlowcontrolV1beta1() flowcontrolv1beta1.FlowcontrolV1beta1Interface
// 	FlowcontrolV1beta2() flowcontrolv1beta2.FlowcontrolV1beta2Interface
// 	FlowcontrolV1beta3() flowcontrolv1beta3.FlowcontrolV1beta3Interface
// 	NetworkingV1() networkingv1.NetworkingV1Interface
// 	NetworkingV1alpha1() networkingv1alpha1.NetworkingV1alpha1Interface
// 	NetworkingV1beta1() networkingv1beta1.NetworkingV1beta1Interface
// 	NodeV1() nodev1.NodeV1Interface
// 	NodeV1alpha1() nodev1alpha1.NodeV1alpha1Interface
// 	NodeV1beta1() nodev1beta1.NodeV1beta1Interface
// 	PolicyV1() policyv1.PolicyV1Interface
// 	PolicyV1beta1() policyv1beta1.PolicyV1beta1Interface
// 	RbacV1() rbacv1.RbacV1Interface
// 	RbacV1beta1() rbacv1beta1.RbacV1beta1Interface
// 	RbacV1alpha1() rbacv1alpha1.RbacV1alpha1Interface
// 	ResourceV1alpha2() resourcev1alpha2.ResourceV1alpha2Interface
// 	SchedulingV1alpha1() schedulingv1alpha1.SchedulingV1alpha1Interface
// 	SchedulingV1beta1() schedulingv1beta1.SchedulingV1beta1Interface
// 	SchedulingV1() schedulingv1.SchedulingV1Interface
// 	StorageV1beta1() storagev1beta1.StorageV1beta1Interface
// 	StorageV1() storagev1.StorageV1Interface
// 	StorageV1alpha1() storagev1alpha1.StorageV1alpha1Interface
// }
