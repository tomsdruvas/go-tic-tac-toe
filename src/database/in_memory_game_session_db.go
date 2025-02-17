package database

import (
	"errors"
	"src/src/models"
	"sync"
)

var (
	instance *InMemoryGameSessionDB
	once     sync.Once
)

func GetInstance() *InMemoryGameSessionDB {
	once.Do(func() {
		instance = NewInMemoryGameSessionDB()
	})
	return instance
}

type InMemoryGameSessionDB struct {
	sessions map[string]models.GameSession
	mu       sync.RWMutex
}

func NewInMemoryGameSessionDB() *InMemoryGameSessionDB {
	return &InMemoryGameSessionDB{
		sessions: make(map[string]models.GameSession),
	}
}

func (db *InMemoryGameSessionDB) StoreSession(session models.GameSession) {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.sessions[session.SessionId] = session
}

func (db *InMemoryGameSessionDB) GetSession(id string) (models.GameSession, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	session, exists := db.sessions[id]
	if !exists {
		return models.GameSession{}, errors.New("session not found")
	}
	return session, nil
}

func (db *InMemoryGameSessionDB) Clear() {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.sessions = make(map[string]models.GameSession)
}
