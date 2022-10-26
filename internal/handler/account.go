package handler

import (
	"Avito/internal/models"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

func (h *Handler) ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

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
	response, err := h.services.BalanceOperations.Deposit(updateBalanceDepositRequest, c)
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
	response, err := h.services.BalanceOperations.Withdrawal(updateBalanceWithdrawRequest, c)
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

func (h *Handler) reserveService(c *gin.Context) {
	var reserveServiceRequest models.ReserveServiceRequest

	if err := c.BindJSON(&reserveServiceRequest); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	response, err := h.services.BalanceOperations.ReserveService(reserveServiceRequest, c)
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

func (h *Handler) approveOrder(c *gin.Context) {
	var statusApproveServiceRequest models.StatusServiceRequest

	if err := c.BindJSON(&statusApproveServiceRequest); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	response, err := h.services.BalanceOperations.ApproveService(statusApproveServiceRequest, c)
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

func (h *Handler) failedService(c *gin.Context) {
	var statusFailedServiceRequest models.StatusServiceRequest

	if err := c.BindJSON(&statusFailedServiceRequest); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	response, err := h.services.BalanceOperations.FailedService(statusFailedServiceRequest, c)
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
