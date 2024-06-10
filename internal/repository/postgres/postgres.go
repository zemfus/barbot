package postgres

import (
	"barbot/internal/repository/gen/bot/public/model"
	"barbot/internal/repository/gen/bot/public/table"
	"database/sql"
	"github.com/go-jet/jet/v2/postgres"
	_ "github.com/lib/pq"
	"go.openly.dev/pointy"
)

type RepositoryPostgres struct {
	DB *sql.DB
}

func New(db *sql.DB) *RepositoryPostgres {
	return &RepositoryPostgres{DB: db}
}

func (r *RepositoryPostgres) SaveID(user int64) bool {
	stmt := table.Users.INSERT(
		table.Users.UserID,
		table.Users.State,
	).
		VALUES(user, 0).ON_CONFLICT().DO_NOTHING()

	_, err := stmt.Exec(r.DB)
	if err != nil {
		return false
	}
	return true
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

func (r *RepositoryPostgres) GetState(userId int64) int64 {
	stmt := table.Users.SELECT(
		table.Users.State,
	).WHERE(table.Users.UserID.EQ(postgres.Int64(userId)))
	var ans model.Users
	err := stmt.Query(r.DB, &ans)
	if err != nil {
		return 0
	}
	return int64(pointy.Int32Value(ans.State, 0))
}

func (r *RepositoryPostgres) SetState(userId int64, state int64) {
	stmt := table.Users.UPDATE(table.Users.State).SET(state).
		WHERE(table.Users.UserID.EQ(postgres.Int64(userId)))
	stmt.Exec(r.DB)
}

func (r *RepositoryPostgres) SetLogin(userId int64, login string) {
	stmt := table.Users.UPDATE(table.Users.Login, table.Users.State).SET(login, 1).
		WHERE(table.Users.UserID.EQ(postgres.Int64(userId)))
	stmt.Exec(r.DB)
}

func (r *RepositoryPostgres) SaveAnswer(user, questionID int64, answer bool) {
	stmt := table.Questions.INSERT(
		table.Questions.UserID,
		table.Questions.QuestionID,
		table.Questions.Answer,
	).
		VALUES(user, questionID, answer)

	stmt.Exec(r.DB)
}
