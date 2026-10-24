package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/adrianopiccinini/roccia-fit-api/pkg/models"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

type CardioController struct {
	DB *gorm.DB
}

// Salva o sovrascrive il protocollo Cardio dell'atleta
func (c *CardioController) SalvaCardio(w http.ResponseWriter, r *http.Request) {
	var input models.ProtocolloCardio
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "JSON non valido", http.StatusBadRequest)
		return
	}

	// Usiamo una transazione: eliminiamo il vecchio protocollo (se esiste) e salviamo il nuovo
	err := c.DB.Transaction(func(tx *gorm.DB) error {
		var existing models.ProtocolloCardio
		if err := tx.Where("id_atleta = ?", input.IDAtleta).First(&existing).Error; err == nil {
			// Se esiste un protocollo precedente, lo cancelliamo (le sessioni verranno rimosse a cascata)
			tx.Delete(&existing)
		}

		// Creiamo il nuovo protocollo da zero
		if err := tx.Create(&input).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		http.Error(w, "Errore salvataggio cardio: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(input)
}

// Recupera il protocollo Cardio di un atleta specifico
func (c *CardioController) GetCardioAtleta(w http.ResponseWriter, r *http.Request) {
	idAtleta := chi.URLParam(r, "idAtleta")
	var protocollo models.ProtocolloCardio

	// Preload per includere l'array delle sessioni
	result := c.DB.Preload("Sessioni").Where("id_atleta = ?", idAtleta).First(&protocollo)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			http.Error(w, "Nessun protocollo cardio trovato", http.StatusNotFound)
			return
		}
		http.Error(w, "Errore database", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(protocollo)
}
