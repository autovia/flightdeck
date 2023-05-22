// Copyright (c) Autovia GmbH
// SPDX-License-Identifier: Apache-2.0

package structs

import "fmt"

type Graph struct {
	Nodes []Node `json:"nodes"`
	Edges []Edge `json:"edges"`
}

type Node struct {
	Id   string `json:"id,omitempty"`
	Data struct {
		Label      string   `json:"label,omitempty"`
		Group      bool     `json:"group,omitempty"`
		Kind       string   `json:"kind,omitempty"`
		Namespace  string   `json:"namespace,omitempty"`
		Containers []string `json:"containers,omitempty"`
		Status     string   `json:"status,omitempty"`
		Age        string   `json:"age,omitempty"`
		Ready      string   `json:"ready,omitempty"`
		Restarts   string   `json:"restarts,omitempty"`
	} `json:"data,omitempty"`
	// optional
	Position string `json:"position,omitempty"`
	Style    struct {
		Width  int `json:"width,omitempty"`
		Height int `json:"height,omitempty"`
	} `json:"style,omitempty"`
	Type        string `json:"type,omitempty"`
	ParentNode  string `json:"parentNode,omitempty"`
	Extent      string `json:"extent,omitempty"`
	Draggable   bool   `json:"draggable"`
	Connectable bool   `json:"connectable"`
}

type NodeOptions struct {
	Type        string
	Draggable   bool
	Connectable bool
	ParentNode  Node
	Extent      string
	Namespace   string
	Group       bool
	Containers  []string
	Status      string
	Age         string
	Ready       string
	Restarts    string
}

type Edge struct {
	Id     string `json:"id,omitempty"`
	Source string `json:"source,omitempty"`
	Target string `json:"target,omitempty"`
}

func (g *Graph) AddNode(kind string, id string, name string, args NodeOptions) Node {
	node := Node{Id: fmt.Sprintf("%s/%s", kind, id)}
	node.Data.Label = name
	node.Data.Kind = kind
	node.Data.Namespace = args.Namespace
	node.Data.Containers = args.Containers
	node.Data.Status = args.Status
	node.Data.Age = args.Age
	node.Data.Ready = args.Ready
	node.Data.Restarts = args.Restarts
	node.Draggable = args.Draggable || false
	node.Connectable = args.Connectable || false
	node.Type = args.Type
	if args.Group {
		node.Style.Width = 0
		node.Style.Height = 0
		node.Data.Group = true
	}
	node.ParentNode = args.ParentNode.Id
	node.Extent = args.Extent
	g.Nodes = append(g.Nodes, node)
	return node
}

func (g *Graph) AddEdge(source Node, target Node) {
	edge := Edge{
		Id:     fmt.Sprintf("%s-%s", source.Id, target.Id),
		Source: source.Id,
		Target: target.Id,
	}
	g.Edges = append(g.Edges, edge)
}

func (g *Graph) Includes(name string) bool {
	for _, n := range g.Nodes {
		if n.Data.Label == name {
			return true
		}
	}
	return false
}
