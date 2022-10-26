package repository

import (
	"Avito/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type BalanceOperations interface {
	GetBalance(userId models.GetBalanceRequest, ctx *gin.Context) (models.GetBalanceResponse, error)
	Deposit(depositReq models.UpdateBalanceRequest, ctx *gin.Context) (models.UpdateBalanceDepositResponse, error)
	Withdrawal(withdrawReq models.UpdateBalanceRequest, ctx *gin.Context) (models.UpdateBalanceWithdrawResponse, error)
	Transfer(transferReq models.TransferRequest, ctx *gin.Context) (models.TransferResponse, error)
	ReserveService(reserveSerFeeReq models.ReserveServiceRequest, ctx *gin.Context) (models.ReserveServiceResponse, error)
	ApproveService(approveSerFeeReq models.StatusServiceRequest, ctx *gin.Context) (models.StatusServiceResponse, error)
	FailedService(failedServiceFeeReq models.StatusServiceRequest, ctx *gin.Context) (models.StatusServiceResponse, error)
}

type AppRepository struct {
	BalanceOperations
}

func NewRepository(db *sqlx.DB) *AppRepository {
	return &AppRepository{
		BalanceOperations: NewPostgres(db),
	}
}
