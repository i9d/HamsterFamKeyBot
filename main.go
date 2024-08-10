package main

import (
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"log"
	"os"
)

var _mainChannelID = os.Getenv("TELEGRAM_MAIN_CHANNEL_ID")
var _enChannelID = os.Getenv("TELEGRAM_EN_CHANNEL_ID")
var _bizonChannelID = os.Getenv("TELEGRAM_BIZON_CHANNEL_ID")
var _playgroundChannelID = os.Getenv("TELEGRAM_PLAYGROUND_CHANNEL_ID")
var _mainChatID = os.Getenv("TELEGRAM_MAIN_CHAT_ID")
var _walletAddress = os.Getenv("WALLET_ADDRESS")

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	log.Println("Test")
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")

	db, err := initializeDB()
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	loadTranslations()

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			handleMessages(update, bot, db)
		}
	}
}
