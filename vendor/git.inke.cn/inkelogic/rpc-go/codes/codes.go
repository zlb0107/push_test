package codes

type Code uint32

const (
	Success              Code = 0
	ParseRequestMessage       = 1
	ParseResponseMessage      = 2
	NotSpecifyMethodName      = 5
	ParseMethodName           = 6
	FoundService              = 7
	FoundMethod               = 8
	ChannelBroken             = 9
	ConnectionClosed          = 10
	RequestTimeout            = 11 // request timeout
	FailedPrecondition        = 100
	Unavailable               = 101
	Canceled                  = 102
	DeadlineExceeded          = 103
	ConfigLb                  = 104
	BreakerOpen               = 105
	FromUser                  = 1000
	UnKnown                   = 1001
)
