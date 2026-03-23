package main

import (
 "fmt"
 "net/http"
 "strings"
 "sync"
)

type Message struct {
 Sender  string
 Content string
}

type Entity struct {
 Type     string
 Name     string
 Username string
 Owner    string
 Members  map[string]bool
 Messages []Message
}

type User struct {
 ID       string
 Email    string
 Password string
 Blocked  map[string]bool
 Contacts []string
}

var (
 users    = make(map[string]*User)
 byID     = make(map[string]*User)
 entities = make(map[string]*Entity)
 mu       sync.Mutex
)

const style = `
<style>
    body { background: #050505; color: #00ff41; font-family: 'Courier New', monospace; margin: 0; padding: 20px; }
    .container { max-width: 450px; margin: 0 auto; border: 1px solid #00ff41; padding: 20px; box-shadow: 0 0 20px #00ff4133; border-radius: 10px; }
    h1, h2, h3 { text-align: center; text-transform: uppercase; letter-spacing: 2px; }
    input, select { width: 100%; padding: 12px; margin: 10px 0; background: #000; color: #00ff41; border: 1px solid #00ff41; box-sizing: border-box; }
    button { width: 100%; padding: 12px; background: #00ff41; color: #000; font-weight: bold; border: none; cursor: pointer; margin-top: 10px; }
    button:hover { background: #00cc33; }
    .msg-box { border: 1px solid #00ff41; height: 300px; overflow-y: auto; padding: 10px; background: #000; margin-bottom: 10px; }
    .msg { margin-bottom: 10px; font-size: 0.9em; border-left: 2px solid #00ff41; padding-left: 8px; }
    .error { color: #ff4444; background: #200; padding: 10px; border: 1px solid #ff4444; margin-bottom: 15px; text-align: center; }
    .nav { display: flex; justify-content: space-between; font-size: 0.8em; margin-bottom: 20px; border-bottom: 1px solid #00ff41; padding-bottom: 5px; }
    a { color: #00ff41; text-decoration: none; }
    .btn-exit { background: #440000; color: #ff4444; padding: 3px 8px; font-size: 0.7em; float: right; border: 1px solid #ff4444; }
</style>`

func main() {
 http.HandleFunc("/login_page", func(w http.ResponseWriter, r *http.Request) {
  err := r.URL.Query().Get("err")
  msg := ""
  if err == "id_exists" {
   msg = "<div class='error'>В САДУ УЖЕ ЕСТЬ ТАКОЙ РАСТОК!<br>Этот ID уже занят.</div>"
  }
  fmt.Fprintf(w, "<html><head>%s</head><body><div class='container'><h1>ZIG GARDEN</h1>%s<form action='/register' method='POST'><input name='userid' placeholder='@username' required><input name='email' type='email' placeholder='Email' required><input name='pass' type='password' placeholder='Пароль' required><button>ПОСАДИТЬ РАСТОК</button></form></div></body></html>", style, msg)
 })

 http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
  id := strings.ToLower(r.FormValue("userid"))
  email, pass := r.FormValue("email"), r.FormValue("pass")
  mu.Lock()
  if _, exists := byID[id]; exists {
   mu.Unlock()
   http.Redirect(w, r, "/login_page?err=id_exists", http.StatusSeeOther)
   return
  }
  u := &User{ID: id, Email: email, Password: pass, Blocked: make(map[string]bool)}
  users[email] = u
  byID[id] = u
  mu.Unlock()
  http.SetCookie(w, &http.Cookie{Name: "user_session", Value: email, Path: "/"})
  http.Redirect(w, r, "/", http.StatusSeeOther)
 })

 http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
  cookie, err := r.Cookie("user_session")
  if err != nil { http.Redirect(w, r, "/login_page", http.StatusSeeOther); return }
  user := users[cookie.Value]
  fmt.Fprintf(w, "<html><head>%s</head><body><div class='container'><div class='nav'><span>@%s</span><span><a href='/search_page'>ПОИСК</a> | <a href='/logout'>ВЫХОД</a></span></div>", style, user.ID)
  fmt.Fprint(w, "<h3>Твои Группы и Каналы</h3><div class='msg-box' style='height:auto; min-height:100px;'>")
  mu.Lock()
  for name, e := range entities {
   if e.Members[user.ID] {
    fmt.Fprintf(w, "<div><b>#%s</b> (%s) <a href='/chat?name=%s'>[ЗАЙТИ]</a> <a href='/leave?name=%s' class='btn-exit'>ВЫЙТИ</a></div><br>", name, e.Type, name, name)
   }
  }
  mu.Unlock()
  fmt.Fprint(w, "</div><button onclick=\"location.href='/create_page'\">+ СОЗДАТЬ</button></div></body></html>")
 })

 http.HandleFunc("/create_page", func(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "<html><head>%s</head><body><div class='container'><h3>Создать расток</h3><form action='/create' method='POST'><select name='type'><option value='channel'>Канал</option><option value='group'>Группа</option></select><input name='name' placeholder='Название' required><button>ПОДТВЕРДИТЬ</button></form><br><a href='/'>Назад</a></div></body></html>", style)
 })

 http.HandleFunc("/create", func(w http.ResponseWriter, r *http.Request) {
  t, name := r.FormValue("type"), strings.ToLower(r.FormValue("name"))
  cookie, _ := r.Cookie("user_session")
  user := users[cookie.Value]
  mu.Lock()
  entities[name] = &Entity{Type: t, Name: name, Owner: user.ID, Members: map[string]bool{user.ID: true}}
  mu.Unlock()
  http.Redirect(w, r, "/", http.StatusSeeOther)
 })

 http.HandleFunc("/search_page", func(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "<html><head>%s</head><body><div class='container'><h3>Поиск</h3><form action='/search' method='GET'><input name='q' placeholder='@id или название' required><button>НАЙТИ</button></form><br><a href='/'>Назад</a></div></body></html>", style)
 })

 http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
  q := strings.ToLower(r.URL.Query().Get("q"))
  fmt.Fprintf(w, "<html><head>%s</head><body><div class='container'><h3>Результаты для: %s</h3>", style, q)
  mu.Lock()
  if u, ok := byID[q]; ok {
   fmt.Fprintf(w, "<div class='msg'>👤 Юзер: @%s <form action='/block' method='POST' style='display:inline;'><input type='hidden' name='target' value='%s'><button style='width:auto; padding:5px; background:red;'>БЛОК</button></form></div>", u.ID, u.ID)
  }
  if e, ok := entities[q]; ok {
   fmt.Fprintf(w, "<div class='msg'>🌐 %s: %s <br><a href='/join?name=%s'><button>ВСТУПИТЬ</button></a></div>", e.Type, e.Name, q)
  }
  mu.Unlock()
  fmt.Fprint(w, "<br><a href='/search_page'>К поиску</a></div></body></html>")
 })

 http.HandleFunc("/join", func(w http.ResponseWriter, r *http.Request) {
  name := r.URL.Query().Get("name")
  cookie, _ := r.Cookie("user_session")
  user := users[cookie.Value]
  mu.Lock()
  if e, ok := entities[name]; ok { e.Members[user.ID] = true }
  mu.Unlock()
  http.Redirect(w, r, "/chat?name="+name, http.StatusSeeOther)
 })

 http.HandleFunc("/block", func(w http.ResponseWriter, r *http.Request) {
  target := r.FormValue("target")
  cookie, _ := r.Cookie("user_session")
  user := users[cookie.Value]
  mu.Lock()
  user.Blocked[target] = true
  mu.Unlock()
  http.Redirect(w, r, "/", http.StatusSeeOther)
 })

 http.HandleFunc("/chat", func(w http.ResponseWriter, r *http.Request) {
  name := r.URL.Query().Get("name")
  fmt.Fprintf(w, "<html><head>%s</head><body><div class='container'><h2>#%s</h2><div class='msg-box'>", style, name)
  mu.Lock()
  ch := entities[name]
  for _, m := range ch.Messages {
   fmt.Fprintf(w, "<div class='msg'><b>@%s</b>: %s</div>", m.Sender, m.Content)
  }
  mu.Unlock()
  fmt.Fprintf(w, "</div><form action='/send' method='POST'><input type='hidden' name='chan' value='%s'><input name='text' placeholder='Сообщение...' required autofocus><button>ОТПРАВИТЬ</button></form><br><a href='/'>Назад</a></div></body></html>", name)
 })

 http.HandleFunc("/send", func(w http.ResponseWriter, r *http.Request) {
  chanName, text := r.FormValue("chan"), r.FormValue("text")
  cookie, _ := r.Cookie("user_session")
  user := users[cookie.Value]
  mu.Lock()
  entities[chanName].Messages = append(entities[chanName].Messages, Message{Sender: user.ID, Content: text})
  mu.Unlock()
  http.Redirect(w, r, "/chat?name="+chanName, http.StatusSeeOther)
 })

 http.HandleFunc("/leave", func(w http.ResponseWriter, r *http.Request) {
  name := r.URL.Query().Get("name")
  cookie, _ := r.Cookie("user_session")
  user := users[cookie.Value]
  mu.Lock()
  if e, ok := entities[name]; ok { delete(e.Members, user.ID) }
  mu.Unlock()
  http.Redirect(w, r, "/", http.StatusSeeOther)
 })

 http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
  http.SetCookie(w, &http.Cookie{Name: "user_session", Value: "", Path: "/", MaxAge: -1})
  http.Redirect(w, r, "/login_page", http.StatusSeeOther)
 })

 http.ListenAndServe(":8080", nil)
}