// Copyright (c) Autovia GmbH
// SPDX-License-Identifier: Apache-2.0

import React, { memo } from "react";
import { Handle, Position } from "@xyflow/react";
import Dropdown from "./Dropdown";
import {
  ArrowRightCircleIcon,
  ArrowLeftCircleIcon,
  XCircleIcon,
} from "@heroicons/react/24/outline";

function CustomNode({ data }) {
  return (
    <div className="px-2 py-2 shadow-md rounded-md bg-white border-2 border-stone-400 text-center">
      {data.position === "left" ? (
        <>
          <div className="flex">
            <div className="flex-none w-12 h-12 flex">
              <img src={"/static/k8s/pod-256.png"} alt="" />
            </div>
            <div className="ml-2">
              <div className="text-s [word-break:break-word]">{data.label}</div>
            </div>
          </div>
          <Handle
            type="target"
            position={Position.Left}
            className="-mt-4 w-6 h-6 bg-white"
            id="ingress"
          >
            <ArrowRightCircleIcon className="h-6 w-6 !bg-lime-200 rounded-xl" />
          </Handle>
          <Handle
            type="source"
            position={Position.Left}
            className="mt-4 w-6  h-6 bg-white"
            id="egress"
          >
            <ArrowLeftCircleIcon className="h-6 w-6 !bg-lime-200 rounded-xl" />
          </Handle>
        </>
      ) : (
        <>
          <div className="flex">
            <div className="flex-none w-12 h-12 flex justify-center items-center">
              <img src={"/static/k8s/pod-256.png"} alt="" />
            </div>
            <div className="ml-2">
              <div className="text-s [word-break:break-word]">{data.label}</div>
            </div>
          </div>
          <Handle
            type="source"
            position={Position.Right}
            className="-mt-4 w-6 h-6 bg-white"
            id="egress"
          >
            <ArrowRightCircleIcon className="h-6 w-6 !bg-lime-200 rounded-xl" />
          </Handle>
          <Handle
            type="target"
            position={Position.Right}
            className="mt-4 w-6 h-6 bg-white"
            id="ingress"
          >
            <ArrowLeftCircleIcon className="h-6 w-6 !bg-lime-200 rounded-xl" />
          </Handle>
        </>
      )}
    </div>
  );
}

export default memo(CustomNode);
