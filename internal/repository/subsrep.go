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
	stmt, err := sr.subs.Prepare("SELECT * FROM usr_subscription WHERE id = $1;")
	defer stmt.Close()
	if err != nil {
		return app.UserSubscription{}, app.Wrap(app.ErrDatabaseProblems, err)
	}

	result, err := stmt.QueryContext(ctx, subId)
	defer result.Close()

	if err != nil {
		return app.UserSubscription{}, app.Wrap(app.ErrDatabaseProblems, err)
	}

	var sub usrSubscription
	result.Next()
	if err = result.Scan(&sub.Id, &sub.Service, &sub.UsrId, &sub.Price, &sub.StartDate, &sub.EndDate); err != nil {
		return app.UserSubscription{}, app.Wrap(app.ErrUserSubscriptionDoesNotExist, err)
	}

	return unmarshalUserSubscription(sub), nil
}

func (sr *SubsRepository) GetSubsByUserId(ctx context.Context, userId string, page, limit int) ([]app.UserSubscription, error) {
	stmt, err := sr.subs.Prepare("SELECT * FROM usr_subscription WHERE usr_id = $1 LIMIT $2 OFFSET $3")
	defer stmt.Close()
	if err != nil {
		return []app.UserSubscription{}, app.Wrap(app.ErrDatabaseProblems, err)
	}
	offset := (page - 1) * limit
	result, err := stmt.QueryContext(ctx, userId, limit, offset)
	defer result.Close()
	if err != nil {
		return []app.UserSubscription{}, app.Wrap(app.ErrDatabaseProblems, err)
	}
	subsRes := make([]app.UserSubscription, 0)
	for result.Next() {
		var sub usrSubscription
		result.Scan(&sub.Id, &sub.Service, &sub.UsrId, &sub.Price, &sub.StartDate, &sub.EndDate)
		subsRes = append(subsRes, unmarshalUserSubscription(sub))
	}
	if len(subsRes) == 0 {
		return []app.UserSubscription{}, app.ErrUserDoesNotExistOrOutOfPage
	}

	return subsRes, nil
}

func (sr *SubsRepository) GetSubsByUserIdServiceNameStarDateEndDate(ctx context.Context, params query.SubsFilterParams) (uint, error) {
	var totalPrice uint
	stmt, err := sr.subs.Prepare("WITH pr AS (SELECT CASE WHEN start_date <= $3 THEN $3 ELSE start_date END AS start_period, CASE WHEN end_date IS NULL OR end_date >= $4 THEN $4 ELSE end_date END AS end_period FROM usr_subscription WHERE usr_id=$1 AND service=$2) SELECT SUM(price*(DATE_PART('year', AGE(pr.end_period, pr.start_period))*12+DATE_PART('month', AGE(pr.end_period, pr.start_period)))) FROM pr, usr_subscription;")
	defer stmt.Close()
	if err != nil {
		return totalPrice, app.Wrap(app.ErrDatabaseProblems, err)
	}

	result, err := stmt.QueryContext(ctx, params.UserId, params.Service, params.StartDate, params.EndDate)
	defer result.Close()
	if err != nil {
		return totalPrice, app.Wrap(app.ErrDatabaseProblems, err)
	}
	result.Next()
	result.Scan(&totalPrice)
	return totalPrice, nil
}

func (sr *SubsRepository) AddSub(ctx context.Context, sub subs.UsrSubscription) (int, error) {
	var id int
	subDb := marshallUsrSubscription(sub)

	err := sr.subs.QueryRow("INSERT INTO usr_subscription (service, usr_id, price, start_date, end_date) VALUES ($1, $2, $3, $4, $5) RETURNING id;",
		subDb.Service, subDb.UsrId, subDb.Price, subDb.StartDate, subDb.EndDate).Scan(&id)
	if err != nil {
		return 0, app.Wrap(app.ErrDatabaseProblems, err)
	}
	return id, nil
}

func (sr *SubsRepository) isExist(ctx context.Context, id int) (bool, error) {
	stmm, err := sr.subs.Prepare("SELECT * FROM usr_subscription WHERE id = $1")
	defer stmm.Close()
	if err != nil {
		return false, err
	}

	result, err := stmm.QueryContext(ctx, id)
	defer result.Close()
	if err != nil {
		return false, err
	}

	return result.Next(), nil
}

func (sr *SubsRepository) UpdateSub(ctx context.Context, sub subs.UsrSubscription) error {
	ok, err := sr.isExist(ctx, sub.ID())
	if err != nil {
		return app.ErrDatabaseProblems
	} else if !ok {
		return app.ErrUserDoesNotExist
	}

	stmt, err := sr.subs.Prepare("UPDATE usr_subscription SET service = $1, usr_id = $2, price = $3, start_date = $4, end_date = $5 WHERE id = $6;")
	defer stmt.Close()
	if err != nil {
		return app.Wrap(app.ErrDatabaseProblems, err)
	}

	dbSub := marshallUsrSubscription(sub)
	_, err = stmt.ExecContext(ctx, dbSub.Service, dbSub.UsrId, dbSub.Price, dbSub.StartDate, dbSub.EndDate, dbSub.Id)
	if err != nil {
		return app.Wrap(app.ErrUserSubscriptionDoesNotExist, err)
	}

	return nil
}

func (sr *SubsRepository) RemoveSub(ctx context.Context, subId int) error {
	ok, err := sr.isExist(ctx, subId)
	if err != nil {
		return app.ErrDatabaseProblems
	} else if !ok {
		return app.ErrUserDoesNotExist
	}

	stmt, err := sr.subs.Prepare("DELETE FROM usr_subscription WHERE id = $1")
	defer stmt.Close()
	if err != nil {
		return app.Wrap(app.ErrDatabaseProblems, err)
	}

	_, err = stmt.ExecContext(ctx, subId)
	if err != nil {
		return app.Wrap(app.ErrUserSubscriptionDoesNotExist, err)
	}

	return nil
}
