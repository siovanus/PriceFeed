package service

const MAX_LENGTH = 60

type Prices struct {
	Prices []uint64
}

func NewPrices() *Prices {
	return &Prices{make([]uint64, 0)}
}

func (this *Prices) Push(price uint64) {
	length := len(this.Prices)
	if length < MAX_LENGTH {
		this.Prices = append(this.Prices, price)
	} else {
		this.Prices = append(this.Prices[1:], price)
	}
}

func (this *Prices) GetPrice() uint64 {
	if len(this.Prices) == 0 {
		return 0
	}
	var sum uint64 = 0
	for _, v := range this.Prices {
		sum += v
	}
	return sum / uint64(len(this.Prices))
}
