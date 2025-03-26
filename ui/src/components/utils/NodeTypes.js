// Copyright (c) Autovia GmbH
// SPDX-License-Identifier: Apache-2.0

import CustomNode from "../partials/CustomNode";
import NetworkPolicyNode from "../partials/NetworkPolicyNode";

export const nodeTypes = {
  pod: CustomNode,
  podedge: CustomNode,
  podnetpol: NetworkPolicyNode,
  service: CustomNode,
  deploy: CustomNode,
  rs: CustomNode,
  vol: CustomNode,
  secret: CustomNode,
  cm: CustomNode,
  sts: CustomNode,
  sa: CustomNode,
  pv: CustomNode,
  pvc: CustomNode,
  ds: CustomNode,
  rb: CustomNode,
  role: CustomNode,
  netpol: CustomNode,
  job: CustomNode,
  cronjob: CustomNode,
  svc: CustomNode,
  node: CustomNode,
  sc: CustomNode,
  cr: CustomNode,
  crb: CustomNode,
  crd: CustomNode,
  ic: CustomNode,
  ev: CustomNode,
};
