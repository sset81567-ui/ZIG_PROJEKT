package main

import (
 "fmt"
 "net/http"
 "os"
)

func main() {
 // 1. Маршруты интерфейса
 http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "text/html; charset=utf-8")
  w.Write([]byte(GetAuthPage()))
 })
 
 http.HandleFunc("/verify-ui", func(w http.ResponseWriter, r *http.Request) {
  email := r.URL.Query().Get("email")
  w.Header().Set("Content-Type", "text/html; charset=utf-8")
  w.Write([]byte(GetVerifyPage(email)))
 })

 http.HandleFunc("/chat", func(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "text/html; charset=utf-8")
  w.Write([]byte(GetAppLayout()))
 })

 http.HandleFunc("/admin", func(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "text/html; charset=utf-8")
  w.Write([]byte(GetAdminPanel()))
 })

 // 2. Маршруты логики (handlers.go)
 http.HandleFunc("/register", HandleRegister)
 http.HandleFunc("/verify", HandleVerify)
 http.HandleFunc("/delete", HandleDelete)

 // 3. Запуск сервера Render
 port := os.Getenv("PORT")
 if port == "" {
  port = "8080"
 }

 fmt.Println("=====================================")
 fmt.Println("🚀 ZIG GLOBAL 4.0 (Calm Dark Edition)")
 fmt.Println("📍 Port:", port)
 fmt.Println("=====================================")

 err := http.ListenAndServe(":"+port, nil)
 if err != nil {
  fmt.Println("Ошибка сервера:", err)
 }
}