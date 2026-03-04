package http

type HttpGateway interface {
	// some methods to call 3rd party services
	CollectCurrencyRates(req gatewayentity.CollectCurrencyRatesRequest) (gatewayentity.CollectCurrencyRatesResponse, error)
}
type httpGateway struct {
	// mostly used for integration with 3rd party services
	thirdPartyServiceUrl string
}

func New(thirdPartyServiceUrl string) HttpGateway {
	return &httpGateway{
		thirdPartyServiceUrl: thirdPartyServiceUrl,
	}
}
