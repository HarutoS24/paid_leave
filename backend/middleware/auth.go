package middleware

import (
	"backend/config"
	"backend/model"
	"backend/service"
	"context"
	"database/sql"
	"net/http"

	"github.com/gorilla/sessions"
)

func RequireLogin(store *sessions.FilesystemStore, db *sql.DB, next http.Handler) http.Handler {
	return RequireAuthSession(store,
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session := r.Context().Value(config.AuthContextKey).(*sessions.Session)
			isLoggedIn, err := service.IsLoggedIn(session)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if !isLoggedIn {
				http.Error(w, "ログインしてください", http.StatusUnauthorized)
				return
			}
			employee, err := service.GetLoggedInEmployee(db, session)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			ctx := context.WithValue(r.Context(), config.LoginUserContextKey, employee)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		}),
	)
}

func RequireAdminLogin(store *sessions.FilesystemStore, db *sql.DB, next http.Handler) http.Handler {
	return RequireLogin(store, db,
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			employee := r.Context().Value(config.LoginUserContextKey).(model.Employee)
			if !employee.IsAdmin {
				http.Error(w, "権限がありません", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		}),
	)
}
