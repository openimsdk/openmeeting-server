package constant

const (
	SuccessCode int32 = 0

	// common error code: -10001 ~ -19999
	ServerErrorCode        int32 = -10001
	RequestDecodeErrorCode int32 = -10002

	// db error code: -20001 ~ -29999
	DbTransactionErrorCode int32 = -20001
)

var RetCode2MsgMapper = map[int32]string{
	SuccessCode:     "success",
	ServerErrorCode: "Service unavailable.",
}
