// Copyright (c) Autovia GmbH
// SPDX-License-Identifier: Apache-2.0

import dagre from 'dagre';

export const layout = {
  graph(nodes, edges, direction = 'TB') {
    const dagreGraph = new dagre.graphlib.Graph();
    dagreGraph.setDefaultEdgeLabel(() => ({}));
  
    const nodeWidth = 250;
    const nodeHeight = 100;
    const isHorizontal = direction === 'LR';
    dagreGraph.setGraph({ rankdir: direction });
  
    nodes.forEach((node) => {
      dagreGraph.setNode(node.id, { width: nodeWidth, height: nodeHeight });
    });
  
    edges.forEach((edge) => {
      dagreGraph.setEdge(edge.source, edge.target);
    });
    
    dagre.layout(dagreGraph);
  
    nodes.forEach((node) => {
      const nodeWithPosition = dagreGraph.node(node.id);
      node.targetPosition = isHorizontal ? 'left' : 'top';
      node.sourcePosition = isHorizontal ? 'right' : 'bottom';
  
      node.position = {
        x: nodeWithPosition.x - nodeWidth / 2,
        y: nodeWithPosition.y - nodeHeight / 2,
      };
      node.zIndex = 1;

      return node;
    });
  
    return {nodes, edges};
  },

  group(nodes, edges, columns) {
    if (nodes.length > 25) {
      columns = 10
    } else {
      columns = 5
    }
    var i = 0;
    var xoffset = 0;
    var yoffset = 0;
    var xmax = 0;
    var ymax = 0;
    nodes.forEach((node) => {
      node.position = {x: 50, y: 50};
      if (node.data.group) {
        i--;
      } else {
        if (i !== 0) {
          if (i % columns === 0) {
            // newline
            i = 0;
            xoffset = 0;
            yoffset += 150;
            node.position = {x: node.position.x + xoffset, y: node.position.y + yoffset};
          } else {
            xoffset += 300;
            node.position = {x: node.position.x + xoffset, y: node.position.y + yoffset};
          }
        }
      }
      xmax = xoffset > xmax ? xoffset : xmax; 
      ymax = yoffset > ymax ? yoffset : ymax;
      i++;
      //console.log(node.id, node.position.x, xoffset, node.position.y, yoffset);
      return node;
    });
    nodes[0].style.width = xmax + 375;
    nodes[0].style.height = ymax + 200;
    nodes[0].type = "default";
    return {nodes, edges};
  }
};
