package main

import (
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)
func main(){
    err:=godotenv.Load()
    if err !=nil {
        log.Fatal("Error loading .env file")
    }
    bot, err:= tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_APITOKEN"))
    if err != nil {
        log.Panic(err)
    } 
    bot.Debug = true
    log.Printf("authorized with name %s", bot.Self.UserName)
    u:= tgbotapi.NewUpdate(0)
    u.Timeout=60
    updates:= bot.GetUpdatesChan(u)
    for update:=range updates {
        if update.Message == nil {
            continue
        }
        log.Printf("[%s]%s", update.Message.From.UserName, update.Message.Text)
        msg:= tgbotapi.NewMessage(update.Message.Chat.ID, "I hear you! You've written: " + update.Message.Text)
        bot.Send(msg)
    }
}