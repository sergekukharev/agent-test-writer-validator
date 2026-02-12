package domain

import (
	"errors"
	"fmt"
	"strings"
)

type Author struct {
	firstName string
	lastName  string
}

func NewAuthor(firstName, lastName string) (Author, error) {
	firstName = strings.TrimSpace(firstName)
	lastName = strings.TrimSpace(lastName)

	if firstName == "" {
		return Author{}, errors.New("first name must not be empty")
	}
	if lastName == "" {
		return Author{}, errors.New("last name must not be empty")
	}
	return Author{firstName: firstName, lastName: lastName}, nil
}

func (a Author) FirstName() string { return a.firstName }
func (a Author) LastName() string  { return a.lastName }

func (a Author) FullName() string {
	return fmt.Sprintf("%s %s", a.firstName, a.lastName)
}

// Initials returns the author's initials, e.g. "J.K." for "Joanne Kathleen".
func (a Author) Initials() string {
	parts := strings.Fields(a.firstName)
	var initials []string
	for _, p := range parts {
		if len(p) > 0 {
			initials = append(initials, string(p[0])+".")
		}
	}
	return strings.Join(initials, "")
}
