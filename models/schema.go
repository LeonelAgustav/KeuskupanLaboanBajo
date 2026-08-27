package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        uint   `gorm:"primaryKey"`
	Nama      string `gorm:"not null"`
	Email     string `gorm:"uniqueIndex;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Paroki struct {
	ID          uint          `gorm:"primaryKey"`
	Nama        string        `gorm:"not null"`
	Alamat      string        `gorm:"not null"`
	KeuskupanID uint          `gorm:"index;not null"`
	Keuskupan   Keuskupan     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;foreignKey:KeuskupanID;references:ID"`
	DetilJurnal []DetilJurnal `gorm:"foreignKey:ParokiID"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Keuskupan struct {
	ID        uint     `gorm:"primaryKey"`
	Nama      string   `gorm:"not null"`
	Alamat    string   `gorm:"not null"`
	Paroki    []Paroki `gorm:"foreignKey:KeuskupanID"`
	Jurnal    []Jurnal `gorm:"foreignKey:KeuskupanID"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Jenis struct {
	ID        uint   `gorm:"primaryKey"`
	Nama      string `gorm:"uniqueIndex;not null"`
	Akun      []Akun `gorm:"foreignKey:JenisID"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Akun struct {
	ID          uint          `gorm:"primaryKey"`
	Kode        string        `gorm:"uniqueIndex;not null"`
	Nama        string        `gorm:"not null"`
	JenisID     uint          `gorm:"not null"`
	Jenis       Jenis         `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;foreignKey:JenisID;references:ID"`
	DetilJurnal []DetilJurnal `gorm:"foreignKey:AkunID"`
	Pembatasan  []Pembatasan  `gorm:"foreignKey:AkunID"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Jurnal struct {
	ID          string        `gorm:"primaryKey;size:36"`
	KeuskupanID uint          `gorm:"index;not null"`
	Keuskupan   Keuskupan     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;foreignKey:KeuskupanID;references:ID"`
	DetilJurnal []DetilJurnal `gorm:"foreignKey:JurnalID"`
	Tanggal     time.Time     `gorm:"index;not null;default:CURRENT_TIMESTAMP"`
	Deskripsi   *string
	NoBukti     *string `gorm:"uniqueIndex"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (j *Jurnal) BeforeCreate(tx *gorm.DB) error {
	if j.ID == "" {
		j.ID = uuid.New().String()
	}
	return nil
}

type DetilJurnal struct {
	ID         uint     `gorm:"primaryKey"`
	JurnalID   string   `gorm:"index;not null;size:36"`
	Jurnal     Jurnal   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;foreignKey:JurnalID;references:ID"`
	AkunID     uint     `gorm:"index;not null"`
	Akun       Akun     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;foreignKey:AkunID;references:ID"`
	ParokiID   uint     `gorm:"index;not null"`
	Paroki     Paroki   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;foreignKey:ParokiID;references:ID"`
	Debit      *float64 `gorm:"type:decimal(15,2);default:0"`
	Kredit     *float64 `gorm:"type:decimal(15,2);default:0"`
	Keterangan *string
	CreatedAt  time.Time `gorm:"index"`
	UpdatedAt  time.Time
}

type Pembatasan struct {
	ID        uint     `gorm:"primaryKey"`
	Tipe      string   `gorm:"not null"`
	Nilai     *float64 `gorm:"type:decimal(15,2)"`
	AkunID    *uint
	Akun      *Akun `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:AkunID;references:ID"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
