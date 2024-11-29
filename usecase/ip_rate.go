package usecase

import (
	"fmt"
	"rate-limiter/internal/entity"
)

type IpRateUsecase struct {
	IpRate            entity.IpRateRepository
	rateLimit         int
	rateLimitInterval int
	blockDuration     int
}

func NewIpRateUsecase(ir entity.IpRateRepository, rateLimit, rateLimitInterval, blockDuration int) *IpRateUsecase {
	return &IpRateUsecase{
		IpRate:            ir,
		rateLimit:         rateLimit,
		rateLimitInterval: rateLimitInterval,
		blockDuration:     blockDuration,
	}
}

func (iru *IpRateUsecase) Execute(ip string) (blocked bool, err error) {

	// Verifica se o IP está bloqueado
	blocked, err = iru.IpRate.IsBlocked(ip)
	if err != nil {
		return false, err
	}
	if blocked {
		return true, nil
	}

	// Checa o número de requisições feitas pelo IP
	reqCount, err := iru.IpRate.Increment(ip)
	if err != nil {
		return false, err
	}
	ttl, err := iru.IpRate.GetTTL(ip)
	if err != nil {
		return false, err
	}

	// Configura o TTL inicial se necessário
	if ttl == entity.TTL_NOT_SET {
		err = iru.IpRate.SetExpiration(ip, iru.rateLimitInterval)
		if err != nil {
			return false, err
		}
	}

	// Verifica se o limite de requisições foi excedido
	ipEntity := entity.NewIpEntity(iru.rateLimit, reqCount)
	if ipEntity.LimitExceeded() {
		fmt.Printf("IP %s has reached the limit of %d requests in %d seconds\n", ip, iru.rateLimit, iru.rateLimitInterval)
		err = iru.IpRate.Block(ip, iru.blockDuration) // Bloqueia o IP pelo período configurado
		if err != nil {
			return true, err
		}
		return true, nil
	}

	// Retorna falso se o IP não foi bloqueado
	return false, nil
}
