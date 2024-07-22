package endpoint

type ProtectedStruct struct {
	Completed   int
	NotRelevant int
	Blocked     int
}

type RiskStruct struct {
	Allowed int
}

type ErrorStruct struct {
	Unexpected      int
	TimeoutExceeded int
}

var Protected = ProtectedStruct{
	Completed:   50,
	NotRelevant: 51,
	Blocked:     52,
}

var Errors = ErrorStruct{
	Unexpected:      49,
	TimeoutExceeded: 48,
}

var Risk = RiskStruct{
	Allowed: 100,
}
