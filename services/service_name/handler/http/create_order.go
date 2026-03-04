package http

import httpentity "awesomeProject/services/service_name/entity/http"

func (h *HttpHandler) CreateOrder(req httpentity.CreateOrderRequest) (httpentity.CreateOrderResponse, error) {
	// validate the request
	// call the service layer to create the order
	return httpentity.CreateOrderResponse{}, nil
}
