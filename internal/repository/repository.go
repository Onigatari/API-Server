package repository

import (
	"Avito/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type BalanceOperations interface {
	GetBalance(userId models.GetBalanceRequest, ctx *gin.Context) (models.GetBalanceResponse, error)
	DepositMoney(depositReq models.UpdateBalanceRequest, ctx *gin.Context) (models.UpdateBalanceDepositResponse, error)
	WithdrawMoney(withdrawReq models.UpdateBalanceRequest, ctx *gin.Context) (models.UpdateBalanceWithdrawResponse, error)
	ReserveServiceFee(reserveSerFeeReq models.ReserveServiceFeeRequest, ctx *gin.Context) (models.ReserveServiceFeeResponse, error)
	ApproveServiceFee(approveSerFeeReq models.StatusServiceFeeRequest, ctx *gin.Context) (models.StatusServiceFeeResponse, error)
	Transfer(transferReq models.TransferRequest, ctx *gin.Context) (models.TransferResponse, error)
	FailedServiceFee(failedServiceFeeReq models.StatusServiceFeeRequest, ctx *gin.Context) (models.StatusServiceFeeResponse, error)
}

type BillingRepo struct {
	BalanceOperations
}

func NewRepo(db *sqlx.DB) *BillingRepo {
	return &BillingRepo{
		BalanceOperations: NewAccPostgres(db),
	}
}
