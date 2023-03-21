package util

import (
	"fmt"
	"github.com/goccy/go-graphviz/cgraph"
	"log"
	"os"
	"topological_graph/model"
)

type GraphHandler struct {
	digraph *cgraph.Graph
	nodes   map[string]*cgraph.Node
	edges   []*cgraph.Edge
	err     error
}

func NewGraphHandler() *GraphHandler {
	return &GraphHandler{
		digraph: nil,
		nodes:   map[string]*cgraph.Node{},
		edges:   []*cgraph.Edge{},
		err:     nil,
	}
}

func (h *GraphHandler) MakeGraphByConnectData(datas []*model.ConnectData) {
	g := graphviz.New()
	defer g.Close()

	// 创建有向图
	h.digraph, h.err = g.Graph()

	// 节点、边
	for _, data := range datas {
		h.CreateEdgeByConnectData(data)
	}

	// 生成图像
	if err := g.RenderFilename(h.digraph, graphviz.PNG, "example.png"); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

}

func (h *GraphHandler) CreateEdgeByConnectData(data *model.ConnectData) {
	SrcNode := h.MakeOrGetNode(fmt.Sprintf(data.ServerIP))
	DstNode := h.MakeOrGetNode(fmt.Sprintf(data.ClientIP))

	h.MakeEdgeByNodes("", DstNode, SrcNode)
}

func (h *GraphHandler) MakeOrGetNode(ip string) *cgraph.Node {
	if _, ok := h.nodes[ip]; ok {
		return h.nodes[ip]
	} else {
		node, err := h.digraph.CreateNode(ip)
		if err != nil {
			h.err = err
			log.Fatalf("Failed to create node, ip = %s", ip)
			return nil
		}
		h.nodes[ip] = node
		return node
	}
}

func (h *GraphHandler) MakeEdgeByNodes(name string, srcNode *cgraph.Node, dstNode *cgraph.Node) {
	edge, err := h.digraph.CreateEdge(name, srcNode, dstNode)
	if err != nil {
		log.Fatalf("Failed to create edge, name = %s, nodeA = %s, nodeB = %s", name, srcNode.Name(), dstNode.Name())
		return
	}
	h.edges = append(h.edges, edge)
}
