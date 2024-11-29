package database

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestIpRate(t *testing.T) {
	// storage := NewRedisStorage(NewRedisPool("localhost:6379"))
	storage := NewInMemoryStorage()
	ipRateStorage := NewIpRateStorage(storage)

	key := "123456"
	key2 := "654321"

	// verifica valor do ttl quando a chave não existe
	ttl, err := ipRateStorage.GetTTL(key)
	assert.Nil(t, err)
	assert.Equal(t, -2, ttl)

	//  verifica se o incremento está funcionando
	reqCount, err := ipRateStorage.Increment(key)
	assert.Nil(t, err)
	assert.Equal(t, 1, reqCount)
	reqCount, err = ipRateStorage.Increment(key)
	assert.Nil(t, err)
	assert.Equal(t, 2, reqCount)
	reqCount, err = ipRateStorage.Increment(key2)
	assert.Nil(t, err)
	assert.Equal(t, 1, reqCount)

	// verifica retorno do ttl quando ainda não está configurado
	ttl, err = ipRateStorage.GetTTL(key)
	assert.Nil(t, err)
	assert.Equal(t, -1, ttl)

	// configura tempo de expiração e verifica se ocorreu no tempo correto e na chave correta
	err = ipRateStorage.SetExpiration(key, 1)
	assert.Nil(t, err)
	time.Sleep(750 * time.Millisecond)
	reqCount, err = ipRateStorage.Increment(key)
	assert.Nil(t, err)
	assert.Equal(t, 3, reqCount)
	time.Sleep(500 * time.Millisecond)
	reqCount, err = ipRateStorage.Increment(key2)
	assert.Nil(t, err)
	assert.Equal(t, 2, reqCount)
	reqCount, err = ipRateStorage.Increment(key)
	assert.Nil(t, err)
	assert.Equal(t, 1, reqCount)

	// remove as chaves usadas para teste
	err = ipRateStorage.SetExpiration(key, 0)
	assert.Nil(t, err)
	err = ipRateStorage.SetExpiration(key2, 0)
	assert.Nil(t, err)

	// verifica se o bloqueio e o desbloqueio estão funcionando
	err = ipRateStorage.Block(key, 1)
	assert.Nil(t, err)
	blocked, err := ipRateStorage.IsBlocked(key)
	assert.Nil(t, err)
	assert.True(t, blocked)
	bloqued, err := ipRateStorage.IsBlocked(key2)
	assert.Nil(t, err)
	assert.False(t, bloqued)
	time.Sleep(750 * time.Millisecond)
	blocked, err = ipRateStorage.IsBlocked(key)
	assert.Nil(t, err)
	assert.True(t, blocked)
	time.Sleep(500 * time.Millisecond)
	blocked, err = ipRateStorage.IsBlocked(key)
	assert.Nil(t, err)
	assert.False(t, blocked)
}
