package main

import (
 "fmt"
 "net/http"
 "strings"
 "sync"
)

// --- СТРУКТУРЫ ДАННЫХ ---

type Message struct {
 Sender  string
 Content string
}

type Entity struct {
 Type     string // "channel" или "group"
 Name     string
 Username string // Для поиска (@name)
 Owner    string
 Members  map[string]bool
 Messages []Message
}

type User struct {
 ID       string
 Email    string
 Password string
 Blocked  map[string]bool // Черный список
 Contacts []string        // Список чатов
}

var (
 users    = make(map[string]*User)   // Email -> User
 byID     = make(map[string]*User)   // ID -> User
 entities = make(map[string]*Entity) // ID/Name -> Entity
 mu       sync.Mutex
)

// --- ДИЗАЙН (CSS) ---

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
 // 1. СТРАНИЦА РЕГИСТРАЦИИ
 http.HandleFunc("/login_page", func(w http.ResponseWriter, r *http.Request) {
  err := r.URL.Query().Get("err")
  msg := ""
  if err == "id_exists" {
   msg = "<div class='error'>В САДУ УЖЕ ЕСТЬ ТАКОЙ РАСТОК!<br>Этот ID уже занят, выбери другой.</div>"
  }
  fmt.Fprintf(w, "<html><head>%s</head><body><div class='container'><h1>ZIG GARDEN</h1>%s<form action='/register' method='POST'><input name='userid' placeholder='@username (напр. neo)' required><input name='email' type='email' placeholder='Email' required><input name='pass' type='password' placeholder='Пароль' required><button>ПОСАДИТЬ РАСТОК</button></form></div></body></html>", style, msg)
 })

 // 2. ЛОГИКА РЕГИСТРАЦИИ
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

 // 3. ГЛАВНЫЙ ХАБ
 http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
  cookie, err := r.Cookie("user_session")
  if err != nil { http.Redirect(w, r, "/login_page", http.StatusSeeOther); return }
  user := users[cookie.Value]

  fmt.Fprintf(w, "<html><head>%s</head><body><div class='container'><div class='nav'><span>@%s</span><span><a href='/search_page'>ПОИСК</a> | <a href='/logout'>ВЫХОД</a></span></div>", style, user.ID)
  package main

import (
 "fmt"
 "net/http"
 "strings"
 "sync"
)

// --- СТРУКТУРЫ ДАННЫХ ---

type Message struct {
 Sender  string
 Content string
}

type Entity struct {
 Type     string // "channel" или "group"
 Name     string
 Username string // Для поиска (@name)
 Owner    string
 Members  map[string]bool
 Messages []Message
}

type User struct {
 ID       string
 Email    string
 Password string
 Blocked  map[string]bool // Черный список
 Contacts []string        // Список чатов
}

var (
 users    = make(map[string]*User)   // Email -> User
 byID     = make(map[string]*User)   // ID -> User
 entities = make(map[string]*Entity) // ID/Name -> Entity
 mu       sync.Mutex
)

// --- ДИЗАЙН (CSS) ---

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
 // 1. СТРАНИЦА РЕГИСТРАЦИИ
 http.HandleFunc("/login_page", func(w http.ResponseWriter, r *http.Request) {
  err := r.URL.Query().Get("err")
  msg := ""
  if err == "id_exists" {
   msg = "<div class='error'>В САДУ УЖЕ ЕСТЬ ТАКОЙ РАСТОК!<br>Этот ID уже занят, выбери другой.</div>"
  }
  fmt.Fprintf(w, "<html><head>%s</head><body><div class='container'><h1>ZIG GARDEN</h1>%s<form action='/register' method='POST'><input name='userid' placeholder='@username (напр. neo)' required><input name='email' type='email' placeholder='Email' required><input name='pass' type='password' placeholder='Пароль' required><button>ПОСАДИТЬ РАСТОК</button></form></div></body></html>", style, msg)
 })

 // 2. ЛОГИКА РЕГИСТРАЦИИ
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

 // 3. ГЛАВНЫЙ ХАБ
 http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
  cookie, err := r.Cookie("user_session")
  if err != nil { http.Redirect(w, r, "/login_page", http.StatusSeeOther); return }
  user := users[cookie.Value]

  fmt.Fprintf(w, "<html><head>%s</head><body><div class='container'><div class='nav'><span>@%s</span><span><a href='/search_page'>ПОИСК</a> | <a href='/logout'>ВЫХОД</a></span></div>", style, user.ID)
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

 http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
  http.SetCookie(w, &http.Cookie{Name: "user_session", Value: "", Path: "/", MaxAge: -1})
  http.Redirect(w, r, "/login_page", http.StatusSeeOther)
 })

 fmt.Println("ZIG CORE ENGINE v3.0 STARTED ON :8080")
 http.ListenAndServe(":8080", nil)
}