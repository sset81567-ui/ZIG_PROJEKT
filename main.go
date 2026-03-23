package main

import (
 "fmt"
 "net/http"
 "strings"
 "sync"
	"fmt"
	"net/http"
	"strings"
	"sync"
)

// Структуры (теперь с Email, как просил Render)
type User struct {
 ID, Email, FirstName, Theme string
	ID    string
	Email string
	Theme string
}

type Msg struct {
 Sender, Text string
	Sender string
	Text   string
}

type Entity struct {
 Type     string
 Members  map[string]bool
 Messages []Msg
	Members  map[string]bool
	Messages []Msg
}

var (
 users    = make(map[string]*User)
 byID     = make(map[string]*User)
 entities = make(map[string]*Entity)
 mu       sync.Mutex
	users    = make(map[string]*User)
	byID     = make(map[string]*User)
	entities = make(map[string]*Entity)
	mu       sync.Mutex
)

func ui(color string, isChatOpen bool) string {
 if color == "" { color = "#0088cc" }
 sDisp, cDisp := "flex", "none"
 if isChatOpen { sDisp, cDisp = "none", "flex" }
// Интерфейс (Добавлена поддержка мобилок и JS-обновление)
func ui(color string, isChat bool) string {
	if color == "" { color = "#0088cc" }
	sD, cD := "flex", "none"
	if isChat { sD, cD = "none", "flex" }

 return fmt.Sprintf(`
	return fmt.Sprintf(`
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<style>
   :root { --main: %s; --bg: #0e1621; --side: #17212b; --text: #fff; }
   body { background: var(--bg); color: var(--text); font-family: sans-serif; margin: 0; display: flex; height: 100vh; overflow: hidden; }
   .sidebar { width: 320px; background: var(--side); border-right: 1px solid #000; display: %s; flex-direction: column; }
   .main-chat { flex: 1; display: %s; flex-direction: column; background: #070d14; }
   @media (min-width: 769px) { .sidebar, .main-chat { display: flex !important; } }
    .chat-list { flex: 1; overflow-y: auto; }
    .chat-item { padding: 15px; display: flex; align-items: center; text-decoration: none; color: #fff; border-bottom: 1px solid #0e1621; }
    .avatar { width: 40px; height: 40px; background: var(--main); border-radius: 50%%; margin-right: 12px; display: flex; align-items: center; justify-content: center; font-weight: bold; }
   .top-nav { padding: 10px 15px; background: var(--side); display: flex; align-items: center; gap: 15px; font-weight: bold; border-bottom: 1px solid #000; }
    .messages { flex: 1; padding: 15px; overflow-y: auto; display: flex; flex-direction: column; gap: 8px; }
    .bubble { max-width: 80%%; padding: 10px; border-radius: 12px; background: #182533; font-size: 15px; }
    .chat-item { padding: 15px; display: flex; align-items: center; text-decoration: none; color: #fff; border-bottom: 1px solid #0e1621; }
    .avatar { width: 40px; height: 40px; background: var(--main); border-radius: 50%%; margin-right: 12px; display: flex; align-items: center; justify-content: center; }
    .messages { flex: 1; padding: 15px; overflow-y: auto; display: flex; flex-direction: column; gap: 10px; }
    .bubble { max-width: 75%%; padding: 10px; border-radius: 10px; background: #182533; font-size: 15px; position: relative; }
   .bubble.me { align-self: flex-end; background: var(--main); }
    .input-bar { padding: 10px; background: var(--side); display: flex; gap: 8px; }
    .input-bar { padding: 10px; background: var(--side); display: flex; gap: 10px; }
   input { flex: 1; padding: 12px; border-radius: 20px; border: none; background: #242f3d; color: #fff; outline: none; }
    button { background: var(--main); color: #fff; border: none; padding: 10px 20px; border-radius: 20px; cursor: pointer; font-weight: bold; }
    .back-btn { display: none; text-decoration: none; color: var(--main); font-size: 24px; }
    @media (max-width: 768px) { .back-btn { display: block; } }
    button { background: var(--main); color: #fff; border: none; padding: 10px 20px; border-radius: 20px; cursor: pointer; }
</style>
<script>
    setInterval(function() {
    // Скрипт для получения новых сообщений без перезагрузки страницы
    setInterval(() => {
       const chat = new URLSearchParams(window.location.search).get('chat');
       if (chat) {
            fetch('/api/messages?chat=' + chat).then(r => r.text()).then(html => {
                const div = document.getElementById('msg-container');
                if (div && div.innerHTML.trim() !== html.trim()) {
                    div.innerHTML = html; div.scrollTop = div.scrollHeight;
                }
            });
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
</script>`, color, sDisp, cDisp)
</script>`, color, sD, cD)
}

func main() {
 http.HandleFunc("/login_page", func(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "<html><head>%s</head><body style='justify-content:center;align-items:center;'><div style='padding:20px;text-align:center;'><h2>ZIG</h2><form action='/register' method='POST'><input name='userid' placeholder='@id' required><br><br><input name='email' type='email' placeholder='Email' required><br><br><button style='width:100%%'>ВОЙТИ</button></form></div></body></html>", ui("", false))
 })

 http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
  id, email := strings.ToLower(r.FormValue("userid")), r.FormValue("email")
  mu.Lock()
  u := &User{ID: id, Email: email, FirstName: "User", Theme: "#0088cc"}
  users[email], byID[id] = u, u
  mu.Unlock()
  http.SetCookie(w, &http.Cookie{Name: "z_sess", Value: email, Path: "/"})
  http.Redirect(w, r, "/", http.StatusSeeOther)
 })

 http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
c, err := r.Cookie("z_sess")
  if err != nil { http.Redirect(w, r, "/login_page", http.StatusSeeOther); return }
  user := users[c.Value]
  chatID := r.URL.Query().Get("chat")
  fmt.Fprintf(w, "<html><head>%s</head><body>", ui(user.Theme, chatID != ""))
  fmt.Fprint(w, "<div class='sidebar'><div class='top-nav'>ZIG</div><div class='chat-list'>")
  mu.Lock()
  for name, e := range entities {
   if e.Members[user.ID] { fmt.Fprintf(w, "<a href='/?chat=%s' class='chat-item'><div class='avatar'>%s</div><b>%s</b></a>", name, strings.ToUpper(name[:1]), name) }
  }
  mu.Unlock()
  fmt.Fprint(w, "</div><div style='padding:15px; border-top:1px solid #000;'><form action='/search'><input name='q' placeholder='Поиск @id'></form><br><a href='/profile' style='color:#aaa'>Настройки</a></div></div>")
  fmt.Fprint(w, "<div class='main-chat'>")
  if chatID != "" && entities[chatID] != nil {
   fmt.Fprintf(w, "<div class='top-nav'><a href='/' class='back-btn'>&larr;</a><span>#%s</span></div><div class='messages' id='msg-container'>", chatID)
   mu.Lock()
   for _, m := range entities[chatID].Messages {
    me := ""; if m.Sender == user.ID { me = "me" }
    fmt.Fprintf(w, "<div class='bubble %s'><b>@%s</b><br>%s</div>", me, m.Sender, m.Text)
   }
   mu.Unlock()
   fmt.Fprintf(w, "</div><form class='input-bar' action='/send' method='POST'><input type='hidden' name='c' value='%s'><input name='t' placeholder='Сообщение...' autofocus required><button>></button></form>", chatID)
  } else { fmt.Fprint(w, "<div style='margin:auto;color:#555;'>Выберите чат</div>") }
  fmt.Fprint(w, "</div></body></html>")
 })

 http.HandleFunc("/api/messages", func(w http.ResponseWriter, r *http.Request) {
  chat := r.URL.Query().Get("chat")
  c, _ := r.Cookie("z_sess"); user := users[c.Value]
  mu.Lock()
  if e, ok := entities[chat]; ok {
   for _, m := range e.Messages {
    me := ""; if m.Sender == user.ID { me = "me" }
    fmt.Fprintf(w, "<div class='bubble %s'><b>@%s</b><br>%s</div>", me, m.Sender, m.Text)
   }
  }
  mu.Unlock()
 })

 http.HandleFunc("/send", func(w http.ResponseWriter, r *http.Request) {
  cn, txt := r.FormValue("c"), r.FormValue("t")
  c, _ := r.Cookie("z_sess"); user := users[c.Value]
  mu.Lock()
  if entities[cn] != nil { entities[cn].Messages = append(entities[cn].Messages, Msg{Sender: user.ID, Text: txt}) }
  mu.Unlock()
  http.Redirect(w, r, "/?chat="+cn, http.StatusSeeOther)
 })

 http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
  q := strings.ToLower(r.URL.Query().Get("q"))
  c, _ := r.Cookie("z_sess"); user := users[c.Value]
  fmt.Fprintf(w, "<html><head>%s</head><body style='flex-direction:column;padding:20px;'>", ui(user.Theme, false))
  mu.Lock()
  if target, ok := byID[q]; ok { fmt.Fprintf(w, "<h3>Найден: @%s</h3><a href='/create_dm?target=%s'><button>Написать</button></a>", target.ID, target.ID)
  } else { fmt.Fprint(w, "<h3>Не найден</h3>") }
  mu.Unlock()
  fmt.Fprint(w, "<br><a href='/' style='color:#0088cc'>Назад</a></body></html>")
 })

 http.HandleFunc("/create_dm", func(w http.ResponseWriter, r *http.Request) {
  tid := r.URL.Query().Get("target")
  c, _ := r.Cookie("z_sess"); user := users[c.Value]
  cn := user.ID + "_" + tid; if user.ID > tid { cn = tid + "_" + user.ID }
  mu.Lock()
  if entities[cn] == nil { entities[cn] = &Entity{Type: "ЛС", Members: map[string]bool{user.ID: true, tid: true}} }
  mu.Unlock()
  http.Redirect(w, r, "/?chat="+cn, http.StatusSeeOther)
 })

 http.HandleFunc("/profile", func(w http.ResponseWriter, r *http.Request) {
  c, _ := r.Cookie("z_sess"); user := users[c.Value]
  fmt.Fprintf(w, "<html><head>%s</head><body style='justify-content:center;align-items:center;'><div style='padding:20px; background: #17212b; border-radius:10px; text-align:center;'><h3>Профиль</h3><form action='/upd' method='POST'>Цвет темы: <input type='color' name='th' value='%s'><br><br><button>СОХРАНИТЬ</button></form><br><a href='/' style='color:#aaa'>Назад</a></div></body></html>", ui(user.Theme,
false), user.Theme)
 })

 http.HandleFunc("/upd", func(w http.ResponseWriter, r *http.Request) {
  c, _ := r.Cookie("z_sess"); user := users[c.Value]
  mu.Lock()
  user.Theme = r.FormValue("th")
  mu.Unlock()
  http.Redirect(w, r, "/", http.StatusSeeOther)
 })

 http.ListenAndServe(":8080", nil)
}
	// Регистрация
	http.HandleFunc("/login_page", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "<html><head>%s</head><body style='justify-content:center;align-items:center;display:flex;'><div style='text-align:center;width:300px;'><h1>ZIG</h1><form action='/register' method='POST'><input name='userid' placeholder='Придумайте @id' required><br><br><input name='email' type='email' placeholder='Ваш Email' required><br><br><button style='width:100%%'>СОЗДАТЬ АККАУНТ</button></form></div></body></html>", ui("", false))
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

	// Главная страница (Список чатов и сами сообщения)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("z_sess")
		if err != nil { http.Redirect(w, r, "/login_page", http.StatusSeeOther); return }
		user, ok := users[c.Value]
		if !ok { http.Redirect(w, r, "/login_page", http.StatusSeeOther); return }
		
		chatID := r.URL.Query().Get("chat")
		fmt.Fprintf(w, "<html><head>%s</head><body>", ui(user.Theme, chatID != ""))

		// Левая панель
		fmt.Fprint(w, "<div class='sidebar'><div class='top-nav'>ZIG Messenger</div><div style='flex:1;overflow-y:auto;'>")
		mu.Lock()
		for name, e := range entities {
			if e.Members[user.ID] {
				fmt.Fprintf(w, "<a href='/?chat=%s' class='chat-item'><div class='avatar'>%s</div><b>%s</b></a>", name, strings.ToUpper(name[:1]), name)
			}