package service

import (
	"Avito/internal/models"
	"Avito/internal/repository"
	"github.com/gin-gonic/gin"
)

type BalanceOperations interface {
	GetBalance(userId models.GetBalanceRequest, ctx *gin.Context) (models.GetBalanceResponse, error)
	Deposit(depositReq models.UpdateBalanceRequest, ctx *gin.Context) (models.UpdateBalanceDepositResponse, error)
	Withdrawal(withdrawReq models.UpdateBalanceRequest, ctx *gin.Context) (models.UpdateBalanceWithdrawResponse, error)
	Transfer(transferReq models.TransferRequest, ctx *gin.Context) (models.TransferResponse, error)
	ReserveService(reserveServiceReq models.ReserveServiceRequest, ctx *gin.Context) (models.ReserveServiceResponse, error)
	ApproveService(approveServiceReq models.StatusServiceRequest, ctx *gin.Context) (models.StatusServiceResponse, error)
	FailedService(failedServiceReq models.StatusServiceRequest, ctx *gin.Context) (models.StatusServiceResponse, error)
}

type BillingService struct {
	BalanceOperations
}

func NewService(repos *repository.AppRepository) *BillingService {
	return &BillingService{
		BalanceOperations: NewBalanceOperationsService(repos.BalanceOperations),
	}
}
