// Copyright (c) Autovia GmbH
// SPDX-License-Identifier: Apache-2.0

package utils

import "k8s.io/apimachinery/pkg/runtime/schema"

var GVR = map[string]schema.GroupVersionResource{
	// cluster resource
	"c-role": {Group: "rbac.authorization.k8s.io", Version: "v1", Resource: "clusterroles"},
	"crb":    {Group: "rbac.authorization.k8s.io", Version: "v1", Resource: "clusterrolebindings"},
	"crd":    {Group: "apiextensions.k8s.io", Version: "v1", Resource: "customresourcedefinitions"},
	"ic":     {Group: "networking.k8s.io", Version: "v1", Resource: "ingressclasses"},
	"no":     {Version: "v1", Resource: "nodes"},
	"pv":     {Version: "v1", Resource: "persistentvolumes"},
	"sc":     {Group: "storage.k8s.io", Version: "v1", Resource: "storageclasses"},
	// namespace resource
	"cm":      {Version: "v1", Resource: "configmaps"},
	"cronjob": {Group: "batch", Version: "v1", Resource: "cronjobs"},
	"ds":      {Group: "apps", Version: "v1", Resource: "daemonsets"},
	"deploy":  {Group: "apps", Version: "v1", Resource: "deployments"},
	"ev":      {Version: "v1", Resource: "events"},
	"ing":     {Group: "networking.k8s.io", Version: "v1", Resource: "ingresses"},
	"job":     {Group: "batch", Version: "v1", Resource: "jobs"},
	"netpol":  {Group: "networking.k8s.io", Version: "v1", Resource: "networkpolicies"},
	"pvc":     {Version: "v1", Resource: "persistentvolumeclaims"},
	"pod":     {Version: "v1", Resource: "pods"},
	"rs":      {Group: "apps", Version: "v1", Resource: "replicasets"},
	"rc":      {Version: "v1", Resource: "replicationcontrollers"},
	"role":    {Group: "rbac.authorization.k8s.io", Version: "v1", Resource: "roles"},
	"rb":      {Group: "rbac.authorization.k8s.io", Version: "v1", Resource: "rolebindings"},
	"secret":  {Version: "v1", Resource: "secrets"},
	"svc":     {Version: "v1", Resource: "services"},
	"sa":      {Version: "v1", Resource: "serviceaccounts"},
	"sts":     {Group: "apps", Version: "v1", Resource: "statefulsets"},
}

// NAME                              SHORTNAMES         APIVERSION                             NAMESPACED   KIND                             VERBS
// bindings                                             v1                                     true         Binding                          [create]
// componentstatuses                 cs                 v1                                     false        ComponentStatus                  [get list]
// configmaps                        cm                 v1                                     true         ConfigMap                        [create delete deletecollection get list patch update watch]
// endpoints                         ep                 v1                                     true         Endpoints                        [create delete deletecollection get list patch update watch]
// events                            ev                 v1                                     true         Event                            [create delete deletecollection get list patch update watch]
// limitranges                       limits             v1                                     true         LimitRange                       [create delete deletecollection get list patch update watch]
// namespaces                        ns                 v1                                     false        Namespace                        [create delete get list patch update watch]
// nodes                             no                 v1                                     false        Node                             [create delete deletecollection get list patch update watch]
// persistentvolumeclaims            pvc                v1                                     true         PersistentVolumeClaim            [create delete deletecollection get list patch update watch]
// persistentvolumes                 pv                 v1                                     false        PersistentVolume                 [create delete deletecollection get list patch update watch]
// pods                              po                 v1                                     true         Pod                              [create delete deletecollection get list patch update watch]
// podtemplates                                         v1                                     true         PodTemplate                      [create delete deletecollection get list patch update watch]
// replicationcontrollers            rc                 v1                                     true         ReplicationController            [create delete deletecollection get list patch update watch]
// resourcequotas                    quota              v1                                     true         ResourceQuota                    [create delete deletecollection get list patch update watch]
// secrets                                              v1                                     true         Secret                           [create delete deletecollection get list patch update watch]
// serviceaccounts                   sa                 v1                                     true         ServiceAccount                   [create delete deletecollection get list patch update watch]
// services                          svc                v1                                     true         Service                          [create delete deletecollection get list patch update watch]
// mutatingwebhookconfigurations                        admissionregistration.k8s.io/v1        false        MutatingWebhookConfiguration     [create delete deletecollection get list patch update watch]
// validatingwebhookconfigurations                      admissionregistration.k8s.io/v1        false        ValidatingWebhookConfiguration   [create delete deletecollection get list patch update watch]
// customresourcedefinitions         crd,crds           apiextensions.k8s.io/v1                false        CustomResourceDefinition         [create delete deletecollection get list patch update watch]
// apiservices                                          apiregistration.k8s.io/v1              false        APIService                       [create delete deletecollection get list patch update watch]
// controllerrevisions                                  apps/v1                                true         ControllerRevision               [create delete deletecollection get list patch update watch]
// daemonsets                        ds                 apps/v1                                true         DaemonSet                        [create delete deletecollection get list patch update watch]
// deployments                       deploy             apps/v1                                true         Deployment                       [create delete deletecollection get list patch update watch]
// replicasets                       rs                 apps/v1                                true         ReplicaSet                       [create delete deletecollection get list patch update watch]
// statefulsets                      sts                apps/v1                                true         StatefulSet                      [create delete deletecollection get list patch update watch]
// applications                      app,apps           argoproj.io/v1alpha1                   true         Application                      [delete deletecollection get list patch create update watch]
// applicationsets                   appset,appsets     argoproj.io/v1alpha1                   true         ApplicationSet                   [delete deletecollection get list patch create update watch]
// appprojects                       appproj,appprojs   argoproj.io/v1alpha1                   true         AppProject                       [delete deletecollection get list patch create update watch]
// tokenreviews                                         authentication.k8s.io/v1               false        TokenReview                      [create]
// localsubjectaccessreviews                            authorization.k8s.io/v1                true         LocalSubjectAccessReview         [create]
// selfsubjectaccessreviews                             authorization.k8s.io/v1                false        SelfSubjectAccessReview          [create]
// selfsubjectrulesreviews                              authorization.k8s.io/v1                false        SelfSubjectRulesReview           [create]
// subjectaccessreviews                                 authorization.k8s.io/v1                false        SubjectAccessReview              [create]
// horizontalpodautoscalers          hpa                autoscaling/v2                         true         HorizontalPodAutoscaler          [create delete deletecollection get list patch update watch]
// cronjobs                          cj                 batch/v1                               true         CronJob                          [create delete deletecollection get list patch update watch]
// jobs                                                 batch/v1                               true         Job                              [create delete deletecollection get list patch update watch]
// certificatesigningrequests        csr                certificates.k8s.io/v1                 false        CertificateSigningRequest        [create delete deletecollection get list patch update watch]
// leases                                               coordination.k8s.io/v1                 true         Lease                            [create delete deletecollection get list patch update watch]
// bgpconfigurations                                    crd.projectcalico.org/v1               false        BGPConfiguration                 [delete deletecollection get list patch create update watch]
// bgppeers                                             crd.projectcalico.org/v1               false        BGPPeer                          [delete deletecollection get list patch create update watch]
// blockaffinities                                      crd.projectcalico.org/v1               false        BlockAffinity                    [delete deletecollection get list patch create update watch]
// caliconodestatuses                                   crd.projectcalico.org/v1               false        CalicoNodeStatus                 [delete deletecollection get list patch create update watch]
// clusterinformations                                  crd.projectcalico.org/v1               false        ClusterInformation               [delete deletecollection get list patch create update watch]
// felixconfigurations                                  crd.projectcalico.org/v1               false        FelixConfiguration               [delete deletecollection get list patch create update watch]
// globalnetworkpolicies                                crd.projectcalico.org/v1               false        GlobalNetworkPolicy              [delete deletecollection get list patch create update watch]
// globalnetworksets                                    crd.projectcalico.org/v1               false        GlobalNetworkSet                 [delete deletecollection get list patch create update watch]
// hostendpoints                                        crd.projectcalico.org/v1               false        HostEndpoint                     [delete deletecollection get list patch create update watch]
// ipamblocks                                           crd.projectcalico.org/v1               false        IPAMBlock                        [delete deletecollection get list patch create update watch]
// ipamconfigs                                          crd.projectcalico.org/v1               false        IPAMConfig                       [delete deletecollection get list patch create update watch]
// ipamhandles                                          crd.projectcalico.org/v1               false        IPAMHandle                       [delete deletecollection get list patch create update watch]
// ippools                                              crd.projectcalico.org/v1               false        IPPool                           [delete deletecollection get list patch create update watch]
// ipreservations                                       crd.projectcalico.org/v1               false        IPReservation                    [delete deletecollection get list patch create update watch]
// kubecontrollersconfigurations                        crd.projectcalico.org/v1               false        KubeControllersConfiguration     [delete deletecollection get list patch create update watch]
// networkpolicies                                      crd.projectcalico.org/v1               true         NetworkPolicy                    [delete deletecollection get list patch create update watch]
// networksets                                          crd.projectcalico.org/v1               true         NetworkSet                       [delete deletecollection get list patch create update watch]
// endpointslices                                       discovery.k8s.io/v1                    true         EndpointSlice                    [create delete deletecollection get list patch update watch]
// events                            ev                 events.k8s.io/v1                       true         Event                            [create delete deletecollection get list patch update watch]
// flowschemas                                          flowcontrol.apiserver.k8s.io/v1beta2   false        FlowSchema                       [create delete deletecollection get list patch update watch]
// prioritylevelconfigurations                          flowcontrol.apiserver.k8s.io/v1beta2   false        PriorityLevelConfiguration       [create delete deletecollection get list patch update watch]
// apprepositories                   apprepos           kubeapps.com/v1alpha1                  true         AppRepository                    [delete deletecollection get list patch create update watch]
// nodes                                                metrics.k8s.io/v1beta1                 false        NodeMetrics                      [get list]
// pods                                                 metrics.k8s.io/v1beta1                 true         PodMetrics                       [get list]
// ingressclasses                                       networking.k8s.io/v1                   false        IngressClass                     [create delete deletecollection get list patch update watch]
// ingresses                         ing                networking.k8s.io/v1                   true         Ingress                          [create delete deletecollection get list patch update watch]
// networkpolicies                   netpol             networking.k8s.io/v1                   true         NetworkPolicy                    [create delete deletecollection get list patch update watch]
// runtimeclasses                                       node.k8s.io/v1                         false        RuntimeClass                     [create delete deletecollection get list patch update watch]
// poddisruptionbudgets              pdb                policy/v1                              true         PodDisruptionBudget              [create delete deletecollection get list patch update watch]
// clusterrolebindings                                  rbac.authorization.k8s.io/v1           false        ClusterRoleBinding               [create delete deletecollection get list patch update watch]
// clusterroles                                         rbac.authorization.k8s.io/v1           false        ClusterRole                      [create delete deletecollection get list patch update watch]
// rolebindings                                         rbac.authorization.k8s.io/v1           true         RoleBinding                      [create delete deletecollection get list patch update watch]
// roles                                                rbac.authorization.k8s.io/v1           true         Role                             [create delete deletecollection get list patch update watch]
// priorityclasses                   pc                 scheduling.k8s.io/v1                   false        PriorityClass                    [create delete deletecollection get list patch update watch]
// csidrivers                                           storage.k8s.io/v1                      false        CSIDriver                        [create delete deletecollection get list patch update watch]
// csinodes                                             storage.k8s.io/v1                      false        CSINode                          [create delete deletecollection get list patch update watch]
// csistoragecapacities                                 storage.k8s.io/v1                      true         CSIStorageCapacity               [create delete deletecollection get list patch update watch]
// storageclasses                    sc                 storage.k8s.io/v1                      false        StorageClass                     [create delete deletecollection get list patch update watch]
// volumeattachments                                    storage.k8s.io/v1                      false        VolumeAttachment                 [create delete deletecollection get list patch update watch]
