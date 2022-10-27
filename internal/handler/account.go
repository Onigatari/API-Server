package handler

import (
	"Avito/internal/models"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

// @Summary Ping
// @Description Ping-Pong!
// @Produce json
// @Router /ping/ [get]
func (h *Handler) ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

// @Summary Balance
// @Tags api
// @Description Метод начисления средств на баланс
// @Produce json
// @Param id path integer true "user id"
// @Success 200 {object} models.GetBalanceResponse
// @Failure 500 {object} errorMessage
// @Router /api/balance/{id} [get]
func (h *Handler) balance(c *gin.Context) {
	idStringInput := c.Param("id")
	log.Printf("[Balance] Input read: %v %T", idStringInput, idStringInput)

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

// @Summary Deposit
// @Tags api
// @Description Метод пополнение баланса
// @Produce json
// @Param input body models.UpdateBalanceRequest true "JSON с ID пользователя и суммой денег для вноса"
// @Success 200 {object} models.UpdateBalanceRequest
// @Failure 500 {object} errorMessage
// @Router /api/deposit/ [post]
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

// @Summary Withdrawal
// @Tags api
// @Description Метод снятие средств со счета
// @Produce json
// @Param input body models.UpdateBalanceRequest true "JSON с ID пользователя и суммой денег для вывода"
// @Success 200 {object} models.UpdateBalanceRequest
// @Failure 500 {object} errorMessage
// @Router /api/withdrawal/ [post]
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

// @Summary Transfer
// @Tags api
// @Description Метод перевода между пользователями
// @Produce json
// @Param input body models.TransferRequest true "JSON с ID отправителя, ID получателя и суммой денег"
// @Success 200 {object} models.TransferRequest
// @Failure 500 {object} errorMessage
// @Router /api/transfer/ [post]
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

// @Summary Reserve Service
// @Tags api
// @Description Метод перевода между пользователями
// @Produce json
// @Param input body models.ReserveServiceRequest true "JSON с ID пользователя, ID услуги, ID заказа и суммой комиссии"
// @Success 200 {object} models.ReserveServiceRequest
// @Failure 500 {object} errorMessage
// @Router /api/reserveService/ [post]
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

// @Summary Approve Order
// @Tags api
// @Description Метод признания выручки: списывает из резерва деньги, добавляет данные в отчет для бухгалтерии
// @Produce json
// @Param input body models.StatusServiceRequest true "JSON с ID пользователя, ID услуги, ID заказа и суммой комиссии"
// @Success 200 {object} models.StatusServiceRequest
// @Failure 500 {object} errorMessage
// @Router /api/approveOrder/ [post]
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

// @Summary Failed Service
// @Tags api
// @Description Метод признания выручки: списывает из резерва деньги, добавляет данные в отчет для бухгалтерии
// @Produce json
// @Param input body models.StatusServiceRequest true "JSON с ID пользователя, ID услуги, ID заказа и суммой комиссии"
// @Success 200 {object} models.StatusServiceRequest
// @Failure 500 {object} errorMessage
// @Router /api/failedService/ [post]
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
