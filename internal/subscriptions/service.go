package subscriptions

import (
	"context"
	"errors"
	"time"

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
	ErrEmptyStartDate   = errors.New("start_date is required")
	ErrEmptyEndDate     = errors.New("end_date is required")
	ErrInvalidDate      = errors.New("date must be in MM-YYYY format")
	ErrInvalidPeriod    = errors.New("end_date must not be before start_date")
)

func (s *Service) List(ctx context.Context, p ListParams) (ListResult, error) {
	if p.PageSize <= 0 {
		p.PageSize = 10
	}
	if p.Page <= 0 {
		p.Page = 1
	}

	return s.r.List(ctx, p)
}

func (s *Service) Sum(ctx context.Context, p SumParams) (SumResult, error) {
	if p.StartDate == "" {
		return SumResult{}, ErrEmptyStartDate
	}
	if p.EndDate == "" {
		return SumResult{}, ErrEmptyEndDate
	}

	startDate, err := parseDate(p.StartDate)
	if err != nil {
		return SumResult{}, ErrInvalidDate
	}

	endDate, err := parseDate(p.EndDate)
	if err != nil {
		return SumResult{}, ErrInvalidDate
	}

	if endDate.Before(startDate) {
		return SumResult{}, ErrInvalidPeriod
	}

	return s.r.Sum(ctx, p)
}

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

	return s.r.Delete(ctx, ID)
}

func validateSubscription(sub Subscription) error {
	if sub.ServiceName == "" {
		return ErrEmptyServiceName
	}

	if sub.Price <= 0 {
		return ErrInvalidPrice
	}

	if sub.StartDate == "" {
		return ErrEmptyStartDate
	}

	if _, err := parseDate(sub.StartDate); err != nil {
		return ErrInvalidDate
	}

	if sub.EndDate != nil {
		if _, err := parseDate(*sub.EndDate); err != nil {
			return ErrInvalidDate
		}
	}

	return nil
}

func parseDate(value string) (time.Time, error) {
	return time.Parse("01-2006", value)
}

func formatDate(value time.Time) string {
	return value.Format("01-2006")
}
