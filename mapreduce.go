package riak

import (
	"encoding/json"
)

// An object to build a MapReduce job similar to how the Ruby client can
// build it by adding different stages.
type MapReduce struct {
	client  *Client
	inputs  [][]string
	phases  []string
	request string
}

func (c *Client) MapReduce() *MapReduce {
	return &MapReduce{client: c, inputs: make([][]string, 0), phases: make([]string, 0), request: ""}
}

func (mr *MapReduce) Add(bucket string, key string) {
	mr.inputs = append(mr.inputs, []string{bucket, key})
}

func (mr *MapReduce) LinkBucket(name string, keep bool) {
	link := "{\"link\":{\"bucket\":\"" + name + "\", \"tag\":\"_\",\"keep\":"
	if keep {
		link = link + "true}}"
	} else {
		link = link + "false}}"
	}
	mr.phases = append(mr.phases, link)
}

func (mr *MapReduce) Map(fun string, keep bool) {
	m := "{\"map\":{\"language\":\"javascript\",\"keep\":"
	if keep {
		m = m + "true,"
	} else {
		m = m + "false,"
	}
	m = m + "\"source\":\"" + fun + "\"}}"
	mr.phases = append(mr.phases, m)
}

// Generate the Query string
func (mr *MapReduce) Query() (query []byte, err error) {
	inputs, err := json.Marshal(mr.inputs)
	if err != nil {
		return nil, err
	}
	q := "{\"inputs\":" + string(inputs) + ", \"query\":["
	for i, s := range mr.phases {
		if i > 0 {
			q = q + ","
		}
		q = q + s
	}
	q = q + "]}"
	return []byte(q), nil
}

func (mr *MapReduce) Run() (resp [][]byte, err error) {
	query, err := mr.Query()
	if err != nil {
		return nil, err
	}
	req := &RpbMapRedReq{
		Request:     query,
		ContentType: []byte("application/json"),
	}
	err = mr.client.request(req, "RpbMapRedReq")
	if err != nil {
		return nil, err
	}
	resp, err = mr.client.mr_response()
	if err != nil {
		return nil, err
	}
	return
}

// Run a MapReduce query
func (c *Client) RunMapReduce(query string) (resp [][]byte, err error) {
	req := &RpbMapRedReq{
		Request:     []byte(query),
		ContentType: []byte("application/json"),
	}
	err = c.request(req, "RpbMapRedReq")
	if err != nil {
		return nil, err
	}
	resp, err = c.mr_response()
	if err != nil {
		return nil, err
	}
	return
}
