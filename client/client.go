package client

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"log"
	mesos "mesos-framework-sdk/include/mesos"
	sched "mesos-framework-sdk/include/scheduler"
	"net"
	"net/http"
	"time"
)

const (
	subscribeRetry = 2
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
			Transport: &http.Transport{
				Dial: (&net.Dialer{
					Timeout:   10 * time.Second,
					KeepAlive: 30 * time.Second,
				}).Dial,
			},
		},
	}
}

// TODO We should pass in a Request object since the request headers will be different
// Makes a new request with data and sends it to the server.
func (c *Client) Request(req *http.Request) (*http.Response, error) {
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
	// We will only get the stream ID after a SUBSCRIBE call.
	streamID := resp.Header.Get("Mesos-Stream-Id")
	if streamID != "" {
		c.streamID = streamID
	}

	return resp, nil
}

//
func (c *Client) Subscribe(frameworkInfo *mesos.FrameworkInfo) {
	call := &sched.Call{
		FrameworkId: frameworkInfo.GetId(),
		Type:        sched.Call_SUBSCRIBE.Enum(),
		Subscribe: &sched.Call_Subscribe{
			FrameworkInfo: frameworkInfo,
		},
	}
	// Marshal the scheduler protobuf.
	data, err := proto.Marshal(call)
	if err != nil {
		log.Println(err.Error())
	}
	// Make a new http request from the subscribe call.
	req, err := NewSubscribeRequest(c, data)
	if err != nil {
		log.Println(err.Error())
	}
	// Make the request.
	for {
		resp, err := c.Request(req)
		if err != nil {
			log.Println(err.Error())
		} else {
			// TODO need to spin off from here and handle/decode events
			// Once connected the client should set our framework ID for all outgoing calls after successful subscribe.
			fmt.Println(resp)
			break
		}

		time.Sleep(time.Duration(subscribeRetry) * time.Second)
	}

}

// TODO this should be moved to the scheduler struct when it is made.
func NewSubscribeRequest(c *Client, data []byte) (*http.Request, error) {
	// Make a new subscribe request.
	req, err := http.NewRequest("POST", c.master, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-protobuf")
	// We need to keep the initial request alive for subscribe.
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Accept", "application/x-protobuf")
	req.Header.Set("User-Agent", "mesos-framework-sdk")
	if c.streamID != "" {
		req.Header.Set("Mesos-Stream-Id", c.streamID)
	}
	return req, nil
}
