package main

import (
 "fmt"
 "net/http"
 "strings"
)

func main() {
 // 1. Страница входа
 http.HandleFunc("/login_page", func(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "<html><head>%s</head><body style='display:flex;justify-content:center;align-items:center;height:100vh;background:#0e1621;color:white;font-family:sans-serif;'><div style='text-align:center;'><h1>ZIG</h1><form action='/register' method='POST'><input name='userid' placeholder='@username' required style='padding:10px;border-radius:5px;border:none;'><br><br><input name='email' type='email' placeholder='Email' required style='padding:10px;border-radius:5px;border:none;'><br><br><button style='padding:10px 20px;background:#0088cc;color:white;border:none;border-radius:5px;cursor:pointer;'>ВОЙТИ</button></form></div></body></html>", ui("", false))
 })

 // 2. Регистрация пользователя
 http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
  id, email := strings.ToLower(r.FormValue("userid")), r.FormValue("email")
  mu.Lock()
  u := &User{ID: id, Email: email, Theme: "#0088cc"}
  users[email], byID[id] = u, u
  mu.Unlock()
  http.SetCookie(w, &http.Cookie{Name: "z_sess", Value: email, Path: "/"})
  http.Redirect(w, r, "/", http.StatusSeeOther)
 })

 // 3. Главная страница мессенджера
 http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
  c, err := r.Cookie("z_sess")
  if err != nil || users[c.Value] == nil { http.Redirect(w, r, "/login_page", http.StatusSeeOther); return }
  u := users[c.Value]
  chatID := r.URL.Query().Get("chat")
  fmt.Fprintf(w, "<html><head>%s</head><body>", ui(u.Theme, chatID != ""))
  
  // Список чатов (Sidebar)
  fmt.Fprint(w, "<div class='sidebar'><div class='top-nav'>ZIG Messenger</div><div style='flex:1;overflow-y:auto;'>")
  mu.Lock()
  for name, e := range entities {
   if e.Members[u.ID] {
    fmt.Fprintf(w, "<a href='/?chat=%s' class='chat-item'><div class='avatar'>%s</div><b>%s</b></a>", name, strings.ToUpper(name[:1]), name)
   }
  }
  mu.Unlock()
  fmt.Fprint(w, "</div><div style='padding:15px;'><form action='/search'><input name='q' placeholder='Поиск @id...' style='width:100%%'></form></div></div>")

  // Окно переписки
  fmt.Fprint(w, "<div class='main-chat'>")
  if chatID != "" && entities[chatID] != nil {
   fmt.Fprintf(w, "<div class='top-nav'><a href='/' style='color:#fff;text-decoration:none;'>&larr;</a><span>#%s</span></div><div class='messages' id='msg-box'>", chatID)
   mu.Lock()
   for _, m := range entities[chatID].Messages {
    me := ""; if m.Sender == u.ID { me = "me" }
    fmt.Fprintf(w, "<div class='bubble %s'><b>@%s</b><br>%s</div>", me, m.Sender, m.Text)
   }
   mu.Unlock()
   fmt.Fprintf(w, "</div><form class='input-bar' action='/send' method='POST'><input type='hidden' name='c' value='%s'><input name='t' placeholder='Написать...' required autocomplete='off'><button>></button></form>", chatID)
  } else {
   fmt.Fprint(w, "<div style='margin:auto;opacity:0.5;'>Выберите чат, чтобы начать общение</div>")
  }
  fmt.Fprint(w, "</div></body></html>")
 })

 // Подключение внешних функций
 http.HandleFunc("/api/messages", handleMessagesAPI)
 http.HandleFunc("/search", handleSearch)
 
 // Отправка сообщения
 http.HandleFunc("/send", func(w http.ResponseWriter, r *http.Request) {
  cn, txt := r.FormValue("c"), r.FormValue("t")
  ck, _ := r.Cookie("z_sess")
  mu.Lock()
  if e, ok := entities[cn]; ok {
   e.Messages = append(e.Messages, Msg{Sender: users[ck.Value].ID, Text: txt})
  }
  mu.Unlock()
  http.Redirect(w, r, "/?chat="+cn, http.StatusSeeOther)
 })

 // Создание диалога
 http.HandleFunc("/create_dm", func(w http.ResponseWriter, r *http.Request) {
  tid := r.URL.Query().Get("target")
  ck, _ := r.Cookie("z_sess")
  uid := users[ck.Value].ID
  cn := uid + "_" + tid
  if uid > tid { cn = tid + "_" + uid }
  mu.Lock()
  if entities[cn] == nil {
   entities[cn] = &Entity{Members: map[string]bool{uid: true, tid: true}}
  }
  mu.Unlock()
  http.Redirect(w, r, "/?chat="+cn, http.StatusSeeOther)
 })

 fmt.Println("ZIG Server is running on :8080")
 http.ListenAndServe(":8080", nil)
}