package subscriptions

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

type Service struct {
	r *Repository
}

func NewService(r *Repository) *Service {
	return &Service{r: r}
}

var (
	ErrEmptyServiceName = errors.New("service name is required")
	ErrInvalidPrice     = errors.New("price must be greater than zero")
	ErrInvalidUserID    = errors.New("user id is required")
)

func (s *Service) Create(ctx context.Context, sub Subscription) (Subscription, error) {
	if sub.ID == uuid.Nil {
		sub.ID = uuid.New()
	}

	if err := validateSubscription(sub); err != nil {
		return Subscription{}, err
	}

	if err := s.r.Create(ctx, sub); err != nil {
		return Subscription{}, err
	}

	return sub, nil
}

func (s *Service) GetOneByID(ctx context.Context, id uuid.UUID) (Subscription, error) {
	if id == uuid.Nil {
		return Subscription{}, ErrInvalidUserID
	}

	return s.r.GetOneByID(ctx, id)
}

func (s *Service) Update(ctx context.Context, sub Subscription) error {
	if sub.ID == uuid.Nil {
		return ErrInvalidUserID
	}

	if err := validateSubscription(sub); err != nil {
		return err
	}

	return s.r.Update(ctx, sub)
}
func (s *Service) Delete(ctx context.Context, ID uuid.UUID) error {
	if ID == uuid.Nil {
		return ErrInvalidUserID
	}

	return s.Delete(ctx, ID)
}

func validateSubscription(sub Subscription) error {
	if sub.ServiceName == "" {
		return ErrEmptyServiceName
	}

	if sub.Price <= 0 {
		return ErrInvalidPrice
	}

	return nil
}
