/*
 * Copyright (C) 2018 The ontology Authors
 * This file is part of The ontology library.
 *
 * The ontology is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The ontology is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public License
 * along with The ontology.  If not, see <http://www.gnu.org/licenses/>.
 */

package fetcher

import (
	"encoding/json"
	"math"
	"strconv"
)

const DECIMAL = 9

type OkexResp struct {
	Last string `json:"last"`
}

func OkexParse(input []byte) (uint64, error) {
	okexResp := new(OkexResp)
	err := json.Unmarshal(input, okexResp)
	if err != nil {
		return 0, err
	}

	r, err := strconv.ParseFloat(okexResp.Last, 64)
	if err != nil {
		return 0, err
	}
	price := uint64(r * math.Pow10(DECIMAL))
	return price, nil
}

type BinanceResp struct {
	Price string `json:"price"`
}

func BinanceParse(input []byte) (uint64, error) {
	binanceResp := new(BinanceResp)
	err := json.Unmarshal(input, binanceResp)
	if err != nil {
		return 0, err
	}

	r, err := strconv.ParseFloat(binanceResp.Price, 64)
	if err != nil {
		return 0, err
	}
	price := uint64(r * math.Pow10(DECIMAL))
	return price, nil
}

type HuobiResp struct {
	Tick *Tick `json:"tick"`
}

type Tick struct {
	Data []*Data `json:"data"`
}

type Data struct {
	Price float64 `json:"price"`
}

func HuobiParse(input []byte) (uint64, error) {
	huobiResp := new(HuobiResp)
	err := json.Unmarshal(input, huobiResp)
	if err != nil {
		return 0, err
	}

	r := huobiResp.Tick.Data[0].Price
	price := uint64(r * math.Pow10(DECIMAL))
	return price, nil
}
