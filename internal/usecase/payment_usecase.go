package usecase

import (
	"context"
	"errors"
	"fmt"
	"log"

	"Payments/internal/domain"
	"Payments/internal/messaging"

	"github.com/google/uuid"
)

type PaymentRepository interface {
	Save(ctx context.Context, p *domain.Payment) error
	FindByOrderID(ctx context.Context, orderID string) (*domain.Payment, error)
	FindByAmountRange(ctx context.Context, min, max int64) ([]*domain.Payment, error)
}

type PaymentUseCase struct {
	repo      PaymentRepository
	publisher messaging.EventPublisher
}

func NewPaymentUseCase(repo PaymentRepository, publisher messaging.EventPublisher) *PaymentUseCase {
	return &PaymentUseCase{repo: repo, publisher: publisher}
}

func (uc *PaymentUseCase) Authorize(ctx context.Context, orderID string, amount int64) (*domain.Payment, error) {
	payStatus := domain.StatusAuthorized
	if amount > domain.MaxPaymentAmount {
		payStatus = domain.StatusDeclined
	}

	p := &domain.Payment{
		ID:            uuid.New().String(),
		OrderID:       orderID,
		TransactionID: uuid.New().String(),
		Amount:        amount,
		Status:        payStatus,
	}

	if err := uc.repo.Save(ctx, p); err != nil {
		return nil, err
	}

	if payStatus == domain.StatusAuthorized {
		event := messaging.PaymentEvent{
			MessageID:     uuid.New().String(),
			OrderID:       orderID,
			Amount:        float64(amount),
			CustomerEmail: fmt.Sprintf("customer-%s@example.com", orderID[:8]),
			Status:        payStatus,
		}
		if err := uc.publisher.PublishPaymentCompleted(ctx, event); err != nil {
			log.Printf("[Payment] WARNING: failed to publish event: %v", err)
		}
	}

	return p, nil
}

func (uc *PaymentUseCase) GetByOrderID(ctx context.Context, orderID string) (*domain.Payment, error) {
	return uc.repo.FindByOrderID(ctx, orderID)
}

func (uc *PaymentUseCase) ListPayments(ctx context.Context, min, max int64) ([]*domain.Payment, error) {
	if min > 0 && max > 0 && min > max {
		return nil, errors.New("min cannot be greater than max")
	}
	return uc.repo.FindByAmountRange(ctx, min, max)
}
