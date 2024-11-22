package ibapi

const (
	NO_VALID_ID int64 = -1
)

// CodeMsgPair is IB internal errors.
type CodeMsgPair struct {
	Code int64
	Msg  string
}

func (cmp CodeMsgPair) Error() string {
	return cmp.Msg
}

func (cmp CodeMsgPair) Equal(other CodeMsgPair) bool {
	return cmp.Code != 0 && other.Code != 0 && cmp.Code == other.Code
}

var (
	ALREADY_CONNECTED        = CodeMsgPair{501, "Already connected."}
	CONNECT_FAIL             = CodeMsgPair{502, "Couldn't connect to TWS. Confirm that 'Enable ActiveX and Socket EClients' is enabled and connection port is the same as 'Socket Port' on the TWS 'Edit->Global Configuration...->API->Settings' menu. Live Trading ports: TWS: 7496; IB Gateway: 4001. Simulated Trading ports for new installations of version 954.1 or newer:  TWS: 7497; IB Gateway: 4002"}
	UPDATE_TWS               = CodeMsgPair{503, "The TWS is out of date and must be upgraded."}
	NOT_CONNECTED            = CodeMsgPair{504, "Not connected"}
	UNKNOWN_ID               = CodeMsgPair{505, "Fatal Error: Unknown message id."}
	UNSUPPORTED_VERSION      = CodeMsgPair{506, "Unsupported version"}
	BAD_LENGTH               = CodeMsgPair{507, "Bad message length"}
	BAD_MESSAGE              = CodeMsgPair{508, "Bad message"}
	SOCKET_EXCEPTION         = CodeMsgPair{509, "Exception caught while reading socket - "}
	FAIL_CREATE_SOCK         = CodeMsgPair{520, "Failed to create socket"}
	SSL_FAIL                 = CodeMsgPair{530, "SSL specific error: "}
	INVALID_SYMBOL           = CodeMsgPair{579, "Invalid symbol in string - "}
	FA_PROFILE_NOT_SUPPORTED = CodeMsgPair{585, "FA Profile is not supported anymore, use FA Group instead - "}
)
