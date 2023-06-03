// Copyright (c) Autovia GmbH
// SPDX-License-Identifier: Apache-2.0

import React, {Component} from 'react';
import ReactFlow, { useReactFlow, ConnectionLineType, Panel, Controls, Background, MiniMap} from 'reactflow';
import {useParams} from 'react-router-dom';
import dagre from 'dagre';
import CustomNodeEdge from './CustomNodeEdge';
import CustomNode from './CustomNode';
import {nodeTypes} from './NodeTypes';
import ResourceOverlay from './ResourceOverlay';
import ListOverlay from './ListView';
import Nav from './Nav';
import { XMarkIcon, ChevronDoubleDownIcon } from '@heroicons/react/24/outline'

import 'reactflow/dist/style.css';
//import Nav from './Nav';
//import './tailwind-config.js';

import './Pod.css';

class Pod extends Component {

  constructor(props) {
    super(props);
    this.state = {
      nodes: [],
      edges: [],
      namespaces: [],
      pods: [],
      resourceOverlay: false,
      listOverlay: false,
      listOverlayMeta: {},
      data: [],
      params: {},
      menu: false
    }
  }

  componentDidMount() {
    console.log("Pod this.props", this.props);
    console.log(window.innerWidth, window.innerHeight);
    
    fetch('/api/v1/namespaces')
    .then(res => res.json())
    .then(d => {
      console.log('/api/v1/namespaces', d);
      //const namespaces = d.nodes.map((n) => n.data.label);
      if (d.nodes.length > 0) {
        this.setState((state, props) => ({
          nodes: state.nodes,
          edges: state.edges,
          namespaces: d.nodes,
          pods: state.pods
        }));
      }
    });

    fetch('/api/v1/graph/' + this.props.params.kind + '/' + this.props.params.namespace + '/' + this.props.params.pod)
    .then(res => res.json())
    .then(d => {
      console.log('/api/v1/graph/' + this.props.params.kind, d);
      // d.nodes = d.nodes.map(node => {
      //   const newLabel = (<div className="on-hover">
      //     {process(node.data.label)}
      //     {node.data.pathMappings && <div className="on-hover-child"><div>pathMappings:</div>{process(node.data.pathMappings)}</div>}
      //   </div>);
      //   return { ...node, data: { ...node.data, label: newLabel } };
      // });

      const {nodes: layoutedNodes, edges: layoutedEdges} = this.getLayoutedGraph(d.nodes, d.edges);
      this.setState((state, props) => ({
        nodes: layoutedNodes,
        edges: layoutedEdges,
        namespaces: state.namespaces,
        pods: state.pods
      }));
    });

    fetch('/api/v1/pods?namespace=' + this.props.params.namespace)
      .then(res => res.json())
      .then(d => {
        console.log('/api/v1/namespace/', d);
        //const pods = d.nodes.filter((f) => f.data.kind === "pod").map(m => m.data.label);
        const pods = d.nodes.filter((f) => f.data.kind === "pod");//.map(m => m.data.label);
        console.log("pods", pods);
        this.setState((state, props) => ({
          nodes: state.nodes,
          edges: state.edges,
          namespaces: state.namespaces,
          pods: pods
        }));
      });
  }

  // onPodClick = (e, node) => {
  //   if (node.data.kind === "pod") {
  //       window.open("/namespace/" + this.props.params.namespace + "/pod/" + node.data.label, "_self");
  //   }
  // };

  getLayoutedGraph(nodes, edges, direction = 'TB') {
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
  };

  goHome = (e) => {
    window.open("/", "_self");
  }

  closeResourceOverlay = (e) => {
    this.setState((state, props) => ({
      resourceOverlay: !this.state.resourceOverlay
    }));
  }

  closeListOverlay = (e) => {
    this.setState((state, props) => ({
      listOverlay: !this.state.listOverlay
    }));
  }

  openListOverlay = (e) => {
    this.setState((state, props) => ({
      listOverlay: !this.state.listOverlay,
      listOverlayMeta: {kind: e.id, label: e.name},
    }));
  }

  onNodeClick = (e,node) => {
    console.log(node)
    if (node.data.kind === "vol") {
      this.setState((state, props) => ({
        resourceOverlay: !this.state.resourceOverlay,
        //url: '/api/v1/' + node.data.kind + '/'  + this.props.params.namespace + '/' + this.props.params.pod + '/'  + node.data.label,
        params: this.props.params,
        data: node.data
      }));
    } else {
      this.setState((state, props) => ({
        resourceOverlay: !this.state.resourceOverlay,
        //url: '/api/v1/' + node.data.kind + '/'  + this.props.params.namespace + '/' + node.data.label,
        params: this.props.params,
        data: node.data
      }));
    }
  }

  onNamespaceChange(e) {
    window.open("/namespace/" + e.target.value + "/pod", "_self");
  }

  onPodSelect = (label) => {
    window.open("/namespace/" + this.props.params.namespace + "/pod/" + label, "_self");
  }

  componentDidUpdate(prevProps, prevState) {
    if(prevState !== this.state) {
      console.log('pod componentDidUpdate: ', this.state);
      this.forceUpdate();
    }
  }

  openMenu = () => {
    this.setState((state, props) => ({
      menu: !this.state.menu
    }));
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
          fitView={!this.state.menu}
          //fitViewOptions={{ nodes: [{id: "svc/02724bc7-b20f-40ae-b8b4-aa39cb6cf58b"}] }}
          //defaultViewport={this.defaultViewport()}
          className="bg-sky-50"
        >
          <Panel position="top-left" className="w-full p-0 m-0">
            <Nav params={this.props.params} onClick={this.openListOverlay} />
          </Panel>
          <Controls />
          <MiniMap />
          <Background variant="dots" gap={12} size={1} />
        </ReactFlow>
        {this.state.resourceOverlay ? <ResourceOverlay data={this.state.data} params={this.state.params} close={this.closeResourceOverlay} /> : ""}
        {this.state.listOverlay ? <ListOverlay data={this.state.data} meta={this.state.listOverlayMeta} close={this.closeListOverlay} /> : ""}
      </div>
    );
  }
}

function withParams(Component) {
  return (props) => <Component {...props} params={useParams()} />;
}

export default withParams(Pod);
