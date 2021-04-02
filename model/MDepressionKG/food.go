package MDepressionKG

import (
	"encoding/json"
	"fmt"
	"github.com/NLPMicrobeKG-CCNU/NLPMicrobeKG-backend/log"
	"github.com/NLPMicrobeKG-CCNU/NLPMicrobeKG-backend/service/graphDB"
	"go.uber.org/zap"
)

type NFoodQueryResponse struct {
	Head    FoodHead     `json:"head"`
	Results NFoodResults `json:"results"`
}

type PFoodQueryResponse struct {
	Head    FoodHead     `json:"head"`
	Results PFoodResults `json:"results"`
}

type FoodHead struct {
	Vars []string `json:"vars"`
}

type FoodBacname struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}
type Compoundname struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}
type FoodBac struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}
type Compound struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}
type Food struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}
type NFoodBindings struct {
	Bacname      FoodBacname  `json:"bacname"`
	Compoundname Compoundname `json:"compoundname"`
	Bac          FoodBac      `json:"bac"`
	Depression   Depression   `json:"depression"`
	Compound     Compound     `json:"compound"`
	Food         Food         `json:"food"`
}
type NFoodResults struct {
	Bindings []NFoodBindings `json:"bindings"`
}

type Depression struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type PFoodBindings struct {
	Compoundname Compoundname `json:"compoundname"`
	Depression   Depression   `json:"depression"`
	Compound     Compound     `json:"compound"`
	Food         Food         `json:"food"`
}
type PFoodResults struct {
	Bindings []PFoodBindings `json:"bindings"`
}

type FoodResponse struct {
	Bacname      string  `json:"bacname"`
	Bac          string      `json:"bac"`
	CompoundName string `json:"compoundname"`
	Compound     string     `json:"compound"`
	Food         string       `json:"food"`
	Type         string       `json:"type"`
}

// FoodQuery returns results of food query.
func FoodQuery(pQuery, nQuery, foodName string, limit, offset int) ([]*FoodResponse, error) {
	pRes, err := GetPFoodQueryRes(pQuery, limit, offset)
	if err != nil {
		log.Error("food query(P) error", zap.String("cause", err.Error()))
		return []*FoodResponse{}, err
	}
	nRes, err := GetNFoodQueryRes(nQuery, limit, offset)
	if err != nil {
		log.Error("food query(N) error", zap.String("cause", err.Error()))
		return []*FoodResponse{}, err
	}

	return TransformToFood(pRes, nRes, foodName)
}

// GetNFoodQueryRes returns raw graphdb query response in type of food on negative.
func GetNFoodQueryRes(query string, limit, offset int) (*NFoodQueryResponse, error) {
	var res *NFoodQueryResponse

	raw, err := graphDB.QueryInfo(query, limit, offset)
	if err != nil {
		return res, err
	}

	err = json.Unmarshal(raw, &res)
	if err != nil {
		return res, err
	}

	return res, nil
}

// GetPFoodQueryRes returns raw graphdb query response in type of food on positive.
func GetPFoodQueryRes(query string, limit, offset int) (*PFoodQueryResponse, error) {
	var res *PFoodQueryResponse

	raw, err := graphDB.QueryInfo(query, limit, offset)
	if err != nil {
		return res, err
	}

	err = json.Unmarshal(raw, &res)
	if err != nil {
		return res, err
	}

	return res, nil
}

// TransformToData transform graphdb query response into service disease response.
func TransformToFood(pReq *PFoodQueryResponse, nReq *NFoodQueryResponse, foodName string) ([]*FoodResponse, error) {
	var resp []*FoodResponse
	for _, item := range pReq.Results.Bindings {
		// 积极部分不需要病菌信息，故返回空结构体
		resp = append(resp, &FoodResponse{
			Bacname:      "",
			Bac:          "",
			CompoundName: item.Compoundname.Value,
			Compound:     item.Compound.Value,
			Food:         foodName,
			Type:         "Positive",
		})
	}
	for _, item := range nReq.Results.Bindings {
		resp = append(resp, &FoodResponse{
			Bacname:      item.Bacname.Value,
			Bac:          item.Bac.Value,
			CompoundName: item.Compoundname.Value,
			Compound:     item.Compound.Value,
			Food:         foodName,
			Type:         "Negative",
		})
	}

	return resp, nil
}
