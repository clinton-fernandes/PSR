package ratelimiter

import (
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"
)

var (
	endPoints = []string{"/GetBooks"}
	rateLimit = 3
)

type APIStats struct {
	RateLimit   int64
	RequestRate map[string]int64 //string is RemoteAddr
}

type RateLimiter struct {
	ApiDetails map[string]APIStats //string is api endpoint
	mux        sync.RWMutex
}

func setApiDetailsMap() map[string]APIStats {
	apiStats := make(map[string]APIStats)
	for _, val := range endPoints {
		apiStats[val] = APIStats{
			RateLimit:   int64(rateLimit),
			RequestRate: map[string]int64{},
		}
	}
	return apiStats
}

func InitializeRateLimiter() RateLimiter {
	return RateLimiter{
		ApiDetails: setApiDetailsMap(),
		mux:        sync.RWMutex{},
	}
}

func (r *RateLimiter) resetRequestRate() {
	r.mux.Lock()
	defer r.mux.Unlock()

	r.ApiDetails = setApiDetailsMap()
	//fmt.Println("Rate Limiter was Reset ...>>>")
}

func (l *RateLimiter) DefineResetInterval() {
	ticker := time.NewTicker(5 * time.Second) // reset after 5 sec
	go func() {
		for {
			select {
			case _ = <-ticker.C:
				l.resetRequestRate()
			}
		}
	}()
}

func (l *RateLimiter) inc(r *http.Request) error {
	l.mux.Lock()
	defer l.mux.Unlock()

	api := l.ApiDetails[r.RequestURI]
	if int64(api.RequestRate[r.RemoteAddr]) >= api.RateLimit {
		fmt.Println("Exceeded no. of requests for: " + r.RequestURI)
		return errors.New("Exceeded no. of requests for: " + r.RequestURI)
	}
	api.RequestRate[r.RemoteAddr] += 1
	return nil
}

func (l *RateLimiter) MiddleWare(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := l.inc(r); err != nil {
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte(err.Error()))
			return
		}
		handler.ServeHTTP(w, r)
	}
}

func StartRateLimiterService() {
	rateLimiter := InitializeRateLimiter()
	fmt.Println("rate limiter initialized...")

	rateLimiter.DefineResetInterval()

	http.HandleFunc("/GetBooks", rateLimiter.MiddleWare(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("/GetBooks executed.....")
		w.Write([]byte("/GetBooks executed....."))
	}))

	http.ListenAndServe(":9090", nil)
}
