package service

import (
	"time"

	"github.com/siovanus/PriceFeed/config"
	"github.com/siovanus/PriceFeed/fetcher"
	"github.com/siovanus/PriceFeed/log"
)

func parseOntData() {
	okexUrl := "https://www.okex.com/api/spot/v3/instruments/ONT-USDT/ticker"
	resp, err := fetcher.Get(okexUrl)
	if err != nil {
		log.Errorf("parseOntData, fetcher.Get %s error: %s", okexUrl, err)
	}
	_, err = fetcher.OkexParse(resp)
	if err != nil {
		log.Errorf("parseOntData, fetcher.OkexParse error: %s", err)
	}

	time.Sleep(time.Duration(config.DefConfig.Interval))
}
