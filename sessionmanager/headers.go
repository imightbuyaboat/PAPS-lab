package sessionmanager

import (
	"context"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type Session struct {
	Login      string
	Useragent  string
	Priveleged bool
}

type SessionID uuid.UUID

type SessionManager struct {
	*redis.Client
}

var ctx = context.Background()
