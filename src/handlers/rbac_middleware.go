package handlers

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"strconv"
)

func (s *APIServer) rbacMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("token")
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		token, tokenErr := jwt.Parse(c.Value, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})
		if tokenErr != nil {
			if tokenErr == jwt.ErrSignatureInvalid {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		val, muxOk := mux.Vars(r)["userID"]
		id, _ := strconv.ParseUint(val, 10, 32)
		if claims, ok := token.Claims.(jwt.MapClaims); !ok || (muxOk && (!claims["IsAdmin"].(bool) && uint(claims["UserID"].(float64)) != uint(id))) || !muxOk && !claims["IsAdmin"].(bool) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
