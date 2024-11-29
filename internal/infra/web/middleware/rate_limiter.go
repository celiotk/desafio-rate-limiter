package middleware

import (
	"net/http"
	"rate-limiter/usecase"
	"strings"
)

const BLOCKED_MESSAGE = "you have reached the maximum number of requests or actions allowed within a certain time frame"

type rateLimitMiddleware struct {
	IpRateUsecase    usecase.IpRateUsecase
	TokenRateUsecase usecase.TokenRateUsecase
}

func NewRateLimiterMiddleware(iru usecase.IpRateUsecase, tru usecase.TokenRateUsecase) *rateLimitMiddleware {
	return &rateLimitMiddleware{
		IpRateUsecase:    iru,
		TokenRateUsecase: tru,
	}
}

func (rl rateLimitMiddleware) RateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("API_KEY")
		if token != "" {
			blocked, err := rl.TokenRateUsecase.Execute(token)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if blocked {
				http.Error(w, BLOCKED_MESSAGE, http.StatusTooManyRequests)
				return
			}
		} else {
			ip := strings.Split(r.RemoteAddr, ":")[0]
			blocked, err := rl.IpRateUsecase.Execute(ip)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if blocked {
				http.Error(w, BLOCKED_MESSAGE, http.StatusTooManyRequests)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}
