package session

import (
	"errors"
	"net/http"
	"share-profile-allocator/internal/utils"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// SessionData represents the data stored in a session
type SessionData struct {
	Data         map[string]any
	LastAccessed time.Time
	mu           sync.RWMutex
}

// SessionManager handles server-side session storage
type SessionManager struct {
	sessions map[string]*SessionData
	mu       sync.RWMutex
	maxAge   time.Duration

	cleanupSessions sync.Once
}

// NewSessionManager creates a new session manager
func NewSessionManager(maxAge time.Duration) *SessionManager {
	sm := &SessionManager{
		sessions: make(map[string]*SessionData),
		maxAge:   maxAge,
	}

	sm.cleanupSessionJob()

	return sm
}

// StartNewSession creates a new session and returns the session ID
// returns the sessionID/cookie
func (sm *SessionManager) StartNewSession() string {
	id := uuid.New().String()

	sm.mu.Lock()
	sm.sessions[id] = &SessionData{
		Data:         make(map[string]any),
		LastAccessed: time.Now(),
	}
	sm.mu.Unlock()

	utils.Log("7b960e5b").Info("Created new session", "id", id)
	return id
}

func (sm *SessionManager) GetSession(c echo.Context) (map[string]any, error) {
	id := c.Get("session_id").(string)

	sm.mu.RLock()
	session, ok := sm.sessions[id]
	sm.mu.RUnlock()

	if !ok {
		utils.Log("02386eab").Warn("Attempted to retrieve session that does not exist", "id", id)
		return nil, errors.New("session not found")
	}

	session.mu.Lock()
	now := time.Now()
	session.LastAccessed = now
	session.mu.Unlock()

	utils.Log("aeb640c8").Info("Retrieved session", "id", id, "LastAccessed", now)

	return session.Data, nil
}

// cleanupSessionJob removes expired sessions periodically
func (sm *SessionManager) cleanupSessionJob() {
	sm.cleanupSessions.Do(
		func() {
			cleanUpFreq := sm.maxAge / 2
			utils.Log("b684c505").Info("Starting session clean up job", "Clean up every", cleanUpFreq)

			go func() {
				ticker := time.NewTicker(cleanUpFreq)
				for range ticker.C {
					sm.mu.Lock()

					for id, session := range sm.sessions {
						if time.Since(session.LastAccessed) > sm.maxAge {
							utils.Log("39c3e699").Info("Deleting expired session", "id", id, "LastAccessed", session.LastAccessed)
							delete(sm.sessions, id)
						}
					}

					sm.mu.Unlock()
				}
			}()
		},
	)
}

// Middleware creates Echo middleware for session management
func (sm *SessionManager) Middleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			sessionID, err := c.Cookie("session_id")

			// Check if the cookie is stale
			if err == nil {
				sm.mu.RLock()
				_, ok := sm.sessions[sessionID.Value]
				sm.mu.RUnlock()

				if ok {
					c.Set("session_id", sessionID.Value)
					return next(c)
				}
			}

			// Need to create a new cookie and session
			newSessionID := sm.StartNewSession()

			cookie := new(http.Cookie)
			cookie.Name = "session_id"
			cookie.Value = newSessionID
			cookie.HttpOnly = true
			cookie.Secure = true // Enable in production
			cookie.SameSite = http.SameSiteStrictMode
			cookie.Path = "/"

			c.SetCookie(cookie)
			c.Set("session_id", newSessionID)

			return next(c)
		}
	}
}
