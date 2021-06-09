package sessionstore

import (
	"errors"

	"github.com/google/uuid"
)

var (
	// ErrRecordNotFound ...
	ErrSessionNotFound = errors.New("session not found")
)

type SessionStore struct {
	WEBsessions   *[]Session
	ADMINsessions *[]Session
}

func New() *SessionStore {
	WEBsessions := &[]Session{}
	ADMINsessions := &[]Session{}
	return &SessionStore{
		WEBsessions:   WEBsessions,
		ADMINsessions: ADMINsessions,
	}
}

func (store *SessionStore) Sessions(isAdmin bool) *[]Session {
	if isAdmin {
		return store.ADMINsessions
	}
	return store.WEBsessions
}

func (store *SessionStore) Create(userID string, isAdmin bool) string {
	if isAdmin {
		session := Session{
			ID:     store.generateSessionID(store.getAllSessionID(true)),
			UserID: userID,
		}
		*store.Sessions(true) = append(*store.Sessions(true), session)
		return session.ID
	}
	session := Session{
		ID:     store.generateSessionID(store.getAllSessionID(false)),
		UserID: userID,
	}
	*store.Sessions(false) = append(*store.Sessions(false), session)
	return session.ID
}

func (store *SessionStore) getAllSessionID(isAdmin bool) *[]string {
	if isAdmin {
		sessionID := []string{}
		for _, session := range *store.ADMINsessions {
			sessionID = append(sessionID, session.ID)
		}
		return &sessionID
	}
	sessionID := []string{}
	for _, session := range *store.WEBsessions {
		sessionID = append(sessionID, session.ID)
	}
	return &sessionID
}

func (store *SessionStore) generateSessionID(idArray *[]string) string {
	uid := uuid.New().String()
	for _, id := range *idArray {
		if uid == id {
			uid = store.generateSessionID(idArray)
		}
	}
	return uid
}

func (store *SessionStore) Get(sessionID string, isAdmin bool) (*Session, error) {
	if isAdmin {
		for _, session := range *store.ADMINsessions {
			if session.ID == sessionID {
				return &session, nil
			}
		}
		return nil, ErrSessionNotFound
	}
	for _, session := range *store.WEBsessions {
		if session.ID == sessionID {
			return &session, nil
		}
	}
	return nil, ErrSessionNotFound
}

func (store *SessionStore) Delete(sessionID string, isAdmin bool) {
	if isAdmin {
		SessionStore := store.Sessions(true)
		for i := 0; i < len(*SessionStore); i++ {
			if sessionID == (*SessionStore)[i].ID {
				*SessionStore = append((*SessionStore)[:i], (*SessionStore)[i+1:]...)
			}
		}
	}
	SessionStore := store.Sessions(false)
	for i := 0; i < len(*SessionStore); i++ {
		if sessionID == (*SessionStore)[i].ID {
			*SessionStore = append((*SessionStore)[:i], (*SessionStore)[i+1:]...)
		}
	}
}
