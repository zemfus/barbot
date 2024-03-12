package postgres

import (
	"bot21/internal/repository/gen/bot/public/table"
	"database/sql"
	_ "github.com/lib/pq"
)

type RepositoryPostgres struct {
	DB *sql.DB
}

func New(db *sql.DB) *RepositoryPostgres {
	return &RepositoryPostgres{DB: db}
}

func (r *RepositoryPostgres) SaveUser(user int64) {
	stmt := table.Users.INSERT(
		table.Users.UserID,
	).
		VALUES(user).ON_CONFLICT().DO_NOTHING()

	stmt.Exec(r.DB)

}

func (r *RepositoryPostgres) GetUsers() []int64 {
	stmt := table.Users.SELECT(
		table.Users.UserID,
	)
	ans := []int64{}
	err := stmt.Query(r.DB, &ans)
	if err != nil {
		return nil
	}
	return ans
}
