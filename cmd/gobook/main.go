package main

import (
	"database/sql"
	_ "modernc.org/sqlite"
	"gobooks/internal/service"
	"gobooks/internal/web"
	"net/http"
)

func main() {
	db, err := sql.Open("sqlite", "./books.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	bookService := service.NewBookService(db)
	bookHandlers := web.NewBookHandlers(bookService)

	router := http.NewServeMux()
	router.HandleFunc("GET /books", bookHandlers.GetBooks)
	router.HandleFunc("GET /books/{id}", bookHandlers.GetBook)
	router.HandleFunc("POST /books", bookHandlers.CreateBook)
	router.HandleFunc("PUT /books/{id}", bookHandlers.UpdateBook)
	router.HandleFunc("DELETE /books/{id}", bookHandlers.DeleteBook)
	
	http.ListenAndServe(":777", router)
}