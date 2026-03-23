package main

import (
 "fmt"
 "net/http"
 "strings"
 "sync"
)

// --- МОДЕЛИ ---
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

// --- ДИЗАЙН ---
func getStyle(color string) string {
 if color == "" { color = "#0088cc" }
 return fmt.Sprintf(`
<style>
    :root { --main: %s; --bg: #0e1621; --side: #17212b; --text: #fff; }
    body { background: var(--bg); color: var(--text); font-family: sans-serif; margin: 0; display: flex; height: 100vh; }
    .sidebar { width: 300px; background: var(--side); border-right: 1px solid #000; display: flex; flex-direction: column; }
    .chat-item { padding: 15px; border-bottom: 1px solid #000; text-decoration: none; color: #fff; display: flex; align-items: center; }
    .chat-item:hover { background: #242f3d; }
    .avatar { width: 40px; height: 40px; background: var(--main); border-radius: 50%%; margin-right: 10px; display: flex; align-items: center; justify-content: center; font-weight: bold; }
    .main { flex: 1; display: flex; flex-direction: column; }
    .top { padding: 15px; background: var(--side); font-weight: bold; box-shadow: 0 2px 5px #000; }
    .msgs { flex: 1; padding: 20px; overflow-y: auto; display: flex; flex-direction: column; background: #070d14; }
    .bubble { max-width: 70%%; padding: 10px; border-radius: 10px; margin-bottom: 10px; background: #182533; }
    .bubble.me { align-self: flex-end; background: var(--main); }
    .input-area { padding: 15px; background: var(--side); display: flex; gap: 10px; }
    input { flex: 1; padding: 10px; border-radius: 5px; border: none; background: #242f3d; color: #fff; }
    button { background: var(--main); color: #fff; border: none; padding: 10px 15px; border-radius: 5px; cursor: pointer; }
    .profile-box { padding: 20px; background: #242f3d; border-radius: 10px; margin: 20px; }
</style>`, color)
}

func main() {
 // Инициализация общего канала
 entities["общий"] = &Entity{Type: "канал", Name: "общий", Members: make(map[string]bool)}

 http.HandleFunc("/login_page", func(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "<html><head>%s</head><body style='justify-content:center;align-items:center;'><div class='profile-box'><h2>ZIG GARDEN</h2><form action='/register' method='POST'><input name='userid' placeholder='@username' required><br><br><input name='email' type='email' placeholder='Email' required><br><br><button style='width:100%%'>ВОЙТИ</button></form></div></body></html>", getStyle(""))
 })

 http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
  id := strings.ToLower(r.FormValue("userid"))
  email := r.FormValue("email")
  mu.Lock()
  u := &User{ID: id, Email: email, FirstName: "Новый", LastName: "Росток", Theme: "#0088cc"}
  users[email] = u
  byID[id] = u
  entities["общий"].Members[id] = true
  mu.Unlock()
  http.SetCookie(w, &http.Cookie{Name: "user_session", Value: email, Path: "/"})
  http.Redirect(w, r, "/", http.StatusSeeOther)
 })

 http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
  cookie, err := r.Cookie("user_session")
  if err != nil { http.Redirect(w, r, "/login_page", http.StatusSeeOther); return }
  user := users[cookie.Value]

  fmt.Fprintf(w, "<html><head>%s</head><body>", getStyle(user.Theme))
  fmt.Fprint(w, "<div class='sidebar'><div style='padding:15px'><b>ZIG</b></div><div class='chat-list'>")
  mu.Lock()
  for name := range entities {
   if entities[name].Members[user.ID] {
    fmt.Fprintf(w, "<a href='/?chat=%s' class='chat-item'><div class='avatar'>%s</div>%s</a>", name, strings.ToUpper(name[:1]), name)
   }
  }
  mu.Unlock()
  fmt.Fprint(w, "</div><div style='margin-top:auto;padding:15px;'><a href='/profile' style='color:#fff'>⚙ Настройки</a></div></div>")

  chatName := r.URL.Query().Get("chat")
  fmt.Fprint(w, "<div class='main'>")
  if chatName != "" && entities[chatName] != nil {
   e := entities[chatName]
   fmt.Fprintf(w, "<div class='top'>#%s</div><div class='msgs'>", chatName)
   for _, m := range e.Messages {
    cls := ""
    if m.Sender == user.ID { cls = "me" }
    fmt.Fprintf(w, "<div class='bubble %s'><b>@%s</b><br>%s</div>", cls, m.Sender, m.Text)
   }
   fmt.Fprintf(w, "</div><form class='input-area' action='/send' method='POST'><input type='hidden' name='chan' value='%s'><input name='text' placeholder='Сообщение...'><button>ОТПРАВИТЬ</button></form>", chatName)
  } else {
   fmt.Fprint(w, "<div style='display:flex;height:100%%;align-items:center;justify-content:center;'>Выберите чат</div>")
  }
  fmt.Fprint(w, "</div></body></html>")
 })

 http.HandleFunc("/send", func(w http.ResponseWriter, r *http.Request) {
  chanName, text := r.FormValue("chan"), r.FormValue("text")
  cookie, _ := r.Cookie("user_session")
  user := users[cookie.Value]
  mu.Lock()
  entities[chanName].Messages = append(entities[chanName].Messages, Msg{Sender: user.ID, Text: text})
  mu.Unlock()
  http.Redirect(w, r, "/?chat="+chanName, http.StatusSeeOther)
 })

 http.HandleFunc("/profile", func(w http.ResponseWriter, r *http.Request) {
  cookie, _ := r.Cookie("user_session")
  user := users[cookie.Value]
  fmt.Fprintf(w, "<html><head>%s</head><body style='justify-content:center;align-items:center;'><div class='profile-box'><h3>Настройки профиля</h3><form action='/update' method='POST'>Имя: <input name='f' value='%s'><br><br>Цвет темы: <input type='color' name='t' value='%s'><br><br><button>СОХРАНИТЬ</button></form><br><a href='/' style='color:#fff'>Назад</a></div></body></html>", getStyle(user.Theme), user.FirstName, user.Theme)
 })

 http.HandleFunc("/update", func(w http.ResponseWriter, r *http.Request) {
  cookie, _ := r.Cookie("user_session")
  user := users[cookie.Value]
  mu.Lock()
  user.FirstName = r.FormValue("f")
  user.Theme = r.FormValue("t")
  mu.Unlock()
  http.Redirect(w, r, "/", http.StatusSeeOther)
 })

 http.ListenAndServe(":8080", nil)
}