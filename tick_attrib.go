package ibapi

import "fmt"

// TickAttrib .
type TickAttrib struct {
	CanAutoExecute bool
	PastLimit      bool
	PreOpen        bool
}

func NewTickAttrib() TickAttrib {
	return TickAttrib{}
}

func (t TickAttrib) String() string {
	return fmt.Sprintf("CanAutoExecute: %t, PastLimit: %t, PreOpen: %t", t.CanAutoExecute, t.PastLimit, t.PreOpen)
}
