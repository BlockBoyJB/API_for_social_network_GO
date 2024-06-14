package pgdb

import (
	"API_for_SN_go/internal/model/pgmodel"
	"API_for_SN_go/internal/repo/pgerrs"
	"API_for_SN_go/pkg/postgres"
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	log "github.com/sirupsen/logrus"
)

const userPrefixLog = "/pgdb/user"

type UserRepo struct {
	*postgres.Postgres
}

func NewUserRepo(pg *postgres.Postgres) *UserRepo {
	return &UserRepo{pg}
}

func (r *UserRepo) CreateUser(ctx context.Context, u pgmodel.User) error {
	sql, args, _ := r.Builder.
		Insert("\"user\"").
		Columns("username, first_name", "last_name", "email", "password").
		Values(u.Username, u.FirstName, u.LastName, u.Email, u.Password).
		ToSql()
	if _, err := r.Pool.Exec(ctx, sql, args...); err != nil {
		var pgErr *pgconn.PgError
		if ok := errors.As(err, &pgErr); ok {
			if pgErr.Code == "23505" {
				return pgerrs.ErrAlreadyExists
			}
		}
		log.Errorf("%s/CreateUser error exec stmt: %s", userPrefixLog, err)
		return err
	}
	return nil

}

func (r *UserRepo) GetUserByUsername(ctx context.Context, username string) (pgmodel.User, error) {
	sql, args, _ := r.Builder.
		Select("*").
		From("\"user\"").
		Where("username = ?", username).
		ToSql()

	var user pgmodel.User

	err := r.Pool.QueryRow(ctx, sql, args...).Scan(
		&user.Id,
		&user.Username,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return pgmodel.User{}, pgerrs.ErrNotFound
		}
		log.Errorf("%s/GetUserByUsername error finding user: %s", userPrefixLog, err)
		return pgmodel.User{}, err
	}
	return user, nil
}

func (r *UserRepo) UpdateUsername(ctx context.Context, username, newUsername string) error {
	sql, args, _ := r.Builder.
		Update("\"user\"").
		Set("username", newUsername).
		Where("username = ?", username).
		ToSql()
	if _, err := r.Pool.Exec(ctx, sql, args...); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return pgerrs.ErrNotFound
		}
		log.Errorf("%s/UpdateUsername error exec stmt: %s", userPrefixLog, err)
		return err
	}
	return nil
}

func (r *UserRepo) UpdateFullName(ctx context.Context, username, firstName, lastName string) error {
	sql, args, _ := r.Builder.
		Update("\"user\"").
		Set("first_name = ?", firstName).
		Set("last_name = ?", lastName).
		Where("username = ?", username).
		ToSql()
	if _, err := r.Pool.Exec(ctx, sql, args...); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return pgerrs.ErrNotFound
		}
		log.Errorf("%s/UpdateFullName error exec stmt: %s", userPrefixLog, err)
		return err
	}
	return nil
}

func (r *UserRepo) DeleteUser(ctx context.Context, username string) error {
	sql, args, _ := r.Builder.
		Delete("\"user\"").
		Where("username = ?", username).
		ToSql()
	if _, err := r.Pool.Exec(ctx, sql, args...); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return pgerrs.ErrNotFound
		}
		log.Errorf("%s/DeleteUser error exec stmt: %s", userPrefixLog, err)
		return err
	}
	return nil
}
