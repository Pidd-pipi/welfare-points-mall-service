package dto

// OrderCreateRequest 兑换下单请求。
type OrderCreateRequest struct {
	ProductID    uint   `json:"product_id" binding:"required"`
	Quantity     int    `json:"quantity" binding:"required,min=1"`
	ReceiverInfo string `json:"receiver_info" binding:"max=255"`
}

// ShipOrderRequest 发货请求。
type ShipOrderRequest struct {
	LogisticsCompany string `json:"logistics_company" binding:"required,max=64"`
	LogisticsNo      string `json:"logistics_no" binding:"required,max=64"`
}
