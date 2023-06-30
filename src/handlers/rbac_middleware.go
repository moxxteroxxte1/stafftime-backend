package handlers

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"github.com/moxxteroxxte1/stafftime-backend/src/models"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func (s *APIServer) rbacMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Has("key") {
			key := models.Key{}
			s.database.Where(&models.Key{Key: r.URL.Query().Get("key")}).First(&key)
			log.Println(key.ExpiresAt)
			if key.ExpiresAt.Before(time.Now()) {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
			return
		}

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
