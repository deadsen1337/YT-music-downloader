package model

import "gorm.io/gorm"

type FowardMessage struct {
	gorm.Model
	BotID      int    `json:"bot_id" db:"bot_id"`
	FromChatID int64  `json:"from_chat_id" db:"from_chat_id"`
	ChatID     int64  `json:"chat_id" db:"chat_id"`
	MessageID  int    `json:"message_id" db:"message_id"`
	TextType   bool   `json:"text_type" db:"text_type" gorm:"default:false"`
	Message    string `json:"message" db:"message"`
}
