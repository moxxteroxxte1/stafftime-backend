package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/moxxteroxxte1/stafftime-backend/src/models"
)

// GET/POST/PUT/DELETE ALL

func (s *APIServer) handlePayments(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		makeHTPPHandler(s.GetAllPayments)(w, r)
	case http.MethodPost:
		makeHTPPHandler(s.CreatePayment)(w, r)
	case http.MethodPut:
		makeHTPPHandler(s.UpdateAllPayments)(w, r)
	case http.MethodDelete:
		makeHTPPHandler(s.DeleteAllPayments)(w, r)
	default:
		WriteJSON(w, http.StatusOK, map[string]string{"message": "/payments"})
	}
}

func (s *APIServer) GetAllPayments(w http.ResponseWriter, r *http.Request) error {
	payments := []models.Payment{}
	s.database.Find(&payments)
	return WriteJSON(w, http.StatusOK, payments)
}

func (s *APIServer) CreatePayment(w http.ResponseWriter, r *http.Request) error {
	payment := new(models.Payment)
	if err := json.NewDecoder(r.Body).Decode(payment); err != nil {
		return WriteJSON(w, http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("failed to decode payment: %s", err)})
	}

	result := s.database.Create(&payment)
	if result.Error != nil {
		return WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("failed to create payment: %s", result.Error)})
	}

	return WriteJSON(w, http.StatusCreated, payment)
}

func (s *APIServer) UpdateAllPayments(w http.ResponseWriter, r *http.Request) error {
	a := map[string]any{}
	if err := json.NewDecoder(r.Body).Decode(&a); err != nil {
		return WriteJSON(w, http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("failed to decode payment: %s", err)})
	}

	result := s.database.Model(&models.Payment{}).Where("true").Updates(a)
	if result.Error != nil {
		return WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("failed to update payment: %s", result.Error)})
	}

	return WriteJSON(w, http.StatusNoContent, nil)
}

func (s *APIServer) DeleteAllPayments(w http.ResponseWriter, r *http.Request) error {
	payments := []models.Payment{}

	result := s.database.Find(&payments).Delete(&payments)
	if result.Error != nil {
		return WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("failed to delete all payments: %s", result.Error)})
	}
	return WriteJSON(w, http.StatusNoContent, nil)
}

// GET/PUT/DELETE BY USER (TODO)

func (s *APIServer) HandlePaymentByUserID(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		makeHTPPHandler(s.GetPaymentsByUserID)(w, r)
	case http.MethodPost:
		makeHTPPHandler(s.CreatePaymentByUserID)(w, r)
	case http.MethodPut:
		makeHTPPHandler(s.UpdatePaymentByUserID)(w, r)
	case http.MethodDelete:
		makeHTPPHandler(s.DeletePaymentByUserID)(w, r)
	default:
		WriteJSON(w, http.StatusOK, map[string]string{"message": "/user/{id}/payments"})
	}
}

func (s *APIServer) GetPaymentsByUserID(w http.ResponseWriter, r *http.Request) error {
	id := mux.Vars(r)["userID"]
	payments := []models.Payment{}

	result := s.database.Where("user_id = ?", id).Find(&payments)
	if result.Error != nil {
		return WriteJSON(w, http.StatusNotFound, map[string]string{"error": fmt.Sprintf("failed to load payment by userId: %s", result.Error)})
	}

	return WriteJSON(w, http.StatusOK, payments)
}

func (s *APIServer) CreatePaymentByUserID(w http.ResponseWriter, r *http.Request) error {
	payment := new(models.Payment)

	i := mux.Vars(r)["userID"]
	userID, err := strconv.Atoi(i)
	if err != nil {
		return WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("failed to create payment: %s", err)})
	}

	data := []byte(fmt.Sprintf(`{"userID": %d}`, uint(userID)))
	jsonErr := json.Unmarshal(data, &payment)
	if jsonErr != nil {
		return WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("failed to create payment: %s", jsonErr)})
	}

	if err1 := json.NewDecoder(r.Body).Decode(payment); err1 != nil {
		return WriteJSON(w, http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("failed to decode payment: %s", err1)})
	}

	result := s.database.Create(&payment)
	if result.Error != nil {
		return WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("failed to create payment: %s", result.Error)})
	}

	return WriteJSON(w, http.StatusCreated, payment)
}

func (s *APIServer) UpdatePaymentByUserID(w http.ResponseWriter, r *http.Request) error {
	id := mux.Vars(r)["userID"]
	a := map[string]any{}
	if err := json.NewDecoder(r.Body).Decode(&a); err != nil {
		return WriteJSON(w, http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("failed to decode payment: %s", err)})
	}

	result := s.database.Model(&models.Payment{}).Where("user_id = ?", id).Updates(a)
	if result.Error != nil {
		return WriteJSON(w, http.StatusNotFound, map[string]string{"error": fmt.Sprintf("failed to update payment: %s", result.Error)})
	}

	return WriteJSON(w, http.StatusNoContent, nil)
}

func (s *APIServer) DeletePaymentByUserID(w http.ResponseWriter, r *http.Request) error {
	id := mux.Vars(r)["userID"]

	result := s.database.Delete(&models.Payment{}, "user_id = ?", id)
	if result.Error != nil {
		return WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("failed to delete all payments for user %s: %s", id, result.Error)})
	}
	return WriteJSON(w, http.StatusNoContent, nil)
}

// GET/PULL/DELETE BY ID

func (s *APIServer) HandlePaymentByID(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		makeHTPPHandler(s.GetPaymentByID)(w, r)
	case http.MethodPost:
		WriteJSON(w, http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("invalid method %s", r.Method)})
	case http.MethodPut:
		makeHTPPHandler(s.UpdatePaymentByID)(w, r)
	case http.MethodDelete:
		makeHTPPHandler(s.DeletePaymentByID)(w, r)
	default:
		WriteJSON(w, http.StatusOK, map[string]string{"message": "/payments/{id}"})
	}
}

func (s *APIServer) GetPaymentByID(w http.ResponseWriter, r *http.Request) error {
	id := mux.Vars(r)["paymentID"]
	payment := models.Payment{}

	result := s.database.Where("ID = ?", id).First(&payment)
	if result.Error != nil {
		return WriteJSON(w, http.StatusNotFound, map[string]string{"error": fmt.Sprintf("failed to load payment: %s", result.Error)})
	}

	return WriteJSON(w, http.StatusOK, payment)
}

func (s *APIServer) UpdatePaymentByID(w http.ResponseWriter, r *http.Request) error {
	id := mux.Vars(r)["paymentID"]
	a := map[string]any{}
	if err := json.NewDecoder(r.Body).Decode(&a); err != nil {
		return WriteJSON(w, http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("failed to decode payment: %s", err)})
	}

	result := s.database.Model(&models.Payment{}).Where("id = ?", id).Updates(a)
	if result.Error != nil {
		return WriteJSON(w, http.StatusNotFound, map[string]string{"error": fmt.Sprintf("failed to update payment: %s", result.Error)})
	}

	return WriteJSON(w, http.StatusNoContent, nil)
}

func (s *APIServer) DeletePaymentByID(w http.ResponseWriter, r *http.Request) error {
	id := mux.Vars(r)["paymentID"]

	result := s.database.Delete(&models.Payment{}, id)
	if result.Error != nil {
		return WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("failed to delete all payment: %s", result.Error)})
	}
	return WriteJSON(w, http.StatusNoContent, nil)
}
