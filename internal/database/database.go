package database

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func Init() {
	var err error

	// Crée le dossier migrations s'il n'existe pas pour lire le fichier SQL
	migrationsDir := "migrations"
	if _, err := os.Stat(migrationsDir); os.IsNotExist(err) {
		log.Fatal("Le dossier migrations/ est introuvable. Assure-toi qu'il existe avec 001_init.sql")
	}

	// Ouvre la base (créée automatiquement)
	DB, err = sql.Open("sqlite", "./familiz.db")
	if err != nil {
		log.Fatal("Erreur d'ouverture de la base:", err)
	}

	// Lit le fichier de migration
	migrationPath := filepath.Join(migrationsDir, "001_init.sql")
	schemaBytes, err := os.ReadFile(migrationPath)
	if err != nil {
		log.Fatal("Impossible de lire le fichier de migration:", err)
	}

	// Exécute les instructions SQL
	_, err = DB.Exec(string(schemaBytes))
	if err != nil {
		log.Fatal("Erreur lors de l'exécution du schema:", err)
	}

	log.Println("✅ Base de données SQLite initialisée avec succès")
}
