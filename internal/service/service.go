package service

import (
	"Avito/internal/models"
	"Avito/internal/repository"
	"github.com/gin-gonic/gin"
)

type BalanceOperations interface {
	GetBalance(userId models.GetBalanceRequest, ctx *gin.Context) (models.GetBalanceResponse, error)
	DepositMoney(depositReq models.UpdateBalanceRequest, ctx *gin.Context) (models.UpdateBalanceDepositResponse, error)
	WithdrawMoney(withdrawReq models.UpdateBalanceRequest, ctx *gin.Context) (models.UpdateBalanceWithdrawResponse, error)
	Transfer(transferReq models.TransferRequest, ctx *gin.Context) (models.TransferResponse, error)
	ReserveServiceFee(reserveSerFeeReq models.ReserveServiceFeeRequest, ctx *gin.Context) (models.ReserveServiceFeeResponse, error)
	ApproveServiceFee(approveSerFeeReq models.StatusServiceFeeRequest, ctx *gin.Context) (models.StatusServiceFeeResponse, error)
	FailedServiceFee(failedServiceFeeReq models.StatusServiceFeeRequest, ctx *gin.Context) (models.StatusServiceFeeResponse, error)
}

type BillingService struct {
	BalanceOperations
}

func NewService(repos *repository.BillingRepo) *BillingService {
	return &BillingService{
		BalanceOperations: NewBalanceOperationsService(repos.BalanceOperations),
	}
}
