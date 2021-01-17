package query

import (
	"encoding/json"
	"fmt"

	"github.com/NLPMicrobeKG-CCNU/NLPMicrobeKG-backend/service/graphDB"
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
	Bacname   string `json:"bacname"`
	Reference string `json:"reference"`
	Bac4Name  string `json:"bac4name"`
	Ref4      string `json:"ref4"`
	Ref3      string `json:"ref3"`
	Bac2Name  string `json:"bac2name"`
	Ref2      string `json:"ref2"`
	Bac3Name  string `json:"bac3name"`
	Disname   string `json:"disname"`
}

// TextQuery returns results of text query.
func TextQuery(query, bacname string, limit, offset int) ([]*TextResponse, error) {
	res, err := GetTextQueryRes(query, limit, offset)
	if err != nil {
		fmt.Println(err)
		return []*TextResponse{}, err
	}

	return TransformToText(res, bacname)
}

// GetTextQueryRes returns raw graphdb query response in type of text.
func GetTextQueryRes(query string, limit, offset int) (*TextQueryResponse, error) {
	var res *TextQueryResponse

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

// TransformToText transform graphdb query response into service text response.
func TransformToText(req *TextQueryResponse, bacname string) ([]*TextResponse, error) {
	var resp []*TextResponse
	for _, item := range req.Results.Bindings {
		resp = append(resp, &TextResponse{
			Bacname:   bacname,
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
