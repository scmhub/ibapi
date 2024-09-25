package ibapi

import "fmt"

// OrderCancel .
type OrderCancel struct {
	ManualOrderCancelTime string
}

func NewOrderCancel() OrderCancel {
	return OrderCancel{}
}

func (o OrderCancel) String() string {
	return fmt.Sprintf("manualOrderCancelTime: %s", o.ManualOrderCancelTime)
}
