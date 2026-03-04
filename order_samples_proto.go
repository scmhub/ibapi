package ibapi

import "github.com/scmhub/ibapi/protobuf"

// CreatePlaceOrderRequest .
func CreatePlaceOrderRequest(orderID int64, contractProto *protobuf.Contract, orderProto *protobuf.Order) *protobuf.PlaceOrderRequest {

	placeOrderRequestProto := protobuf.PlaceOrderRequest{
		OrderId:  new(int32(orderID)),
		Contract: contractProto,
		Order:    orderProto,
	}

	return &placeOrderRequestProto
}

// LimitOrderProto .
func LimitOrderProto(action string, quantity Decimal, limitPrice float64, transmit bool) *protobuf.Order {

	orderProto := protobuf.Order{
		Action:        new(action),
		OrderType:     new("LMT"),
		TotalQuantity: new(DecimalToString(quantity)),
		LmtPrice:      new(limitPrice),
		Tif:           new("DAY"),
		Transmit:      new(transmit),
	}

	return &orderProto
}

// BetaHedgeOrder .
func BetaHedgeOrder(parentID int64, action string, hedgeParam string, hedgeMaxSize int64, transmit bool) *protobuf.Order {

	orderProto := protobuf.Order{
		ParentId:     new(int32(parentID)),
		Action:       new(action),
		OrderType:    new("MKT"),
		Tif:          new("DAY"),
		HedgeType:    new("B"),
		HedgeParam:   new(hedgeParam),
		HedgeMaxSize: new(int32(hedgeMaxSize)),
		Transmit:     new(transmit),
	}

	return &orderProto
}
