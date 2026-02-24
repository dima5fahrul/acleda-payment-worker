package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"payment-airpay/domain/entities"
	"payment-airpay/infrastructure/database/connectors"
	"payment-airpay/infrastructure/database/repositories"
	"payment-airpay/infrastructure/gateway/acleda"
)

type PaymentAcleda struct {
	masterDataRepo *repositories.MasterDataRepositoryYugabyteDB
	paymentRepo    *repositories.PaymentRepositoryYugabyteDB
	acledaRepo     *repositories.AcledaRepositoryYugabyteDB
	db             *connectors.YugabyteConnector
}

type CreatePaymentInput struct {
	Amount        string
	Currency      string
	PaymentMethod string
	ReferenceID   string
	Description   string
	CustomerName  string
	CustomerEmail string
	CustomerPhone string
	ReturnURL     string
	NotifyURL     string
}

type CreatePaymentOutput struct {
	Status        string `json:"status"`
	Message       string `json:"message"`
	TransactionID string `json:"transaction_id"`
	PaymentURL    string `json:"payment_url"`
	ReferenceID   string `json:"reference_id"`
	Amount        string `json:"amount"`
	Currency      string `json:"currency"`
	CreatedAt     string `json:"created_at"`
}

type PaymentStatusInput struct {
	TransactionID string
}

type PaymentStatusOutput struct {
	Status        string `json:"status"`
	Message       string `json:"message"`
	TransactionID string `json:"transaction_id"`
	ReferenceID   string `json:"reference_id"`
	Amount        string `json:"amount"`
	Currency      string `json:"currency"`
	PaymentStatus string `json:"payment_status"`
	PaidAt        string `json:"paid_at,omitempty"`
	FailureReason string `json:"failure_reason,omitempty"`
}

func NewPaymentAcleda(
	masterDataRepo *repositories.MasterDataRepositoryYugabyteDB,
	paymentRepo *repositories.PaymentRepositoryYugabyteDB,
	acledaRepo *repositories.AcledaRepositoryYugabyteDB,
	db *connectors.YugabyteConnector,
) *PaymentAcleda {
	return &PaymentAcleda{
		masterDataRepo: masterDataRepo,
		paymentRepo:    paymentRepo,
		acledaRepo:     acledaRepo,
		db:             db,
	}
}

func (s *PaymentAcleda) CreatePayment(ctx context.Context, in CreatePaymentInput) (*CreatePaymentOutput, error) {
	log.Printf("Creating Acleda payment with reference ID: %s", in.ReferenceID)

	// Validate input
	if in.Amount == "" {
		return nil, fmt.Errorf("amount is required")
	}
	if in.ReferenceID == "" {
		return nil, fmt.Errorf("reference ID is required")
	}
	if in.CustomerEmail == "" {
		return nil, fmt.Errorf("customer email is required")
	}

	// Create gateway request
	gatewayReq := acleda.CreatePaymentRequest{
		Amount:        in.Amount,
		Currency:      in.Currency,
		PaymentMethod: in.PaymentMethod,
		ReferenceID:   in.ReferenceID,
		Description:   in.Description,
		CustomerName:  in.CustomerName,
		CustomerEmail: in.CustomerEmail,
		CustomerPhone: in.CustomerPhone,
		ReturnURL:     in.ReturnURL,
		NotifyURL:     in.NotifyURL,
	}

	// Set default currency if not provided
	if gatewayReq.Currency == "" {
		gatewayReq.Currency = "USD"
	}

	// Set default payment method if not provided
	if gatewayReq.PaymentMethod == "" {
		gatewayReq.PaymentMethod = "credit_card"
	}

	// Call Acleda gateway
	gatewayResp, err := s.callAcledaGateway(gatewayReq)
	if err != nil {
		log.Printf("Failed to call Acleda gateway: %v", err)
		return nil, fmt.Errorf("failed to create payment: %w", err)
	}

	// Save to database
	err = s.savePaymentToDatabase(gatewayResp)
	if err != nil {
		log.Printf("Failed to save payment to database: %v", err)
		// Continue even if database save fails
	}

	// Prepare output
	out := &CreatePaymentOutput{
		Status:        gatewayResp.Status,
		Message:       gatewayResp.Message,
		TransactionID: gatewayResp.TransactionID,
		PaymentURL:    gatewayResp.PaymentURL,
		ReferenceID:   gatewayResp.ReferenceID,
		Amount:        gatewayResp.Amount,
		Currency:      gatewayResp.Currency,
		CreatedAt:     gatewayResp.CreatedAt,
	}

	log.Printf("Successfully created Acleda payment with transaction ID: %s", gatewayResp.TransactionID)
	return out, nil
}

// Save implements TransactionService interface
func (s *PaymentAcleda) Save(ctx context.Context, payment entities.Payment, payload map[string]interface{}) error {
	log.Printf("Saving Acleda payment to database: %s", payment.PaymentRequestID)

	// Placeholder for actual database save operation
	// This would save the payment to the database using repositories

	log.Printf("Successfully saved Acleda payment: %s", payment.PaymentRequestID)
	return nil
}

func (s *PaymentAcleda) GetPaymentStatus(ctx context.Context, in PaymentStatusInput) (*PaymentStatusOutput, error) {
	log.Printf("Getting Acleda payment status for transaction ID: %s", in.TransactionID)

	// Validate input
	if in.TransactionID == "" {
		return nil, fmt.Errorf("transaction ID is required")
	}

	// Create gateway request
	gatewayReq := acleda.PaymentStatusRequest{
		TransactionID: in.TransactionID,
	}

	// Call Acleda gateway
	gatewayResp, err := s.callAcledaStatusGateway(gatewayReq)
	if err != nil {
		log.Printf("Failed to call Acleda status gateway: %v", err)
		return nil, fmt.Errorf("failed to get payment status: %w", err)
	}

	// Update database with latest status
	err = s.updatePaymentStatusInDatabase(gatewayResp)
	if err != nil {
		log.Printf("Failed to update payment status in database: %v", err)
		// Continue even if database update fails
	}

	// Prepare output
	out := &PaymentStatusOutput{
		Status:        gatewayResp.Status,
		Message:       gatewayResp.Message,
		TransactionID: gatewayResp.TransactionID,
		ReferenceID:   gatewayResp.ReferenceID,
		Amount:        gatewayResp.Amount,
		Currency:      gatewayResp.Currency,
		PaymentStatus: gatewayResp.PaymentStatus,
		PaidAt:        gatewayResp.PaidAt,
		FailureReason: gatewayResp.FailureReason,
	}

	log.Printf("Successfully retrieved Acleda payment status for transaction ID: %s", gatewayResp.TransactionID)
	return out, nil
}

func (s *PaymentAcleda) callAcledaGateway(req acleda.CreatePaymentRequest) (*acleda.CreatePaymentResponse, error) {
	// This will be implemented with actual Acleda gateway client
	// For now, return a mock response
	return &acleda.CreatePaymentResponse{
		Status:        "success",
		Message:       "Payment created successfully",
		TransactionID: fmt.Sprintf("TXN-%d", time.Now().Unix()),
		PaymentURL:    "https://payment.acleda.com/pay/mock",
		ReferenceID:   req.ReferenceID,
		Amount:        req.Amount,
		Currency:      req.Currency,
		CreatedAt:     time.Now().Format("2006-01-02T15:04:05Z"),
	}, nil
}

func (s *PaymentAcleda) callAcledaStatusGateway(req acleda.PaymentStatusRequest) (*acleda.PaymentStatusResponse, error) {
	// This will be implemented with actual Acleda gateway client
	// For now, return a mock response
	return &acleda.PaymentStatusResponse{
		Status:        "success",
		Message:       "Payment status retrieved successfully",
		TransactionID: req.TransactionID,
		ReferenceID:   fmt.Sprintf("REF-%d", time.Now().Unix()),
		Amount:        "100.00",
		Currency:      "USD",
		PaymentStatus: "pending",
	}, nil
}

func (s *PaymentAcleda) savePaymentToDatabase(resp *acleda.CreatePaymentResponse) error {
	// Placeholder for database save operation
	log.Printf("Saving payment to database: %s", resp.TransactionID)
	return nil
}

func (s *PaymentAcleda) updatePaymentStatusInDatabase(resp *acleda.PaymentStatusResponse) error {
	// Placeholder for database update operation
	log.Printf("Updating payment status in database: %s", resp.TransactionID)
	return nil
}
