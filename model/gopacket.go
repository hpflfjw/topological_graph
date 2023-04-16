package model

type ConnectData struct {
	ServerPort string
	ClientPort string
	ServerIP   string
	ClientIP   string
}

type InterfaceInfo struct {
	Name       string
	MacAddress string
	IP         string
	Netmask    string
}
