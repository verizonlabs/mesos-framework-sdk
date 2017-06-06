package client

import (
	"bytes"
	"errors"
	"github.com/golang/protobuf/proto"
	"io/ioutil"
	"mesos-framework-sdk/include/mesos_v1_executor"
	"mesos-framework-sdk/include/mesos_v1_scheduler"
	"mesos-framework-sdk/logging"
	"net"
	"net/http"
	"strings"
	"time"
)

type Client interface {
	Request(interface{}) (*http.Response, error)
	StreamID() string
	SetStreamID(string) Client
}

// HTTP client.
type DefaultClient struct {
	streamID string
	master   string
	client   *http.Client
	logger   logging.Logger
}

// Return a new HTTP client.
func NewClient(master string, logger logging.Logger) Client {
	return &DefaultClient{
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
func (c *DefaultClient) Request(call interface{}) (*http.Response, error) {
	var data []byte
	var err error
	var executorCall bool

	switch call := call.(type) {
	case *mesos_v1_scheduler.Call:
		data, err = proto.Marshal(call)
	case *mesos_v1_executor.Call:
		data, err = proto.Marshal(call)
		executorCall = true
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

	// Executors do not use stream IDs.
	if !executorCall && c.streamID != "" {
		req.Header.Set("Mesos-Stream-Id", c.streamID)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		msg, _ := ioutil.ReadAll(resp.Body)
		return nil, errors.New(string(msg))
	}

	// Our master detection only applies to the scheduler.
	if !executorCall {

		// We will only get the stream ID after a SUBSCRIBE call.
		streamID := resp.Header.Get("Mesos-Stream-Id")
		if streamID != "" {
			c.streamID = streamID
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
	}

	return resp, nil
}

// Gets our stream ID.
func (c *DefaultClient) StreamID() string {
	return c.streamID
}

// Sets our stream ID.
func (c *DefaultClient) SetStreamID(id string) Client {
	c.streamID = id

	return c
}
