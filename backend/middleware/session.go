package middleware

import (
	"backend/config"
	"context"
	"net/http"

	"github.com/gorilla/sessions"
)

func RequireAuthSession(store *sessions.FilesystemStore, next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			session, err := store.Get(r, config.AuthSessionName)
			if err != nil {
				session.Values = make(map[interface{}]interface{})
				session.Options.MaxAge = -1
				session.Save(r, w)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			ctx := context.WithValue(r.Context(), config.AuthContextKey, session)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		},
	)
}
