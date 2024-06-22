package service

import (
	"barbot/internal/bots"
	repoP "barbot/internal/repository/postgres"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"runtime"
	"sync"
)

type Service struct {
	db       *repoP.RepositoryPostgres
	bots     *bots.Bots
	AdminID  int64
	BarmenID int64
	ChatID   int64
}

func NewService(
	db *repoP.RepositoryPostgres,
	bots *bots.Bots,
	adminID int64,
	barmenID int64,
	chatID int64,
) *Service {
	return &Service{db: db, bots: bots, AdminID: adminID, BarmenID: barmenID, ChatID: chatID}
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
	if update.FromChat().ID == s.ChatID {
		return
	}
	switch update.SentFrom().ID {
	case s.AdminID:
		s.handleAdmin(update)
	case s.BarmenID:
		s.handleBarmen(update)
	default:
		s.handleGuest(update)

	}
}
