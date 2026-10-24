package controllers

import (
	"encoding/json" // serve a lavorare con i json
	"net/http"      // crea il server http gestendo post, get, put e delete
	"os"            // permette di leggere le variabili d'ambiente
	"time"          // date e orari

	"github.com/adrianopiccinini/roccia-fit-api/pkg/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthController struct {
	DB *gorm.DB
}

// Struttura per lettura dati login
type Credenziali struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CambioPswReq struct {
	Email      string `json:"email"`
	VecchiaPsw string `json:"vecchia_psw"`
	NuovaPsw   string `json:"nuova_psw"`
}

// Struttura del Token JWT (carta d'identita contenuta nel token) Json Web Token
type Claims struct {
	UtenteID string `json:"utente_id"`
	jwt.RegisteredClaims
}

// RegistrazioneSicura godoc
// @Summary Registrazione nuovo utente
// @Description Crea un account con password crittografata
// @Tags Autenticazione
// @Accept json
// @Produce json
// @Param utente body models.Utente true "Dati utente"
// @Success 201 {object} models.Utente
// @Failure 400 {string} string "Errore nella decodifica"
// @Router /auth/registrati [post]

func (c *AuthController) RegistrazioneSicura(w http.ResponseWriter, r *http.Request) {
	var nuovoUtente models.Utente

	if err := json.NewDecoder(r.Body).Decode(&nuovoUtente); err != nil {
		http.Error(w, "Errore nella decodifica", http.StatusBadRequest)
		return
	}

	nuovoUtente.CF = uuid.New().String()

	// Bcrypt prende la password in chiaro e genera l'hash irreversibile
	hash, err := bcrypt.GenerateFromPassword([]byte(nuovoUtente.Psw), 10)
	if err != nil {
		http.Error(w, "Errore nella crittografia della password", http.StatusInternalServerError)
		return
	}

	nuovoUtente.ID = uuid.New().String()
	nuovoUtente.Psw = string(hash) //Salviamo l'hash nel db non la password

	if risultato := c.DB.Create(&nuovoUtente); risultato.Error != nil {
		http.Error(w, "Errore durante la creazione dell'utente", http.StatusInternalServerError)
		return
	}

	nuovoUtente.Psw = ""

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(nuovoUtente)
}

// Login godoc
// @Summary Autenticazione utente
// @Description Verifica email e password e restituisce un token JWT
// @Tags Autenticazione
// @Accept json
// @Produce json
// @Param credenziali body Credenziali true "Dati di accesso"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {string} string "Richiesta non valida"
// @Failure 401 {string} string "Credenziali errate"
// @Router /auth/login [post]
func (c *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	var creds Credenziali
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Richiesta non valida", http.StatusBadRequest)
		return
	}

	// --- NUOVO CONTROLLO DI SICUREZZA ---
	// Blocchiamo immediatamente chi tenta di accedere con campi vuoti
	if creds.Email == "" || creds.Password == "" {
		http.Error(w, "Email e password sono obbligatori", http.StatusBadRequest)
		return
	}
	// ------------------------------------

	var utente models.Utente
	// 1. Cerchiamo l'utente tramite l'email
	if err := c.DB.Where("email = ?", creds.Email).First(&utente).Error; err != nil {
		http.Error(w, "Credenziali errate", http.StatusUnauthorized)
		return
	}

	// 2. Confrontiamo la password inserita con l'hash salvato nel database
	if err := bcrypt.CompareHashAndPassword([]byte(utente.Psw), []byte(creds.Password)); err != nil {
		http.Error(w, "Credenziali Errate", http.StatusUnauthorized)
		return
	}

	// Applichiamo una scadenza a tale token e la inseriamo in claims dove ci sono tutti i dati NON CIFRATI
	scadenza := time.Now().Add(7 * 24 * time.Hour)
	claims := &Claims{
		UtenteID: utente.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(scadenza),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	//JWT_SECRET e' una chiave che conosce solo il mio backend
	chiaveSegreta := os.Getenv("JWT_SECRET")
	tokenString, err := token.SignedString([]byte(chiaveSegreta))

	if err != nil {
		http.Error(w, "Errore nella generazione del token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"token": tokenString,
		"utente": map[string]string{
			"id":    utente.ID,
			"nome":  utente.Nome,
			"ruolo": utente.Ruolo,
		},
	})
}

// CambiaPassword godoc
// @Summary Modifica password
// @Description Aggiorna la password verificando la precedente
// @Tags Autenticazione
// @Accept json
// @Produce json
// @Param credenziali body CambioPswReq true "Dati cambio password"
// @Security BearerAuth
// @Success 200 {string} string "Password aggiornata con successo"
// @Router /auth/password [put]
func (c *AuthController) CambiaPassword(w http.ResponseWriter, r *http.Request) {
	var req CambioPswReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Richiesta non valida", http.StatusBadRequest)
		return
	}

	var utente models.Utente
	if err := c.DB.Where("email = ?", req.Email).First(&utente).Error; err != nil {
		http.Error(w, "Utente non trovato", http.StatusNotFound)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(utente.Psw), []byte(req.VecchiaPsw)); err != nil {
		http.Error(w, "La vecchia password non e' valida", http.StatusUnauthorized)
		return
	}

	nuovoHash, err := bcrypt.GenerateFromPassword([]byte(req.NuovaPsw), 10)
	if err != nil {
		http.Error(w, "Errore durante l'aggiornamento della password", http.StatusInternalServerError)
		return
	}

	utente.Psw = string(nuovoHash)
	if err := c.DB.Save(&utente).Error; err != nil {
		http.Error(w, "Errore durante l'aggironamento della password", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Password aggiornata con successo"}`))
}
