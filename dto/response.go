package dto

import (
	"time"
)

type AkunResponse struct {
	ID string           `json:"id"`
	Kode      string         `json:"kode"`
	Nama      string         `json:"nama"`
	JenisID string           `json:"jenis_id"`
	Jenis     *JenisResponse `json:"jenis,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

type DetilJurnalResp struct {
	ID string            `json:"id"`
	JurnalID   string          `json:"jurnal_id"`
	AkunID string            `json:"akun_id"`
	Akun       *AkunResponse   `json:"akun,omitempty"`
	ParokiID string            `json:"paroki_id"`
	Paroki     *ParokiResponse `json:"paroki,omitempty"`
	Debit      *float64        `json:"debit,omitempty"`
	Kredit     *float64        `json:"kredit,omitempty"`
	Keterangan *string         `json:"keterangan,omitempty"`
	CreatedAt  time.Time       `json:"created_at"`
	UpdatedAt  time.Time       `json:"updated_at"`
}

type JenisResponse struct {
	ID string      `json:"id"`
	Nama      string    `json:"nama"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type JurnalResponse struct {
	ID          string             `json:"id"`
	KeuskupanID string               `json:"keuskupan_id"`
	Keuskupan   *KeuskupanResponse `json:"keuskupan,omitempty"`
	Tanggal     time.Time          `json:"tanggal"`
	Deskripsi   *string            `json:"deskripsi,omitempty"`
	NoBukti     *string            `json:"no_bukti,omitempty"`
	DetilJurnal []DetilJurnalResp  `json:"detil_jurnal,omitempty"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
}

type KeuskupanResponse struct {
	ID string      `json:"id"`
	Nama      string    `json:"nama"`
	Alamat    string    `json:"alamat"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ParokiResponse struct {
	ID string               `json:"id"`
	Nama        string             `json:"nama"`
	Alamat      string             `json:"alamat"`
	KeuskupanID string               `json:"keuskupan_id"`
	Keuskupan   *KeuskupanResponse `json:"keuskupan,omitempty"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
}

type PembatasanResponse struct {
	ID string          `json:"id"`
	Tipe      string        `json:"tipe"`
	Nilai     *float64      `json:"nilai,omitempty"`
	AkunID    *string         `json:"akun_id,omitempty"`
	Akun      *AkunResponse `json:"akun,omitempty"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}

type UserResponse struct {
	ID string      `json:"id"`
	Nama      string    `json:"nama"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ListResponse struct {
	Data       interface{} `json:"data"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	Total      int64       `json:"total"`
	TotalPages int         `json:"total_pages"`
}

type ErrorResponse struct {
	Code    string        `json:"code"`
	Message string        `json:"message"`
	Details []ErrorDetail `json:"details,omitempty"`
}

type ErrorDetail struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type CreateResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
