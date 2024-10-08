package service

import (
	"database/sql"
	"fmt"
	"time"
)

type Book struct {
	ID int
	Title string
	Author string
	Genre string
}

type BookService struct {
	db *sql.DB
}

func NewBookService(db *sql.DB) *BookService {
	return &BookService{db: db}
}

func (s *BookService) Create(b *Book) error {
	result, err := s.db.Exec("INSERT INTO books (title, author, genre) VALUES ($1, $2, $3)", b.Title, b.Author, b.Genre)
	if err != nil {
		return err
	}
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return err
	}
	b.ID = int(lastInsertID)
	return nil
}

func (s *BookService) GetBooks() ([]Book, error) {
	row, err := s.db.Query("SELECT * FROM books")
	if err != nil {
		return nil, err
	}
	var books []Book
	for row.Next() {
		var book Book
		err := row.Scan(&book.ID, &book.Title, &book.Author, &book.Genre)
		if err != nil {
			return nil, err
		}
		books = append(books, book)
	}
	return books, nil
}

func (s *BookService) GetBookId(id int) (*Book, error) {
	row := s.db.QueryRow(`select * from books where id = ?`, id)
	var book Book
	err := row.Scan(&book.ID, &book.Title, &book.Author, &book.Genre)
	if err != nil {
		return nil, err
	}
	return &book, nil
}

func (s *BookService) UpdateBook(b *Book) error {
	_, err := s.db.Exec("UPDATE books SET title = ?, author = ?, genre = ? WHERE id = ?", b.Title, b.Author, b.Genre, b.ID)
	if err != nil {
		return err
	}
	return nil
}

func (s *BookService) DeleteBook(id int) error {
	_, err := s.db.Exec("DELETE FROM books WHERE id = ?", id)
	if err != nil {
		return err
	}
	return nil
}

func (s *BookService) SimulateReading(bookID int, duration time.Duration, results chan<-string) {
	book, err := s.GetBookId(bookID)
	if err != nil || book == nil {
		results <- fmt.Sprintf("Book with ID %d not found", bookID)
		return
	}
	time.Sleep(duration)
	results <- fmt.Sprintf("Book %s read", book.Title)
}

func (s *BookService) SimulateMultiplesReading(bookIDs []int, duration time.Duration) []string {
	results := make(chan string, len(bookIDs))
	for _, id := range bookIDs {
		go func(bookID int) {
			s.SimulateReading(bookID, duration, results)
		}(id)
	}
	var responses []string
	for range(bookIDs) {
		responses = append(responses, <-results)
	}
	close(results)
	return responses
}

func (s *BookService) SearchBooksByName(name string) ([]Book, error) {
	row, err := s.db.Query("SELECT * FROM books WHERE title LIKE $1", "%"+name+"%")
	if err != nil {
		return nil, err
	}
	var books []Book
	for row.Next() {
		var book Book
		err := row.Scan(&book.ID, &book.Title, &book.Author, &book.Genre)
		if err != nil {
			return nil, err
		}
		books = append(books, book)
	}
	return books, nil
}