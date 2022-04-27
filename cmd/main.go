package main

import (
	"bot/app"
	"os"
	"os/signal"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func main() {
	if errLoadEnv := godotenv.Load(); errLoadEnv != nil {
		logrus.Fatal(errLoadEnv, "error loading env variables")
	}

	app := app.NewApp()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	<-sig
	app.Close()
}