package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/adrianopiccinini/roccia-fit-api/pkg/models"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

type StatsController struct {
	DB *gorm.DB
}

// GetTonnellaggio godoc
// @Summary Calcolo tonnellaggio
// @Description Calcola il volume totale sollevato (Kg * Reps) in una sessione
// @Tags Statistiche
// @Produce json
// @Param idAllenamento path string true "ID Allenamento"
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /stats/tonnellaggio/{idAllenamento} [get]
func (c *StatsController) GetTonnellaggio(w http.ResponseWriter, r *http.Request) {
	idAllenamento := chi.URLParam(r, "idAllenamento")

	var serie []models.SerieEseguita
	if err := c.DB.Where("allenamento_id = ?", idAllenamento).Find(&serie).Error; err != nil {
		http.Error(w, "Errore nel recupero delle serie", http.StatusInternalServerError)
		return
	}

	var tonnellaggioTot float32
	for _, s := range serie {
		tonnellaggioTot += float32(s.Peso) * float32(s.Reps)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"allenamento_id":  idAllenamento,
		"tonnellaggio_kg": tonnellaggioTot,
	})
}

// GetProiezione1RM godoc
// @Summary Stima Massimale
// @Description Calcola l'1RM teorico di una serie tramite l'equazione di Brzycki
// @Tags Statistiche
// @Produce json
// @Param idSerie path string true "ID Serie"
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /stats/massimale/{idSerie} [get]
func (c *StatsController) GetProiezione1RM(w http.ResponseWriter, r *http.Request) {
	idSerie := chi.URLParam(r, "idSerie")

	var serie models.SerieEseguita
	if err := c.DB.First(&serie, "id = ?", idSerie).Error; err != nil {
		http.Error(w, "Serie non trovata", http.StatusNotFound)
		return
	}

	var proiezione float32
	if serie.Reps == 1 {
		proiezione = float32(serie.Peso)
	} else {
		proiezione = float32(serie.Peso) / (1.0278 - (0.0278 * float32(serie.Reps)))
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"serie_id":       idSerie,
		"esercizio_id":   serie.IDEsercizioScheda,
		"proiezione_1rm": proiezione,
	})
}

// GetVolumeSettimanale godoc
// @Summary Analisi volume settimanale
// @Description Calcola il totale delle serie allenanti per gruppo muscolare degli ultimi 7 giorni
// @Tags Statistiche
// @Produce json
// @Param idAtleta path string true "ID Atleta"
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /stats/volume/{idAtleta} [get]
func (c *StatsController) GetVolumeSettimanale(w http.ResponseWriter, r *http.Request) {
	idAtleta := chi.URLParam(r, "idAtleta")

	type RisultatoVolume struct {
		GruppoMuscolare string `json:"gruppo_muscolare"`
		TotaleSerie     int    `json:"totale_serie"`
	}

	var volumi []RisultatoVolume

	// Utilizzo della query corretta con le associazioni JOIN precise
	err := c.DB.Table("serie_eseguite").
		Select("esercizios.gruppo_muscolare, count(serie_eseguite.id) as totale_serie").
		Joins("JOIN allenamentos ON allenamentos.id = serie_eseguite.allenamento_id").
		Joins("JOIN esercizios ON esercizios.id = serie_eseguite.esercizio_id").
		Where("allenamentos.id_atleta = ? AND allenamentos.data >= current_date - interval '7 days'", idAtleta).
		Group("esercizios.gruppo_muscolare").
		Scan(&volumi).Error

	if err != nil {
		http.Error(w, "Errore nel calcolo dei volumi", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"atleta_id":       idAtleta,
		"volumi_7_giorni": volumi,
	})
}
