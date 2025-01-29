package handlers

import (
	"encoding/json"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"tg-notifier/models"
)

func WebAppHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.New("webapp").Parse(`
<!DOCTYPE html>
<html>
<head>
    <title>Timer Setup</title>
    <script src="https://telegram.org/js/telegram-web-app.js"></script>
    <script src="https://unpkg.com/htmx.org"></script>
    <style>
        body {
            font-family: -apple-system, system-ui;
            padding: 20px;
            background: var(--tg-theme-bg-color, #ffffff);
            color: var(--tg-theme-text-color, #000000);
        }
        .container {
            max-width: 400px;
            margin: 0 auto;
        }
        input, button {
            width: 100%;
            padding: 12px;
            margin: 8px 0;
            border-radius: 8px;
            font-size: 16px;
        }
        button {
            background: var(--tg-theme-button-color, #2481cc);
            color: var(--tg-theme-button-text-color, #ffffff);
            border: none;
        }
    </style>
</head>
<body>
    <div class="container">
        <h2>⏰ Set Timer</h2>
        <form 
            hx-post="/create-timer" 
            hx-target="#result"
            hx-on::after-request="if(event.detail.successful) Telegram.WebApp.close()"
        >
            <input 
                type="number" 
                name="days" 
                placeholder="Number of days" 
                min="1" 
                required
            >
            <input type="hidden" name="user_id" id="user_id">
            <button type="submit">Set Timer</button>
        </form>
        <div id="result"></div>
    </div>

    <script>
        const webApp = Telegram.WebApp;
        webApp.ready();
        
        document.getElementById('user_id').value = webApp.initDataUnsafe.user?.id || '';
        
        document.body.style.backgroundColor = webApp.colorScheme === 'dark' ? '#212121' : '#ffffff';
    </script>
</body>
</html>
`))

	tmpl.Execute(w, nil)
}

func CreateTimerHandler(app *models.App, w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	days, _ := strconv.Atoi(r.FormValue("days"))
	userID, _ := strconv.ParseInt(r.FormValue("user_id"), 10, 64)

	// Валидация данных
	if days < 1 || userID == 0 {
		http.Error(w, "Invalid data", http.StatusBadRequest)
		return
	}

	_, err := app.DB.Exec(
		"INSERT INTO timers (user_id, days, created_at) VALUES (?, ?, ?)",
		userID,
		days,
		time.Now(),
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": "Таймер создан! Возвращайтесь к своим делам!",
	})
}
