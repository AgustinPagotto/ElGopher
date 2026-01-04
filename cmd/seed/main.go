package main

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL not set")
	}
	email := os.Getenv("ADMIN_EMAIL")
	name := os.Getenv("ADMIN_NAME")
	password := os.Getenv("ADMIN_PASSWORD")
	if email == "" || name == "" || password == "" {
		log.Fatal("ENV variables missing")
	}
	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		log.Fatal("Unable to connect to database: ", err)
	}
	defer conn.Close(context.Background())
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		log.Fatal(err)
	}
	sqlQuery := `INSERT INTO users (name, email, hashed_password) 
		VALUES ($1,$2,$3)`
	_, err = conn.Exec(context.Background(), sqlQuery, name, email, string(hashedPassword))
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Admin user seeded successfully")
}
