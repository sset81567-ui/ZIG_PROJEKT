package main

import (
 "fmt"
 "net/http"
 "time"
)

func HandleRegister(w http.ResponseWriter, r *http.Request) {
 if r.Method != http.MethodPost { return }
 
 email, name, username := r.FormValue("email"), r.FormValue("name"), r.FormValue("username")

 DataMutex.Lock()
 user, exists := Users[email]
 
 // Проверки на баны
 if exists && time.Now().Before(user.DeletedUntil) {
  DataMutex.Unlock()
  w.Write([]byte("<script>alert('Этот E-mail в бане на 24 часа после удаления!'); window.history.back();</script>"))
  return
 }
 if exists && time.Now().Before(user.BlockedUntil) {
  DataMutex.Unlock()
  w.Write([]byte("<script>alert('Бан на 10 минут за неверные коды!'); window.history.back();</script>"))
  return
 }

 if !exists {
  user = &User{FullName: name, Username: username, Email: email, Language: "ru", MishkaCount: 0}
  Users[email] = user
 }

 user.VerificationCode = generateCode()
 user.Attempts = 0
 fmt.Println("[ZIG SERVER] Код", user.VerificationCode, "отправлен на", email)
 DataMutex.Unlock()

 http.Redirect(w, r, "/verify-ui?email="+email, http.StatusSeeOther)
}

func HandleVerify(w http.ResponseWriter, r *http.Request) {
 if r.Method != http.MethodPost { return }
 
 email, code, password := r.FormValue("email"), r.FormValue("code"), r.FormValue("password")

 DataMutex.Lock()
 defer DataMutex.Unlock()

 user, exists := Users[email]
 if !exists {
  http.Redirect(w, r, "/", http.StatusSeeOther)
  return
 }

 // Вход Создателя (Админка)
 if email == "zipsakyra5@gmail.com" && password == user.CloudPassword {
  http.Redirect(w, r, "/admin", http.StatusSeeOther)
  return
 }

 if user.VerificationCode != code {
  user.Attempts++
  if user.Attempts >= 3 {
   user.BlockedUntil = time.Now().Add(10 * time.Minute)
   w.Write([]byte("<script>alert('3 ошибки! Бан на 10 минут.'); window.location.href='/';</script>"))
   return
  }
  w.Write([]byte(fmt.Sprintf("<script>alert('Неверный код! Попыток осталось: %d'); window.history.back();</script>", 3-user.Attempts)))
  return
 }

 http.Redirect(w, r, "/chat", http.StatusSeeOther)
}

func HandleDelete(w http.ResponseWriter, r *http.Request) {
 email := r.URL.Query().Get("email")
 DataMutex.Lock()
 if u, ok := Users[email]; ok {
  u.DeletedUntil = time.Now().Add(24 * time.Hour)
 }
 DataMutex.Unlock()
 w.Write([]byte("<script>alert('Аккаунт удален. Почта заблокирована на 24 часа.'); window.location.href='/';</script>"))
}