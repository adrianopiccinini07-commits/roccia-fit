package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/adrianopiccinini/roccia-fit-api/pkg/models"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PesoController struct {
	DB *gorm.DB
}

// AggiungiPeso godoc
// @Summary Registra pesata
// @Description Salva il peso corporeo giornaliero
// @Tags Peso
// @Accept json
// @Produce json
// @Param peso body models.PesoCorporeo true "Dati pesata"
// @Security BearerAuth
// @Success 201 {object} models.PesoCorporeo
// @Router /peso [post]
func (c *PesoController) AggiungiPeso(w http.ResponseWriter, r *http.Request) {
	var nuovaPesata models.PesoCorporeo

	if err := json.NewDecoder(r.Body).Decode(&nuovaPesata); err != nil {
		http.Error(w, "Errore nella lettura dei dati", http.StatusBadRequest)
		return
	}

	nuovaPesata.ID = uuid.New().String()
	if nuovaPesata.Data.IsZero() {
		nuovaPesata.Data = time.Now()
	}

	if risultato := c.DB.Create(&nuovaPesata); risultato.Error != nil {
		http.Error(w, "Errpre durante il salvataggio dei dati", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(nuovaPesata)
}

// GetStoricoPeso godoc
// @Summary Storico peso
// @Description Recupera tutte le pesate di un atleta in ordine cronologico
// @Tags Peso
// @Produce json
// @Param idAtleta path string true "ID Atleta"
// @Security BearerAuth
// @Success 200 {array} models.PesoCorporeo
// @Router /peso/atleta/{idAtleta} [get]
func (c *PesoController) GetStoricoPeso(w http.ResponseWriter, r *http.Request) {

	idAtleta := chi.URLParam(r, "idAtleta")

	var storico []models.PesoCorporeo

	if err := c.DB.Where("id_atleta = ?", idAtleta).Order("data desc").Find(&storico).Error; err != nil {
		http.Error(w, "Errore durante il recupero dello storico", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(storico)
}

// DeletePeso godoc
// @Summary Elimina pesata
// @Description Rimuove una singola pesata dal database
// @Tags Peso
// @Produce json
// @Param id path string true "ID Pesata"
// @Security BearerAuth
// @Success 200 {string} string "Pesata eliminata"
// @Router /peso/{id} [delete]
func (c *PesoController) DeletePeso(w http.ResponseWriter, r *http.Request) {
	idAtleta := chi.URLParam(r, "id")

	if err := c.DB.Delete(&models.PesoCorporeo{}, "id = ?", idAtleta).Error; err != nil {
		http.Error(w, "Errore durante l'eliminazione", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Pesata eliminata con successo"}`))
}

// UpdatePeso godoc
// @Summary Aggiorna pesata
// @Description Modifica i kg o la data di una pesata esistente
// @Tags Peso
// @Accept json
// @Produce json
// @Param id path string true "ID Pesata"
// @Param peso body models.PesoCorporeo true "Dati aggiornati"
// @Security BearerAuth
// @Success 200 {object} models.PesoCorporeo
// @Router /peso/{id} [put]
func (c *PesoController) UpdatePeso(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var pesata models.PesoCorporeo
	if err := c.DB.First(&pesata, "id = ?", id).Error; err != nil {
		http.Error(w, "Pesata non trovata", http.StatusNotFound)
		return
	}

	var datiAggiornati models.PesoCorporeo
	if err := json.NewDecoder(r.Body).Decode(&datiAggiornati); err != nil {
		http.Error(w, "Dati non validi", http.StatusBadRequest)
		return
	}

	c.DB.Model(&pesata).Updates(datiAggiornati)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(pesata)
}

// GetUltimoPeso godoc
// @Summary Ultimo peso
// @Description Recupera solo la pesata più recente
// @Tags Peso
// @Produce json
// @Param idAtleta path string true "ID Atleta"
// @Security BearerAuth
// @Success 200 {object} models.PesoCorporeo
// @Router /peso/ultimo/{idAtleta} [get]
func (c *PesoController) GetUltimoPeso(w http.ResponseWriter, r *http.Request) {
	idAtleta := chi.URLParam(r, "idAtleta")

	var ultimaPesata models.PesoCorporeo
	if err := c.DB.Where("id_atleta = ?", idAtleta).Order("data desc").First(&ultimaPesata).Error; err != nil {
		http.Error(w, "Nessuna pesata trovata", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ultimaPesata)
}

// AggiungiPesiMultipli godoc
// @Summary Inserimento massivo pesi
// @Description Carica un array di pesate storiche in un colpo solo
// @Tags Peso
// @Accept json
// @Produce json
// @Param pesate body []models.PesoCorporeo true "Array di pesate"
// @Security BearerAuth
// @Success 201 {array} models.PesoCorporeo
// @Router /peso/bulk [post]
func (c *PesoController) AggiungiPesiMultipli(w http.ResponseWriter, r *http.Request) {
	var nuovePesate []models.PesoCorporeo
	if err := json.NewDecoder(r.Body).Decode(&nuovePesate); err != nil {
		http.Error(w, "Errore nella decodifica del JSON", http.StatusBadRequest)
		return
	}

	for i := range nuovePesate {
		nuovePesate[i].ID = uuid.New().String()
		if nuovePesate[i].Data.IsZero() {
			nuovePesate[i].Data = time.Now()
		}
	}

	if risultato := c.DB.Create(&nuovePesate); risultato.Error != nil {
		http.Error(w, "Errore durante il salvataggio dei dati multiplo", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(nuovePesate)
}

// GetPesoByID godoc
// @Summary Singola pesata
// @Description Recupera una pesata specifica tramite ID
// @Tags Peso
// @Produce json
// @Param id path string true "ID Pesata"
// @Security BearerAuth
// @Success 200 {object} models.PesoCorporeo
// @Router /peso/{id} [get]
func (c *PesoController) GetPesoByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var pesata models.PesoCorporeo
	if err := c.DB.First(&pesata, "id = ?", id).Error; err != nil {
		http.Error(w, "Peso non trovato", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(pesata)
}
