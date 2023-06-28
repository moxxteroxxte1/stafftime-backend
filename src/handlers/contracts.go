package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/moxxteroxxte1/stafftime-backend/src/models"
)

// GET/POST/PUT/DELETE ALL

func (s *APIServer) handleContracts(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("token")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
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

	if claims, ok := token.Claims.(jwt.MapClaims); !ok || !claims["IsAdmin"].(bool) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	switch r.Method {
	case http.MethodGet:
		makeHTPPHandler(s.GetAllContracts)(w, r)
	case http.MethodPost:
		makeHTPPHandler(s.CreateContract)(w, r)
	case http.MethodPut:
		makeHTPPHandler(s.UpdateAllContracts)(w, r)
	case http.MethodDelete:
		makeHTPPHandler(s.DeleteAllContracts)(w, r)
	default:
		WriteJSON(w, http.StatusOK, map[string]string{"message": "/contracts"})
	}
}

func (s *APIServer) GetAllContracts(w http.ResponseWriter, r *http.Request) error {
	contracts := []models.Contract{}
	s.database.Find(&contracts)
	return WriteJSON(w, http.StatusOK, contracts)
}

func (s *APIServer) CreateContract(w http.ResponseWriter, r *http.Request) error {
	contract := new(models.Contract)
	if err := json.NewDecoder(r.Body).Decode(contract); err != nil {
		return WriteJSON(w, http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("failed to decode contract: %s", err)})
	}

	result := s.database.Create(&contract)
	if result.Error != nil {
		return WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("failed to create contract: %s", result.Error)})
	}

	return WriteJSON(w, http.StatusCreated, contract)
}

func (s *APIServer) UpdateAllContracts(w http.ResponseWriter, r *http.Request) error {
	a := map[string]any{}
	if err := json.NewDecoder(r.Body).Decode(&a); err != nil {
		return WriteJSON(w, http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("failed to decode contract: %s", err)})
	}

	result := s.database.Model(&models.Contract{}).Where("true").Updates(a)
	if result.Error != nil {
		return WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("failed to update contract: %s", result.Error)})
	}

	return WriteJSON(w, http.StatusNoContent, nil)
}

func (s *APIServer) DeleteAllContracts(w http.ResponseWriter, r *http.Request) error {
	contracts := []models.Contract{}

	result := s.database.Find(&contracts).Delete(&contracts)
	if result.Error != nil {
		return WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("failed to delete all contracts: %s", result.Error)})
	}
	return WriteJSON(w, http.StatusNoContent, nil)
}

// GET/PUT/DELETE BY USER

func (s *APIServer) HandleContractByUserID(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("token")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
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

	id, convErr := strconv.ParseUint(mux.Vars(r)["userID"], 10, 32)
	if convErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if claims, ok := token.Claims.(jwt.MapClaims); !ok || (!claims["IsAdmin"].(bool) && uint(claims["UserID"].(float64)) != uint(id)) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	switch r.Method {
	case http.MethodGet:
		makeHTPPHandler(s.GetContractsByUserID)(w, r)
	case http.MethodPost:
		makeHTPPHandler(s.CreateContractByUserID)(w, r)
	case http.MethodPut:
		makeHTPPHandler(s.UpdateContractByUserID)(w, r)
	case http.MethodDelete:
		makeHTPPHandler(s.DeleteContractByUserID)(w, r)
	default:
		WriteJSON(w, http.StatusOK, map[string]string{"message": "/user/{id}/contracts"})
	}
}

func (s *APIServer) GetContractsByUserID(w http.ResponseWriter, r *http.Request) error {
	id := mux.Vars(r)["userID"]
	contracts := []models.Contract{}

	result := s.database.Where("user_id = ?", id).Find(&contracts)
	if result.Error != nil {
		return WriteJSON(w, http.StatusNotFound, map[string]string{"error": fmt.Sprintf("failed to load contract by userId: %s", result.Error)})
	}

	return WriteJSON(w, http.StatusOK, contracts)
}

func (s *APIServer) CreateContractByUserID(w http.ResponseWriter, r *http.Request) error {
	contract := new(models.Contract)

	i := mux.Vars(r)["userID"]
	userID, err := strconv.Atoi(i)
	if err != nil {
		return WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("failed to create contract: %s", err)})
	}

	data := []byte(fmt.Sprintf(`{"userID": %d}`, uint(userID)))
	jsonErr := json.Unmarshal(data, &contract)
	if jsonErr != nil {
		return WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("failed to create contract: %s", jsonErr)})
	}

	if err1 := json.NewDecoder(r.Body).Decode(contract); err1 != nil {
		return WriteJSON(w, http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("failed to decode contract: %s", err1)})
	}

	result := s.database.Create(&contract)
	if result.Error != nil {
		return WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("failed to create contract: %s", result.Error)})
	}

	return WriteJSON(w, http.StatusCreated, contract)
}

func (s *APIServer) UpdateContractByUserID(w http.ResponseWriter, r *http.Request) error {
	id := mux.Vars(r)["userID"]
	a := map[string]any{}
	if err := json.NewDecoder(r.Body).Decode(&a); err != nil {
		return WriteJSON(w, http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("failed to decode contract: %s", err)})
	}

	result := s.database.Model(&models.Contract{}).Where("user_id = ?", id).Updates(a)
	if result.Error != nil {
		return WriteJSON(w, http.StatusNotFound, map[string]string{"error": fmt.Sprintf("failed to update contract: %s", result.Error)})
	}

	return WriteJSON(w, http.StatusNoContent, nil)
}

func (s *APIServer) DeleteContractByUserID(w http.ResponseWriter, r *http.Request) error {
	id := mux.Vars(r)["userID"]

	result := s.database.Delete(&models.Contract{}, "user_id = ?", id)
	if result.Error != nil {
		return WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("failed to delete all contracts for user %s: %s", id, result.Error)})
	}
	return WriteJSON(w, http.StatusNoContent, nil)
}

// GET/PULL/DELETE BY ID

func (s *APIServer) HandleContractByID(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("token")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
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

	if claims, ok := token.Claims.(jwt.MapClaims); !ok || !claims["IsAdmin"].(bool) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	switch r.Method {
	case http.MethodGet:
		makeHTPPHandler(s.GetContractByID)(w, r)
	case http.MethodPost:
		WriteJSON(w, http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("invalid method %s", r.Method)})
	case http.MethodPut:
		makeHTPPHandler(s.UpdateContractByID)(w, r)
	case http.MethodDelete:
		makeHTPPHandler(s.DeleteContractByID)(w, r)
	default:
		WriteJSON(w, http.StatusOK, map[string]string{"message": "/contracts/{id}"})
	}
}

func (s *APIServer) GetContractByID(w http.ResponseWriter, r *http.Request) error {
	id := mux.Vars(r)["contractID"]
	contract := models.Contract{}

	result := s.database.Where("ID = ?", id).First(&contract)
	if result.Error != nil {
		return WriteJSON(w, http.StatusNotFound, map[string]string{"error": fmt.Sprintf("failed to load contract: %s", result.Error)})
	}

	return WriteJSON(w, http.StatusOK, contract)
}

func (s *APIServer) UpdateContractByID(w http.ResponseWriter, r *http.Request) error {
	id := mux.Vars(r)["contractID"]
	a := map[string]any{}
	if err := json.NewDecoder(r.Body).Decode(&a); err != nil {
		return WriteJSON(w, http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("failed to decode contract: %s", err)})
	}

	result := s.database.Model(&models.Contract{}).Where("id = ?", id).Updates(a)
	if result.Error != nil {
		return WriteJSON(w, http.StatusNotFound, map[string]string{"error": fmt.Sprintf("failed to update contract: %s", result.Error)})
	}

	return WriteJSON(w, http.StatusNoContent, nil)
}

func (s *APIServer) DeleteContractByID(w http.ResponseWriter, r *http.Request) error {
	id := mux.Vars(r)["contractID"]

	result := s.database.Delete(&models.Contract{}, id)
	if result.Error != nil {
		return WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("failed to delete all contract: %s", result.Error)})
	}
	return WriteJSON(w, http.StatusNoContent, nil)
}
