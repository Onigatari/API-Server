package handler

import (
	"Avito/internal/models"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

func (h *Handler) balance(c *gin.Context) {
	idStringInput := c.Param("id")
	log.Printf("Input read: %v %T", idStringInput, idStringInput)

	idNumberInput, err := strconv.ParseInt(idStringInput, 10, 64)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	input := models.GetBalanceRequest{UserId: idNumberInput}

	response, newErr := h.services.BalanceOperations.GetBalance(input, c)
	if newErr != nil {
		NewErrorResponse(c, http.StatusInternalServerError, newErr.Error())
	} else {
		c.JSON(http.StatusOK, gin.H{
			"user-balance":        response.Balance,
			"user-pending-amount": response.Pending,
		})
	}
}

func (h *Handler) deposit(c *gin.Context) {
	var updateBalanceDepositRequest models.UpdateBalanceRequest

	if err := c.BindJSON(&updateBalanceDepositRequest); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	response, err := h.services.BalanceOperations.DepositMoney(updateBalanceDepositRequest, c)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
	} else {
		c.JSON(http.StatusOK, gin.H{
			"account-id":       response.UserId,
			"sum-deposited":    response.Sum,
			"operation-status": response.Status,
			"operation-event":  response.EventType,
			"created-at":       response.CreatedAt,
		})
	}
}

func (h *Handler) withdrawal(c *gin.Context) {
	var updateBalanceWithdrawRequest models.UpdateBalanceRequest

	if err := c.BindJSON(&updateBalanceWithdrawRequest); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	response, err := h.services.BalanceOperations.WithdrawMoney(updateBalanceWithdrawRequest, c)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
	} else {
		c.JSON(http.StatusOK, gin.H{
			"account-id":       response.UserId,
			"sum-withdrawn":    response.Sum,
			"operation-status": response.Status,
			"operation-event":  response.EventType,
			"created-at":       response.CreatedAt,
		})
	}
}

// @Summary reserveServiceFee
// @Tags account
// @Description "Put specified amount of money in reservation for a given account"
// @Accept json
// @Produce json
// @Param input body models.ReserveServiceFeeRequest true "JSON object with used ID, service ID, order ID and fee amount"
// @Success 200 {object} models.ReserveServiceFeeResponse
// @Failure 500 {object} errorAcc
// @Router /account/reserveServiceFee [post]
func (h *Handler) reserveServiceFee(c *gin.Context) {
	var reserveServiceFeeRequest models.ReserveServiceFeeRequest

	if err := c.BindJSON(&reserveServiceFeeRequest); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	response, err := h.services.BalanceOperations.ReserveServiceFee(reserveServiceFeeRequest, c)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
	} else {
		c.JSON(http.StatusOK, gin.H{
			"account-id": response.UserId,
			"service-id": response.ServiceId,
			"order-id":   response.OrderId,
			"invoice":    response.Invoice,
			"status":     response.Status,
			"created-at": response.CreatedAt,
			"updated-at": response.UpdatedAt,
		})
	}
}

// @Summary approveOrderFee
// @Tags account
// @Description "Approve specified reservation"
// @Accept json
// @Produce json
// @Param input body models.StatusServiceFeeRequest true "JSON object with used ID, service ID, order ID and fee amount"
// @Success 200 {object} models.StatusServiceFeeResponse
// @Failure 500 {object} errorAcc
// @Router /account/approveOrderFee [post]
func (h *Handler) approveOrderFee(c *gin.Context) {
	var statusApproveServiceFeeRequest models.StatusServiceFeeRequest

	if err := c.BindJSON(&statusApproveServiceFeeRequest); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	response, err := h.services.BalanceOperations.ApproveServiceFee(statusApproveServiceFeeRequest, c)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
	} else {
		c.JSON(http.StatusOK, gin.H{
			"account-id": response.UserId,
			"service-id": response.ServiceId,
			"order-id":   response.OrderId,
			"invoice":    response.Invoice,
			"status":     response.Status,
			"created-at": response.CreatedAt,
			"updated-at": response.UpdatedAt,
		})
	}
}

// @Summary failedServiceFee
// @Tags account
// @Description "Mark reservation as failed and release funds"
// @Accept json
// @Produce json
// @Param input body models.StatusServiceFeeRequest true "JSON object with used ID, service ID, order ID and fee amount"
// @Success 200 {object} models.StatusServiceFeeResponse
// @Failure 500 {object} errorAcc
// @Router /account/failedServiceFee [post]
func (h *Handler) failedServiceFee(c *gin.Context) {
	var statusFailedServiceFeeRequest models.StatusServiceFeeRequest

	if err := c.BindJSON(&statusFailedServiceFeeRequest); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	response, err := h.services.BalanceOperations.FailedServiceFee(statusFailedServiceFeeRequest, c)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
	} else {
		c.JSON(http.StatusOK, gin.H{
			"account-id": response.UserId,
			"service-id": response.ServiceId,
			"order-id":   response.OrderId,
			"invoice":    response.Invoice,
			"status":     response.Status,
			"created-at": response.CreatedAt,
			"updated-at": response.UpdatedAt,
		})
	}
}

// @Summary transfer
// @Tags account
// @Description "Transfer funds from one account to another"
// @Accept json
// @Produce json
// @Param input body models.TransferRequest true "JSON object with sender ID, receiver ID and money amount"
// @Success 200 {object} models.TransferResponse
// @Failure 500 {object} errorAcc
// @Router /account/transfer [post]
func (h *Handler) transfer(c *gin.Context) {
	var transferRequest models.TransferRequest

	if err := c.BindJSON(&transferRequest); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	response, err := h.services.BalanceOperations.Transfer(transferRequest, c)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
	} else {
		c.JSON(http.StatusOK, gin.H{
			"receive-account":  response.UserIdTo,
			"transfer-account": response.UserIdFrom,
			"amount":           response.Amount,
			"status":           response.Status,
			"event-type":       response.EventType,
			"created-at":       response.Timecode,
		})
	}
}

func (h *Handler) pong(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "ping",
	})
}
