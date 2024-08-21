package cli

import (
	"fmt"
	"gobooks/internal/service"
	"os"
	"strconv"
	"time"
)

type BookCLI struct {
	service *service.BookService	
}

func NewBookCLI(service *service.BookService) *BookCLI {
	return &BookCLI{service: service}
}

func (c *BookCLI) Run() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: books <command> [arguments]")
		return
	}

	command := os.Args[1]
	
	switch command {
	case "search":
		if len(os.Args) < 3 {
			fmt.Println("Usage: books search <query>")
			return
		}
		bookName := os.Args[2]
		c.searchBooks(bookName)
	case "simulate":
		if len (os.Args) < 3 {
			fmt.Println("Usage: books simulate <id>")
			return
		}
		booksID := os.Args[2:]
		c.simulateReading(booksID)
	}
}

func (c *BookCLI) searchBooks(name string) {
	books, err := c.service.SearchBooksByName(name)
	if err != nil {
		fmt.Println("failed to get books:", err)
		return
	}

	if len(books) == 0 {
		fmt.Println("no books found")
	}
	
	fmt.Printf("%d books found:\n", len(books))
	for _, book := range books {
		fmt.Printf("ID: %d, Title: %s, Author: %s, Genre: %s\n", book.ID, book.Title, book.Author, book.Genre)
	}
}

func (c *BookCLI) simulateReading(bookIDsStr []string) {
	var bookIDs []int
	for _, idStr := range bookIDsStr {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			fmt.Println("invalid book ID:", idStr)
			return
		}
		bookIDs = append(bookIDs, id)
	}
	responses := c.service.SimulateMultiplesReading(bookIDs, 2*time.Second)
	for _, response := range responses {
		fmt.Println(response)
	}
}