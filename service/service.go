package service

import (
	"math"

	sdk "github.com/ontio/ontology-go-sdk"
	"github.com/siovanus/PriceFeed/fetcher"
)

type PriceFeedService struct {
	account     *sdk.Account
	ontologySdk *sdk.OntologySdk
	prices      map[string]*Prices
	failSum     uint64
}

func NewPriceFeedService(account *sdk.Account, ontologySdk *sdk.OntologySdk) *PriceFeedService {
	svr := &PriceFeedService{
		account:     account,
		ontologySdk: ontologySdk,
		prices:      make(map[string]*Prices),
	}
	svr.prices[ONT] = NewPrices()
	svr.prices[ONTD] = NewPrices()
	svr.prices[BTC] = NewPrices()
	svr.prices[ETH] = NewPrices()
	svr.prices[DAI] = NewPrices()
	svr.prices[USDT] = &Prices{[]uint64{uint64(USDTPRICE * math.Pow10(fetcher.DECIMAL))}}
	return svr
}

func (this *PriceFeedService) Run() {
	go this.parseOntData()
	go this.parseBtcData()
	go this.parseEthData()
	go this.parseDaiData()
	go this.fulfillOracle()
	go this.checkFail()
}
