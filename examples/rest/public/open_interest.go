package main

import (
	"log"

	"github.com/iaping/go-okx/examples/rest"
	"github.com/iaping/go-okx/rest/api"
	"github.com/iaping/go-okx/rest/api/public"
)

func main() {
	param := &public.OpenInterestParam{
		InstType: api.InstTypeSWAP,
		InstId:   "ZBT-USDT-SWAP",
	}
	req, resp := public.NewOpenInterest(param)
	if err := rest.TestClient.Do(req, resp); err != nil {
		panic(err)
	}
	log.Println(req, resp.(*public.OpenInterestResponse))
}
