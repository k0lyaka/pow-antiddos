package session

import (
	"errors"
	"sync"

	"github.com/k0lyaka/pow-antiddos/internal/utils"
)

var sessionMu sync.Mutex = sync.Mutex{}
var sessions map[string]*Session = map[string]*Session{}

type Session struct {
	Authorized bool

	Prefix string
}

func NewSession() (string, *Session) {
	sessionMu.Lock()
	defer sessionMu.Unlock()

	sid := utils.RandStringBytes(32)
	prefix := utils.RandStringBytes(64)

	session := Session{
		Authorized: false,
		Prefix:     prefix,
	}
	sessions[sid] = &session

	return sid, &session
}

func GetSession(sid string) (*Session, error) {
	sessionMu.Lock()
	defer sessionMu.Unlock()

	session, exists := sessions[sid]

	if !exists {
		return nil, errors.New("session not found")
	}

	return session, nil
}

func AuthorizeSession(sid string) error {
	sessionMu.Lock()
	defer sessionMu.Unlock()

	session, exists := sessions[sid]

	if !exists {
		return errors.New("session not found")
	}

	session.Authorized = true
	return nil
}
