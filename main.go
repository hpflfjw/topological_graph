package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"topological_graph/handler"
	"topological_graph/util"
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

	gHandler := util.NewGraphHandler()

	// 扫描所有本机的端口
	ports := handler.ScanPort()
	fmt.Println(len(ports))

	// 抓包
	connectData := handler.GetPacketData("80")

	// 生成有向图
	gHandler.MakeGraphByConnectData(connectData)

	////转成csv
	//util.ToCsvFile(80, connectData)

}
