package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/adrianopiccinini/roccia-fit-api/pkg/models"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EsercizioController struct {
	DB *gorm.DB
}

// CreaEsercizio godoc
// @Summary Crea esercizio
// @Description Aggiunge un nuovo esercizio al catalogo globale
// @Tags Esercizi
// @Accept json
// @Produce json
// @Param esercizio body models.Esercizio true "Dati esercizio"
// @Security BearerAuth
// @Success 201 {object} models.Esercizio
// @Router /esercizi [post]
func (c *EsercizioController) CreaEsercizio(w http.ResponseWriter, r *http.Request) {
	var nuovoEsercizio models.Esercizio

	if err := json.NewDecoder(r.Body).Decode(&nuovoEsercizio); err != nil {
		http.Error(w, "Errore nella decodifica del JSON", http.StatusBadRequest)
		return
	}

	nuovoEsercizio.ID = uuid.New().String()
	if risultato := c.DB.Create(&nuovoEsercizio); risultato.Error != nil {
		http.Error(w, "Errore durante il salvataggio", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(nuovoEsercizio)
}

// CreaEserciziMultipli godoc
// @Summary Creazione massiva esercizi
// @Description Inserisce un array di esercizi in un colpo solo
// @Tags Esercizi
// @Accept json
// @Produce json
// @Param esercizi body []models.Esercizio true "Array di esercizi"
// @Security BearerAuth
// @Success 201 {array} models.Esercizio
// @Router /esercizi/bulk [post]
func (c *EsercizioController) CreaEserciziMultipli(w http.ResponseWriter, r *http.Request) {
	var nuoviEsercizi []models.Esercizio

	if err := json.NewDecoder(r.Body).Decode(&nuoviEsercizi); err != nil {
		http.Error(w, "Errore nella decodifica del JSon", http.StatusBadRequest)
		return
	}

	for i := range nuoviEsercizi {
		nuoviEsercizi[i].ID = uuid.New().String()
	}

	if risultato := c.DB.Create(&nuoviEsercizi); risultato.Error != nil {
		http.Error(w, "Errore durante il salvataggio nl database", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(nuoviEsercizi)
}

// GetEsercizi godoc
// @Summary Catalogo esercizi
// @Description Recupera la lista completa degli esercizi disponibili
// @Tags Esercizi
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.Esercizio
// @Router /esercizi [get]
func (c *EsercizioController) GetEsercizi(w http.ResponseWriter, r *http.Request) {
	var esercizi []models.Esercizio

	if err := c.DB.Find(&esercizi).Error; err != nil {
		http.Error(w, "Errore durante il recupero degli esercizi", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(esercizi)
}

// GetEsercizioByID godoc
// @Summary Singolo esercizio
// @Description Recupera i dettagli di uno specifico esercizio tramite ID
// @Tags Esercizi
// @Produce json
// @Param id path string true "ID Esercizio"
// @Security BearerAuth
// @Success 200 {object} models.Esercizio
// @Router /esercizi/{id} [get]
func (c *EsercizioController) GetEsercizioByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var esercizio models.Esercizio
	if err := c.DB.First(&esercizio, "id = ?", id).Error; err != nil {
		http.Error(w, "Esercizio non trovato", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(esercizio)
}

// UpdateEsercizio godoc
// @Summary Aggiorna esercizio
// @Description Modifica i dati (es. nome o gruppo muscolare) di un esercizio
// @Tags Esercizi
// @Accept json
// @Produce json
// @Param id path string true "ID Esercizio"
// @Param esercizio body models.Esercizio true "Dati da aggiornare"
// @Security BearerAuth
// @Success 200 {object} models.Esercizio
// @Router /esercizi/{id} [put]
func (c *EsercizioController) UpdateEsercizio(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var esercizio models.Esercizio
	if err := c.DB.First(&esercizio, "id = ?", id).Error; err != nil {
		http.Error(w, "Esercizio non trovato", http.StatusNotFound)
		return
	}

	var datiAggiornati models.Esercizio
	if err := json.NewDecoder(r.Body).Decode(&datiAggiornati); err != nil {
		http.Error(w, "Dati non validi", http.StatusInternalServerError)
		return
	}

	c.DB.Model(&esercizio).Updates(datiAggiornati)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(esercizio)
}

// DeleteEsercizio godoc
// @Summary Elimina esercizio
// @Description Rimuove logicamente (Soft Delete) un esercizio dal catalogo
// @Tags Esercizi
// @Produce json
// @Param id path string true "ID Esercizio"
// @Security BearerAuth
// @Success 200 {string} string "Esercizio eliminato"
// @Router /esercizi/{id} [delete]
func (c *EsercizioController) DeleteEsercizio(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := c.DB.Delete(&models.Esercizio{}, "id = ?", id).Error; err != nil {
		http.Error(w, "Errore durante l'eliminazione", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Esercizio eliminato con successo"}`))
}
