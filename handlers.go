package main

import (
 "net/http"
 "regexp"
)

// Валидация ника: только латиница, цифры и _, 4-24 символа
func isValidUsername(u string) bool {
 match, _ := regexp.MatchString("^[a-zA-Z0-9_]{4,24}$", u)
 return match
}

// Главный обработчик домашней страницы
func HandleHome(w http.ResponseWriter, r *http.Request) {
 if r.URL.Path != "/" {
  http.NotFound(w, r)
  return
 }

 // Временный юзер для теста интерфейса
 u := &User{
  FullName:    "ZIG Developer",
  Username:    "sset81567",
  Bio:         "Разработка лучшего мессенджера ZIG GLOBAL",
  ThemeColor:  "#007AFF", // Стильный голубой
  MishkaCount: 5,         // Уже есть подарки!
 }

 w.Header().Set("Content-Type", "text/html; charset=utf-8")
 w.Write([]byte(GetLayout("Главная", u)))
}

// Обработка формы обновления
func HandleUpdate(w http.ResponseWriter, r *http.Request) {
 if r.Method != http.MethodPost {
  http.Redirect(w, r, "/", http.StatusSeeOther)
  return
 }

 // Тут будет логика сохранения в Users
 http.Redirect(w, r, "/", http.StatusSeeOther)
}