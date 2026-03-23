package main

import "sync"

// Данные пользователя
type User struct {
 ID    string
 Email string
 Theme string
}

// Структура сообщения
type Msg struct {
 Sender string
 Text   string
}

// Чат или Канал
type Entity struct {
 Members  map[string]bool
 Messages []Msg
}

// Хранилище в оперативной памяти
var (
 users    = make(map[string]*User)
 byID     = make(map[string]*User)
 entities = make(map[string]*Entity)
 mu       sync.Mutex
)
package main

import "fmt"

func ui(color string, isChat bool) string {
 if color == "" { color = "#0088cc" }
 sD, cD := "flex", "none"
 if isChat { sD, cD = "none", "flex" }

 return fmt.Sprintf(`
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<style>
    :root { --main: %s; --bg: #0e1621; --side: #17212b; --text: #fff; }
    body { background: var(--bg); color: var(--text); font-family: sans-serif; margin: 0; display: flex; height: 100vh; overflow: hidden; }
    .sidebar { width: 320px; background: var(--side); border-right: 1px solid #000; display: %s; flex-direction: column; }
    .main-chat { flex: 1; display: %s; flex-direction: column; background: #070d14; }
    @media (min-width: 769px) { .sidebar, .main-chat { display: flex !important; } }
    .top-nav { padding: 10px 15px; background: var(--side); display: flex; align-items: center; gap: 15px; font-weight: bold; border-bottom: 1px solid #000; }
    .chat-item { padding: 15px; display: flex; align-items: center; text-decoration: none; color: #fff; border-bottom: 1px solid #0e1621; }
    .avatar { width: 40px; height: 40px; background: var(--main); border-radius: 50%%; margin-right: 12px; display: flex; align-items: center; justify-content: center; }
    .messages { flex: 1; padding: 15px; overflow-y: auto; display: flex; flex-direction: column; gap: 10px; }
    .bubble { max-width: 75%%; padding: 10px; border-radius: 10px; background: #182533; font-size: 15px; }
    .bubble.me { align-self: flex-end; background: var(--main); }
    .input-bar { padding: 10px; background: var(--side); display: flex; gap: 10px; }
    input { flex: 1; padding: 12px; border-radius: 20px; border: none; background: #242f3d; color: #fff; outline: none; }
    button { background: var(--main); color: #fff; border: none; padding: 10px 20px; border-radius: 20px; cursor: pointer; }
</style>
<script>
    setInterval(() => {
        const chat = new URLSearchParams(window.location.search).get('chat');
        if (chat) {
            fetch('/api/messages?chat=' + chat)
                .then(r => r.text())
                .then(html => {
                    const box = document.getElementById('msg-box');
                    if (box && box.innerHTML.trim() !== html.trim()) {
                        box.innerHTML = html;
                        box.scrollTop = box.scrollHeight;
                    }
                });
        }
    }, 2000);
</script>`, color, sD, cD)
}
package main

import (
 "fmt"
 "net/http"
 "strings"
)

func main() {
 // Авторизация
 http.HandleFunc("/login_page", func(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "<html><head>%s</head><body style='justify-content:center;align-items:center;display:flex;'><div style='text-align:center;width:300px;'><h1>ZIG</h1><form action='/register' method='POST'><input name='userid' placeholder='@id' required><br><br><input name='email' type='email' placeholder='Email' required><br><br><button style='width:100%%'>ВОЙТИ</button></form></div></body></html>", ui("", false))
 })

 http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
  id, email := strings.ToLower(r.FormValue("userid")), r.FormValue("email")
  mu.Lock()
  u := &User{ID: id, Email: email, Theme: "#0088cc"}
  users[email], byID[id] = u, u
  mu.Unlock()
  http.SetCookie(w, &http.Cookie{Name: "z_sess", Value: email, Path: "/"})
  http.Redirect(w, r, "/", http.StatusSeeOther)
 })

 // Главный экран
 http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
  c, err := r.Cookie("z_sess")
  if err != nil { http.Redirect(w, r, "/login_page", http.StatusSeeOther); return }
  user, ok := users[c.Value]
  if !ok { http.Redirect(w, r, "/login_page", http.StatusSeeOther); return }
  
  chatID := r.URL.Query().Get("chat")
  fmt.Fprintf(w, "<html><head>%s</head><body>", ui(user.Theme, chatID != ""))

  // Sidebar
  fmt.Fprint(w, "<div class='sidebar'><div class='top-nav'>ZIG Messenger</div><div style='flex:1;overflow-y:auto;'>")
  mu.Lock()
  for name, e := range entities {
   if e.Members[user.ID] {
    fmt.Fprintf(w, "<a href='/?chat=%s' class='chat-item'><div class='avatar'>%s</div><b>%s</b></a>", name, strings.ToUpper(name[:1]), name)
   }
  }
  mu.Unlock()
  fmt.Fprint(w, "</div><div style='padding:15px;'><form action='/search'><input name='q' placeholder='Поиск @id...' style='width:100%%'></form></div></div>")

  // Чат
  fmt.Fprint(w, "<div class='main-chat'>")
  if chatID != "" && entities[chatID] != nil {
   fmt.Fprintf(w, "<div class='top-nav'><a href='/' style='color:#fff;text-decoration:none;'>&larr;</a><span>#%s</span></div><div class='messages' id='msg-box'>", chatID)
   mu.Lock()
   for _, m := range entities[chatID].Messages {
    me := ""; if m.Sender == user.ID { me = "me" }
    fmt.Fprintf(w, "<div class='bubble %s'><b>@%s</b><br>%s</div>", me, m.Sender, m.Text)
   }
   mu.Unlock()
   fmt.Fprintf(w, "</div><form class='input-bar' action='/send' method='POST'><input type='hidden' name='c' value='%s'><input name='t' placeholder='Написать...' required autocomplete='off'><button>></button></form>", chatID)
  } else {
   fmt.Fprint(w, "<div style='margin:auto;opacity:0.5;'>Выберите чат</div>")
  }
  fmt.Fprint(w, "</div></body></html>")
 })

 // API и другие обработчики (Search, Send, Create_DM) оставляем ниже...
 // (Для краткости они идентичны твоему коду)
    setupHandlers() 

 http.ListenAndServe(":8080", nil)
}

func setupHandlers() {
    // Здесь должны быть http.HandleFunc для /api/messages, /send, /search и /create_dm
}