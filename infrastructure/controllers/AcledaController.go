package controllers

import (
	"net/http"

	"payment-airpay/application/services"
	"payment-airpay/domain/entities"
	"payment-airpay/infrastructure/common"
	"payment-airpay/infrastructure/configuration"
	"payment-airpay/infrastructure/database/repositories"
	"payment-airpay/infrastructure/gateway/acleda"
	"payment-airpay/infrastructure/service"

	"github.com/go-resty/resty/v2"
	"github.com/gofiber/fiber/v2"
)

type AcledaController struct {
	paymentLinkService *services.CreateAcledaPaymentLinkService
}

func NewAcledaController(
	gateway *acleda.AcledaGateway,
	service *service.PaymentAcleda,
	repo *repositories.PaymentAcledaRepositoryYugabyteDB,
	client *resty.Client,
) *AcledaController {
	return &AcledaController{
		paymentLinkService: services.NewCreateAcledaPaymentLinkService(gateway, service, repo, client),
	}
}

// CreatePaymentLink handles payment link creation
func (c *AcledaController) CreatePaymentLink(ctx *fiber.Ctx) error {
	incoming := ctx.Locals("incoming").(*entities.Incoming)
	incoming.Save = true
	var req services.CreateAcledaPaymentLinkInput
	if err := ctx.BodyParser(&req); err != nil {
		return common.ErrorResponse(ctx, http.StatusBadRequest, "Invalid request body", err, req, "")
	}

	if err := services.ValidateRequest(&req); err != nil {
		return common.ErrorResponse(ctx, http.StatusBadRequest, "Validation error", err, req, "")
	}

	// Create payment link
	result, err := c.paymentLinkService.Execute(ctx.Context(), req, *incoming)
	if err != nil {
		return common.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to create payment link", err, req, incoming.TransactionID)
	}

	return common.SuccessResponse(ctx, http.StatusOK, "Payment link created successfully", result, incoming.TransactionID)
}

// PaymentPage shows the Acleda payment page
func (c *AcledaController) PaymentPage(ctx *fiber.Ctx) error {
	transactionID := ctx.Params("id")
	sessionID := ctx.Query("sid")
	ptID := ctx.Query("ptid")

	if transactionID == "" {
		return common.ErrorResponse(ctx, http.StatusBadRequest, "Transaction ID is required", nil, nil, "")
	}

	// Get payment link data
	paymentLink, err := c.paymentLinkService.GetByTransactionID(ctx.Context(), transactionID)
	if err != nil {
		return common.ErrorResponse(ctx, http.StatusNotFound, "Payment link not found", err, nil, transactionID)
	}

	// Render HTML template
	return ctx.Render("payment-page-acleda", fiber.Map{
		"sid":          sessionID,
		"data":         paymentLink,
		"merchant_id":  configuration.AppConfig.AcledaMerchantID,
		"ptid":         ptID,
		"desc":         paymentLink.Description,
		"amount":       paymentLink.Amount,
		"invoice_id":   paymentLink.InvoiceID,
		"return_url":   paymentLink.ReturnURL,
		"error_url":    paymentLink.ErrorURL,
		"currency":     paymentLink.Currency,
		"expired_time": paymentLink.ExpiryTime,
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
