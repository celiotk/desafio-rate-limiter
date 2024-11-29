package database

type RateLimiterStorage interface {
	Increment(key string) (int, error)                                  // Incrementa e retorna o contador
	SetExpiration(key string, ttl int) error                            // Define o TTL para uma chave
	GetTTL(key string) (int, error)                                     // Retorna o TTL da chave
	Exists(key string) (bool, error)                                    // Verifica se uma chave existe
	Block(key string, duration int) error                               // Bloqueia uma chave por um período
	SetTokenRateConfig(token string, limit int, interval int) error     // Define a configuração de rate limit para um token
	GetTokenRateConfig(key string) (limit int, interval int, err error) // Retorna a configuração de rate limit para um token
}
