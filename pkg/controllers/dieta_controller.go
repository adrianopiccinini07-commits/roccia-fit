package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/adrianopiccinini/roccia-fit-api/pkg/models"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

type DietaController struct {
	DB *gorm.DB
}

func (c *DietaController) CreaDieta(w http.ResponseWriter, r *http.Request) {
	var dieta models.Dieta
	if err := json.NewDecoder(r.Body).Decode(&dieta); err != nil {
		fmt.Println("Errore DECODE JSON DIETA: ", err)
		http.Error(w, "Dati JSON non validi", http.StatusBadRequest)
		return
	}

	// GORM salverà automaticamente a cascata GiorniDieta, Pasti, Alimenti e Alternative
	if err := c.DB.Create(&dieta).Error; err != nil {
		fmt.Println("ERRORE GORM SALVATAGGIO DIETA: ", err)
		http.Error(w, "Errore salvataggio: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(dieta)
}

func (c *DietaController) GetDietaAttivaAtleta(w http.ResponseWriter, r *http.Request) {
	idAtleta := chi.URLParam(r, "idAtleta")
	var dieta models.Dieta

	result := c.DB.Preload("GiorniDieta").
		Preload("GiorniDieta.Pasti").
		Preload("GiorniDieta.Pasti.Alimenti").
		Preload("GiorniDieta.Pasti.Alimenti.Alternative"). // Aggiunto per caricare le alternative
		Where("id_atleta = ? AND attiva = ?", idAtleta, true).
		First(&dieta)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			http.Error(w, "Nessuna dieta attiva trovata per questo atleta", http.StatusNotFound)
			return
		}
		http.Error(w, "Errore nel recupero della dieta", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dieta)
}

func (c *DietaController) GetDietaByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var dieta models.Dieta

	result := c.DB.Preload("GiorniDieta").
		Preload("GiorniDieta.Pasti").
		Preload("GiorniDieta.Pasti.Alimenti").
		Preload("GiorniDieta.Pasti.Alimenti.Alternative"). // Aggiunto per caricare le alternative
		First(&dieta, "id = ?", id)

	if result.Error != nil {
		http.Error(w, "Dieta non trovata", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dieta)
}

func (c *DietaController) UpdateDieta(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var datiAggiornati models.Dieta

	if err := json.NewDecoder(r.Body).Decode(&datiAggiornati); err != nil {
		http.Error(w, "Dati JSON non trovati", http.StatusBadRequest)
		return
	}

	var dieta models.Dieta
	if err := c.DB.First(&dieta, "id = ?", id).Error; err != nil {
		http.Error(w, "Dieta non trovata", http.StatusNotFound)
		return
	}

	// Aggiornamento dei campi base
	dieta.Nome = datiAggiornati.Nome
	dieta.DataInizio = datiAggiornati.DataInizio
	dieta.DataFine = datiAggiornati.DataFine
	dieta.Attiva = datiAggiornati.Attiva
	dieta.Acqua = datiAggiornati.Acqua // Aggiunto nuovo campo
	dieta.Sale = datiAggiornati.Sale   // Aggiunto nuovo campo

	if err := c.DB.Save(&dieta).Error; err != nil {
		http.Error(w, "Errore durante l'aggiornamento", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dieta)
}

func (c *DietaController) DeleteDieta(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := c.DB.Delete(&models.Dieta{}, "id = ?", id).Error; err != nil {
		http.Error(w, "Errore durante l'eliminazione", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (c *DietaController) GetDieteAtleta(w http.ResponseWriter, r *http.Request) {
	idAtleta := chi.URLParam(r, "idAtleta")
	var diete []models.Dieta

	c.DB.Preload("GiorniDieta").
		Preload("GiorniDieta.Pasti").
		Preload("GiorniDieta.Pasti.Alimenti").
		Preload("GiorniDieta.Pasti.Alimenti.Alternative").
		Where("id_atleta = ?", idAtleta).
		Order("created_at desc").
		Find(&diete)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(diete)
}
