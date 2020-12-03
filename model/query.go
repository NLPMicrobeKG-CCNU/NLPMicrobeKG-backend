package model

import (
	"encoding/json"
	"github.com/NLPMicrobeKG-CCNU/NLPMicrobeKG-backend/service"
)

type QueryResponse struct {
	Info string `json:"info"`
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
