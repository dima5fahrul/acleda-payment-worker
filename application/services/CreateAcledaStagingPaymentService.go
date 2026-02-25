package services

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"payment-airpay/domain/entities"
	"payment-airpay/infrastructure/gateway/acleda"
)

type CreateAcledaStagingPaymentService struct {
	gateway *acleda.AcledaGateway
}

func NewCreateAcledaStagingPaymentService(gateway *acleda.AcledaGateway) *CreateAcledaStagingPaymentService {
	return &CreateAcledaStagingPaymentService{
		gateway: gateway,
	}
}

type CreateAcledaStagingPaymentInput struct {
	Amount            string `json:"amount" validate:"required"`
	Msisdn            string `json:"msisdn" validate:"required"`
	Country           string `json:"country" validate:"required"`
	Description       string `json:"description"`
	PaymentMethod     string `json:"payment_method" validate:"required"`
	BankCode          string `json:"bank_code"`
	EwalletSuccessURL string `json:"ewallet_success_url"`
	EwalletFailureURL string `json:"ewallet_failuure_url"`
	EwalletCancelURL  string `json:"ewallet_cancel_url"`
	VACustomerName    string `json:"va_customer_name"`
	CallbackURL       string `json:"callback_url"`
	Email             string `json:"email"`
	Username          string `json:"username"`
	ReturnURL         string `json:"return_url"`
	Currency          string `json:"currency" validate:"required"`
}

type CreateAcledaStagingPaymentOutput struct {
	TransactionID string  `json:"transaction_id"`
	PaymentMethod string  `json:"payment_method"`
	Provider      string  `json:"provider"`
	Bank          string  `json:"bank"`
	PaymentLink   string  `json:"payment_link"`
	PaymentCode   string  `json:"payment_code"`
	Name          string  `json:"name"`
	Email         string  `json:"email"`
	Amount        float64 `json:"amount"`
	Currency      string  `json:"currency"`
	Status        string  `json:"status"`
	ExpiresAt     string  `json:"expires_at"`
}

func (s *CreateAcledaStagingPaymentService) Execute(ctx context.Context, in *CreateAcledaStagingPaymentInput) (*CreateAcledaStagingPaymentOutput, error) {
	log.Printf("Creating Acleda staging payment for amount: %s, msisdn: %s", in.Amount, in.Msisdn)

	// Generate transaction ID
	transactionID := fmt.Sprintf("LINKIT%d", time.Now().Unix())

	// Call Acleda staging gateway
	resp, err := s.gateway.CreateStagingPayment(ctx, &acleda.StagingPaymentRequest{
		Amount:            in.Amount,
		Msisdn:            in.Msisdn,
		Country:           in.Country,
		Description:       in.Description,
		PaymentMethod:     in.PaymentMethod,
		BankCode:          in.BankCode,
		EwalletSuccessURL: in.EwalletSuccessURL,
		EwalletFailureURL: in.EwalletFailureURL,
		EwalletCancelURL:  in.EwalletCancelURL,
		VACustomerName:    in.VACustomerName,
		CallbackURL:       in.CallbackURL,
		Email:             in.Email,
		Username:          in.Username,
		ReturnURL:         in.ReturnURL,
		Currency:          in.Currency,
		TransactionID:     transactionID,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create staging payment: %w", err)
	}

	// Parse amount
	amount, err := strconv.ParseFloat(in.Amount, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid amount format: %w", err)
	}

	// Calculate expiry time (30 minutes from now)
	expiresAt := time.Now().Add(30 * time.Minute).Format("2006-01-02 03:04:05.000000000 -0700")

	// Return response
	out := &CreateAcledaStagingPaymentOutput{
		TransactionID: resp.TransactionID,
		PaymentMethod: resp.PaymentMethod,
		Provider:      resp.Provider,
		Bank:          resp.Bank,
		PaymentLink:   resp.PaymentLink,
		PaymentCode:   resp.PaymentCode,
		Name:          resp.Name,
		Email:         resp.Email,
		Amount:        amount,
		Currency:      resp.Currency,
		Status:        resp.Status,
		ExpiresAt:     expiresAt,
	}

	log.Printf("Successfully created Acleda staging payment: %s", transactionID)
	return out, nil
}

func (s *CreateAcledaStagingPaymentService) GetByTransactionID(ctx context.Context, transactionID string) (*entities.PaymentAcledaPaymentLink, error) {
	// For now, return nil since staging service doesn't need database operations
	return nil, fmt.Errorf("staging service does not support database operations")
}
