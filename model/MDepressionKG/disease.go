package MDepressionKG

import (
	"encoding/json"
	"github.com/NLPMicrobeKG-CCNU/NLPMicrobeKG-backend/service/graphDB"
)

type DiseaseQueryResponse struct {
	Head    DiseaseHead    `json:"head"`
	Results DiseaseResults `json:"results"`
}
type DiseaseHead struct {
	Vars []string `json:"vars"`
}
type DiseaseBacname struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}
type DiseaseBac struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}
type Syndrome struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}
type DiseaseDepression struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}
type DiseaseBindings struct {
	Bacname    DiseaseBacname    `json:"bacname"`
	Bac        DiseaseBac        `json:"bac"`
	Syndrome   Syndrome          `json:"syndrome"`
	Depression DiseaseDepression `json:"depression"`
}
type DiseaseResults struct {
	Bindings []DiseaseBindings `json:"bindings"`
}

type DiseaseResponse struct {
	Bacname         DiseaseBacname `json:"bacname"`
	Bac             DiseaseBac     `json:"bac"`
	Syndrome        Syndrome       `json:"syndrome"`
	Type            string         `json:"type"`
	RelevantDisease string         `json:"relevant_disease"`
}

// DiseaseQuery returns results of disease query.
func DiseaseQuery(pQuery, nQuery, relevantDiseaseName string, limit, offset int) ([]*DiseaseResponse, error) {
	pRes, err := GetDiseaseQueryRes(pQuery, limit, offset)
	if err != nil {
		return []*DiseaseResponse{}, err
	}
	nRes, err := GetDiseaseQueryRes(nQuery, limit, offset)
	if err != nil {
		return []*DiseaseResponse{}, err
	}

	return TransformToDisease(pRes, nRes, relevantDiseaseName)
}

// GetDiseaseQueryRes returns raw graphdb query response in type of disease.
func GetDiseaseQueryRes(query string, limit, offset int) (*DiseaseQueryResponse, error) {
	var res *DiseaseQueryResponse

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
func TransformToDisease(pReq *DiseaseQueryResponse, nReq *DiseaseQueryResponse, relevantDiseaseName string) ([]*DiseaseResponse, error) {
	var resp []*DiseaseResponse
	for _, item := range pReq.Results.Bindings {
		resp = append(resp, &DiseaseResponse{
			Bacname:         item.Bacname,
			Bac:             item.Bac,
			Syndrome:        item.Syndrome,
			RelevantDisease: relevantDiseaseName,
			Type:            "Positive",
		})
	}
	for _, item := range nReq.Results.Bindings {
		resp = append(resp, &DiseaseResponse{
			Bacname:         item.Bacname,
			Bac:             item.Bac,
			Syndrome:        item.Syndrome,
			RelevantDisease: relevantDiseaseName,
			Type:            "Negative",
		})
	}

	return resp, nil
}
