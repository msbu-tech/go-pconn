package main

import (
	"bytes"
	"fmt"

	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/msbu-tech/go-pconn/msg"
)

var (
	push_url = "http://127.0.0.1:8078/push"
	push_req = msg.PushMsgReq{
		Push_type:  0,
		Channel_id: "baiduyiliaoshiyebu",
		Device_ids: []string{"hanpeng", "zhangwanlong"},
		Content:    "Hello World",
	}
)

func main() {
	push_req_json, err := json.Marshal(push_req)
	if err != nil {
		panic(err)
	}
	req, err := http.NewRequest("POST", push_url, bytes.NewBuffer(push_req_json))
	//req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
}
