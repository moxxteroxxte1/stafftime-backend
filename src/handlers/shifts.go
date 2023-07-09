package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/moxxteroxxte1/stafftime-backend/src/models"
	"net/http"
	"strconv"
)

// GET/POST/PUT/DELETE ALL

func (s *APIServer) handleShifts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		makeHTPPHandler(s.GetAllShifts)(w, r)
	case http.MethodPost:
		makeHTPPHandler(s.CreateShift)(w, r)
	case http.MethodPut:
		makeHTPPHandler(s.UpdateAllShifts)(w, r)
	case http.MethodDelete:
		makeHTPPHandler(s.DeleteAllShifts)(w, r)
	default:
		err := WriteJSON(w, http.StatusOK, map[string]string{"message": "/payments"})
		if err != nil {
			return
		}
	}
}

func (s *APIServer) GetAllShifts(w http.ResponseWriter, r *http.Request) error {
	shifts := []models.Shift{}
	s.database.Find(&shifts)
	return WriteJSON(w, http.StatusOK, shifts)
}

func (s *APIServer) CreateShift(w http.ResponseWriter, r *http.Request) error {
	shift := new(models.Shift)
	if err := json.NewDecoder(r.Body).Decode(shift); err != nil {
		return WriteJSON(w, http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("failed to decode shift: %s", err)})
	}

	diff := shift.EndTime.Sub(shift.StartTime).Hours()
	data := []byte(fmt.Sprintf(`{"hours": %f}`, diff))
	jsonErr := json.Unmarshal(data, &shift)
	if jsonErr != nil {
		return WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("failed to create token: %s", jsonErr)})
	}

	result := s.database.Create(&shift)
	if result.Error != nil {
		return WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("failed to create shift: %s", result.Error)})
	}

	return WriteJSON(w, http.StatusCreated, shift)
}

func (s *APIServer) UpdateAllShifts(w http.ResponseWriter, r *http.Request) error {
	a := map[string]any{}
	if err := json.NewDecoder(r.Body).Decode(&a); err != nil {
		return WriteJSON(w, http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("failed to decode shift: %s", err)})
	}

	result := s.database.Model(&models.Shift{}).Where("true").Updates(a)
	if result.Error != nil {
		return WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("failed to update shift: %s", result.Error)})
	}

	return WriteJSON(w, http.StatusNoContent, nil)
}

func (s *APIServer) DeleteAllShifts(w http.ResponseWriter, r *http.Request) error {
	shifts := []models.Shift{}

	result := s.database.Find(&shifts).Delete(&shifts)
	if result.Error != nil {
		return WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("failed to delete all shift: %s", result.Error)})
	}
	return WriteJSON(w, http.StatusNoContent, nil)
}

// GET/POST/PUT/DELETE BY USER (TODO)

func (s *APIServer) HandleShiftByUserID(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		makeHTPPHandler(s.GetShiftByUserID)(w, r)
	case http.MethodPost:
		makeHTPPHandler(s.CreateShifByUserID)(w, r)
	case http.MethodPut:
		makeHTPPHandler(s.UpdateShiftByUserID)(w, r)
	case http.MethodDelete:
		makeHTPPHandler(s.DeleteShiftByUserID)(w, r)
	default:
		err := WriteJSON(w, http.StatusOK, map[string]string{"message": "/user/{id}/payments"})
		if err != nil {
			return
		}
	}
}

func (s *APIServer) GetShiftByUserID(w http.ResponseWriter, r *http.Request) error {
	id := mux.Vars(r)["userID"]
	shifts := []models.Shift{}

	result := s.database.Where("user_id = ?", id).Find(&shifts)
	if result.Error != nil {
		return WriteJSON(w, http.StatusNotFound, map[string]string{"error": fmt.Sprintf("failed to load shift by userId: %s", result.Error)})
	}

	return WriteJSON(w, http.StatusOK, shifts)
}

func (s *APIServer) CreateShifByUserID(w http.ResponseWriter, r *http.Request) error {
	shift := new(models.Shift)

	if err1 := json.NewDecoder(r.Body).Decode(shift); err1 != nil {
		return WriteJSON(w, http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("failed to decode shift: %s", err1)})
	}

	i := mux.Vars(r)["userID"]
	userID, err := strconv.Atoi(i)
	if err != nil {
		return WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("failed to create shift: %s", err)})
	}

	diff := shift.EndTime.Sub(shift.StartTime).Hours()

	data := []byte(fmt.Sprintf(`{"userID": %d, "hours": %f}`, uint(userID), diff))
	jsonErr := json.Unmarshal(data, &shift)
	if jsonErr != nil {
		return WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("failed to create shift: %s", jsonErr)})
	}

	result := s.database.Create(&shift)
	if result.Error != nil {
		return WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("failed to create shift: %s", result.Error)})
	}

	return WriteJSON(w, http.StatusCreated, shift)
}

func (s *APIServer) UpdateShiftByUserID(w http.ResponseWriter, r *http.Request) error {
	id := mux.Vars(r)["userID"]
	a := map[string]any{}
	if err := json.NewDecoder(r.Body).Decode(&a); err != nil {
		return WriteJSON(w, http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("failed to decode shift: %s", err)})
	}

	result := s.database.Model(&models.Shift{}).Where("user_id = ?", id).Updates(a)
	if result.Error != nil {
		return WriteJSON(w, http.StatusNotFound, map[string]string{"error": fmt.Sprintf("failed to update shift: %s", result.Error)})
	}

	return WriteJSON(w, http.StatusNoContent, nil)
}

func (s *APIServer) DeleteShiftByUserID(w http.ResponseWriter, r *http.Request) error {
	id := mux.Vars(r)["userID"]

	result := s.database.Delete(&models.Shift{}, "user_id = ?", id)
	if result.Error != nil {
		return WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("failed to delete all shift: %s", result.Error)})
	}
	return WriteJSON(w, http.StatusNoContent, nil)
}

// GET/PULL/DELETE BY ID

func (s *APIServer) HandleShiftByID(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		makeHTPPHandler(s.GetShiftByID)(w, r)
	case http.MethodPost:
		err := WriteJSON(w, http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("invalid method %s", r.Method)})
		if err != nil {
			return
		}
	case http.MethodPut:
		makeHTPPHandler(s.UpdateShiftByID)(w, r)
	case http.MethodDelete:
		makeHTPPHandler(s.DeleteShiftByID)(w, r)
	default:
		err := WriteJSON(w, http.StatusOK, map[string]string{"message": "/shifts/{id}"})
		if err != nil {
			return
		}
	}
}

func (s *APIServer) GetShiftByID(w http.ResponseWriter, r *http.Request) error {
	id := mux.Vars(r)["shiftID"]
	shift := models.Shift{}

	result := s.database.Where("ID = ?", id).First(&shift)
	if result.Error != nil {
		return WriteJSON(w, http.StatusNotFound, map[string]string{"error": fmt.Sprintf("failed to load shift: %s", result.Error)})
	}

	return WriteJSON(w, http.StatusOK, shift)
}

func (s *APIServer) UpdateShiftByID(w http.ResponseWriter, r *http.Request) error {
	id := mux.Vars(r)["shiftID"]
	a := map[string]any{}
	if err := json.NewDecoder(r.Body).Decode(&a); err != nil {
		return WriteJSON(w, http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("failed to decode shift: %s", err)})
	}

	result := s.database.Model(&models.Shift{}).Where("id = ?", id).Updates(a)
	if result.Error != nil {
		return WriteJSON(w, http.StatusNotFound, map[string]string{"error": fmt.Sprintf("failed to update shift: %s", result.Error)})
	}

	return WriteJSON(w, http.StatusNoContent, nil)
}

func (s *APIServer) DeleteShiftByID(w http.ResponseWriter, r *http.Request) error {
	id := mux.Vars(r)["shiftID"]

	result := s.database.Delete(&models.Shift{}, id)
	if result.Error != nil {
		return WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("failed to delete all shift: %s", result.Error)})
	}
	return WriteJSON(w, http.StatusNoContent, nil)
}
