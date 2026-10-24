package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/adrianopiccinini/roccia-fit-api/pkg/models"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

type FeedbackController struct {
	DB *gorm.DB
}

// Funzione di utilità per inviare notifiche push tramite l'API di Expo
func InviaNotificaExpo(pushToken string, titolo string, corpo string) {
	if pushToken == "" {
		return
	}

	payload := map[string]interface{}{
		"to":    pushToken,
		"sound": "default",
		"title": titolo,
		"body":  corpo,
	}

	bodyBytes, _ := json.Marshal(payload)

	// Chiamata HTTP verso il servizio di Expo Push
	resp, err := http.Post("https://exp.host/--/api/v2/push/send", "application/json", bytes.NewBuffer(bodyBytes))
	if err != nil {
		return
	}
	defer resp.Body.Close()
}

// Salva un nuovo messaggio del coach e invia la notifica push
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

	// --- 🚀 INVIO DELLA NOTIFICA PUSH ALL'ATLETA ---
	var atleta models.Utente
	if err := c.DB.First(&atleta, "id = ?", feedback.IDAtleta).Error; err == nil {
		if atleta.PushToken != "" {
			titoloNotifica := feedback.Titolo
			if titoloNotifica == "" {
				titoloNotifica = "Nuovo messaggio dal Coach 🦅"
			}
			// Eseguiamo l'invio (anche in goroutine per non rallentare la risposta HTTP)
			go InviaNotificaExpo(atleta.PushToken, titoloNotifica, feedback.Testo)
		}
	}
	// ----------------------------------------------

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
