package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/moxxteroxxte1/stafftime-backend/src/models"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

func (s *APIServer) handleUsers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		makeHTPPHandler(s.GetAllUsers)(w, r)
	case http.MethodPost:
		makeHTPPHandler(s.CreateUser)(w, r)
	case http.MethodPut:
		makeHTPPHandler(s.UpdateAllUsers)(w, r)
	case http.MethodDelete:
		makeHTPPHandler(s.DeleteAllUsers)(w, r)
	default:
		WriteJSON(w, http.StatusOK, map[string]string{"message": "/user"})
	}
}

func (s *APIServer) GetAllUsers(w http.ResponseWriter, r *http.Request) error {
	users := []models.User{}
	s.database.Find(&users)
	return WriteJSON(w, http.StatusOK, users)
}

func (s *APIServer) CreateUser(w http.ResponseWriter, r *http.Request) error {
	user := new(models.User)
	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		return WriteJSON(w, http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("failed to decode user: %s", err)})
	}

	bytes, bcryptErr := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	if bcryptErr != nil {
		return WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("failed to create token: %s", bcryptErr)})
	}

	data := []byte(fmt.Sprintf(`{"password": "%s"}`, string(bytes)))
	jsonErr := json.Unmarshal(data, &user)
	if jsonErr != nil {
		return WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("failed to create token: %s", jsonErr)})
	}

	result := s.database.Create(&user)
	if result.Error != nil {
		return WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("failed to create user: %s", result.Error)})
	}

	return WriteJSON(w, http.StatusCreated, user)
}

func (s *APIServer) UpdateAllUsers(w http.ResponseWriter, r *http.Request) error {
	a := map[string]any{}
	if err := json.NewDecoder(r.Body).Decode(&a); err != nil {
		return WriteJSON(w, http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("failed to decode users: %s", err)})
	}

	fmt.Println(a)

	if a["password"] != nil {
		bytes, bcryptErr := bcrypt.GenerateFromPassword([]byte(a["password"].(string)), 14)
		if bcryptErr != nil {
			return WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("failed to create token: %s", bcryptErr)})
		}
		a["password"] = string(bytes)
	}

	user := models.User{}
	if err := json.NewDecoder(r.Body).Decode(&user); err == nil {
		log.Println(user)
	}

	result := s.database.Model(&models.User{}).Where("true").Updates(a)
	if result.Error != nil {
		return WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("failed to update users: %s", result.Error)})
	}

	return WriteJSON(w, http.StatusNoContent, nil)
}

func (s *APIServer) DeleteAllUsers(w http.ResponseWriter, r *http.Request) error {
	users := []models.User{}

	result := s.database.Find(&users).Delete(&users)
	if result.Error != nil {
		return WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("failed to delete all users: %s", result.Error)})
	}
	return WriteJSON(w, http.StatusNoContent, nil)
}

func (s *APIServer) HandleUserByID(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		makeHTPPHandler(s.GetUserByID)(w, r)
	case http.MethodPost:
		WriteJSON(w, http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("invalid method %s", r.Method)})
	case http.MethodPut:
		makeHTPPHandler(s.UpdateUserByID)(w, r)
	case http.MethodDelete:
		makeHTPPHandler(s.DeleteUserByID)(w, r)
	default:
		WriteJSON(w, http.StatusOK, map[string]string{"message": "/user/{id}"})
	}
}

func (s *APIServer) GetUserByID(w http.ResponseWriter, r *http.Request) error {
	id := mux.Vars(r)["userID"]
	user := models.User{}

	result := s.database.Where("id = ?", id).First(&user)
	if result.Error != nil {
		return WriteJSON(w, http.StatusNotFound, map[string]string{"error": fmt.Sprintf("failed to decode users: %s", result.Error)})
	}

	return WriteJSON(w, http.StatusOK, user)
}

func (s *APIServer) UpdateUserByID(w http.ResponseWriter, r *http.Request) error {
	id := mux.Vars(r)["userID"]
	a := map[string]any{}
	if err := json.NewDecoder(r.Body).Decode(&a); err != nil {
		return WriteJSON(w, http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("failed to decode users: %s", err)})
	}

	if a["password"] != nil {
		bytes, bcryptErr := bcrypt.GenerateFromPassword([]byte(a["password"].(string)), 14)
		if bcryptErr != nil {
			return WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("failed to create token: %s", bcryptErr)})
		}
		a["password"] = string(bytes)
	}

	result := s.database.Model(&models.User{}).Where("id = ?", id).Updates(a)
	if result.Error != nil {
		return WriteJSON(w, http.StatusNotFound, map[string]string{"error": fmt.Sprintf("failed to update users: %s", result.Error)})
	}

	return WriteJSON(w, http.StatusNoContent, nil)
}

func (s *APIServer) DeleteUserByID(w http.ResponseWriter, r *http.Request) error {
	id := mux.Vars(r)["userID"]

	result := s.database.Delete(&models.User{}, id)
	if result.Error != nil {
		return WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("failed to delete all users: %s", result.Error)})
	}
	return WriteJSON(w, http.StatusNoContent, nil)
}
