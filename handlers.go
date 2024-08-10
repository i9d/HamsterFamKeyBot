package main

import (
	"database/sql"
	"log"
	"strconv"
	"time"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func handleMessages(update tgbotapi.Update, bot *tgbotapi.BotAPI, db *sql.DB) {
	if update.Message.Chat.Type != "private" {
		log.Printf("Received message from non-private chat: %v", update.Message.Chat.Type)
		return
	}

	user := User{
		ID:       update.Message.From.ID,
		ChatID:   update.Message.Chat.ID,
		Username: update.Message.From.UserName,
		Name:     update.Message.From.FirstName + " " + update.Message.From.LastName,
	}
	if update.Message.From.LanguageCode == "ru" {
		user.Lang = update.Message.From.LanguageCode
	} else if update.Message.From.LanguageCode == "uk" {
		user.Lang = update.Message.From.LanguageCode
	} else {
		user.Lang = "en"
	}

	// Handle commands
	if update.Message.IsCommand() {
		switch update.Message.Command() {
		case "start":
			startCommand(user, bot)
			//miniAppKeyboard(user, bot)
		}
	}

	// Handle messages
	switch update.Message.Text {

	case "Получить код":
		sendUpdateMessage(user, bot)
	case "Получить коды":
		sendUpdateMessage(user, bot)
	case "Отримати коди":
		sendUpdateMessage(user, bot)
	case "Get codes":
		sendUpdateMessage(user, bot)
	case "Bike Ride 3D":
		user.Game = "bike"
		getCodesCommand(user, bot, db)
	case "My Clone Army":
		user.Game = "clone"
		getCodesCommand(user, bot, db)
	case "Chain Cube 2048":
		user.Game = "cube"
		getCodesCommand(user, bot, db)
	case "Train Miner":
		user.Game = "train"
		getCodesCommand(user, bot, db)
	}
}

func miniAppKeyboard(user User, bot *tgbotapi.BotAPI) {
	message := tgbotapi.NewMessage(user.ChatID, getTranslation(user.Lang, "start_message"))
	miniappButton := tgbotapi.NewInlineKeyboardButtonURL(getTranslation(user.Lang, "get_codes"), "https://t.me/hamster_fam_bot?startapp=hamster_fam_bot")
	var inlineKeyboardRows [][]tgbotapi.InlineKeyboardButton
	var inlineKeyboard tgbotapi.InlineKeyboardMarkup
	inlineKeyboardRows = append(inlineKeyboardRows, tgbotapi.NewInlineKeyboardRow(miniappButton))
	inlineKeyboard = tgbotapi.NewInlineKeyboardMarkup(inlineKeyboardRows...)
	message.ReplyMarkup = inlineKeyboard
	bot.Send(message)
	message2 := tgbotapi.NewMessage(user.ChatID, getTranslation(user.Lang, "donate"))
	message2.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	bot.Send(message2)
	bot.Send(tgbotapi.NewMessage(user.ChatID, _walletAddress))
}

func startCommand(user User, bot *tgbotapi.BotAPI) {
	// Greeting message
	greeting := tgbotapi.NewMessage(user.ChatID, getTranslation(user.Lang, "start_message"))
	greeting.ReplyMarkup = createKeyboard(user)
	bot.Send(greeting)
	// Subscribe requirement message
	requirement := tgbotapi.NewMessage(user.ChatID, getTranslation(user.Lang, "subscribe_message"))
	requirement.ReplyMarkup = getSubscribeButtons(user)
	bot.Send(requirement)
}

func sendUpdateMessage(user User, bot *tgbotapi.BotAPI) {
	bot.Send(tgbotapi.NewMessage(user.ChatID, getTranslation(user.Lang, "start_again")))
}

func getCodesCommand(user User, bot *tgbotapi.BotAPI, db *sql.DB) {
	var isSubscriber bool
	switch user.Lang {
	case "en":
		isSubscriber = checkSubscription(_enChannelID, user.ChatID, bot)
	case "ru":
		isSubscriber = checkSubscription(_bizonChannelID, user.ChatID, bot) && checkSubscription(_mainChannelID, user.ChatID, bot)
	case "uk":
		isSubscriber = checkSubscription(_bizonChannelID, user.ChatID, bot) && checkSubscription(_mainChannelID, user.ChatID, bot)
	}

	if isSubscriber {
		increasedLimit := 2
		// Increase Limit if member of the Playground channel
		isPlaygroundSubscriber := checkSubscription(_playgroundChannelID, user.ChatID, bot)
		if isPlaygroundSubscriber {
			increasedLimit += 2
		}
		// Increase Limit if member of the main chat
		isChatSubscriber := checkSubscription(_mainChatID, user.ChatID, bot)
		if isChatSubscriber {
			increasedLimit += 2
		}
		// Increase Limit if member of the main channel
		isMainChannelSubscriber := checkSubscription(_mainChannelID, user.ChatID, bot)
		if isMainChannelSubscriber {
			increasedLimit += 2
		}

		codesCount := availableCodes(&user, increasedLimit, db)
		if codesCount < 1 {
			sendMessage(user, getTranslation(user.Lang, "limit_reached"), bot)
		} else {
			if codesCount > 2 {
				codesCount = 2
			}
			codes := getCodes(db, codesCount, user.Game)

			if len(codes) == 0 {
				sendMessage(user, getTranslation(user.Lang, "no_codes"), bot)
			} else {
				markCodesUsed(codes, db)
				user.AlreadyGot += len(codes)
				updateDailyLimit(user.ChatID, user.AlreadyGot, db)
				for _, code := range codes {
					code.UsedBy = user.ID
					bot.Send(tgbotapi.NewMessage(user.ChatID, code.Code))
				}
				if user.UserLimit <= 2 {
					switch user.Lang {
					case "en":
						if !checkSubscription(_playgroundChannelID, user.ChatID, bot) {
							sendMessage(user, getTranslation(user.Lang, "additional_codes"), bot)
						}
					default:
						if !checkSubscription(_playgroundChannelID, user.ChatID, bot) || !checkSubscription(_mainChatID, user.ChatID, bot) {
							sendMessage(user, getTranslation(user.Lang, "additional_codes"), bot)
						}
					}
				} else {
					bot.Send(tgbotapi.NewMessage(user.ChatID, getTranslation(user.Lang, "donate")))
					bot.Send(tgbotapi.NewMessage(user.ChatID, _walletAddress))
				}
			}
		}
	} else {
		// Subscribe requirement message
		requirement := tgbotapi.NewMessage(user.ChatID, getTranslation(user.Lang, "not_subscribed"))
		requirement.ReplyMarkup = getSubscribeButtons(user)
		bot.Send(requirement)
	}
}

func createKeyboard(user User) tgbotapi.ReplyKeyboardMarkup {
	//getCodeButton := tgbotapi.NewKeyboardButton(getTranslation(user.Lang, "get_codes"))
	getBikeButton := tgbotapi.NewKeyboardButton("Bike Ride 3D")
	getCubeButton := tgbotapi.NewKeyboardButton("Chain Cube 2048")
	getTrainButton := tgbotapi.NewKeyboardButton("Train Miner")
	getCloneButton := tgbotapi.NewKeyboardButton("My Clone Army")

	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(getBikeButton, getCubeButton),
		tgbotapi.NewKeyboardButtonRow(getTrainButton, getCloneButton),
	)
	return keyboard
}

func checkSubscription(channelID string, chatID int64, bot *tgbotapi.BotAPI) bool {
	channelChatID, _ := strconv.ParseInt(channelID, 10, 64)

	memberConfig := tgbotapi.GetChatMemberConfig{
		ChatConfigWithUser: tgbotapi.ChatConfigWithUser{
			ChatID: channelChatID,
			UserID: chatID,
		},
	}
	member, err := bot.GetChatMember(memberConfig)
	if err != nil {
		log.Printf("Error checking membership: %v", err)
		return false
	}
	if member.Status == "member" || member.Status == "administrator" || member.Status == "creator" {
		return true
	}
	return false
}

func availableCodes(user *User, increasedLimit int, db *sql.DB) int {
	err := db.QueryRow("SELECT last_check, already_got, user_limit FROM users WHERE chat_id = $1", user.ChatID).Scan(&user.LastCheck, &user.AlreadyGot, &user.UserLimit)
	if err != nil && err != sql.ErrNoRows {
		log.Printf("Error checking user: %v", err)
		return 0
	}

	if err == sql.ErrNoRows {
		_, err = db.Exec("INSERT INTO users (chat_id, user_id, username, firstname, last_check) VALUES ($1, $2, $3, $4, $5)", user.ChatID, user.ID, user.Username, user.Name, time.Now())
		if err != nil {
			log.Printf("Error inserting new user: %v", err)
			return 0
		}
		return increasedLimit
	}

	if user.UserLimit == 0 {
		user.UserLimit = increasedLimit
	} else if increasedLimit < 8 {
		user.UserLimit = increasedLimit
	}

	if user.UserLimit > 16 {
		user.UserLimit = 16
	}

	// Compare with midnight UTC
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	if now.Before(startOfDay) {
		startOfDay = startOfDay.Add(-24 * time.Hour)
	}

	if user.LastCheck.Before(startOfDay) {
		updateDailyLimit(user.ChatID, 0, db)
		user.AlreadyGot = 0
	}
	if (user.UserLimit - user.AlreadyGot) >= 0 {
		return user.UserLimit - user.AlreadyGot
	} else {
		return 0
	}
}

func getCodes(db *sql.DB, limit int, game string) []Code {
	rows, err := db.Query("SELECT id, code FROM codes WHERE used = $1 AND game_type = $2 LIMIT $3", "false", game, limit)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var codes []Code
	for rows.Next() {
		var code Code
		if err := rows.Scan(&code.ID, &code.Code); err != nil {
			return nil
		}
		codes = append(codes, code)
	}

	return codes
}

func markCodesUsed(codes []Code, db *sql.DB) {
	for _, code := range codes {
		_, err := db.Exec("UPDATE codes SET used = TRUE, used_by = $1 WHERE id = $2", code.UsedBy, code.ID)
		if err != nil {
			log.Printf("Error marking code as used: %v", err)
		}
	}
}

func updateDailyLimit(chatID int64, gotCodesCount int, db *sql.DB) {
	// if not me
	if chatID != 307766739 {
		_, err := db.Exec("UPDATE users SET last_check = $1, already_got = $2 WHERE chat_id = $3", time.Now(), gotCodesCount, chatID)
		if err != nil {
			log.Printf("Error updating last check time: %v", err)
		}
	}
}

func sendMessage(user User, text string, bot *tgbotapi.BotAPI) {
	message := tgbotapi.NewMessage(user.ChatID, text)

	var inlineKeyboard tgbotapi.InlineKeyboardMarkup
	var inlineKeyboardRows [][]tgbotapi.InlineKeyboardButton

	if user.Lang == "ru" || user.Lang == "uk" {
		if !checkSubscription(_mainChannelID, user.ChatID, bot) {
			mainChannelButton := tgbotapi.NewInlineKeyboardButtonURL(getTranslation(user.Lang, "main_channel_ad"), getTranslation(user.Lang, "main_channel_link"))
			inlineKeyboardRows = append(inlineKeyboardRows, tgbotapi.NewInlineKeyboardRow(mainChannelButton))
		}
		if !checkSubscription(_mainChatID, user.ChatID, bot) {
			chatButton := tgbotapi.NewInlineKeyboardButtonURL(getTranslation(user.Lang, "chat_ad"), getTranslation(user.Lang, "chat_link"))
			inlineKeyboardRows = append(inlineKeyboardRows, tgbotapi.NewInlineKeyboardRow(chatButton))
		}
		if !checkSubscription(_playgroundChannelID, user.ChatID, bot) {
			playgroundButton := tgbotapi.NewInlineKeyboardButtonURL(getTranslation(user.Lang, "playground_channel_ad"), getTranslation(user.Lang, "playground_channel_link"))
			inlineKeyboardRows = append(inlineKeyboardRows, tgbotapi.NewInlineKeyboardRow(playgroundButton))
		}
	} else if user.Lang == "en" {
		if !checkSubscription(_playgroundChannelID, user.ChatID, bot) {
			playgroundButton := tgbotapi.NewInlineKeyboardButtonURL(getTranslation(user.Lang, "playground_channel_ad"), getTranslation(user.Lang, "playground_channel_link"))
			inlineKeyboardRows = append(inlineKeyboardRows, tgbotapi.NewInlineKeyboardRow(playgroundButton))
		}
	}
	if len(inlineKeyboardRows) > 0 {
		inlineKeyboard = tgbotapi.NewInlineKeyboardMarkup(inlineKeyboardRows...)
		message.ReplyMarkup = inlineKeyboard
		bot.Send(message)
	} else {
		bot.Send(message)
	}
}

func getSubscribeButtons(user User) tgbotapi.InlineKeyboardMarkup {
	switch user.Lang {
	case "ru":
		bizonChannelButton := tgbotapi.NewInlineKeyboardButtonURL(getTranslation(user.Lang, "bizon_channel_ad"), getTranslation(user.Lang, "bizon_channel_link"))
		mainChannelButton := tgbotapi.NewInlineKeyboardButtonURL(getTranslation(user.Lang, "main_channel_ad"), getTranslation(user.Lang, "main_channel_link"))
		//playgroundChannelButton := tgbotapi.NewInlineKeyboardButtonURL(getTranslation(user.Lang, "playground_channel_ad"), getTranslation(user.Lang, "playground_channel_link"))
		return tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(bizonChannelButton),
			tgbotapi.NewInlineKeyboardRow(mainChannelButton),
			//tgbotapi.NewInlineKeyboardRow(playgroundChannelButton),
		)
	case "uk":
		bizonChannelButton := tgbotapi.NewInlineKeyboardButtonURL(getTranslation(user.Lang, "bizon_channel_ad"), getTranslation(user.Lang, "bizon_channel_link"))
		mainChannelButton := tgbotapi.NewInlineKeyboardButtonURL(getTranslation(user.Lang, "main_channel_ad"), getTranslation(user.Lang, "main_channel_link"))
		//playgroundChannelButton := tgbotapi.NewInlineKeyboardButtonURL(getTranslation(user.Lang, "playground_channel_ad"), getTranslation(user.Lang, "playground_channel_link"))
		return tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(bizonChannelButton),
			tgbotapi.NewInlineKeyboardRow(mainChannelButton),
			//tgbotapi.NewInlineKeyboardRow(playgroundChannelButton),
		)
	default:
		mainChannelButton := tgbotapi.NewInlineKeyboardButtonURL(getTranslation(user.Lang, "main_channel_ad"), getTranslation(user.Lang, "main_channel_link"))
		playgroundChannelButton := tgbotapi.NewInlineKeyboardButtonURL(getTranslation(user.Lang, "playground_channel_ad"), getTranslation(user.Lang, "playground_channel_link"))
		return tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(mainChannelButton),
			tgbotapi.NewInlineKeyboardRow(playgroundChannelButton),
			//tgbotapi.NewInlineKeyboardRow(newsChannelButton),
		)
	}
}
