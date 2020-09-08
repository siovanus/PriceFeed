package service

import (
	"time"

	"github.com/ontio/ontology/common"
	"github.com/ontio/ontology/core/payload"
	"github.com/ontio/ontology/core/types"
	"github.com/ontio/ontology/smartcontract/states"

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

	PUTUNDERLYINGPRICE = "putUnderlyingPrice"
	GETUNDERLYINGPRICE = "getUnderlyingPrice"
)

func (this *PriceFeedService) parseOntData() {
	for {
		var sum uint64 = 0
		var length uint64 = 0
		okexUrl := "https://www.okex.com/api/spot/v3/instruments/ONT-USDT/ticker"
		okexPrice, err := fetcher.FetchOkex(okexUrl)
		if err != nil {
			log.Errorf("parseOntData, fetcher.FetchOkex %s error: %s", okexUrl, err)
			this.failSum += 1
		} else {
			sum += okexPrice
			length += 1
		}

		binanceUrl := "https://api.binance.com/api/v3/ticker/price?symbol=ONTUSDT"
		binancePrice, err := fetcher.FetchBinance(binanceUrl)
		if err != nil {
			log.Errorf("parseOntData, fetcher.FetchBinance %s error: %s", binanceUrl, err)
			this.failSum += 1
		} else {
			sum += binancePrice
			length += 1
		}

		huobiUrl := "https://api.huobi.pro/market/trade?symbol=ontusdt"
		huobiPrice, err := fetcher.FetchHuobi(huobiUrl)
		if err != nil {
			log.Errorf("parseOntData, fetcher.FetchHuobi %s error: %s", huobiUrl, err)
			this.failSum += 1
		} else {
			sum += huobiPrice
			length += 1
		}

		price := sum / length
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
	time.Sleep(time.Duration(10*config.DefConfig.ScanInterval) * time.Second)
	for {
		time.Sleep(time.Duration(config.DefConfig.ScanInterval) * time.Second)

		contractAddress, err := common.AddressFromHexString(config.DefConfig.OracleAddress)
		if err != nil {
			log.Errorf("fulfillOracle, oracle contract address format error")
			continue
		}

		allKeys := []string{ONT, BTC, ETH, DAI, USDT}
		keys := make([]string, 0)
		values := make([]uint64, 0)
		for _, v := range allKeys {
			result, err := this.ontologySdk.WasmVM.PreExecInvokeWasmVMContract(contractAddress, GETUNDERLYINGPRICE, []interface{}{v})
			if err != nil {
				log.Errorf("fulfillOracle, this.ontologySdk.WasmVM.PreExecInvokeWasmVMContract error")
				continue
			}
			r, err := result.Result.ToByteArray()
			if err != nil {
				log.Errorf("fulfillOracle, result.Result.ToByteArray error")
				continue
			}
			source := common.NewZeroCopySource(r)
			p, eof := source.NextI128()
			if eof {
				log.Errorf("fulfillOracle, source.NextI128 error")
			}
			prePrice := p.ToBigInt().Uint64()
			currentPrice := this.prices[v].GetPrice()
			var delta uint64 = 0
			if prePrice > currentPrice {
				delta = prePrice - currentPrice
			} else {
				delta = currentPrice - prePrice
			}
			if delta >= prePrice/100 {
				keys = append(keys, v)
				values = append(values, currentPrice)
				log.Infof("fulfillOracle, token %s, price %d", v, currentPrice)
			}
		}

		if len(keys) > 0 {
			sink := common.NewZeroCopySink(nil)
			sink.WriteString(PUTUNDERLYINGPRICE)
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
			err = this.ontologySdk.SignToTransaction(tx, this.account)
			if err != nil {
				log.Errorf("fulfillOracle, this.ontologySdk.SignToTransaction error: %s", err)
				continue
			}

			txHash, err := this.ontologySdk.SendTransaction(tx)
			if err != nil {
				log.Errorf("fulfillOracle, this.ontologySdk.SendTransaction error: %s", err)
				continue
			}
			log.Infof("fulfillOracle success, txHash is: %s", txHash.ToHexString())
		}
	}
}

func (this *PriceFeedService) checkFail() {
	for {
		time.Sleep(time.Duration(config.DefConfig.CheckFailInterval) * time.Second)

		if this.failSum >= config.DefConfig.CheckFailMax {
			log.Errorf("Waring! Sum of failed request is too much")
		}
		this.failSum = 0
	}
}
