package authorize

import (
	"crypto/rsa"
	"goapp/db"
	"goapp/myRedis"
)

type AuthHandlerInterface interface {
	GetPublicKey() *rsa.PublicKey
	GetAuthDatabase() db.UserDB
	GetRedis() myRedis.RedisPool
}

type authHandler struct {
	publicKey *rsa.PublicKey
	// インターフェースから定義してモックが作れるような設計する
	db    db.UserDB
	redis myRedis.RedisPool
}

func NewAuthHandlerInterface(publicKey *rsa.PublicKey, userDB db.UserDB, redis myRedis.RedisPool) AuthHandlerInterface {
	return &authHandler{
		publicKey: publicKey,
		db:        userDB,
		redis:     redis,
	}
}

func (h *authHandler) GetPublicKey() *rsa.PublicKey {
	return h.publicKey
}

func (h *authHandler) GetAuthDatabase() db.UserDB {
	return h.db
}

func (h *authHandler) GetRedis() myRedis.RedisPool {
	return h.redis
}
