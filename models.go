package main

import (
 "sync"
 "time"
)

// User - Структура пользователя со всеми твоими идеями
type User struct {
 ID           string
 FullName     string
 Username     string    // Лимит 4-24 символа
 Email        string
 Password     string
 Bio          string    // О себе
 AvatarURL    string
 Language     string    // "ru" или "en"
 ThemeColor   string    // Твой любимый цвет
 IsVerified   bool      // Подтверждена ли почта
 VerifyCode   string    // 6-значный код
 BlockedUsers []string  // Черный список
 CreatedAt    time.Time
}

// Message - Сообщения (и для лички, и для Избранного)
type Message struct {
 ID         string
 SenderID   string
 ReceiverID string    // Если SenderID == ReceiverID, это Избранное
 Text       string
 IsRead     bool
 CreatedAt  time.Time
}

// Entity - Группы и Каналы
type Entity struct {
 ID          string
 Name        string
 Type        string    // "group" или "channel"
 Description string
 IsPublic    bool
 OwnerID     string
 Admins      []string
 Members     []string
}

// Глобальное хранилище данных
var (
 mu       sync.RWMutex
 Users    = make(map[string]*User)
 Entities = make(map[string]*Entity)
 Messages []Message
)

// Вспомогательные функции (Логика данных)
func (u *User) GetDisplayName() string {
 if u.Username != "" {
  return "@" + u.Username
 }
 return u.FullName
}

func (e *Entity) IsAdmin(userID string) bool {
 if e.OwnerID == userID { return true }
 for _, id := range e.Admins {
  if id == userID { return true }
 }
 return false
}

func (e *Entity) AddMember(userID string) {
 for _, id := range e.Members {
  if id == userID { return }
 }
 e.Members = append(e.Members, userID)
}