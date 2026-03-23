package main

import "sync"

type User struct {
 ID    string
 Email string
 Theme string
}

type Msg struct {
 Sender string
 Text   string
}

type Entity struct {
 Name      string
 Members   map[string]bool
 Messages  []Msg
}

var (
 users    = make(map[string]*User)
 byID     = make(map[string]*User)
 entities = make(map[string]*Entity)
 mu       sync.Mutex
)