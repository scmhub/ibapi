package ibapi

import "fmt"

// HistoricalTick is the historical tick's description.
// Used when requesting historical tick data with whatToShow = MIDPOINT.
type HistoricalTick struct {
	Time  int64
	Price float64
	Size  Decimal
}

func NewHistoricalTick() HistoricalTick {
	ht := HistoricalTick{}
	ht.Size = UNSET_DECIMAL
	return ht
}

func (h HistoricalTick) String() string {
	return fmt.Sprintf("Time: %d, Price: %f, Size: %s", h.Time, h.Price, DecimalMaxString(h.Size))
}
