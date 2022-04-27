package app

import (
	"os"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type App struct {
	bot   *tgbotapi.BotAPI
	close chan struct{}

	wg sync.WaitGroup
}

func CreateAndRun() *App {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_APITOKEN"))
	if err != nil {
		panic(err)
	}

	bot.Debug = true

	app := &App{
		bot: bot,
		close: make(chan struct{}),
	}

	app.wg.Add(1)
	go app.Run()

	return app
}

func (a *App) Run() {
	defer a.wg.Done()
}

func (a *App) Close() {
	close(a.close)
	a.wg.Wait()
}
