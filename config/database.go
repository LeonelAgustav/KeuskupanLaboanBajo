package config

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"

	"KeuskupanLaboanBajo_BE/models"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func ensureDatabaseExists(user, pass, host, port, name string) {
	// konek tanpa nama DB untuk CREATE DATABASE IF NOT EXISTS
	dsnRoot := fmt.Sprintf("%s:%s@tcp(%s:%s)/?charset=utf8mb4&parseTime=True&loc=Local",
		user, pass, host, port)
	dbRoot, err := gorm.Open(mysql.Open(dsnRoot), &gorm.Config{})
	if err != nil {
		log.Fatalf("Gagal konek MySQL server untuk create DB: %v", err)
	}
	// ponytail: raw SQL paling simpel, tidak perlu driver tambahan
	sqlDB, _ := dbRoot.DB()
	defer sqlDB.Close()
	createSQL := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci", name)
	if err := dbRoot.Exec(createSQL).Error; err != nil {
		log.Fatalf("Gagal create database %s: %v", name, err)
	}
	log.Printf("Database `%s` siap (created if not exists)", name)
}

func parseDBURL(raw string) (user, pass, host, port, name string, ok bool) {
	if raw == "" {
		return "", "", "", "", "", false
	}
	u, err := url.Parse(raw)
	if err != nil {
		return "", "", "", "", "", false
	}
	if u.Scheme != "mysql" {
		return "", "", "", "", "", false
	}
	user = u.User.Username()
	pass, _ = u.User.Password()
	host = u.Hostname()
	port = u.Port()
	if port == "" {
		port = "3306"
	}
	name = strings.TrimPrefix(u.Path, "/")
	if user == "" || host == "" || name == "" {
		return "", "", "", "", "", false
	}
	return user, pass, host, port, name, true
}

func ConnectDB() *gorm.DB {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env tidak ditemukan, pakai env sistem")
	}

	// ponytail: dukung DB_URL Railway mysql://user:pass@host:port/db
	var dsn string
	var user, pass, host, port, name string
	var fromURL bool

	if dbURL := os.Getenv("DB_URL"); dbURL != "" {
		var ok bool
		user, pass, host, port, name, ok = parseDBURL(dbURL)
		if !ok {
			log.Fatalf("DB_URL format salah, harus mysql://user:pass@host:port/dbname, got: %s", dbURL)
		}
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			user, pass, host, port, name)
		fromURL = true
		log.Printf("Pakai DB_URL -> %s:***@%s:%s/%s", user, host, port, name)
	} else {
		user = os.Getenv("DB_USER")
		pass = os.Getenv("DB_PASS")
		host = os.Getenv("DB_HOST")
		port = os.Getenv("DB_PORT")
		name = os.Getenv("DB_NAME")

		if user == "" || host == "" || name == "" {
			log.Fatal("DB_USER/DB_HOST/DB_NAME kosong, cek .env atau set DB_URL")
		}
		if port == "" {
			port = "3306"
		}
		// 1. Pastikan DB ada (auto create) - skip untuk Railway (DB sudah ada)
		ensureDatabaseExists(user, pass, host, port, name)
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			user, pass, host, port, name)
	}

	// Railway tidak perlu CREATE DATABASE, skip ensure
	if fromURL {
		log.Printf("Skip ensureDatabaseExists untuk DB_URL (Railway)")
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		log.Fatal("Gagal koneksi database:", err)
	}

	// 3. AutoMigrate untuk local-first (UUID + DeletedAt + UpdatedAt)
	if err := db.AutoMigrate(
		&models.User{},
		&models.Keuskupan{},
		&models.Paroki{},
		&models.Jenis{},
		&models.Akun{},
		&models.Jurnal{},
		&models.DetilJurnal{},
		&models.Pembatasan{},
	); err != nil {
		log.Fatal("Gagal AutoMigrate:", err)
	}

	log.Println("Koneksi database berhasil + AutoMigrate OK!")
	return db
}
