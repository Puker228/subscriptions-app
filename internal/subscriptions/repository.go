package subscriptions

import (
	"context"

	"github.com/Puker228/subscriptions-app/internal/db"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Repository struct {
	q *db.Queries
}

func NewRepository(conn db.DBTX) *Repository {
	return &Repository{
		q: db.New(conn),
	}
}

func (r *Repository) Create(ctx context.Context, s Subscription) error {
	endDate := pgtype.Date{}
	if s.EndDate != nil {
		endDate = pgtype.Date{
			Time:  *s.EndDate,
			Valid: true,
		}
	}

	err := r.q.CreateSubscription(ctx, db.CreateSubscriptionParams{
		ID:          uuidToPgtype(s.ID),
		ServiceName: s.ServiceName,
		Price:       int32(s.Price),
		UserID:      uuidToPgtype(s.UserID),
		StartDate: pgtype.Timestamptz{
			Time:  s.StartDate,
			Valid: true,
		},
		EndDate: endDate,
	})

	return err
}

func (r *Repository) GetOneByID(ctx context.Context, sID uuid.UUID) (Subscription, error) {
	s, err := r.q.GetSubscription(ctx, uuidToPgtype(sID))
	if err != nil {
		return Subscription{}, err
	}

	subscription := Subscription{
		ID:          uuid.UUID(s.ID.Bytes),
		ServiceName: s.ServiceName,
		Price:       int(s.Price),
		UserID:      uuid.UUID(s.UserID.Bytes),
		StartDate:   s.StartDate.Time,
	}

	if s.EndDate.Valid {
		subscription.EndDate = &s.EndDate.Time
	}

	return subscription, nil
}

func (r *Repository) Update(ctx context.Context, s Subscription) error {
	endDate := pgtype.Date{}
	if s.EndDate != nil {
		endDate = pgtype.Date{
			Time:  *s.EndDate,
			Valid: true,
		}
	}

	err := r.q.UpdateSubscription(ctx, db.UpdateSubscriptionParams{
		ID:          uuidToPgtype(s.ID),
		ServiceName: s.ServiceName,
		Price:       int32(s.Price),
		UserID:      uuidToPgtype(s.UserID),
		StartDate: pgtype.Timestamptz{
			Time:  s.StartDate,
			Valid: true,
		},
		EndDate: endDate,
	})

	return err
}

func (r *Repository) Delete(ctx context.Context, sID uuid.UUID) error {
	err := r.q.DeleteSubscription(ctx, uuidToPgtype(sID))
	return err
}

func uuidToPgtype(id uuid.UUID) pgtype.UUID {
	return pgtype.UUID{
		Bytes: [16]byte(id),
		Valid: true,
	}
}
