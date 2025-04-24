package ibapi

import "fmt"

// TickAttribBidAsk .
type TickAttribBidAsk struct {
	BidPastLow  bool
	AskPastHigh bool
}

func NewTickAttribBidAsk() TickAttribBidAsk {
	return TickAttribBidAsk{}
}

func (t TickAttribBidAsk) String() string {
	return fmt.Sprintf("BidPastLow: %t, AskPastHigh: %t", t.BidPastLow, t.AskPastHigh)
}
