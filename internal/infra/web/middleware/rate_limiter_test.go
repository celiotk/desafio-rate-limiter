package middleware

import (
	"net/http"
	"net/http/httptest"
	"rate-limiter/internal/infra/database"
	"rate-limiter/usecase"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRateLimitMiddleware(t *testing.T) {

	// redisPool := database.NewRedisPool("localhost:6379")
	// storage := database.NewRedisStorage(redisPool)
	storage := database.NewInMemoryStorage()

	tkStorage := database.NewTokenStorage(storage)
	ipStorage := database.NewIpRateStorage(storage)
	iru := usecase.NewIpRateUsecase(ipStorage, 10, 1, 2)
	tru := usecase.NewTokenRateUsecase(tkStorage, 10, 2, 1)

	regToken := usecase.NewRegisterTokenUseCase(tkStorage)
	regToken.AddToken(usecase.RegisterTokenUseCaseInputDTO{
		Token:    "token1",
		Limit:    10,
		Interval: 1,
	})
	regToken.AddToken(usecase.RegisterTokenUseCaseInputDTO{
		Token:    "token2",
		Limit:    20,
		Interval: 2,
	})

	middleware := NewRateLimiterMiddleware(*iru, *tru)

	handler := middleware.RateLimit(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	token1HttpReq := httptest.NewRequest(http.MethodGet, "/", nil)
	token1HttpReq.RemoteAddr = "192.168.1.1:12345"
	token1HttpReq.Header.Set("API_KEY", "token1")

	token2HttpReq := httptest.NewRequest(http.MethodGet, "/", nil)
	token2HttpReq.RemoteAddr = "192.168.1.1:12345"
	token2HttpReq.Header.Set("API_KEY", "token2")

	unknowTokenHttpReq := httptest.NewRequest(http.MethodGet, "/", nil)
	unknowTokenHttpReq.RemoteAddr = "192.168.1.1:12345"
	unknowTokenHttpReq.Header.Set("API_KEY", "unknowToken")

	ip1HttpReq := httptest.NewRequest(http.MethodGet, "/", nil)
	ip1HttpReq.RemoteAddr = "192.168.1.1:12345"

	ip2HttpReq := httptest.NewRequest(http.MethodGet, "/", nil)
	ip2HttpReq.RemoteAddr = "192.168.1.2:12345"

	// inicia a janela de tempo de contagem de requisições
	simulateSimultaneousRequests(handler, token1HttpReq, 1)
	simulateSimultaneousRequests(handler, token2HttpReq, 1)
	simulateSimultaneousRequests(handler, unknowTokenHttpReq, 1)
	simulateSimultaneousRequests(handler, ip1HttpReq, 1)
	simulateSimultaneousRequests(handler, ip2HttpReq, 1)

	time.Sleep(1500 * time.Millisecond)

	// estes deveriam estar na mesma janela de tempo
	countBlocked, countSucess := simulateSimultaneousRequests(handler, unknowTokenHttpReq, 30)
	assert.Equal(t, int32(21), countBlocked)
	assert.Equal(t, int32(9), countSucess)

	countBlocked, countSucess = simulateSimultaneousRequests(handler, token2HttpReq, 30)
	assert.Equal(t, int32(11), countBlocked)
	assert.Equal(t, int32(19), countSucess)

	// estes deveriam ter finalizado a janela de tempo e devem iniciar uma nova
	countBlocked, countSucess = simulateSimultaneousRequests(handler, token1HttpReq, 20)
	assert.Equal(t, int32(10), countBlocked)
	assert.Equal(t, int32(10), countSucess)

	countBlocked, countSucess = simulateSimultaneousRequests(handler, ip1HttpReq, 20)
	assert.Equal(t, int32(10), countBlocked)
	assert.Equal(t, int32(10), countSucess)

	countBlocked, countSucess = simulateSimultaneousRequests(handler, ip2HttpReq, 20)
	assert.Equal(t, int32(10), countBlocked)
	assert.Equal(t, int32(10), countSucess)

	time.Sleep(1500 * time.Millisecond)

	// este deveria estar bloqueado
	countBlocked, countSucess = simulateSimultaneousRequests(handler, ip1HttpReq, 20)
	assert.Equal(t, int32(20), countBlocked)
	assert.Equal(t, int32(0), countSucess)

	// estes deveriam estar desbloqueados
	countBlocked, countSucess = simulateSimultaneousRequests(handler, token1HttpReq, 10)
	assert.Equal(t, int32(0), countBlocked)
	assert.Equal(t, int32(10), countSucess)

	countBlocked, countSucess = simulateSimultaneousRequests(handler, token2HttpReq, 20)
	assert.Equal(t, int32(0), countBlocked)
	assert.Equal(t, int32(20), countSucess)

	countBlocked, countSucess = simulateSimultaneousRequests(handler, unknowTokenHttpReq, 10)
	assert.Equal(t, int32(0), countBlocked)
	assert.Equal(t, int32(10), countSucess)

	// verifica a mensagem de erro
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, token1HttpReq)
	assert.Equal(t, BLOCKED_MESSAGE, strings.TrimSpace(rr.Body.String()))
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, ip1HttpReq)
	assert.Equal(t, BLOCKED_MESSAGE, strings.TrimSpace(rr.Body.String()))
}

func simulateSimultaneousRequests(handler http.Handler, req *http.Request, count int) (countBlocked int32, countSucess int32) {
	wg := sync.WaitGroup{}
	wg.Add(count)
	for i := 0; i < count; i++ {
		go func() {
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)
			if rr.Code == http.StatusTooManyRequests {
				atomic.AddInt32(&countBlocked, 1)
			} else {
				atomic.AddInt32(&countSucess, 1)
			}
			wg.Done()
		}()
	}
	wg.Wait()
	return
}
