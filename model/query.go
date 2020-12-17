package model

import (
	"encoding/json"
	"github.com/NLPMicrobeKG-CCNU/NLPMicrobeKG-backend/service"
)

type QueryResponse struct {
	Info string `json:"info"`
}

type Sentence struct {
	Source        string `json:"source"`
	Target        string `json:"target"`
	RawPredicates string `json:"rawPredicates"`
	Predicates    string `json:"predicates"`
}

func QueryInfo(query string, limit, offset int) ([]*QueryResponse, error) {
	var res []*QueryResponse
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

func ParseInfo(info string) ([]*Sentence, error) {
	var res []*Sentence
	if  err := json.Unmarshal([]byte(info), &res); err != nil {
		return nil, err
	}
	return res, nil
}
