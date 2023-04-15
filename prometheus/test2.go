package prometheus

import (
	"fmt"
	"os"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cms"
)

func init() {
	os.Setenv("ALICLOUD_ACCESS_KEY", "LTAI5tCZNq2ukH1kTpDEpSHp")
	os.Setenv("ALICLOUD_SECRET_KEY", "PLWTkDLa2GRLskGjQSFNyHZHOHJFqa")
}

func ProTest() {
	accessKeyID := os.Getenv("ALICLOUD_ACCESS_KEY")
	accessKeySecret := os.Getenv("ALICLOUD_SECRET_KEY")

	// 创建阿里云客户端
	//client, err := sdk.NewClientWithAccessKey("cn-hangzhou", accessKeyID, accessKeySecret)
	//if err != nil {
	//	panic(err)
	//}

	// 创建 CloudMonitor 客户端
	cmsClient, err := cms.NewClientWithAccessKey("cn-beijing", accessKeyID, accessKeySecret)
	if err != nil {
		panic(err)
	}

	// 获取实例的网络流量信息
	getInstanceNetworkTraffic(cmsClient)
}

func getInstanceNetworkTraffic(cmsClient *cms.Client) {
	// 设置查询时间范围
	now := time.Now()
	endTime := now.Format("2006-01-02T15:04:05Z")
	startTime := now.Add(-1 * time.Hour).Format("2006-01-02T15:04:05Z")

	// 构建请求以获取实例的入流量
	inTrafficRequest := cms.CreateDescribeMetricDataRequest()
	inTrafficRequest.MetricName = "IntranetIn"
	inTrafficRequest.Namespace = "acs_ecs_dashboard"
	inTrafficRequest.StartTime = startTime
	inTrafficRequest.EndTime = endTime
	inTrafficRequest.Period = "60"
	inTrafficRequest.Scheme = requests.HTTPS

	// 发送请求
	inTrafficResponse, err := cmsClient.DescribeMetricData(inTrafficRequest)
	if err != nil {
		fmt.Println("Failed to fetch IntranetIn traffic data")
		return
	}

	// 处理入流量数据
	fmt.Println(fmt.Sprintf("IntranetIn traffic data: %+v", inTrafficResponse))
	fmt.Println(inTrafficResponse.Datapoints)

	// 构建请求以获取实例的出流量
	outTrafficRequest := cms.CreateDescribeMetricDataRequest()
	outTrafficRequest.MetricName = "IntranetOut"
	outTrafficRequest.Namespace = "acs_ecs_dashboard"
	outTrafficRequest.StartTime = startTime
	outTrafficRequest.EndTime = endTime
	outTrafficRequest.Period = "60"
	outTrafficRequest.Scheme = requests.HTTPS

	// 发送请求
	outTrafficResponse, err := cmsClient.DescribeMetricData(outTrafficRequest)
	if err != nil {
		fmt.Println("Failed to fetch IntranetOut traffic data")
		return
	}

	// 处理出流量数据
	fmt.Println(fmt.Sprintf("IntranetOut traffic data: %+v", outTrafficResponse))
	fmt.Println(outTrafficResponse.Datapoints)
}
