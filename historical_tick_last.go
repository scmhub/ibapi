package ibapi

import (
	"fmt"
	"time"
)

// HistoricalTickLast is the historical last tick's description.
// Used when requesting historical tick data with whatToShow = TRADES.
type HistoricalTickLast struct {
	Time              int64
	TickAttribLast    TickAttribLast
	Price             float64
	Size              Decimal
	Exchange          string
	SpecialConditions string
}

func NewHistoricalTickLast() HistoricalTickLast {
	htl := HistoricalTickLast{}
	htl.Size = UNSET_DECIMAL
	return htl
}

func (h HistoricalTickLast) String() string {
	return fmt.Sprintf("Time: %d, TickAttribLast: %s, Price: %f, Size: %s, Exchange: %s, SpecialConditions: %s",
		h.Time, h.TickAttribLast, h.Price, DecimalMaxString(h.Size), h.Exchange, h.SpecialConditions)
}

func (h HistoricalTickLast) Timestamp() time.Time {
	return time.Unix(h.Time, 0)
}
