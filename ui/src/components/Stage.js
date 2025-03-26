// Copyright (c) Autovia GmbH
// SPDX-License-Identifier: Apache-2.0

import React, { Component } from "react";
import {
  ReactFlow,
  ConnectionLineType,
  Panel,
  Controls,
  Background,
  MiniMap,
} from "@xyflow/react";

import { useParams } from "react-router-dom";
import { nodeTypes } from "./utils/NodeTypes";
import { edgeTypes } from "./utils/EdgeTypes";
import ListView from "./List";
import SearchView from "./Search";
import ResourceOverlay from "./Resource";
import Nav from "./partials/Nav";
import * as k8s from "./utils/K8s";
import { layout } from "./utils/GraphLayout";

//import "@xyflow/react/dist/base.css";
import "@xyflow/react/dist/style.css";
//import "./css/Stage.css";

class Stage extends Component {
  constructor(props) {
    super(props);
    this.state = {
      nodes: [],
      edges: [],
      namespaces: [],
      list: { view: false, kind: "", label: "" },
      search: { view: false, filter: "" },
      overlay: { view: false },
      data: [],
      params: {
        kind:
          typeof props.params.kind === "undefined" ? "pod" : props.params.kind,
        namespace:
          typeof props.params.namespace === "undefined"
            ? "default"
            : props.params.namespace,
        cluster:
          typeof props.params.cluster === "undefined"
            ? ""
            : props.params.cluster,
        resource:
          typeof props.params.resource === "undefined"
            ? ""
            : props.params.resource,
        type: this.isView("clusterresource") ? "clusterresource" : "resource",
      },
      clusterInfo: {},
    };
  }

  componentDidMount() {
    console.log("Stage componentDidMount", window.location.pathname.split("/"));
    switch (true) {
      case this.isView("cluster"):
        this.openListView({
          id: this.props.params.cluster,
          name: k8s.cluster.filter((f) => f.id === this.props.params.cluster)[0]
            .name,
        });
        break;
      case this.isView("clusterresource"):
        this.setState((state, props) => ({
          overlay: { view: true },
          data: [],
        }));
        break;
      case this.isView("node"):
        this.getGroup("/api/v1/node/" + this.props.params.node);
        break;
      case this.isView("resource"):
        this.getGraph(
          "/api/v1/graph/" +
            this.props.params.kind +
            "/" +
            this.props.params.namespace +
            "/" +
            this.props.params.resource,
        );
        break;
      default:
        this.getGroup(
          "/api/v1/resources/" +
            this.state.params.namespace +
            "/" +
            this.state.params.kind,
        );
    }
    this.getClusterInfo();
  }

  isView(name) {
    let path = window.location.pathname.split("/");
    return path.length > 1 && path[1] === name ? true : false;
  }

  getClusterInfo() {
    fetch("/api/v1/cluster/info")
      .then((res) => res.json())
      .then((d) => {
        this.setState((state, props) => ({
          clusterInfo: d,
        }));
      });
  }

  getGraph(url) {
    fetch(url)
      .then((res) => res.json())
      .then((d) => {
        console.log("getGraph", url, d);
        // d.nodes = d.nodes.map(node => {
        //   const newLabel = (<div className="on-hover">
        //     {process(node.data.label)}
        //     {node.data.pathMappings && <div className="on-hover-child"><div>pathMappings:</div>{process(node.data.pathMappings)}</div>}
        //   </div>);
        //   return { ...node, data: { ...node.data, label: newLabel } };
        // });
        var direction = "TB";
        if (
          this.props.params.kind === "netpol" &&
          this.props.params.resources !== ""
        ) {
          direction = d.direction;
        }

        const { nodes: layoutedNodes, edges: layoutedEdges } = layout.graph(
          d.nodes,
          d.edges,
          direction,
        );
        this.setState((state, props) => ({
          nodes: layoutedNodes,
          edges: layoutedEdges,
        }));
      });
  }

  getGroup(url, filter = "") {
    if (filter != "") {
      url += "&filter=" + filter;
    }

    fetch(url, {
      method: "GET",
      headers: {
        Accept: "application/json",
        "Content-Type": "application/json",
      },
    })
      .then((res) => res.json())
      .then((d) => {
        console.log("getGroup", url, d);
        const { nodes: layoutedNodes, edges: layoutedEdges } = layout.group(
          d.nodes,
          d.edges,
        );
        this.setState((state, props) => ({
          nodes: layoutedNodes,
          edges: layoutedEdges,
        }));
      });
  }

  openListView = (e) => {
    this.setState((state, props) => ({
      list: { view: true, kind: e.id, label: e.name },
      search: { view: false, filter: "" },
    }));
  };

  closeListView = (e) => {
    this.setState((state, props) => ({
      list: { view: false },
    }));
  };

  openSearchView = (e) => {
    this.setState((state, props) => ({
      search: { view: e.filter === "" ? false : true, filter: e.filter },
      list: { view: false },
    }));
  };

  closeSearchView = (e) => {
    console.log("closeSearchView", e);
    this.setState((state, props) => ({
      search: { view: false, filter: "" },
    }));
  };

  onNodeClick = (e, node) => {
    if (node.data.kind === "vol") {
      this.setState((state, props) => ({
        overlay: { view: true },
        data: node.data,
      }));
    } else {
      this.setState((state, props) => ({
        overlay: { view: true },
        data: node.data,
      }));
    }
  };

  closeResourceOverlay = (e) => {
    this.setState((state, props) => ({
      overlay: { view: false },
    }));
  };

  componentDidUpdate(prevProps, prevState) {
    if (prevState !== this.state) {
      console.log("Stage componentDidUpdate: ", this.state);
    }
  }

  render() {
    return (
      <div style={{ width: "100vw", height: "90vh" }}>
        <Nav
          params={this.state.params}
          onListClick={this.openListView}
          filter={this.state.search.filter}
          close={this.closeSearchView}
          onSearchClick={this.openSearchView}
        />
        {this.state.list.view ? (
          <ListView meta={this.state.list} close={this.closeListView} />
        ) : this.state.search.view ? (
          <SearchView
            filter={this.state.search.filter}
            close={this.closeSearchView}
          />
        ) : (
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
            edgeTypes={edgeTypes}
            fitView
            className="bg-sky-50"
          >
            <Panel position="top-left" className="w-full p-0 m-0"></Panel>
            <Controls />
            <MiniMap />
            <Background variant="dots" gap={12} size={1} />
          </ReactFlow>
        )}{" "}
        {this.state.overlay.view ? (
          <ResourceOverlay
            overlay={true}
            data={this.state.data}
            close={this.closeResourceOverlay}
          />
        ) : (
          ""
        )}
        <div className="text-sm px-4">
          Cluster: {this.state.clusterInfo.gitVersion} connected
        </div>
      </div>
    );
  }
}

function withParams(Component) {
  return (props) => <Component {...props} params={useParams()} />;
}

export default withParams(Stage);
