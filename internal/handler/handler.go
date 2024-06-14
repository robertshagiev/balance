package handler

import (
	"balance/internal/logger"
	"balance/internal/model"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Handler struct {
	usecase ucase
	logger  *logger.Logger
}

type ucase interface {
	Introduction(userID int, amount float64) error
	Debit(userID int, amount float64) error
	Transfer(fromUserID, toUserID int, amount float64) error
	GetBalance(userID int) (model.Balance, error)
}

func NewHandler(u ucase, log *logger.Logger) *Handler {
	return &Handler{usecase: u, logger: log}
}

func (h *Handler) Introduction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.logger.Warning("Invalid method for Introduction")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		UserID int     `json:"user_id"`
		Amount float64 `json:"amount"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("Invalid request body for Introduction")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.usecase.Introduction(req.UserID, req.Amount); err != nil {
		h.logger.Error(fmt.Sprintf("Introduction failed for user %d: %v", req.UserID, err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.logger.Info(fmt.Sprintf("Introduction completed for user %d with amount %f", req.UserID, req.Amount))
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Introduction completed successfully"))
}

func (h *Handler) Debit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.logger.Warning("Invalid method for Debit")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		UserID int     `json:"user_id"`
		Amount float64 `json:"amount"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("Invalid request body for Debit")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.usecase.Debit(req.UserID, req.Amount); err != nil {
		h.logger.Error(fmt.Sprintf("Debit failed for user %d: %v", req.UserID, err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.logger.Info(fmt.Sprintf("Debit completed for user %d with amount %f", req.UserID, req.Amount))
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Debit completed successfully"))
}

func (h *Handler) Transfer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.logger.Warning("Invalid method for Transfer")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		FromUserID int     `json:"from_user_id"`
		ToUserID   int     `json:"to_user_id"`
		Amount     float64 `json:"amount"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("Invalid request body for Transfer")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.usecase.Transfer(req.FromUserID, req.ToUserID, req.Amount); err != nil {
		h.logger.Error(fmt.Sprintf("Transfer failed from user %d to user %d: %v", req.FromUserID, req.ToUserID, err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.logger.Info(fmt.Sprintf("Transfer completed from user %d to user %d with amount %f", req.FromUserID, req.ToUserID, req.Amount))
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Transfer completed successfully"))
}

func (h *Handler) GetBalance(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["user_id"]
	userID, err := strconv.Atoi(idStr)
	if err != nil {
		h.logger.Error("Invalid user ID for GetBalance")
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	balance, err := h.usecase.GetBalance(userID)
	if err != nil {
		if err.Error() == "user not found" {
			h.logger.Warning(fmt.Sprintf("User not found for GetBalance: %d", userID))
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			h.logger.Error(fmt.Sprintf("GetBalance failed for user %d: %v", userID, err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	h.logger.Info(fmt.Sprintf("GetBalance completed for user %d", userID))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(balance)
}
