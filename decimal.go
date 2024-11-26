package ibapi

import (
	"errors"

	"github.com/robaho/fixed"
)

var UNSET_DECIMAL = Decimal(fixed.NaN)
var ZERO = Decimal(fixed.ZERO)
var ONE = Decimal(fixed.NewF(1))

// Decimal implements Fixed from "github.com/robaho/fixed".
type Decimal fixed.Fixed

func (d Decimal) String() string {
	return fixed.Fixed(d).String()
}

func (d Decimal) Int() int64 {
	return fixed.Fixed(d).Int()
}

func (d Decimal) Float() float64 {
	return fixed.Fixed(d).Float()
}

func (d Decimal) MarshalBinary() ([]byte, error) {
	return fixed.Fixed(d).MarshalBinary()
}

func (d *Decimal) UnmarshalBinary(data []byte) error {
	var f fixed.Fixed
	err := f.UnmarshalBinary(data)
	if err != nil {
		return err
	}
	*d = Decimal(f)
	return nil
}

func StringToDecimal(s string) Decimal {
	d, _ := StringToDecimalErr(s)
	return d
}

func StringToDecimalErr(s string) (Decimal, error) {
	if s == "" || s == "2147483647" || s == "9223372036854775807" || s == "1.7976931348623157E308" || s == "-9223372036854775808" {
		return UNSET_DECIMAL, errors.New("unset decimal")
	}
	f, err := fixed.NewSErr(s)
	if err != nil {
		return UNSET_DECIMAL, err
	}
	return Decimal(f), nil
}

func DecimalToString(d Decimal) string {
	return fixed.Fixed(d).String()
}
