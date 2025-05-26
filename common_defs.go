package ibapi

import (
	"fmt"
	"math"
)

const (
	UNSET_INT       int64   = math.MaxInt32
	UNSET_LONG      int64   = math.MaxInt64
	UNSET_FLOAT     float64 = math.MaxFloat64
	INFINITY_STRING string  = "Infinity"
)

var INFINITY_FLOAT float64 = math.Inf(1)

type TickerID = int64

type OrderID = int64

type FaDataType int64

const (
	GROUPS  FaDataType = 1
	ALIASES FaDataType = 3
)

func (fa FaDataType) String() string {
	switch fa {
	case GROUPS:
		return "GROUPS"
	case ALIASES:
		return "ALIASES"
	default:
		return ""
	}
}

type MarketDataType int64

const (
	REALTIME       MarketDataType = 1
	FROZEN         MarketDataType = 2
	DELAYED        MarketDataType = 3
	DELAYED_FROZEN MarketDataType = 4
)

func (mdt MarketDataType) String() string {
	switch mdt {
	case REALTIME:
		return "REALTIME"
	case FROZEN:
		return "FROZEN"
	case DELAYED:
		return "DELAYED"
	case DELAYED_FROZEN:
		return "DELAYED_FROZEN"
	}
	return ""
}

// SmartComponent .
type SmartComponent struct {
	BitNumber      int64
	Exchange       string
	ExchangeLetter string
}

func NewSmartComponent() SmartComponent {
	return SmartComponent{}
}

func (s SmartComponent) String() string {
	return fmt.Sprintf("BitNumber: %d, Exchange: %s, ExchangeLetter: %s", s.BitNumber, s.Exchange, s.ExchangeLetter)
}

// FundAssetType .
type FundAssetType string

const (
	FundAssetTypeNone        FundAssetType = ""
	FundAssetTypeOthers      FundAssetType = "000"
	FundAssetTypeMoneyMarket FundAssetType = "001"
	FundAssetTypeFixedIncome FundAssetType = "002"
	FundAssetTypeMultiAsset  FundAssetType = "003"
	FundAssetTypeEquity      FundAssetType = "004"
	FundAssetTypeSector      FundAssetType = "005"
	FundAssetTypeGuaranteed  FundAssetType = "006"
	FundAssetTypeAlternative FundAssetType = "007"
)

func getFundAssetType(fat string) FundAssetType {
	switch fat {
	case string(FundAssetTypeOthers):
		return FundAssetTypeOthers
	case string(FundAssetTypeMoneyMarket):
		return FundAssetTypeMoneyMarket
	case string(FundAssetTypeFixedIncome):
		return FundAssetTypeFixedIncome
	case string(FundAssetTypeMultiAsset):
		return FundAssetTypeMultiAsset
	case string(FundAssetTypeEquity):
		return FundAssetTypeEquity
	case string(FundAssetTypeSector):
		return FundAssetTypeSector
	case string(FundAssetTypeGuaranteed):
		return FundAssetTypeGuaranteed
	case string(FundAssetTypeAlternative):
		return FundAssetTypeAlternative
	default:
		return FundAssetTypeNone
	}
}

// FundDistributionPolicyIndicator .
type FundDistributionPolicyIndicator string

const (
	FundDistributionPolicyIndicatorNone             FundDistributionPolicyIndicator = ""
	FundDistributionPolicyIndicatorAccumulationFund FundDistributionPolicyIndicator = "N"
	FundDistributionPolicyIndicatorIncomeFund       FundDistributionPolicyIndicator = "Y"
)

func getFundDistributionPolicyIndicator(fat string) FundDistributionPolicyIndicator {
	switch fat {
	case string(FundDistributionPolicyIndicatorAccumulationFund):
		return FundDistributionPolicyIndicatorAccumulationFund
	case string(FundDistributionPolicyIndicatorIncomeFund):
		return FundDistributionPolicyIndicatorIncomeFund
	default:
		return FundDistributionPolicyIndicatorNone
	}
}

type OptionExerciseType int

const (
	OptionExerciseTypeNone                 OptionExerciseType = -1
	OptionExerciseTypeExercise             OptionExerciseType = 1
	OptionExerciseTypeLapse                OptionExerciseType = 2
	OptionExerciseTypeDoNothing            OptionExerciseType = 3
	OptionExerciseTypeAssigned             OptionExerciseType = 100
	OptionExerciseTypeAutoexerciseClearing OptionExerciseType = 101
	OptionExerciseTypeExpired              OptionExerciseType = 102
	OptionExerciseTypeNetting              OptionExerciseType = 103
	OptionExerciseTypeAutoexerciseTrading  OptionExerciseType = 104
)

func (e OptionExerciseType) String() string {
	switch e {
	case OptionExerciseTypeNone:
		return "None"
	case OptionExerciseTypeExercise:
		return "Exercise"
	case OptionExerciseTypeLapse:
		return "Lapse"
	case OptionExerciseTypeDoNothing:
		return "DoNothing"
	case OptionExerciseTypeAssigned:
		return "Assigned"
	case OptionExerciseTypeAutoexerciseClearing:
		return "AutoexerciseClearing"
	case OptionExerciseTypeExpired:
		return "Expired"
	case OptionExerciseTypeNetting:
		return "Netting"
	case OptionExerciseTypeAutoexerciseTrading:
		return "AutoexerciseTrading"
	default:
		return "Unknown"
	}
}
