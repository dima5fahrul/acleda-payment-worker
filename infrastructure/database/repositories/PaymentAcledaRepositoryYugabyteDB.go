package repositories

import (
	"context"
	"encoding/json"

	"payment-airpay/domain/entities"
	"payment-airpay/infrastructure/database/clients"
	"payment-airpay/infrastructure/database/models"
)

type PaymentAcledaRepositoryYugabyteDB struct {
	db clients.YugabyteClient
}

func NewPaymentAcledaRepositoryYugabyteDB(db clients.YugabyteClient) *PaymentAcledaRepositoryYugabyteDB {
	return &PaymentAcledaRepositoryYugabyteDB{db: db}
}

func (r *PaymentAcledaRepositoryYugabyteDB) Create(ctx context.Context, paymentLink entities.PaymentAcledaPaymentLink) error {
	if r == nil || r.db == nil || r.db.GetDB() == nil {
		return nil
	}

	// Convert entity to model
	model := models.PaymentAcledaPaymentLinksDataModel{
		ID:              paymentLink.ID,
		TransactionID:   paymentLink.TransactionID,
		MerchantID:      paymentLink.MerchantID,
		SessionID:       paymentLink.SessionID,
		PaymentTokenID:  paymentLink.PaymentTokenID,
		Description:     paymentLink.Description,
		Amount:          paymentLink.Amount,
		PaymentCurrency: paymentLink.Currency,
		InvoiceID:       paymentLink.InvoiceID,
		Status:          paymentLink.Status,
		ExpiryTime:      paymentLink.ExpiryTime,
		CreatedAt:       paymentLink.CreatedAt,
		UpdatedAt:       paymentLink.UpdatedAt,
		PurchaseAmount:  paymentLink.PurchaseAmount,
		PurchaseDate:    paymentLink.PurchaseDate,
		Quantity:        paymentLink.Quantity,
		ConfirmDate:     paymentLink.ConfirmDate,
		PurchaseType:    paymentLink.PurchaseType,
		SaveToken:       paymentLink.SaveToken,
		FeeAmount:       paymentLink.FeeAmount,
		TxDirection:     paymentLink.TxDirection,
		ReturnURL:       paymentLink.ReturnURL,
		ErrorURL:        paymentLink.ErrorURL,
		RequestJSON:     toJSON(paymentLink),
		ResponseJSON:    toJSON(paymentLink),
	}

	return r.db.GetDB().WithContext(ctx).Create(&model).Error
}

func (r *PaymentAcledaRepositoryYugabyteDB) GetByTransactionID(ctx context.Context, transactionID string) (*entities.PaymentAcledaPaymentLink, error) {
	if r == nil || r.db == nil || r.db.GetDB() == nil {
		return nil, nil
	}

	var model models.PaymentAcledaPaymentLinksDataModel
	err := r.db.GetDB().WithContext(ctx).Where("transaction_id = ?", transactionID).First(&model).Error
	if err != nil {
		return nil, err
	}

	entity := model.ToEntity()
	return &entity, nil
}

func (r *PaymentAcledaRepositoryYugabyteDB) UpdateStatus(ctx context.Context, transactionID, status string) error {
	if r == nil || r.db == nil || r.db.GetDB() == nil {
		return nil
	}

	return r.db.GetDB().WithContext(ctx).Model(&models.PaymentAcledaPaymentLinksDataModel{}).
		Where("transaction_id = ?", transactionID).
		Update("status", status).Error
}

func toJSON(v interface{}) string {
	if bytes, err := json.Marshal(v); err == nil {
		return string(bytes)
	}
	return ""
}
