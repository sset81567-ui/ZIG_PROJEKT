package main

import (
 "math/rand"
 "sync"
 "time"
)

// Структура пользователя со всеми фишками ZIG
type User struct {
 ID               string
 FullName         string
 Username         string
 Email            string
 Bio              string
 AvatarURL        string
 Language         string // ru, uk, be, en
 ThemeColor       string
 
 // Безопасность и авторизация
 VerificationCode string
 Attempts         int
 BlockedUntil     time.Time // Бан 10 мин за неверные коды
 DeletedUntil     time.Time // Бан 24 часа после удаления
 CloudPassword    string
 PasswordHint     string
 
 // Статусы
 IsPro            bool
 IsAdmin          bool
 HideLastSeen     bool
 HideReadStatus   bool
 MishkaCount      int
}

// Структура чатов и каналов
type Chat struct {
 ID         string
 Name       string
 Username   string
 AvatarURL  string
 IsChannel  bool
 IsPrivate  bool
 InviteLink string
 OwnerID    string
 Members    []string
}

var (
 DataMutex  sync.RWMutex
 Users      = make(map[string]*User)
 Chats      = make(map[string]*Chat)
 PromoCodes = map[string]bool{"ZIG_PRO_2026": true, "CREATOR_GIFT": true}
)

// Генерация случайного 6-значного кода
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
 // Создаем твой аккаунт Бога (Создателя)
 Users["zipsakyra5@gmail.com"] = &User{
  FullName:      "Создатель",
  Username:      "admin",
  Email:         "zipsakyra5@gmail.com",
  CloudPassword: "1D467fd67kk",
  IsAdmin:       true,
  IsPro:         true,
  ThemeColor:    "#FFD700", // Золотой акцент
  Language:      "ru",
 }
}