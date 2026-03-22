package main

import (
 "fmt"
 "net/http"
)

const htmlPage = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>ZIG GARDEN</title>
    <style>
        body { background-color: #1a1a1a; color: #00ff00; font-family: 'Courier New', monospace; display: flex; flex-direction: column; align-items: center; justify-content: center; height: 100vh; margin: 0; }
        .eye { font-size: 100px; margin-bottom: 20px; }
        input { padding: 10px; border: 2px solid #00ff00; background: #000; color: #00ff00; width: 250px; outline: none; border-radius: 5px; }
        button { padding: 10px 20px; background: #00ff00; color: #000; border: none; font-weight: bold; cursor: pointer; margin-top: 10px; border-radius: 5px; }
    </style>
</head>
<body>
    <div class="eye">👁</div>
    <h1>ZIG GARDEN</h1>
    <form action="/send" method="POST">
        <input type="text" name="message" placeholder="Напиши что-нибудь..." required>
        <br>
        <button type="submit">ОТПРАВИТЬ В САД</button>
    </form>
</body>
</html>
`

func main() {
 http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
  fmt.Fprint(w, htmlPage)
 })

 http.HandleFunc("/send", func(w http.ResponseWriter, r *http.Request) {
  if r.Method == http.MethodPost {
   msg := r.FormValue("message")
   fmt.Printf("[ZIG LOG] Новое сообщение: %s\n", msg) 
   fmt.Fprintf(w, "<html><body style='background:#000;color:#0f0;text-align:center;padding-top:50px;font-family:monospace;'>")
   fmt.Fprintf(w, "<h1>Сообщение доставлено: %s</h1><br><a href='/' style='color:#fff;'>Назад</a></body></html>", msg)
  }
 })

 fmt.Println("[ZIG] Сервер запущен...")
 http.ListenAndServe(":8080", nil)
}