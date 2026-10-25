package controllers

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/adrianopiccinini/roccia-fit-api/pkg/models" // Adatta questo path alla tua struttura

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

type MessaggioController struct {
	DB *gorm.DB
}

// Strutture WebSocket
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

var clients = make(map[string]*websocket.Conn)
var clientsMu sync.Mutex

// ⚡ 1. Rotta per inizializzare il WebSocket
func (c *MessaggioController) ConnettiWS(w http.ResponseWriter, r *http.Request) {
	idUtente := chi.URLParam(r, "id")

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	clientsMu.Lock()
	clients[idUtente] = conn
	clientsMu.Unlock()

	go func() {
		defer func() {
			clientsMu.Lock()
			delete(clients, idUtente)
			clientsMu.Unlock()
			conn.Close()
		}()
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				break
			}
		}
	}()
}

// ⚡ 2. Invio e salvataggio
func (c *MessaggioController) InviaMessaggio(w http.ResponseWriter, r *http.Request) {
	var messaggio models.Messaggio
	if err := json.NewDecoder(r.Body).Decode(&messaggio); err != nil {
		http.Error(w, "Payload invalido", http.StatusBadRequest)
		return
	}

	if err := c.DB.Create(&messaggio).Error; err != nil {
		http.Error(w, "Errore salvataggio: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Notifica in tempo reale al destinatario se online
	clientsMu.Lock()
	if destConn, online := clients[messaggio.IDDestinatario]; online {
		destConn.WriteJSON(messaggio)
	}
	clientsMu.Unlock()

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(messaggio)
}

// ⚡ 3. Caricamento cronologia (Al primo avvio)
func (c *MessaggioController) GetCronologiaChat(w http.ResponseWriter, r *http.Request) {
	idUtente1 := chi.URLParam(r, "id1")
	idUtente2 := chi.URLParam(r, "id2")

	var messaggi []models.Messaggio
	err := c.DB.Where(
		"(id_mittente = ? AND id_destinatario = ?) OR (id_mittente = ? AND id_destinatario = ?)",
		idUtente1, idUtente2, idUtente2, idUtente1,
	).Order("data_invio asc").Find(&messaggi).Error

	if err != nil {
		http.Error(w, "Errore DB: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messaggi)
}
