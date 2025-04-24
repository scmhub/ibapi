package ibapi

import "fmt"

// Bar .
type Bar struct {
	Date     string
	Open     float64
	High     float64
	Low      float64
	Close    float64
	Volume   Decimal
	Wap      Decimal
	BarCount int64
}

func NewBar() Bar {
	b := Bar{}
	b.Volume = UNSET_DECIMAL
	b.Wap = UNSET_DECIMAL
	return b
}

func (b Bar) String() string {
	return fmt.Sprintf("Date: %s, Open: %f, High: %f, Low: %f, Close: %f, Volume: %s, WAP: %s, BarCount: %d",
		b.Date, b.Open, b.High, b.Low, b.Close, DecimalMaxString(b.Volume), DecimalMaxString(b.Wap), b.BarCount)
}

// RealTimeBar .
type RealTimeBar struct {
	Time    int64
	EndTime int64
	Open    float64
	High    float64
	Low     float64
	Close   float64
	Volume  Decimal
	Wap     Decimal
	Count   int64
}

func NewRealTimeBar() RealTimeBar {
	rtb := RealTimeBar{}
	rtb.Volume = UNSET_DECIMAL
	rtb.Wap = UNSET_DECIMAL
	return rtb
}

func (rb RealTimeBar) String() string {
	return fmt.Sprintf("Time: %d, Open: %f, High: %f, Low: %f, Close: %f, Volume: %s, Wap: %s, Count: %d",
		rb.Time, rb.Open, rb.High, rb.Low, rb.Close, DecimalMaxString(rb.Volume), DecimalMaxString(rb.Wap), rb.Count)
}
