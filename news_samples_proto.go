package ibapi

import (
	"time"

	"github.com/scmhub/ibapi/protobuf"
)

// HistoricalNewsRequestWithEndTime .
func HistoricalNewsRequestWithEndTime(reqID int32) *protobuf.HistoricalNewsRequest {
	endDateTime := time.Now().AddDate(0, 0, -10).Format("2006-01-02 15:04:05.0")

	return &protobuf.HistoricalNewsRequest{
		ReqId:         new(reqID),
		ConId:         new(int32(8314)),
		ProviderCodes: new("BRFUPDN+BRFG"),
		EndDateTime:   new(endDateTime),
		TotalResults:  new(int32(10)),
	}
}

// HistoricalNewsRequestWithStartTime .
func HistoricalNewsRequestWithStartTime(reqID int32) *protobuf.HistoricalNewsRequest {
	startDateTime := time.Now().AddDate(0, 0, -10).Format("2006-01-02 15:04:05.0")

	return &protobuf.HistoricalNewsRequest{
		ReqId:         new(reqID),
		ConId:         new(int32(8314)),
		ProviderCodes: new("BRFUPDN+BRFG"),
		StartDateTime: new(startDateTime),
		TotalResults:  new(int32(10)),
	}
}
