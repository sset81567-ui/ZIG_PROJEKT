package main

import (
 "fmt"
 "net/http"
 "os"
)

func main() {
 // Инициализируем тестового юзера, чтобы не было пустых переменных
 testUser := &User{FullName: "Admin", ThemeColor: "#3498db", Language: "ru"}

 http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
  content := `<div class="empty-state"><h3>Выберите чат</h3><p>Или создайте новый канал</p></div>`
  fmt.Fprint(w, GetLayout(content, testUser))
 })

 http.HandleFunc("/register", RegisterHandler)

 port := os.Getenv("PORT")
 if port == "" { port = "8080" }

 fmt.Println("ZIG запущен! Порт:", port)
 http.ListenAndServe(":"+port, nil)
}