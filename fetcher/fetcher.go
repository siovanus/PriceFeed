package fetcher

func FetchOkex(url string) (uint64, error) {
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
