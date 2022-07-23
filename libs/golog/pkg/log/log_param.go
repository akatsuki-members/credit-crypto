package log

// Log Params with key and value
type LogParam struct {
	Key   string
	Value interface{}
}

func NewParamLog(key string, value interface{}) LogParam {
	return LogParam{Key: key, Value: value}
}
