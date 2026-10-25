package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/adrianopiccinini/roccia-fit-api/pkg/models"
	"gorm.io/gorm"
)

type ChatAIController struct {
	DB *gorm.DB
}

// Struttura per ricevere la richiesta dall'app mobile
type AIRequest struct {
	IDAtleta  string `json:"id_atleta"`
	Messaggio string `json:"messaggio"`
}

// Strutture per comunicare con l'API di OpenAI
type OpenAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenAIRequest struct {
	Model    string          `json:"model"`
	Messages []OpenAIMessage `json:"messages"`
}

type OpenAIResponse struct {
	Choices []struct {
		Message OpenAIMessage `json:"message"`
	} `json:"choices"`
}

// GeneraRispostaAI riceve la domanda, inietta il contesto fisico e chiama l'AI
func (c *ChatAIController) GeneraRispostaAI(w http.ResponseWriter, r *http.Request) {
	var req AIRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Richiesta JSON non valida", http.StatusBadRequest)
		return
	}

	// 1. RECUPERO IL CONTESTO DELL'ATLETA DAL DATABASE
	var utente models.Utente
	c.DB.First(&utente, "id = ?", req.IDAtleta)

	var ultimaPesata models.PesoCorporeo
	c.DB.Where("id_atleta = ?", req.IDAtleta).Order("data desc").First(&ultimaPesata)

	var dietaAttiva models.Dieta
	c.DB.Preload("GiorniDieta").Where("id_atleta = ? AND attiva = true", req.IDAtleta).First(&dietaAttiva)

	kcal := 0
	if len(dietaAttiva.GiorniDieta) > 0 {
		kcal = dietaAttiva.GiorniDieta[0].TargetKcal
	}

	// 2. CREAZIONE DEL SYSTEM PROMPT INVISIBILE (La Magia)
	systemPrompt := fmt.Sprintf(`Sei il Coach virtuale AI dell'app RocciaFit. 
	Parli in modo diretto, motivante e tecnico (stile bodybuilding). 
	Ecco i dati attuali dell'atleta con cui stai parlando:
	- Nome: %s
	- Categoria: %s
	- Ultimo peso registrato: %.1f kg
	- Calorie giornaliere target: %d kcal.
	Rispondi alla sua domanda basandoti su questi dati in modo conciso, senza superare le 4-5 righe.`,
		utente.Nome, utente.Categoria, ultimaPesata.Peso, kcal)

	// 3. PREPARAZIONE DEL PAYLOAD PER OPENAI
	openAIReq := OpenAIRequest{
		Model: "gpt-4o-mini", // Il modello economico ma intelligente
		Messages: []OpenAIMessage{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: req.Messaggio},
		},
	}

	bodyBytes, _ := json.Marshal(openAIReq)

	// 4. CHIAMATA HTTP ALL'API DELL'AI
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		http.Error(w, "API Key mancante nel server", http.StatusInternalServerError)
		return
	}

	httpReq, _ := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(bodyBytes))
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil || resp.StatusCode != 200 {
		http.Error(w, "Errore di comunicazione con l'intelligenza artificiale", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// 5. DECODIFICA DELLA RISPOSTA
	bodyData, _ := io.ReadAll(resp.Body)
	var openAIResp OpenAIResponse
	json.Unmarshal(bodyData, &openAIResp)

	rispostaFinale := "Scusa Roccia, non sono riuscito a elaborare il pensiero."
	if len(openAIResp.Choices) > 0 {
		rispostaFinale = openAIResp.Choices[0].Message.Content
	}

	// 6. INVIO AL FRONTEND
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"risposta": rispostaFinale,
	})
}
