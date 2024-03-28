package v1

import (
	t "EXAM3/api-gateway/api/handlers/v1/tokens"
	"EXAM3/api-gateway/config"
	"EXAM3/api-gateway/pkg/logger"
	"EXAM3/api-gateway/services"
	"EXAM3/api-gateway/storage/repo"

	"github.com/casbin/casbin/v2"
)

const (
	ErrorCodeInvalidURL          = "INVALID_URL"
	ErrorCodeInvalidJSON         = "INVALID_JSON"
	ErrorCodeInvalidParams       = "INVALID_PARAMS"
	ErrorCodeInternalServerError = "INTERNAL_SERVER_ERROR"
	ErrorCodeUnauthorized        = "UNAUTHORIZED"
	ErrorCodeAlreadyExists       = "ALREADY_EXISTS"
	ErrorCodeNotFound            = "NOT_FOUND"
	ErrorCodeInvalidCode         = "INVALID_CODE"
	ErrorBadRequest              = "BAD_REQUEST"
	ErrorInvalidCredentials      = "INVALID_CREDENTIALS"
)

type handlerV1 struct {
	log            logger.Logger
	serviceManager services.IServiceManager
	cfg            config.Config
	reds           repo.RedisStorageI
	casbin         *casbin.Enforcer
	jwtHandler     t.JWTHandler
}

type HandlerV1Config struct {
	Logger         logger.Logger
	ServiceManager services.IServiceManager
	Cfg            config.Config
	Reds           repo.RedisStorageI
	Casbin         *casbin.Enforcer
	JWTHandler     t.JWTHandler
}

func New(c *HandlerV1Config) *handlerV1 {
	return &handlerV1{
		log:            c.Logger,
		serviceManager: c.ServiceManager,
		cfg:            c.Cfg,
		reds:           c.Reds,
		casbin:         c.Casbin,
		jwtHandler:     c.JWTHandler,
	}
}
