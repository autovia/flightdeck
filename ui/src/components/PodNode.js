// Copyright (c) Autovia GmbH
// SPDX-License-Identifier: Apache-2.0

import React, { memo } from 'react';
import { Handle, Position } from 'reactflow';

function PodNode({ data }) {
  return (
    <div className="px-4 py-2 shadow-md rounded-md bg-white border-2 border-stone-400">
        <div className="ml-2">
          <div className="text-gray-500">{data.kind}</div>
          <div className="text-lg truncate">{data.label}</div>
        </div>
    </div>
  );
}

export default memo(PodNode);
