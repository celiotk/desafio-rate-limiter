package database

import (
	"rate-limiter/internal/entity"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestToken(t *testing.T) {
	// storage := NewRedisStorage(NewRedisPool("localhost:6379"))
	storage := NewInMemoryStorage()
	tokenStorage := NewTokenStorage(storage)

	key := "123456"
	key2 := "654321"

	// verifica valor do ttl quando a chave não existe
	ttl, err := tokenStorage.GetTTL(key)
	assert.Nil(t, err)
	assert.Equal(t, -2, ttl)

	//  verifica se o incremento está funcionando
	reqCount, err := tokenStorage.Increment(key)
	assert.Nil(t, err)
	assert.Equal(t, 1, reqCount)
	reqCount, err = tokenStorage.Increment(key)
	assert.Nil(t, err)
	assert.Equal(t, 2, reqCount)
	reqCount, err = tokenStorage.Increment(key2)
	assert.Nil(t, err)
	assert.Equal(t, 1, reqCount)

	// verifica retorno do ttl quando ainda não está configurado
	ttl, err = tokenStorage.GetTTL(key)
	assert.Nil(t, err)
	assert.Equal(t, -1, ttl)

	// configura tempo de expiração e verifica se ocorreu no tempo correto e na chave correta
	err = tokenStorage.SetExpiration(key, 1)
	assert.Nil(t, err)
	time.Sleep(750 * time.Millisecond)
	reqCount, err = tokenStorage.Increment(key)
	assert.Nil(t, err)
	assert.Equal(t, 3, reqCount)
	time.Sleep(500 * time.Millisecond)
	reqCount, err = tokenStorage.Increment(key2)
	assert.Nil(t, err)
	assert.Equal(t, 2, reqCount)
	reqCount, err = tokenStorage.Increment(key)
	assert.Nil(t, err)
	assert.Equal(t, 1, reqCount)

	// remove as chaves usadas para teste
	err = tokenStorage.SetExpiration(key, 0)
	assert.Nil(t, err)
	err = tokenStorage.SetExpiration(key2, 0)
	assert.Nil(t, err)

	// verifica se o bloqueio e o desbloqueio estão funcionando
	err = tokenStorage.Block(key, 1)
	assert.Nil(t, err)
	blocked, err := tokenStorage.IsBlocked(key)
	assert.Nil(t, err)
	assert.True(t, blocked)
	bloqued, err := tokenStorage.IsBlocked(key2)
	assert.Nil(t, err)
	assert.False(t, bloqued)
	time.Sleep(750 * time.Millisecond)
	blocked, err = tokenStorage.IsBlocked(key)
	assert.Nil(t, err)
	assert.True(t, blocked)
	time.Sleep(500 * time.Millisecond)
	blocked, err = tokenStorage.IsBlocked(key)
	assert.Nil(t, err)
	assert.False(t, blocked)

	// verifica se está salvando os parâmetros do token corretamente
	err = tokenStorage.SetTokenRateConfig(key, entity.TokenSettingsParam{Limit: 10, Interval: 10})
	assert.Nil(t, err)
	tokenConfig, err := tokenStorage.GetTokenRateConfig(key)
	assert.Nil(t, err)
	assert.Equal(t, 10, tokenConfig.Limit)
	assert.Equal(t, 10, tokenConfig.Interval)
}
