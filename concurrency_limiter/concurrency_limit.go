package concurrencylimiter

import "github.com/he-end/simproute/routes/response"

type ConcurrenctLimit struct {
	limitReq             chan struct{}
	defaultResponseError *DefaultResErr
	responser            response.ResponseHandler
}
type DefaultResErr struct {
	Status     string
	StatusCode int
	ErrCode    string
	Message    string
	Details    string
}

func NewConcurrencyLimit(capacityPerSecond int64, defErr *DefaultResErr, responser response.ResponseHandler) *ConcurrenctLimit {
	l := &ConcurrenctLimit{
		defaultResponseError: defErr,
		responser:            responser,
	}
	if capacityPerSecond <= 0 {
		capacityPerSecond = 500000
	}
	l.limitReq = make(chan struct{}, capacityPerSecond)
	return l
}
