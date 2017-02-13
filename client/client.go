package client

import (
	"bytes"
	"errors"
	"github.com/golang/protobuf/proto"
	"log"
	sched "mesos-framework-sdk/include/scheduler"
	"net"
	"net/http"
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
			Transport: &http.Transport{
				Dial: (&net.Dialer{
					Timeout:   10 * time.Second,
					KeepAlive: 30 * time.Second,
				}).Dial,
			},
		},
	}
}

// Makes a new request with data and sends it to the server.
func (c *Client) Request(call *sched.Call) (*http.Response, error) {
	data, err := proto.Marshal(call)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", c.master, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Connection", "keep-alive")
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

	// We will only get the stream ID after a SUBSCRIBE call.
	streamID := resp.Header.Get("Mesos-Stream-Id")
	if streamID != "" {
		c.streamID = streamID
	}

	if resp.StatusCode == http.StatusTemporaryRedirect || resp.StatusCode == http.StatusPermanentRedirect {
		log.Println("Old Master:", c.master)

		master := resp.Header.Get("Location")
		c.master = master

		log.Println("New Master:", c.master)

		return nil, errors.New("Redirect encountered, new master found")
	}

	return resp, nil
}
