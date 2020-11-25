// Package files is generated by Handlergenerator tooling
// Make sure to insert real Description here
package files

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/percybolmer/workflow/metric"
	"github.com/percybolmer/workflow/payload"
	"github.com/percybolmer/workflow/property"
	"github.com/percybolmer/workflow/pubsub"
	"github.com/percybolmer/workflow/register"
)

// ReadFile is used to ReadFiles data
type ReadFile struct {
	// Cfg is values needed to properly run the Handle func
	Cfg    *property.Configuration `json:"configs" yaml:"configs"`
	Name   string                  `json:"handler" yaml:"handler_name"`
	remove bool

	subscriptionless bool
	errChan          chan error

	metrics      metric.Provider
	metricPrefix string
	// MetricPayloadOut is how many payloads the processor has outputted
	MetricPayloadOut string
	// MetricPayloadIn is how many payloads the processor has inputted
	MetricPayloadIn string
}

func init() {
	register.Register("ReadFile", NewReadFileHandler())
}

// NewReadFileHandler generates a new ReadFile Handler
func NewReadFileHandler() *ReadFile {
	act := &ReadFile{
		Cfg: &property.Configuration{
			Properties: make([]*property.Property, 0),
		},
		Name:    "ReadFile",
		errChan: make(chan error, 1000),
	}
	act.Cfg.AddProperty("remove_after", "This property is used to configure if files that are read should be removed after", true)
	return act
}

// GetHandlerName is used to retrun a unqiue string name
func (a *ReadFile) GetHandlerName() string {
	return a.Name
}

// Handle is used to Read the content of a file from the former payload
// Expects a filepath in the input payload
func (a *ReadFile) Handle(ctx context.Context, input payload.Payload, topics ...string) error {
	a.metrics.IncrementMetric(a.MetricPayloadIn, 1)
	path := string(input.GetPayload())
	file, err := os.Open(path)

	if err != nil {
		return err
	}
	defer func() {
		file.Close()
		if a.remove {
			os.Remove(path)
		}
	}()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	a.metrics.IncrementMetric(a.MetricPayloadOut, 1)
	errs := pubsub.PublishTopics(topics, payload.BasePayload{
		Payload: data,
		Source:  "ReadFile",
	})
	for _, err := range errs {
		a.errChan <- err
	}

	return nil
}

// ValidateConfiguration is used to see that all needed configurations are assigned before starting
func (a *ReadFile) ValidateConfiguration() (bool, []string) {
	// Check if Cfgs are there as needed
	removeProp := a.Cfg.GetProperty("remove_after")
	missing := make([]string, 0)

	if removeProp == nil && removeProp.Value == nil {
		missing = append(missing, "remove_after")
		return false, missing
	}

	remove, err := removeProp.Bool()
	if err != nil {
		return false, nil
	}

	a.remove = remove

	return true, nil
}

// GetConfiguration will return the CFG for the Handler
func (a *ReadFile) GetConfiguration() *property.Configuration {
	return a.Cfg
}

// Subscriptionless will return true/false if the Handler is genereating payloads itself
func (a *ReadFile) Subscriptionless() bool {
	return a.subscriptionless
}

// GetErrorChannel will return a channel that the Handler can output eventual errors onto
func (a *ReadFile) GetErrorChannel() chan error {
	return a.errChan
}

// SetMetricProvider is used to change what metrics provider is used by the handler
func (a *ReadFile) SetMetricProvider(p metric.Provider, prefix string) error {
	a.metrics = p
	a.metricPrefix = prefix

	a.MetricPayloadIn = fmt.Sprintf("%s_payloads_in", prefix)
	a.MetricPayloadOut = fmt.Sprintf("%s_payloads_out", prefix)
	err := a.metrics.AddMetric(&metric.Metric{
		Name:        a.MetricPayloadOut,
		Description: "keeps track of how many payloads the handler has outputted",
	})
	if err != nil {
		return err
	}
	err = a.metrics.AddMetric(&metric.Metric{
		Name:        a.MetricPayloadIn,
		Description: "keeps track of how many payloads the handler has ingested",
	})
	return err
}
