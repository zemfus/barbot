package service

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func (s *Service) handleGuest(update tgbotapi.Update) {
	if update.Message == nil {
		return
	}
	command := update.Message.Command()

	if command == "start" {
		guest, err := s.db.CheckGuest(update.SentFrom().UserName)
		if err != nil {
			s.bots.Bot.Send(tgbotapi.NewMessage(s.AdminID, "Траблы с запросом /start от пользователя "+update.SentFrom().UserName))
			return
		}
		if guest == "" {
			s.bots.Bot.Send(tgbotapi.NewMessage(s.AdminID, "несанкционированный доступ от "+update.SentFrom().UserName))
			s.bots.Bot.Send(tgbotapi.NewMessage(update.SentFrom().ID, "Проблемка, напиши @isuprun"))
			return
		}
		s.bots.Bot.Send(tgbotapi.NewMessage(update.SentFrom().ID, "Привет, "+guest))
		return
	}
}
