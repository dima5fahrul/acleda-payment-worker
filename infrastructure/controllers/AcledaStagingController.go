package controllers

import (
	"net/http"

	"payment-airpay/application/services"

	"github.com/gofiber/fiber/v2"
)

type AcledaStagingController struct {
	stagingService *services.CreateAcledaStagingPaymentService
}

func NewAcledaStagingController(stagingService *services.CreateAcledaStagingPaymentService) *AcledaStagingController {
	return &AcledaStagingController{
		stagingService: stagingService,
	}
}

// CreateStagingPayment creates a new Acleda staging payment
func (c *AcledaStagingController) CreateStagingPayment(ctx *fiber.Ctx) error {
	var input services.CreateAcledaStagingPaymentInput

	if err := ctx.BodyParser(&input); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   true,
			"message": "Invalid request format",
			"details": err.Error(),
		})
	}

	// Validate required fields
	if input.Amount == "" {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   true,
			"message": "Amount is required",
		})
	}

	if input.Msisdn == "" {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   true,
			"message": "Msisdn is required",
		})
	}

	if input.Country == "" {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   true,
			"message": "Country is required",
		})
	}

	if input.PaymentMethod == "" {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   true,
			"message": "Payment method is required",
		})
	}

	if input.Currency == "" {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   true,
			"message": "Currency is required",
		})
	}

	// Create staging payment
	result, err := c.stagingService.Execute(ctx.Context(), &input)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   true,
			"message": "Failed to create staging payment",
			"details": err.Error(),
		})
	}

	// Return success response
	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"status":  200,
		"error":   false,
		"message": "success",
		"data":    result,
	})
}

// GetStagingPaymentStatus retrieves staging payment status
func (c *AcledaStagingController) GetStagingPaymentStatus(ctx *fiber.Ctx) error {
	transactionID := ctx.Params("id")

	if transactionID == "" {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   true,
			"message": "Transaction ID is required",
		})
	}

	// Get payment status
	paymentLink, err := c.stagingService.GetByTransactionID(ctx.Context(), transactionID)
	if err != nil {
		return ctx.Status(http.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error":   true,
			"message": "Payment not found",
			"details": err.Error(),
		})
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    paymentLink,
	})
}
