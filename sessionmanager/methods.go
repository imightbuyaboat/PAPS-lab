package sessionmanager

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

func NewSessionID() (SessionID, error) {
	id, err := uuid.NewRandom()
	return SessionID(id), err
}

func ParseSessionID(s string) (SessionID, error) {
	id, err := uuid.Parse(s)
	return SessionID(id), err
}

func (sid SessionID) String() string {
	return uuid.UUID(sid).String()
}

func NewSessionManager() (*SessionManager, error) {
	err := godotenv.Load() // Загружаем .env
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

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}
	return &SessionManager{client}, nil
}

func (sm *SessionManager) Create(s *Session) (*SessionID, error) {
	id, err := NewSessionID()
	if err != nil {
		return nil, err
	}

	sm.HSet(ctx, "session:"+id.String(), map[string]interface{}{
		"login":      s.Login,
		"useragent":  s.Useragent,
		"priveleged": strconv.FormatBool(s.Priveleged),
	})
	sm.Expire(ctx, "session:"+id.String(), 24*time.Hour)

	return &id, nil
}

func (sm *SessionManager) Check(id SessionID) (*Session, error) {
	data, err := sm.HGetAll(ctx, "session:"+id.String()).Result()
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
	deleted, err := sm.Del(ctx, "session:"+id.String()).Result()
	if err != nil {
		return err
	}
	if deleted != 1 {
		return fmt.Errorf("ничего не удалено")
	}
	return nil
}
