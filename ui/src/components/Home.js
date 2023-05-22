// Copyright (c) Autovia GmbH
// SPDX-License-Identifier: Apache-2.0

import React, {Component} from 'react';
import ReactFlow, { addEdge, ConnectionLineType, useNodesState, useEdgesState, Panel, Controls, Background, MiniMap} from 'reactflow';
import dagre from 'dagre';
import CustomNode from './CustomNode';

const nodeTypes = {
  pod: CustomNode,
  service: CustomNode,
  deployment: CustomNode,
  replicaset: CustomNode,
  namespace: CustomNode,
};

class Home extends Component {

  constructor(props) {
    super(props);
    this.state = {
      nodes: [],
      edges: [],
      namespaces: [],
      namespace: ""
    }
  }

  componentDidMount() {
    const queryParameters = new URLSearchParams(window.location.search)
    const pod = queryParameters.get("pod")
    const namespace = queryParameters.get("namespace")

    fetch('/api/v1/namespaces')
    .then(res => res.json())
    .then(d => {
      console.log('/api/v1/namespaces', d);
      // d.nodes = d.nodes.map(node => {
      //   const newLabel = (<div><b>{node.data.kind}</b><br />{node.data.label}</div>);
      //   return { ...node, data: { ...node.data, label: newLabel } };
      // });
      const namespaces = d.nodes.map((n) => n.data.label);
      if (namespaces.length > 0) {
        this.setState((state, props) => ({
          nodes: state.nodes,
          edges: state.edges,
          namespaces: namespaces,
          selectedNamespace: namespaces[0]
        }));
      }

      if (!pod && !namespace) {
        const {nodes: layoutedNodes, edges: layoutedEdges} = this.getLayoutedGrid(d.nodes, d.edges);
        this.setState((state, props) => ({
          nodes: layoutedNodes,
          edges: layoutedEdges,
          namespaces: state.namespaces,
          selectedNamespace: state.selectedNamespace
        }));
      }
    });

    if (!pod && namespace) {
      fetch('/api/v1/namespace/' + namespace)
      .then(res => res.json())
      .then(d => {
        console.log('/api/v1/namespace/', d);
        // d.nodes = d.nodes.map(node => {
        //   const newLabel = (<div className="on-hover">
        //     {process(node.data.label)}
        //     {node.data.pathMappings && <div className="on-hover-child"><div>pathMappings:</div>{process(node.data.pathMappings)}</div>}
        //   </div>);
        //   return { ...node, data: { ...node.data, label: newLabel } };
        // });
        this.setNodes(d.nodes);
        this.setEdges(d.edges);
        this.setSelectedNamespace(namespace);
        const { nodes: layoutedNodes, edges: layoutedEdges } = this.getLayoutedGroup(
          d.nodes,
          d.edges
        );
      });
    }

    if (pod && namespace) {
      fetch('/api/v1/pod/' + namespace + '/' + pod)
      .then(res => res.json())
      .then(d => {
        console.log('/api/v1/pod/', d);
        // d.nodes = d.nodes.map(node => {
        //   const newLabel = (<div className="on-hover">
        //     {process(node.data.label)}
        //     {node.data.pathMappings && <div className="on-hover-child"><div>pathMappings:</div>{process(node.data.pathMappings)}</div>}
        //   </div>);
        //   return { ...node, data: { ...node.data, label: newLabel } };
        // });
        this.setNodes(d.nodes);
        this.setEdges(d.edges);
        this.setSelectedNamespace(namespace);
        const { nodes: layoutedNodes, edges: layoutedEdges } = this.getLayoutedGraph(
          d.nodes,
          d.edges
        );
      });
    }    
    //console.log("namespaces: ", namespaces);
    //console.log("selectedNamespace: ", selectedNamespace);  
  }

  onNodeClick(e, node) {
    const queryParameters = new URLSearchParams(window.location.search);
    const namespace = queryParameters.get("namespace");
    if (node) {
      console.log(node.data);
      if (node.data.kind === "ns") {
        window.open("/namespace/" + node.data.label, "_self");
      }
      if (node.data.kind === "pod") {
        window.open("/?namespace=" + namespace + "&" + node.data.kind + "=" + node.data.label, "_self");
      }
    }
  };

  getLayoutedGrid(nodes, edges, groups = [], columns = 3) {
    var i = 0;
    var xoffset = 0;
    var yoffset = 0;
    nodes.forEach((node) => {
      node.position = {x: 0, y: 0};
      if (i % columns == 0) {
        // newline
        i = 0;
        xoffset = 0;
        yoffset += 100;
        node.position = {x: node.position.x + xoffset, y: node.position.y + yoffset};
      } else {
        xoffset += 300;
        node.position = {x: node.position.x + xoffset, y: node.position.y + yoffset};
      }
      i++;
      //console.log(node.id, node.position.x, xoffset, node.position.y, yoffset);
      return node;
    });
    return {nodes, edges};
  };

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
            yoffset += 100;
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
    nodes[0].style.width = xmax + 350;
    nodes[0].style.height = ymax + 200;
    nodes[0].type = "default";
    return {nodes, edges};
  };

  getLayoutedGraph(nodes, edges, direction = 'TB') {
    const dagreGraph = new dagre.graphlib.Graph();
    dagreGraph.setDefaultEdgeLabel(() => ({}));
  
    const nodeWidth = 172;
    const nodeHeight = 36;
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
    
      return node;
    });
  
    return {nodes, edges};
  };
  
  componentDidUpdate(prevProps, prevState) {
    //console.log('componentDidUpdate: ', prevState.state);
    if(prevState !== this.state) {
      console.log('home componentDidUpdate: ', this.state);
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
          onNodeClick={this.onNodeClick}
          proOptions={{ hideAttribution: true }}
          connectionLineType={ConnectionLineType.SmoothStep}
          nodeTypes={nodeTypes}
          fitView
          className="bg-sky-50"
        >
          <Panel position="top-left" className="w-full p-0 m-0">
          <div className="bg-white bg-opacity-80 shadow sm:rounded-lg">
            <div className="px-4 py-5 sm:p-6">
              <div className="hidden sm:-my-px sm:ml-6 sm:flex sm:space-x-8">
                <a key="cluster" href="/nodes" className="border-transparent text-gray-500 hover:border-gray-300 hover:text-gray-700 inline-flex items-center border-b-2 px-1 pt-1 text-sm font-medium">
                  Cluster
                </a>
                <a key="namespaces" href="/" className="border-indigo-500 text-gray-900 inline-flex items-center border-b-2 px-1 pt-1 text-sm font-medium" aria-current="page">
                  Namespaces
                </a>
              </div>
            </div>
          </div>
          </Panel>
          <Controls />
          <MiniMap />
          <Background variant="dots" gap={12} size={1} />
        </ReactFlow>
      </div>
    );
  }
}

export default Home;
