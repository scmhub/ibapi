package ibapi

/*
High level IB message info.
*/

// IN is the incoming msg id's
type IN = int64

const (
	TICK_PRICE                               IN = 1
	TICK_SIZE                                IN = 2
	ORDER_STATUS                             IN = 3
	ERR_MSG                                  IN = 4
	OPEN_ORDER                               IN = 5
	ACCT_VALUE                               IN = 6
	PORTFOLIO_VALUE                          IN = 7
	ACCT_UPDATE_TIME                         IN = 8
	NEXT_VALID_ID                            IN = 9
	CONTRACT_DATA                            IN = 10
	EXECUTION_DATA                           IN = 11
	MARKET_DEPTH                             IN = 12
	MARKET_DEPTH_L2                          IN = 13
	NEWS_BULLETINS                           IN = 14
	MANAGED_ACCTS                            IN = 15
	RECEIVE_FA                               IN = 16
	HISTORICAL_DATA                          IN = 17
	BOND_CONTRACT_DATA                       IN = 18
	SCANNER_PARAMETERS                       IN = 19
	SCANNER_DATA                             IN = 20
	TICK_OPTION_COMPUTATION                  IN = 21
	TICK_GENERIC                             IN = 45
	TICK_STRING                              IN = 46
	TICK_EFP                                 IN = 47
	CURRENT_TIME                             IN = 49
	REAL_TIME_BARS                           IN = 50
	FUNDAMENTAL_DATA                         IN = 51
	CONTRACT_DATA_END                        IN = 52
	OPEN_ORDER_END                           IN = 53
	ACCT_DOWNLOAD_END                        IN = 54
	EXECUTION_DATA_END                       IN = 55
	DELTA_NEUTRAL_VALIDATION                 IN = 56
	TICK_SNAPSHOT_END                        IN = 57
	MARKET_DATA_TYPE                         IN = 58
	COMMISSION_AND_FEES_REPORT               IN = 59
	POSITION_DATA                            IN = 61
	POSITION_END                             IN = 62
	ACCOUNT_SUMMARY                          IN = 63
	ACCOUNT_SUMMARY_END                      IN = 64
	VERIFY_MESSAGE_API                       IN = 65
	VERIFY_COMPLETED                         IN = 66
	DISPLAY_GROUP_LIST                       IN = 67
	DISPLAY_GROUP_UPDATED                    IN = 68
	VERIFY_AND_AUTH_MESSAGE_API              IN = 69
	VERIFY_AND_AUTH_COMPLETED                IN = 70
	POSITION_MULTI                           IN = 71
	POSITION_MULTI_END                       IN = 72
	ACCOUNT_UPDATE_MULTI                     IN = 73
	ACCOUNT_UPDATE_MULTI_END                 IN = 74
	SECURITY_DEFINITION_OPTION_PARAMETER     IN = 75
	SECURITY_DEFINITION_OPTION_PARAMETER_END IN = 76
	SOFT_DOLLAR_TIERS                        IN = 77
	FAMILY_CODES                             IN = 78
	SYMBOL_SAMPLES                           IN = 79
	MKT_DEPTH_EXCHANGES                      IN = 80
	TICK_REQ_PARAMS                          IN = 81
	SMART_COMPONENTS                         IN = 82
	NEWS_ARTICLE                             IN = 83
	TICK_NEWS                                IN = 84
	NEWS_PROVIDERS                           IN = 85
	HISTORICAL_NEWS                          IN = 86
	HISTORICAL_NEWS_END                      IN = 87
	HEAD_TIMESTAMP                           IN = 88
	HISTOGRAM_DATA                           IN = 89
	HISTORICAL_DATA_UPDATE                   IN = 90
	REROUTE_MKT_DATA_REQ                     IN = 91
	REROUTE_MKT_DEPTH_REQ                    IN = 92
	MARKET_RULE                              IN = 93
	PNL                                      IN = 94
	PNL_SINGLE                               IN = 95
	HISTORICAL_TICKS                         IN = 96
	HISTORICAL_TICKS_BID_ASK                 IN = 97
	HISTORICAL_TICKS_LAST                    IN = 98
	TICK_BY_TICK                             IN = 99
	ORDER_BOUND                              IN = 100
	COMPLETED_ORDER                          IN = 101
	COMPLETED_ORDERS_END                     IN = 102
	REPLACE_FA_END                           IN = 103
	WSH_META_DATA                            IN = 104
	WSH_EVENT_DATA                           IN = 105
	HISTORICAL_SCHEDULE                      IN = 106
	USER_INFO                                IN = 107
	HISTORICAL_DATA_END                      IN = 108
	CURRENT_TIME_IN_MILLIS                   IN = 109
)

// OUT is the outgoing msg id's.
type OUT = int64

const (
	REQ_MKT_DATA                  OUT = 1
	CANCEL_MKT_DATA               OUT = 2
	PLACE_ORDER                   OUT = 3
	CANCEL_ORDER                  OUT = 4
	REQ_OPEN_ORDERS               OUT = 5
	REQ_ACCT_DATA                 OUT = 6
	REQ_EXECUTIONS                OUT = 7
	REQ_IDS                       OUT = 8
	REQ_CONTRACT_DATA             OUT = 9
	REQ_MKT_DEPTH                 OUT = 10
	CANCEL_MKT_DEPTH              OUT = 11
	REQ_NEWS_BULLETINS            OUT = 12
	CANCEL_NEWS_BULLETINS         OUT = 13
	SET_SERVER_LOGLEVEL           OUT = 14
	REQ_AUTO_OPEN_ORDERS          OUT = 15
	REQ_ALL_OPEN_ORDERS           OUT = 16
	REQ_MANAGED_ACCTS             OUT = 17
	REQ_FA                        OUT = 18
	REPLACE_FA                    OUT = 19
	REQ_HISTORICAL_DATA           OUT = 20
	EXERCISE_OPTIONS              OUT = 21
	REQ_SCANNER_SUBSCRIPTION      OUT = 22
	CANCEL_SCANNER_SUBSCRIPTION   OUT = 23
	REQ_SCANNER_PARAMETERS        OUT = 24
	CANCEL_HISTORICAL_DATA        OUT = 25
	REQ_CURRENT_TIME              OUT = 49
	REQ_REAL_TIME_BARS            OUT = 50
	CANCEL_REAL_TIME_BARS         OUT = 51
	REQ_FUNDAMENTAL_DATA          OUT = 52
	CANCEL_FUNDAMENTAL_DATA       OUT = 53
	REQ_CALC_IMPLIED_VOLAT        OUT = 54
	REQ_CALC_OPTION_PRICE         OUT = 55
	CANCEL_CALC_IMPLIED_VOLAT     OUT = 56
	CANCEL_CALC_OPTION_PRICE      OUT = 57
	REQ_GLOBAL_CANCEL             OUT = 58
	REQ_MARKET_DATA_TYPE          OUT = 59
	REQ_POSITIONS                 OUT = 61
	REQ_ACCOUNT_SUMMARY           OUT = 62
	CANCEL_ACCOUNT_SUMMARY        OUT = 63
	CANCEL_POSITIONS              OUT = 64
	VERIFY_REQUEST                OUT = 65
	VERIFY_MESSAGE                OUT = 66
	QUERY_DISPLAY_GROUPS          OUT = 67
	SUBSCRIBE_TO_GROUP_EVENTS     OUT = 68
	UPDATE_DISPLAY_GROUP          OUT = 69
	UNSUBSCRIBE_FROM_GROUP_EVENTS OUT = 70
	START_API                     OUT = 71
	VERIFY_AND_AUTH_REQUEST       OUT = 72
	VERIFY_AND_AUTH_MESSAGE       OUT = 73
	REQ_POSITIONS_MULTI           OUT = 74
	CANCEL_POSITIONS_MULTI        OUT = 75
	REQ_ACCOUNT_UPDATES_MULTI     OUT = 76
	CANCEL_ACCOUNT_UPDATES_MULTI  OUT = 77
	REQ_SEC_DEF_OPT_PARAMS        OUT = 78
	REQ_SOFT_DOLLAR_TIERS         OUT = 79
	REQ_FAMILY_CODES              OUT = 80
	REQ_MATCHING_SYMBOLS          OUT = 81
	REQ_MKT_DEPTH_EXCHANGES       OUT = 82
	REQ_SMART_COMPONENTS          OUT = 83
	REQ_NEWS_ARTICLE              OUT = 84
	REQ_NEWS_PROVIDERS            OUT = 85
	REQ_HISTORICAL_NEWS           OUT = 86
	REQ_HEAD_TIMESTAMP            OUT = 87
	REQ_HISTOGRAM_DATA            OUT = 88
	CANCEL_HISTOGRAM_DATA         OUT = 89
	CANCEL_HEAD_TIMESTAMP         OUT = 90
	REQ_MARKET_RULE               OUT = 91
	REQ_PNL                       OUT = 92
	CANCEL_PNL                    OUT = 93
	REQ_PNL_SINGLE                OUT = 94
	CANCEL_PNL_SINGLE             OUT = 95
	REQ_HISTORICAL_TICKS          OUT = 96
	REQ_TICK_BY_TICK_DATA         OUT = 97
	CANCEL_TICK_BY_TICK_DATA      OUT = 98
	REQ_COMPLETED_ORDERS          OUT = 99
	REQ_WSH_META_DATA             OUT = 100
	CANCEL_WSH_META_DATA          OUT = 101
	REQ_WSH_EVENT_DATA            OUT = 102
	CANCEL_WSH_EVENT_DATA         OUT = 103
	REQ_USER_INFO                 OUT = 104
	REQ_CURRENT_TIME_IN_MILLIS    OUT = 105
)

// TWS New Bulletins constants
const NEWS_MSG int64 = 1             // standard IB news bulleting message
const EXCHANGE_AVAIL_MSG int64 = 2   // control message specifying that an exchange is available for trading
const EXCHANGE_UNAVAIL_MSG int64 = 3 // control message specifying that an exchange is unavailable for trading

const PROTOBUF_MSG_ID int64 = 200

var PROTOBUF_MSG_IDS = map[OUT]Version{
	REQ_EXECUTIONS:    MIN_SERVER_VER_PROTOBUF,
	PLACE_ORDER:       MIN_SERVER_VER_PROTOBUF_PLACE_ORDER,
	CANCEL_ORDER:      MIN_SERVER_VER_PROTOBUF_PLACE_ORDER,
	REQ_GLOBAL_CANCEL: MIN_SERVER_VER_PROTOBUF_PLACE_ORDER,
}
