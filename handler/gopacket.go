package handler

import (
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"

	"go.uber.org/zap"
	"topological_graph/model"
)

var (
	snapshotLen int32         = 1024             // 每个数据包读取的最大长度 the maximum size to read for each packet
	promiscuous bool          = true             // 是否将网口设置为混杂模式,即是否接收目的地址不为本机的包
	timeout     time.Duration = -1 * time.Second // 设置抓到包返回的超时，-1为立即返回
	handle      *pcap.Handle
	err         error
	recordData  []*model.ConnectData
	captureTime time.Duration = time.Second * 10
	logger      *zap.Logger
)

func init() {
	logger, err = zap.NewProduction()
	if err != nil {
		fmt.Println("Failed to init logger")
		return
	}
}

func GetPacketData(port string) []*model.ConnectData {

	// 获取网络设备
	deviceInfo, err := GetConfig()
	if err != nil {
		logger.Error("GetConfig error", zap.Error(err))
	}
	if len(deviceInfo) == 0 {
		logger.Error("No device found")
	}

	// 打开网络接口
	handle, err = pcap.OpenLive(deviceInfo[0].Name, snapshotLen, promiscuous, timeout)
	if err != nil {
		logger.Error("Failed to open live", zap.Error(err))
	}
	defer handle.Close()

	// 设置过滤器，只接收端口为port的TCP数据包
	var filter = "tcp"
	err = handle.SetBPFFilter(filter)
	if err != nil {
		logger.Error("Failed to set BPF filter", zap.Error(err))
	}
	logger.Info("Only capturing TCP port " + port + " packets")

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
			logger.Error("nil tcp")
		}
		tcp, _ := tcpLayer.(*layers.TCP)
		if !tcp.SYN || !tcp.ACK || tcp.SrcPort != 80 {
			continue
		}

		// 获取ip层信息
		ipv4Layer := packet.Layer(layers.LayerTypeIPv4)
		if ipv4Layer == nil {
			logger.Error("nil ipv4")
		}
		ipv4, _ := ipv4Layer.(*layers.IPv4)

		// 根据tcp三次握手，syn+ack是服务端到客户端的包
		current := &model.ConnectData{
			ServerPort: strconv.Itoa(int(tcp.SrcPort)),
			ClientPort: strconv.Itoa(int(tcp.DstPort)),
			ServerIP:   fmt.Sprintf("%s", ipv4.SrcIP),
			ClientIP:   fmt.Sprintf("%s", ipv4.DstIP),
		}
		logger.Info("Captured packet",
			zap.String("ServerIP", current.ServerIP),
			zap.String("ServerPort", current.ServerPort),
			zap.String("ClientIP", current.ClientIP),
			zap.String("ClientPort", current.ClientPort),
		)
		recordData = append(recordData, current)
	}

	logger.Info("Capture finished, total " + strconv.Itoa(len(recordData)) + " packets")

	return recordData
}

func GetConfig() ([]*model.InterfaceInfo, error) {
	cmd := exec.Command("ifconfig")
	output, err := cmd.Output()
	if err != nil {
		logger.Error("Command execution failed", zap.Error(err))
		return nil, err
	}

	outputString := string(output)

	re := regexp.MustCompile(`(?s)(\w+):\s.*?ether\s(.+?)\s.*?inet\s(.+?)\s.*?netmask\s(.+?)\s`)
	matches := re.FindAllStringSubmatch(outputString, -1)

	var ret []*model.InterfaceInfo
	for _, match := range matches {
		interfaceInfo := &model.InterfaceInfo{
			Name:       match[1],
			MacAddress: match[2],
			IP:         match[3],
			Netmask:    match[4],
		}
		logger.Info("Interface info",
			zap.String("Name", interfaceInfo.Name),
			zap.String("MacAddress", interfaceInfo.MacAddress),
			zap.String("IP", interfaceInfo.IP),
			zap.String("Netmask", interfaceInfo.Netmask),
		)
	}
	return ret, nil
}
