package entity

type TokenRateRepository interface {
	Increment(token string) (int, error)
	SetExpiration(token string, ttl int) error
	GetTTL(token string) (int, error)
	IsBlocked(token string) (bool, error)
	Block(token string, duration int) error
	SetTokenRateConfig(token string, cfg TokenSettingsParam) error
	GetTokenRateConfig(token string) (TokenSettings, error)
}

type IpRateRepository interface {
	Increment(ip string) (int, error)
	SetExpiration(ip string, ttl int) error
	GetTTL(ip string) (int, error)
	IsBlocked(ip string) (bool, error)
	Block(ip string, duration int) error
}
