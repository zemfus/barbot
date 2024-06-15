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

	s.bots.Bot.Send(tgbotapi.NewMessage(s.AdminID, "Я не знаю как на это реагировать(\nКоманды:\n/newGuest - добавить гостя\n/guests - все приглашения\n/delGuest - удалить гостя по логину"))
}
