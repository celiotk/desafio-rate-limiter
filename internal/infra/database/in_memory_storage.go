package database

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type InMemoryStorage struct {
	data        map[string]int
	expiry      map[string]expireData
	tokenConfig map[string]TokenConfig
	mu          sync.Mutex
}

type expireData struct {
	dateTime time.Time
	timer    *time.Timer
}

type TokenConfig struct {
	Limit    int
	Interval int
}

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		data:        make(map[string]int),
		expiry:      make(map[string]expireData),
		tokenConfig: make(map[string]TokenConfig),
	}
}

func (m *InMemoryStorage) Increment(key string) (int, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.data[key]++
	return m.data[key], nil
}

func (m *InMemoryStorage) SetExpiration(key string, ttl int) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	data, exist := m.expiry[key]
	if exist {
		data.timer.Stop()
	}

	timer := time.AfterFunc(time.Duration(ttl)*time.Second, func() {
		m.mu.Lock()
		defer m.mu.Unlock()
		delete(m.data, key)
		delete(m.expiry, key)
	})

	m.expiry[key] = expireData{dateTime: time.Now().Add(time.Duration(ttl) * time.Second), timer: timer}
	return nil
}

func (m *InMemoryStorage) GetTTL(key string) (int, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.data[key]; !exists {
		return -2, nil
	}

	if data, exists := m.expiry[key]; exists {
		ttl := int(time.Until(data.dateTime).Seconds())
		if ttl < 0 {
			return 0, nil
		}
		return ttl, nil
	}
	return -1, nil
}

func (m *InMemoryStorage) Exists(key string) (bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	_, exists := m.data[key]
	return exists, nil
}

func (m *InMemoryStorage) Block(key string, duration int) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	fmt.Printf("Blocking key: %s for %d seconds\n", key, duration)

	m.data[key] = 1
	timer := time.AfterFunc(time.Duration(duration)*time.Second, func() {
		m.mu.Lock()
		defer m.mu.Unlock()
		delete(m.data, key)
		delete(m.expiry, key)
		fmt.Printf("Unblocking key: %s\n", key)
	})

	m.expiry[key] = expireData{dateTime: time.Now().Add(time.Duration(duration) * time.Second), timer: timer}
	return nil
}

func (m *InMemoryStorage) SetTokenRateConfig(token string, limit int, interval int) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.tokenConfig[token] = TokenConfig{Limit: limit, Interval: interval}
	return nil
}

func (m *InMemoryStorage) GetTokenRateConfig(token string) (int, int, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if config, exists := m.tokenConfig[token]; exists {
		return config.Limit, config.Interval, nil
	}
	return 0, 0, errors.New("token not found")
}
