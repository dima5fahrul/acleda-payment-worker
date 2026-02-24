package controllers

import (
	"fmt"
	"net/http"

	"payment-airpay/application/services"
	"payment-airpay/infrastructure/database/repositories"
	"payment-airpay/infrastructure/gateway/acleda"
	"payment-airpay/infrastructure/service"

	"github.com/gofiber/fiber/v2"
)

type AcledaController struct {
	paymentLinkService *services.CreateAcledaPaymentLinkService
}

func NewAcledaController(
	gateway *acleda.AcledaGateway,
	service *service.PaymentAcleda,
	repo *repositories.PaymentAcledaRepositoryYugabyteDB,
) *AcledaController {
	return &AcledaController{
		paymentLinkService: services.NewCreateAcledaPaymentLinkService(gateway, service, repo),
	}
}

// CreatePaymentLink handles payment link creation
func (c *AcledaController) CreatePaymentLink(ctx *fiber.Ctx) error {
	var req services.CreateAcledaPaymentLinkInput
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
	}

	// Validate required fields
	if req.Amount == "" {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Amount is required",
		})
	}
	if req.Currency == "" {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Currency is required",
		})
	}
	if req.Merchant == "" {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Merchant is required",
		})
	}

	// Create payment link
	result, err := c.paymentLinkService.Execute(ctx.Context(), req)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to create payment link",
			"details": err.Error(),
		})
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    result,
	})
}

// PaymentPage shows the Acleda payment page
func (c *AcledaController) PaymentPage(ctx *fiber.Ctx) error {
	transactionID := ctx.Params("id")
	sessionID := ctx.Query("sid")
	ptID := ctx.Query("ptid")

	fmt.Println(ptID)

	if transactionID == "" {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Transaction ID is required",
		})
	}

	// Get payment link data
	paymentLink, err := c.paymentLinkService.GetByTransactionID(ctx.Context(), transactionID)
	if err != nil {
		return ctx.Status(http.StatusNotFound).JSON(fiber.Map{
			"error":   "Payment link not found",
			"details": err.Error(),
		})
	}

	// Render HTML template
	return ctx.Render("payment-page-acleda", fiber.Map{
		"sid":         sessionID,
		"data":        paymentLink,
		"merchant_id": paymentLink.MerchantID,
		"ptid":        ptID,
		"desc":        paymentLink.Description,
		"amount":      paymentLink.Amount,
		"invoice_id":  paymentLink.InvoiceID,
		"return_url":  paymentLink.ReturnURL,
		"error_url":   paymentLink.ErrorURL,
		"currency":    paymentLink.Currency,
	})
}

// GetPaymentStatus retrieves payment status
func (c *AcledaController) GetPaymentStatus(ctx *fiber.Ctx) error {
	transactionID := ctx.Params("id")

	if transactionID == "" {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Transaction ID is required",
		})
	}

	paymentLink, err := c.paymentLinkService.GetByTransactionID(ctx.Context(), transactionID)
	if err != nil {
		return ctx.Status(http.StatusNotFound).JSON(fiber.Map{
			"error":   "Payment link not found",
			"details": err.Error(),
		})
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    paymentLink,
	})
}
