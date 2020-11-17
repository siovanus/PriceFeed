package fetcher

import (
	"github.com/siovanus/PriceFeed/log"
)

func FetchOkex(url string) (uint64, error) {
	defer func() {
		if r := recover(); r != nil {
			log.Errorf("Recovered FetchHuobi: %s", r)
		}
	}()
	resp, err := Get(url)
	if err != nil {
		return 0, err
	}
	price, err := OkexParse(resp)
	if err != nil {
		return 0, err
	}
	return price, nil
}

func FetchBinance(url string) (uint64, error) {
	defer func() {
		if r := recover(); r != nil {
			log.Errorf("Recovered FetchHuobi: %s", r)
		}
	}()
	resp, err := Get(url)
	if err != nil {
		return 0, err
	}
	price, err := BinanceParse(resp)
	if err != nil {
		return 0, err
	}
	return price, nil
}

func FetchHuobi(url string) (uint64, error) {
	defer func() {
		if r := recover(); r != nil {
			log.Errorf("Recovered FetchHuobi: %s", r)
		}
	}()
	resp, err := Get(url)
	if err != nil {
		return 0, err
	}
	price, err := HuobiParse(resp)
	if err != nil {
		return 0, err
	}
	return price, nil
}
