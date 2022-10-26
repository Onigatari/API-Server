package service

import (
	"Avito/internal/models"
	"Avito/internal/repository"
	"github.com/gin-gonic/gin"
)

type BalanceOperationsService struct {
	repo repository.BalanceOperations
}

func NewBalanceOperationsService(repo repository.BalanceOperations) *BalanceOperationsService {
	return &BalanceOperationsService{repo: repo}
}

func (s *BalanceOperationsService) GetBalance(userid models.GetBalanceRequest, ctx *gin.Context) (models.GetBalanceResponse, error) {
	return s.repo.GetBalance(userid, ctx)
}

func (s *BalanceOperationsService) Deposit(depReq models.UpdateBalanceRequest, ctx *gin.Context) (models.UpdateBalanceDepositResponse, error) {
	return s.repo.Deposit(depReq, ctx)
}

func (s *BalanceOperationsService) Withdrawal(withdrawReq models.UpdateBalanceRequest, ctx *gin.Context) (models.UpdateBalanceWithdrawResponse, error) {
	return s.repo.Withdrawal(withdrawReq, ctx)
}

func (s *BalanceOperationsService) Transfer(transferReq models.TransferRequest, ctx *gin.Context) (models.TransferResponse, error) {
	return s.repo.Transfer(transferReq, ctx)
}

func (s *BalanceOperationsService) ReserveService(reserveReq models.ReserveServiceRequest, ctx *gin.Context) (models.ReserveServiceResponse, error) {
	return s.repo.ReserveService(reserveReq, ctx)
}

func (s *BalanceOperationsService) ApproveService(appServiceReq models.StatusServiceRequest, ctx *gin.Context) (models.StatusServiceResponse, error) {
	return s.repo.ApproveService(appServiceReq, ctx)
}

func (s *BalanceOperationsService) FailedService(failedServiceReq models.StatusServiceRequest, ctx *gin.Context) (models.StatusServiceResponse, error) {
	return s.repo.FailedService(failedServiceReq, ctx)
}
