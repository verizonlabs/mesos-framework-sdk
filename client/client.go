package client

import (
	sched "mesos-framework-sdk/include/scheduler"
	"net/http"
	//mesos "mesos-framework-sdk/include/mesos"
	"bytes"
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"io/ioutil"
	"log"
	"time"
)

// HTTP client.
type Client struct {
	streamID string
	master   string
	client   *http.Client
}

// Return a new HTTP client.
func NewClient(master string) *Client {
	return &Client{
		master: master,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Makes a new request with data and sends it to the server.
func (c *Client) Request(data []byte) (*http.Response, error) {
	req, err := http.NewRequest("POST", c.master, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-protobuf")
	req.Header.Set("Accept", "application/x-protobuf")
	req.Header.Set("User-Agent", "mesos-framework-sdk")
	if c.streamID != "" {
		req.Header.Set("Mesos-Stream-Id", c.streamID)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusTemporaryRedirect || resp.StatusCode == http.StatusPermanentRedirect {
		log.Println("Old Master:", c.master)

		master := resp.Header.Get("Location")
		c.master = master

		log.Println("New Master:", c.master)

		return nil, errors.New("Redirect encountered, new master found")
	}

	streamID := resp.Header.Get("Mesos-Stream-Id")
	if streamID != "" {
		c.streamID = streamID
	}

	return resp, nil
}

func Subscribe_Call(call *sched.Call) {
	client := &http.Client{}
	k, err := proto.Marshal(call)
	if err != nil {
		log.Println(err.Error())
	}
	req, err := http.NewRequest("POST", "http://10.0.0.10:5050/api/v1/scheduler", bytes.NewBuffer(k))
	req.Header.Set("Content-Type", "application/x-protobuf")

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
