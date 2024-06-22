package postgres

import (
	"barbot/internal/repository/gen/bot/public/model"
	"barbot/internal/repository/gen/bot/public/table"
	"database/sql"
	"errors"
	"github.com/go-jet/jet/v2/postgres"
	_ "github.com/lib/pq"
	"go.openly.dev/pointy"
	"strconv"
)

type RepositoryPostgres struct {
	DB *sql.DB
}

func New(db *sql.DB) *RepositoryPostgres {
	return &RepositoryPostgres{DB: db}
}

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
		table.Guests.Photo,
	).VALUES(0, login, name, GuestNone, level, false, false, "").ON_CONFLICT().DO_NOTHING()

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

func (r *RepositoryPostgres) GetWishlist() (ret []model.Wishlist, err error) {
	stmt := table.Wishlist.
		SELECT(table.Wishlist.AllColumns).ORDER_BY(table.Wishlist.ID)
	err = stmt.Query(r.DB, &ret)
	return
}

func (r *RepositoryPostgres) NewGift(description string) bool {
	tmp, err := r.GetWishlist()
	if err != nil {
		return false
	}
	idx := len(tmp) + 1
	stmt := table.Wishlist.INSERT(
		table.Wishlist.ID,
		table.Wishlist.Description,
		table.Wishlist.UserID).VALUES(idx, description, 0).ON_CONFLICT().DO_NOTHING()
	_, err = stmt.Exec(r.DB)
	if err != nil {
		return false
	}
	return true
}

func (r *RepositoryPostgres) DropGift(id int32) error {
	stmt := table.Wishlist.
		DELETE().WHERE(table.Wishlist.ID.EQ(postgres.Int32(id)))
	_, err := stmt.Exec(r.DB)
	if err != nil {
		return err
	}
	return nil
}

func (r *RepositoryPostgres) SetGiftUserID(id int32, user_id int64) (bool, error) {
	ret, _ := r.GetWishlist()
	if pointy.PointerValue(ret[id-1].UserID, 0) != 0 &&
		pointy.PointerValue(ret[id-1].UserID, 0) != user_id && user_id != 0 {
		return false, nil
	}
	stmt := table.Wishlist.UPDATE(table.Wishlist.UserID).SET(user_id).
		WHERE(table.Wishlist.ID.EQ(postgres.Int32(id)))
	_, err := stmt.Exec(r.DB)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *RepositoryPostgres) CheckIn(login string, photo string) error {
	stmt := table.Guests.UPDATE(table.Guests.Photo).SET(photo).
		WHERE(table.Guests.Login.EQ(postgres.String(login)))
	_, err := stmt.Exec(r.DB)
	return err
}

func (r *RepositoryPostgres) NewCocktail(name string, isBarmen bool, level int) bool {
	stmt := table.Cocktails.INSERT(
		table.Cocktails.Name,
		table.Cocktails.Availability,
		table.Cocktails.Level,
		table.Cocktails.Barmen,
	).VALUES(name, true, level, isBarmen).ON_CONFLICT().DO_NOTHING()

	_, err := stmt.Exec(r.DB)
	if err != nil {
		return false
	}
	return true
}

func (r *RepositoryPostgres) SetComposition(name string, composition string) bool {
	stmt := table.Cocktails.UPDATE(table.Cocktails.Composition).SET(composition).
		WHERE(table.Cocktails.Name.EQ(postgres.String(name)))
	_, err := stmt.Exec(r.DB)
	if err != nil {
		return false
	}
	return true
}

func (r *RepositoryPostgres) DelCocktail(name string) bool {
	stmt := table.Guests.
		UPDATE(table.Cocktails.Availability).SET(false).WHERE(table.Cocktails.Name.EQ(postgres.String(name)))
	_, err := stmt.Exec(r.DB)
	if err != nil {
		return false
	}
	return true
}

func (r *RepositoryPostgres) GetCocktails(level int) ([]model.Cocktails, error) {
	stmt := table.Cocktails.SELECT(
		table.Cocktails.AllColumns,
	).WHERE(table.Cocktails.Level.EQ(postgres.Int32(int32(level)))).ORDER_BY(table.Cocktails.Name)
	var cocktails []model.Cocktails
	err := stmt.Query(r.DB, &cocktails)
	if err != nil {
		return nil, err
	}
	return cocktails, nil
}

func (r *RepositoryPostgres) GetMenu(alcohol bool) (string, error) {
	stmt := table.Menu.SELECT(
		table.Menu.AllColumns,
	).WHERE(table.Menu.Alcohol.EQ(postgres.Bool(alcohol)))
	var menu []model.Menu
	err := stmt.Query(r.DB, &menu)
	if err != nil || len(menu) != 1 {
		return "", errors.New(err.Error() + "длина" + strconv.Itoa(len(menu)))
	}
	return *menu[0].Photo, nil
}

func (r *RepositoryPostgres) NewOrder(user_id int64, cocktail_id int) error {
	stmt := table.Orders.INSERT(
		table.Orders.UserID,
		table.Orders.CocktailID,
	).VALUES(user_id, cocktail_id).ON_CONFLICT().DO_NOTHING()

	_, err := stmt.Exec(r.DB)
	return err
}

func (r *RepositoryPostgres) GetGuest(id int64) (model.Guests, error) {
	stmt := table.Guests.SELECT(
		table.Guests.AllColumns).WHERE(table.Guests.UserID.EQ(postgres.Int64(id)))
	var ans []model.Guests
	err := stmt.Query(r.DB, &ans)
	if err != nil || len(ans) != 1 {
		return model.Guests{}, err
	}
	return ans[0], err
}
