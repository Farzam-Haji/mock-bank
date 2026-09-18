package handlers

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/Farzam-Haji/mock-bank/internal/models"
)

type BankService interface {
	CreatePayment(ctx context.Context, orderID int64, amount int64) (*models.Payment, error)
}

type BankHandler struct {
	service 	BankService
}

func NewBankHanler(serviceInterface BankService) (*BankHandler) {
	return &BankHandler{
		service: serviceInterface,
	}
}

//requests
type paymentRequest struct {
	OrderID 	int64 	`json:"order_id" binding:"required,gt=0"`
	Amount 				int64 	`json:"amount" binding:"required,gt=0"`
}

func (h *BankHandler) CreatePayment(c *gin.Context) {
	var request paymentRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{
			"invalid json": err,
		})
		return
	}

	payment, err := h.service.CreatePayment(c.Request.Context(), request.OrderID, request.Amount)
	if err != nil {
		c.JSON(500, err.Error())
		return
	}

	c.JSON(201, gin.H{
		"payment_id": payment.ID,
		"amount": payment.Amount,
		"status": payment.Status,
	})

}