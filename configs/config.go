package configs

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

type config struct {
	WEB_SERVER_PORT          string `mapstructure:"WEB_SERVER_PORT"`
	DbHost                   string `mapstructure:"DB_HOST"`
	DbPort                   string `mapstructure:"DB_PORT"`
	IpRateLimit              int    `mapstructure:"IP_RATE_LIMIT"`
	IpRateInterval           int    `mapstructure:"IP_RATE_INTERVAL"`
	IpBlockTime              int    `mapstructure:"IP_BLOCK_TIME"`
	DefaultTokenRateLimit    int    `mapstructure:"DEFAULT_TOKEN_RATE_LIMIT"`
	DefaultTokenRateInterval int    `mapstructure:"DEFAULT_TOKEN_RATE_INTERVAL"`
	TokenBlockTime           int    `mapstructure:"TOKEN_BLOCK_TIME"`
	Tokens                   string `mapstructure:"TOKENS"`
	TokensParsed             map[string]Token
}

type Token struct {
	RateLimit    int
	RateInterval int
}

func LoadConfig(path string) (*config, error) {
	cfg := &config{}
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	// Parse TokensString into Tokens map
	cfg.TokensParsed, err = parseTokens(cfg.Tokens)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func parseTokens(tokensString string) (map[string]Token, error) {
	tokens := make(map[string]Token)
	pairs := strings.Split(tokensString, ",")
	for _, pair := range pairs {
		kv := strings.Split(pair, ":")
		if len(kv) != 2 {
			return nil, fmt.Errorf("invalid token format: %s", pair)
		}
		values := strings.Split(kv[1], "/")
		if len(values) != 2 {
			return nil, fmt.Errorf("invalid token values format: %s", kv[1])
		}
		rateLimit, err := strconv.Atoi(values[0])
		if err != nil {
			return nil, fmt.Errorf("invalid rate limit value: %s", values[0])
		}
		rateInterval, err := strconv.Atoi(values[1])
		if err != nil {
			return nil, fmt.Errorf("invalid rate interval value: %s", values[1])
		}
		tokens[kv[0]] = Token{
			RateLimit:    rateLimit,
			RateInterval: rateInterval,
		}
	}
	return tokens, nil
}
