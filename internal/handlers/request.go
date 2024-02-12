package handlers

type GetOrderRequest struct {
	OrderUid string `uri:"order_uid" validate:"required"`
}
