package wrapper

type Result struct {
	CheckName string
	Status    ResultStatus
	Expected  interface{}
	Observed  interface{}
}

func (r *Result) Success() bool {
	return r.Status == SUCCESS
}
