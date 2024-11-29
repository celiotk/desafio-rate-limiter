package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"rate-limiter/internal/infra/database"
	"rate-limiter/internal/infra/web"
	"rate-limiter/internal/infra/web/middleware"
	"rate-limiter/internal/infra/web/webserver"
	"rate-limiter/usecase"
	"time"

	"rate-limiter/configs"
)

func main() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	configs, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}
	fmt.Printf("configs: %+v\n", configs)

	redisPool := database.NewRedisPool(configs.DbHost + ":" + configs.DbPort)
	storage := database.NewRedisStorage(redisPool)
	// storage := database.NewInMemoryStorage()

	apiStorage := database.NewTokenStorage(storage)
	ipStorage := database.NewIpRateStorage(storage)
	iru := usecase.NewIpRateUsecase(ipStorage, configs.IpRateLimit, configs.IpRateInterval, configs.IpBlockTime)
	tru := usecase.NewTokenRateUsecase(apiStorage, configs.DefaultTokenRateLimit, configs.DefaultTokenRateInterval, configs.TokenBlockTime)
	regToken := usecase.NewRegisterTokenUseCase(apiStorage)
	for token, par := range configs.TokensParsed {
		regToken.AddToken(usecase.RegisterTokenUseCaseInputDTO{
			Token:    token,
			Limit:    par.RateLimit,
			Interval: par.RateInterval,
		})
	}

	rateLimiterMiddleware := middleware.NewRateLimiterMiddleware(*iru, *tru)
	ws := webserver.NewWebServer(":"+configs.WEB_SERVER_PORT, rateLimiterMiddleware.RateLimit)
	hello := web.NewHelloHandler()
	ws.AddHandler("/hello", hello.Hello, http.MethodGet)
	fmt.Println("Starting web server on port ", configs.WEB_SERVER_PORT)
	go func() {
		if err := ws.Start(); err != nil {
			panic(err)
		}
	}()

	<-sigCh
	ctx2, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := ws.Stop(ctx2); err != nil {
		log.Fatal("Failed to stop server: %w", err)
	}
}
