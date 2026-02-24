package acleda

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"payment-airpay/domain/entities"
	"payment-airpay/infrastructure/configuration"
	"strconv"
	"time"
)

type AcledaGateway struct {
	baseURL    string
	apiKey     string
	merchantID string
	login      string
	password   string
	httpClient *http.Client
}

type CreatePaymentRequest struct {
	Amount        string `json:"amount"`
	Currency      string `json:"currency"`
	PaymentMethod string `json:"payment_method"`
	ReferenceID   string `json:"reference_id"`
	Description   string `json:"description"`
	CustomerName  string `json:"customer_name"`
	CustomerEmail string `json:"customer_email"`
	CustomerPhone string `json:"customer_phone"`
	ReturnURL     string `json:"return_url"`
	NotifyURL     string `json:"notify_url"`
}

type CreatePaymentResponse struct {
	Status        string `json:"status"`
	Message       string `json:"message"`
	TransactionID string `json:"transaction_id"`
	PaymentURL    string `json:"payment_url"`
	ReferenceID   string `json:"reference_id"`
	Amount        string `json:"amount"`
	Currency      string `json:"currency"`
	CreatedAt     string `json:"created_at"`
}

type PaymentStatusRequest struct {
	TransactionID string `json:"transaction_id"`
}

type PaymentStatusResponse struct {
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

func NewAcledaGateway() *AcledaGateway {
	return &AcledaGateway{
		baseURL:    configuration.AppConfig.AcledaAPIURL,
		apiKey:     configuration.AppConfig.AcledaAPIKey,
		merchantID: configuration.AppConfig.AcledaMerchantID,
		login:      configuration.AppConfig.AcledaLogin,
		password:   configuration.AppConfig.AcledaPassword,
		httpClient: &http.Client{
			Timeout: time.Duration(configuration.AppConfig.AcledaTimeout) * time.Millisecond,
		},
	}
}

func (g *AcledaGateway) CreatePayment(ctx context.Context, req CreatePaymentRequest) (*CreatePaymentResponse, error) {
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, g.baseURL+"/payments", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+g.apiKey)

	resp, err := g.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("acleda api error: status %d, body: %s", resp.StatusCode, string(body))
	}

	var response CreatePaymentResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, nil
}

func (g *AcledaGateway) GetPaymentStatus(ctx context.Context, req PaymentStatusRequest) (*PaymentStatusResponse, error) {
	url := fmt.Sprintf("%s/payments/%s/status", g.baseURL, req.TransactionID)

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+g.apiKey)

	resp, err := g.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("acleda api error: status %d, body: %s", resp.StatusCode, string(body))
	}

	var response PaymentStatusResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, nil
}

// OpenSession implements Acleda session opening
func (g *AcledaGateway) OpenSession(ctx context.Context, req OpenSessionRequest) (*OpenSessionResponse, error) {
	// Create request with credentials
	sessionReq := OpenSessionRequest{
		LoginID:    g.login,
		Password:   g.password,
		MerchantID: g.merchantID,
		Signature:  generateSignature(req),
		TxID:       req.TxID,
		XPayTransaction: XPayTransaction{
			TxID:             req.TxID,
			PurchaseAmount:   req.XPayTransaction.PurchaseAmount,
			PurchaseCurrency: req.XPayTransaction.PurchaseCurrency,
			PurchaseDate:     req.XPayTransaction.PurchaseDate,
			PurchaseDesc:     req.XPayTransaction.PurchaseDesc,
			InvoiceID:        req.XPayTransaction.InvoiceID,
			Item:             req.XPayTransaction.Item,
			Quantity:         req.XPayTransaction.Quantity,
			ExpiryTime:       req.XPayTransaction.ExpiryTime,
		},
	}

	jsonData, err := json.Marshal(sessionReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, g.baseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := g.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("acleda api error: status %d, body: %s", resp.StatusCode, string(body))
	}

	var response OpenSessionResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	fmt.Println(response.Result.XTran.PaymentTokenID)
	fmt.Println(response.Result.SessionID)

	return &response, nil
}

func generateSignature(req OpenSessionRequest) string {
	// Use the secret key from configuration
	return configuration.AppConfig.AcledaAPIKey
}

// Request structures for OpenSession
type OpenSessionRequest struct {
	LoginID         string          `json:"loginId"`
	Password        string          `json:"password"`
	MerchantID      string          `json:"merchantID"`
	Signature       string          `json:"signature"`
	TxID            string          `json:"txId"`
	XPayTransaction XPayTransaction `json:"xpayTransaction"`
}

type XPayTransaction struct {
	TxID             string `json:"txid"`
	PurchaseAmount   string `json:"purchaseAmount"`
	PurchaseCurrency string `json:"purchaseCurrency"`
	PurchaseDate     string `json:"purchaseDate"`
	PurchaseDesc     string `json:"purchaseDesc"`
	InvoiceID        string `json:"invoiceid"`
	Item             string `json:"item"`
	Quantity         string `json:"quantity"`
	ExpiryTime       string `json:"expiryTime"`
}

type OpenSessionResponse struct {
	Result ResultDTO `json:"result"`
}

type ResultDTO struct {
	Code         int      `json:"code"`
	ErrorDetails string   `json:"errorDetails"`
	SessionID    string   `json:"sessionid"`
	XTran        XTranDTO `json:"xTran"`
	TxDirection  int      `json:"TxDirection"`
}

type XTranDTO struct {
	PurchaseAmount float64 `json:"purchaseAmount"`
	PurchaseDate   int64   `json:"purchaseDate"`
	Quantity       int     `json:"quantity"`
	PaymentTokenID string  `json:"paymentTokenid"`
	ExpiryTime     int     `json:"expiryTime"`
	ConfirmDate    int64   `json:"confirmDate"`
	PurchaseType   int     `json:"purchaseType"`
	SaveToken      int     `json:"savetoken"`
	FeeAmount      float64 `json:"feeAmount"`
}

// Create implements PaymentGateway interface
func (g *AcledaGateway) Create(ctx context.Context, payload map[string]interface{}) (entities.Payment, error) {
	// Extract required fields from payload
	amount, ok := payload["amount"].(string)
	if !ok {
		return entities.Payment{}, fmt.Errorf("amount is required and must be string")
	}

	referenceID, ok := payload["reference_id"].(string)
	if !ok {
		return entities.Payment{}, fmt.Errorf("reference_id is required and must be string")
	}

	// Build Acleda request
	req := CreatePaymentRequest{
		Amount:        amount,
		ReferenceID:   referenceID,
		Currency:      getString(payload, "currency", "USD"),
		Description:   getString(payload, "description", ""),
		CustomerName:  getString(payload, "customer_name", ""),
		CustomerEmail: getString(payload, "customer_email", ""),
		CustomerPhone: getString(payload, "customer_phone", ""),
		ReturnURL:     getString(payload, "return_url", ""),
		NotifyURL:     getString(payload, "notify_url", ""),
		PaymentMethod: getString(payload, "payment_method", "credit_card"),
	}

	// Call Acleda API
	resp, err := g.CreatePayment(ctx, req)
	if err != nil {
		return entities.Payment{}, fmt.Errorf("failed to create Acleda payment: %w", err)
	}

	// Convert to entities.Payment
	payment := entities.Payment{
		BusinessID:       getString(payload, "business_id", ""),
		ReferenceID:      resp.ReferenceID,
		PaymentRequestID: resp.TransactionID,
		Type:             getString(payload, "type", "PAYMENT"),
		Country:          getString(payload, "country", "KH"),
		Currency:         resp.Currency,
		RequestAmount:    parseFloat(resp.Amount),
		CaptureMethod:    getString(payload, "capture_method", "FULL_CAPTURE"),
		ChannelCode:      "ACLEDA",
		ChannelProps: map[string]interface{}{
			"payment_url": resp.PaymentURL,
			"status":      resp.Status,
		},
		Actions: []entities.PaymentAction{
			{
				Type:       "PAYMENT_URL",
				Value:      resp.PaymentURL,
				Descriptor: "Complete payment using Acleda",
			},
		},
		Status:      entities.PaymentStatus(resp.Status),
		Description: getString(payload, "description", ""),
		Metadata:    payload,
		Created:     resp.CreatedAt,
		Updated:     resp.CreatedAt,
	}

	return payment, nil
}

// Helper functions
func getString(m map[string]interface{}, key, defaultValue string) string {
	if val, ok := m[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return defaultValue
}

func parseFloat(s string) float64 {
	if s == "" {
		return 0
	}
	// Simple float parsing - in production, use strconv.ParseFloat with error handling
	if f, err := strconv.ParseFloat(s, 64); err == nil {
		return f
	}
	return 0
}
