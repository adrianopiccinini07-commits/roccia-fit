package main

import (
	"fmt"
	"log"
	"net/http"
	"os" // lettura delle variabili
	"path/filepath"

	"github.com/adrianopiccinini/roccia-fit-api/pkg/controllers"
	"github.com/adrianopiccinini/roccia-fit-api/pkg/database"
	"github.com/adrianopiccinini/roccia-fit-api/pkg/middleware"
	"github.com/adrianopiccinini/roccia-fit-api/pkg/models"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "github.com/adrianopiccinini/roccia-fit-api/docs"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv" // Metodo per caricare il file .env
)

// @title RocciaFit API
// @version 1.0
// @description API per gestione la gestione dei parametri d'allenamento e condizione fisica
// @termsOfService http://swagger.io/terms/
// @host localhost
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorzation

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("Nessun file .env trovato, verrano  usate le variabili di sistema")
	}

	dsn := os.Getenv("DB_DSN")

	if dsn == "" {
		log.Fatal("DB_DSN non impostato nelle variabili d'ambiente")
	}

	//inizializziamo la connessione al db
	db := database.InitDB(dsn)

	//defer mette in coda questa istruzione istruzione: chiudere la connesione
	// che verra attivata solo quando chiudero il server go
	// Estraiamo l'oggetto SQL generico per poter usare il defer Close
	sqlDB, err := db.DB()
	if err == nil {
		defer sqlDB.Close()
	}

	// AUTO-MIGRAZIONE
	//GORM legge le struct e crea la tabella su PostgreSQL
	fmt.Println("Esecuzione della migrazione del database...")
	if err := db.AutoMigrate(
		&models.Utente{},
		&models.SchedaAllenamento{},
		&models.GiornoAllenamento{},
		&models.EsercizioScheda{},
		&models.Esercizio{},
		&models.Allenamento{},
		&models.SerieEseguita{},
		&models.PesoCorporeo{},
		&models.CheckFoto{},
		&models.Dieta{},
		&models.GiornoDieta{},
		&models.Pasto{},
		&models.AlimentoPasto{},
		&models.AlternativaAlimento{},
		&models.ProtocolloCardio{},
		&models.SessioneCardio{},
		&models.Feedback{},
		&models.CheckIn{},
	); err != nil {
		log.Fatalf("Errore durante la migrazione: %v", err) // %v e' un segnaposto usato per stampare log
	}
	fmt.Println("Migrazione completata con successo!")

	utenteController := &controllers.UtenteController{DB: db}
	schedeController := &controllers.SchedaController{DB: db}
	esercizioController := &controllers.EsercizioController{DB: db}
	allenamentoController := &controllers.AllenamentoController{DB: db}
	pesoController := &controllers.PesoController{DB: db}
	checkFotoController := &controllers.CheckFotoController{DB: db}
	authController := &controllers.AuthController{DB: db}
	statsController := &controllers.StatsController{DB: db}
	dietaController := &controllers.DietaController{DB: db}
	cardioController := &controllers.CardioController{DB: db}
	feedbackController := &controllers.FeedbackController{DB: db}

	// Creiamo il ROUTER chi
	r := chi.NewRouter()

	// Leggi i domini permessi dal file .env, altrimenti usa "*" solo come fallback di sviluppo
	allowedOrigin := os.Getenv("ALLOWED_ORIGIN")
	if allowedOrigin == "" {
		allowedOrigin = "*"
	}

	// CONFIGURAZIONE PER L'APP MOBILE
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{allowedOrigin},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Use(chiMiddleware.Logger)    // stampa a terminale ogni richiesta HTTP
	r.Use(chiMiddleware.Recoverer) // Impedisce al server di crashare se ce un errore grave

	// ESPOSIZIONE DEI FILE STATICI: permette di visualizzare le foto tramite URL
	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, "uploads"))
	r.Handle("/uploads/*", http.StripPrefix("/uploads", http.FileServer(filesDir)))

	r.Get("/health", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "dio bessia" , "message": "API operative"}`))
	})

	r.Get("/swagger/*", httpSwagger.WrapHandler)

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Post("/registrati", authController.RegistrazioneSicura)
			r.Post("/login", authController.Login)
		})

		r.Group(func(r chi.Router) {
			//applica il middleware a tutto quello che si trova nel middleware
			r.Use(middleware.AuthMiddleware)

			r.Route("/utenti", func(r chi.Router) {
				r.Post("/", utenteController.CreaUtente)
				r.Post("/bulk", utenteController.CreaUtentiMultipli)
				r.Get("/", utenteController.GetUtenti)
				r.Get("/classifica", utenteController.GetClassifica)

				r.Route("/{id}", func(r chi.Router) {
					r.Get("/", utenteController.GetUtenteByID)
					r.Put("/", utenteController.UpdateUtente)
					r.Delete("/", utenteController.DeleteUtente)
					r.Post("/xp", utenteController.AggiungiXP)
					r.Post("/checkin", utenteController.GestisciCheckIn)
					r.Post("/acqua", utenteController.GestisciAcqua)
				})
			})
			r.Route("/schede", func(r chi.Router) {
				r.Get("/test", func(w http.ResponseWriter, req *http.Request) {
					w.Write([]byte("Ecco le schede di allenamento"))
				})

				r.Post("/", schedeController.CreaScheda)
				r.Get("/atleta/{idAtleta}", schedeController.GetSchedeAtleta)

				r.Route("/{id}", func(r chi.Router) {
					r.Put("/", schedeController.UpdateScheda)
					r.Get("/", schedeController.GetSchedaByID)
					r.Delete("/", schedeController.DeleteScheda)
				})
			})
			r.Route("/esercizi", func(r chi.Router) {
				r.Get("/", esercizioController.GetEsercizi)
				r.Post("/", esercizioController.CreaEsercizio)
				r.Post("/bulk", esercizioController.CreaEserciziMultipli)
				r.Route("/{id}", func(r chi.Router) {
					r.Get("/", esercizioController.GetEsercizioByID)
					r.Put("/", esercizioController.UpdateEsercizio)
					r.Delete("/", esercizioController.DeleteEsercizio)
				})
			})
			r.Route("/allenamenti", func(r chi.Router) {

				r.Post("/inizia", allenamentoController.IniziaAllenamento)
				r.Get("/atleta/{idAtleta}", allenamentoController.GetAllenamentiAtleta)

				r.Route("/serie", func(r chi.Router) {
					r.Post("/", allenamentoController.AggiungiSerie)
					r.Put("/{id}", allenamentoController.UpdateSerie)
					r.Delete("/{id}", allenamentoController.DeleteSerie)
				})

				r.Route("/{id}", func(r chi.Router) {
					r.Get("/", allenamentoController.GetAllenamentoById)
					r.Delete("/", allenamentoController.DeleteAllenamento)
				})
			})
			r.Route("/peso", func(r chi.Router) {
				r.Post("/", pesoController.AggiungiPeso)
				r.Post("/bulk", pesoController.AggiungiPesiMultipli)
				r.Get("/atleta/{idAtleta}", pesoController.GetStoricoPeso)
				r.Get("/ultimo/{idAtleta}", pesoController.GetUltimoPeso)

				r.Route("/{id}", func(r chi.Router) {
					r.Get("/", pesoController.GetPesoByID)
					r.Put("/", pesoController.UpdatePeso)
					r.Delete("/", pesoController.DeletePeso)
				})
			})
			r.Route("/check", func(r chi.Router) {
				r.Post("/", checkFotoController.AggiungiCheckFoto)
				r.Post("/upload", checkFotoController.UploadImmagine)
				r.Get("/atleta/{idAtleta}", checkFotoController.GetCheckAtleta)

				r.Route("/{id}", func(r chi.Router) {
					r.Get("/", checkFotoController.GetCheckFotoByID)
					r.Put("/", checkFotoController.UpdateCheckFoto)
					r.Delete("/", checkFotoController.DeleteCheckFoto)
				})
			})
			r.Route("/stats", func(r chi.Router) {
				r.Get("/tonnellaggio/{idAllenamento}", statsController.GetTonnellaggio)
				r.Get("/massimale/{idSerie}", statsController.GetProiezione1RM)
				r.Get("/volume/{idAtleta}", statsController.GetVolumeSettimanale)
			})
			r.Route("/dieta", func(r chi.Router) {
				r.Get("/atleta/{idAtleta}", dietaController.GetDietaAttivaAtleta)
				r.Get("/atleta/{idAtleta}/all", dietaController.GetDieteAtleta)
				r.Post("/", dietaController.CreaDieta)

				r.Route("/{id}", func(r chi.Router) {
					r.Get("/", dietaController.GetDietaByID)
					r.Put("/", dietaController.UpdateDieta)
					r.Delete("/", dietaController.DeleteDieta)
				})
			})
			r.Route("/storico", func(r chi.Router) {
				r.Get("/atleta/{id}", utenteController.GetStoricoAtleta)
			})

			// --- ROTTE CARDIO ---
			r.Get("/cardio/atleta/{idAtleta}", cardioController.GetCardioAtleta)
			r.Post("/cardio", cardioController.SalvaCardio)
			r.Post("/cardio/", cardioController.SalvaCardio)

			// --- ROTTE FEEDBACK ---
			r.Get("/feedback/atleta/{idAtleta}", feedbackController.GetFeedbackAtleta)
			r.Post("/feedback", feedbackController.CreaFeedback)
			r.Post("/feedback/", feedbackController.CreaFeedback)

			// --- CAMBIO PASSWORD ---
			r.Put("/auth/password", authController.CambiaPassword)
		})
	})

	//qui caricheremo le variabili d'ambiente
	port := os.Getenv("PORT")
	if port == "" {
		port = ":8080"
	}

	fmt.Printf("Server in ascolto sulla porta %s...\n", port)

	// accendiamo il server e vediamo se produce un errore(siamo in TCP)
	if err := http.ListenAndServe(port, r); err != nil {
		log.Fatalf("Errore durante l'avvio del server: %v", err) // stampa l'errore e poi chiama os.Exit(1) chiude il programma
	}
}
