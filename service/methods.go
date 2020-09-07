package service

import (
	"github.com/ontio/ontology/common"
	"github.com/ontio/ontology/core/payload"
	"github.com/ontio/ontology/core/types"
	"github.com/ontio/ontology/smartcontract/states"
	"time"

	"github.com/siovanus/PriceFeed/config"
	"github.com/siovanus/PriceFeed/fetcher"
	"github.com/siovanus/PriceFeed/log"
)

const (
	ONT  = "ONT"
	BTC  = "BTC"
	ETH  = "ETH"
	DAI  = "DAI"
	USDT = "USDT"

	USDTPRICE = 1

	FULFILLORACLE = "putUnderlyingPrice"
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

		sink := common.NewZeroCopySink(nil)
		sink.WriteString(FULFILLORACLE)
		keys := []string{ONT, BTC, ETH, DAI, USDT}
		values := []uint64{this.prices[ONT].GetPrice(), this.prices[BTC].GetPrice(),
			this.prices[ETH].GetPrice(), this.prices[DAI].GetPrice(), USDTPRICE}
		length := uint64(len(keys))
		sink.WriteVarUint(length)
		for _, v := range keys {
			sink.WriteString(v)
		}
		sink.WriteVarUint(length)
		for _, v := range values {
			sink.WriteI128(common.I128FromUint64(v))
		}

		contract := &states.WasmContractParam{}
		contract.Address = contractAddress
		//bf := bytes.NewBuffer(nil)
		argbytes := sink.Bytes()
		contract.Args = argbytes

		invokePayload := &payload.InvokeCode{
			Code: common.SerializeToBytes(contract),
		}
		tx := &types.MutableTransaction{
			Payer:    this.account.Address,
			GasPrice: 2500,
			GasLimit: 300000,
			TxType:   types.InvokeWasm,
			Nonce:    uint32(time.Now().Unix()),
			Payload:  invokePayload,
			Sigs:     nil,
		}
		this.ontologySdk.SignToTransaction(tx, this.account)

		txHash, err := this.ontologySdk.SendTransaction(tx)
		if err != nil {
			log.Errorf("fulfillOracle, this.ontologySdk.SendTransaction error: %s", err)
			continue
		}
		log.Infof("fulfillOracle success, txHash is: %s", txHash.ToHexString())
	}
}
