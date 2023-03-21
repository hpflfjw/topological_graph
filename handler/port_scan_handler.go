package handler

import (
	"flag"
	"fmt"
	"net"
	"sort"
	"sync"
)

var startPort = flag.Int("start", 1, "Start port")
var endPort = flag.Int("end", 65535, "End port")
var ipAddress = flag.String("ip", "127.0.0.1", "ipAddress")

func GetIpAddress() []string {
	// 获取本地主机的所有 IP 地址
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	ipSlices := []string{}

	// 遍历所有 IP 地址
	for _, a := range addrs {
		// 将 IP 地址转换为字符串
		if ipNet, ok := a.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				fmt.Println(ipNet.IP.String())
				ipSlices = append(ipSlices, ipNet.IP.String())
			}
		}
	}
	return ipSlices
}

func ScanPort() []int {

	openPorts := []int{}
	var wg sync.WaitGroup

	for i := 1; i <= 10000; i++ {
		wg.Add(1)
		go func(port int) {
			defer wg.Done()
			address := fmt.Sprintf("localhost:%d", port)
			conn, err := net.Dial("tcp", address)
			if err != nil {
				return
			}
			conn.Close()
			openPorts = append(openPorts, port)
			fmt.Println(address, "is open")
		}(i)
	}

	wg.Wait()
	// 排序
	sort.Ints(openPorts)

	return openPorts
}
