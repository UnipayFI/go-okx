package public

import "github.com/iaping/go-okx/rest/api"

func NewOpenInterest(param *OpenInterestParam) (api.IRequest, api.IResponse) {
	return &api.Request{
		Path:   "/api/v5/public/open-interest",
		Method: api.MethodGet,
		Param:  param,
	}, &OpenInterestResponse{}
}

type OpenInterestParam struct {
	InstType   string `url:"instType"`
	InstFamily string `json:"instFamily,omitempty"`
	InstId     string `url:"instId,omitempty"`
}

type OpenInterestResponse struct {
	api.Response
	Data []OpenInterest `json:"data"`
}

type OpenInterest struct {
	InstType string `json:"instType"`
	InstId   string `json:"instId"`
	Oi       string `json:"oi"`
	OiCcy    string `json:"oiCcy"`
	OiUsd    string `json:"oiUsd"`
	Ts       int64  `json:"ts,string"`
}
