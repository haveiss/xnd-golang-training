package main

import (
	"encoding/json"
	"fmt"
)

type Article struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func main() {
	articleJson := `{"title":"Golang Training - Part 2","content":"What next after hello world?"}`

	var article Article
	json.Unmarshal([]byte(articleJson), &article)

	fmt.Printf("Title: %s\nContent: %s\n", article.Title, article.Content)
}
