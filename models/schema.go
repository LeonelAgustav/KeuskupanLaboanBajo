package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// BaseSyncFields: kolom wajib local-first
// - ID UUID string biar tidak tabrakan saat offline
// - DeletedAt soft delete untuk sync hapus
// - UpdatedAt untuk conflict resolution (last-write-wins)

type User struct {
	ID           string         `gorm:"primaryKey;size:36"`
	Nama         string         `gorm:"not null;size:255"`
	Email        string         `gorm:"uniqueIndex;not null;size:255"`
	PasswordHash string         `gorm:"not null;size:255"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	return nil
}

type Keuskupan struct {
	ID        string         `gorm:"primaryKey;size:36"`
	Nama      string         `gorm:"not null"`
	Alamat    string         `gorm:"not null"`
	Paroki    []Paroki       `gorm:"foreignKey:KeuskupanID"`
	Jurnal    []Jurnal       `gorm:"foreignKey:KeuskupanID"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (k *Keuskupan) BeforeCreate(tx *gorm.DB) error {
	if k.ID == "" {
		k.ID = uuid.New().String()
	}
	return nil
}

type Paroki struct {
	ID          string         `gorm:"primaryKey;size:36"`
	Nama        string         `gorm:"not null"`
	Alamat      string         `gorm:"not null"`
	KeuskupanID string         `gorm:"index;not null;size:36"`
	Keuskupan   Keuskupan      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;foreignKey:KeuskupanID;references:ID"`
	DetilJurnal []DetilJurnal  `gorm:"foreignKey:ParokiID"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func (p *Paroki) BeforeCreate(tx *gorm.DB) error {
	if p.ID == "" {
		p.ID = uuid.New().String()
	}
	return nil
}

type Jenis struct {
	ID        string         `gorm:"primaryKey;size:36"`
	Nama      string         `gorm:"uniqueIndex;not null;size:100"`
	Akun      []Akun         `gorm:"foreignKey:JenisID"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (j *Jenis) BeforeCreate(tx *gorm.DB) error {
	if j.ID == "" {
		j.ID = uuid.New().String()
	}
	return nil
}

type Akun struct {
	ID          string         `gorm:"primaryKey;size:36"`
	Kode        string         `gorm:"uniqueIndex;not null;size:20"`
	Nama        string         `gorm:"not null;size:255"`
	JenisID     string         `gorm:"index;not null;size:36"`
	Jenis       Jenis          `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;foreignKey:JenisID;references:ID"`
	DetilJurnal []DetilJurnal  `gorm:"foreignKey:AkunID"`
	Pembatasan  []Pembatasan   `gorm:"foreignKey:AkunID"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func (a *Akun) BeforeCreate(tx *gorm.DB) error {
	if a.ID == "" {
		a.ID = uuid.New().String()
	}
	return nil
}

type Jurnal struct {
	ID          string         `gorm:"primaryKey;size:36"`
	KeuskupanID string         `gorm:"index;not null;size:36"`
	Keuskupan   Keuskupan      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;foreignKey:KeuskupanID;references:ID"`
	DetilJurnal []DetilJurnal  `gorm:"foreignKey:JurnalID"`
	Tanggal     time.Time      `gorm:"index;not null"`
	Deskripsi   *string        `gorm:"size:500"`
	NoBukti     *string        `gorm:"uniqueIndex;size:50"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func (j *Jurnal) BeforeCreate(tx *gorm.DB) error {
	if j.ID == "" {
		j.ID = uuid.New().String()
	}
	if j.Tanggal.IsZero() {
		j.Tanggal = time.Now()
	}
	return nil
}

type DetilJurnal struct {
	ID         string         `gorm:"primaryKey;size:36"`
	JurnalID   string         `gorm:"index;not null;size:36"`
	Jurnal     Jurnal         `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;foreignKey:JurnalID;references:ID"`
	AkunID     string         `gorm:"index;not null;size:36"`
	Akun       Akun           `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;foreignKey:AkunID;references:ID"`
	ParokiID   string         `gorm:"index;not null;size:36"`
	Paroki     Paroki         `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;foreignKey:ParokiID;references:ID"`
	Debit      *float64       `gorm:"type:decimal(15,2);default:0"`
	Kredit     *float64       `gorm:"type:decimal(15,2);default:0"`
	Keterangan *string
	CreatedAt  time.Time      `gorm:"index"`
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

func (d *DetilJurnal) BeforeCreate(tx *gorm.DB) error {
	if d.ID == "" {
		d.ID = uuid.New().String()
	}
	return nil
}

type Pembatasan struct {
	ID        string         `gorm:"primaryKey;size:36"`
	Tipe      string         `gorm:"not null"`
	Nilai     *float64       `gorm:"type:decimal(15,2)"`
	AkunID    *string        `gorm:"size:36"`
	Akun      *Akun          `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:AkunID;references:ID"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (p *Pembatasan) BeforeCreate(tx *gorm.DB) error {
	if p.ID == "" {
		p.ID = uuid.New().String()
	}
	return nil
}
