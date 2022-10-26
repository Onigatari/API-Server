package models

import "time"

type Order struct {
	Id        int64     `json:"order-id"`
	UserId    int64     `json:"user-id"`
	ServiceId int64     `json:"service-id"`
	OrderId   int64     `json:"order-id"`
	Amount    int64     `json:"money-amount"`
	Status    string    `json:"status"`
	Timecode  time.Time `json:"timecode"`
}

type ReserveServiceRequest struct {
	UserId    int64 `json:"user-id"`
	ServiceId int64 `json:"service-id"`
	OrderId   int64 `json:"order-id"`
	Payment   int64 `json:"payment"`
}

type ReserveServiceResponse struct {
	UserId    int64     `json:"user-id"`
	ServiceId int64     `json:"service-id"`
	OrderId   int64     `json:"order-id"`
	Invoice   int64     `json:"invoice"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created-at"`
	UpdatedAt time.Time `json:"updated-at"`
}

type StatusServiceRequest struct {
	UserId    int64 `json:"user-id"`
	ServiceId int64 `json:"service-id"`
	OrderId   int64 `json:"order-id"`
	Payment   int64 `json:"payment"`
}

type StatusServiceResponse struct {
	UserId    int64     `json:"user-id"`
	ServiceId int64     `json:"service-id"`
	OrderId   int64     `json:"order-id"`
	Invoice   int64     `json:"invoice"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created-at"`
	UpdatedAt time.Time `json:"updated-at"`
}
