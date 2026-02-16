package redis

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/labstack/gommon/log"
	"github.com/redis/go-redis/v9"
)

type Service struct {
	redis *redis.Client
}

var (
	redisAddress = os.Getenv("REDIS_ADDRESS")
	redisUser    = os.Getenv("REDIS_USER")
	redisPass    = os.Getenv("REDIS_PASS")
	redisDB      = os.Getenv("REDIS_DB")
)

func New() *Service {
	db, err := strconv.Atoi(redisDB)
	if err != nil {
		db = 0 // Default to DB 0 if conversion fails
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s", redisAddress),
		Username: redisUser,
		Password: redisPass,
		DB:       db,
	})

	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("failed to connect to Redis: %v", err)
	}

	log.Info("connected to Redis")

	return &Service{
		redis: rdb,
	}
}

func (s *Service) Health() map[string]string {
	ctx := context.Background()

	_, err := s.redis.Ping(ctx).Result()
	if err != nil {
		log.Error(err.Error())
		return map[string]string{
			"status": "unhealthy",
			"error":  err.Error(),
		}
	}

	return map[string]string{
		"status": "healthy",
	}
}

func (s *Service) Close() error {
	return s.redis.Close()
}
