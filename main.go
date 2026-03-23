package main

import (
 "fmt"
 "net/http"
 "strings"
 "sync"
)

// --- ДАННЫЕ СЕРВЕРА ---
type User struct {
 ID, FirstName, LastName, Bio, Email, Theme string
}

type Msg struct {
 Sender, Text string
}

type Entity struct {
 Type, Name string
 Members    map[string]bool
 Messages   []Msg
}

var (
 users    = make(map[string]*User)
 byID     = make(map[string]*User)
 entities = make(map[string]*Entity)
 mu       sync.Mutex
)

// --- ДИЗАЙН (TG STYLE) ---
func ui(color string) string {
 if color == "" { color = "#0088cc" }
 return fmt.Sprintf(`
<style>
    :root { --main: %s; --bg: #0e1621; --side: #17212b; --text: #fff; --active: #2b3948; }
    body { background: var(--bg); color: var(--text); font-family: 'Segoe UI', sans-serif; margin: 0; display: flex; height: 100vh; overflow: hidden; }
    .sidebar { width: 320px; background: var(--side); border-right: 1px solid #000; display: flex; flex-direction: column; }
    .search-area { padding: 15px; border-bottom: 1px solid #000; }
    .search-area input { width: 100%%; padding: 10px; border-radius: 8px; border: none; background: #242f3d; color: #fff; }
    .chat-list { flex: 1; overflow-y: auto; }
    .chat-item { padding: 15px; display: flex; align-items: center; text-decoration: none; color: #fff; border-bottom: 1px solid #0e1621; }
    .chat-item:hover { background: var(--active); }
    .avatar { width: 48px; height: 48px; background: var(--main); border-radius: 50%%; margin-right: 12px; display: flex; align-items: center; justify-content: center; font-weight: bold; font-size: 20px; }
    .main-chat { flex: 1; display: flex; flex-direction: column; background: #070d14; }
    .top-nav { padding: 15px 20px; background: var(--side); font-weight: bold; display: flex; justify-content: space-between; box-shadow: 0 2px 10px #000; }
    .messages { flex: 1; padding: 20px; overflow-y: auto; display: flex; flex-direction: column; }
    .bubble { max-width: 70%%; padding: 10px 15px; border-radius: 12px; margin-bottom: 8px; background: #182533; line-height: 1.4; }
    .bubble.me { align-self: flex-end; background: var(--main); border-bottom-right-radius: 2px; }
    .input-bar { padding: 15px; background: var(--side); display: flex; gap: 10px; }
    .input-bar input { flex: 1; padding: 12px; border-radius: 8px; border: none; background: #242f3d; color: #fff; }
    .btn { background: var(--main); color: #fff; border: none; padding: 10px 20px; border-radius: 8px; cursor: pointer; font-weight: bold; }
    .profile-card { padding: 30px; background: var(--side); border-radius: 15px; margin: auto; width: 350px; text-align: center; }
</style>`, color)
}

func main() {
 // Создаем стартовую группу
 entities["world"] = &Entity{Type: "Глобальный Сад", Name: "world", Members: make(map[string]bool)}

 http.HandleFunc("/login_page", func(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "<html><head>%s</head><body style='align-items:center;justify-content:center;'><div class='profile-card'><h1>ZIG PRO</h1><form action='/register' method='POST'><input name='userid' placeholder='@id' required style='width:100%%;padding:10px;margin-bottom:10px;background:#242f3d;border:none;color:#fff;'><input name='email' type='email' placeholder='Email' required style='width:100%%;padding:10px;margin-bottom:10px;background:#242f3d;border:none;color:#fff;'><button class='btn' style='width:100%%;'>ВОЙТИ В СЕТЬ</button></form></div></body></html>", ui(""))
 })

 http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
  id := strings.ToLower(r.FormValue("userid"))
  email := r.FormValue("email")
  mu.Lock()
  u := &User{ID: id, Email: email, FirstName: "User", Theme: "#0088cc"}
  users[email] = u
  byID[id] = u
  entities["world"].Members[id] = true
  mu.Unlock()
  http.SetCookie(w, &http.Cookie{Name: "z_sess", Value: email, Path: "/"})
  http.Redirect(w, r, "/", http.StatusSeeOther)
 })

 http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
  c, err := r.Cookie("z_sess")
  if err != nil { http.Redirect(w, r, "/login_page", http.StatusSeeOther); return }
  user := users[c.Value]
  
  fmt.Fprintf(w, "<html><head>%s</head><body>", ui(user.Theme))
  
  // Sidebar
  fmt.Fprint(w, "<div class='sidebar'><div class='search-area'><form action='/search'><input name='q' placeholder='Поиск...'></form></div><div class='chat-list'>")
  mu.Lock()
  for name, e := range entities {
   if e.Members[user.ID] {
    fmt.Fprintf(w, "<a href='/?chat=%s' class='chat-item'><div class='avatar'>%s</div><div><b>%s</b><br><small>%s</small></div></a>", name, strings.ToUpper(name[:1]), name, e.Type)
   }
  }
  mu.Unlock()
  fmt.Fprint(w, "</div><div style='padding:15px;'><a href='/profile' style='color:#aaa;text-decoration:none;'>⚙ Настройки</a></div></div>")

  // Chat Area
  chatID := r.URL.Query().Get("chat")
  fmt.Fprint(w, "<div class='main-chat'>")
  if chatID != "" && entities[chatID] != nil {
   e := entities[chatID]
   fmt.Fprintf(w, "<div class='top-nav'><span>#%s</span><a href='/' style='color:var(--main)'>Закрыть</a></div><div class='messages'>", chatID)
   for _, m := range e.Messages {
    me := ""
    if m.Sender == user.ID { me = "me" }
    fmt.Fprintf(w, "<div class='bubble %s'><b>@%s</b><br>%s</div>", me, m.Sender, m.Text)
   }
   fmt.Fprintf(w, "</div><form class='input-bar' action='/send' method='POST'><input type='hidden' name='c' value='%s'><input name='t' placeholder='Напишите сообщение...' autofocus required><button class='btn'>></button></form>", chatID)
  } else {
   fmt.Fprint(w, "<div style='margin:auto;color:#555;'>Выберите чат слева</div>")
  }
  fmt.Fprint(w, "</div></body></html>")
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
  fmt.Fprintf(w, "<html><head>%s</head><body><div class='profile-card'><h3>Результаты для '%s'</h3>", ui(user.Theme), q)
  mu.Lock()
  found := false
  for id := range byID {
   if strings.Contains(id, q) {
    fmt.Fprintf(w, "<p>👤 @%s <button class='btn' onclick='alert(\"Заглушка ЛС\")'>Чат</button></p>", id)
    found = true
   }
  }
  mu.Unlock()
  if !found { fmt.Fprint(w, "<p>Ничего не найдено</p>") }
  fmt.Fprint(w, "<br><a href='/' style='color:var(--main)'>Назад</a></div></body></html>")
 })

 http.HandleFunc("/profile", func(w http.ResponseWriter, r *http.Request) {
  c, _ := r.Cookie("z_sess")
  user := users[c.Value]
  fmt.Fprintf(w, "<html><head>%s</head><body style='align-items:center;justify-content:center;'><div class='profile-card'><h2>Настройки</h2><form action='/upd' method='POST'>Имя: <input name='f' value='%s' style='width:100%%;padding:8px;margin:10px 0;background:#242f3d;border:none;color:#fff;'><br>Цвет темы: <input type='color' name='th' value='%s' style='width:100%%;height:40px;border:none;background:none;'><br><button class='btn' style='width:100%%;'>СОХРАНИТЬ</button></form><br><a href='/' style='color:#aaa;'>Назад</a></div></body></html>", ui(user.Theme), user.FirstName, user.Theme)
 })

 http.HandleFunc("/upd", func(w http.ResponseWriter, r *http.Request) {
  c, _ := r.Cookie("z_sess")
  user := users[c.Value]
  mu.Lock()
  user.FirstName, user.Theme = r.FormValue("f"), r.FormValue("th")
  mu.Unlock()
  http.Redirect(w, r, "/", http.StatusSeeOther)
 })

 http.ListenAndServe(":8080", nil)
}