package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/adrianopiccinini/roccia-fit-api/pkg/models" // Assicurati che il path sia quello del tuo modulo
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

type MessaggioController struct {
	DB *gorm.DB
}

// Invia un nuovo messaggio
func (c *MessaggioController) InviaMessaggio(w http.ResponseWriter, r *http.Request) {
	var msg models.Messaggio
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		http.Error(w, "JSON non valido", http.StatusBadRequest)
		return
	}

	if err := c.DB.Create(&msg).Error; err != nil {
		http.Error(w, "Errore salvataggio messaggio", http.StatusInternalServerError)
		return
	}

	// 🚀 BONUS: Inviamo la notifica push al destinatario
	var destinatario models.Utente
	if err := c.DB.First(&destinatario, "id = ?", msg.IDDestinatario).Error; err == nil {
		if destinatario.PushToken != "" {
			go InviaNotificaExpo(destinatario.PushToken, "Nuovo Messaggio 💬", msg.Testo)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(msg)
}

// Recupera la chat tra due utenti specifici (es. Coach e Atleta)
func (c *MessaggioController) GetCronologiaChat(w http.ResponseWriter, r *http.Request) {
	idUtente1 := chi.URLParam(r, "id1")
	idUtente2 := chi.URLParam(r, "id2")

	var messaggi []models.Messaggio

	// Cerchiamo i messaggi dove (Mittente=1 E Destinatario=2) OPPURE (Mittente=2 E Destinatario=1)
	err := c.DB.Where(
		"(id_mittente = ? AND id_destinatario = ?) OR (id_mittente = ? AND id_destinatario = ?)",
		idUtente1, idUtente2, idUtente2, idUtente1,
	).Order("data_invio asc").Find(&messaggi).Error

	if err != nil {
		http.Error(w, "Errore nel recupero della chat", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messaggi)
}
