package ibapi

import "fmt"

// HistoricalSession .
type HistoricalSession struct {
	StartDateTime string
	EndDateTime   string
	RefDate       string
}

func NewHistoricalSession() HistoricalSession {
	return HistoricalSession{}
}

func (h HistoricalSession) String() string {
	return fmt.Sprintf("Start: %s, End: %s, Ref Date: %s", h.StartDateTime, h.EndDateTime, h.RefDate)
}
