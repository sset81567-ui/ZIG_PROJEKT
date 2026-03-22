package main

import (
 "fmt" 
 "net/http"
)

func main() {
 // Сообщение в твоем терминале (на сервере)
 fmt.Println("🌿 [ZIG] Сад начинает цвести...")
 fmt.Println("🚀 Сервер ZIG запущен на http://localhost:8080")

 // Главная страница нашего мессенджера (то, что увидит пользователь)
 http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "text/html; charset=utf-8")
  fmt.Fprintf(w, `
   <body style="background: #0d1117; color: #769d6d; font-family: sans-serif; display: flex; align-items: center; justify-content: center; height: 100vh; margin: 0;">
    <div style="text-align: center; border: 2px solid #2d333b; padding: 40px; border-radius: 20px;">
     <h1 style="font-size: 3em; margin-bottom: 10px;">ZIG</h1>
     <p style="color: #c9d1d9;">Zoom In Garden — ваш цифровой сад общения.</p>
     <div style="background: #238636; color: white; padding: 10px 20px; border-radius: 5px; display: inline-block; margin-top: 20px;">
      Сервер работает стабильно
     </div>
    </div>
   </body>
  `)
 })

 // Запуск
 err := http.ListenAndServe(":8080", nil)
 if err != nil {
  fmt.Println("❌ Ошибка запуска сада:", err)
 }
}
package main

import (
 "fmt"
 "net/http"
 "sync"
)

// Хранилище для наших сообщений
var (
 messages []string
 mu       sync.Mutex
)

func main() {
 // Обработка главной страницы
 http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
  // Если кто-то отправил сообщение (нажал кнопку)
  if r.Method == http.MethodPost {
   msg := r.FormValue("message")
   if msg != "" {
    mu.Lock()
    messages = append(messages, msg) // Добавляем в список
    mu.Unlock()
   }
   // Перенаправляем обратно на главную, чтобы страница обновилась
   http.Redirect(w, r, "/", http.StatusSeeOther)
   return
  }

  // Формируем HTML прямо здесь
  w.Header().Set("Content-Type", "text/html; charset=utf-8")
  fmt.Fprint(w, `
   <style>
    body { font-family: sans-serif; background: #121212; color: white; text-align: center; padding: 50px; }
    .chat-box { background: #1e1e1e; border: 1px solid #333; border-radius: 10px; padding: 20px; max-width: 400px; margin: 0 auto; }
    .messages { text-align: left; height: 200px; overflow-y: auto; border-bottom: 1px solid #333; margin-bottom: 20px; padding: 10px; }
    input { padding: 10px; border-radius: 5px; border: none; width: 70%; }
    button { padding: 10px; border-radius: 5px; border: none; background: #28a745; color: white; cursor: pointer; }
   </style>
   <div class="chat-box">
    <h2>ZIG Chat</h2>
    <div class="messages">`)
  
  // Выводим все сообщения из списка
  mu.Lock()
  for _, m := range messages {
   fmt.Fprintf(w, "
