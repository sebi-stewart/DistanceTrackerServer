package router

import (
	"awesomeProject/utils"
	"context"
	"fmt"
	"net/http"
	"time"
)

var (
	logIncomingRequest   = LogIncomingRequest
	logRequestTime       = LogRequestTime
	addLoggingMiddleware = AddLoggingMiddleware
	sugarFromContext     = utils.SugarFromContext
)

func formatLogMessage(r *http.Request) string {
	return fmt.Sprintf("%s \t%s \t%s", r.RemoteAddr, r.Method, r.URL)
}

func LogIncomingRequest(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sugar, _ := sugarFromContext(r.Context())
		sugar.Info(formatLogMessage(r))
		next.ServeHTTP(w, r)
	}
}

func LogRequestTime(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		duration := time.Since(start)
		sugar, _ := sugarFromContext(r.Context())
		sugar.Info(fmt.Sprintf("%s \tRequest took %s", formatLogMessage(r), duration))
	}
}

func AddLoggingMiddleware(ctx context.Context, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r = r.Clone(ctx)
		next.ServeHTTP(w, r)
	}
}

func HealthCheckHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("OK"))
}

func Init(ctx context.Context) error {

	mux := http.NewServeMux()
	var handler http.Handler = mux
	handler = addLoggingMiddleware(ctx, logIncomingRequest(handler))
	handler = addLoggingMiddleware(ctx, logRequestTime(handler))

	server := &http.Server{
		Addr:           ":8080",
		Handler:        handler,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	mux.HandleFunc("/healthcheck", HealthCheckHandler)
	mux.HandleFunc("GET /")
	err := server.ListenAndServe()
	if err != nil {
		return err
	}
	return nil
}
