package handler

import (
	"balance/internal/logger"
	"balance/internal/model"
	"encoding/json"
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
		h.logger.Warning("Method not allowed")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		UserID int     `json:"user_id"`
		Amount float64 `json:"amount"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("Invalid request body")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.usecase.Introduction(req.UserID, req.Amount); err != nil {
		h.logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.logger.Info("Introduction completed successfully")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Introduction completed successfully"))
}

func (h *Handler) Debit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.logger.Warning("Method not allowed")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		UserID int     `json:"user_id"`
		Amount float64 `json:"amount"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("Invalid request body")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.usecase.Debit(req.UserID, req.Amount); err != nil {
		h.logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.logger.Info("Debit completed successfully")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Debit completed successfully"))
}

func (h *Handler) Transfer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.logger.Warning("Method not allowed")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		FromUserID int     `json:"from_user_id"`
		ToUserID   int     `json:"to_user_id"`
		Amount     float64 `json:"amount"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("Invalid request body")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.usecase.Transfer(req.FromUserID, req.ToUserID, req.Amount); err != nil {
		h.logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.logger.Info("Transfer completed successfully")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Transfer completed successfully"))
}

func (h *Handler) GetBalance(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["user_id"]
	userID, err := strconv.Atoi(idStr)
	if err != nil {
		h.logger.Error("Invalid user ID")
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	balance, err := h.usecase.GetBalance(userID)
	if err != nil {
		if err.Error() == "user not found" {
			h.logger.Warning("User not found")
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			h.logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	h.logger.Info("Get balance completed successfully")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(balance)
}
