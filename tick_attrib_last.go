package ibapi

import "fmt"

// TickAttribLast .
type TickAttribLast struct {
	PastLimit  bool
	Unreported bool
}

func NewTickAttribLast() TickAttribLast {
	return TickAttribLast{}
}

func (t TickAttribLast) String() string {
	return fmt.Sprintf("PastLimit: %t, Unreported: %t", t.PastLimit, t.Unreported)
}
