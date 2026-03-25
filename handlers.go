package main

import (
 "fmt"
 "net/http"
 "time"
)

// Обработка регистрации и проверки банов
func HandleRegister(w http.ResponseWriter, r *http.Request) {
 if r.Method != http.MethodPost {
  http.Redirect(w, r, "/", http.StatusSeeOther)
  return
 }

 email := r.FormValue("email")
 name := r.FormValue("name")
 username := r.FormValue("username")

 DataMutex.Lock()
 defer DataMutex.Unlock()

 user, exists := Users[email]
 
 // Проверка на 24-часовой бан после удаления
 if exists && time.Now().Before(user.DeletedUntil) {
  http.Error(w, "Этот E-mail заблокирован на 24 часа после удаления.", http.StatusForbidden)
  return
 }

 // Проверка на 10-минутный бан за ошибки
 if exists && time.Now().Before(user.BlockedUntil) {
  http.Error(w, "Слишком много попыток. Подождите 10 минут.", http.StatusTooManyRequests)
  return
 }

 // Если аккаунта нет - создаем заготовку
 if !exists {
  user = &User{
   FullName: name,
   Username: username,
   Email:    email,
   Language: "ru",
  }
  Users[email] = user
 }

 // Генерация и "отправка" кода
 user.VerificationCode = generateCode()
 user.Attempts = 0 // Сбрасываем попытки
 fmt.Printf("[ZIG MAIL] Код %s отправлен на %s\n", user.VerificationCode, email)

 http.Redirect(w, r, "/verify-ui?email="+email, http.StatusSeeOther)
}

// Проверка 6-значного кода
func HandleVerify(w http.ResponseWriter, r *http.Request) {
 if r.Method != http.MethodPost { return }
 
 email := r.FormValue("email")
 code := r.FormValue("code")

 DataMutex.Lock()
 defer DataMutex.Unlock()

 user, exists := Users[email]
 if !exists {
  http.Error(w, "Пользователь не найден", 404)
  return
 }

 // Если это админ - пускаем в панель
 if email == "zipsakyra5@gmail.com" && r.FormValue("password") == user.CloudPassword {
  http.Redirect(w, r, "/admin", http.StatusSeeOther)
  return
 }

 if user.VerificationCode != code {
  user.Attempts++
  if user.Attempts >= 3 {
   user.BlockedUntil = time.Now().Add(10 * time.Minute)
   http.Error(w, "3 неверных попытки! Вы заблокированы на 10 минут.", 403)
   return
  }
  w.Write([]byte(fmt.Sprintf("<script>alert('Неверный код! Осталось попыток: %d'); window.history.back();</script>", 3-user.Attempts)))
  return
 }

 // Успешный вход
 http.Redirect(w, r, "/chat", http.StatusSeeOther)
}

// Удаление аккаунта
func HandleDelete(w http.ResponseWriter, r *http.Request) {
 email := r.URL.Query().Get("email")
 DataMutex.Lock()
 if user, exists := Users[email]; exists {
  // Ставим метку удаления на 24 часа
  user.DeletedUntil = time.Now().Add(24 * time.Hour)
 }
 DataMutex.Unlock()
 http.Redirect(w, r, "/", http.StatusSeeOther)
}