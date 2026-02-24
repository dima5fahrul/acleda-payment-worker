package entities

import (
	"time"
)

type PaymentAcledaPaymentLink struct {
	ID              string    `json:"id"`
	TransactionID   string    `json:"transaction_id"`
	MerchantID      string    `json:"merchant_id"`
	SessionID       string    `json:"session_id"`
	PaymentTokenID  string    `json:"payment_token_id"`
	Description     string    `json:"description"`
	Amount          float64   `json:"amount"`
	Currency        string    `json:"currency"`
	InvoiceID       string    `json:"invoice_id"`
	Status          string    `json:"status"`
	ExpiryTime      int       `json:"expiry_time"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	
	// Foreign Keys
	PaymentID       string    `json:"payment_id"`
	PaymentMethodID string    `json:"payment_method_id"`
	CountryID       string    `json:"country_id"`
	MerchantCode    string    `json:"merchant_code"`
	CurrencyID      string    `json:"currency_id"`
	
	// Additional fields from Acleda response
	PurchaseAmount   float64   `json:"purchase_amount"`
	PurchaseDate     int64     `json:"purchase_date"`
	Quantity         int       `json:"quantity"`
	ConfirmDate      int64     `json:"confirm_date"`
	PurchaseType     int       `json:"purchase_type"`
	SaveToken        int       `json:"save_token"`
	FeeAmount        float64   `json:"fee_amount"`
	TxDirection      int       `json:"tx_direction"`
	
	// URLs
	ReturnURL        string    `json:"return_url"`
	ErrorURL         string    `json:"error_url"`
	
	// Request/Response JSON for debugging
	RequestJSON      string    `json:"request_json"`
	ResponseJSON     string    `json:"response_json"`
}
