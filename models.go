package main

import (
 "math/rand"
 "sync"
 "time"
)

type User struct {
 FullName, Username, Email, Bio, AvatarURL, Language string
 VerificationCode string
 Attempts         int
 BlockedUntil     time.Time // Бан на 10 минут
 DeletedUntil     time.Time // Бан на 24 часа после удаления
 CloudPassword    string
 PasswordHint     string
 IsPro, IsAdmin   bool
 ThemeColor       string
 MishkaCount      int
}

var (
 DataMutex  sync.RWMutex
 Users      = make(map[string]*User)
 PromoCodes = map[string]bool{"ZIG_PRO_2026": true}
)

// Твоя функция добавления Мишек
func (u *User) AddMishka() {
 DataMutex.Lock()
 u.MishkaCount++
 DataMutex.Unlock()
}

func generateCode() string {
 const charset = "0123456789"
 b := make([]byte, 6)
 for i := range b {
  b[i] = charset[rand.Intn(len(charset))]
 }
 return string(b)
}

func init() {
 rand.Seed(time.Now().UnixNano())
 // Твой аккаунт Создателя
 Users["zipsakyra5@gmail.com"] = &User{
  FullName:      "Создатель",
  Username:      "admin",
  Email:         "zipsakyra5@gmail.com",
  CloudPassword: "1D467fd67kk",
  IsAdmin:       true,
  IsPro:         true,
  ThemeColor:    "#FFD700",
  MishkaCount:   999, // Бесконечные Мишки для админа
  Bio:           "Разработка лучшего мессенджера ZIG GLOBAL",
  Language:      "ru",
 }
}