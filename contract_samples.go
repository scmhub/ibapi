package ibapi

/*
Contracts can be defined in multiple ways. The TWS/IB Gateway will always perform a query on the available contracts
and find which one is the best candidate:
 - More than a single candidate will yield an ambiguity error message.
 - No suitable candidates will produce a "contract not found" message.
How do I find my contract though?
 - Often the quickest way is by looking for it in the TWS and looking at its description there (double click).
 - The TWS' symbol corresponds to the API's localSymbol. Keep this in mind when defining Futures or Options.
 - The TWS' underlying's symbol can usually be mapped to the API's symbol.

Any stock or option symbols displayed are for illustrative purposes only and are not intended to portray a recommendation.

Usually, the easiest way to define a Stock/CASH contract is through these four attributes.
*/

// IBMBond .
func IBMBond() *Contract {

	contract := NewContract()
	contract.Symbol = "IBM"
	contract.SecType = "BOND"
	contract.Currency = "USD"
	contract.Exchange = "SMART"

	return contract
}

// IBKRStk .
func IBKRStk() *Contract {

	contract := NewContract()
	contract.Symbol = "IBKR"
	contract.SecType = "STK"
	contract.Currency = "USD"
	contract.Exchange = "SMART"

	return contract
}

// HKStk
func HKStk() *Contract {

	contract := NewContract()
	contract.Symbol = "1"
	contract.SecType = "STK"
	contract.Currency = "HKD"
	contract.Exchange = "SEHK"

	return contract
}

// EurGbpFx .
func EurGbpFx() *Contract {

	contract := NewContract()
	contract.Symbol = "EUR"
	contract.SecType = "CASH"
	contract.Currency = "GBP"
	contract.Exchange = "IDEALPRO"

	return contract
}

// Index .
func Index() *Contract {

	contract := NewContract()
	contract.Symbol = "DAX"
	contract.SecType = "IND"
	contract.Currency = "EUR"
	contract.Exchange = "EUREX"

	return contract
}

// CFD .
func CFD() *Contract {

	contract := NewContract()
	contract.Symbol = "IBDE30"
	contract.SecType = "CFD"
	contract.Currency = "EUR"
	contract.Exchange = "SMART"

	return contract
}

// USStockCFD .
func USStockCFD() *Contract {

	contract := NewContract()
	contract.Symbol = "IBM"
	contract.SecType = "CFD"
	contract.Currency = "USD"
	contract.Exchange = "SMART"

	return contract
}

// EuropeanStockCFD .
func EuropeanStockCFD() *Contract {

	contract := NewContract()
	contract.Symbol = "BMW"
	contract.SecType = "CFD"
	contract.Currency = "EUR"
	contract.Exchange = "SMART"

	return contract
}

// CashCFD .
func CashCFD() *Contract {

	contract := NewContract()
	contract.Symbol = "EUR"
	contract.SecType = "CFD"
	contract.Currency = "USD"
	contract.Exchange = "SMART"

	return contract
}

// EuropeanStock .
func EuropeanStock() *Contract {

	contract := NewContract()
	contract.Symbol = "NOKIA"
	contract.SecType = "STK"
	contract.Currency = "EUR"
	contract.Exchange = "SMART"
	contract.PrimaryExchange = "HEX"

	return contract
}

// OptionAtIse .
func OptionAtIse() *Contract {

	contract := NewContract()
	contract.Symbol = "BPX"
	contract.SecType = "OPT"
	contract.Currency = "USD"
	contract.Exchange = "ISE"
	contract.LastTradeDateOrContractMonth = "20160916"
	contract.Right = "C"
	contract.Strike = 65
	contract.Multiplier = "100"

	return contract
}

// USStock .
func USStock() *Contract {

	contract := NewContract()
	contract.Symbol = "SPY"
	contract.SecType = "STK"
	contract.Currency = "USD"
	contract.Exchange = "ARCA"

	return contract
}

// ETF .
func ETF() *Contract {

	contract := NewContract()
	contract.Symbol = "QQQ"
	contract.SecType = "STK"
	contract.Currency = "USD"
	contract.Exchange = "SMART"

	return contract
}

// USStockAtSmart .
func USStockAtSmart() *Contract {

	contract := NewContract()
	contract.Symbol = "IBM"
	contract.SecType = "STK"
	contract.Currency = "USD"
	contract.Exchange = "SMART"

	return contract
}

// IBMUSStockAtSmart .
func IBMUSStockAtSmart() *Contract {

	contract := NewContract()
	contract.Symbol = "IBM"
	contract.SecType = "STK"
	contract.Currency = "USD"
	contract.Exchange = "SMART"

	return contract
}

// USStockWithPrimaryExch .
func USStockWithPrimaryExch() *Contract {

	contract := NewContract()
	contract.Symbol = "SPY"
	contract.SecType = "STK"
	contract.Currency = "USD"
	contract.Exchange = "SMART"
	contract.PrimaryExchange = "ARCA"

	return contract
}

// BondWithCusip .
func BondWithCusip() *Contract {

	contract := NewContract()
	// enter CUSIP as symbol
	contract.Symbol = "449276AA2"
	contract.SecType = "BOND"
	contract.Exchange = "SMART"
	contract.Currency = "USD"

	return contract
}

// Bond .
func Bond() *Contract {

	contract := NewContract()
	contract.ConID = 456467716
	contract.Exchange = "SMART"

	return contract
}

// MutualFund .
func MutualFund() *Contract {

	contract := NewContract()
	contract.Symbol = "VINIX"
	contract.SecType = "FUND"
	contract.Exchange = "FUNDSERV"
	contract.Currency = "USD"

	return contract
}

// Commodity
func Commodity() *Contract {

	contract := NewContract()
	contract.Symbol = "XAUUSD"
	contract.SecType = "CMDTY"
	contract.Exchange = "SMART"
	contract.Currency = "USD"

	return contract
}

// USOptionContract .
func USOptionContract() *Contract {
	contract := NewContract()
	contract.Symbol = "GOOG"
	contract.SecType = "OPT"
	contract.Exchange = "SMART"
	contract.Currency = "USD"
	contract.LastTradeDateOrContractMonth = "20170120"
	contract.Strike = 615
	contract.Right = "C"
	contract.Multiplier = "100"
	return contract
}

// OptionAtBox
func OptionAtBox() *Contract {

	contract := NewContract()
	contract.Symbol = "GOOG"
	contract.SecType = "OPT"
	contract.Exchange = "BOX"
	contract.Currency = "USD"
	contract.LastTradeDateOrContractMonth = "20170120"
	contract.Strike = 615
	contract.Right = "C"
	contract.Multiplier = "100"

	return contract
}

// OptionWithTradingClass .
// Option contracts require far more information since there are many contracts having the exact same
// attributes such as symbol, currency, strike, etc.
func OptionWithTradingClass() *Contract {

	contract := NewContract()
	contract.Symbol = "SANT"
	contract.SecType = "OPT"
	contract.Exchange = "MEFFRV"
	contract.Currency = "EUR"
	contract.LastTradeDateOrContractMonth = "20190621"
	contract.Strike = 7.5
	contract.Right = "C"
	contract.Multiplier = "100"
	contract.TradingClass = "SANEU"

	return contract
}

// OptionWithLocalSymbol .
// Using the contract's own symbol (localSymbol) can greatly simplify a contract description
func OptionWithLocalSymbol() *Contract {

	contract := NewContract()
	//Watch out for the spaces within the local symbol!
	contract.LocalSymbol = "P BMW  20221216 72 M"
	contract.SecType = "OPT"
	contract.Exchange = "EUREX"
	contract.Currency = "EUR"
	//! [optcontract_localsymbol]
	return contract
}

// DutchWarrant .
// Dutch Warrants (IOPTs) can be defined using the local symbol or conid
func DutchWarrant() *Contract {

	contract := NewContract()
	contract.LocalSymbol = "B881G"
	contract.SecType = "IOPT"
	contract.Exchange = "SBF"
	contract.Currency = "EUR"

	return contract
}

// SimpleFuture .
// Future contracts also require an expiration date but are less complicated than options.
func SimpleFuture() *Contract {

	contract := NewContract()
	contract.Symbol = "GBL"
	contract.SecType = "FUT"
	contract.Exchange = "EUREX"
	contract.Currency = "EUR"
	contract.LastTradeDateOrContractMonth = "202303"

	return contract
}

// FutureWithLocalSymbol .
// Rather than giving expiration dates we can also provide the local symbol
// attributes such as symbol, currency, strike, etc.
func FutureWithLocalSymbol() *Contract {

	contract := NewContract()
	contract.SecType = "FUT"
	contract.Exchange = "EUREX"
	contract.Currency = "EUR"
	contract.LocalSymbol = "FGBL MAR 23"

	return contract
}

// FutureWithMultiplier .
func FutureWithMultiplier() *Contract {

	contract := NewContract()
	contract.Symbol = "DAX"
	contract.SecType = "FUT"
	contract.Exchange = "EUREX"
	contract.Currency = "EUR"
	contract.LastTradeDateOrContractMonth = "202303"
	contract.Multiplier = "1"

	return contract
}

// WrongContract .
// Note the space in the symbol!
func WrongContract() *Contract {

	contract := NewContract()
	contract.Symbol = " IJR "
	contract.ConID = 9579976
	contract.SecType = "STK"
	contract.Exchange = "SMART"
	contract.Currency = "USD"

	return contract
}

// FuturesOnOptions .
func FuturesOnOptions() *Contract {

	contract := NewContract()
	contract.Symbol = "GBL"
	contract.SecType = "FOP"
	contract.Exchange = "EUREX"
	contract.Currency = "EUR"
	contract.LastTradeDateOrContractMonth = "20230224"
	contract.Strike = 138
	contract.Right = "C"
	contract.Multiplier = "1000"

	return contract
}

// Warrants .
func Warrants() *Contract {

	contract := NewContract()
	contract.Symbol = "GOOG"
	contract.SecType = "WAR"
	contract.Exchange = "FWB"
	contract.Currency = "EUR"
	contract.LastTradeDateOrContractMonth = "20201117"
	contract.Strike = 1500.0
	contract.Right = "C"
	contract.Multiplier = "0.01"

	return contract
}

// ByISIN .
// It is also possible to define contracts based on their ISIN (IBKR STK sample).
func ByISIN() *Contract {

	contract := NewContract()
	contract.SecIDType = "ISIN"
	contract.SecID = "US45841N1072"
	contract.Exchange = "SMART"
	contract.Currency = "USD"
	contract.SecType = "STK"

	return contract
}

// ByConId .
// It is also possible to define contracts based on their ConID (EUR.USD sample).
// Note: passing a contract containing the conId can cause problems if one of the other provided
// attributes does not match 100% with what is in IB's database. This is particularly important
// for contracts such as Bonds which may change their description from one day to another.
// If the conId is provided, it is best not to give too much information as in the example below.
func ByConId() *Contract {

	contract := NewContract()
	contract.ConID = 12087792
	contract.Exchange = "IDEALPRO"
	contract.SecType = "CASH"

	return contract
}

// OptionForQuery .
// Ambiguous contracts are great to use with reqContractDetails. This way you can
// query the whole option chain for an underlying. Bear in mind that there are
// pacing mechanisms in place which will delay any further responses from the TWS
// to prevent abuse.
func OptionForQuery() *Contract {

	contract := NewContract()
	contract.Symbol = "FISV"
	contract.SecType = "OPT"
	contract.Exchange = "SMART"
	contract.Currency = "USD"

	return contract
}

// OptionComboContract .
func OptionComboContract() *Contract {

	contract := NewContract()
	contract.Symbol = "DBK"
	contract.SecType = "BAG"
	contract.Currency = "EUR"
	contract.Exchange = "EUREX"

	leg1 := NewComboLeg()
	leg1.ConID = 577164786 //DBK Jun21'24 CALL @EUREX
	leg1.Action = "BUY"
	leg1.Ratio = 1
	leg1.Exchange = "EUREX"

	leg2 := NewComboLeg()
	leg2.ConID = 577164767 //DBK Dec15'23 CALL @EUREX
	leg2.Action = "SELL"
	leg2.Ratio = 1
	leg2.Exchange = "EUREX"

	contract.ComboLegs = []ComboLeg{}
	contract.ComboLegs = append(contract.ComboLegs, leg1)
	contract.ComboLegs = append(contract.ComboLegs, leg2)

	return contract
}

// StockComboContract .
// STK Combo contract
// Leg 1: 43645865 - IBKR's STK
// Leg 2: 9408 - McDonald's STK
func StockComboContract() *Contract {

	contract := NewContract()
	contract.Symbol = "MCD"
	contract.SecType = "BAG"
	contract.Currency = "USD"
	contract.Exchange = "SMART"

	leg1 := NewComboLeg()
	leg1.ConID = 43645865
	leg1.Action = "BUY"
	leg1.Ratio = 1
	leg1.Exchange = "SMART"

	leg2 := NewComboLeg()
	leg2.ConID = 9408
	leg2.Action = "SELL"
	leg2.Ratio = 1
	leg2.Exchange = "SMART"

	contract.ComboLegs = []ComboLeg{}
	contract.ComboLegs = append(contract.ComboLegs, leg1)
	contract.ComboLegs = append(contract.ComboLegs, leg2)

	return contract
}

// FutureComboContract .
// CBOE Volatility Index Future combo contract
// Leg 1: 195538625 - FUT expiring 2016/02/17
// Leg 2: 197436571 - FUT expiring 2016/03/16
func FutureComboContract() *Contract {

	contract := NewContract()
	contract.Symbol = "VIX"
	contract.SecType = "BAG"
	contract.Currency = "USD"
	contract.Exchange = "CFE"

	leg1 := NewComboLeg()
	leg1.ConID = 195538625
	leg1.Action = "BUY"
	leg1.Ratio = 1
	leg1.Exchange = "CFE"

	leg2 := NewComboLeg()
	leg2.ConID = 197436571
	leg2.Action = "SELL"
	leg2.Ratio = 1
	leg2.Exchange = "CFE"

	contract.ComboLegs = []ComboLeg{}
	contract.ComboLegs = append(contract.ComboLegs, leg1)
	contract.ComboLegs = append(contract.ComboLegs, leg2)

	return contract
}

// SmartFutureComboContract .
func SmartFutureComboContract() *Contract {

	contract := NewContract()
	contract.Symbol = "WTI" // WTI,COIL spread. Symbol can be defined as first leg symbol ("WTI") or currency ("USD").
	contract.SecType = "BAG"
	contract.Currency = "USD"
	contract.Exchange = "SMART"

	leg1 := NewComboLeg()
	leg1.ConID = 55928698 // WTI future June 2017
	leg1.Action = "BUY"
	leg1.Ratio = 1
	leg1.Exchange = "IPE"

	leg2 := NewComboLeg()
	leg2.ConID = 55850663 // COIL future June 2017
	leg2.Action = "SELL"
	leg2.Ratio = 1
	leg2.Exchange = "IPE"

	contract.ComboLegs = []ComboLeg{}
	contract.ComboLegs = append(contract.ComboLegs, leg1)
	contract.ComboLegs = append(contract.ComboLegs, leg2)

	return contract
}

// InterCmdtyFuturesContract .
func InterCmdtyFuturesContract() *Contract {

	contract := NewContract()
	contract.Symbol = "COIL.WTI"
	contract.SecType = "BAG"
	contract.Currency = "USD"
	contract.Exchange = "IPE"

	leg1 := NewComboLeg()
	leg1.ConID = 183405603 //WTI Dec'23 @IPE
	leg1.Action = "BUY"
	leg1.Ratio = 1
	leg1.Exchange = "IPE"

	leg2 := NewComboLeg()
	leg2.ConID = 254011009 //COIL Dec'23 @IPE
	leg2.Action = "SELL"
	leg2.Ratio = 1
	leg2.Exchange = "IPE"

	contract.ComboLegs = []ComboLeg{}
	contract.ComboLegs = append(contract.ComboLegs, leg1)
	contract.ComboLegs = append(contract.ComboLegs, leg2)

	return contract
}

// NewsFeedForQuery .
func NewsFeedForQuery() *Contract {

	contract := NewContract()
	contract.SecType = "NEWS"
	contract.Exchange = "BRF" //Briefing Trader

	return contract
}

// BTbroadtapeNewsFeed .
func BTbroadtapeNewsFeed() *Contract {

	contract := NewContract()
	contract.Symbol = "BRF:BRF_ALL" //BroadTape All News
	contract.SecType = "NEWS"
	contract.Exchange = "BRF" //Briefing Trader

	return contract
}

// BZbroadtapeNewsFeed .
func BZbroadtapeNewsFeed() *Contract {

	contract := NewContract()
	contract.Symbol = "BZ:BZ_ALL" //BroadTape All News
	contract.SecType = "NEWS"
	contract.Exchange = "BZ" //Benzinga Pro

	return contract
}

// FLYbroadtapeNewsFeed .
func FLYbroadtapeNewsFeed() *Contract {

	contract := NewContract()
	contract.Symbol = "FLY:FLY_ALL" //BroadTape All News
	contract.SecType = "NEWS"
	contract.Exchange = "FLY" //Fly on the Wall

	return contract
}

// ContFut .
func ContFut() *Contract {

	contract := NewContract()
	contract.Symbol = "GBL"
	contract.SecType = "CONTFUT"
	contract.Exchange = "EUREX"

	return contract
}

// ContAndExpiringFut .
func ContAndExpiringFut() *Contract {

	contract := NewContract()
	contract.Symbol = "GBL"
	contract.SecType = "FUT+CONTFUT"
	contract.Exchange = "EUREX"

	return contract
}

// JefferiesContract .
func JefferiesContract() *Contract {

	contract := NewContract()
	contract.Symbol = "AAPL"
	contract.SecType = "STK"
	contract.Exchange = "JEFFALGO" // must be direct-routed to JEFALGO
	contract.Currency = "USD"      // only available for US stocks

	return contract
}

// CSFBContract .
func CSFBContract() *Contract {

	contract := NewContract()
	contract.Symbol = "IBKR"
	contract.SecType = "STK"
	contract.Exchange = "CSFBALGO"
	contract.Currency = "USD"

	return contract
}

// IBKRATSContract .
func IBKRATSContract() *Contract {

	contract := NewContract()
	contract.Symbol = "SPY"
	contract.SecType = "STK"
	contract.Exchange = "IBKRATS"
	contract.Currency = "USD"

	return contract
}

// CryptoContract .
func CryptoContract() *Contract {

	contract := NewContract()
	contract.Symbol = "BTC"
	contract.SecType = "CRYPTO"
	contract.Exchange = "PAXOS"
	contract.Currency = "USD"

	return contract
}

// StockWithIPOPrice .
func StockWithIPOPrice() *Contract {

	contract := NewContract()
	contract.Symbol = "EMCGU"
	contract.SecType = "STK"
	contract.Exchange = "SMART"
	contract.Currency = "USD"

	return contract
}

// ByFIGI .
func ByFIGI() *Contract {

	contract := NewContract()
	contract.SecIDType = "FIGI"
	contract.SecID = "BBG000B9XRY4"
	contract.Exchange = "SMART"

	return contract
}

// ByIssuerId .
func ByIssuerId() *Contract {

	contract := NewContract()
	contract.IssuerID = "e1453318"

	return contract
}

// Fund .
func Fund() *Contract {

	contract := NewContract()
	contract.Symbol = "I406801954"
	contract.SecType = "FUND"
	contract.Exchange = "ALLFUNDS"
	contract.Currency = "USD"

	return contract
}
