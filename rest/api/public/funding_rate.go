package public

import "github.com/iaping/go-okx/rest/api"

func NewFundingRate(param *FundingRateParam) (api.IRequest, api.IResponse) {
	return &api.Request{
		Path:   "/api/v5/public/funding-rate",
		Method: api.MethodGet,
		Param:  param,
	}, &FundingRateResponse{}
}

type FundingRateParam struct {
	InstId string `url:"instId"`
}

type FundingRateResponse struct {
	api.Response
	Data []FundingRate `json:"data"`
}

type FundingRate struct {
	FormulaType     string `json:"formulaType"`
	FundingRate     string `json:"fundingRate"`
	FundingTime     string `json:"fundingTime"`
	ImpactValue     string `json:"impactValue"`
	InstId          string `json:"instId"`
	InstType        string `json:"instType"`
	InterestRate    string `json:"interestRate"`
	MaxFundingRate  string `json:"maxFundingRate"`
	Method          string `json:"method"`
	MinFundingRate  string `json:"minFundingRate"`
	NextFundingRate string `json:"nextFundingRate"`
	NextFundingTime string `json:"nextFundingTime"`
	Premium         string `json:"premium"`
	SettFundingRate string `json:"settFundingRate"`
	SettState       string `json:"settState"`
	Ts              string `json:"ts"`
}
