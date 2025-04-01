package ibapi

import (
	protobf "github.com/scmhub/ibapi/protobuf"
)

func decodeContract(contractProto *protobf.Contract) *Contract {
	contract := NewContract()
	if contractProto.ConId != nil {
		contract.ConID = int64(contractProto.GetConId())
	}
	if contractProto.Symbol != nil {
		contract.Symbol = contractProto.GetSymbol()
	}
	if contractProto.SecType != nil {
		contract.SecType = contractProto.GetSecType()
	}
	if contractProto.LastTradeDateOrContractMonth != nil {
		contract.LastTradeDateOrContractMonth = contractProto.GetLastTradeDateOrContractMonth()
	}
	if contractProto.Strike != nil {
		contract.Strike = contractProto.GetStrike()
	}
	if contractProto.Right != nil {
		contract.Right = contractProto.GetRight()
	}
	if contractProto.Multiplier != nil {
		contract.Multiplier = FloatMaxString(contractProto.GetMultiplier())
	}
	if contractProto.Exchange != nil {
		contract.Exchange = contractProto.GetExchange()
	}
	if contractProto.Currency != nil {
		contract.Currency = contractProto.GetCurrency()
	}
	if contractProto.LocalSymbol != nil {
		contract.LocalSymbol = contractProto.GetLocalSymbol()
	}
	if contractProto.TradingClass != nil {
		contract.TradingClass = contractProto.GetTradingClass()
	}
	return contract
}

func decodeExecution(executionProto *protobf.Execution) *Execution {
	execution := NewExecution()
	if executionProto.OrderId != nil {
		execution.OrderID = int64(executionProto.GetOrderId())
	}
	if executionProto.ClientId != nil {
		execution.ClientID = int64(executionProto.GetClientId())
	}
	if executionProto.ExecId != nil {
		execution.ExecID = executionProto.GetExecId()
	}
	if executionProto.Time != nil {
		execution.Time = executionProto.GetTime()
	}
	if executionProto.AcctNumber != nil {
		execution.AcctNumber = executionProto.GetAcctNumber()
	}
	if executionProto.Exchange != nil {
		execution.Exchange = executionProto.GetExchange()
	}
	if executionProto.Side != nil {
		execution.Side = executionProto.GetSide()
	}
	if executionProto.Shares != nil {
		execution.Shares = StringToDecimal(executionProto.GetShares())
	}
	if executionProto.Price != nil {
		execution.Price = executionProto.GetPrice()
	}
	if executionProto.PermId != nil {
		execution.PermID = int64(executionProto.GetPermId())
	}
	if executionProto.IsLiquidation != nil {
		execution.Liquidation = BoolToInt64(executionProto.GetIsLiquidation())
	}
	if executionProto.CumQty != nil {
		execution.CumQty = StringToDecimal(executionProto.GetCumQty())
	}
	if executionProto.AvgPrice != nil {
		execution.AvgPrice = executionProto.GetAvgPrice()
	}
	if executionProto.OrderRef != nil {
		execution.OrderRef = executionProto.GetOrderRef()
	}
	if executionProto.EvRule != nil {
		execution.EVRule = executionProto.GetEvRule()
	}
	if executionProto.EvMultiplier != nil {
		execution.EVMultiplier = executionProto.GetEvMultiplier()
	}
	if executionProto.ModelCode != nil {
		execution.ModelCode = executionProto.GetModelCode()
	}
	if executionProto.LastLiquidity != nil {
		execution.LastLiquidity = int64(executionProto.GetLastLiquidity())
	}
	if executionProto.IsPriceRevisionPending != nil {
		execution.PendingPriceRevision = executionProto.GetIsPriceRevisionPending()
	}
	if executionProto.Submitter != nil {
		execution.Submitter = executionProto.GetSubmitter()
	}
	return execution
}
