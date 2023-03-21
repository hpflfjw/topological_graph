package handler

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func TestRpc(ctx *gin.Context) {

	resp, err := http.Get("127.0.0.1/ping")
	if err != nil {
		log.Println("Failed to Get resp from rpc test")
		return
	}
	var m map[string]string
	dataByte, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(dataByte, &m)
	if err != nil {
		log.Println("Failed to unmarshal body to byte")
		return
	}
	log.Printf("%+v", m)

	ctx.Data(200, "json", dataByte)

}
