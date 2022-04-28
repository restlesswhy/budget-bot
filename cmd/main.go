package main

import (
	"bot/config"
	"bot/internal/app"
	"bot/internal/repository"
	"bot/pkg/postgres"
	"os"
	"os/signal"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func main() {
	if err := godotenv.Load(); err != nil {
		logrus.Fatal(err, "error loading env")
	}

	cfg, err := config.Load(os.Getenv("CONFIG"))
	if err != nil {
		logrus.Fatal(err, "error loading config")
	}

	pool, err := postgres.Connect(cfg)
	if err != nil {
		logrus.Fatal(err, "error connect to db")
	}
	defer pool.Close()

	repo := repository.NewDataRepo(pool)
	app := app.NewApp(repo)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	<-sig
	app.Close()
}
