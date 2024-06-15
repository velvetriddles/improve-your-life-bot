package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var gBot *tgbotapi.BotAPI
var gToken string
var gChatId int64

var gUsersInChat Users

type User struct {
	id    int64
	name  string
	coins uint16
}
type Users []*User

func init() {
	var err error

	_ = os.Setenv(TOKEN_NAME_IN_OS, "7379647501:AAF1FVgTfmIVpo_Qrh15DFJ99JqIJVKXAdI")
	if gToken = os.Getenv(TOKEN_NAME_IN_OS); gToken == "" {
		panic(fmt.Errorf(`failed to load environment variable "%s"`, TOKEN_NAME_IN_OS))
	}
	if gBot, err = tgbotapi.NewBotAPI(gToken); err != nil {
		log.Panic(err)
	}

	gBot.Debug = false
}

func isStartMessage(update *tgbotapi.Update) bool {
	return update.Message != nil && update.Message.Text == "/start"
}

func isCallbackQuery(update *tgbotapi.Update) bool {
	return update.CallbackQuery != nil && update.CallbackQuery.Data != ""
}

func delay(seconds uint8) {
	time.Sleep(time.Second * time.Duration(seconds))
}

func sendStringMessage(msg string) {
	gBot.Send(tgbotapi.NewMessage(gChatId, msg))
}

func sendMessageWithDelay(delayInSec uint8, message string) {
	sendStringMessage(message)
	delay(delayInSec)
}

func printIntro() {
	sendMessageWithDelay(2, "Hello! "+EMOJI_SUNGLASSES)
	sendMessageWithDelay(7, "There are numerous beneficial actions that, by performing regularly, we improve the quality of our life. However, often it's more fun, easier, or tastier to do something harmful. Isn't that so?")
	sendMessageWithDelay(7, "With greater likelihood, we'll prefer to get lost in YouTube Shorts instead of an English lesson, buy M&M's instead of vegetables, or lie in bed instead of doing yoga.")
	sendMessageWithDelay(1, EMOJI_SAD)
	sendMessageWithDelay(10, "Everyone has played at least one game where you need to level up a character, making them stronger, smarter, or more beautiful. It's enjoyable because each action brings results. In real life, though, systematic actions over time start to become noticeable. Let's change that, shall we?")
	sendMessageWithDelay(1, EMOJI_SMILE)
	sendMessageWithDelay(14, `Before you are two tables: "Useful Activities" and "Rewards". The first table lists simple short activities, and for completing each of them, you'll earn the specified amount of coins. In the second table, you'll see a list of activities you can only do after paying for them with coins earned in the previous step.`)
	sendMessageWithDelay(1, EMOJI_COIN)
	sendMessageWithDelay(10, `For example, you spend half an hour doing yoga, for which you get 2 coins. After that, you have 2 hours of programming study, for which you get 8 coins. Now you can watch 1 episode of "Interns" and break even. It's that simple!`)
	sendMessageWithDelay(6, `Mark completed useful activities to not lose your coins. And don't forget to "purchase" the reward before actually doing it.`)
}

func getKeyboardRow(buttonText, buttonCode string) []tgbotapi.InlineKeyboardButton {
	return tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(buttonText, buttonCode))

}

func askToPrintIntro() {
	msg := tgbotapi.NewMessage(gChatId, "In the introductory messages, you can find the purpose of this bot and the rules of the game. What do you think?")
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		getKeyboardRow(BUTTON_TEXT_PRINT_INTRO, BUTTON_CODE_PRINT_INTRO),
		getKeyboardRow(BUTTON_TEXT_SKIP_INTRO, BUTTON_CODE_SKIP_INTRO),
	)
	gBot.Send(msg)
}

func showMenu() {
	msg := tgbotapi.NewMessage(gChatId, "Choose one of the options:")
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		getKeyboardRow(BUTTON_TEXT_BALANCE, BUTTON_CODE_BALANCE),
		getKeyboardRow(BUTTON_TEXT_USEFUL_ACTIVITIES, BUTTON_CODE_USEFUL_ACTIVITIES),
		getKeyboardRow(BUTTON_TEXT_REWARDS, BUTTON_CODE_REWARDS),
	)
	gBot.Send(msg)
}

func showBalance(user *User) {
	msg := fmt.Sprintf("%s, your wallet is currently empty %s \nTrack a useful activity to earn coins", user.name, EMOJI_DONT_KNOW)
	if coins := user.coins; coins > 0 {
		msg = fmt.Sprintf("%s, you have %d %s", user.name, coins, EMOJI_COIN)
	}
	sendStringMessage(msg)
	showMenu()
}

func callbackQueryIsMissing(update *tgbotapi.Update) bool {
	return update.CallbackQuery == nil || update.CallbackQuery.From == nil
}
func getUserFromUpdate(update *tgbotapi.Update) (user *User, found bool) {
	if callbackQueryIsMissing(update) {
		return
	}

	userId := update.CallbackQuery.From.ID
	for _, userInChat := range gUsersInChat {
		if userId == userInChat.id {
			return userInChat, true
		}
	}
	return nil, false
}

func storeUserFromUpdate(update *tgbotapi.Update) (user *User, found bool) {
	if callbackQueryIsMissing(update) {
		return
	}
	from := update.CallbackQuery.From
	user = &User{id: from.ID, name: strings.TrimSpace(from.FirstName + " " + from.LastName), coins: 0}
	gUsersInChat = append(gUsersInChat, user)
	return user, true
}

func updateProcessing(update *tgbotapi.Update) {
	user, found := getUserFromUpdate(update)
	if !found {
		if user, found = storeUserFromUpdate(update); !found{
			sendStringMessage("Unable to identify the user")
			return
		}
	}
	choiceCode := update.CallbackQuery.Data
	log.Printf("[%T] %s", time.Now(), choiceCode)

	switch choiceCode {
	case BUTTON_CODE_BALANCE:
		showBalance(user)
	case BUTTON_CODE_USEFUL_ACTIVITIES:
		// showUsefulActivities()
	case BUTTON_CODE_REWARDS:
		// showRewards()
	case BUTTON_CODE_PRINT_INTRO:
		printIntro()
		showMenu()
	case BUTTON_CODE_SKIP_INTRO:
		showMenu()
	}
}

func main() {

	log.Printf("Authorized on account %s", gBot.Self.UserName)

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = UPDATE_CONFIG_TIMEOUT

	updates := gBot.GetUpdatesChan(updateConfig)

	for update := range updates {
		if isCallbackQuery(&update) {
			updateProcessing(&update)
		} else if isStartMessage(&update) { // If we got a message
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			gChatId = update.Message.Chat.ID
			askToPrintIntro()
		}
	}
}
