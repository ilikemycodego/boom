package server

import (
	"log"
	"net/http"
	"time"
)

// statusWriter нужен, чтобы поймать код ответа (200/500/404...)
type statusWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

// RequestLogger — логирует каждый HTTP запрос + статус + время
func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sw := &statusWriter{ResponseWriter: w, status: 200}
		start := time.Now()

		next.ServeHTTP(sw, r)

		log.Printf("[HTTP] %s %s -> %d (%s)",
			r.Method, r.URL.RequestURI(), sw.status, time.Since(start))
	})
}
