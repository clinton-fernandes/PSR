package ratelimiter

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestStartRateLimiterService(t *testing.T) {
	req := httptest.NewRequest("GET", "/GetBooks", nil)
	w := httptest.NewRecorder()

	rateLimiter := InitializeRateLimiter()

	rateLimiter.DefineResetInterval()

	handler := rateLimiter.MiddleWare(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("executed test /GetBooks...")
	})

	req.RemoteAddr = "user-1"
	for i := 0; i < 20; i++ {
		time.Sleep(1 * time.Second)
		handler(w, req)
	}

}
