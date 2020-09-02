package service

import (
	sdk "github.com/ontio/ontology-go-sdk"
)

type PriceFeedService struct {
	account     *sdk.Account
	ontologySdk *sdk.OntologySdk
	prices      map[string]*Prices
}

func NewPriceFeedService(account *sdk.Account, ontologySdk *sdk.OntologySdk) *PriceFeedService {
	svr := &PriceFeedService{
		account:     account,
		ontologySdk: ontologySdk,
		prices:      make(map[string]*Prices),
	}
	svr.prices[ONT] = NewPrices()
	svr.prices[BTC] = NewPrices()
	svr.prices[ETH] = NewPrices()
	svr.prices[DAI] = NewPrices()
	return svr
}

func (this *PriceFeedService) Run() {
	go this.parseOntData()
	go this.parseBtcData()
	go this.parseEthData()
	go this.parseDaiData()
	go this.fulfillOracle()
}
