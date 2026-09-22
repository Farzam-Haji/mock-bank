package handlers

import (
	"context"
	"log"
	"errors"
	"fmt"
	"strconv"

	"github.com/Farzam-Haji/mock-bank/internal/models"
	"github.com/Farzam-Haji/mock-bank/internal/services"
	"github.com/gin-gonic/gin"
)

type BankService interface {
	CreatePayment(ctx context.Context, orderID int64, amount int64) (*models.Payment, error)
	GetPayment(ctx context.Context, paymentID int64) (*models.Payment, error)
	ProcessPayment(ctx context.Context, paymentID int64, card string, password string) (*models.Payment, error)
	CallbackClient(paymentID int64, orderID int64, status string) (error)
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

	paymentURL := fmt.Sprintf("http://localhost:8081/api/v1/payments/%v", payment.ID)

	c.JSON(201, gin.H{
		"payment_id": payment.ID,
		"payment_url": paymentURL,
	})

}

func (h *BankHandler) ShowPayment(c *gin.Context) {
	id := c.Param("id")
	paymentID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{
			"invalid param": err,
		})
		return
	}

	payment, err := h.service.GetPayment(c.Request.Context(), paymentID)
	if errors.Is(err, services.ErrInternal) {
		c.JSON(500, err.Error())
		return
	}
	if err != nil {
		c.JSON(500, err.Error())
		return
	}

	amount := fmt.Sprintf("%.2f $", float64(payment.Amount)/float64(100))

	c.HTML(200, "payment.html", gin.H{
		"payment_id": payment.ID,
		"amount": amount,
	})
}

func (h *BankHandler) ProcessPayment(c *gin.Context) {
	id := c.Param("id")
	paymentID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{
			"invalid param": err,
		})
		return
	}

	cardNumber := c.PostForm("card_number")
	password := c.PostForm("password")
	if cardNumber == "" || password == "" {
		c.JSON(400, "html issue")
		return 
	}

	resultCode := 200
	var resultErr error

	payment, err := h.service.ProcessPayment(c.Request.Context(), paymentID, cardNumber, password)

	if errors.Is(err, services.ErrInternal){
		resultCode = 500
		resultErr = services.ErrInternal
	}

	if err != nil {
		resultErr = err
	}

	err = h.service.CallbackClient(payment.ID, payment.MerchantOrderID, payment.Status)
	if err != nil {
		log.Printf(
			"payment succeeded but callback failed: payment_id=%d order_id=%d error=%v",
	        payment.ID,
	        payment.MerchantOrderID,
	        err,
		)
	}

	c.HTML(resultCode, "result.html", gin.H{
		"status": payment.Status,
		"error": resultErr,
	})
}