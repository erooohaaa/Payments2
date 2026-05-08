package repository

import (
	"context"
	"database/sql"

	"Payments/internal/domain"
)

type PostgresPaymentRepository struct {
	db *sql.DB
}

func NewPostgresPaymentRepository(db *sql.DB) *PostgresPaymentRepository {
	return &PostgresPaymentRepository{db: db}
}

func (r *PostgresPaymentRepository) Save(ctx context.Context, p *domain.Payment) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO payments (id, order_id, transaction_id, amount, status)
		 VALUES ($1, $2, $3, $4, $5)`,
		p.ID, p.OrderID, p.TransactionID, p.Amount, p.Status,
	)
	return err
}

func (r *PostgresPaymentRepository) FindByOrderID(ctx context.Context, orderID string) (*domain.Payment, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, order_id, transaction_id, amount, status
		 FROM payments WHERE order_id = $1`, orderID)

	p := &domain.Payment{}
	err := row.Scan(&p.ID, &p.OrderID, &p.TransactionID, &p.Amount, &p.Status)
	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	return p, err

}
func (r *PostgresPaymentRepository) FindByAmountRange(ctx context.Context, min, max int64) ([]*domain.Payment, error) {
	query := `
		SELECT id, order_id, transaction_id, amount, status 
		FROM payments 
		WHERE (amount >= $1 OR $1 = 0) 
		  AND (amount <= $2 OR $2 = 0)
	`
	rows, err := r.db.QueryContext(ctx, query, min, max)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []*domain.Payment
	for rows.Next() {
		p := &domain.Payment{}
		if err := rows.Scan(&p.ID, &p.OrderID, &p.TransactionID, &p.Amount, &p.Status); err != nil {
			return nil, err
		}
		payments = append(payments, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return payments, nil
}
