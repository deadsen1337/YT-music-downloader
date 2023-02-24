package model

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UserName    string         `json:"user_name" db:"user_name"`
	UserId      int64          `json:"user_id" db:"user_id" gorm:"UNIQUE"`
	Root        bool           `json:"root" db:"root" gorm:"default false"`
	Balance     float64        `json:"balance" db:"balance" gorm:"default:0"`
	Admin       bool           `json:"admin" db:"admin" gorm:"default false"`
	BlockedUs   bool           `json:"blocked_us" db:"blocked_us" gorm:"default false"`
	OldUsername pq.StringArray `json:"old_username" db:"old_username" gorm:"type:varchar(64)[]"`
}
