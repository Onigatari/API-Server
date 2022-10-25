package models

import "time"

// Order /* */
type Order struct {
	Id        int64     `json:"order-id"`
	UserId    int64     `json:"user-id"`
	ServiceId int64     `json:"service-id"`
	OrderId   int64     `json:"order-id"`
	Amount    int64     `json:"money-amount"`
	Status    string    `json:"status"`
	Timecode  time.Time `json:"timecode"`
}

// ReserveServiceFeeRequest
// Метод резервирования средств с основного баланса на отдельном счете.
// Принимает id пользователя, ИД услуги, ИД заказа, стоимость.

type ReserveServiceFeeRequest struct {
	UserId    int64 `json:"user-id"`
	ServiceId int64 `json:"service-id"`
	OrderId   int64 `json:"order-id"`
	Fee       int64 `json:"fee"`
}

type ReserveServiceFeeResponse struct {
	UserId    int64     `json:"user-id"`
	ServiceId int64     `json:"service-id"`
	OrderId   int64     `json:"order-id"`
	Invoice   int64     `json:"invoice"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created-at"`
	UpdatedAt time.Time `json:"updated-at"`
}

type StatusServiceFeeRequest struct {
	UserId    int64 `json:"user-id"`
	ServiceId int64 `json:"service-id"`
	OrderId   int64 `json:"order-id"`
	Fee       int64 `json:"fee"`
}

type StatusServiceFeeResponse struct {
	UserId    int64     `json:"user-id"`
	ServiceId int64     `json:"service-id"`
	OrderId   int64     `json:"order-id"`
	Invoice   int64     `json:"invoice"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created-at"`
	UpdatedAt time.Time `json:"updated-at"`
}
