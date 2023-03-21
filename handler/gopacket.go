package handler

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"

	"topological_graph/model"
	"topological_graph/util"
)

var (
	device      string        = "en0"            //	网络设备
	snapshotLen int32         = 1024             // 每个数据包读取的最大长度 the maximum size to read for each packet
	promiscuous bool          = true             // 是否将网口设置为混杂模式,即是否接收目的地址不为本机的包
	timeout     time.Duration = -1 * time.Second // 设置抓到包返回的超时，-1为立即返回
	handle      *pcap.Handle
	err         error
	recordData  []*model.ConnectData
	captureTime time.Duration = time.Second * 10
)

//func FindAllDevice() {
//	// 得到所有的(网络)设备
//	devices, err := pcap.FindAllDevs()
//	if err != nil {
//		log.Fatal(err)
//	}
//	// 打印设备信息
//	fmt.Println("Devices found:")
//	for _, device := range devices {
//		fmt.Println("\nName: ", device.Name)
//		fmt.Println("Description: ", device.Description)
//		fmt.Println("Devices addresses: ", device.Description)
//		for _, address := range device.Addresses {
//			fmt.Println("- IP address: ", address.IP)
//			fmt.Println("- Subnet mask: ", address.Netmask)
//		}
//	}
//}

func GetPacketData(port string) []*model.ConnectData {
	// 打开网络接口
	handle, err = pcap.OpenLive(device, snapshotLen, promiscuous, timeout)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	// 设置过滤器，只接收端口为port的TCP数据包
	var filter = "tcp"
	err = handle.SetBPFFilter(filter)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Only capturing TCP port " + port + " packets")

	// 循环抓包
	endTime := time.Now().Add(captureTime)
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {

		// 结束
		if !endTime.After(time.Now()) {
			break
		}

		// 获取TCP层
		tcpLayer := packet.Layer(layers.LayerTypeTCP)
		if tcpLayer == nil {
			log.Fatalln("nil tcp")
		}
		tcp, _ := tcpLayer.(*layers.TCP)
		if !tcp.SYN || !tcp.ACK || tcp.SrcPort != 80 {
			continue
		}

		// 获取ip层信息
		ipv4Layer := packet.Layer(layers.LayerTypeIPv4)
		if ipv4Layer == nil {
			log.Fatalln("nil ipv4")
		}
		ipv4, _ := ipv4Layer.(*layers.IPv4)
		if fmt.Sprintf("%s", ipv4.DstIP) != "10.21.168.4" {
			fmt.Println("ClientIP = ", fmt.Sprintf("%s", ipv4.DstIP)+"and ServerIP = "+fmt.Sprintf("%s", ipv4.SrcIP))
		}

		// 根据tcp三次握手，syn+ack是服务端到客户端的包
		current := &model.ConnectData{
			ServerPort: strconv.Itoa(int(tcp.SrcPort)),
			ClientPort: strconv.Itoa(int(tcp.DstPort)),
			ServerIP:   fmt.Sprintf("%s", ipv4.SrcIP),
			ClientIP:   fmt.Sprintf("%s", ipv4.DstIP),
		}
		recordData = append(recordData, current)
	}

	fmt.Println("链接数" + strconv.Itoa(len(recordData)))

	util.PrintConnectData(recordData)
	return recordData
}
