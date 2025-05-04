package router

import (
	"DistanceTrackerServer/auth"
	"DistanceTrackerServer/utils"
	"context"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"time"
)

var (
	logRequest           = LogRequest
	addLoggingMiddleware = AddLoggingMiddleware
	sugarFromContext     = utils.SugarFromContext
	register             = auth.RegisterHandler
	login                = auth.LoginHandler
	logout               = auth.LogoutHandler
)

func formatRequestLogMessage(r *http.Request) string {
	return fmt.Sprintf("%s \t%s \t%s", r.RemoteAddr, r.Method, r.URL)
}

type LoggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
	requestID  string
}

func NewLoggingResponseWriter(w http.ResponseWriter) *LoggingResponseWriter {
	return &LoggingResponseWriter{w, http.StatusOK, uuid.New().String()}
}

func (lrw *LoggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func LogRequest(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		lrw := NewLoggingResponseWriter(w)
		sugar, _ := sugarFromContext(r.Context())
		sugar.Infof("--> %s \t%s", lrw.requestID, formatRequestLogMessage(r))
		start := time.Now()

		next.ServeHTTP(lrw, r)

		duration := time.Since(start)

		sugar.Infof("<-- %s \t%s \t%d \tRequest took %s", lrw.requestID, formatRequestLogMessage(r), lrw.statusCode, duration)
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
	handler = addLoggingMiddleware(ctx, logRequest(handler))

	server := &http.Server{
		Addr:           ":8080",
		Handler:        handler,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	mux.HandleFunc("/healthcheck", HealthCheckHandler)
	mux.HandleFunc("POST /register", register)
	mux.HandleFunc("POST /login", login)
	mux.HandleFunc("DELETE /logout", logout)

	err := server.ListenAndServe()
	if err != nil {
		return err
	}
	return nil
}
