package client

import (
	sched "mesos-framework-sdk/include/scheduler"
	"net/http"
	//mesos "mesos-framework-sdk/include/mesos"
	"bytes"
	"fmt"
	"github.com/gogo/protobuf/proto"
	"io/ioutil"
	"log"
)

func Subscribe_Call(call *sched.Call) {
	client := &http.Client{}
	k, err := proto.Marshal(call)
	if err != nil {
		log.Println(err.Error())
	}
	req, err := http.NewRequest("POST", "http://10.0.0.10:5050/api/v1/scheduler", bytes.NewBuffer(k))
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Print("ERROR: ")
		log.Println(err.Error())
	}
	fmt.Println("response Body:", string(body))

}
