package services

import (
	"context"
	"fmt"
	"log"

	"payment-airpay/domain/entities"
)

type CreatePaymentService struct {
	Gateway PaymentGateway
	TxSvc   TransactionService
}

func NewCreatePaymentService(g PaymentGateway, t TransactionService) *CreatePaymentService {
	return &CreatePaymentService{Gateway: g, TxSvc: t}
}

func (s *CreatePaymentService) Execute(ctx context.Context, payload map[string]interface{}) (entities.Payment, error) {
	log.Printf("[CreatePaymentService.Execute] Starting payment creation, channel_code: %v", payload["channel_code"])

	res, err := s.Gateway.Create(ctx, payload)
	if err != nil {
		log.Printf("[CreatePaymentService.Execute] ERROR creating payment via gateway: %v", err)
		return entities.Payment{}, err
	}

	log.Printf("[CreatePaymentService.Execute] Payment created successfully, payment_request_id: %s, channel_code: %s", res.PaymentRequestID, res.ChannelCode)
	log.Printf("[CreatePaymentService.Execute] Calling TxSvc.Save to persist to database...")

	if err := s.TxSvc.Save(ctx, res, payload); err != nil {
		log.Printf("[CreatePaymentService.Execute] ERROR persisting payment to database: %v", err)
		return entities.Payment{}, fmt.Errorf("failed to persist payment: %w", err)
	}

	log.Printf("[CreatePaymentService.Execute] Payment persisted successfully to database")
	return res, nil
}
