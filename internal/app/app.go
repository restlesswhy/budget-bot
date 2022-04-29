package app

import (
	"bot/internal"
	"bot/internal/models"
	"fmt"
	"os"
	"strconv"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
)

var buttons = tgbotapi.NewReplyKeyboard(
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

var inlineButtons = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("street food", "1"),
		tgbotapi.NewInlineKeyboardButtonData("home food", "2"),
		tgbotapi.NewInlineKeyboardButtonData("home stuff", "3"),
	),
)

type App struct {
	bot   *tgbotapi.BotAPI
	close chan struct{}
	repo  internal.Repository

	wg sync.WaitGroup
}

func NewApp(repo internal.Repository) *App {
	fmt.Println(os.Getenv("TELEGRAM_APITOKEN"))
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_APITOKEN"))
	if err != nil {
		panic(err)
	}

	bot.Debug = true

	app := &App{
		bot:   bot,
		close: make(chan struct{}),
		repo:  repo,
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

				switch u.Message.Text {
				// case "hello":
				// 	if err := a.WriteMessage(u.Message); err != nil {
				// 		logrus.Warn(fmt.Sprintf("err write msg: %v", err))
				// 	}

				default:
					_, err := strconv.Atoi(u.Message.Text)
					if err != nil {
						answ := tgbotapi.NewMessage(u.Message.Chat.ID, "wrong input format, buddy :(")

						if _, err := a.bot.Send(answ); err != nil {
							logrus.Fatal(fmt.Sprintf("err send msg: %v", err))
						}

						continue
					}

					if err := a.WriteMessage(u.Message); err != nil {
						logrus.Warn(fmt.Sprintf("err write msg: %v", err))
					}

					answ := tgbotapi.NewMessage(u.Message.Chat.ID, "choose category")
					answ.ReplyMarkup = inlineButtons
					answ.ReplyToMessageID = u.Message.MessageID
					
					send, err := a.bot.Send(answ)
					if err != nil {
						logrus.Fatal(fmt.Sprintf("err send msg: %v", err))
					}

					if err := a.repo.WriteButton(&models.Buttons{
						ID: send.MessageID,
						MessageRelationID: u.Message.MessageID,
					}); err != nil {
						logrus.Fatal(fmt.Sprintf("err write button: %v", err))
					}
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
			// // msg := tgbotapi.NewMessage(b.ID, "a")
			// x := tgbotapi.NewRemoveKeyboard(true)

			// if _, err := a.bot.Send(x); err != nil {
			// 	logrus.Error(err)
			// }

			return
		}
	}
}

func (a *App) WriteMessage(msg *tgbotapi.Message) error {
	err := a.repo.WriteMessage(&models.Message{
		ID:        msg.MessageID,
		Text:      msg.Text,
		Firstname: msg.From.FirstName,
		Lastname:  msg.From.LastName,
		Username:  msg.From.UserName,
	})
	if err != nil {
		logrus.Fatal(fmt.Sprintf("err write msg: %v", err))
	}

	return err
}

func (a *App) Close() {
	close(a.close)
	a.wg.Wait()
}
