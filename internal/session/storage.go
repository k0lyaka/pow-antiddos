package session

import (
	"context"
	"errors"
	"time"

	"github.com/k0lyaka/pow-antiddos/internal/config"
	"github.com/k0lyaka/pow-antiddos/internal/redis"
	"github.com/k0lyaka/pow-antiddos/internal/utils"
)

type Session struct {
	Authorized bool `redis:"authorized"`

	Prefix string `redis:"prefix"`
}

func NewSession() (string, *Session) {
	sid := utils.RandStringBytes(32)
	prefix := utils.RandStringBytes(64)

	session := Session{
		Authorized: false,
		Prefix:     prefix,
	}

	ctx := context.Background()

	redis.Client.HSet(ctx, "session:"+sid, session).Err()
	redis.Client.Expire(ctx, "session:"+sid, time.Duration(config.Config.SessionTTL)*time.Second)

	return sid, &session
}

func GetSession(sid string) (*Session, error) {
	ctx := context.Background()

	data, err := redis.Client.HGetAll(ctx, "session:"+sid).Result()

	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return nil, errors.New("session not found")
	}

	session := &Session{
		Prefix:     data["prefix"],
		Authorized: data["authorized"] == "true",
	}

	return session, nil
}

func AuthorizeSession(sid string) error {
	ctx := context.Background()

	err := redis.Client.HSet(ctx, "session:"+sid, "authorized", "true").Err()
	redis.Client.Expire(ctx, "session:"+sid, time.Duration(config.Config.SessionTTL)*time.Second)

	if err != nil {
		return err
	}

	return nil
}
