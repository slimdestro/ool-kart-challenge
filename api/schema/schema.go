package schema

// Product represents a food item available for ordering.
type Product struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Category string  `json:"category"`
}

// OrderReqItem is a single item in the incoming order request.
type OrderReqItem struct {
	ProductID string `json:"productId"`
	Quantity  int    `json:"quantity"`
}

// OrderReq is the payload for placing a new order.
type OrderReq struct {
	CouponCode string         `json:"couponCode,omitempty"`
	Items      []OrderReqItem `json:"items"`
}

// OrderItem represents an item that was successfully ordered.
type OrderItem struct {
	ProductID string `json:"productId"`
	Quantity  int    `json:"quantity"`
}

type Order struct {
	ID       string      `json:"id"`
	Items    []OrderItem `json:"items"`
	Products []Product   `json:"products"`
}

// ApiResponse defines the standard error response structure used for HTTP responses.
type ApiResponse struct {
	Code    int    `json:"code"`
	Type    string `json:"type"`
	Message string `json:"message"`
}
