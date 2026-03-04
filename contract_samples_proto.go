package ibapi

import "github.com/scmhub/ibapi/protobuf"

// IBMStockAtSmart .
func IBMStockAtSmart() *protobuf.Contract {

	contractProto := protobuf.Contract{
		Symbol:   new("IBM"),
		SecType:  new("STK"),
		Exchange: new("SMART"),
		Currency: new("USD"),
	}

	return &contractProto
}

// MSFTStockAtSmart .
func MSFTStockAtSmart() *protobuf.Contract {

	contractProto := protobuf.Contract{
		Symbol:   new("MSFT"),
		SecType:  new("STK"),
		Exchange: new("SMART"),
		Currency: new("USD"),
	}

	return &contractProto
}
