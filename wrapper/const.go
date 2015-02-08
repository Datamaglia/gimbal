package wrapper

type ResultStatus int

const (
	FAILURE ResultStatus = iota
	SUCCESS
	WARNING
	UNKNOWN
)
