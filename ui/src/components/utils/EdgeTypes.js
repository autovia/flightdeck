// Copyright (c) Autovia GmbH
// SPDX-License-Identifier: Apache-2.0

import CustomEdge from "../partials/CustomEdge";
import CustomEgressEdge from "../partials/CustomEgressEdge";
import CustomIngressEdge from "../partials/CustomIngressEdge";

export const edgeTypes = {
  default: CustomEdge,
  netpolegress: CustomEgressEdge,
  netpolingress: CustomIngressEdge,
};
