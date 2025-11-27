package middleware

import (
	"backend/config"
	"context"
	"fmt"
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
				if saveErr := session.Save(r, w); saveErr != nil {
					err = fmt.Errorf("%w reset error: %s", err, saveErr)
				}
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			ctx := context.WithValue(r.Context(), config.AuthContextKey, session)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		},
	)
}
