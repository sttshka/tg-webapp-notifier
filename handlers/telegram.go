package handlers

import (
	"fmt"
	"time"

	tb "gopkg.in/telebot.v3"
	"tg-notifier/models"
)

func CheckTimers(app *models.App) {
	for {
		rows, err := app.DB.Query(`
            SELECT id, user_id, days, created_at 
            FROM timers 
            WHERE datetime(created_at, '+' || days || ' days') <= datetime('now')
        `)
		if err != nil {
			time.Sleep(1 * time.Hour)
			continue
		}

		for rows.Next() {
			var t models.Timer
			if err := rows.Scan(&t.ID, &t.UserID, &t.Days, &t.CreatedAt); err != nil {
				continue
			}

			app.Bot.Send(&tb.User{ID: t.UserID},
				fmt.Sprintf("🔔 Уведомляю! %d дней прошло!", t.Days))

			app.DB.Exec("DELETE FROM timers WHERE id = ?", t.ID)
		}

		rows.Close()
		time.Sleep(1 * time.Hour)
	}
}
