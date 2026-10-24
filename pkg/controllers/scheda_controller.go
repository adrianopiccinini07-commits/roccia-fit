package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/adrianopiccinini/roccia-fit-api/pkg/models"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SchedaController struct {
	DB *gorm.DB
}

// CreaScheda godoc
// @Summary Crea scheda di allenamento
// @Description Genera una nuova scheda con giorni ed esercizi nidificati
// @Tags Schede
// @Accept json
// @Produce json
// @Param scheda body models.SchedaAllenamento true "Dati scheda"
// @Security BearerAuth
// @Success 201 {object} models.SchedaAllenamento
// @Router /schede [post]
func (c *SchedaController) CreaScheda(w http.ResponseWriter, r *http.Request) {
	var nuovaScheda models.SchedaAllenamento

	if err := json.NewDecoder(r.Body).Decode(&nuovaScheda); err != nil {
		// STAMPA L'ERRORE ESATTO NEL TERMINALE GO
		fmt.Println("❌ ERRORE DECODIFICA JSON:", err.Error())
		http.Error(w, "Errore nella lettura dei dati: "+err.Error(), http.StatusBadRequest)
		return
	}

	nuovaScheda.ID = uuid.New().String()

	// Assegnamo gli UUID, colleghiamo le chiavi esterne e verifichiamo l'esistenza dell'esercizio
	for i := range nuovaScheda.Giorni {
		nuovaScheda.Giorni[i].ID = uuid.New().String()
		nuovaScheda.Giorni[i].IDScheda = nuovaScheda.ID

		for j := range nuovaScheda.Giorni[i].Esercizi {
			if nuovaScheda.Giorni[i].Esercizi[j].IDEsercizio == "" {
				http.Error(w, "Errore: Nel giorno '"+nuovaScheda.Giorni[i].Giorno+"' c'è un esercizio senza ID.", http.StatusBadRequest)
				return
			}

			nuovaScheda.Giorni[i].Esercizi[j].ID = uuid.New().String()
			nuovaScheda.Giorni[i].Esercizi[j].IDGiorno = nuovaScheda.Giorni[i].ID

			idCercato := nuovaScheda.Giorni[i].Esercizi[j].IDEsercizio

			// CONTROLLO DIRETTO SULLA TABELLA "esercizi"
			var count int64
			c.DB.Model(&models.Esercizio{}).Where("id = ?", idCercato).Count(&count)
			if count == 0 {
				// Stampiamo nel terminale Go l'ID che sta fallendo la ricerca
				println("ESERCIZIO NON TROVATO NEL DB CON ID:", idCercato)
				http.Error(w, "Errore: L'esercizio con ID "+idCercato+" non esiste nella tabella esercizi del DB!", http.StatusBadRequest)
				return
			}
		}
	}

	if risultato := c.DB.Create(&nuovaScheda); risultato.Error != nil {
		http.Error(w, "Errore DB: "+risultato.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(nuovaScheda)
}

// GetSchedeAtleta godoc
// @Summary Elenco schede
// @Description Recupera tutte le schede di un atleta con i relativi esercizi
// @Tags Schede
// @Produce json
// @Param idAtleta path string true "ID Atleta"
// @Security BearerAuth
// @Success 200 {array} models.SchedaAllenamento
// @Router /schede/atleta/{idAtleta} [get]
func (c *SchedaController) GetSchedeAtleta(w http.ResponseWriter, r *http.Request) {
	idAtleta := chi.URLParam(r, "idAtleta")

	var schede []models.SchedaAllenamento

	// L'aggiunta di .Esercizio alla fine del Preload estrae i nomi (es. "Panca Piana")
	// Order("data_inizio desc") assicura che il frontend riceva per prima la scheda attuale
	risultato := c.DB.Preload("Giorni.Esercizi.Esercizio").
		Where("id_atleta = ?", idAtleta).
		Order("data_inizio desc").
		Find(&schede)

	if risultato.Error != nil {
		http.Error(w, "Errore durante il recupero delle schede", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(schede)
}

// GetSchedaByID godoc
// @Summary Dettaglio scheda
// @Description Recupera una specifica scheda tramite ID
// @Tags Schede
// @Produce json
// @Param id path string true "ID Scheda"
// @Security BearerAuth
// @Success 200 {object} models.SchedaAllenamento
// @Router /schede/{id} [get]
func (c *SchedaController) GetSchedaByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var scheda models.SchedaAllenamento

	// Aggiornato il Preload anche qui per coerenza
	if err := c.DB.Preload("Giorni.Esercizi.Esercizio").First(&scheda, "id = ?", id).Error; err != nil {
		http.Error(w, "Scheda non trovata", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(scheda)
}

// UpdateScheda godoc
// @Summary Aggiorna scheda
// @Description Sostituisce i dati di una scheda e delle sue relazioni nidificate
// @Tags Schede
// @Accept json
// @Produce json
// @Param id path string true "ID Scheda"
// @Param scheda body models.SchedaAllenamento true "Dati aggiornati"
// @Security BearerAuth
// @Success 200 {object} models.SchedaAllenamento
// @Router /schede/{id} [put]
func (c *SchedaController) UpdateScheda(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var scheda models.SchedaAllenamento
	if err := c.DB.First(&scheda, "id = ?", id).Error; err != nil {
		http.Error(w, "Dati non validi", http.StatusNotFound)
		return
	}

	var datiAggiornati models.SchedaAllenamento
	if err := json.NewDecoder(r.Body).Decode(&datiAggiornati); err != nil {
		http.Error(w, "Dati non validi", http.StatusBadRequest)
		return
	}

	// Aggiorniamo i dati della scheda
	c.DB.Model(&scheda).Updates(datiAggiornati)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(scheda)
}

// DeleteScheda godoc
// @Summary Elimina scheda
// @Description Rimuove una scheda e applica la cancellazione a cascata
// @Tags Schede
// @Produce json
// @Param id path string true "ID Scheda"
// @Security BearerAuth
// @Success 200 {string} string "Scheda eliminata"
// @Router /schede/{id} [delete]
func (c *SchedaController) DeleteScheda(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")

	var scheda models.SchedaAllenamento
	if err := c.DB.First(&scheda, "id = ?", id).Error; err != nil {
		http.Error(w, "Scheda non trovata", http.StatusNotFound)
		return
	}

	// Usiamo Select(clause.Associations) per dire a GORM di eliminare anche
	// tutti i Giorni e gli EserciziScheda collegati a questo ID
	if err := c.DB.Select(clause.Associations).Delete(&scheda).Error; err != nil {
		http.Error(w, "Errore durante l'eliminazione", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Scheda e relativi giorni eliminati con successo"}`))
}
