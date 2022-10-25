package models

import "time"

// Users /* */
type Users struct {
	Id            int64     `json:"order-id"`
	UserId        int64     `json:"user-id" binding:"required"`
	CurrAmount    float64   `json:"money-amount" binding:"required"`
	PendingAmount string    `json:"pending-amount"`
	LastUpdatedAt time.Time `json:"last-updated-time"`
}

type GetBalanceRequest struct {
	UserId int64 `json:"user-id" binding:"required"`
}

type GetBalanceResponse struct {
	Balance int64 `json:"user-balance"`
	Pending int64 `json:"user-pending-amount"`
}

type DepositMoneyRequest struct {
	UserId        int64 `json:"user-id" binding:"required"`
	DepositAmount int64 `json:"deposit-amount" binding:"required"`
}

type UpdateBalanceRequest struct {
	UserId int64 `json:"user-id" binding:"required"`
	Sum    int64 `json:"update-amount" binding:"required"`
}

type UpdateBalanceDepositResponse struct {
	UserId    int64     `json:"user-id" binding:"required"`
	Sum       int64     `json:"deposit-sum" binding:"required"`
	Status    string    `json:"operation-status" binding:"required"`
	EventType string    `json:"operation-event"`
	CreatedAt time.Time `json:"created-at"`
}

type UpdateBalanceWithdrawResponse struct {
	UserId    int64     `json:"user-id" binding:"required"`
	Sum       int64     `json:"deposit-sum" binding:"required"`
	Status    string    `json:"operation-status" binding:"required"`
	EventType string    `json:"operation-event"`
	CreatedAt time.Time `json:"created-at"`
}
