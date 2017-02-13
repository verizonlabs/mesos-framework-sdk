package client

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"log"
	mesos "mesos-framework-sdk/include/mesos"
	sched "mesos-framework-sdk/include/scheduler"
	"mesos-sdk/recordio"
	"net"
	"net/http"
	"time"
)

const (
	subscribeRetry = 2
)

// HTTP client.
type Client struct {
	streamID    string
	master      string
	client      *http.Client
	frameworkId mesos.FrameworkID
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

// Create a Subscription to mesos.
func (c *Client) Subscribe(frameworkInfo *mesos.FrameworkInfo) {
	// We really want the ID after the call...
	c.frameworkId = *frameworkInfo.GetId()
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
			_ = recordio.NewReader(resp.Body)
			resp.Body.Close()
			break
		}

		time.Sleep(time.Duration(subscribeRetry) * time.Second)
	}
}

// Send a teardown request to mesos master.
func (c *Client) Teardown() {
	if *c.frameworkId.Value != "" {
		teardown := &sched.Call{
			FrameworkId: &c.frameworkId,
			Type:        sched.Call_TEARDOWN.Enum(),
		}
		resp, err := c.DefaultPostRequest(teardown)
		if err != nil {
			log.Println(err.Error())
		}
		fmt.Println(resp)
		return
	}
	fmt.Print("No framework id: ")
	fmt.Println(c.frameworkId.Value)
}

// Skeleton funcs for the rest of the calls.

// Accepts offers from mesos master
func (c *Client) Accept(offerIds []mesos.OfferID, tasks []mesos.Task, filters mesos.Filters) {
	accept := &sched.Call{
		FrameworkId: &c.frameworkId,
		Type:        sched.Call_ACCEPT.Enum(),
		Accept:      sched.Call_Accept{OfferIds: offerIds, Operations: tasks, Filters: filters},
	}

	resp, err := c.DefaultPostRequest(accept)
	if err != nil {
		log.Println(err.Error())
	}
	fmt.Println(resp)
}

func (c *Client) Decline(offerIds []mesos.OfferID, filters mesos.Filters) {
	// Get a list of the offer ids to decline and any filters.
	decline := &sched.Call{
		FrameworkId: &c.frameworkId,
		Type:        sched.Call_DECLINE.Enum(),
		Decline:     sched.Call_Decline{OfferIds: offerIds, Filters: filters},
	}

	resp, err := c.DefaultPostRequest(decline)
	if err != nil {
		log.Println(err.Error())
	}
	fmt.Println(resp)
	return
}

// Sent by the scheduler to remove any/all filters that it has previously set via ACCEPT or DECLINE calls.
func (c *Client) Revive() {

	revive := &sched.Call{
		FrameworkId: &c.frameworkId,
		Type:        sched.Call_REVIVE.Enum(),
	}

	resp, err := c.DefaultPostRequest(revive)
	if err != nil {
		log.Println(err.Error())
	}
	fmt.Println(resp)
	return
}

func (c *Client) Kill(taskId mesos.TaskID, agentid mesos.AgentID) {
	// Probably want some validation that this is a valid task and valid agentid.
	kill := &sched.Call{
		FrameworkId: &c.frameworkId,
		Type:        sched.Call_KILL.Enum(),
		Kill:        sched.Call_Kill{TaskId: taskId, AgentId: agentid},
	}

	resp, err := c.DefaultPostRequest(kill)
	if err != nil {
		log.Println(err.Error())
	}
	fmt.Println(resp)
	return
}

func (c *Client) Shutdown(execId mesos.ExecutorID, agentId mesos.AgentID) {
	shutdown := &sched.Call{
		FrameworkId: &c.frameworkId,
		Type:        sched.Call_SHUTDOWN.Enum(),
		Shutdown: sched.Call_Shutdown{
			ExecutorId: execId,
			AgentId:    agentId,
		},
	}
	resp, err := c.DefaultPostRequest(shutdown)
	if err != nil {
		log.Println(err.Error())
	}
	fmt.Println(resp)
	return
}

func (c *Client) Acknowledge() {

}

func (c *Client) Reconcile() {

}

func (c *Client) Message() {

}

func (c *Client) SchedRequest() {

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
	req, err := http.NewRequest("POST", c.master, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "mesos-framework-sdk")
	return req, nil
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
