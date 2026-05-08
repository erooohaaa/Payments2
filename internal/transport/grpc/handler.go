package grpc

import (
	"context"
	"errors"
	"strconv" // Добавили для конвертации чисел в строки

	"Payments/internal/domain"
	"Payments/internal/usecase"

	api "github.com/erooohaaa/orders-generated"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PaymentGRPCHandler struct {
	api.UnimplementedPaymentServiceServer
	uc *usecase.PaymentUseCase
}

func NewPaymentGRPCHandler(uc *usecase.PaymentUseCase) *PaymentGRPCHandler {
	return &PaymentGRPCHandler{uc: uc}
}

func (h *PaymentGRPCHandler) ProcessPayment(ctx context.Context, req *api.PaymentRequest) (*api.PaymentResponse, error) {
	if req.OrderId == "" {
		return nil, status.Error(codes.InvalidArgument, "order_id is required")
	}
	if req.Amount <= 0 {
		return nil, status.Error(codes.InvalidArgument, "amount must be positive")
	}

	payment, err := h.uc.Authorize(ctx, req.OrderId, req.Amount)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "order not found")
		}
		return nil, status.Error(codes.Internal, "failed to process payment")
	}

	return &api.PaymentResponse{
		TransactionId: payment.TransactionID,
		Status:        payment.Status,
		Amount:        strconv.FormatInt(payment.Amount, 10), // Добавили заполнение Amount
	}, nil
}

func (h *PaymentGRPCHandler) ListPayments(ctx context.Context, req *api.ListPaymentsRequest) (*api.ListPaymentsResponse, error) {
	payments, err := h.uc.ListPayments(ctx, req.MinAmount, req.MaxAmount)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	var respPayments []*api.PaymentResponse
	for _, p := range payments {
		respPayments = append(respPayments, &api.PaymentResponse{
			TransactionId: p.TransactionID,
			Status:        p.Status,
			Amount:        strconv.FormatInt(p.Amount, 10), // ИСПРАВЛЕНО: конвертируем int64 в string
		})
	}

	return &api.ListPaymentsResponse{
		Payments: respPayments,
	}, nil
}
