package main

import (
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var bot *tgbotapi.BotAPI

func initBot() {
	var err error
	bot, err = tgbotapi.NewBotAPI(BotToken)
	if err != nil {
		log.Println(err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)
}

func msg(title, text string) error {
	msg := tgbotapi.NewMessage(ChatID, fmt.Sprintf("*%s*\n%s", title, text))
	msg.ParseMode = "MarkdownV2"
	_, err := bot.Send(msg)
	return err
}
