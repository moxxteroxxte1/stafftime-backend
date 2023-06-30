package handlers

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/moxxteroxxte1/stafftime-backend/src/models"
	"net/http"
	"os"
	"time"
)

func (s *APIServer) jwtAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		tokenString := c.Value

		token, tokenErr := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
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

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			var count int64
			s.database.Model(&models.User{}).Where("id = ?", uint(claims["UserID"].(float64))).Count(&count)
			if count < 1 {
				http.SetCookie(w, &http.Cookie{
					Name:    "token",
					Expires: time.Now(),
					Path:    "/",
				})
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			if c.Expires.Before(time.Now().Add(5*time.Minute)) || claims["ExpiresAt"].(time.Time) != c.Expires {

				exp := time.Now().Add(10 * time.Minute)
				newToken, newTokenErr := createJWT(claims["UserID"].(float64), claims["IsAdmin"].(bool), exp)
				if newTokenErr != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				http.SetCookie(w, &http.Cookie{
					Name:    "token",
					Value:   newToken,
					Expires: exp,
					Path:    "/",
				})
			}
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func createJWT(id float64, admin bool, exp time.Time) (string, error) {
	claims := &jwt.MapClaims{
		"UserID":    id,
		"IsAdmin":   admin,
		"ExpiresAt": jwt.NewNumericDate(exp),
	}

	secret := os.Getenv("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}

func validateToken(token string) (*jwt.Token, error) {
	secret := os.Getenv("JWT_SECRET")

	return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	})
}
