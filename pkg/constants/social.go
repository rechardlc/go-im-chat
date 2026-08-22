package constants
type HandlerResult int

const (
	NoHandlerResult HandlerResult = iota + 1
	Pass
	Refuse
	Cancel
)