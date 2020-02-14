// Package flow is a sub package to workflow.
// Flow package contains interfaces related to and used by workflows
// Its the Flow struct that handles logging,statistics, commiunication between flows
// and the Payload interface, which handles data for each item transmitted
// T
package flow

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"
)

var (
	//DefaultBufferSize is used to set a default buffer to channels
	DefaultBufferSize int = 1000
	//ErrNoIngressChannel is a error that occurs when the flow does not have a ingress
	ErrNoIngressChannel error = errors.New("There is no ingressChannel configured in the flow")
)

// NewFlow is used to correctly initialize a new Flow with all values needed
// Use this instead of creating flows manually to avoid nil pointers etc
func NewFlow(name string, ingress chan Payload, conf json.RawMessage) *Flow {
	return &Flow{
		ProcessorName:  name,
		ingressChannel: ingress,
		Configuration:  conf,
		ErrorChannel:   make(chan error, DefaultBufferSize),
	}
}

// Flow is used to create a flow correct struct that can init a workflow process
type Flow struct {
	ingressChannel chan Payload
	egressChannel  chan Payload
	ErrorChannel   chan error `json:"-"`
	//StopChannel is a channel that should be used by Spawned Goroutines to know when to exit
	StopChannel   chan bool       `json:"-"`
	ProcessorName string          `json:"processor"`
	Configuration json.RawMessage `json:"configuration"`
	wg            *sync.WaitGroup
}

// Close will close all channels that are used in the Flow
// All spawned G
func (nf *Flow) Close() {
	nf.StopChannel <- true
	close(nf.ErrorChannel)
	close(nf.ingressChannel)
	close(nf.egressChannel)
	close(nf.StopChannel)
}

// SetWaitGroup will change the current waitgroup
func (nf *Flow) SetWaitGroup(wg *sync.WaitGroup) {
	nf.wg = wg
}

// GetWaitGroup will return the configured waitgroup to be used by processors that has gorotuines
func (nf *Flow) GetWaitGroup() *sync.WaitGroup {
	return nf.wg
}

//GetIngressChannel is used by processors that require a continous flow of new flows,
//It should return a channel that will keep returning Flows for the duration of the Workflow duration
func (nf *Flow) GetIngressChannel() chan Payload {
	return nf.ingressChannel
}

//GetEgressChannel will return a channel that reports Outgoing Flows from a Flow
func (nf *Flow) GetEgressChannel() chan Payload {
	return nf.egressChannel
}

// SetIngressChannel is used to set a new Channel for ingressing flows, This hsould be the previous channels Egress Channel
// The ingressChannel should commonly be set by the previous Flow executed
// and should be the previous flows EgressChannel
func (nf *Flow) SetIngressChannel(newchan chan Payload) {
	nf.ingressChannel = newchan
}

//SetEgressChannel will change the egress channel into a new one
func (nf *Flow) SetEgressChannel(egress chan Payload) {
	nf.egressChannel = egress
}

//GetType will retutrn the configured type, Type should be the processor name
func (nf *Flow) GetType() string {
	return nf.ProcessorName
}

//SetType is used to change the value of a type
func (nf *Flow) SetType(s string) {
	nf.ProcessorName = s
}

//GetConfiguration will return a raw JSON to be Unmarshalled into propriate struct
func (nf *Flow) GetConfiguration() json.RawMessage {
	return nf.Configuration
}

//SetConfiguration is a way to change the Configs
func (nf *Flow) SetConfiguration(conf json.RawMessage) {
	nf.Configuration = conf
}

//Log should store the error into the configured Logging mechanism
//Should be changed from single value to Channel I guess.
func (nf *Flow) Log(err error) {
	// only for debugging remove this fmt.Pritln before publishing
	fmt.Println(err)
	nf.ErrorChannel <- err
}
