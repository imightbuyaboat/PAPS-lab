package sessionmanager

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

type SessionManager struct {
	client *redis.Client
	ctx    context.Context
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
		Addr:         host + ":" + port,
		Password:     password,
		DB:           0,
		PoolSize:     20,
		MinIdleConns: 5,
	})

	sm.client = client
	sm.ctx = context.Background()

	if err := sm.client.Ping(sm.ctx).Err(); err != nil {
		return nil, err
	}
	return sm, nil
}

func (sm *SessionManager) Create(s *Session) (*SessionID, error) {
	id, err := NewSessionID()
	if err != nil {
		return nil, err
	}

	sm.client.HSet(sm.ctx, "session:"+id.String(), map[string]interface{}{
		"login":      s.Login,
		"useragent":  s.Useragent,
		"priveleged": strconv.FormatBool(s.Priveleged),
	})
	sm.client.Expire(sm.ctx, "session:"+id.String(), 24*time.Hour)

	return &id, nil
}

func (sm *SessionManager) Check(id SessionID) (*Session, error) {
	data, err := sm.client.HGetAll(sm.ctx, "session:"+id.String()).Result()
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

	session := &Session{
		Login:      data["login"],
		Useragent:  data["useragent"],
		Priveleged: priveleged,
	}
	return session, nil
}

func (sm *SessionManager) Delete(id SessionID) error {
	deleted, err := sm.client.Del(sm.ctx, "session:"+id.String()).Result()
	if err != nil {
		return err
	}
	if deleted != 1 {
		return fmt.Errorf("ничего не удалено")
	}
	return nil
}
