package database

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// InitDB inizializza e restituisce la connessione PostgreSQL usando GORM
func InitDB(dsn string) *gorm.DB {

	// 👇 AGGIUNTO: PrepareStmt: false disabilita la cache delle query problematica su Neon DB
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		PrepareStmt: false,
	})

	if err != nil {
		log.Fatalf("Impossibile connettersi al database: %v", err)
	}

	log.Println("Connessione a PostgreSQL stabilita con successo tramite GORM!")

	return db
}
