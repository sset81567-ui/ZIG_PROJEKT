package main

import (
 "sync"
 "time"
)

// User — основной профиль со всеми фишками
type User struct {
 ID           string
 FullName     string
 Username     string
 Email        string
 Bio          string
 ThemeColor   string // Для кастомного дизайна
 Language     string
 IsPro        bool      // Для премиум-статуса
 MishkaCount  int       // Кол-во подарков «Мишка»
 LastSeen     time.Time
 CreatedAt    time.Time
}

// Entity — чаты, группы или каналы
type Entity struct {
 ID        string
 Name      string
 OwnerID   string
 Members   []string
 IsPrivate bool
 AvatarURL string
}

// Глобальное хранилище с защитой от сбоев (Mutex)
var (
 DataMutex sync.RWMutex
 Users     = make(map[string]*User)
 Entities  = make(map[string]*Entity)
)

// GetDisplayName возвращает ник или имя
func (u *User) GetDisplayName() string {
 if u.Username != "" {
  return "@" + u.Username
 }
 return u.FullName
}

// AddMishka — функция для вручения подарка
func (u *User) AddMishka() {
 DataMutex.Lock()
 u.MishkaCount++
 DataMutex.Unlock()
}