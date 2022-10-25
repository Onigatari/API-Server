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

func (s *BalanceOperationsService) DepositMoney(depReq models.UpdateBalanceRequest, ctx *gin.Context) (models.UpdateBalanceDepositResponse, error) {
	return s.repo.DepositMoney(depReq, ctx)
}

func (s *BalanceOperationsService) WithdrawMoney(withdrawReq models.UpdateBalanceRequest, ctx *gin.Context) (models.UpdateBalanceWithdrawResponse, error) {
	return s.repo.WithdrawMoney(withdrawReq, ctx)
}

func (s *BalanceOperationsService) Transfer(transferReq models.TransferRequest, ctx *gin.Context) (models.TransferResponse, error) {
	return s.repo.Transfer(transferReq, ctx)
}

func (s *BalanceOperationsService) ReserveServiceFee(reserveReq models.ReserveServiceFeeRequest, ctx *gin.Context) (models.ReserveServiceFeeResponse, error) {
	return s.repo.ReserveServiceFee(reserveReq, ctx)
}

func (s *BalanceOperationsService) ApproveServiceFee(appSerFeeReq models.StatusServiceFeeRequest, ctx *gin.Context) (models.StatusServiceFeeResponse, error) {
	return s.repo.ApproveServiceFee(appSerFeeReq, ctx)
}

func (s *BalanceOperationsService) FailedServiceFee(failedServiceFeeReq models.StatusServiceFeeRequest, ctx *gin.Context) (models.StatusServiceFeeResponse, error) {
	return s.repo.FailedServiceFee(failedServiceFeeReq, ctx)
}
