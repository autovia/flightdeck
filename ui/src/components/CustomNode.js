// Copyright (c) Autovia GmbH
// SPDX-License-Identifier: Apache-2.0

import React, { memo } from 'react';

function ServiceNode({ data }) {
  return (
    <div className="px-2 py-2 shadow-md rounded-md bg-white border-2 border-stone-500">
      <div className="flex">
        <div className="flex-none w-12 h-12 flex justify-center items-center">
          <img src={'/static/k8s/' + data.kind + '-256.png'} />
        </div>
        <div className="ml-2">
          <div className="text-lg">{data.label}</div>
        </div>
      </div>
    </div>
  );
}

export default memo(ServiceNode);
