package controllers

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/adrianopiccinini/roccia-fit-api/pkg/models"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CheckFotoController struct {
	DB *gorm.DB
}

// UploadImmagine godoc
// @Summary Carica foto condizione
// @Description Salva il file fisico della foto di check (Max 10MB)
// @Tags Check Fotografici
// @Accept multipart/form-data
// @Produce json
// @Param foto formData file true "File immagine"
// @Security BearerAuth
// @Success 201 {object} map[string]string
// @Router /check/upload [post]
func (c *CheckFotoController) UploadImmagine(w http.ResponseWriter, r *http.Request) {
	// Limitiamo la dimensione della richiesta (es. 10 MB massimo)
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "File troppo grande o richiesta non valida", http.StatusBadRequest)
		return
	}

	// Estraiamo il file dalla richiesta e si chiamera foto
	file, handler, err := r.FormFile("foto")
	if err != nil {
		http.Error(w, "Errore nel recupero del file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	//Creazione della cartella degli upload in caso non ci sia
	uploadDir := "./uploads/checks"
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		os.MkdirAll(uploadDir, os.ModePerm)
	}

	// Generiamo un nome univoco uuid + estensione
	estensione := filepath.Ext(handler.Filename)
	nuovoNome := uuid.New().String() + estensione
	percorsoCompleto := filepath.Join(uploadDir, nuovoNome)

	// Creiamo il file fisico sul server
	dst, err := os.Create(percorsoCompleto)
	if err != nil {
		http.Error(w, "Errore nella creazione del file sul server", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// Copiamo i byte caricati nel file del server
	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, "Errore durante il salvataggio dei dati", http.StatusInternalServerError)
		return
	}

	// Restituiamo il percorso relativo da salvare nel database
	urlFile := "/uploads/checks/" + nuovoNome

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"url_foto": urlFile})
}

// AggiungiCheckFoto salva i dati e il link dell'immagine del check

// AggiungiCheckFoto godoc
// @Summary Salva dati check
// @Description Registra nel database il link della foto caricata e la data
// @Tags Check Fotografici
// @Accept json
// @Produce json
// @Param check body models.CheckFoto true "Dati check"
// @Security BearerAuth
// @Success 201 {object} models.CheckFoto
// @Router /check [post]
func (c *CheckFotoController) AggiungiCheckFoto(w http.ResponseWriter, r *http.Request) {
	var nuovoCheck models.CheckFoto

	if err := json.NewDecoder(r.Body).Decode(&nuovoCheck); err != nil {
		http.Error(w, "Errore nella decodifica del JSON", http.StatusBadRequest)
		return
	}

	nuovoCheck.ID = uuid.New().String()

	if nuovoCheck.Data.IsZero() {
		nuovoCheck.Data = time.Now()
	}

	if risultato := c.DB.Create(&nuovoCheck); risultato.Error != nil {
		http.Error(w, "Errore durante il salvataggio del check", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(nuovoCheck)
}

// GetCheckAtleta recupera la galleria storica dei check di un atleta

// GetCheckAtleta godoc
// @Summary Storico check
// @Description Recupera la galleria storica dei check di un atleta
// @Tags Check Fotografici
// @Produce json
// @Param idAtleta path string true "ID Atleta"
// @Security BearerAuth
// @Success 200 {array} models.CheckFoto
// @Router /check/atleta/{idAtleta} [get]
func (c *CheckFotoController) GetCheckAtleta(w http.ResponseWriter, r *http.Request) {
	idAtleta := chi.URLParam(r, "idAtleta")

	var checkFoto []models.CheckFoto

	if err := c.DB.Where("id_atleta = ?", idAtleta).Order("data desc").Find(&checkFoto).Error; err != nil {
		http.Error(w, "Errore durante il recupero dei check", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(checkFoto)
}

// DeleteCheckFoto rimuove un check specifico

// DeleteCheckFoto godoc
// @Summary Elimina check
// @Description Rimuove logicamente (Soft Delete) una foto di check
// @Tags Check Fotografici
// @Produce json
// @Param id path string true "ID Check"
// @Security BearerAuth
// @Success 200 {string} string "Foto eliminata"
// @Router /check/{id} [delete]
func (c *CheckFotoController) DeleteCheckFoto(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := c.DB.Delete(&models.CheckFoto{}, "id = ?", id).Error; err != nil {
		http.Error(w, "Errore durante l'eliminazione", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Foto di check eliminata con successo"}`))
}

// GetCheckFotoByID godoc
// @Summary Singolo check
// @Description Recupera i dettagli di una singola foto
// @Tags Check Fotografici
// @Produce json
// @Param id path string true "ID Check"
// @Security BearerAuth
// @Success 200 {object} models.CheckFoto
// @Router /check/{id} [get]
func (c *CheckFotoController) GetCheckFotoByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var check models.CheckFoto
	if err := c.DB.First(&check, "id = ?", id).Error; err != nil {
		http.Error(w, "Check Fotografico non trovato", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(check)
}

// UpdateCheckFoto godoc
// @Summary Aggiorna check
// @Description Modifica i dati di un check fotografico esistente
// @Tags Check Fotografici
// @Accept json
// @Produce json
// @Param id path string true "ID Check"
// @Param check body models.CheckFoto true "Dati da aggiornare"
// @Security BearerAuth
// @Success 200 {object} models.CheckFoto
// @Router /check/{id} [put]
func (c *CheckFotoController) UpdateCheckFoto(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var check models.CheckFoto
	if err := c.DB.First(&check, "id = ?", id).Error; err != nil {
		http.Error(w, "Check fotografico non trovato", http.StatusNotFound)
		return
	}

	var datiAggiornati models.CheckFoto
	if err := json.NewDecoder(r.Body).Decode(&datiAggiornati); err != nil {
		http.Error(w, "Errore durante la creazione", http.StatusInternalServerError)
		return
	}

	c.DB.Model(&check).Updates(datiAggiornati)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(check)
}
