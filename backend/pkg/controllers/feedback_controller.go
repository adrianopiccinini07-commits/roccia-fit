package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/adrianopiccinini/roccia-fit-api/pkg/models"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

type FeedbackController struct {
	DB *gorm.DB
}

// Salva un nuovo messaggio del coach
func (c *FeedbackController) CreaFeedback(w http.ResponseWriter, r *http.Request) {
	var feedback models.Feedback
	if err := json.NewDecoder(r.Body).Decode(&feedback); err != nil {
		http.Error(w, "JSON non valido", http.StatusBadRequest)
		return
	}

	if err := c.DB.Create(&feedback).Error; err != nil {
		http.Error(w, "Errore salvataggio feedback: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(feedback)
}

// Recupera lo storico dei feedback per l'atleta (ordinati dal più recente)
func (c *FeedbackController) GetFeedbackAtleta(w http.ResponseWriter, r *http.Request) {
	idAtleta := chi.URLParam(r, "idAtleta")
	var feedbacks []models.Feedback

	if err := c.DB.Where("id_atleta = ?", idAtleta).Order("data desc").Find(&feedbacks).Error; err != nil {
		http.Error(w, "Errore nel recupero dei feedback", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(feedbacks)
}
