// Copyright (c) Autovia GmbH
// SPDX-License-Identifier: Apache-2.0

import React, {Component} from 'react';
import ReactFlow, { ConnectionLineType, Panel, Controls, Background, MiniMap} from 'reactflow';
import {useParams} from 'react-router-dom';
import {nodeTypes} from './NodeTypes';
import ListView from './ListView';
import SearchView from './SearchView';
import ResourceOverlay from './ResourceOverlay';
import Nav from './Nav';
import * as k8s from './utils/k8s';
import {layout} from './utils/graph';

import 'reactflow/dist/style.css';
import './Namespace.css';

//import './index.css';
    
class Namespace extends Component {
  constructor(props) {
    super(props);
    this.state = {
      nodes: [],
      edges: [],
      namespaces: [], 
      list: {view: false, kind: "", label: ""},
      search: {view: false, filter: ""},
      overlay: {view: false, kind: "", label: ""},
      data: [],
      params: {
        kind: typeof props.params.kind === "undefined" ? "pod" : props.params.kind, 
        namespace: typeof props.params.namespace === "undefined" ? "default" : props.params.namespace,
        cluster: typeof props.params.cluster === "undefined" ? "" : props.params.cluster,
        resource: typeof props.params.resource === "undefined" ? "" : props.params.resource
      }
    }
  }

  componentDidMount() {
    console.log("Namespace path", window.location.pathname.split("/"));
    
    fetch('/api/v1/namespaces', {
      method: 'GET',
      headers: {
          'Accept': 'application/json',
          'Content-Type': 'application/json'
      }
    })
    .then(res => res.json())
    .then(d => {
      console.log('/api/v1/namespaces', d);
      const namespaces = d.nodes.map((n) => n.data.label);
      if (namespaces.length > 0) {
        this.setState((state, props) => ({
          // nodes: state.nodes,
          // edges: state.edges,
          namespaces: namespaces
        }));
      }
    });
    switch(true) {
      case this.isView("cluster"):
        this.openListView({id: this.props.params.cluster, name: k8s.clusterResources.filter((f) => f.id === this.props.params.cluster)[0].name});
        break;
      case this.isView("node"):
        this.getGroup('/api/v1/node/' + this.props.params.node);
        break;
      case this.isView("resource"):
        this.getGraph('/api/v1/graph/' + this.props.params.kind + '/' + this.props.params.namespace + '/' + this.props.params.resource);
        break;
      default:
        this.getGroup('/api/v1/' + this.state.params.kind + "?namespace=" + this.state.params.namespace);
    } 
  }

  isView(name) {
    let path = window.location.pathname.split("/");
    return path.length > 1 && path[1] === name ? true : false;
  }

  getGraph(url) {
    fetch(url)
    .then(res => res.json())
    .then(d => {
      console.log("getGraph", url, d);
      // d.nodes = d.nodes.map(node => {
      //   const newLabel = (<div className="on-hover">
      //     {process(node.data.label)}
      //     {node.data.pathMappings && <div className="on-hover-child"><div>pathMappings:</div>{process(node.data.pathMappings)}</div>}
      //   </div>);
      //   return { ...node, data: { ...node.data, label: newLabel } };
      // });

      const {nodes: layoutedNodes, edges: layoutedEdges} = layout.graph(d.nodes, d.edges);
      this.setState((state, props) => ({
        nodes: layoutedNodes,
        edges: layoutedEdges,
        // namespaces: state.params.namespaces //state.namespaces,
        // pods: state.pods
      }));
    });
  }

  getGroup(url, filter = "") {
    if (filter != "") {
      url += "&filter=" + filter;
    }

    fetch(url, {
      method: 'GET',
      headers: {
          'Accept': 'application/json',
          'Content-Type': 'application/json'
      }
    })
    .then(res => res.json())
    .then(d => {
      console.log("getGroup", url, d);
      const {nodes: layoutedNodes, edges: layoutedEdges} = layout.group(d.nodes, d.edges);
      this.setState((state, props) => ({
        nodes: layoutedNodes,
        edges: layoutedEdges,
        //namespaces: state.params.namespaces
      }));
    });
  }

  openListView = (e) => {
    this.setState((state, props) => ({
      list: {view: true, kind: e.id, label: e.name},
    }));
  }

  closeListView = (e) => {
    this.setState((state, props) => ({
      list: {view: false}
    }));
  }

  openSearchView = (e) => {
    this.setState((state, props) => ({
      search: {view: e.filter === "" ? false : true, filter: e.filter}
    }));
  }

  closeSearchView = (e) => {
    console.log("closeSearchView", e);
    this.setState((state, props) => ({
      search: {view: false, filter: ""}
    }));
  }

  onNodeClick = (e, node) => {
    if (node.data.kind === "vol") {
      this.setState((state, props) => ({
        //resourceOverlay: !this.state.resourceOverlay,
        overlay: {view: true, kind: "", label: ""},
        //url: '/api/v1/' + node.data.kind + '/'  + this.props.params.namespace + '/' + this.props.params.pod + '/'  + node.data.label,
        //params: this.props.params,
        data: node.data
      }));
    } else {
      this.setState((state, props) => ({
        //resourceOverlay: !this.state.resourceOverlay,
        overlay: {view: true, kind: "", label: ""},
        //url: '/api/v1/' + node.data.kind + '/'  + this.props.params.namespace + '/' + node.data.label,
        //params: this.props.params,
        data: node.data
      }));
    }
    //window.open("/namespace/" + node.data.namespace + "/" + node.data.kind + "/" + node.data.label, "_self");
  };

  goHome() {
    window.open("/", "_self");
  }

  closeResourceOverlay = (e) => {
    this.setState((state, props) => ({
      overlay: {view: false, kind: "", label: ""}
    }));
  }

  componentDidUpdate(prevProps, prevState) {
    if(prevState !== this.state) {
      console.log('namespace componentDidUpdate: ', this.state);
    }
  }

  render() {
    return (
      <div style={{ width: '100vw', height: '100vh' }}>
        <Nav params={this.state.params} onListClick={this.openListView} filter={this.state.search.filter} close={this.closeSearchView} onSearchClick={this.openSearchView} />
        {this.state.list.view 
        ? <ListView meta={this.state.list} close={this.closeListView} /> 
        : this.state.search.view
          ? <SearchView filter={this.state.search.filter} close2={this.closeSearchView} /> 
          : <ReactFlow
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
          </Panel>
          <Controls />
          <MiniMap />
          <Background variant="dots" gap={12} size={1} />
        </ReactFlow>
      } {this.state.overlay.view ? <ResourceOverlay data={this.state.data} params={this.state.params} close={this.closeResourceOverlay} /> : ""}
      </div>
    );
  }
}

function withParams(Component) {
  return (props) => <Component {...props} params={useParams()} />;
}

export default withParams(Namespace);
