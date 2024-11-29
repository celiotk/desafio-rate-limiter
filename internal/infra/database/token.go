package database

import "rate-limiter/internal/entity"

type TokenStorage struct {
	storage RateLimiterStorage
}

func NewTokenStorage(storage RateLimiterStorage) *TokenStorage {
	return &TokenStorage{storage: storage}
}

func (t *TokenStorage) Increment(key string) (int, error) {
	return t.storage.Increment("token:count:" + key)
}

func (t *TokenStorage) SetExpiration(key string, ttl int) error {
	return t.storage.SetExpiration("token:count:"+key, ttl)
}

func (t *TokenStorage) GetTTL(key string) (int, error) {
	return t.storage.GetTTL("token:count:" + key)
}

func (t *TokenStorage) IsBlocked(key string) (bool, error) {
	return t.storage.Exists("token:block:" + key)
}

func (t *TokenStorage) Block(key string, duration int) error {
	return t.storage.Block("token:block:"+key, duration)
}

func (t *TokenStorage) SetTokenRateConfig(token string, cfg entity.TokenSettingsParam) error {
	return t.storage.SetTokenRateConfig("config:token:"+token, cfg.Limit, cfg.Interval)
}

func (t *TokenStorage) GetTokenRateConfig(key string) (entity.TokenSettings, error) {
	limit, interval, err := t.storage.GetTokenRateConfig("config:token:" + key)
	return entity.TokenSettings{Limit: limit, Interval: interval}, err
}
