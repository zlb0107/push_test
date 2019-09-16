package client

var rpcContextKey = "_rpc_context_key"

type rpcContext struct {
	Endpoint string
	Request  interface{}
	Response interface{}
	// for KeepTrying
	retrymax int
	retry    bool
	host     string
}

func (r *rpcContext) KeepTrying(n int, err error) bool {
	if n < r.retrymax && r.retry {
		r.retry = false
		return true
	}
	return false
}

func (r *rpcContext) Retry() {
	r.retry = true
}
