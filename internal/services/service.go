package services

import (
	"context"
	"errors"
	"fmt"
	"subscription-service/internal/models"
	"subscription-service/internal/repository"
	"time"

	"github.com/google/uuid"
)

type SubsService struct {
	repo *repository.SubscriptionRepository
}

func NewSubsService(repo *repository.SubscriptionRepository) *SubsService {
	return &SubsService{
		repo: repo,
	}
}

var ErrInvalidId = errors.New("Invalid id")

func (s *SubsService) Delete(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return ErrInvalidId
	}
	err := s.repo.Delete(ctx, id)
	if errors.Is(err, repository.ErrNotFound) {
		return ErrNotFound
	}
	return err
}

func (s *SubsService) Create(ctx context.Context,
	UserID uuid.UUID,
	serviceName string,
	price uint,
	startDate time.Time,
	endDate *time.Time) (*models.Subscription, error) {

	if UserID == uuid.Nil {
		return nil, errors.New("Invalid user id")
	}

	if price <= 0 {
		return nil, errors.New("Price must be greater than 0")
	}

	subs := &models.Subscription{
		UserId:      UserID,
		ServiceName: serviceName,
		Price:       price,
		StartDate:   startDate,
		EndDate:     endDate,
	}

	err := s.repo.Create(ctx, subs)
	if err != nil {
		return nil, fmt.Errorf("failed to create subscription: %w", err)
	}

	return subs, nil

}

func (s *SubsService) Update(ctx context.Context,
	id uuid.UUID,
	serviceName string,
	price uint,
	endDate *time.Time) (*models.Subscription, error) {
	if id == uuid.Nil {
		return nil, errors.New("Invalid subscription  id")
	}
	if price <= 0 {
		return nil, errors.New("Price must be greater than 0")
	}
	subs := &models.Subscription{
		ID:          id,
		ServiceName: serviceName,
		Price:       price,
		EndDate:     endDate,
	}
	err := s.repo.Update(ctx, subs)
	if err != nil {
		return nil, fmt.Errorf("Failed to update subscription: %w", err)
	}
	return subs, nil
}

func (s *SubsService) GetByID(ctx context.Context, id uuid.UUID) (*models.Subscription, error) {

	if id == uuid.Nil {
		return nil, errors.New("invalid subscription id")
	}

	return s.repo.GetByID(ctx, id)
}

func (s *SubsService) List(
	ctx context.Context,
	userID uuid.UUID,
	serviceName string,
) ([]models.Subscription, error) {

	if userID == uuid.Nil {
		return nil, errors.New("invalid user id")
	}

	return s.repo.List(ctx, userID, serviceName)
}

func (s *SubsService) CalculateTotal(
	ctx context.Context,
	userID uuid.UUID,
	serviceName string,
	start time.Time,
	end time.Time,
) (uint, error) {

	if userID == uuid.Nil {
		return 0, errors.New("invalid user id")
	}

	if start.After(end) {
		return 0, errors.New("start date must be before end date")
	}

	return s.repo.CalculateTotal(ctx, userID, serviceName, start, end)
}
