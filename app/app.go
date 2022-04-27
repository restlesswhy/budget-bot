package app

import (
	"os"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
)

var numericKeyboard = tgbotapi.NewReplyKeyboard(
    tgbotapi.NewKeyboardButtonRow(
        tgbotapi.NewKeyboardButton("1"),
        tgbotapi.NewKeyboardButton("2"),
        tgbotapi.NewKeyboardButton("3"),
    ),
    tgbotapi.NewKeyboardButtonRow(
        tgbotapi.NewKeyboardButton("4"),
        tgbotapi.NewKeyboardButton("5"),
        tgbotapi.NewKeyboardButton("6"),
    ),
)

type App struct {
	bot   *tgbotapi.BotAPI
	close chan struct{}

	wg sync.WaitGroup
}

func NewApp() *App {
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

	updateConfig := tgbotapi.NewUpdate(0)

	updateConfig.Timeout = 30

	updates := a.bot.GetUpdatesChan(updateConfig)

	for {
		select {
		case u := <-updates:
			if u.Message == nil {
				continue
			}
	
			msg := tgbotapi.NewMessage(u.Message.Chat.ID, u.Message.Text)
	
			switch u.Message.Text {
			case "open":
				msg.ReplyMarkup = numericKeyboard
			case "close":
				msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
			}
	
			if _, err := a.bot.Send(msg); err != nil {
				logrus.Error(err)
			}
		case <-a.close:
			return
		}
	}
}

func (a *App) Close() {
	close(a.close)
	a.wg.Wait()
}
