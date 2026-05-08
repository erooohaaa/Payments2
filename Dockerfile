FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY Payments/go.mod Payments/go.sum ./Payments/
COPY orders-generated/ ./orders-generated/
COPY Payments/ ./Payments/
WORKDIR /app/Payments
RUN go mod edit -replace github.com/erooohaaa/orders-generated=/app/orders-generated && go mod tidy && go mod download
RUN go build -o /payment-service ./cmd/payment-service
FROM alpine:3.19
WORKDIR /app
COPY --from=builder /payment-service .
CMD ["./payment-service"]