package domain

import (
	"errors"
	"fmt"
	"time"
)

type Book struct {
	isbn        ISBN
	title       string
	author      Author
	price       Money
	publishedAt time.Time
	genre       Genre
}

type Genre string

const (
	GenreFiction    Genre = "fiction"
	GenreNonFiction Genre = "non-fiction"
	GenreScience    Genre = "science"
	GenreBiography  Genre = "biography"
	GenreChildren   Genre = "children"
)

func NewBook(isbn ISBN, title string, author Author, price Money, publishedAt time.Time, genre Genre) (Book, error) {
	if title == "" {
		return Book{}, errors.New("title must not be empty")
	}
	if !isValidGenre(genre) {
		return Book{}, fmt.Errorf("unknown genre: %s", genre)
	}
	if price.Amount() < 0 {
		return Book{}, errors.New("price must not be negative")
	}
	return Book{
		isbn:        isbn,
		title:       title,
		author:      author,
		price:       price,
		publishedAt: publishedAt,
		genre:       genre,
	}, nil
}

func (b Book) ISBN() ISBN       { return b.isbn }
func (b Book) Title() string    { return b.title }
func (b Book) Author() Author   { return b.author }
func (b Book) Price() Money     { return b.price }
func (b Book) Genre() Genre     { return b.genre }
func (b Book) PublishedAt() time.Time { return b.publishedAt }

// IsClassic returns true if the book was published more than 50 years ago.
func (b Book) IsClassic() bool {
	return time.Since(b.publishedAt) > 50*365*24*time.Hour
}

// IsRecent returns true if the book was published within the last year.
func (b Book) IsRecent() bool {
	return time.Since(b.publishedAt) < 365*24*time.Hour
}

func isValidGenre(g Genre) bool {
	switch g {
	case GenreFiction, GenreNonFiction, GenreScience, GenreBiography, GenreChildren:
		return true
	default:
		return false
	}
}
