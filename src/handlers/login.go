package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"

	"github.com/moxxteroxxte1/stafftime-backend/src/models"
	"github.com/moxxteroxxte1/stafftime-backend/src/types"

	"golang.org/x/crypto/bcrypt"
)

// GET/POST/PUT/DELETE ALL

func (s *APIServer) HandleLogin(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		makeHTPPHandler(s.login)(w, r)
	default:
		WriteJSON(w, http.StatusUnauthorized, "Unauthorized")
	}
}

func (s *APIServer) login(w http.ResponseWriter, r *http.Request) error {
	var req types.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		resetCookie(w)
		return WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("failed to decode request: %s", err)})
	}

	user := new(models.User)
	var result = s.database.Where("username = ?", req.Username).First(&user)
	if result.Error != nil {
		resetCookie(w)
		return WriteJSON(w, http.StatusNotFound, map[string]string{"error": fmt.Sprintf("failed to decode users: %s", result.Error)})
	}

	if !(bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)) == nil) {
		resetCookie(w)
		return WriteJSON(w, http.StatusUnauthorized, "401 Unauthorized")
	}

	var expirationTime time.Time
	if !req.KeepLoggedIn {
		expirationTime = time.Now().Add(10 * time.Minute)
	} else {
		expirationTime = time.Date(9999, 12, 31, 23, 59, 59, 0, time.UTC)
	}

	claims := &jwt.MapClaims{
		"UserID":    user.ID,
		"IsAdmin":   user.IsAdmin,
		"ExpiresAt": jwt.NewNumericDate(expirationTime),
	}

	secret := os.Getenv("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		resetCookie(w)
		return WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("failed to create token: %s", err)})
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})
	return WriteJSON(w, http.StatusNoContent, nil)
}

func resetCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Expires: time.Now(),
	})
}
