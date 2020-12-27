package model

import (
	"encoding/json"
)

type Sentence struct {
	Source        string `json:"source"`
	Target        string `json:"target"`
	RawPredicates string `json:"rawPredicates"`
	Predicates    string `json:"predicates"`
}

type Data struct {
	Edges []Edge `json:"edges"`
	Nodes []Node `json:"nodes"`
	Sum   uint32 `json:"sum"`
}

type Edge struct {
	Source string   `json:"source"`
	Target string   `json:"target"`
	Label  string   `json:"label"`
	Data   EdgeData `json:"data"`
}

type EdgeData struct {
	Properties []KeyValue `json:"properties"`
}

type KeyValue struct {
	Key   string
	Value string
}

type Node struct {
	ID   string   `json:"id"`
	Data NodeData `json:"data"`
}

type NodeData struct {
	ID    string `json:"id"`
	Label string `json:"label"`
}

func ParseInfo(info string) ([]*Sentence, error) {
	var res []*Sentence
	if err := json.Unmarshal([]byte(info), &res); err != nil {
		return nil, err
	}
	return res, nil
}

func QuerySolve(info string) (*Data, error) {
	var resp Data

	req, err := ParseInfo(info)
	if err != nil {
		return nil, err
	}

	resp.Nodes = append(resp.Nodes, Node{ID: req[0].Source})
	for _, item := range req {
		resp.Nodes = append(resp.Nodes, Node{ID: item.Target})
		resp.Edges = append(resp.Edges, Edge{
			Source: item.Source,
			Target: item.Target,
			Label:  item.Predicates,
		})
	}

	return &resp, nil
}
