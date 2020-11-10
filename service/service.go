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
	svr.prices[WBTC] = NewPrices()
	svr.prices[RENBTC] = NewPrices()
	svr.prices[ETH] = NewPrices()
	svr.prices[ETH9] = NewPrices()
	svr.prices[DAI] = NewPrices()
	svr.prices[USDC] = NewPrices()
	svr.prices[WING] = NewPrices()
	svr.prices[NEO] = NewPrices()
	svr.prices[UNI] = NewPrices()
	svr.prices[OKB] = NewPrices()
	svr.prices[ONG] = NewPrices()
	svr.prices[USDT] = &Prices{[]uint64{uint64(USDTPRICE * math.Pow10(fetcher.DECIMAL))}}
	svr.prices[SUSD] = &Prices{[]uint64{uint64(USDTPRICE * math.Pow10(fetcher.DECIMAL))}}
	return svr
}

func (this *PriceFeedService) Run() {
	go this.parseOntData()
	go this.parseBtcData()
	go this.parseEthData()
	go this.parseDaiData()
	go this.parseUsdcData()
	go this.parseWingData()
	go this.parseNeoData()
	go this.parseUniData()
	go this.parseOkbData()
	go this.parseOngData()
	go this.fulfillOracle()
	go this.checkFail()
}
