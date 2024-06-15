package app

import (
	bots2 "barbot/internal/bots"
	"barbot/internal/config"
	"barbot/internal/postgres"
	repoP "barbot/internal/repository/postgres"
	"barbot/internal/service"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func Run() error {
	if len(os.Args) < 2 {
		return errors.New("need file path")
	}

	cfg, err := config.New(os.Args[1])
	if err != nil {
		return err
	}

	bots, err := bots2.New(cfg.Telegram)
	if err != nil {
		return err
	}

	db, err := postgres.New(cfg.Database)
	if err != nil {
		return err
	}
	postgresRepo := repoP.New(db)

	s := service.NewService(postgresRepo, bots, cfg.App.AdminID, cfg.App.BarmenID, cfg.App.ChatID)

	go func() { s.Run() }()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigChan
	fmt.Println("Received signal:", sig)

	return nil
}
