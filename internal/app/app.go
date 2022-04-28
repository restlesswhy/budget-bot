package app

import (
	"bot/internal"
	"os"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var numericKeyboard = tgbotapi.NewInlineKeyboardMarkup(
    tgbotapi.NewInlineKeyboardRow(
        tgbotapi.NewInlineKeyboardButtonURL("1.com", "http://1.com"),
        tgbotapi.NewInlineKeyboardButtonData("2", "hello"),
        tgbotapi.NewInlineKeyboardButtonData("3", "3"),
    ),
    tgbotapi.NewInlineKeyboardRow(
        tgbotapi.NewInlineKeyboardButtonData("4", "4"),
        tgbotapi.NewInlineKeyboardButtonData("5", "5"),
        tgbotapi.NewInlineKeyboardButtonData("6", "6"),
    ),
)

type App struct {
	bot   *tgbotapi.BotAPI
	close chan struct{}
	repo internal.Repository

	wg sync.WaitGroup
}

func NewApp(repo internal.Repository) *App {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_APITOKEN"))
	if err != nil {
		panic(err)
	}

	bot.Debug = true

	app := &App{
		bot: bot,
		close: make(chan struct{}),
		repo: repo,
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
			if u.Message != nil {
				// Construct a new message from the given chat ID and containing
				// the text that we received.
				msg := tgbotapi.NewMessage(u.Message.Chat.ID, u.Message.Text)
	
				// If the message was open, add a copy of our numeric keyboard.
				switch u.Message.Text {
				case "open":
					msg.ReplyMarkup = numericKeyboard
	
				}
	
				// Send the message.
				if _, err := a.bot.Send(msg); err != nil {
					panic(err)
				}
			} else if u.CallbackQuery != nil {
				// Respond to the callback query, telling Telegram to show the user
				// a message with the data received.
				callback := tgbotapi.NewCallback(u.CallbackQuery.ID, u.CallbackQuery.Data)
				if _, err := a.bot.Request(callback); err != nil {
					panic(err)
				}
	
				// And finally, send a message containing the data received.
				msg := tgbotapi.NewMessage(u.CallbackQuery.Message.Chat.ID, u.CallbackQuery.Data)
				if _, err := a.bot.Send(msg); err != nil {
					panic(err)
				}
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
