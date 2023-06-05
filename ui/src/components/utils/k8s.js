// Copyright (c) Autovia GmbH
// SPDX-License-Identifier: Apache-2.0

export const resources = [
    { id: "cm", name: 'Config Maps', type: 'config' },
    { id: "cronjob", name: 'Cron Jobs', type: 'workload' },
    { id: "ds", name: 'DaemonSet', type: 'workload' },
    { id: "deploy", name: 'Deployments', type: 'workload' },
    { id: "ev", name: 'Events', type: 'cluster' },
    { id: "ing", name: 'Ingresses', type: 'service' },
    { id: "job", name: 'Jobs', type: 'workload' },
    { id: "netpol", name: 'Network Policies', type: 'cluster' },
    { id: "pvc", name: 'Persistent Volume Claims', type: 'storage' },
    { id: "pod", name: 'Pods', type: 'workload'}, // default
    { id: "rs", name: 'Replica Sets', type: 'workload' },
    { id: "rc", name: 'Replication Controllers', type: 'workload' },
    { id: "role", name: 'Roles', type: 'cluster' },
    { id: "rb", name: 'Role Bindings', type: 'cluster' },
    { id: "secret", name: 'Secrets', type: 'config' },
    { id: "svc", name: 'Services', type: 'service' },
    { id: "sa", name: 'Service Accounts', type: 'cluster' },
    { id: "sts", name: 'Stateful Sets', type: 'workload' }
  ]
  
export const clusterResources = [
    { id: "c-role", name: 'Cluster Roles', type: 'cluster' },
    { id: "crb", name: 'Cluster Role Bindings', type: 'cluster' },
    { id: "crd", name: 'Custom Resource Definitions', type: 'cluster' },
    { id: "ic", name: 'Ingress Classes', type: 'service' },
    { id: "no", name: 'Nodes', type: 'cluster' }, // default
    { id: "pv", name: 'Persistent Volumes', type: 'storage' },
    { id: "sc", name: 'Storage Classes', type: 'storage' }
  ]