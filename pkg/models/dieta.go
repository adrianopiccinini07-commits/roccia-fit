package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Dieta struct {
	ID         string    `gorm:"primaryKey;type:varchar(36)" json:"ID"`
	IDAtleta   string    `gorm:"column:id_atleta" json:"IDAtleta"`
	Nome       string    `gorm:"size:100;not null" json:"Nome"`
	DataInizio time.Time `json:"DataInizio"`
	DataFine   time.Time `json:"DataFine"`
	Attiva     bool      `gorm:"default:false" json:"Attiva"`
	Acqua      string    `gorm:"size:150" json:"Acqua"`
	Sale       string    `gorm:"size:150" json:"Sale"`
	CreatedAt  time.Time `json:"CreatedAt"`

	GiorniDieta []GiornoDieta `gorm:"foreignKey:IDDieta;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"GiorniDieta"`
}

type GiornoDieta struct {
	ID          string `gorm:"primaryKey;type:varchar(36)" json:"ID"`
	IDDieta     string `gorm:"not null" json:"IDDieta"`
	TipoGiorno  string `gorm:"size:50;not null" json:"TipoGiorno"`
	TargetKcal  int    `gorm:"not null" json:"TargetKcal"`
	TargetPro   int    `gorm:"not null" json:"TargetPro"`
	TargetCarbo int    `gorm:"not null" json:"TargetCarbo"`
	TargetFats  int    `gorm:"not null" json:"TargetFats"`
	TargetFibre int    `gorm:"not null;default:0" json:"TargetFibre"` // ➔ Aggiunto per non perdere i dati delle fibre

	Pasti []Pasto `gorm:"foreignKey:IDGiornoDieta;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"Pasti"`
}

type Pasto struct {
	ID            string `gorm:"primaryKey;type:varchar(36)" json:"ID"`
	IDGiornoDieta string `gorm:"not null" json:"IDGiornoDieta"`
	NomePasto     string `gorm:"size:50;not null" json:"NomePasto"`
	Orario        string `gorm:"size:10" json:"Orario"`
	Ordine        int    `gorm:"not null" json:"Ordine"`

	Alimenti []AlimentoPasto `gorm:"foreignKey:IDPasto;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"Alimenti"`
}

type AlimentoPasto struct {
	ID       string  `gorm:"primaryKey;type:varchar(36)" json:"ID"`
	IDPasto  string  `gorm:"not null" json:"IDPasto"`
	NomeCibo string  `gorm:"size:100;not null" json:"NomeCibo"`
	Grammi   float64 `gorm:"not null" json:"Grammi"`
	Kcal     float64 `json:"Kcal"`
	Protein  float64 `json:"Protein"`
	Carbo    float64 `json:"Carbo"`
	Fats     float64 `json:"Fats"`

	Alternative []AlternativaAlimento `gorm:"foreignKey:IDAlimentoPasto;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"Alternative"`
}

type AlternativaAlimento struct {
	ID              string  `gorm:"primaryKey;type:varchar(36)" json:"ID"`
	IDAlimentoPasto string  `gorm:"not null" json:"IDAlimentoPasto"`
	NomeCibo        string  `gorm:"size:100;not null" json:"NomeCibo"`
	Grammi          float64 `gorm:"not null" json:"Grammi"`
}

func (d *Dieta) BeforeCreate(tx *gorm.DB) (err error) {
	if d.ID == "" {
		d.ID = uuid.New().String()
	}
	return
}

func (g *GiornoDieta) BeforeCreate(tx *gorm.DB) (err error) {
	if g.ID == "" {
		g.ID = uuid.New().String()
	}
	return
}

func (p *Pasto) BeforeCreate(tx *gorm.DB) (err error) {
	if p.ID == "" {
		p.ID = uuid.New().String()
	}
	return
}

func (a *AlimentoPasto) BeforeCreate(tx *gorm.DB) (err error) {
	if a.ID == "" {
		a.ID = uuid.New().String()
	}
	return
}

func (a *AlternativaAlimento) BeforeCreate(tx *gorm.DB) (err error) {
	if a.ID == "" {
		a.ID = uuid.New().String()
	}
	return
}
