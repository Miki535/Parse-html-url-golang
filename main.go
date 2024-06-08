package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

func main() {
	botToken := ""

	bot, err := telego.NewBot(botToken, telego.WithDefaultDebugLogger())

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	updates, _ := bot.UpdatesViaLongPolling(nil)

	bh, _ := th.NewBotHandler(bot, updates)

	defer bh.Stop()
	defer bot.StopLongPolling()

	bh.Handle(func(bot *telego.Bot, update telego.Update) {
		if update.Message != nil {
			chatId := tu.ID(update.Message.Chat.ID)
			msg := tu.Message(
				chatId,
				"Привіт! За допомогою цього бота ви дізнаєтесь актуальний курс на такі крипто валюти як bitcoin та ethereum!",
			)
			bot.SendMessage(msg)
		}
	}, th.CommandEqual("start"))

	bh.Handle(func(bot *telego.Bot, update telego.Update) {
		if update.Message != nil {
			chatId := tu.ID(update.Message.Chat.ID)
			newMessage := update.Message.Text
			HTMLparse(bot, newMessage, chatId)
		}
	})

	bh.Start()
}

func HTMLparse(bot *telego.Bot, url string, chatId telego.ChatID) {
	resp, err := http.Get(url)
	if err != nil {
		bot.SendMessage(tu.Message(chatId, "Error1"))
		log.Println(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		bot.SendMessage(tu.Message(chatId, "Error2"))
		log.Println(err)
	}
	_, err = bot.SendMessage(tu.Message(chatId, string(body)))
	if err != nil {
		fmt.Println(err)
	}
}
