package repository

import (
	"context"
	"errors"
	"subscription-service/internal/models"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

var ErrNotFound = errors.New("subscription not found")

type SubscriptionRepository struct {
	db *pgxpool.Pool
}

func NewSubsRepo(db *pgxpool.Pool) *SubscriptionRepository {
	return &SubscriptionRepository{db: db}
}

func (r *SubscriptionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(
		ctx,
		"DELETE FROM subscriptions WHERE id=$1",
		id,
	)
	return err
}

func (r *SubscriptionRepository) Create(ctx context.Context, subs *models.Subscription) error {
	query := `INSERT INTO subscriptions 
	(user_id, service_name, price, start_date, end_date)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id, created_at`

	return r.db.QueryRow(ctx, query,
		subs.UserId,
		subs.ServiceName,
		subs.Price,
		subs.StartDate,
		subs.EndDate,
	).Scan(
		&subs.ID,
		&subs.CreatedAt,
	)
}

func (r *SubscriptionRepository) Update(ctx context.Context, subs *models.Subscription) error {
	query := `UPDATE subscriptions
SET service_name = $1,
    price = $2,
    end_date = $3,
    updated_at = NOW()
WHERE id = $4
RETURNING id, service_name, price, user_id, start_date, end_date, created_at, updated_at)`

	return r.db.QueryRow(ctx, query,
		subs.ServiceName,
		subs.Price,
		subs.EndDate,
		subs.ID,
	).Scan(
		&subs.ID,
		&subs.ServiceName,
		&subs.Price,
		&subs.UserId,
		&subs.StartDate,
		&subs.EndDate,
		&subs.CreatedAt,
		&subs.UpdatedAt,
	)
}

func (r *SubscriptionRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Subscription, error) {

	query := `
	SELECT id, service_name, price, user_id, start_date, end_date, created_at, updated_at
	FROM subscriptions
	WHERE id = $1
	`

	subs := &models.Subscription{}

	err := r.db.QueryRow(ctx, query, id).Scan(
		&subs.ID,
		&subs.ServiceName,
		&subs.Price,
		&subs.UserId,
		&subs.StartDate,
		&subs.EndDate,
		&subs.CreatedAt,
		&subs.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return subs, nil
}

func (r *SubscriptionRepository) List(
	ctx context.Context,
	userID uuid.UUID,
	serviceName string,
) ([]models.Subscription, error) {

	query := `
	SELECT id, service_name, price, user_id, start_date, end_date, created_at, updated_at
	FROM subscriptions
	WHERE user_id = $1
	`

	args := []interface{}{userID}

	if serviceName != "" {
		query += " AND service_name = $2"
		args = append(args, serviceName)
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subs []models.Subscription

	for rows.Next() {
		var s models.Subscription
		if err := rows.Scan(
			&s.ID,
			&s.ServiceName,
			&s.Price,
			&s.UserId,
			&s.StartDate,
			&s.EndDate,
			&s.CreatedAt,
			&s.UpdatedAt,
		); err != nil {
			return nil, err
		}
		subs = append(subs, s)
	}

	return subs, nil
}

func (r *SubscriptionRepository) CalculateTotal(
	ctx context.Context,
	userID uuid.UUID,
	serviceName string,
	start time.Time,
	end time.Time,
) (uint, error) {

	query := `
	SELECT COALESCE(SUM(price), 0)
	FROM subscriptions
	WHERE user_id = $1
	  AND service_name = $2
	  AND start_date <= $3
	  AND (end_date IS NULL OR end_date >= $4)
	`

	var total uint

	err := r.db.QueryRow(
		ctx,
		query,
		userID,
		serviceName,
		end,
		start,
	).Scan(&total)

	if err != nil {
		return 0, err
	}

	return total, nil
}
