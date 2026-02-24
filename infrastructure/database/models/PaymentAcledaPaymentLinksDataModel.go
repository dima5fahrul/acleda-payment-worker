package models

import (
	"time"

	"payment-airpay/domain/entities"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PaymentAcledaPaymentLinksDataModel struct {
	ID              string    `gorm:"primaryKey;column:id;type:varchar(255)"`
	TransactionID   string    `gorm:"column:transaction_id;uniqueIndex;type:varchar(255)"`
	MerchantID      string    `gorm:"column:merchant_id;type:varchar(255)"`
	SessionID       string    `gorm:"column:session_id;type:varchar(255)"`
	PaymentTokenID  string    `gorm:"column:payment_token_id;type:varchar(255)"`
	Description     string    `gorm:"column:description;type:text"`
	Amount          float64   `gorm:"column:amount"`
	PaymentCurrency string    `gorm:"column:payment_currency;type:varchar(10)"`
	InvoiceID       string    `gorm:"column:invoice_id;type:varchar(255)"`
	Status          string    `gorm:"column:status;type:varchar(50)"`
	ExpiryTime      int       `gorm:"column:expiry_time"`
	CreatedAt       time.Time `gorm:"column:created_at"`
	UpdatedAt       time.Time `gorm:"column:updated_at"`

	// Additional fields from Acleda response
	PurchaseAmount float64 `gorm:"column:purchase_amount"`
	PurchaseDate   int64   `gorm:"column:purchase_date"`
	Quantity       int     `gorm:"column:quantity"`
	ConfirmDate    int64   `gorm:"column:confirm_date"`
	PurchaseType   int     `gorm:"column:purchase_type"`
	SaveToken      int     `gorm:"column:save_token"`
	FeeAmount      float64 `gorm:"column:fee_amount"`
	TxDirection    int     `gorm:"column:tx_direction"`

	// URLs
	ReturnURL string `gorm:"column:return_url;type:text"`
	ErrorURL  string `gorm:"column:error_url;type:text"`

	// Request/Response JSON for debugging
	RequestJSON  string `gorm:"column:request_json;type:text"`
	ResponseJSON string `gorm:"column:response_json;type:text"`
}

func (PaymentAcledaPaymentLinksDataModel) TableName() string {
	return "payment_acleda_payment_links"
}

func (p *PaymentAcledaPaymentLinksDataModel) BeforeCreate(tx *gorm.DB) error {
	if p.ID == "" {
		p.ID = uuid.New().String()
	}
	return nil
}

// Convert to entity
func (p *PaymentAcledaPaymentLinksDataModel) ToEntity() entities.PaymentAcledaPaymentLink {
	return entities.PaymentAcledaPaymentLink{
		ID:              p.ID,
		TransactionID:   p.TransactionID,
		MerchantID:      p.MerchantID,
		SessionID:       p.SessionID,
		PaymentTokenID:  p.PaymentTokenID,
		Description:     p.Description,
		Amount:          p.Amount,
		Currency:        p.PaymentCurrency,
		InvoiceID:       p.InvoiceID,
		Status:          p.Status,
		ExpiryTime:      p.ExpiryTime,
		CreatedAt:       p.CreatedAt,
		UpdatedAt:       p.UpdatedAt,
		PaymentID:       "",
		PaymentMethodID: "",
		CountryID:       "",
		MerchantCode:    p.MerchantID,
		CurrencyID:      "",
		PurchaseAmount:  p.PurchaseAmount,
		PurchaseDate:    p.PurchaseDate,
		Quantity:        p.Quantity,
		ConfirmDate:     p.ConfirmDate,
		PurchaseType:    p.PurchaseType,
		SaveToken:       p.SaveToken,
		FeeAmount:       p.FeeAmount,
		TxDirection:     p.TxDirection,
		ReturnURL:       p.ReturnURL,
		ErrorURL:        p.ErrorURL,
		RequestJSON:     p.RequestJSON,
		ResponseJSON:    p.ResponseJSON,
	}
}
