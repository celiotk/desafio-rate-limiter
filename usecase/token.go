package usecase

import "rate-limiter/internal/entity"

type RegisterTokenUseCase struct {
	repository entity.TokenRateRepository
}

type RegisterTokenUseCaseInputDTO struct {
	Token    string
	Limit    int
	Interval int
}

func NewRegisterTokenUseCase(repository entity.TokenRateRepository) *RegisterTokenUseCase {
	return &RegisterTokenUseCase{repository: repository}
}

func (uc *RegisterTokenUseCase) AddToken(input RegisterTokenUseCaseInputDTO) error {
	return uc.repository.SetTokenRateConfig(input.Token, entity.TokenSettingsParam{
		Token:    input.Token,
		Limit:    input.Limit,
		Interval: input.Interval,
	})
}
