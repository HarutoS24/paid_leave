package middleware

import (
	"backend/config"
	"context"
	"database/sql"
	"net/http"
)

func UseDB(db *sql.DB, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), config.DBContextKey, db)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
