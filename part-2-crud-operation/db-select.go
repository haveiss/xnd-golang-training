package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

type Article struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
}

func main() {
	db := initDB()

	rows, err := db.QueryContext(context.Background(), "SELECT id, title, content, created_at FROM article")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	articles := make([]Article, 0)

	for rows.Next() {
		var article Article
		if err := rows.Scan(
			&article.ID,
			&article.Title,
			&article.Content,
			&article.CreatedAt,
		); err != nil {
			log.Fatal(err)
		}
		articles = append(articles, article)
	}

	fmt.Println(articles)
}

func initDB() *sql.DB {
	dbHost := os.Getenv("DATABASE_HOST")
	dbPort := os.Getenv("DATABASE_PORT")
	dbUser := os.Getenv("DATABASE_USER")
	dbPass := os.Getenv("DATABASE_PASS")
	dbName := os.Getenv("DATABASE_NAME")

	pgDsn := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPass, dbName)

	db, err := sql.Open("postgres", pgDsn)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")
	return db
}
