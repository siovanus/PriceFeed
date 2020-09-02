package service

import (
	"github.com/ontio/ontology/common"
	"time"

	"github.com/siovanus/PriceFeed/config"
	"github.com/siovanus/PriceFeed/fetcher"
	"github.com/siovanus/PriceFeed/log"
)

const (
	ONT = "ONT"
	BTC = "BTC"
	ETH = "ETH"
	DAI = "DAI"

	FULFILLORACLE = "fulfillOracle"
)

func (this *PriceFeedService) parseOntData() {
	for {
		okexUrl := "https://www.okex.com/api/spot/v3/instruments/ONT-USDT/ticker"
		okexPrice, err := fetcher.FetchOkex(okexUrl)
		if err != nil {
			log.Errorf("parseOntData, fetcher.FetchOkex %s error: %s", okexUrl, err)
			continue
		}

		binanceUrl := "https://api.binance.com/api/v3/ticker/price?symbol=ONTUSDT"
		binancePrice, err := fetcher.FetchBinance(binanceUrl)
		if err != nil {
			log.Errorf("parseOntData, fetcher.FetchBinance %s error: %s", binanceUrl, err)
			continue
		}

		huobiUrl := "https://api.huobi.pro/market/trade?symbol=ontusdt"
		huobiPrice, err := fetcher.FetchHuobi(huobiUrl)
		if err != nil {
			log.Errorf("parseOntData, fetcher.FetchHuobi %s error: %s", huobiUrl, err)
			continue
		}

		price := (okexPrice + binancePrice + huobiPrice) / 3
		this.prices[ONT].Push(price)

		time.Sleep(time.Duration(config.DefConfig.ScanInterval))
	}
}

func (this *PriceFeedService) parseBtcData() {
	for {
		okexUrl := "https://www.okex.com/api/spot/v3/instruments/BTC-USDT/ticker"
		okexPrice, err := fetcher.FetchOkex(okexUrl)
		if err != nil {
			log.Errorf("parseBtcData, fetcher.FetchOkex %s error: %s", okexUrl, err)
			continue
		}

		binanceUrl := "https://api.binance.com/api/v3/ticker/price?symbol=BTCUSDT"
		binancePrice, err := fetcher.FetchBinance(binanceUrl)
		if err != nil {
			log.Errorf("parseBtcData, fetcher.FetchBinance %s error: %s", binanceUrl, err)
			continue
		}

		huobiUrl := "https://api.huobi.pro/market/trade?symbol=btcusdt"
		huobiPrice, err := fetcher.FetchHuobi(huobiUrl)
		if err != nil {
			log.Errorf("parseBtcData, fetcher.FetchHuobi %s error: %s", huobiUrl, err)
			continue
		}

		price := (okexPrice + binancePrice + huobiPrice) / 3
		this.prices[BTC].Push(price)

		time.Sleep(time.Duration(config.DefConfig.ScanInterval))
	}
}

func (this *PriceFeedService) parseEthData() {
	for {
		okexUrl := "https://www.okex.com/api/spot/v3/instruments/ETH-USDT/ticker"
		okexPrice, err := fetcher.FetchOkex(okexUrl)
		if err != nil {
			log.Errorf("parseEthData, fetcher.FetchOkex %s error: %s", okexUrl, err)
			continue
		}

		binanceUrl := "https://api.binance.com/api/v3/ticker/price?symbol=ETHUSDT"
		binancePrice, err := fetcher.FetchBinance(binanceUrl)
		if err != nil {
			log.Errorf("parseEthData, fetcher.FetchBinance %s error: %s", binanceUrl, err)
			continue
		}

		huobiUrl := "https://api.huobi.pro/market/trade?symbol=ethusdt"
		huobiPrice, err := fetcher.FetchHuobi(huobiUrl)
		if err != nil {
			log.Errorf("parseEthData, fetcher.FetchHuobi %s error: %s", huobiUrl, err)
			continue
		}

		price := (okexPrice + binancePrice + huobiPrice) / 3
		this.prices[ETH].Push(price)

		time.Sleep(time.Duration(config.DefConfig.ScanInterval))
	}
}

func (this *PriceFeedService) parseDaiData() {
	for {
		okexUrl := "https://www.okex.com/api/spot/v3/instruments/DAI-USDT/ticker"
		okexPrice, err := fetcher.FetchOkex(okexUrl)
		if err != nil {
			log.Errorf("parseDaiData, fetcher.FetchOkex %s error: %s", okexUrl, err)
			continue
		}

		binanceUrl := "https://api.binance.com/api/v3/ticker/price?symbol=DAIUSDT"
		binancePrice, err := fetcher.FetchBinance(binanceUrl)
		if err != nil {
			log.Errorf("parseDaiData, fetcher.FetchBinance %s error: %s", binanceUrl, err)
			continue
		}

		huobiUrl := "https://api.huobi.pro/market/trade?symbol=daiusdt"
		huobiPrice, err := fetcher.FetchHuobi(huobiUrl)
		if err != nil {
			log.Errorf("parseDaiData, fetcher.FetchHuobi %s error: %s", huobiUrl, err)
			continue
		}

		price := (okexPrice + binancePrice + huobiPrice) / 3
		this.prices[DAI].Push(price)

		time.Sleep(time.Duration(config.DefConfig.ScanInterval) * time.Second)
	}
}

func (this *PriceFeedService) fulfillOracle() {
	for {
		time.Sleep(time.Duration(config.DefConfig.FulfillInterval) * time.Second)

		contractAddress, err := common.AddressFromHexString(config.DefConfig.OracleAddress)
		if err != nil {
			log.Errorf("fulfillOracle, oracle contract address format error")
			continue
		}
		txHash, err := this.ontologySdk.WasmVM.InvokeWasmVMSmartContract(config.DefConfig.GasPrice, config.DefConfig.GasLimit,
			this.account, this.account, contractAddress, FULFILLORACLE,
			[]interface{}{[]string{ONT, BTC, ETH, DAI}, []uint64{this.prices[ONT].GetPrice(), this.prices[BTC].GetPrice(),
				this.prices[ETH].GetPrice(), this.prices[DAI].GetPrice()}})
		if err != nil {
			log.Errorf("fulfillOracle, invoke oracle contract error: %s", err)
			continue
		}
		log.Infof("fulfillOracle success, txHash is: %s", txHash.ToHexString())
	}
}
