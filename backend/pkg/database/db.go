package database

//File che rappresenta il motore per comunicare con il db

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// InitDB inizializza e restituisce la connessione PostgreSQL usando GORM
func InitDB(dsn string) *gorm.DB {

	// Aprimao la connessione con GORM
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("Impossibile connettersi al database: %v", err)
	}

	log.Println("Connessione a PostgreSQL stabilita con successo tramite GORM!")

	return db
}
