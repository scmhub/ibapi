package ibapi

import "github.com/scmhub/ibapi/protobuf"

// OddLotBidAskQuotesRequest .
func OddLotBidAskQuotesRequest(reqID int32, contractProto *protobuf.Contract) *protobuf.MarketDataRequest {
	return &protobuf.MarketDataRequest{
		ReqId:           new(reqID),
		Contract:        contractProto,
		GenericTickList: new("787"),
	}
}

// RegulatorySnapshotRequest .
func RegulatorySnapshotRequest(reqID int32, contractProto *protobuf.Contract) *protobuf.MarketDataRequest {
	return &protobuf.MarketDataRequest{
		ReqId:              new(reqID),
		Contract:           contractProto,
		RegulatorySnapshot: new(true),
	}
}

// CancelMarketDataRequest .
func CancelMarketDataRequest(reqID int32) *protobuf.CancelMarketData {
	return &protobuf.CancelMarketData{
		ReqId: new(reqID),
	}
}
