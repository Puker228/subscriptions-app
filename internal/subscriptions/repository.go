package subscriptions

import (
	"context"
	"fmt"
	"strings"

	"github.com/Puker228/subscriptions-app/internal/db"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Repository struct {
	conn db.DBTX
	q    *db.Queries
}

func NewRepository(conn db.DBTX) *Repository {
	return &Repository{
		conn: conn,
		q:    db.New(conn),
	}
}

func (r *Repository) List(ctx context.Context, p ListParams) (ListResult, error) {
	if p.PageSize <= 0 {
		p.PageSize = 10
	}
	if p.Page <= 0 {
		p.Page = 1
	}

	sortCol := "start_date"
	switch p.Sort {
	case "service_name", "price", "start_date", "end_date":
		sortCol = p.Sort
	}
	sortOrder := "DESC"
	if strings.ToUpper(p.Order) == "ASC" {
		sortOrder = "ASC"
	}

	filters := []string{"1 = 1"}
	args := []any{}
	argPos := 1
	if p.ServiceName != "" {
		filters = append(filters, fmt.Sprintf("service_name ILIKE $%d", argPos))
		args = append(args, "%"+p.ServiceName+"%")
		argPos++
	}
	if p.UserID != uuid.Nil {
		filters = append(filters, fmt.Sprintf("user_id = $%d", argPos))
		args = append(args, uuidToPgtype(p.UserID))
		argPos++
	}
	whereClause := strings.Join(filters, " AND ")

	var total int64
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM subscriptions
		WHERE %s
	`, whereClause)
	if err := r.conn.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return ListResult{}, err
	}

	query := fmt.Sprintf(`
		SELECT id, service_name, price, user_id, start_date, end_date
		FROM subscriptions
		WHERE %s
		ORDER BY %s %s
		LIMIT $%d OFFSET $%d
	`, whereClause, sortCol, sortOrder, argPos, argPos+1)

	offset := (p.Page - 1) * p.PageSize
	args = append(args, p.PageSize, offset)
	rows, err := r.conn.Query(ctx, query, args...)
	if err != nil {
		return ListResult{}, err
	}
	defer rows.Close()

	subscriptions := make([]Subscription, 0)
	for rows.Next() {
		sub, err := scanSubscription(rows)
		if err != nil {
			return ListResult{}, err
		}
		subscriptions = append(subscriptions, sub)
	}
	if err := rows.Err(); err != nil {
		return ListResult{}, err
	}

	totalCount := int(total)
	totalPages := (totalCount + p.PageSize - 1) / p.PageSize
	if totalPages == 0 {
		totalPages = 1
	}

	return ListResult{
		Subscriptions: subscriptions,
		Total:         totalCount,
		Page:          p.Page,
		PageSize:      p.PageSize,
		TotalPages:    totalPages,
		HasPrev:       p.Page > 1,
		HasNext:       p.Page < totalPages,
		PrevPage:      p.Page - 1,
		NextPage:      p.Page + 1,
	}, nil
}

func (r *Repository) Create(ctx context.Context, s Subscription) error {
	endDate := pgtype.Date{}
	if s.EndDate != nil {
		parsedEndDate, err := parseDate(*s.EndDate)
		if err != nil {
			return err
		}

		endDate = pgtype.Date{
			Time:  parsedEndDate,
			Valid: true,
		}
	}

	startDate, err := parseDate(s.StartDate)
	if err != nil {
		return err
	}

	err = r.q.CreateSubscription(ctx, db.CreateSubscriptionParams{
		ID:          uuidToPgtype(s.ID),
		ServiceName: s.ServiceName,
		Price:       int32(s.Price),
		UserID:      uuidToPgtype(s.UserID),
		StartDate: pgtype.Timestamptz{
			Time:  startDate,
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
		StartDate:   formatDate(s.StartDate.Time),
	}

	if s.EndDate.Valid {
		endDate := formatDate(s.EndDate.Time)
		subscription.EndDate = &endDate
	}

	return subscription, nil
}

func (r *Repository) Update(ctx context.Context, s Subscription) error {
	endDate := pgtype.Date{}
	if s.EndDate != nil {
		parsedEndDate, err := parseDate(*s.EndDate)
		if err != nil {
			return err
		}

		endDate = pgtype.Date{
			Time:  parsedEndDate,
			Valid: true,
		}
	}

	startDate, err := parseDate(s.StartDate)
	if err != nil {
		return err
	}

	err = r.q.UpdateSubscription(ctx, db.UpdateSubscriptionParams{
		ID:          uuidToPgtype(s.ID),
		ServiceName: s.ServiceName,
		Price:       int32(s.Price),
		UserID:      uuidToPgtype(s.UserID),
		StartDate: pgtype.Timestamptz{
			Time:  startDate,
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

type subscriptionScanner interface {
	Scan(dest ...any) error
}

func scanSubscription(scanner subscriptionScanner) (Subscription, error) {
	var (
		id          pgtype.UUID
		serviceName string
		price       int32
		userID      pgtype.UUID
		startDate   pgtype.Timestamptz
		endDate     pgtype.Date
	)

	if err := scanner.Scan(&id, &serviceName, &price, &userID, &startDate, &endDate); err != nil {
		return Subscription{}, err
	}

	sub := Subscription{
		ID:          uuid.UUID(id.Bytes),
		ServiceName: serviceName,
		Price:       int(price),
		UserID:      uuid.UUID(userID.Bytes),
		StartDate:   formatDate(startDate.Time),
	}
	if endDate.Valid {
		end := formatDate(endDate.Time)
		sub.EndDate = &end
	}

	return sub, nil
}

func uuidToPgtype(id uuid.UUID) pgtype.UUID {
	return pgtype.UUID{
		Bytes: [16]byte(id),
		Valid: true,
	}
}
