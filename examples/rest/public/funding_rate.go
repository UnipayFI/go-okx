package main

import (
	"log"

	"github.com/iaping/go-okx/examples/rest"
	"github.com/iaping/go-okx/rest/api/public"
)

func main() {
	param := &public.FundingRateParam{
		InstId: "BTC-USDT-SWAP",
	}
	req, resp := public.NewFundingRate(param)
	if err := rest.TestClient.Do(req, resp); err != nil {
		panic(err)
	}
	log.Println(req, resp.(*public.FundingRateResponse))
}
