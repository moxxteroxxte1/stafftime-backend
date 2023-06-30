package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/moxxteroxxte1/stafftime-backend/src/models"
	"log"
	"net/http"
)

func (s *APIServer) handleKeys(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		makeHTPPHandler(s.GetAllKeys)(w, r)
	case http.MethodPost:
		makeHTPPHandler(s.CreateKey)(w, r)
	case http.MethodPut:
		makeHTPPHandler(s.UpdateAllKeys)(w, r)
	case http.MethodDelete:
		makeHTPPHandler(s.DeleteAllKeys)(w, r)
	default:
		WriteJSON(w, http.StatusOK, map[string]string{"message": "/keys"})
	}
}

func (s *APIServer) GetAllKeys(w http.ResponseWriter, r *http.Request) error {
	keys := []models.Key{}
	s.database.Find(&keys)
	return WriteJSON(w, http.StatusOK, keys)
}

func (s *APIServer) CreateKey(w http.ResponseWriter, r *http.Request) error {
	key := new(models.Key)
	if err := json.NewDecoder(r.Body).Decode(key); err != nil {
		return WriteJSON(w, http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("failed to decode key: %s", err)})
	}

	id := uuid.New()
	data := []byte(fmt.Sprintf(`{"key": "%s"}`, id.String()))
	jsonErr := json.Unmarshal(data, &key)
	if jsonErr != nil {
		return WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("failed to create token: %s", jsonErr)})
	}

	result := s.database.Create(&key)
	if result.Error != nil {
		return WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("failed to create key: %s", result.Error)})
	}

	return WriteJSON(w, http.StatusCreated, key)
}

func (s *APIServer) UpdateAllKeys(w http.ResponseWriter, r *http.Request) error {
	a := map[string]any{}
	if err := json.NewDecoder(r.Body).Decode(&a); err != nil {
		return WriteJSON(w, http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("failed to decode keys: %s", err)})
	}

	key := models.Key{}
	if err := json.NewDecoder(r.Body).Decode(&key); err == nil {
		log.Println(key)
	}

	result := s.database.Model(&models.Key{}).Where("true").Updates(a)
	if result.Error != nil {
		return WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("failed to update keys: %s", result.Error)})
	}

	return WriteJSON(w, http.StatusNoContent, nil)
}

func (s *APIServer) DeleteAllKeys(w http.ResponseWriter, r *http.Request) error {
	keys := []models.Key{}

	result := s.database.Find(&keys).Delete(&keys)
	if result.Error != nil {
		return WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("failed to delete all keys: %s", result.Error)})
	}
	return WriteJSON(w, http.StatusNoContent, nil)
}

func (s *APIServer) HandleKeysByID(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		makeHTPPHandler(s.GetKeyByID)(w, r)
	case http.MethodPost:
		WriteJSON(w, http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("invalid method %s", r.Method)})
	case http.MethodPut:
		makeHTPPHandler(s.UpdateKeyByID)(w, r)
	case http.MethodDelete:
		makeHTPPHandler(s.DeleteKeyByID)(w, r)
	default:
		WriteJSON(w, http.StatusOK, map[string]string{"message": "/key/{id}"})
	}
}

func (s *APIServer) GetKeyByID(w http.ResponseWriter, r *http.Request) error {
	id := mux.Vars(r)["keyID"]
	key := models.Key{}

	result := s.database.Where("id = ?", id).First(&key)
	if result.Error != nil {
		return WriteJSON(w, http.StatusNotFound, map[string]string{"error": fmt.Sprintf("failed to decode keys: %s", result.Error)})
	}

	return WriteJSON(w, http.StatusOK, key)
}

func (s *APIServer) UpdateKeyByID(w http.ResponseWriter, r *http.Request) error {
	id := mux.Vars(r)["keyID"]
	a := map[string]any{}
	if err := json.NewDecoder(r.Body).Decode(&a); err != nil {
		return WriteJSON(w, http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("failed to decode keys: %s", err)})
	}

	result := s.database.Model(&models.Key{}).Where("id = ?", id).Updates(a)
	if result.Error != nil {
		return WriteJSON(w, http.StatusNotFound, map[string]string{"error": fmt.Sprintf("failed to update keys: %s", result.Error)})
	}

	return WriteJSON(w, http.StatusNoContent, nil)
}

func (s *APIServer) DeleteKeyByID(w http.ResponseWriter, r *http.Request) error {
	id := mux.Vars(r)["keyID"]

	result := s.database.Delete(&models.Key{}, id)
	if result.Error != nil {
		return WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("failed to delete all keys: %s", result.Error)})
	}
	return WriteJSON(w, http.StatusNoContent, nil)
}
