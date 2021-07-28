package exception

// Exception : Exception message.
type Exception struct {
	Scope     ExceptionScope
	InnerCode int
	Message   string
}

// ExceptionScope : A exception's scope.
// this element's :
//  Information is 0;
// Warning is 1;
// Error is 2.
type ExceptionScope int

const (
	// Information : Is not a exception.
	Information ExceptionScope = 0 << iota
	// Warning : Is a exception, but program can continue to run.
	Warning
	// Error : This will cause get result is not ture.
	Error
	// FatalError ::This will cause the run to stop.
	FatalError
)

// ResetCode : Re write the innercode of this exception.
func (ex *Exception) ResetCode(newCode int) *Exception {
	ex.InnerCode = newCode
	return ex
}

// NewException is init a error info entity.
func NewException(scope ExceptionScope, innerCode int, msg string) *Exception {
	return &Exception{scope, innerCode, msg}
}
