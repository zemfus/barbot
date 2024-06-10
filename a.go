package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strconv"
	"strings"
)

var questions = []string{
	"1. Близнецы по знаку зодиака.",
	"2. Есть опыт в программировании.",
	"3. Студент(ка) ВУЗа.",
	"4. Есть домашний питомец.",
	"5. Занимаюсь танцами.",
	"6. Закончил(а) школу с золотой медалью.",
	"7. Выполнил(а) первый день бассейна.",
	"8. Играю на музыкальном инструменте.",
	"9. Не подписался(лась) на второй день.",
	"10. Говорю на двух и более языках.",
	"11. Был(а) за границей.",
	"12. Любимый напиток по утрам - кофе.",
	"13. Карий цвет глаз.",
	"14. Умею кататься на сноуборде.",
	"15. Есть водительские права.",
}

func main() {
	bot, err := tgbotapi.NewBotAPI("token")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for u := range updates {
		Handle(bot, u)
	}
}

func Handle(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	if update.CallbackQuery != nil {
		// Обработка нажатия на inline кнопку
		handleCallbackQuery(bot, update.CallbackQuery)
	} else if update.Message != nil {
		// Обработка обычных сообщений
		handleMessage(bot, update.Message)
	}
}

func handleMessage(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	id := message.From.ID

	// проверяем есть ли пользователь в базе
	if !exist(id) {
		if message.Text == "/start" {
			// добавляем нового пользователя
			// newUser(id)
			msg := tgbotapi.NewMessage(message.Chat.ID, "Введи свой ник на платформе:")
			bot.Send(msg)
		} else {
			msg := tgbotapi.NewMessage(message.Chat.ID, "я не знаю такой команды(")
			bot.Send(msg)
		}
		return
	}

	// проверяем записан ли логин у человека
	if getLogin(id) == "" {
		// добавлемя логин к пользователю
		setLogin(id, message.Text)
		msg := tgbotapi.NewMessage(message.Chat.ID, "Привет, "+message.Text)
		bot.Send(msg)
		return
	}

	defer sendGrid(bot, message.Chat.ID, id)
	// в бд должен храниться номер вопроса который сейчас заполняется.
	// 0, если на данный момент чтение пользователя не активно
	answer := readAnswer(id)
	if answer > 0 {
		// проверяем есть ли логин в бд
		// если есть то возвращается id иначе 0
		anotherID := checkLogin(message.Text)
		if anotherID != 0 {
			msg := tgbotapi.NewMessage(message.Chat.ID, "я не знаю такого человека(")
			bot.Send(msg)
			return
		}
		// получаем ответы данного пользователя (логины, точнее id, которые он уже упоминал при ответе)
		// в виде массива int
		users := getAnswers(id)
		if findNum(users, anotherID) {
			msg := tgbotapi.NewMessage(message.Chat.ID, "Ты уже упоминал этого человека ранее)")
			bot.Send(msg)
			return
		}

		// проверяем по логину стоит ли true у ддругого пользователя на вопросе answer
		if checkAnswer(anotherID, answer) {
			// если стоит то помечаем у человека ответ на вопрос answer
			setAnswerTrue(id, answer, anotherID)
			msg := tgbotapi.NewMessage(message.Chat.ID, "Ответ записан")
			bot.Send(msg)
			return
		} else {
			msg := tgbotapi.NewMessage(message.Chat.ID, "Неправильно)")
			bot.Send(msg)
			return
		}
	} else {
		msg := tgbotapi.NewMessage(message.Chat.ID, "я не знаю такой команды(")
		bot.Send(msg)
		return
	}
}

func handleCallbackQuery(bot *tgbotapi.BotAPI, callbackQuery *tgbotapi.CallbackQuery) {
	id := callbackQuery.From.ID
	chatID := callbackQuery.Message.Chat.ID
	if callbackQuery.Data[0] == 't' || callbackQuery.Data[0] == 'f' {
		answer, _ := strconv.Atoi(callbackQuery.Data[1:])
		// помечаем у человека ответ
		setAnswer(callbackQuery.From.ID, answer, callbackQuery.Data[0] == 't')
		if answer < len(questions) {
			sendQuestion(bot, chatID, id, answer+1)
		} else {
			sendGrid(bot, chatID, id)
		}
		return
	}

	if callbackQuery.Data[0] == 'a' {
		answer, _ := strconv.Atoi(callbackQuery.Data[1:])
		// помечаем какой ответ сейчас будет читаться
		setReadAnswer(id, answer)
		msg := tgbotapi.NewMessage(
			callbackQuery.Message.Chat.ID,
			"Введите ник пользователя про которого данный факт является правдой: "+questions[answer-1])
		bot.Send(msg)
	}
}

func findNum(arr []int, target int) bool {
	for _, v := range arr {
		if v == target {
			return true
		}
	}
	return false
}

func sendGrid(bot *tgbotapi.BotAPI, chatID int64, id int) {
	msg := tgbotapi.NewMessage(chatID, strings.Join(questions, "\n"))
	var rows [][]tgbotapi.InlineKeyboardButton

	// Создание кнопок и добавление их в строки
	for i := 1; i <= len(questions); i++ {
		var button tgbotapi.InlineKeyboardButton
		users := getAnswers(id)
		if users[i-1] == "" {
			button = tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%d", i), fmt.Sprintf("a%d", i))
		} else {
			button = tgbotapi.NewInlineKeyboardButtonData("_", "_")
		}
		if i%5 == 1 {
			// Добавление новой строки
			rows = append(rows, tgbotapi.NewInlineKeyboardRow(button))
		} else {
			// Добавление кнопки в последнюю строку
			rows[len(rows)-1] = append(rows[len(rows)-1], button)
		}
	}

	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(rows...)
	bot.Send(msg)
}

func sendQuestion(bot *tgbotapi.BotAPI, chatID int64, id int, q int) {
	msg := tgbotapi.NewMessage(chatID, "Является ли следующий факт о тебе правдой? \n\n"+questions[q])
	// Создаем встроенную клавиатуру с кнопками
	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Правда", "t"+strconv.Itoa(q)),
			tgbotapi.NewInlineKeyboardButtonData("Ложь", "f"+strconv.Itoa(q)),
		),
	)
	msg.ReplyMarkup = inlineKeyboard
	bot.Send(msg)
}
