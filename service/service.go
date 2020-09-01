package service

import (
	sdk "github.com/ontio/ontology-go-sdk"
	"github.com/siovanus/PriceFeed/config"
)

type PriceFeedService struct {
	account     *sdk.Account
	ontologySdk *sdk.OntologySdk
	config      *config.Config
}

func NewPriceFeedService(account *sdk.Account, ontologySdk *sdk.OntologySdk) *PriceFeedService {
	svr := &PriceFeedService{
		account:     account,
		ontologySdk: ontologySdk,
		config:      config.DefConfig,
	}
	return svr
}

func (this *PriceFeedService) Run() {
	go this.SideToAlliance()
	go this.AllianceToSide()
	go this.ProcessToAllianceCheckAndRetry()
}
