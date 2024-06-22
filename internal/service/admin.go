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
	CheckIn
	NewCocktail
	AddLevel
	AddComposition
	DelCocktail
)

var adminState int = None

var AllCommands string = "Команды:\n" +
	"/menu - показать доступны\n\n" +
	"/newGuest - добавить гостя\n" +
	"/delGuest - удалить гостя\n" +
	"/guests - показать всех гостей\n\n" +
	"/newGift - добавить подарок в вишлист\n" +
	"/delGift - удалить подарок из вишлиста\n" +
	"/wishlist - показать вишлист\n\n" +
	"/checkIn - регистрация гостя\n\n" +
	"/newCocktail - новый коктейль\n" +
	"/delCocktail - удалить коктейль\n" +
	"/cocktails - посмотреть все коктейли\n" +
	"/newOrder - создать заказ"

func (s *Service) handleAdmin(update tgbotapi.Update) {
	if s.handleAdminState(update) {
		return
	}
	if s.handleAdminCommand(update) {
		return
	}
	if update.Message.Photo != nil {
		s.bots.Bot.Send(tgbotapi.NewPhoto(s.AdminID, tgbotapi.FileID(update.Message.Photo[len(update.Message.Photo)-1].FileID)))
		s.bots.Bot.Send(tgbotapi.NewMessage(s.AdminID, update.Message.Photo[len(update.Message.Photo)-1].FileID))
	}
	s.bots.Bot.Send(tgbotapi.NewMessage(s.AdminID,
		"Я не знаю как на это реагировать(\n"+AllCommands))
}

func (s *Service) handleAdminCommand(update tgbotapi.Update) bool {
	command := update.Message.Command()
	switch command {
	case "menu":
		s.adminCommandMenu()
	case "start":
		s.adminCommandStart()
	case "newGuest":
		s.adminCommandNewGuest()
	case "delGuest":
		s.adminCommandDelGuest()
	case "guests":
		s.adminCommandGuests()
	case "newGift":
		s.adminCommandNewGift()
	case "delGift":
		s.adminCommandDelGift()
	case "wishlist":
		s.adminCommandWishlist()
	case "checkIn":
		s.adminCommandCheckIn()
	case "newCocktail":
		s.adminCommandNewCocktail()
	case "delCocktail":
		s.adminCommandDelCocktail()
	case "cocktails":
		s.adminCommandCocktails()
	default:
		return false
	}
	return true
}

func (s *Service) handleAdminState(update tgbotapi.Update) bool {
	switch adminState {
	case NewGuest:
		s.adminHandleNewGuest(update)
	case DelGuest:
		s.adminHandleDelGuest(update)
	case NewGift:
		s.adminHandleNewGift(update)
	case DelGift:
		s.adminHandleDelGift(update)
	case CheckIn:
		s.adminHandleCheckIn(update)
	case NewCocktail:
		s.adminHandleNewCocktail(update)
	case AddLevel:
		s.adminHandleAddLevel(update)
	case AddComposition:
		s.adminHandleAddComposition(update)
	case DelCocktail:
		s.adminHandleDelCocktail(update)

	default:
		return false
	}
	return true
}

func (s *Service) adminCommandMenu() {
	s.bots.Bot.Send(tgbotapi.NewMessage(s.AdminID, AllCommands))
}

func (s *Service) adminCommandStart() {
	s.bots.Bot.Send(tgbotapi.NewMessage(s.AdminID, AllCommands))
}

func (s *Service) adminCommandNewGuest() {
	s.bots.Bot.Send(tgbotapi.NewMessage(s.AdminID, "введи через пробел: логин (без @), имя, уровень (0 - без алко, 1 - алко, 2 - спешл)"))
	adminState = NewGuest
}

func (s *Service) adminCommandDelGuest() {
	s.bots.Bot.Send(tgbotapi.NewMessage(s.AdminID, "введи логин (без @)"))
	adminState = DelGuest
}
func (s *Service) adminCommandGuests() {
	guests, err := s.db.GetGuests()
	if err != nil {
		s.bots.Bot.Send(tgbotapi.NewMessage(s.AdminID, "Траблы с запросом"))
		return
	}
	str := "Guests:\n"
	for i, g := range guests {
		str += fmt.Sprint(
			strconv.Itoa(i+1), ") ",
			pointy.PointerValue(g.UserID, 0), " @",
			pointy.PointerValue(g.Login, ""), " ",
			pointy.PointerValue(g.Name, ""), " ",
			pointy.PointerValue(g.State, 0), " ",
			pointy.PointerValue(g.Level, 0), " ",
			pointy.PointerValue(g.Participation, false), " ",
			pointy.PointerValue(g.CheckIn, false), "\n",
			pointy.PointerValue(g.Photo, ""), "\n",
		)
	}
	s.bots.Bot.Send(tgbotapi.NewMessage(s.AdminID, str))
}

func (s *Service) adminCommandNewGift() {
	s.bots.Bot.Send(tgbotapi.NewMessage(s.AdminID, "Напиши что ты хочешь"))
	adminState = NewGift
}

func (s *Service) adminCommandDelGift() {
	s.bots.Bot.Send(tgbotapi.NewMessage(s.AdminID, "Введи номер"))
	adminState = DelGift
}

func (s *Service) adminCommandWishlist() {
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
	msg := tgbotapi.NewMessage(s.AdminID, str)
	msg.DisableWebPagePreview = true
	s.bots.Bot.Send(msg)
	return
}

func (s *Service) adminCommandCheckIn() {
	s.bots.Bot.Send(tgbotapi.NewMessage(s.AdminID, "введи через пробел: логин (без @) и прикрепи фото"))
	adminState = CheckIn
}

func (s *Service) adminCommandNewCocktail() {
	s.bots.Bot.Send(tgbotapi.NewMessage(s.AdminID, "введи название нового коктейля"))
	adminState = NewCocktail
}

func (s *Service) adminCommandDelCocktail() {
	s.bots.Bot.Send(tgbotapi.NewMessage(s.AdminID, "введи название коктейля для удаления"))
	adminState = DelCocktail
}

func (s *Service) adminCommandCocktails() {
	str := "Cocktails\n\n"
	for i := 0; i < 6; i++ {
		cocktails, err := s.db.GetCocktails(i)
		if err != nil {
			s.bots.Bot.Send(tgbotapi.NewMessage(s.AdminID, "что-то с запросами"))
			return
		}
		for _, v := range cocktails {
			str += strconv.Itoa(int(v.ID)) + " " +
				*v.Name + ": " +
				*v.Composition + " \n" +
				"доступность: " + strconv.FormatBool(*v.Availability) + "\n" +
				"создал бармен:" + strconv.FormatBool(*v.Barmen) + "\n" +
				"уровень: " + strconv.Itoa(int(*v.Level)) + "\n"
		}
		str += "\n"
	}
	s.bots.Bot.Send(tgbotapi.NewMessage(s.AdminID, str))
}

func (s *Service) adminCommandNewOrder() {

}

func (s *Service) adminHandleNewGuest(update tgbotapi.Update) {
	adminState = None
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
}

func (s *Service) adminHandleDelGuest(update tgbotapi.Update) {
	adminState = None
	err := s.db.DropGuest(update.Message.Text)
	if err != nil {
		s.bots.Bot.Send(tgbotapi.NewMessage(s.AdminID, "Что-то пошло не так в удалении"))
		return
	}
	s.bots.Bot.Send(tgbotapi.NewMessage(s.AdminID, "Без ошибок прошло удаление"))
}

func (s *Service) adminHandleNewGift(update tgbotapi.Update) {
	adminState = None
	if s.db.NewGift(update.Message.Text) {
		s.bots.Bot.Send(tgbotapi.NewMessage(s.AdminID, "Добавили"))
	} else {
		s.bots.Bot.Send(tgbotapi.NewMessage(s.AdminID, "Косяк с добавлением"))
	}
}

func (s *Service) adminHandleDelGift(update tgbotapi.Update) {
	adminState = None
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
}

func (s *Service) adminHandleCheckIn(update tgbotapi.Update) {
	adminState = None
	if update.Message.Caption == "" {
		s.bots.Bot.Send(tgbotapi.NewMessage(s.AdminID, "А где ник?"))
		return
	}
	if len(update.Message.Photo) == 0 {
		s.bots.Bot.Send(tgbotapi.NewMessage(s.AdminID, "А где фото?"))
		return
	}
	photo := update.Message.Photo[len(update.Message.Photo)-1].FileID
	err := s.db.CheckIn(update.Message.Caption, photo)
	if err != nil {
		s.bots.Bot.Send(tgbotapi.NewMessage(s.AdminID, "Что-то пошло не так check in"))
		return
	}
	guest, _ := s.db.CheckGuest(update.Message.Caption)
	if err != nil {
		s.bots.Bot.Send(tgbotapi.NewMessage(s.AdminID, "проблемы с отправкой меню"))
		return
	}
	s.sendCocktails(*guest.UserID, *guest.Level, false)
	s.bots.Bot.Send(tgbotapi.NewMessage(s.AdminID, "check-in успешно, меню успешно отправлено"))
}

var AdminCocktailName string

func (s *Service) adminHandleNewCocktail(update tgbotapi.Update) {
	adminState = None
	msg := tgbotapi.NewMessage(s.AdminID, "Тип коктейля")
	AdminCocktailName = update.Message.Text
	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("без алко кордилы", "0"+update.Message.Text),
			tgbotapi.NewInlineKeyboardButtonData("без алко", "1"+update.Message.Text),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("алко кордилы", "2"+update.Message.Text),
			tgbotapi.NewInlineKeyboardButtonData("алко", "3"+update.Message.Text),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("антивирус", "5"+update.Message.Text),
			tgbotapi.NewInlineKeyboardButtonData("special", "6"+update.Message.Text),
		),
	)
	adminState = AddLevel
	msg.ReplyMarkup = inlineKeyboard
	s.bots.Bot.Send(msg)
}

func (s *Service) adminHandleAddLevel(update tgbotapi.Update) {
	adminState = None
	if update.CallbackQuery == nil {
		return
	}
	s.bots.Bot.Send(tgbotapi.NewDeleteMessage(s.AdminID, update.CallbackQuery.Message.MessageID))
	if !s.db.NewCocktail(
		update.CallbackQuery.Data[1:],
		false,
		int(update.CallbackQuery.Data[0]-'0'),
	) {
		s.bots.Bot.Send(tgbotapi.NewMessage(s.AdminID, "Что-то пошло не так с добавлением"))
		return
	}
	adminState = AddComposition
	s.bots.Bot.Send(tgbotapi.NewMessage(s.AdminID, "Введи описание для коктейля"))
}

func (s *Service) adminHandleAddComposition(update tgbotapi.Update) {
	adminState = None
	if !s.db.SetComposition(
		AdminCocktailName,
		update.Message.Text,
	) {
		s.bots.Bot.Send(tgbotapi.NewMessage(s.AdminID, "Что-то пошло не так с добавлением описания"))
		return
	}
	s.bots.Bot.Send(tgbotapi.NewMessage(s.AdminID, "Добавлено"))
}

func (s *Service) adminHandleDelCocktail(update tgbotapi.Update) {
	adminState = None
	if s.db.DelCocktail(update.Message.Text) {
		s.bots.Bot.Send(tgbotapi.NewMessage(s.AdminID, "Удален"))
	} else {
		s.bots.Bot.Send(tgbotapi.NewMessage(s.AdminID, "Проблемы с удалением"))
	}
}
