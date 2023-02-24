package utils

import (
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"runtime"
	"strconv"
	"strings"
	"testIlnlineKbrd/internal/model"
	"time"
)

var SendMsgFirst = make(chan tgbotapi.MessageConfig)

func MsgDealer(botID int, bot *tgbotapi.BotAPI, c *chan tgbotapi.MessageConfig, bd *gorm.DB) {
	for {
		select {
		case msg := <-*c:
			{
				if len(msg.Text) == 0 {
					log.Error("msg.ChatID ", msg.ChatID)
					continue
				}
				_, err := bot.Send(msg)
				if err != nil {
					log.Error(err)
					if err.Error() == "Bad Request: chat not found" {
						err1 := bd.Model(model.User{}).Where("user_id = ?", msg.ChatID).Update("blocked_us", true).Error
						if err1 != nil {
							UppendErrorWithPath(err1, bd)
						}
					} else if err.Error() == "Forbidden: bot was blocked by the user" {
						err1 := bd.Model(model.User{}).Where("user_id = ?", msg.ChatID).Update("blocked_us", true).Error
						if err1 != nil {
							UppendErrorWithPath(err1, bd)
						}
						// clearUser(msg.ChatID)
					} else if err.Error() == "Forbidden: user is deactivated" {
						err1 := bd.Model(model.User{}).Where("user_id = ?", msg.ChatID).Update("blocked_us", true).Error
						if err1 != nil {
							UppendErrorWithPath(err1, bd)
						}
						// clearUser(msg.ChatID)
					} else if err.Error() == "Bad Request: message is too long" {

					} else if err.Error() == "Gateway Timeout" {
					} else if strings.Contains(err.Error(), "Too Many Requests") {
						time.Sleep(500 * time.Millisecond)
					}
				}
			}
		default:
			{
				var fowardMsg model.FowardMessage
				resulterr := bd.Model(model.FowardMessage{}).Where("bot_id = ?", botID).Order("created_at ASC").First(&fowardMsg)
				if resulterr.Error != nil {
					if resulterr.Error.Error() == "record not found" {
						time.Sleep(time.Duration(2000) * time.Millisecond)
						continue
					} else {
						time.Sleep(time.Duration(200) * time.Millisecond)
						if resulterr.Error != gorm.ErrRecordNotFound {
							UppendErrorWithPath(resulterr.Error, bd)
						}
						continue
					}
				} else {
					var err error
					if fowardMsg.TextType {
						msghi := tgbotapi.NewMessage(int64(fowardMsg.ChatID), fowardMsg.Message)
						time.Sleep(time.Duration(550) * time.Millisecond)
						_, err = bot.Send(msghi)
					} else {
						fwdmsg := tgbotapi.NewForward(fowardMsg.ChatID, fowardMsg.FromChatID, fowardMsg.MessageID)

						time.Sleep(time.Duration(550) * time.Millisecond)
						_, err = bot.Send(fwdmsg)
					}

					if err != nil {
						log.Error(err)
						if err.Error() == "Bad Request: chat not found" {
							err1 := bd.Model(model.User{}).Where("user_id = ?", fowardMsg.ChatID).Update("blocked_us", true).Error
							if err1 != nil {
								UppendErrorWithPath(err1, bd)
							}
							bd.Unscoped().Delete(&fowardMsg)
						} else if err.Error() == "Forbidden: user is deactivated" {
							err1 := bd.Model(model.User{}).Where("user_id = ?", fowardMsg.ChatID).Update("blocked_us", true).Error
							if err1 != nil {
								UppendErrorWithPath(err1, bd)
							}
							// clearUser(fowardMsg.ChatID)
							bd.Unscoped().Delete(&fowardMsg)
						} else if strings.Contains(err.Error(), "message text is empty") {
							log.Error("message text is empty,", fowardMsg.ChatID)
							if fowardMsg.MessageID != 0 {
								var try []model.FowardMessage
								bd.Where("from_chat_id = ? and message_id = ?", fowardMsg.ChatID, fowardMsg.MessageID, botID).Find(&try)
								for _, tr := range try {
									bd.Unscoped().Delete(&tr)
								}
							}
							var try2 []model.FowardMessage
							bd.Where("message_id = ? and message = ?", 0, "").Find(&try2)
							for _, tr := range try2 {
								bd.Unscoped().Delete(&tr)
							}
						} else if err.Error() == "Forbidden: bot can't initiate conversation with a user" {
							err1 := bd.Model(model.User{}).Where("user_id = ?", fowardMsg.ChatID).Update("blocked_us", true).Error
							if err1 != nil {
								UppendErrorWithPath(err1, bd)
							}
							// if fowardMsg.TextType {
							// 	msghi := tgbotapi.NewMessage(int64(fowardMsg.ChatID), fowardMsg.Message)
							// 	msghi.ReplyMarkup = cancel
							// 	time.Sleep(time.Duration(550) * time.Millisecond)
							// 	_, err = bot.Send(msghi)
							// }
							bd.Unscoped().Delete(&fowardMsg)
						} else if err.Error() == "Forbidden: bot was blocked by the user" {
							err1 := bd.Model(model.User{}).Where("user_id = ?", fowardMsg.ChatID).Update("blocked_us", true).Error
							if err1 != nil {
								UppendErrorWithPath(err1, bd)
							}
							// clearUser(fowardMsg.ChatID)
							bd.Unscoped().Delete(&fowardMsg)
						} else if err.Error() == "Forbidden: bot can't send messages to bots" {

							bd.Unscoped().Delete(&fowardMsg)
						}
						bd.Unscoped().Delete(&fowardMsg)
					} else {
						err = bd.Unscoped().Delete(&fowardMsg).Error
						if err != nil {
							time.Sleep(time.Duration(650) * time.Millisecond)
							continue
						}
					}
				}
				time.Sleep(time.Duration(300) * time.Millisecond)
			}
		}
		time.Sleep(time.Duration(650) * time.Millisecond)
	}
}

func SendTextMsg(MFID int64, text string, mybtn interface{}, botID int, bd *gorm.DB) {
	if len(strings.TrimSpace(text)) == 0 {
		return
	}
	msg := tgbotapi.NewMessage(int64(MFID), text)
	if mybtn != nil {
		msg.ReplyMarkup = mybtn
	}
	msg.DisableWebPagePreview = true
	switch botID {
	case 0:
		{
			go uppendTextMsg(msg, &SendMsgFirst, bd)
		}
	default:
		{
			go UppendErrorWithPath(errors.New("Какой то неизвестный айди бота оО? "+strconv.Itoa(botID)), bd)
		}
	}
}
func uppendTextMsg(msg tgbotapi.MessageConfig, myChan *chan tgbotapi.MessageConfig, bd *gorm.DB) {
	var nextmsg tgbotapi.MessageConfig
	var nexFlag = false
	if len(msg.Text) > 3800 {
		mytext := strings.Split(msg.Text, "\n")
		chatID := msg.ChatID
		nexFlag = true
		myLen := len(mytext)
		myLen /= 2
		firstText := strings.Join(mytext[myLen:], "\n")
		secondText := strings.Join(mytext[:myLen], "\n")

		nextmsg = tgbotapi.NewMessage(chatID, firstText)
		nextmsg.ReplyMarkup = msg.ReplyMarkup
		msg = tgbotapi.NewMessage(chatID, secondText)
	}

	//////////
	if len(msg.Text) == 0 {
		badway := ""
		for stepCount := 1; stepCount <= 10; stepCount++ {
			// получаем через рантайм указатель на положение stepCount в цепочке шаг наверх, а так же строку вызова внутри функции.
			pc, _, line, _ := runtime.Caller(stepCount)
			// получаем через указатель полный путь к функции записанный на этом шагу
			fullFuncPath := runtime.FuncForPC(pc).Name()
			// дробим строку через точку, что бы получить массив строк, в котором последним элементом останется имя функции.
			splitedFuncPath := strings.Split(fullFuncPath, ".")
			// выбираем имя функции из массива в новую переменную, для лёгкого чтения дальнейших строк.
			funcName := splitedFuncPath[len(splitedFuncPath)-1]
			// если за stepCount шагов мы дошли до роута, то заканчиваем сбор имён в цепочке функций.
			if funcName == "ServeHTTP" {
				break
			} else {
				// исключаем пустые имена функций, что периодически могут появляться если шаги ушли достаточно далеко.
				if funcName != "" {
					// сохраняем/дописываем в crudLog.FuncPath имя и строку вызова функции.
					badway += funcName + "(" + strconv.Itoa(line) + ") | "
				}
			}
		}
		UppendError(badway, "zeroLentetx", bd)
		return
	}

	if nexFlag {
		uppendTextMsg(msg, myChan, bd)
		uppendTextMsg(nextmsg, myChan, bd)
		return
	}
	msg.DisableWebPagePreview = true
	*myChan <- msg
}
