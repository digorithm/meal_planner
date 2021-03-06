// Package middlewares provides common middleware handlers.
package middlewares

import (
	"context"
	"github.com/digorithm/meal_planner/finchgo"
	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"
	"strings"
	"time"
)

func SetDB(db *sqlx.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			req = req.WithContext(context.WithValue(req.Context(), "db", db))

			next.ServeHTTP(res, req)
		})
	}
}

func SetSessionStore(sessionStore sessions.Store) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			req = req.WithContext(context.WithValue(req.Context(), "sessionStore", sessionStore))

			next.ServeHTTP(res, req)
		})
	}
}

// MustLogin is a middleware that checks existence of current user.
func MustLogin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		sessionStore := req.Context().Value("sessionStore").(sessions.Store)
		session, _ := sessionStore.Get(req, "meal_planner-session")
		userRowInterface := session.Values["user"]

		if userRowInterface == nil {
			http.Redirect(res, req, "/login", 302)
			return
		}

		next.ServeHTTP(res, req)
	})
}

func Log(Finch *finchgo.Finch) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			startTime := time.Now()

			next.ServeHTTP(w, r)

			duration := time.Now().Sub(startTime)
			// Convert to ms
			durationMS := float64(duration.Nanoseconds() / int64(1000000))
			basePath := strings.Split(r.URL.Path, "/")[1]

			// Inject monitoring loop
			Finch.MonitorWorkload(r.Method, basePath)
			Finch.MonitorLatency(r.Method, basePath, durationMS)
			Finch.MonitorKnobs()

			log.Printf("Address:: %s -- '%s' '%s' -- response time:: %v", r.RemoteAddr, r.Method, r.URL, duration)
		})
	}
}
