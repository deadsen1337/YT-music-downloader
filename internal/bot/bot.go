package bot

import (
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"os"
	"strings"
	"testIlnlineKbrd/internal/bot/keyboard"
	"testIlnlineKbrd/internal/model"
	"testIlnlineKbrd/internal/utils"
)

type Bot struct {
	bot *tgbotapi.BotAPI
}

func NewBot(bot *tgbotapi.BotAPI) *Bot {
	return &Bot{bot: bot}
}

var messageMap = make(map[int64]string)

func (b *Bot) Start(db *gorm.DB) error {
	log.Printf("Authorized on account %s", b.bot.Self.UserName)

	updates, err := b.initUpdatesChannel()
	if err != nil {
		log.Fatal(err)
		return err
	}

	b.handleUpdates(updates, db)
	utils.UppendErrorWithPath(errors.New("выключился хандл апдейт цикл"), db)
	go utils.SendTextMsg(426010190, "Выключился хандл апдейт цикл", nil, 0, db)
	return nil
}

func (b *Bot) initUpdatesChannel() (tgbotapi.UpdatesChannel, error) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	return b.bot.GetUpdatesChan(u), nil
}

func (b *Bot) handleUpdates(updates tgbotapi.UpdatesChannel, db *gorm.DB) {
	go utils.MsgDealer(0, b.bot, &utils.SendMsgFirst, db)

	var myMsg model.FowardMessage
	myMsg.TextType = true
	myMsg.BotID = 0
	myMsg.Message = "Это сообщение отправляется каждый раз как я запускаюсь.\nОно не всегда говорит об ощибке, возможно я просто был перезапущен."
	myMsg.ChatID = 426010190
	db.Create(&myMsg)

	for update := range updates {
		MFID := update.Message.From.ID
		if update.Message != nil {
			if update.Message.From.IsBot {
				continue
			}
			var user model.User
			err := db.Model(model.User{}).Where("user_id = ?", update.Message.From.ID).First(&user).Error
			if err != nil {
				if err == gorm.ErrRecordNotFound {
					user.UserId = update.Message.From.ID
					user.UserName = update.Message.From.UserName
					err = db.Model(model.User{}).Create(&user).Error
					if err != nil {
						utils.UppendErrorWithPath(err, db)
					}
				} else {
					utils.UppendErrorWithPath(err, db)
					continue
				}
			}
			if user.UserName != update.Message.From.UserName {
				user.OldUsername = append(user.OldUsername, user.UserName)
				user.UserName = update.Message.From.UserName
				err = db.Model(model.User{}).Where("user_id = ?", MFID).Save(&user).Error
				if err != nil {
					utils.UppendErrorWithPath(err, db)
				}

			}
			if update.Message.IsCommand() {

				switch update.Message.Command() {
				case "start":
					go utils.SendTextMsg(MFID, "Hi! I can download music from youtube", keyboard.MainMenuKeyboard, 0, db)
					continue
				}
			}
			switch update.Message.Text {
			case "Download sound":
				if user.Root {
					go utils.SendTextMsg(MFID, "Send YouTube Link", nil, 0, db)
					messageMap[MFID] = "yt_link"
				}
				continue
			}

			switch messageMap[MFID] {
			case "yt_link":
				filename, err := utils.YouTubeDownload(strings.TrimSpace(update.Message.Text))
				if err != nil {
					utils.UppendErrorWithPath(err, db)
					go utils.SendTextMsg(MFID, "Something went wrong..", nil, 0, db)
					messageMap[MFID] = ""
					continue
				}

				file, err := os.ReadFile(filename)
				if err != nil {
					utils.UppendErrorWithPath(err, db)
					go utils.SendTextMsg(MFID, "Something went wrong..", nil, 0, db)
					messageMap[MFID] = ""
					continue
				}

				fileBytes := tgbotapi.FileBytes{
					Name:  filename,
					Bytes: file,
				}

				msg := tgbotapi.NewAudio(MFID, fileBytes)
				msg.Caption = "Here's your audio"
				_, err = b.bot.Send(msg)
				if err != nil {
					utils.UppendErrorWithPath(err, db)
					go utils.SendTextMsg(MFID, "Something went wrong..", nil, 0, db)
					messageMap[MFID] = ""
					continue
				}
				messageMap[MFID] = ""
				continue
			}
		}
	}
}
