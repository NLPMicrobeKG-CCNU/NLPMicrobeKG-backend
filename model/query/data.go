package query

import (
	"encoding/json"
	"github.com/NLPMicrobeKG-CCNU/NLPMicrobeKG-backend/service"
)

type DataQueryResponse struct {
	Head DataHead `json:"head"`
	Results DataResults `json:"results"`
}

type DataHead struct {
	Vars []string `json:"vars"`
}

type Compoundname struct {
	Type string `json:"type"`
	Value string `json:"value"`
}

type Unit struct {
	Type string `json:"type"`
	Value string `json:"value"`
}

type Foodname struct {
	Type string `json:"type"`
	Value string `json:"value"`
}

type Modulename struct {
	Type string `json:"type"`
	Value string `json:"value"`
}

type Mount struct {
	Type string `json:"type"`
	Value string `json:"value"`
}

type DataBindings struct {
	Compoundname Compoundname `json:"compoundname"`
	Unit Unit `json:"unit"`
	Foodname Foodname `json:"foodname"`
	Modulename Modulename `json:"modulename"`
	Mount Mount `json:"mount"`
}

type DataResults struct {
	Bindings []DataBindings `json:"bindings"`
}

type DataResponse struct {
	Compoundname string `json:"compoundname"`
	Unit string `json:"unit"`
	Foodname string `json:"foodname"`
	Modulename string `json:"modulename"`
	Mount string `json:"mount"`
}

func DataQueryInfo(query string, limit, offset int) ([]*DataQueryResponse, error) {
	var res []*DataQueryResponse

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


func TransformInData(req []*DataQueryResponse) ([]*DataResponse, error) {
	var resp []*DataResponse
	list := req[0]
	for _, item := range list.Results.Bindings {
		resp = append(resp, &DataResponse{
			Compoundname: item.Compoundname.Value,
			Unit:         item.Unit.Value,
			Foodname:     item.Foodname.Value,
			Modulename:   item.Modulename.Value,
			Mount:        item.Mount.Value,
		})
	}
	return resp, nil
}