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
var gUsersInChat = make(map[int64]*User)
var gChatId int64

var gUsefulActivities = Activities{
	{"yoga", "Yoga (15 minutes)", 1},
	{"meditation", "Meditation (15 minutes)", 1},
	{"language", "Learning a foreign language (15 minutes)", 1},
	{"swimming", "Swimming (15 minutes)", 1},
	{"walk", "Walk (15 minutes)", 1},
	{"chores", "Chores", 1},
	{"work_learning", "Studying work materials (15 minutes)", 1},
	{"portfolio_work", "Working on a portfolio project (15 minutes)", 1},
	{"resume_edit", "Resume editing (15 minutes)", 1},
	{"creative", "Creative creation (15 minutes)", 1},
	{"reading", "Reading fiction literature (15 minutes)", 1},
}

var gRewards = Activities{
	{"watch_series", "Watching a series (1 episode)", 10},
	{"watch_movie", "Watching a movie (1 item)", 30},
	{"social_nets", "Browsing social networks (30 minutes)", 10},
	{"eat_sweets", "300 kcal of sweets", 60},
}

type User struct {
	id    int64
	name  string
	coins uint16
}

type Activity struct {
	code, name string
	coins      uint16
}

type Activities []*Activity

func init() {
	_ = os.Setenv(TOKEN_NAME_IN_OS, "7379647501:AAF1FVgTfmIVpo_Qrh15DFJ99JqIJVKXAdI")

	if gToken = os.Getenv(TOKEN_NAME_IN_OS); gToken == "" {
		panic(fmt.Errorf(`failed to load environment variable "%s"`, TOKEN_NAME_IN_OS))
	}

	var err error
	if gBot, err = tgbotapi.NewBotAPI(gToken); err != nil {
		log.Panic(err)
	}
	gBot.Debug = true
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

func sendStringMessage(chatId int64, msg string) {
	gBot.Send(tgbotapi.NewMessage(chatId, msg))
}

func sendMessageWithDelay(chatId int64, delayInSec uint8, message string) {
	sendStringMessage(chatId, message)
	delay(delayInSec)
}

func printIntro(chatId int64) {
	sendMessageWithDelay(chatId, 2, "Hello! "+EMOJI_SUNGLASSES)
	sendMessageWithDelay(chatId, 7, "There are numerous beneficial actions that, by performing regularly, we improve the quality of our life. However, often it's more fun, easier, or tastier to do something harmful. Isn't that so?")
	sendMessageWithDelay(chatId, 7, "With greater likelihood, we'll prefer to get lost in YouTube Shorts instead of an English lesson, buy M&M's instead of vegetables, or lie in bed instead of doing yoga.")
	sendMessageWithDelay(chatId, 1, EMOJI_SAD)
	sendMessageWithDelay(chatId, 10, "Everyone has played at least one game where you need to level up a character, making them stronger, smarter, or more beautiful. It's enjoyable because each action brings results. In real life, though, systematic actions over time start to become noticeable. Let's change that, shall we?")
	sendMessageWithDelay(chatId, 1, EMOJI_SMILE)
	sendMessageWithDelay(chatId, 14, `Before you are two tables: "Useful Activities" and "Rewards". The first table lists simple short activities, and for completing each of them, you'll earn the specified amount of coins. In the second table, you'll see a list of activities you can only do after paying for them with coins earned in the previous step.`)
	sendMessageWithDelay(chatId, 1, EMOJI_COIN)
	sendMessageWithDelay(chatId, 10, `For example, you spend half an hour doing yoga, for which you get 2 coins. After that, you have 2 hours of programming study, for which you get 8 coins. Now you can watch 1 episode of "Interns" and break even. It's that simple!`)
	sendMessageWithDelay(chatId, 6, `Mark completed useful activities to not lose your coins. And don't forget to "purchase" the reward before actually doing it.`)
}

func getKeyboardRow(buttonText, buttonCode string) []tgbotapi.InlineKeyboardButton {
	return tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(buttonText, buttonCode))
}

func askToPrintIntro(chatId int64) {
	msg := tgbotapi.NewMessage(chatId, "In the introductory messages, you can find the purpose of this bot and the rules of the game. What do you think?")
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		getKeyboardRow(BUTTON_TEXT_PRINT_INTRO, BUTTON_CODE_PRINT_INTRO),
		getKeyboardRow(BUTTON_TEXT_SKIP_INTRO, BUTTON_CODE_SKIP_INTRO),
	)
	gBot.Send(msg)
}

func showMenu(chatId int64) {
	msg := tgbotapi.NewMessage(chatId, "Choose one of the options:")
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		getKeyboardRow(BUTTON_TEXT_BALANCE, BUTTON_CODE_BALANCE),
		getKeyboardRow(BUTTON_TEXT_USEFUL_ACTIVITIES, BUTTON_CODE_USEFUL_ACTIVITIES),
		getKeyboardRow(BUTTON_TEXT_REWARDS, BUTTON_CODE_REWARDS),
	)
	gBot.Send(msg)
}

func showBalance(user *User) {
	msg := fmt.Sprintf("%s, your wallet is currently empty %s \nTrack a useful activity to earn coins", user.name, EMOJI_DONT_KNOW)
	if user.coins > 0 {
		msg = fmt.Sprintf("%s, you have %d %s", user.name, user.coins, EMOJI_COIN)
	}
	sendStringMessage(user.id, msg)
	showMenu(user.id)
}

func getUserFromUpdate(update *tgbotapi.Update) (*User, bool) {
	userId := update.CallbackQuery.From.ID
	user, found := gUsersInChat[userId]
	return user, found
}

func storeUserFromUpdate(update *tgbotapi.Update) (*User, bool) {
	from := update.CallbackQuery.From
	user := &User{id: from.ID, name: strings.TrimSpace(from.FirstName + " " + from.LastName), coins: 0}
	gUsersInChat[user.id] = user
	return user, true
}

func showActivities(chatId int64, activities Activities, message string, isUseful bool) {
	activitiesButtonsRows := make([]([]tgbotapi.InlineKeyboardButton), 0, len(activities)+1)
	for _, activity := range activities {
		activityDescription := ""
		if isUseful {
			activityDescription = fmt.Sprintf("+ %d %s: %s", activity.coins, EMOJI_COIN, activity.name)
		} else {
			activityDescription = fmt.Sprintf("- %d %s: %s", activity.coins, EMOJI_COIN, activity.name)
		}
		activitiesButtonsRows = append(activitiesButtonsRows, getKeyboardRow(activityDescription, activity.code))
	}
	activitiesButtonsRows = append(activitiesButtonsRows, getKeyboardRow(BUTTON_TEXT_PRINT_MENU, BUTTON_CODE_PRINT_MENU))

	msg := tgbotapi.NewMessage(chatId, message)
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(activitiesButtonsRows...)
	gBot.Send(msg)
}

func showUsefulActivities(chatId int64) {
	showActivities(chatId, gUsefulActivities, "Track a useful activity or return to the main menu:", true)
}

func showRewards(chatId int64) {
	showActivities(chatId, gRewards, "Purchase a reward or return to the main menu:", false)
}

func findActivity(activities Activities, choiceCode string) (*Activity, bool) {
	for _, activity := range activities {
		if choiceCode == activity.code {
			return activity, true
		}
	}
	return nil, false
}

func processUsefulActivity(activity *Activity, user *User) {
	errorMsg := ""
	if activity.coins == 0 {
		errorMsg = fmt.Sprintf(`the activity "%s" doesn't have a specified cost`, activity.name)
	} else if user.coins+activity.coins > MAX_USER_COINS {
		errorMsg = fmt.Sprintf("you cannot have more than %d %s", MAX_USER_COINS, EMOJI_COIN)
	}

	resultMessage := ""
	if errorMsg != "" {
		resultMessage = fmt.Sprintf("%s, I'm sorry, but %s %s Your balance remains unchanged.", user.name, errorMsg, EMOJI_SAD)
	} else {
		user.coins += activity.coins
		resultMessage = fmt.Sprintf(`%s, the "%s" activity is completed! %d %s has been added to your account. Keep it up! %s%s Now you have %d %s`,
			user.name, activity.name, activity.coins, EMOJI_COIN, EMOJI_BICEPS, EMOJI_SUNGLASSES, user.coins, EMOJI_COIN)
	}
	sendStringMessage(user.id, resultMessage)
}

func processReward(activity *Activity, user *User) {
	errorMsg := ""
	if activity.coins == 0 {
		errorMsg = fmt.Sprintf(`the reward "%s" doesn't have a specified cost`, activity.name)
	} else if user.coins < activity.coins {
		errorMsg = fmt.Sprintf(`you currently have %d %s. You cannot afford "%s" for %d %s`, user.coins, EMOJI_COIN, activity.name, activity.coins, EMOJI_COIN)
	}

	resultMessage := ""
	if errorMsg != "" {
		resultMessage = fmt.Sprintf("%s, I'm sorry, but %s %s Your balance remains unchanged, the reward is unavailable %s", user.name, errorMsg, EMOJI_SAD, EMOJI_DONT_KNOW)
	} else {
		user.coins -= activity.coins
		resultMessage = fmt.Sprintf(`%s, the reward "%s" has been paid for, get started! %d %s has been deducted from your account. Now you have %d %s`, user.name, activity.name, activity.coins, EMOJI_COIN, user.coins, EMOJI_COIN)
	}
	sendStringMessage(user.id, resultMessage)
}

func updateProcessing(update *tgbotapi.Update) {
	user, found := getUserFromUpdate(update)
	if !found {
		if user, found = storeUserFromUpdate(update); !found {
			sendStringMessage(update.CallbackQuery.Message.Chat.ID, "Unable to identify the user")
			return
		}
	}

	choiceCode := update.CallbackQuery.Data
	log.Printf("[%T] %s", time.Now(), choiceCode)

	switch choiceCode {
	case BUTTON_CODE_BALANCE:
		showBalance(user)
	case BUTTON_CODE_USEFUL_ACTIVITIES:
		showUsefulActivities(user.id)
	case BUTTON_CODE_REWARDS:
		showRewards(user.id)
	case BUTTON_CODE_PRINT_INTRO:
		printIntro(user.id)
		showMenu(user.id)
	case BUTTON_CODE_SKIP_INTRO:
		showMenu(user.id)
	case BUTTON_CODE_PRINT_MENU:
		showMenu(user.id)
	default:
		if usefulActivity, found := findActivity(gUsefulActivities, choiceCode); found {
			processUsefulActivity(usefulActivity, user)
			delay(2)
			showUsefulActivities(user.id)
			return
		}

		if reward, found := findActivity(gRewards, choiceCode); found {
			processReward(reward, user)
			delay(2)
			showRewards(user.id)
			return
		}

		log.Printf(`[%T] !!!!!!!!! ERROR: Unknown code "%s"`, time.Now(), choiceCode)
		msg := fmt.Sprintf("%s, I'm sorry, I don't recognize code '%s' %s Please report this error to my creator.", user.name, choiceCode, EMOJI_SAD)
		sendStringMessage(user.id, msg)
	}
}

func processUpdate(update tgbotapi.Update) {
	if isCallbackQuery(&update) {
		updateProcessing(&update)
	} else if isStartMessage(&update) {
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		gChatId = update.Message.Chat.ID
		from := update.Message.From
		gUsersInChat[from.ID] = &User{id: from.ID, name: strings.TrimSpace(from.FirstName + " " + from.LastName), coins: 0}
		askToPrintIntro(gChatId)
	}
}

func main() {
	log.Printf("Authorized on account %s", gBot.Self.UserName)

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = UPDATE_CONFIG_TIMEOUT

	updates := gBot.GetUpdatesChan(updateConfig)

	for update := range updates {
		go processUpdate(update)
	}
}
