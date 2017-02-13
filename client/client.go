package client

import (
	"bytes"
	"errors"
	"github.com/golang/protobuf/proto"
	"log"
	mesos "mesos-framework-sdk/include/mesos"
	sched "mesos-framework-sdk/include/scheduler"
	"net"
	"net/http"
	"time"
)

// HTTP client.
type Client struct {
	StreamID    string
	Master      string
	Client      *http.Client
	FrameworkId mesos.FrameworkID
}

// Return a new HTTP client.
func NewClient(master string) *Client {
	return &Client{
		Master: master,
		Client: &http.Client{
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
	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusTemporaryRedirect || resp.StatusCode == http.StatusPermanentRedirect {
		log.Println("Old Master:", c.Master)

		master := resp.Header.Get("Location")
		c.Master = master

		log.Println("New Master:", c.Master)

		return nil, errors.New("Redirect encountered, new master found")
	}
	// We will only get the stream ID after a SUBSCRIBE call.
	streamID := resp.Header.Get("Mesos-Stream-Id")
	if streamID != "" {
		c.StreamID = streamID
	}

	return resp, nil
}

// Func that marshals the call, wraps it up in a http.request and sends it off.
func (c *Client) DefaultPostRequest(call *sched.Call) (*http.Response, error) {
	data, err := proto.Marshal(call)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	req, err := NewDefaultPostHeaders(c, data)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	resp, err := c.Request(req)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return resp, nil
}

// Default headers to set for a post request for mesos.
func NewDefaultPostHeaders(c *Client, data []byte) (*http.Request, error) {
	req, err := http.NewRequest("POST", c.Master, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-protobuf")
	req.Header.Set("User-Agent", "mesos-framework-sdk")
	return req, nil
}

// TODO this should be moved to the scheduler struct when it is made.
func NewSubscribeRequest(c *Client, data []byte) (*http.Request, error) {
	// Make a new subscribe request.
	req, err := http.NewRequest("POST", c.Master, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-protobuf")
	// We need to keep the initial request alive for subscribe.
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Accept", "application/x-protobuf")
	req.Header.Set("User-Agent", "mesos-framework-sdk")
	if c.StreamID != "" {
		req.Header.Set("Mesos-Stream-Id", c.StreamID)
	}
	return req, nil
}
