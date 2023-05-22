// Copyright (c) Autovia GmbH
// SPDX-License-Identifier: Apache-2.0

import CustomNodeEdge from './CustomNodeEdge';

export const nodeTypes = {
    pod: CustomNodeEdge,
    podedge: CustomNodeEdge,
    service: CustomNodeEdge,
    deploy: CustomNodeEdge,
    rs: CustomNodeEdge,
    sts: CustomNodeEdge,
    vol: CustomNodeEdge,
    secret: CustomNodeEdge,
    cm: CustomNodeEdge,
    sts: CustomNodeEdge,
    sa: CustomNodeEdge,
    pv: CustomNodeEdge,
    pvc: CustomNodeEdge,
    ds: CustomNodeEdge,
    rb: CustomNodeEdge,
    role: CustomNodeEdge,
    netpol: CustomNodeEdge,
    job: CustomNodeEdge,
    cronjob: CustomNodeEdge,
    svc: CustomNodeEdge,
    node: CustomNodeEdge,
    sc: CustomNodeEdge,
    cr: CustomNodeEdge,
    crb: CustomNodeEdge,
    crd: CustomNodeEdge,
    ic: CustomNodeEdge,
    ev: CustomNodeEdge
  };