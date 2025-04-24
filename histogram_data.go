package ibapi

import "fmt"

// HistogramData .
type HistogramData struct {
	Price float64
	Size  Decimal
}

func NewHistogramData() HistogramData {
	hd := HistogramData{}
	hd.Size = UNSET_DECIMAL
	return hd
}

func (hd HistogramData) String() string {
	return fmt.Sprintf("Price: %v, Size: %v", hd.Price, hd.Size)
}
