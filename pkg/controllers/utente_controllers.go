package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/adrianopiccinini/roccia-fit-api/pkg/models"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UtenteController struct {
	DB *gorm.DB
}

// CreaUtente godoc
// @Summary Crea utente (Admin)
// @Description Inserisce un nuovo utente manualmente nel database
// @Tags Utenti
// @Accept json
// @Produce json
// @Param utente body models.Utente true "Dati utente"
// @Security BearerAuth
// @Success 201 {object} models.Utente
// @Router /utenti [post]
func (c *UtenteController) CreaUtente(w http.ResponseWriter, r *http.Request) {

	var nuovoUtente models.Utente

	// 1. Decodifica il JSON inviato dall'app e lo inserisce nella struct
	err := json.NewDecoder(r.Body).Decode(&nuovoUtente)
	if err != nil {
		http.Error(w, "Errore nella lettura dei dati", http.StatusBadRequest)
		return
	}

	// creiamo il nuovo id del utente da inserire
	nuovoUtente.ID = uuid.New().String()

	// salviamo il nuovo Utente nel database tramite GORM
	risultato := c.DB.Create(&nuovoUtente)
	if risultato.Error != nil {
		http.Error(w, "Errore durante il salvataggio nel database", http.StatusBadRequest)
		return
	}

	// risponde all app con il nuovo utente appena creato
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(nuovoUtente)
}

// CreaUtentiMultipli godoc
// @Summary Creazione massiva utenti
// @Description Inserisce un array di utenti in un colpo solo
// @Tags Utenti
// @Accept json
// @Produce json
// @Param utenti body []models.Utente true "Array di utenti"
// @Security BearerAuth
// @Success 201 {array} models.Utente
// @Router /utenti/bulk [post]
func (c *UtenteController) CreaUtentiMultipli(w http.ResponseWriter, r *http.Request) {

	var nuoviUtenti []models.Utente

	if err := json.NewDecoder(r.Body).Decode(&nuoviUtenti); err != nil {
		http.Error(w, "Errore nella lettura dei dati multipli", http.StatusBadRequest)
		return
	}

	for i := range nuoviUtenti {
		nuoviUtenti[i].ID = uuid.New().String()
	}

	if risultato := c.DB.Create(&nuoviUtenti); risultato.Error != nil {
		http.Error(w, "Errore durante il salvataggio multiplo", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(nuoviUtenti)
}

// GetUtenti godoc
// @Summary Elenco utenti
// @Description Recupera la lista di tutti gli utenti registrati
// @Tags Utenti
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.Utente
// @Router /utenti [get]
func (c *UtenteController) GetUtenti(w http.ResponseWriter, r *http.Request) {
	var utenti []models.Utente

	risultato := c.DB.Find(&utenti)
	if risultato.Error != nil {
		http.Error(w, "Errore durante il recupero degli utenti", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(utenti)
}

// GetUtenteByID godoc
// @Summary Singolo utente
// @Description Recupera i dettagli di uno specifico utente tramite ID
// @Tags Utenti
// @Produce json
// @Param id path string true "ID Utente"
// @Security BearerAuth
// @Success 200 {object} models.Utente
// @Router /utenti/{id} [get]
func (c *UtenteController) GetUtenteByID(w http.ResponseWriter, r *http.Request) {

	// Estraiamo il valore id dall'URL
	id := chi.URLParam(r, "id")

	var utente models.Utente
	if err := c.DB.First(&utente, "id = ?", id).Error; err != nil {
		http.Error(w, "Utente non trovato", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(utente)
}

// UpdateUtente godoc
// @Summary Aggiorna utente
// @Description Modifica i dati di un utente esistente
// @Tags Utenti
// @Accept json
// @Produce json
// @Param id path string true "ID Utente"
// @Param utente body models.Utente true "Dati da aggiornare"
// @Security BearerAuth
// @Success 200 {object} models.Utente
// @Router /utenti/{id} [put]
func (c *UtenteController) UpdateUtente(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")
	var utente models.Utente
	if err := c.DB.First(&utente, "id = ?", id).Error; err != nil {
		http.Error(w, "Utente non trovato", http.StatusNotFound)
		return
	}

	// decodifichiamo le modifiche inviate dal client
	var datiAggiornati models.Utente
	if err := json.NewDecoder(r.Body).Decode(&datiAggiornati); err != nil {
		http.Error(w, "Dati non validi", http.StatusBadRequest)
		return
	}

	// applichiamo i nuovi dati al DB
	c.DB.Model(&utente).Updates(datiAggiornati)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(utente)
}

// DeleteUtente godoc
// @Summary Elimina utente
// @Description Rimuove un utente dal database
// @Tags Utenti
// @Produce json
// @Param id path string true "ID Utente"
// @Security BearerAuth
// @Success 200 {string} string "Utente eliminato"
// @Router /utenti/{id} [delete]
func (c *UtenteController) DeleteUtente(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")
	if err := c.DB.Delete(&models.Utente{}, "id = ?", id).Error; err != nil {
		http.Error(w, "Errore durante l'eliminazione", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Utente eliminato con successo"}`))
}

func (c *UtenteController) AggiungiXP(w http.ResponseWriter, r *http.Request) {
	idAtleta := chi.URLParam(r, "id")
	var input struct {
		XP int `json:"xp"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "JSON non valido", http.StatusBadRequest)
		return
	}

	var utente models.Utente
	if err := c.DB.First(&utente, "id = ?", idAtleta).Error; err != nil {
		http.Error(w, "Utente non trovato", http.StatusNotFound)
		return
	}

	// Aggiungiamo i punti
	utente.PuntiXP += input.XP

	// Logica di level-up (ogni livello richiede livello * 100 XP)
	xpNecessario := utente.Livello * 100
	if utente.PuntiXP >= xpNecessario {
		utente.PuntiXP -= xpNecessario
		utente.Livello++
	}

	c.DB.Save(&utente)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(utente)
}

// 1. GESTIONE DEI CHECK-IN E XP
func (c *UtenteController) GestisciCheckIn(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	idAtleta, _ := strconv.Atoi(idStr)

	var payload struct {
		TipoTask string `json:"tipo_task"`
		XPGained int    `json:"xp_gained"`
		Data     string `json:"data"`
	}
	json.NewDecoder(r.Body).Decode(&payload)

	dataTask, _ := time.Parse("2006-01-02", payload.Data)

	var checkin models.CheckIn
	c.DB.Where("atleta_id = ? AND data = ?", idAtleta, dataTask).FirstOrCreate(&checkin, models.CheckIn{
		AtletaID: idStr,
		Data:     dataTask,
	})

	// Toggle della spunta in base al task ricevuto
	switch payload.TipoTask {
	case "allenamento":
		checkin.AllenamentoCompletato = !checkin.AllenamentoCompletato
	case "dieta":
		checkin.DietaCompletata = !checkin.DietaCompletata
	case "cardio":
		checkin.CardioCompletato = !checkin.CardioCompletato
	}

	checkin.XPGained += payload.XPGained
	c.DB.Save(&checkin)

	// Aggiornamento livello Utente
	var utente models.Utente
	c.DB.First(&utente, idAtleta)
	utente.PuntiXP += payload.XPGained

	// Logica Level-Up
	xpNecessari := utente.Livello * 100
	if utente.PuntiXP >= xpNecessari {
		utente.PuntiXP -= xpNecessari
		utente.Livello++
	}
	c.DB.Save(&utente)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(utente)
}

// 2. GESTIONE DELL'ACQUA
func (c *UtenteController) GestisciAcqua(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	idAtleta, _ := strconv.Atoi(idStr)

	var payload struct {
		Litri float64 `json:"litri"`
		Data  string  `json:"data"`
	}
	json.NewDecoder(r.Body).Decode(&payload)

	dataTask, _ := time.Parse("2006-01-02", payload.Data)

	var checkin models.CheckIn
	c.DB.Where("atleta_id = ? AND data = ?", idAtleta, dataTask).FirstOrCreate(&checkin, models.CheckIn{
		AtletaID: idStr,
		Data:     dataTask,
	})

	checkin.LitriAcqua = payload.Litri
	c.DB.Save(&checkin)

	w.WriteHeader(http.StatusOK)
}

// 3. RECUPERO DELLO STORICO (PER CALENDARIO E PAGINA STORICO)
func (c *UtenteController) GetStoricoAtleta(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	var checkins []models.CheckIn

	// Ordina dal più recente (Oggi) al più vecchio
	c.DB.Where("atleta_id = ?", idStr).Order("data desc").Find(&checkins)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(checkins)
}

type AtletaRank struct {
	ID    string `json:"id"`
	Nome  string `json:"nome"`
	Punti int    `json:"punti"`
	Rank  int    `json:"rank"`
}

// GetClassifica calcola e restituisce la leaderboard globale
func (c *UtenteController) GetClassifica(w http.ResponseWriter, r *http.Request) {
	var utenti []models.Utente

	// Recuperiamo solo gli atleti, ordinati per XP
	if err := c.DB.Where("LOWER(ruolo) = ?", "atleta").Order("punti_xp DESC").Limit(50).Find(&utenti).Error; err != nil {
		http.Error(w, "Errore nel recupero della classifica", http.StatusInternalServerError)
		return
	}

	// Inizializza una slice vuota (evita di restituire "null" al frontend)
	classifica := make([]AtletaRank, 0)

	for i, u := range utenti {
		classifica = append(classifica, AtletaRank{
			// Convertiamo l'ID direttamente in stringa senza controlli nil
			ID:   fmt.Sprintf("%v", u.ID),
			Nome: u.Nome,
			// ⚠️ ATTENZIONE QUI: se ti dà errore su "u.Punti_xp",
			// guardalo nel tuo file models.go e cambialo in u.PuntiXP o u.PuntiXp
			Punti: u.PuntiXP,
			Rank:  i + 1,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(classifica)
}

func (c *UtenteController) SalvaPushToken(w http.ResponseWriter, r *http.Request) {
	idAtleta := chi.URLParam(r, "id")

	// Contenuto della notifica
	var payload struct {
		PushToken string `json:"push_token"`
	}

	// Controlliamo l'integrita del payload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "JSON non valido", http.StatusBadRequest)
		return
	}

	var utente models.Utente
	if err := c.DB.First(&utente, "id = ?", idAtleta).Error; err != nil {
		http.Error(w, "Utente non trovato", http.StatusNotFound)
		return
	}

	utente.PushToken = payload.PushToken
	c.DB.Save(&utente)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`"message" : "Push token salvato con successo"`))
}
