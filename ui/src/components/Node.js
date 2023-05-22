// Copyright (c) Autovia GmbH
// SPDX-License-Identifier: Apache-2.0

import React, {Component} from 'react';
import ReactFlow, { ConnectionLineType, Panel, Controls, Background, MiniMap} from 'reactflow';
import {useParams} from 'react-router-dom';
import PodNode from './PodNode';
import CustomNode from './CustomNode';
import Nav from './Nav';
//import 'reactflow/dist/style.css';

//import './index.css';

const nodeTypes = {
  pod: CustomNode,
  service: CustomNode,
  deployment: CustomNode,
  replicaset: CustomNode,
  node: CustomNode
};


class Node extends Component {

  constructor(props) {
    super(props);
    this.state = {
      nodes: [],
      edges: [],
      namespaces: []
    }
  }

  componentDidMount() {
    console.log("this.props", this.props);
    
    fetch('/api/v1/namespaces')
    .then(res => res.json())
    .then(d => {
      console.log('/api/v1/namespaces', d);
      const namespaces = d.nodes.map((n) => n.data.label);
      if (namespaces.length > 0) {
        this.setState((state, props) => ({
          nodes: state.nodes,
          edges: state.edges,
          namespaces: namespaces
        }));
      }
    });

    fetch('/api/v1/node/' + this.props.params.node)
    .then(res => res.json())
    .then(d => {
      console.log('/api/v1/node/', d);
      // d.nodes = d.nodes.map(node => {
      //   const newLabel = (<div className="on-hover">
      //     {process(node.data.label)}
      //     {node.data.pathMappings && <div className="on-hover-child"><div>pathMappings:</div>{process(node.data.pathMappings)}</div>}
      //   </div>);
      //   return { ...node, data: { ...node.data, label: newLabel } };
      // });

      const {nodes: layoutedNodes, edges: layoutedEdges} = this.getLayoutedGroup(d.nodes, d.edges);
      this.setState((state, props) => ({
        nodes: layoutedNodes,
        edges: layoutedEdges,
        namespaces: state.namespaces
      }));
    });
  }

  onPodClick = (e, node) => {
    if (node.data.kind === "pod") {
        window.open("/namespace/" + node.data.namespace + "/pod/" + node.data.label, "_self");
    }
  };

  goHome() {
    window.open("/", "_self");
  }

  onNamespaceChange(e) {
      window.open("/namespace/" + e.target.value, "_self");
  }

  getLayoutedGroup(nodes, edges, columns = 5) {
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
  };

  componentDidUpdate(prevProps, prevState) {
    if(prevState !== this.state) {
      console.log('node componentDidUpdate: ', this.state);
    }
  }

  render() {
    return (
      <div style={{ width: '100vw', height: '100vh' }}>
        <ReactFlow
          nodes={this.state.nodes}
          edges={this.state.edges}
          //onNodesChange={onNodesChange}
          //onEdgesChange={onEdgesChange}
          //onConnect={onConnect}
          onNodeClick={this.onPodClick}
          proOptions={{ hideAttribution: true }}
          connectionLineType={ConnectionLineType.SmoothStep}
          nodeTypes={nodeTypes}
          fitView
          className="bg-sky-50"
        >
          <Panel position="top-left" className="w-full p-0 m-0">
            <Nav params={this.props.params} />
          </Panel>
          <Controls />
          <MiniMap />
          <Background variant="dots" gap={12} size={1} />
        </ReactFlow>
      </div>
    );
  }
}

function withParams(Component) {
  return (props) => <Component {...props} params={useParams()} />;
}

export default withParams(Node);
