// Copyright (C) 2016 Great Beyond AB - All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential
// Written by David Högborg <d@greatbeyond.se>, 2016

package mailchimp

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	log "github.com/Sirupsen/logrus"
)

const (
	batchURL = "/batches"
)

// BatchQueue is a collecton of operations that are run with a
// single netowrk call.
// When adding batch commands you should ignore the returned object as it
// will always be a successful placeholder.
// Example:
//  // Clone the client to create a parallel client that will
//  // queue the commands until they are sent to the server.
//  // A clone is prefered to prevent the main client beeing held up
//  // by a batch.
//  client := client.Clone()
//  client.NewBatch()
//
//  // ignore the output, always successful
//  _, _ = client.CreateList(&mailchimp.CreateList{
//      Name: "Batched created list",
//  })
//
//  // run the batch of commands.
//  result, err = client.RunBatch()
//  if err != nil {
//      return nil, err
//  }
//
//  // operations are run in the backgorund
//  println(result.FinishedOperations)
type BatchQueue struct {
	// An array of objects that describes operations to perform.
	Operations []*BatchOperation `json:"operations"`

	// internal
	client *Client
}

type Batch struct {
	// A string that uniquely identifies this batch request.
	ID string `json:"id"`

	// The status of the batch call.
	Status string `json:"status"`

	// The total number of operations to complete as part of this batch request. For GET requests requiring pagination,
	TotalOperations int `json:"total_operations"`

	// The number of completed operations. This includes operations that returned an error.
	FinishedOperations int `json:"finished_operations"`

	// The number of completed operations that returned an error.
	ErroredOperations int `json:"errored_operations"`

	// The date and time when the server received the batch request.
	SubmittedAt string `json:"submitted_at"`

	// The date and time when all operations in the batch request completed.
	CompletedAt string `json:"completed_at"`

	// The URL of the gzipped archive of the results of all the operations.
	ResponseBodyURL string `json:"response_body_url"`
}

// BatchOperation is a single operation part of a batch
type BatchOperation struct {
	// The HTTP method to use for the operation.
	Method string `json:"method,omitempty"`

	// The relative path to use for the operation.
	Path string `json:"path,omitempty"`

	// Any URL params to use, only applies to GET operations.
	Params map[string]string `json:"params,omitempty"`

	// A string containing the JSON body to use with the request.
	Body string `json:"body,omitempty"`

	// An optional client-supplied id returned with the operation results.
	OperationID string `json:"operation_id,omitempty"`
}

// Do adds the operation to the queue
func (b *BatchQueue) Do(request *http.Request) ([]byte, error) {

	var body string
	if request.Body != nil {
		str, _ := ioutil.ReadAll(request.Body)
		body = string(str)
	}

	// pretty.Println(request)
	var params map[string]string
	if request.URL.RawQuery != "" {
		params = b.paramMap(request.URL.RawQuery)
	}

	path := strings.Replace(request.URL.Path, "/3.0/", "/", 1)

	op := &BatchOperation{
		Method: request.Method,
		Path:   path,
		Params: params,
		Body:   body,
	}

	if b.Operations == nil {
		b.Operations = []*BatchOperation{}
	}
	b.Operations = append(b.Operations, op)

	return []byte("{}"), nil
}

func (b *BatchQueue) paramMap(params string) map[string]string {
	splits := strings.Split(params, "&")
	paramMap := map[string]string{}
	for _, part := range splits {
		kv := strings.Split(part, "=")
		if len(kv) == 2 {
			paramMap[kv[0]] = kv[1]
		}
	}
	return paramMap
}

// Run executes a batch and resets the batch queue
func (b *BatchQueue) Run() (*Batch, error) {

	response, err := b.client.post(batchURL, nil, b)
	if err != nil {
		log.Debug(err.Error, caller())
		return nil, err
	}

	var br *Batch
	err = json.Unmarshal(response, &br)
	if err != nil {
		log.Debug(err.Error, caller())
		return nil, err
	}
	return br, nil
}

// Get retrieves a single batch status object
func (b *BatchQueue) Get(id string) (*BatchOperation, error) {
	response, err := b.client.get(slashJoin(batchURL, id), nil)
	if err != nil {
		log.WithFields(log.Fields{
			"batch_id": id,
			"error":    err.Error(),
		}).Warn("response error", caller())
		return nil, err
	}

	var batch *BatchOperation
	err = json.Unmarshal(response, &batch)
	if err != nil {
		log.WithFields(log.Fields{
			"batch_id": id,
			"error":    err.Error(),
		}).Warn("response error", caller())
		return nil, err
	}

	return batch, nil
}
