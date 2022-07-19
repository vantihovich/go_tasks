package middleware

import (
	"context"
	"net/http"

	"github.com/golang-jwt/jwt"
	log "github.com/sirupsen/logrus"
	usershandlers "github.com/vantihovich/go_tasks/tree/master/swagger/usershandlers"
)

type claims struct {
	UserID int `json:"user_id"`
	jwt.StandardClaims
}

func Authorize(cfg string, f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			log.Error("empty header authorization data")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		claims := &claims{}
		tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(cfg), nil
		})
		if err != nil {
			if err.Error() == jwt.ErrSignatureInvalid.Error() {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			log.WithError(err).Info("Error parsing the token")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if !tkn.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		contextKeyUserID := usershandlers.ContextKeyUserID
		ctx := context.WithValue(r.Context(), contextKeyUserID, claims.UserID)
		f.ServeHTTP(w, r.WithContext(ctx))
	}
}
