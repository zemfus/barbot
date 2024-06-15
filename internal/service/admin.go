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
		s.bots.Bot.Send(tgbotapi.NewMessage(s.AdminID, "Команды:\n/new_guest - добавить гостя\n/invitations - все приглашения"))
		return
	}

	if command == "new_guest" {
		s.bots.Bot.Send(tgbotapi.NewMessage(s.AdminID, "введи через пробел: логин (без @), имя, уровень (0 - без алко, 1 - алко, 2 - спешл)"))
		state = NewGuest
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

	if command == "invitations" {
		allInv, err := s.db.GetInvitations()
		if err != nil {
			s.bots.Bot.Send(tgbotapi.NewMessage(s.AdminID, "Траблы с запросом"))
			return
		}
		str := ""
		for i, v := range allInv {
			str += strconv.Itoa(i) + ") " + pointy.PointerValue(v.Login, "") + " " + pointy.PointerValue(v.Name, "") + " " + fmt.Sprint(pointy.PointerValue(v.Level, 0)) + "\n"
		}
		s.bots.Bot.Send(tgbotapi.NewMessage(s.AdminID, str))
		return
	}

}
