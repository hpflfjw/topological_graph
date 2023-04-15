package prometheus

import (
	"fmt"
	"net/http"
	"os"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// 全局变量，用于存储指标
var (
	vpcCount = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "vpc_count",
		Help: "Number of VPCs in the account",
	})
)

func init() {
	// 注册指标
	prometheus.MustRegister(vpcCount)
}

func main() {
	// 设置环境变量
	accessKeyID := os.Getenv("ALICLOUD_ACCESS_KEY")
	accessKeySecret := os.Getenv("ALICLOUD_SECRET_KEY")

	// 创建阿里云客户端
	client, err := sdk.NewClientWithAccessKey("cn-hangzhou", accessKeyID, accessKeySecret)
	if err != nil {
		panic(err)
	}

	// 创建 VPC 客户端
	vpcClient, err := vpc.NewClientWithAccessKey("cn-hangzhou", accessKeyID, accessKeySecret)
	if err != nil {
		panic(err)
	}

	// 创建一个 HTTP 服务器，提供 /metrics 端点
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/probe", func(w http.ResponseWriter, r *http.Request) {
		probeHandler(w, r, client, vpcClient)
	})

	// 启动 HTTP 服务器
	fmt.Println("Starting server on :2112")
	http.ListenAndServe(":2112", nil)
}

func probeHandler(w http.ResponseWriter, r *http.Request, client *sdk.Client, vpcClient *vpc.Client) {
	// 获取 VPC 信息
	vpcRequest := vpc.CreateDescribeVpcsRequest()
	vpcRequest.Scheme = requests.HTTPS
	vpcRequest.PageSize = requests.NewInteger(50)

	vpcResponse, err := vpcClient.DescribeVpcs(vpcRequest)
	if err != nil {
		fmt.Println("Failed to fetch VPC information")
		return
	}

	// 设置 VPC 数量指标
	vpcCount.Set(float64(len(vpcResponse.Vpcs.Vpc)))

	// 写入指标
	promhttp.Handler().ServeHTTP(w, r)
}
