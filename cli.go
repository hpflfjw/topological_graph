package main

import (
	"fmt"
	"strconv"
	"topological_graph/handler"
	"topological_graph/util"
)

type CliHandler struct {
	Command string
	Args    []string
}

func NewCliHandler(command string, args []string) *CliHandler {
	return &CliHandler{
		Command: command,
		Args:    args,
	}
}

func (h *CliHandler) Run() error {
	switch h.Command {
	case "help":
		return h.RunHelp()
	case "port_scan":
		return h.RunPort()
	case "connect_data":
		return h.RunConnectData()
	case "graph":
		return h.RunGraph()
	default:
		fmt.Println("topological: no command specified. Use help for help.")
		return fmt.Errorf("command not found")
	}
}

func (h *CliHandler) RunHelp() error {
	fmt.Println("Usage:")
	fmt.Println("  topological command [arguments]")
	fmt.Println("")
	fmt.Println("Commands:")
	fmt.Println("  port_scan -s -e      展示端口开启的信息")
	fmt.Println("  connect_data      展示连接的信息，包括源IP，源端口，目的IP，目的端口，通过csv文件展示")
	fmt.Println("  graph      绘制拓扑图信息，通过图片展示")
	return nil
}

func (h *CliHandler) RunPort() error {

	if len(h.Args) < 2 {
		fmt.Println("topological: no command specified. Use help.")
		return nil
	}
	startPort, err := strconv.Atoi(h.Args[0])
	if err != nil {
		fmt.Println(h.Args[0] + " is not a number")
		return err
	}

	endPort, err := strconv.Atoi(h.Args[1])
	if err != nil {
		fmt.Println(h.Args[1] + " is not a number")
		return err
	}

	fmt.Println(fmt.Sprintf("port_scan 开启的端口信息从%d - %d:", startPort, endPort))

	handler.ScanPort(startPort, endPort)

	return nil

}

func (h *CliHandler) RunConnectData() error {

	port := h.Args[0]

	fmt.Println(fmt.Sprintf("connect_data 扫描的端口为%s,连接信息:", port))

	handler.GetPacketData(port)
	return nil
}

func (h *CliHandler) RunGraph() error {

	port := h.Args[0]

	fmt.Println(fmt.Sprintf("graph 扫描的端口为:%s， 将会生成拓扑图", port))

	connectData := handler.GetPacketData(port)

	// 生成有向图
	util.NewGraphHandler().MakeGraphByConnectData(connectData)
	return nil
}
