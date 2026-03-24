package main

import (
 "sync"
 "time"
)

// --- КОНСТАНТЫ И НАСТРОЙКИ ---
const (
 MinUsernameLen = 4
 MaxUsernameLen = 24
 MinPasswordLen = 8
 DefaultTheme   = "#3498db"
)

// --- ОСНОВНЫЕ СТРУКТУРЫ ---

// User - Полный профиль пользователя со всеми твоими идеями
type User struct {
 ID           string    `json:"id"`
 FullName     string    `json:"full_name"`
 Username     string    `json:"username"`   // Тот самый 4-24 символа
 Email        string    `json:"email"`
 Password     string    `json:"-"`          // Пароль скрыт для безопасности
 Bio          string    `json:"bio"`        // "О себе"
 AvatarURL    string    `json:"avatar_url"`
 
 // Настройки интерфейса
 Language     string    `json:"language"`    // "ru" или "en"
 ThemeColor   string    `json:"theme_color"` // Например, "#000000"
 
 // Безопасность и верификация
 IsVerified   bool      `json:"is_verified"`
 VerifyCode   string    `json:"-"`           // 6-значный код для почты
 VerifyExpiry time.Time `json:"-"`           // Срок действия кода
 
 // Социальные функции
 BlockedIDs   []string  `json:"blocked_ids"` // Черный список
 LastSeen     time.Time `json:"last_seen"`   // Статус "был в сети"
 CreatedAt    time.Time `json:"created_at"`
}

// Message - Структура сообщения (подходит для лички и групп)
type Message struct {
 ID         string    `json:"id"`
 SenderID   string    `json:"sender_id"`
 ReceiverID string    `json:"receiver_id"` // Если SenderID == ReceiverID, это Избранное
 Text       string    `json:"text"`
 Type       string    `json:"type"`        // "text", "image", "voice"
 IsRead     bool      `json:"is_read"`
 CreatedAt  time.Time `json:"created_at"`
}

// Entity - Группы и Каналы
type Entity struct {
 ID          string    `json:"id"`
 Name        string    `json:"name"`
 Username    string    `json:"username"`    // Публичная ссылка @channel
 Type        string    `json:"type"`        // "group" или "channel"
 Description string    `json:"description"`
 IsPublic    bool      `json:"is_public"`   // Открытый/Закрытый
 OwnerID     string    `json:"owner_id"`
 Admins      []string  `json:"admins"`      // Кто может удалять сообщения
 Members     []string  `json:"members"`     // Список всех ID участников
 CreatedAt   time.Time `json:"created_at"`
}

// --- ХРАНИЛИЩЕ (БАЗА ДАННЫХ В ПАМЯТИ) ---

var (
 // mu - мьютекс для защиты данных при сотнях одновременных запросов
 mu sync.RWMutex

 Users    = make(map[string]*User)   // Ключ: ID (сессия)
 Entities = make(map[string]*Entity) // Все группы и каналы проекта ZIG
 Messages []Message                  // История всех переписок
)

// --- ВСПОМОГАТЕЛЬНЫЕ ФУНКЦИИ (ЛОГИКА СЕРДЦА) ---

// GetDisplayName возвращает ник через @ или полное имя
func (u *User) GetDisplayName() string {
 if u.Username != "" {
  return "@" + u.Username
 }
 return u.FullName
}

// IsAdmin проверяет, является ли пользователь админом в группе/канале
func (e *Entity) IsAdmin(userID string) bool {
 if e.OwnerID == userID {
  return true
 }
 for _, id := range e.Admins {
  if id == userID {
   return true
  }
 }
 return false
}

// AddMember добавляет человека в группу, если его там еще нет
func (e *Entity) AddMember(userID string) {
 for _, id := range e.Members {
  if id == userID {
   return
  }
 }
 e.Members = append(e.Members, userID)
}