package postgres

import (
	"barbot/internal/repository/gen/bot/public/model"
	"barbot/internal/repository/gen/bot/public/table"
	"database/sql"
	"github.com/go-jet/jet/v2/postgres"
	_ "github.com/lib/pq"
)

type RepositoryPostgres struct {
	DB *sql.DB
}

func New(db *sql.DB) *RepositoryPostgres {
	return &RepositoryPostgres{DB: db}
}

//func (r *RepositoryPostgres) SaveID(user int64) bool {
//	stmt := table.Users.INSERT(
//		table.Users.UserID,
//		table.Users.State,
//	).
//		VALUES(user, 0).ON_CONFLICT().DO_NOTHING()
//
//	_, err := stmt.Exec(r.DB)
//	if err != nil {
//		return false
//	}
//	return true
//}
//
//func (r *RepositoryPostgres) GetUsers() []int64 {
//	stmt := table.Users.SELECT(
//		table.Users.UserID,
//	)
//	ans := []int64{}
//	err := stmt.Query(r.DB, &ans)
//	if err != nil {
//		return nil
//	}
//	return ans
//}
//
//func (r *RepositoryPostgres) GetState(userId int64) int64 {
//	stmt := table.Users.SELECT(
//		table.Users.State,
//	).WHERE(table.Users.UserID.EQ(postgres.Int64(userId)))
//	var ans model.Users
//	err := stmt.Query(r.DB, &ans)
//	if err != nil {
//		return 0
//	}
//	return int64(pointy.Int32Value(ans.State, 0))
//}
//
//func (r *RepositoryPostgres) SetState(userId int64, state int64) {
//	stmt := table.Users.UPDATE(table.Users.State).SET(state).
//		WHERE(table.Users.UserID.EQ(postgres.Int64(userId)))
//	stmt.Exec(r.DB)
//}
//
//func (r *RepositoryPostgres) SetLogin(userId int64, login string) {
//	stmt := table.Users.UPDATE(table.Users.Login, table.Users.State).SET(login, 1).
//		WHERE(table.Users.UserID.EQ(postgres.Int64(userId)))
//	stmt.Exec(r.DB)
//}
//
//func (r *RepositoryPostgres) SaveAnswer(user, questionID int64, answer bool) {
//	stmt := table.Questions.INSERT(
//		table.Questions.UserID,
//		table.Questions.QuestionID,
//		table.Questions.Answer,
//	).
//		VALUES(user, questionID, answer)
//
//	stmt.Exec(r.DB)
//}
//
//func (r *RepositoryPostgres) NewUser(login, name string, level int) bool {
//	// todo
//	return true
//}

// --------------------------------------------------------------------------------

const (
	GuestNone = iota
	GuestAlcohol
	GuestMusic
	GuestFood
)

func (r *RepositoryPostgres) NewGuest(login, name string, level int) bool {
	stmt := table.Guests.INSERT(
		table.Guests.UserID,
		table.Guests.Login,
		table.Guests.Name,
		table.Guests.State,
		table.Guests.Level,
		table.Guests.Participation,
		table.Guests.CheckIn,
	).VALUES(0, login, name, GuestNone, level, false, false).ON_CONFLICT().DO_NOTHING()

	_, err := stmt.Exec(r.DB)
	if err != nil {
		return false
	}
	return true
}

func (r *RepositoryPostgres) GetGuests() ([]model.Guests, error) {
	stmt := table.Guests.SELECT(
		table.Guests.AllColumns,
	)
	var guests []model.Guests
	err := stmt.Query(r.DB, &guests)
	if err != nil {
		return nil, err
	}
	return guests, nil
}

func (r *RepositoryPostgres) CheckGuest(login string) (model.Guests, error) {
	stmt := table.Guests.SELECT(
		table.Guests.AllColumns).WHERE(table.Guests.Login.EQ(postgres.String(login)))
	var ans []model.Guests
	err := stmt.Query(r.DB, &ans)
	if err != nil || len(ans) != 1 {
		return model.Guests{}, err
	}
	return ans[0], err
}

func (r *RepositoryPostgres) SetID(login string, user_id int64) error {
	stmt := table.Guests.UPDATE(table.Guests.UserID).SET(user_id).
		WHERE(table.Guests.Login.EQ(postgres.String(login)))
	_, err := stmt.Exec(r.DB)
	if err != nil {
		return err
	}
	return nil
}

func (r *RepositoryPostgres) SetParticipation(userId int64, ans bool) error {
	stmt := table.Guests.UPDATE(table.Guests.Participation).SET(ans).
		WHERE(table.Guests.UserID.EQ(postgres.Int64(userId)))
	_, err := stmt.Exec(r.DB)
	if err != nil {
		return err
	}
	return nil
}

func (r *RepositoryPostgres) SetState(userId int64, state int) error {
	stmt := table.Guests.UPDATE(table.Guests.State).SET(state).
		WHERE(table.Guests.UserID.EQ(postgres.Int64(userId)))
	_, err := stmt.Exec(r.DB)
	if err != nil {
		return err
	}
	return nil
}

func (r *RepositoryPostgres) GetState(userId int64) (int, error) {
	stmt := table.Guests.SELECT(
		table.Guests.AllColumns).WHERE(table.Guests.UserID.EQ(postgres.Int64(userId)))
	var ans []model.Guests
	err := stmt.Query(r.DB, &ans)
	if err != nil || len(ans) != 1 {
		return -1, err
	}
	return int(*ans[0].State), err
}

func (r *RepositoryPostgres) DropGuest(login string) error {
	stmt := table.Guests.
		DELETE().WHERE(table.Guests.Login.EQ(postgres.String(login)))
	_, err := stmt.Exec(r.DB)
	if err != nil {
		return err
	}
	return nil
}
