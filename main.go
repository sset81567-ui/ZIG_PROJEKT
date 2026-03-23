package main

import (
 "fmt"
 "net/http"
 "strings"
 "sync"
)

type User struct {
 ID, FirstName, Theme string
}

type Msg struct {
 Sender, Text string
}

type Entity struct {
 Type     string
 Members  map[string]bool
 Messages []Msg
}

var (
 users    = make(map[string]*User)
 byID     = make(map[string]*User)
 entities = make(map[string]*Entity)
 mu       sync.Mutex
)

func ui(color string) string {
 if color == "" { color = "#0088cc" }
 return fmt.Sprintf(`
<style>
    :root { --main: %s; --bg: #0e1621; --side: #17212b; --text: #fff; }
    body { background: var(--bg); color: var(--text); font-family: sans-serif; margin: 0; display: flex; height: 100vh; overflow: hidden; }
    .sidebar { width: 320px; background: var(--side); border-right: 1px solid #000; display: flex; flex-direction: column; }
    .chat-list { flex: 1; overflow-y: auto; }
    .chat-item { padding: 15px; display: flex; align-items: center; text-decoration: none; color: #fff; border-bottom: 1px solid #0e1621; }
    .chat-item:hover { background: #2b3948; }
    .avatar { width: 45px; height: 45px; background: var(--main); border-radius: 50%%; margin-right: 12px; display: flex; align-items: center; justify-content: center; font-weight: bold; }
    .main-chat { flex: 1; display: flex; flex-direction: column; }
    .messages { flex: 1; padding: 20px; overflow-y: auto; display: flex; flex-direction: column; background: #070d14; }
    .bubble { max-width: 70%%; padding: 10px; border-radius: 12px; margin-bottom: 8px; background: #182533; position: relative; }
    .bubble.me { align-self: flex-end; background: var(--main); }
    .input-bar { padding: 15px; background: var(--side); display: flex; gap: 10px; }
    input { flex: 1; padding: 12px; border-radius: 8px; border: none; background: #242f3d; color: #fff; }
    button { background: var(--main); color: #fff; border: none; padding: 10px 20px; border-radius: 8px; cursor: pointer; }
    .top-nav { padding: 15px; background: var(--side); display: flex; justify-content: space-between; font-weight: bold; }
</style>
<script>
    // Скрипт для авто-обновления только зоны сообщений
    setInterval(function() {
        const urlParams = new URLSearchParams(window.location.search);
        const chat = urlParams.get('chat');
        if (chat) {
            fetch('/api/messages?chat=' + chat)
                .then(response => response.text())
                .then(html => {
                    const msgDiv = document.getElementById('msg-container');
                    if (msgDiv.innerHTML !== html) {
                        msgDiv.innerHTML = html;
                        msgDiv.scrollTop = msgDiv.scrollHeight;
                    }
                });
        }
    }, 2000); // Проверка каждые 2 секунды
</script>`, color)
}

func main() {
 http.HandleFunc("/login_page", func(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "<html><head>%s</head><body style='display:flex;justify-content:center;align-items:center;'><div style='background:#17212b;padding:30px;border-radius:15px;text-align:center;'><h1>ZIG</h1><form action='/register' method='POST'><input name='userid' placeholder='@username' required style='width:100%%;padding:10px;margin-bottom:10px;background:#242f3d;border:none;color:#fff;'><input name='email' type='email' placeholder='Email' required style='width:100%%;padding:10px;margin-bottom:10px;background:#242f3d;border:none;color:#fff;'><button style='width:100%%;padding:10px;background:#0088cc;color:#fff;border:none;border-radius:5px;'>ВОЙТИ</button></form></div></body></html>", ui(""))
 })

 http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
  id := strings.ToLower(r.FormValue("userid"))
  email := r.FormValue("email")
  mu.Lock()
  u := &User{ID: id, Email: email, FirstName: "User", Theme: "#0088cc"}
  users[email] = u
  byID[id] = u
  mu.Unlock()
  http.SetCookie(w, &http.Cookie{Name: "z_sess", Value: email, Path: "/"})
  http.Redirect(w, r, "/", http.StatusSeeOther)
 })
 http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
  c, err := r.Cookie("z_sess")
  if err != nil { http.Redirect(w, r, "/login_page", http.StatusSeeOther); return }
  user := users[c.Value]
  
  fmt.Fprintf(w, "<html><head>%s</head><body>", ui(user.Theme))
  fmt.Fprint(w, "<div class='sidebar'><div style='padding:15px; border-bottom: 1px solid #000;'><form action='/search'><input name='q' placeholder='Поиск @id...' style='width:100%%;padding:8px;background:#242f3d;border:none;color:#fff;border-radius:5px;'></form></div><div class='chat-list'>")
  mu.Lock()
  for name, e := range entities {
   if e.Members[user.ID] {
    fmt.Fprintf(w, "<a href='/?chat=%s' class='chat-item'><div class='avatar'>%s</div><b>%s</b></a>", name, strings.ToUpper(name[:1]), name)
   }
  }
  mu.Unlock()
  fmt.Fprint(w, "</div><div style='padding:15px;'><a href='/profile' style='color:#aaa;'>Настройки</a></div></div>")

  chatID := r.URL.Query().Get("chat")
  fmt.Fprint(w, "<div class='main-chat'>")
  if chatID != "" && entities[chatID] != nil {
   fmt.Fprintf(w, "<div class='top-nav'><span>#%s</span><a href='/' style='color:#0088cc;'>Закрыть</a></div>", chatID)
   fmt.Fprint(w, "<div class='messages' id='msg-container'>")
   // Содержимое подгрузится само или при первой загрузке
   for _, m := range entities[chatID].Messages {
    me := ""
    if m.Sender == user.ID { me = "me" }
    fmt.Fprintf(w, "<div class='bubble %s'><b>@%s</b><br>%s</div>", me, m.Sender, m.Text)
   }
   fmt.Fprintf(w, "</div><form class='input-bar' action='/send' method='POST'><input type='hidden' name='c' value='%s'><input name='t' placeholder='Напишите сообщение...' autofocus required><button>></button></form>", chatID)
  } else {
   fmt.Fprint(w, "<div style='margin:auto;color:#555;'>Выберите чат</div>")
  }
  fmt.Fprint(w, "</div></body></html>")
 })

 // API для получения сообщений без перезагрузки
 http.HandleFunc("/api/messages", func(w http.ResponseWriter, r *http.Request) {
  chat := r.URL.Query().Get("chat")
  c, _ := r.Cookie("z_sess")
  user := users[c.Value]
  mu.Lock()
  if e, ok := entities[chat]; ok {
   for _, m := range e.Messages {
    me := ""
    if m.Sender == user.ID { me = "me" }
    fmt.Fprintf(w, "<div class='bubble %s'><b>@%s</b><br>%s</div>", me, m.Sender, m.Text)
   }
  }
  mu.Unlock()
 })

 http.HandleFunc("/send", func(w http.ResponseWriter, r *http.Request) {
  cn, txt := r.FormValue("c"), r.FormValue("t")
  c, _ := r.Cookie("z_sess")
  user := users[c.Value]
  mu.Lock()
  entities[cn].Messages = append(entities[cn].Messages, Msg{Sender: user.ID, Text: txt})
  mu.Unlock()
  http.Redirect(w, r, "/?chat="+cn, http.StatusSeeOther)
 })

 http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
  q := strings.ToLower(r.URL.Query().Get("q"))
  c, _ := r.Cookie("z_sess")
  user := users[c.Value]
  fmt.Fprintf(w, "<html><head>%s</head><body><div style='background:#17212b;padding:30px;border-radius:15px;margin:auto;'><h3>Результаты для '%s'</h3>", ui(user.Theme), q)
  mu.Lock()
  if target, ok := byID[q]; ok {
   chatName := user.ID + "_" + target.ID
   if user.ID > target.ID { chatName = target.ID + "_" + user.ID }
   fmt.Fprintf(w, "<p>👤 @%s <a href='/create_dm?target=%s'><button>Написать в ЛС</button></a></p>", target.ID, target.ID)
  } else {
   fmt.Fprint(w, "<p>Никто не найден</p>")
  }
  mu.Unlock()
  fmt.Fprint(w, "<br><a href='/' style='color:#0088cc;'>Назад</a></div></body></html>")
 })

 http.HandleFunc("/create_dm", func(w http.ResponseWriter, r *http.Request) {
  targetID := r.URL.Query().Get("target")
  c, _ := r.Cookie("z_sess")
  user := users[c.Value]
  
  chatName := user.ID + "_" + targetID
  if user.ID > targetID { chatName = targetID + "_" + user.ID }
  
  mu.Lock()
  if entities[chatName] == nil {
   entities[chatName] = &Entity{Type: "ЛС", Members: map[string]bool{user.ID: true, targetID: true}}
  }
  mu.Unlock()
  http.Redirect(w, r, "/?chat="+chatName, http.StatusSeeOther)
 })
http.ListenAndServe(":8080", nil)
}