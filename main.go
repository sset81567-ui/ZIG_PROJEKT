package main

import (
 "fmt"
 "net/http"
 "os"
)

func main() {
 http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
  w.Write([]byte(GetAuthPage()))
 })
 
 http.HandleFunc("/verify-ui", func(w http.ResponseWriter, r *http.Request) {
  email := r.URL.Query().Get("email")
  w.Write([]byte(GetVerifyPage(email)))
 })

 http.HandleFunc("/chat", func(w http.ResponseWriter, r *http.Request) {
  w.Write([]byte(GetAppLayout()))
 })

 http.HandleFunc("/admin", func(w http.ResponseWriter, r *http.Request) {
  DataMutex.RLock()
  count := len(Users)
  DataMutex.RUnlock()
  w.Write([]byte(GetAdminPanel(count)))
 })

 http.HandleFunc("/register", HandleRegister)
 http.HandleFunc("/verify", HandleVerify)
 http.HandleFunc("/delete", HandleDelete)

 port := os.Getenv("PORT")
 if port == "" { port = "8080" }

 fmt.Println("🚀 ZIG GLOBAL запущен на порту:", port)
 http.ListenAndServe(":"+port, nil)
}