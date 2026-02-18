package handlers

import (
	"context"
	"subscription-service/internal/models"
	"time"

	"github.com/google/uuid"
)

type ISubsServiceHandler interface {
	Create(
		ctx context.Context,
		userID uuid.UUID,
		serviceName string,
		price uint,
		startDate time.Time,
		endDate *time.Time,
	) (*models.Subscription, error)

	GetByID(ctx context.Context, id uuid.UUID) (*models.Subscription, error)

	Update(
		ctx context.Context,
		id uuid.UUID,
		serviceName string,
		price uint,
		endDate *time.Time,
	) (*models.Subscription, error)

	Delete(ctx context.Context, id uuid.UUID) error

	List(
		ctx context.Context,
		userID uuid.UUID,
		serviceName string,
	) ([]models.Subscription, error)

	CalculateTotal(
		ctx context.Context,
		userID uuid.UUID,
		serviceName string,
		start, end time.Time,
	) (uint, error)
}
