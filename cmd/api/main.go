package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"

	"familiz/internal/apps/auth"
	"familiz/internal/apps/events"
	"familiz/internal/apps/members"
	"familiz/internal/apps/profile"
	"familiz/internal/apps/settings"
	"familiz/internal/apps/transactions"
	"familiz/internal/database"
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
		r.Put("/members/{id}", members.UpdateMemberHandler)
		r.Delete("/members/{id}", members.DeleteMemberHandler)

		// PROFILE
		r.Get("/profile/{id}", profile.GetMemberProfileHandler)

		// Transactions
		r.Post("/transactions", transactions.CreateTransactionHandler)
		r.Get("/transactions", transactions.ListTransactionsHandler) // Gère ?member_id= X
		r.Put("/transactions/{id}", transactions.UpdateTransactionHandler)
		r.Delete("/transactions/{id}", transactions.DeleteTransactionHandler)

		// Events
		r.Post("/events", events.CreateEventHandler)
		r.Get("/events", events.ListEventsHandler) // ?member_id=1
		r.Put("/events/{id}", events.UpdateEventHandler)
		r.Delete("/events/{id}", events.DeleteEventHandler)

		// SETTINGS
		r.Get("/settings/contributions", settings.GetContributionSettingsHandler)
		r.Put("/settings/contributions", settings.UpdateContributionSettingsHandler)
		r.Get("/settings/events", settings.GetEventSettingsHandler)
		r.Put("/settings/events/{type}", settings.UpdateEventSettingHandler)
	})

	log.Println("🚀FAMILIZ dispo sur http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
