package main

import (
	"database/sql"
	"net/http"
	"os"
	"time"

	"tg-notifier/handlers"
	"tg-notifier/models"

	_ "github.com/mattn/go-sqlite3"
	tb "gopkg.in/telebot.v3"
)

func main() {
	db := initDB()
	defer db.Close()

	bot := initBot()
	app := &models.App{DB: db, Bot: bot}

	go handlers.CheckTimers(app)

	// Веб-роуты
	http.HandleFunc("/webapp", handlers.WebAppHandler)
	http.HandleFunc("/create-timer", func(w http.ResponseWriter, r *http.Request) {
		handlers.CreateTimerHandler(app, w, r)
	})

	// Обработчик команды /start с кнопкой Web App
	bot.Handle("/start", func(c tb.Context) error {
		menu := &tb.ReplyMarkup{}
		btn := menu.WebApp("Задать таймер", &tb.WebApp{URL: os.Getenv("CLIENT_URL")})
		menu.Inline(menu.Row(btn))
		return c.Send("Добро пожаловать! Ниже кнопка для создания таймера:", menu)
	})

	go bot.Start()
	http.ListenAndServe(":8080", nil)
}

func initDB() *sql.DB {
	db, err := sql.Open("sqlite3", "storage.db")
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS timers (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        user_id INTEGER,
        days INTEGER,
        created_at DATETIME
    )`)
	if err != nil {
		panic(err)
	}

	return db
}

func initBot() *tb.Bot {
	bot, err := tb.NewBot(tb.Settings{
		Token:  os.Getenv("TELEGRAM_TOKEN"),
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		panic(err)
	}
	return bot
}
