package main

import (
	"log"
	"net/http"
)

type Logger struct {
	log           *log.Logger
	wrapedHandler http.Handler
}

func (l *Logger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	l.wrapedHandler.ServeHTTP(w, r)
	l.log.Print(r.Method, r.URL.Path)

}

func NewLoggingMiddleWare(next http.Handler) *Logger {

	return &Logger{
		log:           log.Default(),
		wrapedHandler: next,
	}

}
