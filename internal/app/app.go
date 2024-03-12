package app

import (
	bots2 "bot21/internal/bots"
	"bot21/internal/config"
	"bot21/internal/postgres"
	repoP "bot21/internal/repository/postgres"
	"bot21/internal/service"
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

	s := service.NewService(postgresRepo, bots, cfg.App.SuperUserID)

	go func() { s.Run() }()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigChan
	fmt.Println("Received signal:", sig)

	return nil
}
