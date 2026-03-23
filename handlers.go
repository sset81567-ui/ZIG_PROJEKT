package main

import (
 "fmt"
 "net/http"
 "strings"
)

func handleMessagesAPI(w http.ResponseWriter, r *http.Request) {
 chat := r.URL.Query().Get("chat")
 ck, _ := r.Cookie("z_sess")
 u := users[ck.Value]
 mu.Lock()
 if e, exists := entities[chat]; exists {
  for _, m := range e.Messages {
   me := ""; if m.Sender == u.ID { me = "me" }
   fmt.Fprintf(w, "<div class='bubble %s'><b>@%s</b><br>%s</div>", me, m.Sender, m.Text)
  }
 }
 mu.Unlock()
}

func handleSearch(w http.ResponseWriter, r *http.Request) {
 q := strings.ToLower(r.URL.Query().Get("q"))
 ck, _ := r.Cookie("z_sess")
 fmt.Fprintf(w, "<html><head>%s</head><body style='display:flex;flex-direction:column;align-items:center;padding:50px;background:#0e1621;color:white;'>", ui(users[ck.Value].Theme, false))
 mu.Lock()
 if t, ok := byID[q]; ok && t.ID != users[ck.Value].ID {
  fmt.Fprintf(w, "<h1>@%s</h1><a href='/create_dm?target=%s'><button>НАПИСАТЬ</button></a>", t.ID, t.ID)
 } else { fmt.Fprint(w, "<h3>Не найден</h3>") }
 mu.Unlock()
 fmt.Fprint(w, "<br><a href='/'>Назад</a></body></html>")
}