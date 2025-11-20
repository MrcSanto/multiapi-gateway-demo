package db

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func ConnectDB() (*gorm.DB, error) {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	schema := os.Getenv("DB_SCHEMA_GO")

	log.Printf(
		"DB Config - Host: %s, Port: %s, User: %s, DBName: %s, Schema: %s",
		host, port, user, dbname, schema,
	)

	// Validar se as variáveis foram carregadas
	if host == "" || port == "" || user == "" || dbname == "" || schema == "" {
		return nil, fmt.Errorf("variáveis de ambiente do banco não configuradas corretamente")
	}

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable search_path=%s",
		host, port, user, password, dbname, schema,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatal("Erro ao conectar ao banco de dados:", err)
		return nil, err
	}

	log.Println("Conexão com banco de dados estabelecida!")
	return db, nil
}
