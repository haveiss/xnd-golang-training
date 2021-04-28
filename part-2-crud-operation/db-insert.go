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

	article, err := NewArticle("Awesome title", "With awesome content")

	res, err := db.ExecContext(context.Background(), "INSERT INTO article (title, content, created_at) VALUES($1 ,$2 ,$3)", article.Title, article.Content, article.CreatedAt)

	if err != nil {
		log.Fatal(err)
	}

	rowAffected, err := res.RowsAffected()

	if err != nil {
		log.Fatal(err)
	}

	if rowAffected != 1 {
		log.Fatalf("Expected row affected equal to 1 but got %d", rowAffected)
	}

	fmt.Println("Insert success")

}

func NewArticle(title string, content string) (Article, error) {
	if len(title) <= 0 {
		return Article{}, fmt.Errorf("title cannot be empty")
	}

	if len(title) > 100 {
		return Article{}, fmt.Errorf("title length maximum 100 chars")
	}

	return Article{
		Title:     title,
		Content:   content,
		CreatedAt: time.Now(),
	}, nil
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
