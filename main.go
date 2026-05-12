package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

type MainInfo struct {
	Temp float64 `json:"temp"`
}

type WeatherDescription struct {
	Description string `json:"description"`
}

type WeatherResponse struct {
	Main    MainInfo             `json:"main"`
	Weather []WeatherDescription `json:"weather"`
	Name    string               `json:"name"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_APITOKEN"))
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true
	log.Printf("authorized with name %s", bot.Self.UserName)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message == nil {
			continue
		}
		chatId := update.Message.Chat.ID
		if update.Message.Location != nil {
			lat := update.Message.Location.Latitude
			lon := update.Message.Location.Longitude
			weather, err := getWeather(lat, lon)
			if err != nil {
				log.Printf("Weather error: %v", err)
				weather = "I can't getweather, sorry!"
			}
			//reply:= fmt.Sprintf("I got your lan %f and lon %f", lat, lon)
			msg := tgbotapi.NewMessage(chatId, weather)
			bot.Send(msg)
			continue
		}
		log.Printf("[%s]%s", update.Message.From.UserName, update.Message.Text)
		btn := tgbotapi.NewKeyboardButtonLocation("Please send location")
		keyboard := tgbotapi.NewReplyKeyboard(tgbotapi.NewKeyboardButtonRow(btn))
		keyboard.ResizeKeyboard = true
		msg := tgbotapi.NewMessage(chatId, "Click button for weather")
		msg.ReplyMarkup = keyboard
		bot.Send(msg)
	}
}
func getWeather(lat, lon float64) (string, error) {
	apiKey := os.Getenv("API_KEY")
	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?lat=%f&lon=%f&appid=%s&units=metric", lat, lon, apiKey)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var data WeatherResponse
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return "", err
	}
	if len(data.Weather) > 0 {
		description := data.Weather[0].Description
		return fmt.Sprintf("City: %s Temperature %.1f°C\nSky: %s", data.Name, data.Main.Temp, description), nil
	}
	return "Weather data is empty", nil
}
