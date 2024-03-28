package main

import (
	"fmt"

	"EXAM3/api-gateway/api"
	"EXAM3/api-gateway/config"
	"EXAM3/api-gateway/pkg/logger"
	"EXAM3/api-gateway/services"
	reds "EXAM3/api-gateway/storage/redis"

	"github.com/gomodule/redigo/redis"
)

// @title EXAM
// @version 0.1
// @description application description
// @securityDefinitions.apikey Token
// @in header
// @name Authorization
func main() {
	cfg := config.Load()
	log := logger.New(cfg.LogLevel, "api_gateway")
	pool := redis.Pool{
		MaxIdle:   80,
		MaxActive: 12000,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", fmt.Sprintf("%s:%d", cfg.RedisHost, cfg.RedisPort))
			if err != nil {
				panic(err.Error())
			}
			return c, err
		},
	}

	serviceManager, err := services.NewServiceManager(&cfg)
	if err != nil {
		log.Error("gRPC dial error", logger.Error(err))
	}

	server := api.New(api.Option{
		Conf:           cfg,
		Logger:         log,
		ServiceManager: serviceManager,
		Reds:           reds.NewRedisRepo(&pool),
	})

	if err := server.Run(cfg.HTTPPort); err != nil {
		log.Fatal("failed to run http server", logger.Error(err))
		panic(err)
	}
}
