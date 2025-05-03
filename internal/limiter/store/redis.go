package store

import (
	"context"
	"fmt"
	"time"

	"rate-limiter/internal/config"

	"github.com/redis/go-redis/v9"
)

type RedisStore struct {
	client *redis.Client
}

func NewRedisStore(cfg *config.Config) (*RedisStore, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &RedisStore{client: rdb}, nil
}

func (r *RedisStore) Increment(key string, limit int, blockSeconds int) (bool, time.Duration, error) {
	ctx := context.Background()

	// Verifica se está bloqueado
	blockKey := fmt.Sprintf("block:%s", key)
	blockExists, err := r.client.Exists(ctx, blockKey).Result()
	if err != nil {
		return false, 0, err
	}
	if blockExists == 1 {
		blockTTL, _ := r.client.TTL(ctx, blockKey).Result()
		return false, blockTTL, nil
	}

	// Incrementa o contador
	count, err := r.client.Incr(ctx, key).Result()
	if err != nil {
		return false, 0, err
	}

	// Expira em 1 segundo se for a primeira requisição
	if count == 1 {
		err := r.client.Expire(ctx, key, time.Second).Err()
		if err != nil {
			return false, 0, err
		}
	}

	// Se passou do limite, aplica bloqueio com SetNX
	if count > int64(limit) {
		ok, err := r.client.SetNX(ctx, blockKey, "1", time.Duration(blockSeconds)*time.Second).Result()
		if err != nil {
			return false, 0, err
		}
		if ok {
			// Bloqueio foi aplicado com sucesso
			return false, time.Duration(blockSeconds) * time.Second, nil
		}

		// Já estava bloqueado
		blockTTL, _ := r.client.TTL(ctx, blockKey).Result()
		return false, blockTTL, nil
	}

	// Ainda dentro do limite
	ttl, _ := r.client.TTL(ctx, key).Result()
	return true, ttl, nil
}

func (r *RedisStore) IsBlocked(key string) (bool, time.Duration, error) {
	ctx := context.Background()
	blockKey := fmt.Sprintf("block:%s", key)
	exists, err := r.client.Exists(ctx, blockKey).Result()
	if err != nil {
		return false, 0, err
	}
	if exists == 1 {
		ttl, _ := r.client.TTL(ctx, blockKey).Result()
		return true, ttl, nil
	}
	return false, 0, nil
}

func (s *RedisStore) GetRequestCount(key string) (int64, error) {
	count, err := s.client.Get(context.Background(), key).Int64()
	if err != nil && err != redis.Nil {
		return 0, fmt.Errorf("error getting request count: %v", err)
	}
	return count, nil
}

func (r *RedisStore) ResetKey(key string) error {
	ctx := context.Background()
	blockKey := fmt.Sprintf("block:%s", key)
	return r.client.Del(ctx, key, blockKey).Err()
}
