package main

import (
 "fmt"
 "net/http"
 "strings"
 "sync"
)

type User struct {
 ID        string
 FirstName string
 LastName  string
 Bio       string
 Email     string
 Theme     string
}

type Entity struct {
 Type     string
 Name     string
 Members  map[string]bool
 Messages []Msg
}

type Msg struct {
 Sender string
 Text   string
}

var (
 users    = make(map[string]*User)
 byID     = make(map[string]*User)
 entities = make(map[string]*Entity)
 mu       sync.Mutex
)

// Функция для генерации CSS с учетом выбранной темы
func getStyle(color string) string {
 if color == "" { color = "#0088cc" } // Цвет Telegram по умолчанию
 return fmt.Sprintf(`
<style>
    :root { --main-color: %s; --bg: #17212b; --sidebar: #0e1621; --text: #ffffff; }
    body { background: var(--bg); color: var(--text); font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif; margin: 0; display: flex; height: 100vh; overflow: hidden; }
    
    /* Сайдбар */
    .sidebar { width: 300px; background: var(--sidebar); border-right: 1px solid #000; display: flex; flex-direction: column; }
    .search-box { padding: 15px; border-bottom: 1px solid #111; }
    .search-box input { width: 100%%; padding: 8px; background: #242f3d; border: none; color: #fff; border-radius: 5px; }
    .chat-list { flex: 1; overflow-y: auto; }
    .chat-item { padding: 15px; cursor: pointer; border-bottom: 1px solid #111; display: flex; align-items: center; text-decoration: none; color: #fff; }
    .chat-item:hover { background: #2b3948; }
    .avatar { width: 45px; height: 45px; background: var(--main-color); border-radius: 50%%; display: flex; align-items: center; justify-content: center; margin-right: 12px; font-weight: bold; }

    /* Основная зона */
    .main-area { flex: 1; display: flex; flex-direction: column; background: url('https://user-images.githubusercontent.com/15075759/28719144-86dc0f70-73b1-11e7-911d-60d70fcded21.png'); }
    .top-bar { padding: 10px 20px; background: var(--sidebar); display: flex; justify-content: space-between; align-items: center; box-shadow: 0 2px 5px rgba(0,0,0,0.2); }
    .messages { flex: 1; padding: 20px; overflow-y: auto; display: flex; flex-direction: column; }
    .bubble { max-width: 70%%; padding: 10px 15px; border-radius: 15px; margin-bottom: 10px; background: #182533; position: relative; }
    .bubble.me { align-self: flex-end; background: var(--main-color); }
    .input-area { padding: 20px; background: var(--sidebar); display: flex; gap: 10px; }
    .input-area input { flex: 1; padding: 12px; border-radius: 8px; border: none; background: #242f3d; color: #fff; }
    
    button { background: var(--main-color); color: #fff; border: none; padding: 10px 20px; border-radius: 8px; cursor: pointer; font-weight: bold; }
    .settings-panel { padding: 20px; background: #242f3d; border-radius: 10px; margin: 20px; }
</style>`, color)
}

func main() {
 // Регистрация (упрощенно для примера)
 http.HandleFunc("/login_page", func(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "<html><head>%s</head><body style='justify-content:center; align-items:center;'><div class='settings-panel' style='width:300px;'><h2>ZIG Welcome</h2><form action='/register' method='POST'><input name='userid' placeholder='@username' required style='width:100%%; margin-bottom:10px;'><input name='email' type='email' placeholder='Email' required style='width:100%%; margin-bottom:10px;'><button style='width:100%%;'>ВОЙТИ В САД</button></form></div></body></html>", getStyle(""))
 })

 http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
  id := strings.ToLower(r.FormValue("userid"))
  email := r.FormValue("email")
  mu.Lock()
  u := &User{ID: id, Email: email, FirstName: "Новый", LastName: "Росток", Theme: "#0088cc"}
  users[email] = u
  byID[id] = u
  mu.Unlock()
  http.SetCookie(w, &http.Cookie{Name: "user_session", Value: email, Path: "/"})
  http.Redirect(w, r, "/", http.StatusSeeOther)
 })

 // Главный интерфейс
 http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
  cookie, err := r.Cookie("user_session")
  if err != nil { http.Redirect(w, r, "/login_page", http.StatusSeeOther); return }
  user := users[cookie.Value]

  fmt.Fprintf(w, "<html><head>%s</head><body>", getStyle(user.Theme))
  
  // Сайдбар
  fmt.Fprint(w, "<div class='sidebar'><div class='search-box'><form action='/search' method='GET'><input name='q' placeholder='Поиск чатов или людей...'></form></div><div class='chat-list'>")
  for name, e := range entities {
   if e.Members[user.ID] {
    fmt.Fprintf(w, "<a href='/?chat=%s' class='chat-item'><div class='avatar'>%s</div><div><b>%s</b><br><small>%s</small></div></a>", name, strings.ToUpper(name[:1]), name, e.Type)
   }
  }
  fmt.Fprint(w, "</div><div style='padding:10px;'><a href='/profile'>⚙ Настройки профиля</a></div></div>")

  // Основная область чата
  chatName := r.URL.Query().Get("chat")
  fmt.Fprint(w, "<div class='main-area'>")
  if chatName != "" {
   e := entities[chatName]
   fmt.Fprintf(w, "<div class='top-bar'><b>#%s</b> <span>%d участников</span></div>", chatName, len(e.Members))
   fmt.Fprint(w, "<div class='messages'>")
   for _, m := range e.Messages {
    class := ""
    if m.Sender == user.ID { class = "me" }
    fmt.Fprintf(w, "<div class='bubble %s'><b>@%s</b><br>%s</div>", class, m.Sender, m.Text)
   }
   fmt.Fprintf(w, "</div><form class='input-area' action='/send' method='POST'><input type='hidden' name='chan' value='%s'><input name='text' placeholder='Напишите сообщение...' autofocus><button>ОТПРАВИТЬ</button></form>", chatName)
  } else {
   fmt.Fprint(w, "<div style='display:flex; height:100%%; align-items:center; justify-content:center; color:#555;'>Выберите чат, чтобы начать общение</div>")
  }
  fmt.Fprint(w, "</div></body></html>")
 })

 // Настройки профиля
 http.HandleFunc("/profile", func(w http.ResponseWriter, r *http.Request) {
  cookie, _ := r.Cookie("user_session")
  user := users[cookie.Value]
  fmt.Fprintf(w, "<html><head>%s</head><body style='justify-content:center; align-items:center;'><div class='settings-panel'><h2>Профиль</h2><form action='/update_profile' method='POST'>Имя: <input name='first' value='%s'><br>Фамилия: <input name='last' value='%s'><br>О себе: <input name='bio' value='%s'><br>Цвет темы: <input type='color' name='theme' value='%s'><br><button>СОХРАНИТЬ</button></form><br><a href='/'>← Назад в чаты</a></div></body></html>", getStyle(user.Theme), user.FirstName, user.LastName, user.Bio, user.Theme)
 })

 http.HandleFunc("/update_profile", func(w http.ResponseWriter, r *http.Request) {
  cookie, _ := r.Cookie("user_session")
  user := users[cookie.Value]
  mu.Lock()
  user.FirstName = r.FormValue("first")
  user.LastName = r.FormValue("last")
  user.Bio = r.FormValue("bio")
  user.Theme = r.FormValue("theme")
  mu.Unlock()
  http.Redirect(w, r, "/", http.StatusSeeOther)
 })

    // Поиск (исправленный)
    http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
        q := strings.ToLower(r.URL.Query().Get("q"))
        cookie, _ := r.Cookie("user_session")
        user := users[cookie.Value]
        fmt.Fprintf(w, "<html><head>%s</head><body><div class='container' style='margin:20px auto; width:400px;'><h3>Результаты поиска: '%s'</h3>", getStyle(user.Theme), q)
        mu.Lock()
        for id, u := range byID {
            if strings.Contains(id, q) {
                fmt.Fprintf(w, "<div class='chat-item'>👤 @%s (%s %s)</div>", u.ID, u.FirstName, u.LastName)
            }
        }
        for name, e := range entities {
            if strings.Contains(name, q) {
                fmt.Fprintf(w, "<div class='chat-item'>🌐 %s: %s <a href='/join?name=%s'>[ВСТУПИТЬ]</a></div>", e.Type, name, name)
            }
        }
        mu.Unlock()
        fmt.Fprint(w, "<br><a href='/'>Назад</a></div></body></html>")
    })

    // Остальные функции (send, join, leave) оставляем как в прошлый раз...
    // (Для краткости они тут подразумеваются)

 http.ListenAndServe(":8080", nil)
}