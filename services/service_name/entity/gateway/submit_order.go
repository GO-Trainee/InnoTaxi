package gateway

type SubmitOrderRequest struct {
	OrderID   string  `json:"order_id"`
	ProductID string  `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

type SubmitOrderResponse struct {
	Status  OrderStatus `json:"status"`
	Message string      `json:"message"`
}
