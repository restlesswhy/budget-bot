package app

import (
	"bot/internal"
	"bot/internal/models"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
)

const (
	StreetFood = "1"
	HomeFood   = "2"
	HomeStuff  = "3"
	Income     = "4"
)

var ctgs map[string]string = map[string]string{
	StreetFood: "Street Food",
	HomeFood:   "Home Food",
	HomeStuff:  "Home Stuff",
	Income:     "Income",
}

var buttons = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("get monthly report"),
		tgbotapi.NewKeyboardButton("get day report"),
	),
)

var inlineButtons = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("street food", StreetFood),
		tgbotapi.NewInlineKeyboardButtonData("home food", HomeFood),
		tgbotapi.NewInlineKeyboardButtonData("home stuff", HomeStuff),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("income", Income),
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
				case "/start":

					answ := tgbotapi.NewMessage(u.Message.Chat.ID, "welcome to ur finance friend!")
					answ.ReplyMarkup = buttons

					if _, err := a.bot.Send(answ); err != nil {
						logrus.Fatal(fmt.Sprintf("err send msg: %v", err))
					}

				case "get monthly report":
					answ := tgbotapi.NewMessage(u.Message.Chat.ID, "wait a minuttee... maybe more :(")
					answ.ReplyMarkup = buttons

					if _, err := a.bot.Send(answ); err != nil {
						logrus.Fatal(fmt.Sprintf("err send msg: %v", err))
					}

				case "get day report":
					answ := tgbotapi.NewMessage(u.Message.Chat.ID, "oh, sorry, buddy, we r working on it :(")
					answ.ReplyMarkup = buttons

					if _, err := a.bot.Send(answ); err != nil {
						logrus.Fatal(fmt.Sprintf("err send msg: %v", err))
					}

				default:
					amount, err := strconv.Atoi(u.Message.Text)
					if err != nil {
						answ := tgbotapi.NewMessage(u.Message.Chat.ID, "wrong input format, buddy. only integers are support :(")

						if _, err := a.bot.Send(answ); err != nil {
							logrus.Fatal(fmt.Sprintf("err send msg: %v", err))
						}

						continue
					}

					answ := tgbotapi.NewMessage(u.Message.Chat.ID, "choose category")
					answ.ReplyMarkup = inlineButtons

					btn, err := a.bot.Send(answ)
					if err != nil {
						logrus.Fatal(fmt.Sprintf("err send msg: %v", err))
					}

					if err := a.repo.WriteButton(&models.Buttons{
						ID:        btn.MessageID,
						MessageID: u.Message.MessageID,
						Amount:    amount,
						Firstname: u.Message.From.FirstName,
						Lastname:  u.Message.From.LastName,
						Username:  u.Message.From.UserName,
					}); err != nil {
						logrus.Fatal(fmt.Sprintf("err write button: %v", err))
					}
				}

			} else if u.CallbackQuery != nil {

				switch u.CallbackQuery.Data {
				case StreetFood:
					if err := a.WriteTransaction(ctgs[StreetFood], u.CallbackQuery); err != nil {
						logrus.Fatal(fmt.Sprintf("err write tx: %v", err))
					}

				case HomeFood:
					if err := a.WriteTransaction(ctgs[HomeFood], u.CallbackQuery); err != nil {
						logrus.Fatal(fmt.Sprintf("err write tx: %v", err))
					}

				case HomeStuff:
					if err := a.WriteTransaction(ctgs[HomeStuff], u.CallbackQuery); err != nil {
						logrus.Fatal(fmt.Sprintf("err write tx: %v", err))
					}

				case Income:
					if err := a.WriteTransaction(ctgs[Income], u.CallbackQuery); err != nil {
						logrus.Fatal(fmt.Sprintf("err write tx: %v", err))
					}

				default:
					answ := tgbotapi.NewMessage(u.CallbackQuery.Message.Chat.ID, "we working on it...")

					_, err := a.bot.Send(answ)
					if err != nil {
						logrus.Fatal(fmt.Sprintf("err send msg: %v", err))
					}
					
					continue
				}
				// Respond to the callback query, telling Telegram to show the user
				// a message with the data received.
				// callback := tgbotapi.NewCallback(u.CallbackQuery.ID, u.CallbackQuery.Data)
				// if _, err := a.bot.Request(callback); err != nil {
				// 	panic(err)
				// }

				// And finally, send a message containing the data received.
				// msg := tgbotapi.NewMessage(u.CallbackQuery.Message.Chat.ID, "alll rightt!!! succesfully add this transaction :)")
				// if _, err := a.bot.Send(msg); err != nil {
				// 	panic(err)
				// }
				edit := tgbotapi.NewEditMessageTextAndMarkup(
					u.CallbackQuery.Message.Chat.ID,
					u.CallbackQuery.Message.MessageID,
					"alll rightt!!! succesfully add this transaction :)",
					tgbotapi.NewInlineKeyboardMarkup(
						tgbotapi.NewInlineKeyboardRow(
							tgbotapi.NewInlineKeyboardButtonData("cancel txn", "5"),
						)))
				_, err := a.bot.Request(edit)
				if err != nil {
					logrus.Fatal(fmt.Sprintf("err edit message tx: %v", err))
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

func (a *App) WriteTransaction(ctg string, data *tgbotapi.CallbackQuery) error {
	err := a.repo.WriteTransaction(&models.Transaction{
		ButtonID: data.Message.MessageID,
		Category: ctg,
		Time:     time.Now().UTC(),
	})
	if err != nil {
		logrus.Fatal(fmt.Sprintf("err write msg: %v", err))
	}

	return nil
}

func (a *App) Close() {
	close(a.close)
	a.wg.Wait()
}
