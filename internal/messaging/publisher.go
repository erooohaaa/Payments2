package messaging

import "context"

type PaymentEvent struct {
	MessageID     string  `json:"message_id"`
	OrderID       string  `json:"order_id"`
	Amount        float64 `json:"amount"`
	CustomerEmail string  `json:"customer_email"`
	Status        string  `json:"status"`
}

type EventPublisher interface {
	PublishPaymentCompleted(ctx context.Context, event PaymentEvent) error
	Close()
}
