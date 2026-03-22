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
package main

import (
 "fmt"
 "net/http"
 "sync"
)

// Хранилище данных (пока в памяти сервера)
var (
 users = make(map[string]string) // Почта -> Пароль
 mu    sync.Mutex
)

const style = `
<style>
    body { background-color: #1a1a1a; color: #00ff00; font-family: 'Courier New', monospace; display: flex; flex-direction: column; align-items: center; justify-content: center; height: 100vh; margin: 0; }
    .box { border: 2px solid #00ff00; padding: 20px; border-radius: 10px; text-align: center; background: #000; }
    input { display: block; margin: 10px auto; padding: 10px; border: 1px solid #00ff00; background: #000; color: #00ff00; width: 200px; }
    button { padding: 10px 20px; background: #00ff00; color: #000; border: none; font-weight: bold; cursor: pointer; border-radius: 5px; }
    .eye { font-size: 50px; }
</style>`

func main() {
 // Главная страница
 http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
  cookie, err := r.Cookie("user_session")
  
  // Если куки нет — показываем форму регистрации
  if err != nil {
   fmt.Fprintf(w, "<html><head>%s</head><body><div class='box'><div class='eye'>👁</div><h1>ZIG: Регистрация</h1><form action='/register' method='POST'><input type='email' name='email' placeholder='Почта' required><input type='password' name='password' placeholder='Пароль' required><button type='submit'>СОЗДАТЬ АККАУНТ</button></form></div></body></html>", style)
   return
  }

  // Если вошел — показываем сад
  fmt.Fprintf(w, "<html><head>%s</head><body><div class='box'><h1>Добро пожаловать, %s!</h1><p>Ты в Саду ZIG.</p><a href='/logout' style='color:red;'>Выйти</a></div></body></html>", style, cookie.Value)
 })

 // Обработка регистрации
 http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
  if r.Method == http.MethodPost {
   email := r.FormValue("email")
   password := r.FormValue("password")

   mu.Lock()
   // Проверка: одна почта — один аккаунт
   if _, exists := users[email]; exists {
    mu.Unlock()
    fmt.Fprint(w, "Этот аккаунт уже существует!")
    return
   }
   users[email] = password
   mu.Unlock()

   // Ставим "Печеньку", чтобы запомнить устройство
   http.SetCookie(w, &http.Cookie{
    Name:  "user_session",
    Value: email,
    Path:  "/",
   })

   fmt.Printf("[ZIG] Новый пользователь: %s\n", email)
   http.Redirect(w, r, "/", http.StatusSeeOther)
  }
 })

 // Выход (удаление куки)
 http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
  http.SetCookie(w, &http.Cookie{
   Name:   "user_session",
   Value:  "",
   Path:   "/",
   MaxAge: -1,
  })
  http.Redirect(w, r, "/", http.StatusSeeOther)
 })

 fmt.Println("[ZIG] Сервер с аккаунтами запущен на :8080")
 http.ListenAndServe(":8080", nil)
}