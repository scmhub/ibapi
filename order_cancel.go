package ibapi

import "fmt"

// OrderCancel .
type OrderCancel struct {
	ManualOrderCancelTime string
	ExtOperator           string
	ManualOrderIndicator  int64
}

func NewOrderCancel() OrderCancel {
	return OrderCancel{
		ManualOrderIndicator: UNSET_INT,
	}
}

func (o OrderCancel) String() string {
	return fmt.Sprintf("ManualOrderCancelTime: %s, ManualOrderIndicator: %s", o.ManualOrderCancelTime, intMaxString(o.ManualOrderIndicator))
}
