package session_manager

import (
	"context"
	"fmt"
	"os"
	"papslab/session"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

type RedisSessionManager struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisSessionManager() (*RedisSessionManager, error) {
	sm := &RedisSessionManager{}

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

func (sm *RedisSessionManager) Create(s *session.Session) (*session.SessionID, error) {
	id, err := session.NewSessionID()
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

func (sm *RedisSessionManager) Check(id session.SessionID) (*session.Session, error) {
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

	session := &session.Session{
		Login:      data["login"],
		Useragent:  data["useragent"],
		Priveleged: priveleged,
	}
	return session, nil
}

func (sm *RedisSessionManager) Delete(id session.SessionID) error {
	deleted, err := sm.client.Del(sm.ctx, "session:"+id.String()).Result()
	if err != nil {
		return err
	}
	if deleted != 1 {
		return fmt.Errorf("ничего не удалено")
	}
	return nil
}
