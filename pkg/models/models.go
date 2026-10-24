package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Utente struct {
	ID        string `gorm:"primaryKey;type:varchar(36)"`
	Nome      string
	Cognome   string
	CF        string `gorm:"unique;not null"`
	Telefono  string
	Email     string `gorm:"unique;not null"` // l'email non puo essere vuota
	Psw       string `json:"password"`
	Ruolo     string // "COACH" , "ATLETA" o "CLIENTE"
	Categoria string `gorm:"size:50" json:"categoria"`
	PuntiXP   int    `gorm:"default:0" json:"punti_xp"`
	Livello   int    `gorm:"default:1" json:"livello"`

	PushToken string `gorm:"size:255" json:"livello"`

	SchedaRicevute []SchedaAllenamento `gorm:"foreignKey:IDAtleta"`
	SchedaCreate   []SchedaAllenamento `gorm:"foreignKey:IDCoach"`

	Allenamenti  []Allenamento  `gorm:"foreignKey:IDAtleta"`
	PesiCorporeo []PesoCorporeo `gorm:"foreignKey:IDAtleta"`
	CheckFoto    []CheckFoto    `gorm:"foreignKey:IDAtleta"`

	AtletiSeguiti []*Utente `gorm:"many2many:coach_atleta;joinForeignKey:coach_id;joinReferences:atleta_id"`
}

type SchedaAllenamento struct {
	ID         string `gorm:"primaryKey;type:varchar(36)"`
	Nome       string
	DataInizio time.Time
	DataFine   time.Time
	IDCoach    string
	IDAtleta   string

	// Relazione 1 a N
	Giorni []GiornoAllenamento `gorm:"foreignKey:IDScheda"`
}

type GiornoAllenamento struct {
	ID       string `gorm:"primaryKey;type:varchar(36)"`
	IDScheda string
	Giorno   string

	Esercizi []EsercizioScheda `gorm:"foreignKey:IDGiorno"`
}

type EsercizioScheda struct {
	ID          string `json:"id" gorm:"primaryKey;type:varchar(36)"`
	IDGiorno    string `json:"id_giorno"`
	IDEsercizio string `json:"id_esercizio"`
	Ordine      int    `json:"ordine"`

	// Sostituiamo i numeri fissi con stringhe o campi descrittivi per gestire "6/10", "4/6", ecc.
	Serie       string `json:"serie" gorm:"size:20"`         // Es: "2x", "3x", "1x"
	Reps        string `json:"reps" gorm:"size:50"`          // Es: "6/10", "4/6 • 1xmax", "max"
	Rir         string `json:"rir" gorm:"size:20"`           // Es: "RiR2", "1ª serie RIR1"
	NoteTecnica string `json:"note_tecnica" gorm:"size:150"` // Es: "F in partenza 2” SQ2”"
	Recupero    int    `json:"recupero"`

	Esercizio Esercizio `json:"esercizio" gorm:"foreignKey:IDEsercizio"`
}

type Esercizio struct {
	ID              string `gorm:"primaryKey;type:varchar(36)"`
	Nome            string
	GruppoMuscolare string
	Descrizione     string
	Video           string
}

type Allenamento struct {
	ID       string `gorm:"primaryKey;type:varchar(36)"`
	IDAtleta string // chiave esterna verso Utente
	Data     time.Time
	Durata   float64
	Note     string

	SerieEseguite []SerieEseguita `gorm:"foreignKey:IDAllenamento"`
}

type SerieEseguita struct {
	ID                string `gorm:"primaryKey;type:varchar(36)"`
	IDAllenamento     string // chiave verso Allenamento
	IDEsercizioScheda string // chiave verso EsercizioScheda
	Peso              float64
	Reps              int
	Rir               int
}

type PesoCorporeo struct {
	ID       string `gorm:"primaryKey;type:varchar(36)"`
	IDAtleta string // chiave verso utente
	Data     time.Time
	Peso     float64
	Fame     int `gorm:"check:fame >= 1 AND fame <= 5"`
	Sonno    int `gorm:"check:sonno >= 1 AND sonno <= 5"`
}

type CheckFoto struct {
	ID       string `gorm:"primaryKey;type:varchar(36)"`
	IDAtleta string
	Data     time.Time
	UrlFoto  string

	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// --- STRUTTURE CARDIO E NEAT ---

type ProtocolloCardio struct {
	ID               string `gorm:"primaryKey;type:varchar(36)"`
	IDAtleta         string `gorm:"uniqueIndex;not null"` // UniqueIndex: 1 atleta = 1 protocollo cardio attivo
	PassiGiornalieri int
	NoteNeat         string `gorm:"size:255"`

	Sessioni []SessioneCardio `gorm:"foreignKey:IDProtocollo;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type SessioneCardio struct {
	ID           string `gorm:"primaryKey;type:varchar(36)"`
	IDProtocollo string `gorm:"not null"`
	Giorno       string `gorm:"size:100"`
	Macchinario  string `gorm:"size:100"`
	Durata       string `gorm:"size:50"`
	Intensita    string `gorm:"size:50"`
}

type Feedback struct {
	ID       string    `gorm:"primaryKey;type:varchar(36)"`
	IDAtleta string    `gorm:"not null;index"`
	IDCoach  string    `gorm:"not null"`
	Titolo   string    `gorm:"size:150"`
	Testo    string    `gorm:"type:text;not null"`
	Data     time.Time `gorm:"autoCreateTime"` // Inserisce in automatico il timestamp
}

type Messaggio struct {
	ID             string    `gorm:"primaryKey;type:varchar(36)" json:"id"`
	IDMittente     string    `gorm:"type:varchar(36);not null;index" json:"id_mittente"`
	IDDestinatario string    `gorm:"type:varchar(36);not null;index" json:"id_destinatario"`
	Testo          string    `gorm:"type:text;not null" json:"testo"`
	Letto          bool      `gorm:"default:false" json:"letto"`
	DataInvio      time.Time `gorm:"autoCreateTime" json:"data_invio"`
}

func (m *Messaggio) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == "" {
		m.ID = uuid.NewString()
	}
	return
}

func (f *Feedback) BeforeCreate(tx *gorm.DB) (err error) {
	if f.ID == "" {
		f.ID = uuid.New().String()
	}
	return
}

func (p *ProtocolloCardio) BeforeCreate(tx *gorm.DB) (err error) {
	if p.ID == "" {
		p.ID = uuid.New().String()
	}
	return
}

func (s *SessioneCardio) BeforeCreate(tx *gorm.DB) (err error) {
	if s.ID == "" {
		s.ID = uuid.New().String()
	}
	return
}

func (s *SchedaAllenamento) BeforeCreate(tx *gorm.DB) (err error) {
	if s.ID == "" {
		s.ID = uuid.New().String()
	}
	return
}

func (g *GiornoAllenamento) BeforeCreate(tx *gorm.DB) (err error) {
	if g.ID == "" {
		g.ID = uuid.New().String()
	}
	return
}

func (e *EsercizioScheda) BeforeCreate(tx *gorm.DB) (err error) {
	if e.ID == "" {
		e.ID = uuid.New().String()
	}
	return
}
