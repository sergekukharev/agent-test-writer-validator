package domain

import (
	"errors"
	"fmt"
	"strings"
)

// ISBN represents a validated ISBN-13 identifier.
type ISBN struct {
	value string
}

func NewISBN(raw string) (ISBN, error) {
	cleaned := strings.ReplaceAll(raw, "-", "")
	cleaned = strings.ReplaceAll(cleaned, " ", "")

	if len(cleaned) != 13 {
		return ISBN{}, fmt.Errorf("ISBN must be 13 digits, got %d", len(cleaned))
	}

	for _, c := range cleaned {
		if c < '0' || c > '9' {
			return ISBN{}, errors.New("ISBN must contain only digits")
		}
	}

	if !validISBN13Checksum(cleaned) {
		return ISBN{}, errors.New("invalid ISBN-13 checksum")
	}

	return ISBN{value: cleaned}, nil
}

func (i ISBN) String() string { return i.value }

// Formatted returns the ISBN in grouped format: 978-X-XXXX-XXXX-X.
func (i ISBN) Formatted() string {
	if len(i.value) != 13 {
		return i.value
	}
	return fmt.Sprintf("%s-%s-%s-%s-%s",
		i.value[0:3],
		i.value[3:4],
		i.value[4:8],
		i.value[8:12],
		i.value[12:13],
	)
}

func validISBN13Checksum(digits string) bool {
	var sum int
	for i, c := range digits {
		d := int(c - '0')
		if i%2 == 0 {
			sum += d
		} else {
			sum += d * 3
		}
	}
	return sum%10 == 0
}
