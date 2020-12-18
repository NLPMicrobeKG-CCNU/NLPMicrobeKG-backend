package query

import (
	"encoding/json"
	"github.com/NLPMicrobeKG-CCNU/NLPMicrobeKG-backend/service"
)

type TextQueryResponse struct {
	Head    TextHead    `json:"head"`
	Results TextResults `json:"results"`
}

type TextHead struct {
	Vars []string `json:"vars"`
}

type Reference struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type Bac4Name struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type Ref4 struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type Ref3 struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type Bac2Name struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type Ref2 struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type Bac3Name struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type Disname struct {
	XMLLang string `json:"xml:lang"`
	Type    string `json:"type"`
	Value   string `json:"value"`
}

type TextBindings struct {
	Reference Reference `json:"reference"`
	Bac4Name  Bac4Name  `json:"bac4name"`
	Ref4      Ref4      `json:"ref4"`
	Ref3      Ref3      `json:"ref3"`
	Bac2Name  Bac2Name  `json:"bac2name"`
	Ref2      Ref2      `json:"ref2"`
	Bac3Name  Bac3Name  `json:"bac3name"`
	Disname   Disname   `json:"disname"`
}

type TextResults struct {
	Bindings []TextBindings `json:"bindings"`
}

type TextResponse struct {
	Reference string `json:"reference"`
	Bac4Name  string `json:"bac4name"`
	Ref4      string `json:"ref4"`
	Ref3      string `json:"ref3"`
	Bac2Name  string `json:"bac2name"`
	Ref2      string `json:"ref2"`
	Bac3Name  string `json:"bac3name"`
	Disname   string `json:"disname"`
}

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
	Source       string `json:"source"`
	Target       string `json:"target"`
	Relationship string `json:"label"`
}

type Node struct {
	ID string `json:"id"`
}

func TextQueryInfo(query string, limit, offset int) ([]*TextQueryResponse, error) {
	var res []*TextQueryResponse

	raw, err := service.QueryInfo(query, limit, offset)
	if err != nil {
		return res, err
	}

	err = json.Unmarshal(raw, &res)
	if err != nil {
		return res, err
	}

	return res, nil
}

func TransformInText(req []*TextQueryResponse) ([]*TextResponse, error) {
	var resp []*TextResponse
	list := req[0]
	for _, item := range list.Results.Bindings{
		resp = append(resp, &TextResponse{
			Reference: item.Reference.Value,
			Bac4Name:  item.Bac4Name.Value,
			Ref4:      item.Ref4.Value,
			Ref3:      item.Ref3.Value,
			Bac2Name:  item.Bac2Name.Value,
			Ref2:      item.Ref2.Value,
			Bac3Name:  item.Bac3Name.Value,
			Disname:   item.Disname.Value,
		})
	}
	return resp, nil
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
	count := 1
	for _, item := range req {
		resp.Nodes = append(resp.Nodes, Node{ID: item.Target})
		count++
		resp.Edges = append(resp.Edges, Edge{
			Source:       item.Source,
			Target:       item.Target,
			Relationship: item.Predicates,
		})
	}
	resp.Sum = uint32(count)

	return &resp, nil
}
