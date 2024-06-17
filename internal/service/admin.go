package service

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.openly.dev/pointy"
	"strconv"
	"strings"
)

const (
	None = iota
	NewGuest
	DelGuest
	NewGift
	DelGift
)

var state int = None

func (s *Service) handleAdmin(update tgbotapi.Update) {
	if update.Message == nil {
		return
	}

	if update.Message.Photo != nil {
		s.bots.Bot.Send(tgbotapi.NewPhoto(s.AdminID, tgbotapi.FileID(update.Message.Photo[len(update.Message.Photo)-1].FileID)))
		return
	}
	command := update.Message.Command()

	if command == "start" {
		s.bots.Bot.Send(tgbotapi.NewMessage(s.AdminID, "Команды:\n/newGuest - добавить гостя\n/guests - все приглашения\n/delGuest - удалить гостя по логину"))
		return
	}

	if command == "newGuest" {
		s.bots.Bot.Send(tgbotapi.NewMessage(s.AdminID, "введи через пробел: логин (без @), имя, уровень (0 - без алко, 1 - алко, 2 - спешл)"))
		state = NewGuest
		return
	}

	if command == "delGuest" {
		s.bots.Bot.Send(tgbotapi.NewMessage(s.AdminID, "введи логин (без @)"))
		state = DelGuest
		return
	}

	if state == NewGuest {
		state = None
		data := strings.Fields(update.Message.Text)
		if len(data) != 3 {
			s.bots.Bot.Send(tgbotapi.NewMessage(s.AdminID, "Неверное количество"))
			return
		}
		level, err := strconv.Atoi(data[2])
		if err != nil || (level < 0 || level > 2) {
			s.bots.Bot.Send(tgbotapi.NewMessage(s.AdminID, "Что-то не так с уровнем"))
			return
		}

		if s.db.NewGuest(data[0], data[1], level) {
			s.bots.Bot.Send(tgbotapi.NewMessage(s.AdminID, "Добавлено"))
		} else {
			s.bots.Bot.Send(tgbotapi.NewMessage(s.AdminID, "Какие-то траблы с добавлением((("))
		}
		return
	}

	if command == "guests" {
		guests, err := s.db.GetGuests()
		if err != nil {
			s.bots.Bot.Send(tgbotapi.NewMessage(s.AdminID, "Траблы с запросом"))
			return
		}
		str := "Guests:\n"
		for i, g := range guests {
			str += fmt.Sprint(
				strconv.Itoa(i+1), ") ",
				pointy.PointerValue(g.UserID, 0), " ",
				pointy.PointerValue(g.Login, ""), " ",
				pointy.PointerValue(g.Name, ""), " ",
				pointy.PointerValue(g.State, 0), " ",
				pointy.PointerValue(g.Level, 0), " ",
				pointy.PointerValue(g.Participation, false), " ",
				pointy.PointerValue(g.CheckIn, false), "\n",
			)
		}
		s.bots.Bot.Send(tgbotapi.NewMessage(s.AdminID, str))
		return
	}

	if state == DelGuest {
		state = None
		err := s.db.DropGuest(update.Message.Text)
		if err != nil {
			s.bots.Bot.Send(tgbotapi.NewMessage(s.AdminID, "Что-то пошло не так в удалении"))
			return
		}
		s.bots.Bot.Send(tgbotapi.NewMessage(s.AdminID, "Без ошибок прошло удаление"))
		return
	}

	if command == "newGift" {
		s.bots.Bot.Send(tgbotapi.NewMessage(s.AdminID, "Напиши что ты хочешь"))
		state = NewGift
		return
	}

	if state == NewGift {
		state = None
		if s.db.NewGift(update.Message.Text) {
			s.bots.Bot.Send(tgbotapi.NewMessage(s.AdminID, "Добавили"))
		} else {
			s.bots.Bot.Send(tgbotapi.NewMessage(s.AdminID, "Косяк с добавлением"))
		}
		return
	}

	if command == "delGift" {
		s.bots.Bot.Send(tgbotapi.NewMessage(s.AdminID, "Введи номер"))
		state = DelGift
		return
	}

	if state == DelGift {
		state = None
		id, err := strconv.Atoi(update.Message.Text)
		if err != nil {
			s.bots.Bot.Send(tgbotapi.NewMessage(s.AdminID, "Это не число"))
			return
		}
		err = s.db.DropGift(int32(id))
		if err != nil {
			s.bots.Bot.Send(tgbotapi.NewMessage(s.AdminID, "Ошибка"+err.Error()))
		} else {
			s.bots.Bot.Send(tgbotapi.NewMessage(s.AdminID, "Удалили"))
		}
		return
	}

	if command == "wishlist" {
		str := "Wish list:\n\n0) Деньги)\n"
		gifts, err := s.db.GetWishlist()
		if err != nil {
			s.bots.Bot.Send(tgbotapi.NewMessage(s.AdminID, "Ошибка"+err.Error()))
			return
		}
		for _, g := range gifts {
			str += strconv.Itoa(int(pointy.PointerValue(g.ID, 0))) + ") " +
				pointy.PointerValue(g.Description, "") + " " + fmt.Sprint(*g.UserID) + "\n"
		}
		s.bots.Bot.Send(tgbotapi.NewMessage(s.AdminID, str))
		return
	}

	if command == "refresh" {
		guests, err := s.db.GetGuests()
		if err != nil {
			s.bots.Bot.Send(tgbotapi.NewMessage(s.AdminID, "Ошибка"+err.Error()))
			return
		}
		for _, g := range guests {
			if pointy.PointerValue(g.UserID, 0) != 0 {
				s.sendInfo(*g.UserID)
			}
		}
	}

	s.bots.Bot.Send(tgbotapi.NewMessage(s.AdminID, "Я не знаю как на это реагировать(\nКоманды:\n/newGuest - добавить гостя\n/guests - все приглашения\n/delGuest - удалить гостя по логину"))
}
