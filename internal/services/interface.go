package services

import (
	"context"
	"subscription-service/internal/models"
	"time"

	"github.com/google/uuid"
)

type ISubsServiceRepo interface {
	Create(ctx context.Context, sub *models.Subscription) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Subscription, error)
	Update(ctx context.Context, sub *models.Subscription) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, id uuid.UUID, serviceName string) ([]models.Subscription, error)
	CalculateTotal(ctx context.Context, id uuid.UUID, serviceName string, start, end time.Time) (uint, error)
}
