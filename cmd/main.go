package main

import (
	"fmt"
	"github.com/BurntSushi/toml"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"testIlnlineKbrd/internal/bot"
	"testIlnlineKbrd/internal/model"
	"testIlnlineKbrd/internal/utils"
)

type botToken struct {
	Token string
}

var db *gorm.DB
var telegramBot *bot.Bot

func main() {
	var err error
	dsn := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=%s",
		utils.GoDotEnvVariable("DB_IP"),
		5555,
		"postgres",
		utils.GoDotEnvVariable("DB_NAME"),
		utils.GoDotEnvVariable("DB_PASS"),
		"disable")
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("bd connect ffs error: ", err)
		return
	}

	db.AutoMigrate(model.User{})
	db.AutoMigrate(utils.ErrLogs{})
	db.AutoMigrate(model.FowardMessage{})
	var token botToken
	if _, err := toml.DecodeFile("config/token.toml", &token); err != nil {
		log.Fatal("Decoding token.toml error: ", err)
		return
	}
	newBot, err := tgbotapi.NewBotAPI(token.Token)
	if err != nil {
		log.Fatal("New bot API error: ", err)
		return
	}

	newBot.Debug = false

	telegramBot = bot.NewBot(newBot)

	for {
		err = telegramBot.Start(db)
		if err != nil {
			utils.UppendErrorWithPath(err, db)
		}
		log.Error("StartFirstBot die?")
	}

}
