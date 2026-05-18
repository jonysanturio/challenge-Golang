package config

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/joho/godotenv"
)

var DB *gorm.DB
// Inicializa init para configurar la conexion a la base de datos
func Init(){
	if err := godotenv.Load(); err != nil{
		log.Println("No .env file found, relying on environment variables", err)
	}

	// Conexion de la base de datos
	dsn := os.Getenv("DATABASE_URL")
	if dsn == ""{
		// Construct DSN
		host := os.Getenv("DB_HOST")
		if host == "" {
			host = "localhost"
		}
		port := os.Getenv("DB_PORT")
		if port == "" {
			port = "5432"
		}
		user := os.Getenv("DB_USER")
		if user == "" {
			user = "postgres"
		}
		password := os.Getenv("DB_PASSWORD")
		if password == "" {
			password = "[PASSWORD]"
		}
		dbname := os.Getenv("DB_NAME")
		if dbname == "" {
			dbname = "product_api"
		}
		sslmode := os.Getenv("DB_SSLMODE")
		if sslmode == "" {
			sslmode = "disable"
		}

		dsn = "host=" + host + " port=" + port + " user=" + user + " password=" + password + " dbname=" + dbname + " sslmode=" + sslmode
	}

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}
}
// GetDB retorna la instancia de la base de datos
func GetDB() *gorm.DB {
	return DB
}