package main

import (
 "fmt"
 "net/http"
 "os"
)

func main() {
 // Маршруты
 http.HandleFunc("/", HandleHome)
 http.HandleFunc("/update", HandleUpdate)

 // Порт для Render
 port := os.Getenv("PORT")
 if port == "" { port = "8080" }

 fmt.Println("🚀 ZIG GLOBAL запущен!")
 fmt.Println("🔗 Адрес: http://localhost:" + port)
 
 if err := http.ListenAndServe(":"+port, nil); err != nil {
  fmt.Printf("Ошибка: %v\n", err)
 }
}