package client

import (
	"bytes"
	"errors"
	"github.com/golang/protobuf/proto"
	"io/ioutil"
	"log"
	mesos "mesos-framework-sdk/include/mesos"
	sched "mesos-framework-sdk/include/scheduler"
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

func (c *Client) Subscribe(frameworkInfo *mesos.FrameworkInfo) {
	call := &sched.Call{
		FrameworkId: frameworkInfo.Id,
		Type:        sched.Call_SUBSCRIBE.Enum(),
		Subscribe: &sched.Call_Subscribe{
			FrameworkInfo: frameworkInfo,
		},
	}
	data, err := proto.Marshal(call)
	if err != nil {
		log.Println(err.Error())
	}

	for {
		resp, err := c.Request(data)
		if err != nil {
			log.Println(err.Error())
		} else {
			// TODO need to spin off from here and handle/decode events
			break
		}

		time.Sleep(time.Duration(subscribeRetry) * time.Second)
	}

}
