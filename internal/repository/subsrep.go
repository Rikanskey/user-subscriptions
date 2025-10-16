package repository

import (
	"context"
	"database/sql"
	"user-subscriptions/internal/app"
	"user-subscriptions/internal/app/query"
	"user-subscriptions/internal/domain/subs"
)

type SubsRepository struct {
	subs *sql.DB
}

func NewSubsRepository(db *sql.DB) *SubsRepository {
	return &SubsRepository{subs: db}
}

func (sr *SubsRepository) GetSub(ctx context.Context, subId int) (app.UserSubscription, error) {
	stmt, err := sr.subs.Prepare(("SELECT * FROM usr_subscriptions WHERE sub_id = $1"))
	defer stmt.Close()
	if err != nil {
		return app.UserSubscription{}, err
	}

	result, err := stmt.QueryContext(ctx, subId)
	defer result.Close()

	if err != nil {
		return app.UserSubscription{}, err
	}

	var sub usrSubscription
	result.Next()
	if err := result.Scan(&sub.Id, &sub.Service, &sub.UsrId, &sub.Price, &sub.StartDate, &sub.EndDate); err != nil {
		return app.UserSubscription{}, err
	}

	return unmarshalUserSubscription(sub), nil
}

func (sr *SubsRepository) GetSubsByUserId(ctx context.Context, userId string) ([]app.UserSubscription, error) {
	stmt, err := sr.subs.Prepare("SELECT * FROM usr_subscriptions WHERE user_id = $1")
	defer stmt.Close()
	if err != nil {
		return []app.UserSubscription{}, err
	}

	result, err := stmt.QueryContext(ctx, userId)
	defer result.Close()
	if err != nil {
		return []app.UserSubscription{}, err
	}
	subs := make([]app.UserSubscription, 0)
	for result.Next() {
		var sub usrSubscription
		result.Scan(&sub)
		subs = append(subs, unmarshalUserSubscription(sub))
	}
	return subs, nil
}

func (sr *SubsRepository) GetSubsByUserIdServiceNameStarDateEndDate(ctx context.Context, params query.SubsFilterParams) ([]app.UserSubscription, error) {
	stmt, err := sr.subs.Prepare("SELECT * FROM usr_subscriptions WHERE user_id = $1 AND service_name = $2 AND start_date >= $3 AND end_date <= $4;")
	defer stmt.Close()
	if err != nil {
		return []app.UserSubscription{}, err
	}

	result, err := stmt.QueryContext(ctx, params.UserId, params.Service, params.StartDate, params.EndDate)
	defer result.Close()
	if err != nil {
		return []app.UserSubscription{}, err
	}
	subs := make([]app.UserSubscription, 0)
	for result.Next() {
		var sub usrSubscription
		result.Scan(&sub)
		subs = append(subs, unmarshalUserSubscription(sub))
	}

	return subs, nil
}

func (sr *SubsRepository) AddSub(ctx context.Context, sub subs.UsrSubscription) (int, error) {
	stmt, err := sr.subs.Prepare("INSERT INTO usr_subscriptions (service, usr_id, price, start_date, end_date VALUES ($1, $2, $3, $4, $5)")
	defer stmt.Close()
	if err != nil {
		return 0, err
	}

	result, err := stmt.ExecContext(ctx, marshallUsrSubscription(sub))

	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (sr *SubsRepository) UpdateSub(ctx context.Context, sub subs.UsrSubscription) error {
	stmt, err := sr.subs.Prepare("UPDATE usr_subscriptions SET service = $1, usr_id = $2, price = $4, start_date = $5, end_date = $6 WHERE id = $7;")
	defer stmt.Close()
	if err != nil {
		return err
	}
	dbSub := marshallUsrSubscription(sub)
	_, err = stmt.QueryContext(ctx, dbSub.Service, dbSub.UsrId, sub.Price, dbSub.StartDate, dbSub.EndDate)

	return err
}

func (sr *SubsRepository) RemoveSub(ctx context.Context, subId int) error {
	stmt, err := sr.subs.Prepare("DELETE FROM usr_subscriptions WHERE sub_id = $1")
	defer stmt.Close()
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx, subId)

	return err
}
