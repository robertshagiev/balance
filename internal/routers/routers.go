package routers

import (
	"balance/internal/handler"

	"github.com/gorilla/mux"
)

func NewRouter(h *handler.Handler) *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/balance/introduction", h.Introduction).Methods("POST")
	r.HandleFunc("/balance/debit", h.Debit).Methods("POST")
	r.HandleFunc("/balance/transfer", h.Transfer).Methods("POST")
	r.HandleFunc("/balance/get/{user_id}", h.GetBalance).Methods("GET")

	return r
}
