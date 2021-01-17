package query

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/NLPMicrobeKG-CCNU/NLPMicrobeKG-backend/service/graphDB"
	"github.com/NLPMicrobeKG-CCNU/NLPMicrobeKG-backend/util"
)

// solve graph query

// Search results.
type Data struct {
	Edges []Edge `json:"edges"`
	Nodes []Node `json:"nodes"`
	Sum   uint32 `json:"sum"`
}

// Sentence is raw data from graphdb.
type Sentence struct {
	Source        string   `json:"source"`
	Target        string   `json:"target"`
	RawPredicates []string `json:"rawPredicates"`
	Predicates    []string `json:"predicates"`
}

// Edge.
type Edge struct {
	Source       string   `json:"source"`
	Target       string   `json:"target"`
	Relationship string   `json:"label"`
	Data         struct{} `json:"data"` // empty realize
}

// NodeProperty is raw node property.
type NodeProperty struct {
	Type  string `json:"t"`
	Value string `json:"v"`
}

// Node properties
type NodeProperties struct {
	Title string   `json:"title"`
	Value []string `json:"value"`
}

// Node
type Node struct {
	ID         string           `json:"id"`
	Labels     string           `json:"label"`
	Types      []string         `json:"type"`
	RDFRank    float64          `json:"rank"`
	Properties []NodeProperties `json:"node_properties"`
	Color      int              `json:"color"`
	Data       struct{}         `json:"data"` // empty realize
}

// Node details from graph db
type RawNodeDetails struct {
	Iri     string      `json:"iri"`
	Labels  []Labels    `json:"labels"`
	Comment interface{} `json:"comment"`
	Types   []string    `json:"types"`
	Size    int         `json:"size"`
	RdfRank float64     `json:"rdfRank"`
}

// node labels
type Labels struct {
	Lang     string `json:"lang"`
	Priority int    `json:"priority"`
	Label    string `json:"label"`
}

// ParseGraphInfo parses information to struct.
func ParseGraphInfo(data []byte) ([]*Sentence, error) {
	var res []*Sentence
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, err
	}

	return res, nil
}

// ParseNodeProperties parses node properties.
func ParseNodeProperties(data []byte) ([]NodeProperties, error) {
	var res []NodeProperties
	var transfer map[string][]NodeProperty
	if err := json.Unmarshal(data, &transfer); err != nil {
		fmt.Println("err!", err)
		return nil, err
	}

	for k, v := range transfer {
		var title string
		var value []string
		if util.FormatNodeProperty(v[0].Type) {
			title = "rdfs:label"
		} else {
			title = k
		}
		for _, s := range v {
			value = append(value, s.Value)
		}
		res = append(res, NodeProperties{
			Title: title,
			Value: value,
		})
	}

	return res, nil
}

// ParseNodeDetails parses node details.
func ParseNodeDetails(data []byte) (RawNodeDetails, error) {
	var res RawNodeDetails
	if err := json.Unmarshal(data, &res); err != nil {
		return res, err
	}

	return res, nil
}

// GraphQuery return query results from visual graph.
func GraphQuery(query string) (*Data, error) {
	var resp Data

	rawResponse, err := graphDB.GetVisualGraphInfo(query)
	if err != nil {
		return &resp, err
	}

	req, err := ParseGraphInfo(rawResponse)
	if err != nil {
		return nil, err
	}

	//resp.Nodes = append(resp.Nodes, Node{
	//	ID: req[0].Source,
	//})

	count := 1

	var typeCount int
	var nodes []string
	typeRecord := make(map[string]int)
	nodeRecord := make(map[string]bool)
	for _, item := range req {
		count++
		resp.Edges = append(resp.Edges, Edge{
			Source:       item.Source,
			Target:       item.Target,
			Relationship: strings.Join(item.Predicates, ", "),
		})

		if !nodeRecord[item.Target] {
			nodes = append(nodes, item.Target)
			nodeRecord[item.Target] = true
		}
		if !nodeRecord[item.Source] {
			nodes = append(nodes, item.Source)
			nodeRecord[item.Source] = true
		}
	}

	for _, node := range nodes {
		var v Node
		v.ID = node

		nodeProperties, err := graphDB.GetNodeProperties(v.ID)
		if err != nil {
			fmt.Println(err)
			return &resp, err
		}
		v.Properties, err = ParseNodeProperties(nodeProperties)

		rawNodeDetails, err := graphDB.GetNodeDetails(v.ID)
		if err != nil {
			fmt.Println(err)
			return &resp, err
		}
		nodeDetails, err := ParseNodeDetails(rawNodeDetails)
		if err != nil {
			fmt.Println(err)
			return &resp, err
		}

		v.Types = util.FormatNodeTypeStr(nodeDetails.Types)

		typeStr := strings.Join(v.Types, ",")
		color := typeRecord[typeStr]
		if typeRecord[typeStr] == 0 {
			typeCount++
			color = typeCount
			typeRecord[typeStr] = color
		}

		v.RDFRank = nodeDetails.RdfRank

		var labels []string
		for _, label := range nodeDetails.Labels {
			labels = append(labels, label.Label)
		}
		v.Labels = strings.Join(labels, ",")

		v.Color = color

		resp.Nodes = append(resp.Nodes, v)
	}

	resp.Sum = uint32(count)

	return &resp, nil
}
