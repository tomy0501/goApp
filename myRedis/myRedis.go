package myRedis

import (
	"context"
	"encoding/json"
	"log/slog"

	redigo "github.com/gomodule/redigo/redis"
)

type RedisPool interface {
	Save(id string, values map[string]any) error
	Load(id string) (map[string]any, error)
	Key(id string) string
	Delete(ctx context.Context, id string) error
}

type prefixRedisStoreWithPool struct {
	pool      *redigo.Pool
	prefix    string
	ttl       int
	secretKey []byte
}

func NewPrefixRedisStoreWithPool(pool *redigo.Pool, prefix string, ttl int, secretKey []byte) RedisPool {
	return &prefixRedisStoreWithPool{
		pool:      pool,
		prefix:    prefix,
		ttl:       ttl,
		secretKey: secretKey,
	}
}

func (s *prefixRedisStoreWithPool) Key(id string) string {
	return s.prefix + ":" + id
}

func (s *prefixRedisStoreWithPool) Save(id string, values map[string]any) error {
	data, err := json.Marshal(values)
	if err != nil {
		return err
	}
	encrytData, err := encrypt(data, s.secretKey)
	if err != nil {
		return err
	}
	conn := s.pool.Get()
	_, err = conn.Do("SET", id, encrytData, "EX", s.ttl)
	if err != nil {
		return err
	}
	return nil
}

func (s *prefixRedisStoreWithPool) Load(id string) (map[string]any, error) {
	conn := s.pool.Get()

	reply, err := redigo.Bytes(conn.Do("GET", id))
	if err == redigo.ErrNil {
		slog.Info("reply is nil")
		return make(map[string]any), nil
	} else if err != nil {
		slog.Info("reply is not nil")
		return nil, err
	}
	decryptReply, err := decrypt(reply, s.secretKey)
	if err != nil {
		return nil, err
	}
	var rawMap map[string]any
	err = json.Unmarshal(decryptReply, &rawMap)
	if err != nil {
		return nil, err
	}
	return rawMap, nil
}

func (s *prefixRedisStoreWithPool) Delete(ctx context.Context, id string) error {
	con := s.pool.Get()
	_, err := con.Do("DEL", id)
	if err != nil {
		slog.Error("DELETE Redis session is failed", "error", err)
		return err
	}

	return nil
}
