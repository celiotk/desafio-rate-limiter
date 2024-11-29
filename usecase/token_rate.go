package usecase

import (
	"fmt"
	"rate-limiter/internal/entity"
)

type TokenRateUsecase struct {
	rateRepository           entity.TokenRateRepository
	defaultRateLimit         int
	defaultRateLimitInterval int
	blockDuration            int
}

func NewTokenRateUsecase(tr entity.TokenRateRepository, defaultRateLimit, defaultRateLimitInterval, blockDuration int) *TokenRateUsecase {
	return &TokenRateUsecase{
		rateRepository:           tr,
		defaultRateLimit:         defaultRateLimit,
		defaultRateLimitInterval: defaultRateLimitInterval,
		blockDuration:            blockDuration,
	}
}

func (tru *TokenRateUsecase) Execute(token string) (blocked bool, err error) {
	rateLimit, rateLimitInterval := tru.defaultRateLimit, tru.defaultRateLimitInterval
	cfg, err := tru.rateRepository.GetTokenRateConfig(token)
	if err == nil {
		rateLimit = cfg.Limit
		rateLimitInterval = cfg.Interval
	}

	// Verifica se o token está bloqueado
	blocked, err = tru.rateRepository.IsBlocked(token)
	if err != nil {
		return false, err
	}
	if blocked {
		return true, nil
	}

	// Checa o número de requisições feitas pelo token
	reqCount, err := tru.rateRepository.Increment(token)
	if err != nil {
		return false, err
	}
	ttl, err := tru.rateRepository.GetTTL(token)
	if err != nil {
		return false, err
	}

	// Configura o TTL inicial se necessário
	if ttl == entity.TTL_NOT_SET {
		err = tru.rateRepository.SetExpiration(token, rateLimitInterval)
		if err != nil {
			return false, err
		}
	}

	// Verifica se o limite de requisições foi excedido
	tokenEntity := entity.NewTokenEntity(rateLimit, reqCount)
	if tokenEntity.LimitExceeded() {
		fmt.Printf("Token %s has reached the limit of %d requests in %d seconds\n", token, rateLimit, rateLimitInterval)
		err = tru.rateRepository.Block(token, tru.blockDuration) // Bloqueia o token pelo período configurado
		if err != nil {
			return true, err
		}
		return true, nil
	}

	// Token está dentro do limite
	return false, nil
}
