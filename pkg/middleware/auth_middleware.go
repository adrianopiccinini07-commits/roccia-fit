package middleware

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Accesso Negato: Token mancante", http.StatusUnauthorized)
			return
		}

		// Il formato standard e' "Bearer <tuo_token>"
		parti := strings.Split(authHeader, " ")
		if len(parti) != 2 || parti[0] != "Bearer" {
			http.Error(w, "Accesso Negato: Formato token non valido", http.StatusUnauthorized)
			return
		}

		tokenString := parti[1]
		chiaveSegreta := os.Getenv("JWT_SECRET")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, http.ErrAbortHandler
			}
			return []byte(chiaveSegreta), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Accesso Negato: Token non valido o scaduto", http.StatusUnauthorized)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			ctx := context.WithValue(r.Context(), "utente_id", claims["utente_id"])
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			http.Error(w, "Accesso negato: Dati token illeggibili", http.StatusUnauthorized)
		}

	})
}
