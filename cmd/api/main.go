package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"

	"familiz/internal/auth"
	"familiz/internal/database"
	"familiz/internal/events"
	"familiz/internal/members"
	"familiz/internal/transactions"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  Aucun fichier .env trouvé")
	}

	database.Init()
	defer database.DB.Close()

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Routes publiques
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Bienvenue sur Familiz API V1",
			"status":  "running",
		})
	})

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		if err := database.DB.Ping(); err != nil {
			http.Error(w, "DB not reachable", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
	})

	// Routes d'authentification (publiques)
	r.Post("/register", auth.Register)
	r.Post("/login", auth.Login)

	// Routes protégées par JWT
	r.Group(func(r chi.Router) {
		r.Use(auth.Authenticate)

		// MEMBRES (CRUD complet)
		r.Post("/members", members.CreateMemberHandler)
		r.Get("/members", members.ListMembersHandler)
		r.Get("/members/{id}", members.GetMemberHandler)
		r.Put("/members/{id}", members.UpdateMemberHandler)
		r.Delete("/members/{id}", members.DeleteMemberHandler)

		// Transactions
		r.Post("/transactions", transactions.CreateTransaction)
		r.Get("/transactions", transactions.GetMemberTransactions)

		// Events
		r.Post("/events", events.CreateEvent)
		r.Get("/events", events.GetMemberEvents)
	})

	log.Println("🚀FAMILIZ dispo sur http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
