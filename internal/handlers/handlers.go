package handlers

import (
	"auth-microservice/internal/crypto"
	"auth-microservice/internal/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Handlers struct {
	Mux  *chi.Mux
	Stor storage
}

type storage interface {
	UserByLogin(login string) *models.User
	AddUser(user models.User)
}

func (h *Handlers) HandleLogin(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Println(err)
	}
	user2 := storage.UserByLogin(h.Stor, user.Login)
	crypto.CheckPwd(user.HashedPwd, user2.HashedPwd)
}

func (h *Handlers) HandleRegister(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error")
		log.Println(err)
	}
	storage.AddUser(h.Stor, user)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User created"))
}
func (h *Handlers) SetupHandlers() {
	h.Mux.Post("/login", h.HandleLogin)
	h.Mux.Post("/register", h.HandleRegister)
}
