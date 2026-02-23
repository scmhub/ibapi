package ibapi

import (
	"github.com/scmhub/ibapi/protobuf"
	"google.golang.org/protobuf/proto"
)

// UpdateConfigAPISettings .
func UpdateConfigAPISettings(reqID int32) *protobuf.UpdateConfigRequest {
	settings := &protobuf.ApiSettingsConfig{
		TotalQuantityForMutualFunds:            proto.Bool(true),
		DownloadOpenOrdersOnConnection:         proto.Bool(true),
		IncludeVirtualFxPositions:              proto.Bool(true),
		PrepareDailyPnL:                        proto.Bool(true),
		SendStatusUpdatesForVolatilityOrders:   proto.Bool(true),
		EncodeApiMessages:                      proto.String("osCodePage"),
		SocketPort:                             proto.Int32(7497),
		UseNegativeAutoRange:                   proto.Bool(true),
		CreateApiMessageLogFile:                proto.Bool(true),
		IncludeMarketDataInLogFile:             proto.Bool(true),
		ExposeTradingScheduleToApi:             proto.Bool(true),
		SplitInsuredDepositFromCashBalance:     proto.Bool(true),
		SendZeroPositionsForTodayOnly:          proto.Bool(true),
		UseAccountGroupsWithAllocationMethods:  proto.Bool(true),
		LoggingLevel:                           proto.String("error"),
		MasterClientId:                         proto.Int32(3),
		BulkDataTimeout:                        proto.Int32(25),
		ComponentExchSeparator:                 proto.String("#"),
		RoundAccountValuesToNearestWholeNumber: proto.Bool(true),
		ShowAdvancedOrderRejectInUi:            proto.Bool(true),
		RejectMessagesAboveMaxRate:             proto.Bool(true),
		MaintainConnectionOnIncorrectFields:    proto.Bool(true),
		CompatibilityModeNasdaqStocks:          proto.Bool(true),
		SendInstrumentTimezone:                 proto.String("utc"),
		SendForexDataInCompatibilityMode:       proto.Bool(true),
		MaintainAndResubmitOrdersOnReconnect:   proto.Bool(true),
		HistoricalDataMaxSize:                  proto.Int32(4),
		AutoReportNettingEventContractTrades:   proto.Bool(true),
		OptionExerciseRequestType:              proto.String("final"),
		TrustedIPs:                             []string{"127.0.0.1"},
	}

	return &protobuf.UpdateConfigRequest{
		ReqId: proto.Int32(reqID),
		Api: &protobuf.ApiConfig{
			Settings: settings,
		},
	}
}

// UpdateOrdersConfig .
func UpdateOrdersConfig(reqID int32) *protobuf.UpdateConfigRequest {
	orders := &protobuf.OrdersConfig{
		SmartRouting: &protobuf.OrdersSmartRoutingConfig{
			SeekPriceImprovement:  proto.Bool(true),
			DoNotRouteToDarkPools: proto.Bool(true),
		},
	}

	return &protobuf.UpdateConfigRequest{
		ReqId:  proto.Int32(reqID),
		Orders: orders,
	}
}

// UpdateMessageConfigConfirmMandatoryCapPriceAccepted .
func UpdateMessageConfigConfirmMandatoryCapPriceAccepted(reqID int32) *protobuf.UpdateConfigRequest {
	messageID := proto.Int32(131)
	messageConfig := &protobuf.MessageConfig{
		Id:      messageID,
		Enabled: proto.Bool(false),
	}
	warning := &protobuf.UpdateConfigWarning{MessageId: messageID}

	return &protobuf.UpdateConfigRequest{
		ReqId:            proto.Int32(reqID),
		Messages:         []*protobuf.MessageConfig{messageConfig},
		AcceptedWarnings: []*protobuf.UpdateConfigWarning{warning},
	}
}

// UpdateConfigOrderIDReset .
func UpdateConfigOrderIDReset(reqID int32) *protobuf.UpdateConfigRequest {
	return &protobuf.UpdateConfigRequest{
		ReqId:                 proto.Int32(reqID),
		ResetAPIOrderSequence: proto.Bool(true),
	}
}
