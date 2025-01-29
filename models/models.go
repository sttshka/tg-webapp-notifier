package models

import (
	"database/sql"
	"time"

	tb "gopkg.in/telebot.v3"
)

type Timer struct {
	ID        int
	UserID    int64
	Days      int
	CreatedAt time.Time
}

type App struct {
	DB  *sql.DB
	Bot *tb.Bot
}
