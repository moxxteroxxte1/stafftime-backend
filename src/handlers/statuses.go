package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/moxxteroxxte1/stafftime-backend/src/models"
	"log"
	"net/http"
)

func (s *APIServer) handleStatus(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		makeHTPPHandler(s.GetAllStatuses)(w, r)
	case http.MethodPost:
		makeHTPPHandler(s.CreateStatus)(w, r)
	case http.MethodPut:
		makeHTPPHandler(s.UpdateAllStatuses)(w, r)
	case http.MethodDelete:
		makeHTPPHandler(s.DeleteAllStatuses)(w, r)
	default:
		WriteJSON(w, http.StatusOK, map[string]string{"message": "/status"})
	}
}

func (s *APIServer) GetAllStatuses(w http.ResponseWriter, r *http.Request) error {
	statuses := []models.Status{}
	s.database.Find(&statuses)
	return WriteJSON(w, http.StatusOK, statuses)
}

func (s *APIServer) CreateStatus(w http.ResponseWriter, r *http.Request) error {
	status := new(models.Status)
	if err := json.NewDecoder(r.Body).Decode(status); err != nil {
		return WriteJSON(w, http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("failed to decode status: %s", err)})
	}

	result := s.database.Create(&status)
	if result.Error != nil {
		return WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("failed to create status: %s", result.Error)})
	}

	return WriteJSON(w, http.StatusCreated, status)
}

func (s *APIServer) UpdateAllStatuses(w http.ResponseWriter, r *http.Request) error {
	a := map[string]any{}
	if err := json.NewDecoder(r.Body).Decode(&a); err != nil {
		return WriteJSON(w, http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("failed to decode statuses: %s", err)})
	}

	status := models.Status{}
	if err := json.NewDecoder(r.Body).Decode(&status); err == nil {
		log.Println(status)
	}

	result := s.database.Model(&models.Status{}).Where("true").Updates(a)
	if result.Error != nil {
		return WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("failed to update statuses: %s", result.Error)})
	}

	return WriteJSON(w, http.StatusNoContent, nil)
}

func (s *APIServer) DeleteAllStatuses(w http.ResponseWriter, r *http.Request) error {
	statuses := []models.Status{}

	result := s.database.Find(&statuses).Delete(&statuses)
	if result.Error != nil {
		return WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("failed to delete all statuses: %s", result.Error)})
	}
	return WriteJSON(w, http.StatusNoContent, nil)
}

func (s *APIServer) HandleStatusByID(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		makeHTPPHandler(s.GetStatusByID)(w, r)
	case http.MethodPost:
		WriteJSON(w, http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("invalid method %s", r.Method)})
	case http.MethodPut:
		makeHTPPHandler(s.UpdateStatusByID)(w, r)
	case http.MethodDelete:
		makeHTPPHandler(s.DeleteStatusByID)(w, r)
	default:
		WriteJSON(w, http.StatusOK, map[string]string{"message": "/status/{id}"})
	}
}

func (s *APIServer) GetStatusByID(w http.ResponseWriter, r *http.Request) error {
	id := mux.Vars(r)["statusID"]
	status := models.Status{}

	result := s.database.Where("id = ?", id).First(&status)
	if result.Error != nil {
		return WriteJSON(w, http.StatusNotFound, map[string]string{"error": fmt.Sprintf("failed to decode statuses: %s", result.Error)})
	}

	return WriteJSON(w, http.StatusOK, status)
}

func (s *APIServer) UpdateStatusByID(w http.ResponseWriter, r *http.Request) error {
	id := mux.Vars(r)["statusID"]
	a := map[string]any{}
	if err := json.NewDecoder(r.Body).Decode(&a); err != nil {
		return WriteJSON(w, http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("failed to decode statuses: %s", err)})
	}

	result := s.database.Model(&models.Status{}).Where("id = ?", id).Updates(a)
	if result.Error != nil {
		return WriteJSON(w, http.StatusNotFound, map[string]string{"error": fmt.Sprintf("failed to update statuses: %s", result.Error)})
	}

	return WriteJSON(w, http.StatusNoContent, nil)
}

func (s *APIServer) DeleteStatusByID(w http.ResponseWriter, r *http.Request) error {
	id := mux.Vars(r)["statusID"]

	result := s.database.Delete(&models.Status{}, id)
	if result.Error != nil {
		return WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("failed to delete all statuses: %s", result.Error)})
	}
	return WriteJSON(w, http.StatusNoContent, nil)
}
