package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"payment-airpay/domain/entities"
	"payment-airpay/infrastructure/database/repositories"
	"payment-airpay/infrastructure/gateway/acleda"
	"payment-airpay/infrastructure/service"
)

type CreateAcledaPaymentLinkService struct {
	gateway *acleda.AcledaGateway
	service *service.PaymentAcleda
	repo    *repositories.PaymentAcledaRepositoryYugabyteDB
}

type CreateAcledaPaymentLinkInput struct {
	Amount        string `json:"amount" validate:"required"`
	Currency      string `json:"currency" validate:"required"`
	Description   string `json:"description"`
	CustomerName  string `json:"customer_name"`
	CustomerEmail string `json:"customer_email"`
	CustomerPhone string `json:"customer_phone"`
	ReturnURL     string `json:"return_url"`
	CallbackURL   string `json:"callback_url"`
	Merchant      string `json:"merchant" validate:"required"`
}

type CreateAcledaPaymentLinkOutput struct {
	TransactionID  string `json:"transaction_id"`
	PaymentURL     string `json:"payment_url"`
	SessionID      string `json:"session_id"`
	PaymentTokenID string `json:"payment_token_id"`
	Amount         string `json:"amount"`
	Currency       string `json:"currency"`
	Status         string `json:"status"`
	ExpiresAt      string `json:"expires_at"`
	CreatedAt      string `json:"created_at"`
}

func NewCreateAcledaPaymentLinkService(gateway *acleda.AcledaGateway, service *service.PaymentAcleda, repo *repositories.PaymentAcledaRepositoryYugabyteDB) *CreateAcledaPaymentLinkService {
	return &CreateAcledaPaymentLinkService{
		gateway: gateway,
		service: service,
		repo:    repo,
	}
}

func (s *CreateAcledaPaymentLinkService) Execute(ctx context.Context, in CreateAcledaPaymentLinkInput) (*CreateAcledaPaymentLinkOutput, error) {
	log.Printf("Creating Acleda payment link for merchant: %s", in.Merchant)

	// Validate input
	if in.Amount == "" {
		return nil, fmt.Errorf("amount is required")
	}
	if in.Currency == "" {
		return nil, fmt.Errorf("currency is required")
	}
	if in.Merchant == "" {
		return nil, fmt.Errorf("merchant is required")
	}

	// Generate transaction ID
	transactionID := fmt.Sprintf("ACL-%d", time.Now().Unix())

	// Step 1: Open Session with Acleda
	sessionResp, err := s.gateway.OpenSession(ctx, acleda.OpenSessionRequest{
		LoginID:    "acleda_login",
		Password:   "acleda_password",
		MerchantID: in.Merchant,
		Signature:  "acleda_signature",
		XPayTransaction: acleda.XPayTransaction{
			TxID:             transactionID,
			PurchaseAmount:   in.Amount,
			PurchaseCurrency: in.Currency,
			PurchaseDate:     time.Now().Format("2006-01-02"),
			PurchaseDesc:     in.Description,
			InvoiceID:        transactionID,
			Item:             "1",
			Quantity:         "1",
			ExpiryTime:       "60",
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to open session: %w", err)
	}

	if sessionResp.Result.ErrorDetails != "SUCCESS" {
		return nil, fmt.Errorf("session failed: %s", sessionResp.Result.ErrorDetails)
	}

	// Step 2: Save to database
	// Set default URLs if empty
	defaultReturnURL := "https://www.google.com/"
	defaultErrorURL := "https://www.google.com/"

	if in.ReturnURL == "" {
		in.ReturnURL = defaultReturnURL
	}
	if in.CallbackURL == "" {
		in.CallbackURL = defaultErrorURL
	}

	paymentLinkEntity := entities.PaymentAcledaPaymentLink{
		ID:             transactionID,
		TransactionID:  transactionID,
		MerchantID:     in.Merchant,
		SessionID:      sessionResp.Result.SessionID,
		PaymentTokenID: sessionResp.Result.XTran.PaymentTokenID,
		Description:    in.Description,
		Amount:         parseFloat(in.Amount),
		Currency:       in.Currency,
		InvoiceID:      transactionID,
		Status:         "PENDING",
		ExpiryTime:     60,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		PurchaseAmount: sessionResp.Result.XTran.PurchaseAmount,
		PurchaseDate:   sessionResp.Result.XTran.PurchaseDate,
		Quantity:       sessionResp.Result.XTran.Quantity,
		ConfirmDate:    sessionResp.Result.XTran.ConfirmDate,
		PurchaseType:   sessionResp.Result.XTran.PurchaseType,
		SaveToken:      sessionResp.Result.XTran.SaveToken,
		FeeAmount:      sessionResp.Result.XTran.FeeAmount,
		TxDirection:    sessionResp.Result.TxDirection,
		ReturnURL:      in.ReturnURL,
		ErrorURL:       in.CallbackURL,
		RequestJSON:    toJSON(sessionResp),
		ResponseJSON:   toJSON(sessionResp),
	}

	err = s.repo.Create(ctx, paymentLinkEntity)
	if err != nil {
		log.Printf("Failed to save payment link to database: %v", err)
		return nil, fmt.Errorf("failed to save payment link: %w", err)
	}

	// Step 3: Save to payments table (using existing logic)
	err = s.saveToPaymentsTable(ctx, in, transactionID)
	if err != nil {
		log.Printf("Failed to save to payments table: %v", err)
		// Continue even if payments table save fails
	}

	// Step 4: Generate payment URL
	paymentURL := fmt.Sprintf("http://localhost:8080/payment-page/acleda/%s?sid=%s&ptid=%s",
		transactionID, sessionResp.Result.SessionID, sessionResp.Result.XTran.PaymentTokenID)

	fmt.Println(paymentURL)

	// Step 5: Return response
	out := &CreateAcledaPaymentLinkOutput{
		TransactionID:  transactionID,
		PaymentURL:     paymentURL,
		SessionID:      sessionResp.Result.SessionID,
		PaymentTokenID: sessionResp.Result.XTran.PaymentTokenID,
		Amount:         in.Amount,
		Currency:       in.Currency,
		Status:         "PENDING",
		ExpiresAt:      time.Now().Add(60 * time.Minute).Format(time.RFC3339),
		CreatedAt:      time.Now().Format(time.RFC3339),
	}

	log.Printf("Successfully created Acleda payment link: %s", transactionID)
	return out, nil
}

func (s *CreateAcledaPaymentLinkService) GetByTransactionID(ctx context.Context, transactionID string) (*entities.PaymentAcledaPaymentLink, error) {
	return s.repo.GetByTransactionID(ctx, transactionID)
}

func (s *CreateAcledaPaymentLinkService) saveToPaymentsTable(ctx context.Context, in CreateAcledaPaymentLinkInput, transactionID string) error {
	// This would use the existing payment repository to save to payments table
	// For now, return nil as placeholder
	log.Printf("Saving to payments table: %s", transactionID)
	return nil
}

func parseFloat(s string) float64 {
	if s == "" {
		return 0
	}
	if f, err := strconv.ParseFloat(s, 64); err == nil {
		return f
	}
	return 0
}

func toJSON(v interface{}) string {
	if bytes, err := json.Marshal(v); err == nil {
		return string(bytes)
	}
	return ""
}
