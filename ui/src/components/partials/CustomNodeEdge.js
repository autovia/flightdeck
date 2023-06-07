// Copyright (c) Autovia GmbH
// SPDX-License-Identifier: Apache-2.0

import React, { memo } from 'react';
import { Handle, Position } from 'reactflow';
import Dropdown from './Dropdown';

function CustomNodeEdge({ data }) {
  return (
    <div className="px-2 py-2 shadow-md rounded-md bg-white border-2 border-stone-400">
      <div className="flex">
        <div className="flex-none w-12 h-12 flex justify-center items-center">
          <img src={'/static/k8s/' + data.kind + '-256.png'} />
        </div>
        <div className="ml-2">
          <div className="text-s [word-break:break-word]">{data.label}</div>
        </div>
        
      </div>

      <Handle type="target" position={Position.Top} className="w-16 !bg-sky-400" />
      <Handle type="source" position={Position.Bottom} className="w-16 !bg-sky-400" />
      <Handle position={Position.Right} className="mr-10"><Dropdown data={data}/></Handle>
    </div>
  );
}

export default memo(CustomNodeEdge);
