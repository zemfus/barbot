package service

import (
	"barbot/internal/bots"
	repoP "barbot/internal/repository/postgres"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"math/rand"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Service struct {
	db      *repoP.RepositoryPostgres
	bots    *bots.Bots
	AdminId []int64
}

var questions = map[int]string{
	1:  "1. Близнецы по знаку зодиака.",
	2:  "2. Есть опыт в программировании.",
	3:  "3. Студент(ка) ВУЗа.",
	4:  "4. Есть домашний питомец.",
	5:  "5. Занимаюсь танцами.",
	6:  "6. Закончил(а) школу с золотой медалью.",
	7:  "7. Выполнил(а) первый день бассейна.",
	8:  "8. Играю на музыкальном инструменте.",
	9:  "9. Не подписался(лась) на второй день.",
	10: "10. Говорю на двух и более языках.",
	11: "11. Был(а) за границей.",
	12: "12. Любимый напиток по утрам - кофе.",
	13: "13. Карий цвет глаз.",
	14: "14. Умею кататься на сноуборде.",
	15: "15. Есть водительские права.",
}

func NewService(
	db *repoP.RepositoryPostgres,
	bots *bots.Bots,
	adminId []int64,
) *Service {
	return &Service{db: db, bots: bots, AdminId: adminId}
}

func (s *Service) Run() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := s.bots.Bot.GetUpdatesChan(u)

	workerChannels := make(map[int64]chan tgbotapi.Update)
	var wg sync.WaitGroup
	workers := int64(runtime.NumCPU())

	for i := int64(0); i < workers; i++ {
		workerChannels[i] = make(chan tgbotapi.Update)
		wg.Add(1)
		go func(id int64, ch chan tgbotapi.Update) {
			defer wg.Done()
			for msg := range ch {
				s.processUpdate(msg)
			}
		}(i, workerChannels[i])
	}

	for update := range updates {
		if u := update.SentFrom(); u != nil {
			workerChannels[u.ID%workers] <- update
		}

	}
	for _, ch := range workerChannels {
		close(ch)
	}
	wg.Wait()
}

func (s *Service) processUpdate(update tgbotapi.Update) {
	userID := update.SentFrom().ID

	if update.CallbackQuery != nil {
		if update.CallbackQuery.Data[0] == 't' || update.CallbackQuery.Data[0] == 'f' {
			answer, _ := strconv.Atoi(update.CallbackQuery.Data[1:])
			if answer < len(questions) {
				s.db.SaveAnswer(userID, int64(answer), update.CallbackQuery.Data[0] == 't')
				s.sendQuestion(userID, answer+1, update.CallbackQuery.Message.MessageID)
			}
		}
		return
	}

	if update.Message != nil && update.Message.Command() == "start" {
		if s.needRegistration(update) {
			m := tgbotapi.NewMessage(userID, "Привет, Дорогой ПИР, введи свой логин")
			s.bots.Bot.Send(m)
		}
		d := tgbotapi.NewDeleteMessage(userID, update.Message.MessageID)
		s.bots.Bot.Send(d)
		return
	}

	if update.Message != nil && s.db.GetState(userID) == 0 {
		arr := strings.Split(update.Message.Text, "@")
		if len(arr) == 0 {
			m := tgbotapi.NewMessage(userID, "Не корректный ник, введи еще раз")
			s.bots.Bot.Send(m)
			return
		}
		s.db.SetLogin(userID, strings.ToLower(arr[0]))
		m := tgbotapi.NewMessage(userID, "Привет, "+arr[0])
		a, _ := s.bots.Bot.Send(m)
		s.sendQuestion(userID, 1, a.MessageID)
	}

	if update.Message != nil && update.Message.Command() == "new" && (s.AdminId[1] == update.SentFrom().ID || s.AdminId[0] == update.SentFrom().ID) {
		users := s.db.GetUsers()
		teamAssignments := distributeTeams(users)
		for id, team := range teamAssignments {
			m1 := tgbotapi.NewMessage(id, fmt.Sprintf("Твоя команда Аксолотлей номер: %d ❤️️", team))
			s.bots.Bot.Send(m1)
		}
		return
	}

	if update.Message != nil && update.Message.Command() == "count" && (s.AdminId[1] == update.SentFrom().ID || s.AdminId[0] == update.SentFrom().ID) {
		users := s.db.GetUsers()
		m1 := tgbotapi.NewMessage(update.SentFrom().ID, fmt.Sprintf("Колличество: %d", len(users)))
		s.bots.Bot.Send(m1)
		return
	}

}

//	func (s *Service) sendGrid(bot *tgbotapi.BotAPI, chatID int64, id int) {
//		msg := tgbotapi.NewMessage(chatID, strings.Join(questions, "\n"))
//		var rows [][]tgbotapi.InlineKeyboardButton
//
//		// Создание кнопок и добавление их в строки
//		for i := 1; i <= len(questions); i++ {
//			var button tgbotapi.InlineKeyboardButton
//			users := getAnswers(id)
//			if users[i-1] == "" {
//				button = tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%d", i), fmt.Sprintf("a%d", i))
//			} else {
//				button = tgbotapi.NewInlineKeyboardButtonData("_", "_")
//			}
//			if i%5 == 1 {
//				// Добавление новой строки
//				rows = append(rows, tgbotapi.NewInlineKeyboardRow(button))
//			} else {
//				// Добавление кнопки в последнюю строку
//				rows[len(rows)-1] = append(rows[len(rows)-1], button)
//			}
//		}
//
//		msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(rows...)
//		bot.Send(msg)
//	}
func (s *Service) sendQuestion(userID int64, q int, msgID int) {
	// Создаем встроенную клавиатуру с кнопками
	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Правда", "t"+strconv.Itoa(q)),
			tgbotapi.NewInlineKeyboardButtonData("Ложь", "f"+strconv.Itoa(q)),
		),
	)
	msg := tgbotapi.NewEditMessageTextAndMarkup(userID, msgID, "Является ли следующий факт о тебе правдой? \n\n"+questions[q], inlineKeyboard)
	s.bots.Bot.Send(msg)
}

func (s *Service) needRegistration(update tgbotapi.Update) bool {
	if u := update.SentFrom(); u != nil {
		return s.db.SaveID(u.ID)
	}
	return false
}

func shuffle(slice []int64) {
	rand.New(rand.NewSource(time.Now().UnixMilli()))
	rand.Shuffle(len(slice), func(i, j int) {
		slice[i], slice[j] = slice[j], slice[i]
	})
}

func distributeTeams(members []int64) map[int64]int64 {
	shuffle(members) // Рандомизируем порядок участников

	teamAssignments := make(map[int64]int64) // ID игрока -> номер команды
	currentTeam := int64(1)
	membersInCurrentTeam := int64(0)

	for _, member := range members {
		teamAssignments[member] = currentTeam
		membersInCurrentTeam++

		if membersInCurrentTeam == 6 {
			currentTeam++
			membersInCurrentTeam = 0
		}
	}

	// Обработка неполных команд
	if membersInCurrentTeam > 0 && membersInCurrentTeam < 4 && membersInCurrentTeam < currentTeam {
		// Перераспределение участников из последней неполной команды
		for member, team := range teamAssignments {
			if team == currentTeam {
				// Ищем команду для перераспределения
				for i := int64(1); i < currentTeam; i++ {
					if countMembers(teamAssignments, i) < 7 {
						teamAssignments[member] = i
						break
					}
				}
			}
		}
	} else if membersInCurrentTeam >= 4 {
		// Оставляем команду с 4 или 5 участниками как есть
		currentTeam++
	}

	return teamAssignments
}

func countMembers(assignments map[int64]int64, team int64) int {
	count := 0
	for _, assignedTeam := range assignments {
		if assignedTeam == team {
			count++
		}
	}
	return count
}
