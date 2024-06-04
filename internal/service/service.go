package service

import (
	"barbot/internal/bots"
	repoP "barbot/internal/repository/postgres"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"math/rand"
	"sync"
	"time"
)

type Service struct {
	db      *repoP.RepositoryPostgres
	bots    *bots.Bots
	AdminId []int64
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

	for i := int64(0); i < 10; i++ {
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
			workerChannels[u.ID%10] <- update
		}
	}
	for _, ch := range workerChannels {
		close(ch)
	}
	wg.Wait()
}

func (s *Service) processUpdate(update tgbotapi.Update) {

	if update.Message != nil && update.Message.Command() == "start" {
		s.needRegistration(update)
		m := tgbotapi.NewMessage(update.Message.From.ID, "Поздравляю, ты в игре, жди дальнейших указаний!")
		s.bots.Bot.Send(m)
		return
	}

	if update.Message != nil && update.Message.Command() == "new" && (s.AdminId[1] == update.SentFrom().ID || s.AdminId[0] == update.SentFrom().ID) {
		users := s.db.GetUsers()
		teamAssignments := distributeTeams(users)
		for id, team := range teamAssignments {
			m1 := tgbotapi.NewMessage(id, fmt.Sprintf("Твоя команда Тихоходок номер: %d ❤️️", team))
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

	//if update.Message != nil && update.Message.Photo != nil {
	//	photo := update.Message.Photo[len(update.Message.Photo)-1]
	//	p := tgbotapi.NewPhoto(7118228041, tgbotapi.FileID(photo.FileID))
	//	s.bots.Bot.Send(p)
	//}

}

func (s *Service) needRegistration(update tgbotapi.Update) {
	if u := update.SentFrom(); u != nil {
		s.db.SaveUser(u.ID)
	}
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
