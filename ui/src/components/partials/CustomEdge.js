// Copyright (c) Autovia GmbH
// SPDX-License-Identifier: Apache-2.0

import React, { memo } from "react";
import {
  BaseEdge,
  EdgeText,
  EdgeLabelRenderer,
  getStraightPath,
  useReactFlow,
  getSimpleBezierPath,
  Position,
} from "@xyflow/react";

function CustomEdge({ id, sourceX, sourceY, targetX, targetY }) {
  const { setEdges } = useReactFlow();
  // const [edgePath, labelX, labelY] = getStraightPath({
  //   sourceX,
  //   sourceY,
  //   targetX,
  //   targetY,
  // });
  const [path, labelX, labelY, offsetX, offsetY] = getSimpleBezierPath({
    sourceX: sourceX,
    sourceY: sourceY,
    sourcePosition: Position.Top,
    targetX: targetX,
    targetY: targetY,
    targetPosition: Position.Bottom,
  });

  return (
    <>
      <BaseEdge id={id} path={path} />
    </>
  );
}

export default memo(CustomEdge);
