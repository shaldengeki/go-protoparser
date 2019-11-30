package scanner

import (
	"unicode/utf8"

	"github.com/yoheimuta/go-protoparser/v4/parser/meta"
)

// Position represents a source position.
type Position struct {
	meta.Position

	// columns is a map which the key is a line number and the value is a column number.
	columns map[int]int
}

// NewPosition creates a new Position.
func NewPosition() *Position {
	return &Position{
		Position: meta.Position{
			Offset: 0,
			Line:   1,
			Column: 1,
		},
		columns: make(map[int]int),
	}
}

// String stringify the position.
func (pos Position) String() string {
	return pos.Position.String()
}

// Advance advances the position value.
func (pos *Position) Advance(r rune) {
	len := utf8.RuneLen(r)
	pos.Offset += len

	if r == '\n' {
		pos.columns[pos.Line] = pos.Column

		pos.Line++
		pos.Column = 1
	} else {
		pos.Column++
	}
}

// Revert reverts the position value.
func (pos *Position) Revert(r rune) {
	len := utf8.RuneLen(r)
	pos.Offset -= len

	if r == '\n' {
		pos.Line--
		pos.Column = pos.columns[pos.Line]
	} else {
		pos.Column--
	}
}
