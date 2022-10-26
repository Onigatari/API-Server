package models

import "time"

type Transaction struct {
	Id         int64     `json:"tx-id"`
	UserIdTo   int64     `json:"receive-user"`
	UserIdFrom int64     `json:"transfer-user"`
	Amount     int64     `json:"money-amount"`
	Status     string    `json:"status"`
	Timecode   time.Time `json:"timecode"`
}

type TransferRequest struct {
	SenderId   int64 `json:"sender-id" binding:"required"`
	ReceiverId int64 `json:"receiver-id" binding:"required"`
	Sum        int64 `json:"transfer-amount" binding:"required"`
}

type TransferResponse struct {
	UserIdTo   int64     `json:"receive-user"`
	UserIdFrom int64     `json:"transfer-user"`
	Amount     int64     `json:"money-amount"`
	Status     string    `json:"status"`
	EventType  string    `json:"event-type"`
	Timecode   time.Time `json:"created-at"`
}
