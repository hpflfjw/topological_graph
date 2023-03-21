package util

import (
	"fmt"

	"topological_graph/model"
)

// 通过将结构体里的所有字段生成一个hash来去重

func DistinctConnectData(connectData []model.ConnectData) []model.ConnectData {
	mapData := map[string]model.ConnectData{}
	distinctData := []model.ConnectData{}

	for _, data := range connectData {
		hashStr := GetHashByConnectData(data)
		if _, ok := mapData[hashStr]; !ok {
			mapData[hashStr] = data
			distinctData = append(distinctData, data)
		}
	}
	return distinctData
}

func GetHashByConnectData(data model.ConnectData) string {
	return fmt.Sprintf("%s%s%s%s", data.ServerPort, data.ClientPort, data.ServerIP, data.ClientIP)
}

func PrintConnectData(data []*model.ConnectData) {
	for _, d := range data {
		fmt.Println(fmt.Sprintf("%+v", d))
	}
}
