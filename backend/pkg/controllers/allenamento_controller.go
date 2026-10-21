package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/adrianopiccinini/roccia-fit-api/pkg/models"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AllenamentoController struct {
	DB *gorm.DB
}

// GetAllenamentiAtleta godoc
// @Summary Storico allenamenti
// @Description Recupera gli allenamenti in formato paginato
// @Tags Allenamenti
// @Produce json
// @Param idAtleta path string true "ID Atleta"
// @Param page query int false "Numero pagina (default: 1)"
// @Param limit query int false "Elementi per pagina (default: 10)"
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /allenamenti/atleta/{idAtleta} [get]
func (c *AllenamentoController) GetAllenamentiAtleta(w http.ResponseWriter, r *http.Request) {
	// GET /atleti/123/allenamenti?page=2&limit=10
	// 123 e' un PATH PARAMETER
	// mentre ?page=2&limit=10 sono dei QUERY PARAMETERS
	idAtleta := chi.URLParam(r, "idAtleta")

	// Lettura dei query parameters con valori di default di sicurezza, prendono delle stringhe non dei numeri
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page := 1
	limit := 10

	// Li convertiamo in numeri perche inizalmente erano stringhe Atoi (ASCII to Integer)
	// p -> numero convertito 	err -> eventuale errore
	if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
		page = p
	}

	if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
		limit = l
	}

	// offset indica da quale record iniziare la lettura
	offset := (page - 1) * limit

	var allenamenti []models.Allenamento
	var totaleRecord int64

	// Conteggio totale dei record , serve per dire al front end che queste pagine esistono
	//? -> permette a GORM/database driver di gestire correttamente il valore e aiuta a evitare SQL injection.
	c.DB.Model(&models.Allenamento{}).Where("idAtleta = ?", idAtleta).Count(&totaleRecord)

	risultato := c.DB.Preload("SerieEseguite").
		Where("id_atleta = ?", idAtleta).
		Order("data desc").
		Limit(limit).
		Offset(offset).
		Find(&allenamenti)

	if risultato.Error != nil {
		http.Error(w, "Errore durante il recupero degli allenamenti", http.StatusInternalServerError)
		return
	}

	// Creazione dell'envelope di risposta , l'envelope restituisce i dati all'interno di una struttura invece di restituire direttamente i dati
	risposta := map[string]interface{}{
		"dati":          allenamenti,
		"pagina":        page,
		"limite":        limit,
		"totale":        totaleRecord,
		"totale_pagine": (int(totaleRecord) + limit - 1) / limit,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(risposta)

	/*
		// Ordiniamo per data decrescente così hai sempre l'ultimo allenamento in cima
		risultato := c.DB.Preload("SerieEseguite").Where("id_atleta = ?", idAtleta).Order("data desc").Find(&allenamenti)
		if risultato.Error != nil {
			http.Error(w, "Errore durante il recupero degli allenamenti", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(allenamenti)
	*/
}

// IniziaAllenamento godoc
// @Summary Inizia un nuovo allenamento
// @Description Crea una nuova sessione vuota
// @Tags Allenamenti
// @Accept json
// @Produce json
// @Param allenamento body models.Allenamento true "Dati allenamento"
// @Security BearerAuth
// @Success 201 {object} models.Allenamento
// @Router /allenamenti/inizia [post]
func (c *AllenamentoController) IniziaAllenamento(w http.ResponseWriter, r *http.Request) {

	var nuovoAllenamento models.Allenamento
	if err := json.NewDecoder(r.Body).Decode(&nuovoAllenamento); err != nil {
		http.Error(w, "Errore nella decodifica", http.StatusBadRequest)
		return
	}

	nuovoAllenamento.ID = uuid.New().String()

	if nuovoAllenamento.Data.IsZero() {
		nuovoAllenamento.Data = time.Now()
	}

	if risultato := c.DB.Create(&nuovoAllenamento); risultato.Error != nil {
		http.Error(w, "Errore durante la creazione dell'allenamento", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(nuovoAllenamento)
}

// AggiungiSerie godoc
// @Summary Registra una serie
// @Description Aggiunge una serie eseguita al database
// @Tags Allenamenti
// @Accept json
// @Produce json
// @Param serie body models.SerieEseguita true "Dati serie"
// @Security BearerAuth
// @Success 201 {object} models.SerieEseguita
// @Router /allenamenti/serie [post]
func (c *AllenamentoController) AggiungiSerie(w http.ResponseWriter, r *http.Request) {
	var nuovaSerie models.SerieEseguita
	if err := json.NewDecoder(r.Body).Decode(&nuovaSerie); err != nil {
		http.Error(w, "Errore nella decodifica", http.StatusBadRequest)
		return
	}

	nuovaSerie.ID = uuid.New().String()

	if risultato := c.DB.Create(&nuovaSerie); risultato.Error != nil {
		http.Error(w, "Errore nel salvataggio della serie", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(nuovaSerie)
}

// GetAllenamentoById godoc
// @Summary Dettagli allenamento
// @Description Recupera un allenamento specifico con tutte le sue serie
// @Tags Allenamenti
// @Produce json
// @Param id path string true "ID Allenamento"
// @Security BearerAuth
// @Success 200 {object} models.Allenamento
// @Router /allenamenti/{id} [get]
func (c *AllenamentoController) GetAllenamentoById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var allenamento models.Allenamento

	if err := c.DB.Preload("SerieEseguite").First(&allenamento, "id = ?", id).Error; err != nil {
		http.Error(w, "Allenamento non trovato", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(allenamento)
}

// DeleteAllenamento godoc
// @Summary Elimina allenamento
// @Description Rimuove un allenamento e le serie associate
// @Tags Allenamenti
// @Produce json
// @Param id path string true "ID Allenamento"
// @Security BearerAuth
// @Success 200 {string} string "Allenamento eliminato"
// @Router /allenamenti/{id} [delete]
func (c *AllenamentoController) DeleteAllenamento(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var allenamento models.Allenamento
	if err := c.DB.First(&allenamento, "id = ?", id).Error; err != nil {
		http.Error(w, "Allenamento non trovato", http.StatusNotFound)
		return
	}

	if err := c.DB.Select(clause.Associations).Delete(&allenamento).Error; err != nil {
		http.Error(w, "Errore durante l'eleminazione", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Allenamento eliminato con successo"}`))
}

// UpdateSerie godoc
// @Summary Aggiorna serie
// @Description Modifica i kg o le ripetizioni di una serie esistente
// @Tags Allenamenti
// @Accept json
// @Produce json
// @Param id path string true "ID Serie"
// @Param serie body models.SerieEseguita true "Dati da aggiornare"
// @Security BearerAuth
// @Success 200 {object} models.SerieEseguita
// @Router /allenamenti/serie/{id} [put]
func (c *AllenamentoController) UpdateSerie(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var serie models.SerieEseguita
	if err := c.DB.First(&serie, "id = ?", id).Error; err != nil {
		http.Error(w, "Serie non trovata", http.StatusNotFound)
		return
	}

	var datiAggiornati models.SerieEseguita
	if err := json.NewDecoder(r.Body).Decode(&datiAggiornati); err != nil {
		http.Error(w, "Dati non validi", http.StatusBadRequest)
		return
	}

	c.DB.Model(&serie).Updates(datiAggiornati)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(serie)
}

// DeleteSerie godoc
// @Summary Elimina serie
// @Description Rimuove una singola serie
// @Tags Allenamenti
// @Produce json
// @Param id path string true "ID Serie"
// @Security BearerAuth
// @Success 200 {string} string "Serie eliminata"
// @Router /allenamenti/serie/{id} [delete]
func (c *AllenamentoController) DeleteSerie(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := c.DB.Delete(&models.SerieEseguita{}, "id = ?", id).Error; err != nil {
		http.Error(w, "Serie non trovata", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Serie eliminata con successo"}`))
}
