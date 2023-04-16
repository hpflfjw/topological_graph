package util

import (
	"encoding/csv"
	"os"
	"topological_graph/model"
)

func ToCsvFile(port string, connectDatas []*model.ConnectData) {

	file, err := os.Create("port " + port + " data.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, connectData := range connectDatas {
		err := writer.Write([]string{connectData.ServerIP, connectData.ServerPort, connectData.ClientIP, connectData.ClientPort})
		if err != nil {
			panic(err)
		}
	}
}
