// Copyright (c) Autovia GmbH
// SPDX-License-Identifier: Apache-2.0

import React, {Component} from 'react';
import ReactFlow, { ConnectionLineType, Panel, Controls, Background, MiniMap} from 'reactflow';
import {useParams} from 'react-router-dom';
import PodNode from './PodNode';
import CustomNode from './CustomNode';
import {nodeTypes} from './NodeTypes';
import ListView from './ListView';
import SearchView from './SearchView';
import Nav from './Nav';
//import 'reactflow/dist/style.css';

//import './index.css';

class Namespace extends Component {
  constructor(props) {
    super(props);
    this.state = {
      nodes: [],
      edges: [],
      namespaces: [], 
      namespace: typeof this.props.params.namespace === "undefined" ? "default" : this.props.params.namespace,
      kind: this.props.params.kind === undefined ? "pod" : this.props.params.kind,
      list: {view: false, kind: "", label: ""},
      search: {view: false, filter: ""},
      data: []
    }
  }

  componentDidMount() {
    console.log("this.props", this.props);
    
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
          nodes: state.nodes,
          edges: state.edges,
          namespaces: namespaces
        }));
      }
    });
    this.getResources();
  }

  getResources(filter = "") {
    var url = '/api/v1/' + this.state.kind + "?namespace=" + this.state.namespace;
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
      console.log('/api/v1/namespace/', d);
      const {nodes: layoutedNodes, edges: layoutedEdges} = this.getLayoutedGroup(d.nodes, d.edges);
      this.setState((state, props) => ({
        nodes: layoutedNodes,
        edges: layoutedEdges,
        namespaces: state.namespaces
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
    window.open("/namespace/" + node.data.namespace + "/" + node.data.kind + "/" + node.data.label, "_self");
  };

  goHome() {
    window.open("/", "_self");
  }

  getLayoutedGroup(nodes, edges, columns) {
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
  };

  componentDidUpdate(prevProps, prevState) {
    if(prevState !== this.state) {
      console.log('namespace componentDidUpdate: ', this.state);
    }
  }

  render() {
    return (
      <div style={{ width: '100vw', height: '100vh' }}>
        <Nav params={this.props.params} onListClick={this.openListView} filter={this.state.search.filter} close={this.closeSearchView} onSearchClick={this.openSearchView} />
        {this.state.list.view 
        ? <ListView data={this.state.data} meta={this.state.list} close={this.closeListView} /> 
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
      }  
      </div>
    );
  }
}

function withParams(Component) {
  return (props) => <Component {...props} params={useParams()} />;
}

export default withParams(Namespace);
