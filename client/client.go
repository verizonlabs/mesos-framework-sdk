// Copyright 2017 Verizon
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package client

import (
	"bytes"
	"errors"
	"io/ioutil"
	"github.com/carlonelong/mesos-framework-sdk/include/mesos/v1/executor"
	"github.com/carlonelong/mesos-framework-sdk/include/mesos/v1/scheduler"
	"github.com/carlonelong/mesos-framework-sdk/logging"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/golang/protobuf/proto"
)

type Client interface {
	Request(interface{}) (*http.Response, error)
	StreamID() string
	SetStreamID(string) Client
}

type ClientData struct {
	Endpoint string
	Auth     string
}

// HTTP client.
type DefaultClient struct {
	streamID string
	data     ClientData
	client   *http.Client
	logger   logging.Logger
}

// Return a new HTTP client.
func NewClient(data ClientData, logger logging.Logger) Client {
	return &DefaultClient{
		data: data,
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
// Determines whether the request/response should be handled for an executor or a scheduler.
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

	req, err := http.NewRequest("POST", c.data.Endpoint, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", c.data.Auth)
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
		if resp.StatusCode == 401 {
			return resp, errors.New("Unauthorized")
		}

		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return resp, err
		}

		return resp, errors.New(string(data))
	}

	// Our master detection only applies to the scheduler.
	if !executorCall {

		// We will only get the stream ID after a SUBSCRIBE call.
		streamID := resp.Header.Get("Mesos-Stream-Id")
		if streamID != "" {
			c.streamID = streamID
		}

		if resp.StatusCode == http.StatusTemporaryRedirect || resp.StatusCode == http.StatusPermanentRedirect {
			c.logger.Emit(logging.INFO, "Old master: %s", c.data.Endpoint)

			master := resp.Header.Get("Location")
			if strings.Contains(master, "http") {
				c.data.Endpoint = master
			} else {
				c.data.Endpoint = resp.Request.URL.Scheme + ":" + master
			}

			c.logger.Emit(logging.INFO, "New master: %s", c.data.Endpoint)

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
