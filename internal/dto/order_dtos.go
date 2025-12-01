package dto

type CreateSingleOrderDTO struct {
	Quantity        int    `json:"quantity" binding:"required,min=1"`
	ShippingAddress string `json:"shipping_address" binding:"required"`
	PaymentMethod   string `json:"payment_method" binding:"required"`
}

type CreateCartOrderDTO struct {
	ShippingAddress string `json:"shipping_address" binding:"required"`
	PaymentMethod   string `json:"payment_method" binding:"required"`
}

type UpdateOrderStatusDTO struct {
	Status string `json:"status" binding:"required"`
}
