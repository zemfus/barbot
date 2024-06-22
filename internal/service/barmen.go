package service

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
	"strings"
)

var barmenState int = BarmenNone

const (
	BarmenNone = iota
)

var AllBarmenCommands string = "Команды:\n" +
	"/menu - показать доступны\n\n" +
	"/newCocktail - добавить коктейль\n" +
	"/delCocktail - удалить коктейль\n" +
	"/cocktails - посмотреть все коктейли\n" +
	"/newOrder - создать заказ"

func (s *Service) handleBarmen(update tgbotapi.Update) {
	s.handleBarmenCommand(update)
}

var tmp string = "BAACAgIAAxkBAAIHLGZy25Uf9PZPG3L9HXZktwABUKIZRwACSU4AAiDtmUsJ1jSLEjtn-DUE"

func (s *Service) handleBarmenCommand(update tgbotapi.Update) {
	if update.CallbackQuery != nil {
		if strings.HasPrefix(update.CallbackQuery.Data, "done") {
			id, _ := strconv.ParseInt(update.CallbackQuery.Data[4:], 10, 64)
			guest, err := s.db.GetGuest(id)
			if err != nil {
				s.bots.Bot.Send(tgbotapi.NewMessage(s.AdminID, "нет человека в базе (момент отправки меню)"))
				return
			}
			s.sendCocktails(id, *guest.Level, false)
			s.bots.Bot.Send(tgbotapi.NewDeleteMessage(s.BarmenID, update.CallbackQuery.Message.MessageID))
			return
		}
		s.bots.Bot.Send(tgbotapi.NewDeleteMessage(s.BarmenID, update.CallbackQuery.Message.MessageID))

		id, err := strconv.ParseInt(update.CallbackQuery.Data, 10, 64)
		if err == nil {
			guest, err := s.db.GetGuest(id)
			if err != nil {
				s.bots.Bot.Send(tgbotapi.NewMessage(s.AdminID, "нет человека в базе (момент отдачи заказа)"))
				return
			}
			msg := tgbotapi.NewPhoto(s.BarmenID, tgbotapi.FileID(*guest.Photo))
			msg.Caption = update.CallbackQuery.Message.Text
			msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("Готово", "done"+update.CallbackQuery.Data)))
			s.bots.Bot.Send(msg)
			return
		}
	}
	command := update.Message.Command()

	if command == "start" {
		s.bots.Bot.Send(tgbotapi.NewVideo(s.BarmenID, tgbotapi.FileID(tmp)))
		s.bots.Bot.Send(tgbotapi.NewMessage(s.BarmenID, "Полетели🧚‍\n\n"+AllBarmenCommands))
		return
	}

}
