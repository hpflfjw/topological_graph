package handler

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	vpcGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "vpc",
			Help: "VPC information",
		},
		[]string{"region", "vpc_id", "vpc_name"},
	)

	subnetGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "subnet",
			Help: "Subnet information",
		},
		[]string{"region", "subnet_id", "subnet_name", "vpc_id", "vpc_name"},
	)

	routeTableGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "route_table",
			Help: "Route table information",
		},
		[]string{"region", "route_table_id", "route_table_name", "vpc_id", "vpc_name"},
	)

	securityGroupGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "security_group",
			Help: "Security group information",
		},
		[]string{"region", "security_group_id", "security_group_name", "vpc_id", "vpc_name"},
	)
)

func init() {
	prometheus.MustRegister(vpcGauge)
	prometheus.MustRegister(subnetGauge)
	prometheus.MustRegister(routeTableGauge)
	prometheus.MustRegister(securityGroupGauge)
}

func main() {
	http.Handle("/metrics", promhttp.Handler())

	go func() {
		for {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			client, err := vpc.NewClientWithAccessKey("cn-hangzhou", os.Getenv("ALIYUN_ACCESS_KEY_ID"), os.Getenv("ALIYUN_ACCESS_KEY_SECRET"))
			if err != nil {
				fmt.Println("Failed to create vpc client:", err)
				continue
			}

			err = updateVpcMetrics(ctx, client)
			if err != nil {
				fmt.Println("Failed to update vpc metrics:", err)
			}

			//err = updateSubnetMetrics(ctx, client)
			//if err != nil {
			//	fmt.Println("Failed to update subnet metrics:", err)
			//}

			err = updateRouteTableMetrics(ctx, client)
			if err != nil {
				fmt.Println("Failed to update route table metrics:", err)
			}

			//err = updateSecurityGroupMetrics(ctx, client)
			//if err != nil {
			//	fmt.Println("Failed to update security group metrics:", err)
			//}

			time.Sleep(30 * time.Second)
		}
	}()

	err := http.ListenAndServe(":9100", nil)
	if err != nil {
		fmt.Println("Failed to start http server:", err)
	}
}

// updateVpcMetrics 更新VPC信息
func updateVpcMetrics(ctx context.Context, client *vpc.Client) error {
	req := vpc.CreateDescribeVpcsRequest()
	resp, err := client.DescribeVpcs(req)
	if err != nil {
		return err
	}

	for _, vpc := range resp.Vpcs.Vpc {
		vpcGauge.WithLabelValues(resp.RequestId, vpc.VpcId, vpc.VpcName).Set(1)
	}

	return nil
}

// updateRouteTableMetrics 更新路由信息
func updateRouteTableMetrics(ctx context.Context, client *vpc.Client) error {
	req := vpc.CreateDescribeRouteTablesRequest()
	resp, err := client.DescribeRouteTables(req)
	if err != nil {
		return err
	}

	for _, routeTable := range resp.RouteTables.RouteTable {
		routeTableGauge.WithLabelValues(resp.RequestId, routeTable.RouteTableId).Set(1)
	}

	return nil
}

// updateNetworkAclMetrics 更新ACL信息
//func updateNetworkAclMetrics(ctx context.Context, client *vpc.Client) error {
//	req := vpc.CreateDescribeNetworkAclAttributesRequest()
//	resp, err := client.DescribeNetworkAclAttributes(req)
//	if err != nil {
//		return err
//	}
//
//	for _, routeTable := range resp.NetworkAclAttribute {
//		routeTableGauge.WithLabelValues(resp.RequestId, routeTable.RouteTableId).Set(1)
//	}
//
//	return nil
//}
