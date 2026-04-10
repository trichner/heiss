package session

import (
	"crypto/subtle"
	"errors"
	"sync/atomic"
)

type Session struct {
	ID int
}

type Manager interface {
	GetSession(token string) *Session
	Login(password string) (*Session, error)
}

func NewInMemorySessionManager(password string) *InMemorySessionManager {
	return &InMemorySessionManager{sessions: make(map[string]*Session), password: password}
}

type InMemorySessionManager struct {
	sessions    map[string]*Session
	idIncrement int32

	password string
}

func (i *InMemorySessionManager) GetSession(token string) *Session {
	return i.sessions[token]
}

func (i *InMemorySessionManager) Login(password string) (*Session, error) {
	// IMPORTANT: constant time comparison to not leak information about the password!
	if subtle.ConstantTimeCompare([]byte(i.password), []byte(password)) != 1 {
		return nil, errors.New("bad password")
	}

	id := atomic.AddInt32(&i.idIncrement, 1)
	s := &Session{ID: int(id)}
	return s, nil
}
