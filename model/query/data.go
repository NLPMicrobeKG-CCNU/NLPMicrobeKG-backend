package query

import (
	"encoding/json"
	"fmt"

	"github.com/NLPMicrobeKG-CCNU/NLPMicrobeKG-backend/service/graphDB"
)

type DataQueryResponse struct {
	Head    DataHead    `json:"head"`
	Results DataResults `json:"results"`
}

type DataHead struct {
	Vars []string `json:"vars"`
}

type Compoundname struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type Unit struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type Foodname struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type Modulename struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type Mount struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type DataBindings struct {
	Compoundname Compoundname `json:"compoundname"`
	Unit         Unit         `json:"unit"`
	Foodname     Foodname     `json:"foodname"`
	Modulename   Modulename   `json:"modulename"`
	Mount        Mount        `json:"mount"`
}

type DataResults struct {
	Bindings []DataBindings `json:"bindings"`
}

type DataResponse struct {
	Compoundname string `json:"compoundname"`
	Unit         string `json:"unit"`
	Foodname     string `json:"foodname"`
	Modulename   string `json:"modulename"`
	Mount        string `json:"mount"`
	Bacname      string `json:"bacname"`
}

// DataQuery returns results of data query.
func DataQuery(query, bacname string, limit, offset int) ([]*DataResponse, error) {
	res, err := GetDataQueryRes(query, limit, offset)
	if err != nil {
		return []*DataResponse{}, err
	}

	return TransformToData(res, bacname)
}

// GetTextQueryRes returns raw graphdb query response in type of data.
func GetDataQueryRes(query string, limit, offset int) (*DataQueryResponse, error) {
	var res *DataQueryResponse

	raw, err := graphDB.QueryInfo(query, limit, offset)
	fmt.Println(string(raw))
	if err != nil {
		return res, err
	}

	err = json.Unmarshal(raw, &res)
	if err != nil {
		return res, err
	}

	return res, nil
}

// TransformToData transform graphdb query response into service data response.
func TransformToData(req *DataQueryResponse, bacname string) ([]*DataResponse, error) {
	var resp []*DataResponse
	for _, item := range req.Results.Bindings {
		resp = append(resp, &DataResponse{
			Bacname:      bacname,
			Compoundname: item.Compoundname.Value,
			Unit:         item.Unit.Value,
			Foodname:     item.Foodname.Value,
			Modulename:   item.Modulename.Value,
			Mount:        item.Mount.Value,
		})
	}

	return resp, nil
}
