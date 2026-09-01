package controllers

import (
	"net/http"
	"time"

	"KeuskupanLaboanBajo_BE/models"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type SyncController struct {
	db *gorm.DB
}

func NewSyncController(db *gorm.DB) *SyncController {
	return &SyncController{db: db}
}

type PushItem struct {
	TargetTable string                 `json:"target_table" validate:"required"`
	RowID       string                 `json:"row_id" validate:"required"`
	Operation   string                 `json:"operation" validate:"required,oneof=insert update delete"`
	Payload     map[string]interface{} `json:"payload"`
}

type PushRequest struct {
	Items []PushItem `json:"items" validate:"required,min=1,dive"`
}

func (s *SyncController) Push(c echo.Context) error {
	var req PushRequest
	if err := c.Bind(&req); err != nil {
		return err
	}
	if err := c.Validate(&req); err != nil {
		return err
	}

	err := s.db.Transaction(func(tx *gorm.DB) error {
		for _, it := range req.Items {
			switch it.TargetTable {
			case "jurnals":
				if it.Operation == "delete" {
					tx.Where("id = ?", it.RowID).Delete(&models.Jurnal{})
					continue
				}
				// skip payload kosong (old queue) biar tidak FK fail
				keuskupanID, _ := it.Payload["keuskupan_id"].(string)
				if keuskupanID == "" {
					continue
				}
				var j models.Jurnal
				if err := tx.Where("id = ?", it.RowID).First(&j).Error; err != nil {
					j.ID = it.RowID
					j.KeuskupanID = keuskupanID
					if v, ok := it.Payload["tanggal"].(string); ok {
						// ponytail: Dart toIso8601String format 2024-08-31T13:35:55.123456 tanpa Z, coba beberapa format
						parsed := false
						for _, layout := range []string{time.RFC3339Nano, time.RFC3339, "2006-01-02T15:04:05.999999", "2006-01-02T15:04:05", "2006-01-02"} {
							if t, err := time.Parse(layout, v); err == nil {
								j.Tanggal = t
								parsed = true
								break
							}
						}
						if !parsed {
							j.Tanggal = time.Now()
						}
					} else {
						j.Tanggal = time.Now()
					}
					if v, ok := it.Payload["deskripsi"].(string); ok {
						j.Deskripsi = &v
					}
					if v, ok := it.Payload["no_bukti"].(string); ok && v != "" {
						j.NoBukti = &v
					}
					if err := tx.Create(&j).Error; err != nil {
						// skip FK fail, jangan abort bulk
						continue
					}
				} else {
					updates := map[string]interface{}{"updated_at": time.Now()}
					if v, ok := it.Payload["deskripsi"].(string); ok {
						updates["deskripsi"] = v
					}
					if v, ok := it.Payload["no_bukti"].(string); ok {
						updates["no_bukti"] = v
					}
					tx.Model(&models.Jurnal{}).Where("id = ?", it.RowID).Updates(updates)
				}
			case "keuskupans":
				if it.Operation == "delete" {
					tx.Where("id = ?", it.RowID).Delete(&models.Keuskupan{})
					continue
				}
				nama, _ := it.Payload["nama"].(string)
				alamat, _ := it.Payload["alamat"].(string)
				if nama == "" {
					continue
				}
				var k models.Keuskupan
				if err := tx.Where("id = ?", it.RowID).First(&k).Error; err != nil {
					k.ID = it.RowID
					k.Nama = nama
					k.Alamat = alamat
					tx.Create(&k)
				} else {
					tx.Model(&models.Keuskupan{}).Where("id = ?", it.RowID).Updates(map[string]interface{}{"nama": nama, "alamat": alamat})
				}
			case "parokis":
				if it.Operation == "delete" {
					tx.Where("id = ?", it.RowID).Delete(&models.Paroki{})
					continue
				}
				nama, _ := it.Payload["nama"].(string)
				alamat, _ := it.Payload["alamat"].(string)
				keuskupanID, _ := it.Payload["keuskupan_id"].(string)
				if nama == "" || keuskupanID == "" {
					continue
				}
				var p models.Paroki
				if err := tx.Where("id = ?", it.RowID).First(&p).Error; err != nil {
					p.ID = it.RowID
					p.Nama = nama
					p.Alamat = alamat
					p.KeuskupanID = keuskupanID
					tx.Create(&p)
				} else {
					tx.Model(&models.Paroki{}).Where("id = ?", it.RowID).Updates(map[string]interface{}{"nama": nama, "alamat": alamat, "keuskupan_id": keuskupanID})
				}
			case "jenis":
				if it.Operation == "delete" {
					tx.Where("id = ?", it.RowID).Delete(&models.Jenis{})
					continue
				}
				nama, _ := it.Payload["nama"].(string)
				if nama == "" {
					continue
				}
				var j models.Jenis
				if err := tx.Where("id = ?", it.RowID).First(&j).Error; err != nil {
					j.ID = it.RowID
					j.Nama = nama
					tx.Create(&j)
				} else {
					tx.Model(&models.Jenis{}).Where("id = ?", it.RowID).Updates(map[string]interface{}{"nama": nama})
				}
			case "akuns":
				if it.Operation == "delete" {
					tx.Where("id = ?", it.RowID).Delete(&models.Akun{})
					continue
				}
				kode, _ := it.Payload["kode"].(string)
				nama, _ := it.Payload["nama"].(string)
				jenisID, _ := it.Payload["jenis_id"].(string)
				if kode == "" || jenisID == "" {
					continue
				}
				var a models.Akun
				if err := tx.Where("id = ?", it.RowID).First(&a).Error; err != nil {
					a.ID = it.RowID
					a.Kode = kode
					a.Nama = nama
					a.JenisID = jenisID
					tx.Create(&a)
				} else {
					tx.Model(&models.Akun{}).Where("id = ?", it.RowID).Updates(map[string]interface{}{"kode": kode, "nama": nama, "jenis_id": jenisID})
				}
			case "detil_jurnals":
				if it.Operation == "delete" {
					tx.Where("id = ?", it.RowID).Delete(&models.DetilJurnal{})
					continue
				}
				jurnalID, _ := it.Payload["jurnal_id"].(string)
				akunID, _ := it.Payload["akun_id"].(string)
				parokiID, _ := it.Payload["paroki_id"].(string)
				if jurnalID == "" || akunID == "" || parokiID == "" {
					continue
				}
				var d models.DetilJurnal
				if err := tx.Where("id = ?", it.RowID).First(&d).Error; err != nil {
					d.ID = it.RowID
					d.JurnalID = jurnalID
					d.AkunID = akunID
					d.ParokiID = parokiID
					if v, ok := it.Payload["debit"].(float64); ok {
						d.Debit = &v
					}
					if v, ok := it.Payload["kredit"].(float64); ok {
						d.Kredit = &v
					}
					if v, ok := it.Payload["keterangan"].(string); ok {
						d.Keterangan = &v
					}
					tx.Create(&d)
				} else {
					updates := map[string]interface{}{"updated_at": time.Now()}
					if v, ok := it.Payload["debit"].(float64); ok {
						updates["debit"] = v
					}
					if v, ok := it.Payload["kredit"].(float64); ok {
						updates["kredit"] = v
					}
					tx.Model(&models.DetilJurnal{}).Where("id = ?", it.RowID).Updates(updates)
				}
			default:
				// ignore unknown
			}
		}
		return nil
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Gagal sync push")
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"message": "Sync push berhasil", "count": len(req.Items)})
}

func (s *SyncController) Pull(c echo.Context) error {
	sinceStr := c.QueryParam("since")
	var since time.Time
	if sinceStr != "" {
		if t, err := time.Parse(time.RFC3339, sinceStr); err == nil {
			since = t
		} else if t, err := time.Parse("2006-01-02", sinceStr); err == nil {
			since = t
		}
	}
	var jurnals []models.Jurnal
	q := s.db.Preload("DetilJurnal").Where("updated_at > ?", since).Order("updated_at ASC").Limit(100)
	if err := q.Find(&jurnals).Error; err != nil {
		return err
	}
	var deleted []models.Jurnal
	s.db.Unscoped().Where("deleted_at IS NOT NULL AND deleted_at > ?", since).Find(&deleted)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"since":   since,
		"jurnals": jurnals,
		"deleted": deleted,
		"count":   len(jurnals),
	})
}
