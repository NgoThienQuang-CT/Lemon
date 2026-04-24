package token

import (
	"strconv"
)

type Position struct {
	Filename string
	Offset   int
	Line     int
	Column   int
}

func (pos *Position) IsValid() bool { return pos.Line > 0 }

func (pos *Position) String() string {
	return pos.Filename + ":" + strconv.Itoa(pos.Line) + ":" + strconv.Itoa(pos.Column)
}
