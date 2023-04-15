package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
)

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		data, _ := json.Marshal(map[string]string{
			"ret": "hello world",
		})
		c.Data(200, "json", data)
	})
	//r.GET("/test_rpc", handler.TestRpc)
	return r
}

func main() {

	//graphFlag := flag.Bool("graph", false, "graph")
	//
	//flag.Parse()
	//
	//if *graphFlag {
	//	fmt.Println("Usage:")
	//	fmt.Println("  topological command [arguments]")
	//	fmt.Println("")
	//	fmt.Println("Commands:")
	//	fmt.Println("  port_scan -s -e      展示端口开启的信息")
	//	fmt.Println("  connect_data      展示连接的信息，包括源IP，源端口，目的IP，目的端口，通过csv文件展示")
	//	fmt.Println("  graph      绘制拓扑图信息，通过图片展示")
	//}

	// 用来接受命令行参数
	args := os.Args[1:]

	if len(args) < 1 {
		fmt.Println("topological: no command specified. Use help.")
		return
	}

	err := NewCliHandler(args[0], args[1:]).Run()
	if err != nil {
		fmt.Println(err)
	}

	//gHandler := util.NewGraphHandler()
	//
	//// 扫描所有本机的端口
	//ports := handler.ScanPort()
	//fmt.Println(len(ports))
	//
	//// 抓包
	//connectData := handler.GetPacketData("80")
	//
	//// 生成有向图
	//gHandler.MakeGraphByConnectData(connectData)

	//转成csv
	//util.ToCsvFile(80, connectData)

	//prometheus.ProTest()

}
