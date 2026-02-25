package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"payment-airpay/domain/entities"
	"payment-airpay/infrastructure/configuration"
	"payment-airpay/infrastructure/database/repositories"
	"payment-airpay/infrastructure/gateway/acleda"
	"payment-airpay/infrastructure/service"

	"github.com/go-resty/resty/v2"
)

type CreateAcledaPaymentLinkService struct {
	gateway *acleda.AcledaGateway
	service *service.PaymentAcleda
	repo    *repositories.PaymentAcledaRepositoryYugabyteDB
	Client  *resty.Client
}

type CreateAcledaPaymentLinkInput struct {
	Amount        string `json:"amount" validate:"required"`
	Currency      string `json:"currency" validate:"required"`
	Description   string `json:"description" validate:"required"`
	CustomerName  string `json:"customer_name" validate:"required"`
	CustomerEmail string `json:"customer_email" validate:"required"`
	CustomerPhone string `json:"customer_phone" validate:"required"`
	ReturnURL     string `json:"return_url" validate:"required"`
	CallbackURL   string `json:"callback_url" validate:"required"`
	ExpiredTime   int    `json:"expired_time" validate:"required"`
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

func NewCreateAcledaPaymentLinkService(
	gateway *acleda.AcledaGateway,
	service *service.PaymentAcleda,
	repo *repositories.PaymentAcledaRepositoryYugabyteDB,
	client *resty.Client,
) *CreateAcledaPaymentLinkService {
	return &CreateAcledaPaymentLinkService{
		gateway: gateway,
		service: service,
		repo:    repo,
		Client:  client,
	}
}

func (s *CreateAcledaPaymentLinkService) Execute(ctx context.Context, in CreateAcledaPaymentLinkInput, incoming entities.Incoming) (*CreateAcledaPaymentLinkOutput, error) {

	// Validate input
	if in.Amount == "" {
		return nil, fmt.Errorf("amount is required")
	}
	if in.Currency == "" {
		return nil, fmt.Errorf("currency is required")
	}

	if in.ReturnURL == "" {
		return nil, fmt.Errorf("return url is required")

	}

	if in.CallbackURL == "" {
		return nil, fmt.Errorf("callback url is required")

	}

	// Generate transaction ID
	transactionID := fmt.Sprintf("ACL-%d", time.Now().Unix())

	// Step 1: Open Session with Acleda
	sessionResp, err := s.gateway.OpenSessionV2(ctx, s.Client, configuration.AppConfig.ACLEDAOPENSESSIONV2URL, acleda.OpenSessionV2RequestDto{
		LoginID:    configuration.AppConfig.AcledaLogin,
		Password:   configuration.AppConfig.AcledaRemotePassword,
		MerchantID: configuration.AppConfig.AcledaMerchantID,
		Signature:  configuration.AppConfig.AcledaSecret,
		XPayTransaction: acleda.XPayTransactionDTO{
			TxID:             transactionID,
			PurchaseAmount:   in.Amount,
			PurchaseCurrency: in.Currency,
			PurchaseDate:     time.Now().Format(time.DateOnly),
			PurchaseDesc:     in.Description,
			InvoiceID:        transactionID,
			Item:             "1",
			Quantity:         "1",
			ExpiryTime:       in.ExpiredTime,
		},
	})

	go SaveAPICall(context.Background(), &sessionResp, incoming.Merchant, err, "acleda", incoming.Path, in.CustomerPhone, incoming.Webtype, incoming.TransactionID)

	if err != nil {
		return nil, fmt.Errorf("failed to open session: %w", err)
	}

	if sessionResp.Result.ErrorDetails != "SUCCESS" {
		return nil, fmt.Errorf("session failed: %s", sessionResp.Result.ErrorDetails)
	}

	paymentLinkEntity := entities.PaymentAcledaPaymentLink{
		ID:             transactionID,
		TransactionID:  transactionID,
		SessionID:      sessionResp.Result.SessionID,
		PaymentTokenID: sessionResp.Result.XTran.PaymentTokenID,
		Description:    in.Description,
		Amount:         parseFloat(in.Amount),
		Currency:       in.Currency,
		InvoiceID:      transactionID,
		Status:         "PENDING",
		ExpiryTime:     in.ExpiredTime,
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
	paymentURL := fmt.Sprintf("%s/payment-page/acleda/%s?sid=%s&ptid=%s", configuration.AppConfig.AcledaBaseURL, transactionID, sessionResp.Result.SessionID, sessionResp.Result.XTran.PaymentTokenID)

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
