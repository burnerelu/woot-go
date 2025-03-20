package woot

import (
	"errors"
)

var (
	ErrInvalidIndex = errors.New("invalid insert index")
	ErrInvalidWChar = errors.New("invalid wcharacter")
	ErrNotFound     = errors.New("element not found")
)

// Represents the ID of a WCharacter
type ID struct {
	site  int
	clock int
}

// Represents a character - or WCharacter, as defined
type WCharacter struct {
	id ID
	v  bool
	c  rune
	cp ID
	cn ID
}

// Represents a WString, which may correspond to the contents of a file
type WString []WCharacter

// Represents the runtime structure for a file
type WOOT struct {
	site  int     // site ID
	clock int     // local clock
	str   WString //
}

func (w *WOOT) Insert(wchar WCharacter) error {
	err := w.str.Insert(wchar)

	if err == nil {
		w.clock += 1
	}

	return err
}

// Inserts at a specific position in the WString. TODO might need to change to insert in between two characters
func (w *WOOT) InsertAt(index int, value rune) error {
	err := w.str.InsertAt(index, WCharacter{id: ID{w.site, w.clock}, v: true, c: value})

	if err == nil {
		w.clock += 1
	}

	return err

}

// Textual representation of a WString
func (s WString) Text() []rune {
	// Skip the start and end markers
	content := make([]rune, 0, len(s)-2)
	for i := 1; i < len(s)-1; i++ {
		if s[i].v { // Only include visible characters
			content = append(content, s[i].c)
		}
	}
	return content
}

// Checks if a WString contains a rune
func (s *WString) Contains(value rune) bool {
	for _, v := range *s {
		if v.c == value {
			return true
		}
	}
	return false
}

// Finds ith visible character
func (s *WString) IthVisible(index int) (rune, error) {
	if index < 0 || index >= len(*s)-2 {
		return 0, ErrInvalidIndex
	}

	c := 0
	for _, v := range *s {
		if v.v {
			c += 1
			if c == index {
				return v.c, nil
			}
		}

	}

	return 0, ErrNotFound
}

func (s *WString) Insert(value WCharacter) error {
	for idx, v := range *s {
		if v.id == value.cp && (*s)[idx+1].id == value.cn {

			*s = append(*s, WCharacter{})
			copy((*s)[idx+2:], (*s)[idx+1:])

			(*s)[idx+1] = value

			return nil
		}
	}

	return ErrInvalidWChar
}

func (s *WString) InsertAt(index int, value WCharacter) error {
	// Take into account the string start and end characters
	if index < 0 || index > len(*s)-2 {
		return ErrInvalidIndex
	}

	// First character is the string start character
	insertPos := index + 1

	// Extend the slice with a dummy character
	*s = append(*s, WCharacter{})
	// Copy all the elements from the insert position to the right
	copy((*s)[insertPos+1:], (*s)[insertPos:])

	(*s)[insertPos] = value
	// Set the prev and next ids
	(*s)[insertPos].cp = (*s)[insertPos-1].id
	(*s)[insertPos].cn = (*s)[insertPos+1].id

	return nil
}

func (w *WOOT) Init(site int) {
	w.site = site
	w.clock = 0
	w.str = []WCharacter{
		{id: ID{site: site, clock: -2}, v: false, c: ' ', cp: ID{site: site, clock: -3}, cn: ID{site: site, clock: -1}},
		{id: ID{site: site, clock: -1}, v: false, c: ' ', cp: ID{site: site, clock: -2}, cn: ID{site: site, clock: -3}},
	}
}
