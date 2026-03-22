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
