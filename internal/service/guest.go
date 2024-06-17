package service

import (
	"barbot/internal/repository/postgres"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.openly.dev/pointy"
	"strconv"
	"strings"
)

func (s *Service) handleGuest(update tgbotapi.Update) {
	if update.CallbackQuery != nil {
		if update.CallbackQuery.Data == "_" {
			return
		}
		s.bots.Bot.Send(tgbotapi.NewDeleteMessage(update.SentFrom().ID, update.CallbackQuery.Message.MessageID))

		if update.CallbackQuery.Data == "approveInvite_" {
			err := s.db.SetParticipation(update.SentFrom().ID, true)
			if err != nil {
				s.bots.Bot.Send(tgbotapi.NewMessage(s.AdminID, "Траблы с подтверждением участия "+update.SentFrom().UserName))
				return
			}
			s.bots.Bot.Send(tgbotapi.NewMessage(update.SentFrom().ID, "Будем ждать 22.06 в 20:00, место будет объявлено накануне"))
			s.bots.Bot.Send(tgbotapi.NewMessage(s.AdminID, "@"+update.SentFrom().UserName+" идет на мероприятие"))
			s.sendInfo(update.SentFrom().ID)
			permission := s.NewPermission(update.SentFrom().ID, true)
			s.bots.Bot.Send(permission)
		}

		if update.CallbackQuery.Data == "refuseInvite_" {
			err := s.db.SetParticipation(update.SentFrom().ID, true)
			if err != nil {
				s.bots.Bot.Send(tgbotapi.NewMessage(s.AdminID, "Траблы с подтверждением участия "+update.SentFrom().UserName))
				return
			}
			s.bots.Bot.Send(tgbotapi.NewMessage(update.SentFrom().ID, "Жаль( но ты можешь вернуться и изменить свое решение"))
			s.sendInvite(update.SentFrom().ID)
			s.bots.Bot.Send(tgbotapi.NewMessage(s.AdminID, "@"+update.SentFrom().UserName+" отказался от мероприятия"))
			permission := s.NewPermission(update.SentFrom().ID, false)
			s.bots.Bot.Send(permission)
		}

		if update.CallbackQuery.Data == "alcohol_" {
			s.bots.Bot.Send(tgbotapi.NewMessage(update.SentFrom().ID, "Напиши пожелания по алкоголю:"))
			s.db.SetState(update.SentFrom().ID, postgres.GuestAlcohol)
		}

		if update.CallbackQuery.Data == "music_" {
			s.bots.Bot.Send(tgbotapi.NewMessage(update.SentFrom().ID, "Напиши пожелания по музыке:"))
			s.db.SetState(update.SentFrom().ID, postgres.GuestMusic)
		}

		if update.CallbackQuery.Data == "food_" {
			s.bots.Bot.Send(tgbotapi.NewMessage(update.SentFrom().ID, "Напиши пожелания по еде:"))
			s.db.SetState(update.SentFrom().ID, postgres.GuestFood)
		}

		if update.CallbackQuery.Data == "invite_" {
			s.sendInvite(update.SentFrom().ID)
		}

		if update.CallbackQuery.Data == "return" {
			s.sendInfo(update.SentFrom().ID)
		}

		if strings.HasPrefix(update.CallbackQuery.Data, "cancel") {
			idx, _ := strconv.Atoi(update.CallbackQuery.Data[7:])
			_, err := s.db.SetGiftUserID(int32(idx), 0)
			if err != nil {
				s.bots.Bot.Send(tgbotapi.NewMessage(s.AdminID, "Траблы с отменой подарка "+update.SentFrom().UserName))
				return
			}
			s.bots.Bot.Send(tgbotapi.NewMessage(update.SentFrom().ID, "ГАЛЯ У НАС ОТМЕНА"))
			s.sendWishlist(update.SentFrom().ID)
		}

		if strings.HasPrefix(update.CallbackQuery.Data, "add") {
			idx, _ := strconv.Atoi(update.CallbackQuery.Data[4:])
			ret, err := s.db.SetGiftUserID(int32(idx), update.SentFrom().ID)
			if err != nil {
				s.bots.Bot.Send(tgbotapi.NewMessage(s.AdminID, "Траблы с добавлением подарка "+update.SentFrom().UserName))
				return
			}
			if !ret {
				s.bots.Bot.Send(tgbotapi.NewMessage(update.SentFrom().ID, "Уже ктото взял("))
			} else {
				s.bots.Bot.Send(tgbotapi.NewMessage(update.SentFrom().ID, "Так и запишем ✍️"))
			}
			s.bots.Bot.Send(tgbotapi.NewDeleteMessage(update.SentFrom().ID, update.CallbackQuery.Message.MessageID))
			s.sendWishlist(update.SentFrom().ID)
		}

		if update.CallbackQuery.Data == "wishlist" {
			s.bots.Bot.Send(tgbotapi.NewDeleteMessage(update.SentFrom().ID, update.CallbackQuery.Message.MessageID))
			s.sendWishlist(update.SentFrom().ID)
		}

		return
	}

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

		if guest.Name == nil {
			s.bots.Bot.Send(tgbotapi.NewMessage(s.AdminID, "несанкционированный доступ от "+update.SentFrom().UserName))
			s.bots.Bot.Send(tgbotapi.NewMessage(update.SentFrom().ID, "Проблемка, напиши @isuprun"))
			return
		}

		s.db.SetID(update.SentFrom().UserName, update.SentFrom().ID)
		s.bots.Bot.Send(tgbotapi.NewMessage(update.SentFrom().ID, "Привет, "+*guest.Name))
		s.sendInvite(update.SentFrom().ID)
		return
	}

	state, err := s.db.GetState(update.SentFrom().ID)
	if err != nil {
		s.bots.Bot.Send(tgbotapi.NewMessage(s.AdminID, "Проблемы со get стэйтом у "+update.SentFrom().UserName))
		return
	}

	if state == postgres.GuestAlcohol || state == postgres.GuestMusic || state == postgres.GuestFood {
		str := "@" + update.SentFrom().UserName + " хочет"
		switch state {
		case postgres.GuestAlcohol:
			str += " Alcohol "
		case postgres.GuestMusic:
			str += " Music "
		case postgres.GuestFood:
			str += " Food "
		}
		err := s.db.SetState(update.SentFrom().ID, postgres.GuestNone)
		if err != nil {
			s.bots.Bot.Send(tgbotapi.NewMessage(s.AdminID, "Проблемы со set стэйтом у "+update.SentFrom().UserName))
			return
		}
		str += update.Message.Text
		s.bots.Bot.Send(tgbotapi.NewMessage(s.AdminID, str))
		s.bots.Bot.Send(tgbotapi.NewMessage(update.SentFrom().ID, "Так и запишем ✍️"))
		s.sendInfo(update.SentFrom().ID)
		return
	}

	s.bots.Bot.Send(tgbotapi.NewMessage(update.SentFrom().ID, "Я не знаю как на это реагировать"))
}

func (s *Service) sendInvite(id int64) {
	msg := tgbotapi.NewMessage(id, "Подтверди свое участие:")
	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Подтверждаю", "approveInvite_"),
			tgbotapi.NewInlineKeyboardButtonData("Не смогу(", "refuseInvite_"),
		),
	)
	msg.ReplyMarkup = inlineKeyboard
	s.bots.Bot.Send(msg)
}

func (s *Service) sendInfo(id int64) {
	msg := tgbotapi.NewMessage(id, "Можешь выбрать либого питомца, которого ты захочешь:")
	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Пожелания по алкоголю", "alcohol_")),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Предпочтения по музыке", "music_")),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Предпочтения по еде", "food_")),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Wishlist", "wishlist")),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Изменить решение", "invite_")),
	)
	msg.ReplyMarkup = inlineKeyboard
	s.bots.Bot.Send(msg)
}

func (s *Service) sendWishlist(id int64) {
	gifts, err := s.db.GetWishlist()
	if err != nil {
		s.bots.Bot.Send(tgbotapi.NewMessage(s.AdminID, "Ошибка"+err.Error()))
		return
	}

	var rows [][]tgbotapi.InlineKeyboardButton
	str := "Wishlist:\n\n"

	for i := 1; i <= len(gifts); i++ {
		var button tgbotapi.InlineKeyboardButton
		switch pointy.PointerValue(gifts[i-1].UserID, 0) {
		case id:
			button = tgbotapi.NewInlineKeyboardButtonData(fmt.Sprint("✅", i, "✅"), fmt.Sprint("cancel ", i))
			str += "🟢 "
		case 0:
			button = tgbotapi.NewInlineKeyboardButtonData(fmt.Sprint(i), fmt.Sprint("add ", i))
			str += "⚪️ "
		default:
			button = tgbotapi.NewInlineKeyboardButtonData(fmt.Sprint("❌", i, "❌"), "_")
			str += "🔴 "
		}
		// Создание кнопок и добавление их в строки
		// Добавление новой строки
		// Добавление кнопки в последнюю строку
		str += strconv.Itoa(int(pointy.PointerValue(gifts[i-1].ID, 0))) + ") "
		str += pointy.PointerValue(gifts[i-1].Description, "") + "\n"

		if i%5 == 1 {
			rows = append(rows, tgbotapi.NewInlineKeyboardRow(button))
		} else {
			rows[len(rows)-1] = append(rows[len(rows)-1], button)
		}

	}
	str += "\n🟢 -- выбранные подарки (при нажатии на ✅ можно отменить бронь)\n⚪️ -- доступные для выбора\n🔴 -- забронированые"
	rows = append(rows,
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Вернуться обратно", "return")))
	msg := tgbotapi.NewMessage(id, str)
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(rows...)
	s.bots.Bot.Send(msg)
}

func (s *Service) NewPermission(id int64, permission bool) tgbotapi.PromoteChatMemberConfig {
	return tgbotapi.PromoteChatMemberConfig{
		ChatMemberConfig: tgbotapi.ChatMemberConfig{
			ChatID:             s.ChatID,
			SuperGroupUsername: "approve",
			ChannelUsername:    "",
			UserID:             id,
		},
		IsAnonymous:         false,
		CanManageChat:       false,
		CanChangeInfo:       false,
		CanPostMessages:     permission,
		CanEditMessages:     false,
		CanDeleteMessages:   false,
		CanManageVoiceChats: permission,
		CanInviteUsers:      false,
		CanRestrictMembers:  false,
		CanPinMessages:      false,
		CanPromoteMembers:   false,
	}
}
