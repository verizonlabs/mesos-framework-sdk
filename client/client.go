package client

import (
	"bytes"
	"errors"
	"github.com/golang/protobuf/proto"
	"io/ioutil"
	"mesos-framework-sdk/include/executor"
	"mesos-framework-sdk/include/scheduler"
	"mesos-framework-sdk/logging"
	"net"
	"net/http"
	"strings"
	"time"
)

// HTTP client.
type Client struct {
	StreamID string
	master   string
	client   *http.Client
	logger   logging.Logger
}

// Return a new HTTP client.
func NewClient(master string, logger logging.Logger) *Client {
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
		logger: logger,
	}
}

// Makes a new request with data and sends it to the server.
func (c *Client) Request(call interface{}) (*http.Response, error) {
	var data []byte
	var err error

	switch call := call.(type) {
	case *mesos_v1_scheduler.Call:
		data, err = proto.Marshal(call)
	case *mesos_v1_executor.Call:
		data, err = proto.Marshal(call)
	}

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
	req.Header.Set("Accept-Encoding", "gzip")
	req.Header.Set("User-Agent", "mesos-framework-sdk")
	if c.StreamID != "" {
		req.Header.Set("Mesos-Stream-Id", c.StreamID)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		msg, _ := ioutil.ReadAll(resp.Body)
		return nil, errors.New(string(msg))
	}

	// We will only get the stream ID after a SUBSCRIBE call.
	streamID := resp.Header.Get("Mesos-Stream-Id")
	if streamID != "" {
		c.StreamID = streamID
	}

	if resp.StatusCode == http.StatusTemporaryRedirect || resp.StatusCode == http.StatusPermanentRedirect {
		c.logger.Emit(logging.INFO, "Old master: %s", c.master)

		master := resp.Header.Get("Location")
		if strings.Contains(master, "http") {
			c.master = master
		} else {
			c.master = "http:" + master
		}

		c.logger.Emit(logging.INFO, "New master: %s", c.master)

		return nil, errors.New("Redirect encountered, new master found")
	}

	return resp, nil
}
