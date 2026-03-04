package http

func (h *httpGateway) CollectCurrencyRates(req gatewayentity.CollectCurrencyRatesRequest) (gatewayentity.CollectCurrencyRatesResponse, error) {
	// need to make http client call to the 3rd party service to get the currency rates
	// for example, we can use the "net/http" package to make the http call

	return gatewayentity.CollectCurrencyRatesResponse{}, nil
}
