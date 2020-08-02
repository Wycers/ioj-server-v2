package sessions

import (
	"encoding/json"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"time"
)

// Data represents the sessions.
type Session struct {
	AccountId uint64 `json:"accountId"`

	ExpTime time.Time `json:"exp_time"`
}

// Save saves the current sessions of the specified context.
func (sd *Session) Save(c *gin.Context) error {
	session := sessions.Default(c)
	sessionDataBytes, err := json.Marshal(sd)
	if err != nil {
		return err
	}
	session.Set("data", string(sessionDataBytes))
	return session.Save()
}

func (sd *Session) Clear(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	sd.ExpTime = time.Now()
	_ = sd.Save(c)
}

// GetSession returns sessions of the specified context.
func GetSession(c *gin.Context) *Session {
	session := sessions.Default(c)
	sessionDataStr := session.Get("data")
	if str, ok := sessionDataStr.(string); ok {
		res := &Session{}
		if err := json.Unmarshal([]byte(str), res); err != nil {
			return nil
		}
		if time.Now().After(res.ExpTime) {
			res.Clear(c)
			return nil
		}
		res.ExpTime = time.Now().Add(12 * time.Hour)
		return res
	} else {
		return nil
	}
}

func New() *Session {
	return &Session{
		ExpTime: time.Now().Add(12 * time.Hour),
	}
}
