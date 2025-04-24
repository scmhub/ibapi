package ibapi

import "fmt"

// WshEventData .
type WshEventData struct {
	ConID           int64 // UNSET_INT
	Filter          string
	FillWatchList   bool
	FillPortfolio   bool
	FillCompetitors bool
	StartDate       string
	EndDate         string
	TotalLimit      int64 // UNSET_INT
}

func NewWshEventData() WshEventData {
	wed := WshEventData{}
	wed.ConID = UNSET_INT
	wed.TotalLimit = UNSET_INT
	return wed
}

func (w WshEventData) String() string {
	return fmt.Sprintf("WshEventData. ConId: %s, Filter: %s, Fill Watchlist: %t, Fill Portfolio: %t, Fill Competitors: %t",
		IntMaxString(w.ConID), w.Filter, w.FillWatchList, w.FillPortfolio, w.FillCompetitors)
}
