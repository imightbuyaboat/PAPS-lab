package sessionmanager

import (
	"context"
	"fmt"
	"os"
	bt "papslab/basic_types"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

type SessionManager struct {
	Client *redis.Client
	ctx    context.Context
}

func NewSessionID() (bt.SessionID, error) {
	id, err := uuid.NewRandom()
	return bt.SessionID(id), err
}

func ParseSessionID(s string) (bt.SessionID, error) {
	id, err := uuid.Parse(s)
	return bt.SessionID(id), err
}

func NewSessionManager() (*SessionManager, error) {
	sm := &SessionManager{}

	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	host := os.Getenv("REDIS_HOST")
	port := os.Getenv("REDIS_PORT")
	password := os.Getenv("REDIS_PASSWORD")

	client := redis.NewClient(&redis.Options{
		Addr:     host + ":" + port,
		Password: password,
		DB:       0,
	})

	sm.Client = client
	sm.ctx = context.Background()

	if err := sm.Client.Ping(sm.ctx).Err(); err != nil {
		return nil, err
	}
	return sm, nil
}

func (sm *SessionManager) Create(s *bt.Session) (*bt.SessionID, error) {
	id, err := NewSessionID()
	if err != nil {
		return nil, err
	}

	sm.Client.HSet(sm.ctx, "session:"+id.String(), map[string]interface{}{
		"login":      s.Login,
		"useragent":  s.Useragent,
		"priveleged": strconv.FormatBool(s.Priveleged),
	})
	sm.Client.Expire(sm.ctx, "session:"+id.String(), 24*time.Hour)

	return &id, nil
}

func (sm *SessionManager) Check(id bt.SessionID) (*bt.Session, error) {
	data, err := sm.Client.HGetAll(sm.ctx, "session:"+id.String()).Result()
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, fmt.Errorf("сессия не найдена")
	}

	priveleged, err := strconv.ParseBool(data["priveleged"])
	if err != nil {
		return nil, fmt.Errorf("не удалось прочитать priveleged: %v", err)
	}

	session := &bt.Session{
		Login:      data["login"],
		Useragent:  data["useragent"],
		Priveleged: priveleged,
	}
	return session, nil
}

func (sm *SessionManager) Delete(id bt.SessionID) error {
	deleted, err := sm.Client.Del(sm.ctx, "session:"+id.String()).Result()
	if err != nil {
		return err
	}
	if deleted != 1 {
		return fmt.Errorf("ничего не удалено")
	}
	return nil
}
