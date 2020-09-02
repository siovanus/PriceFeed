package service

type Prices struct {
	Prices []uint64
}

func NewPrices() Prices {
	return Prices{make([]uint64, 0)}
}
