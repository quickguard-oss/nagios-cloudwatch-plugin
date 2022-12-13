package alert

type ReturnCode int

const (
	OK ReturnCode = iota
	Warning
	Critical
	Unknown
)

func (r ReturnCode) String() string {
	switch r {
	case OK:
		return "OK"
	case Warning:
		return "WARNING"
	case Critical:
		return "CRITICAL"
	case Unknown:
		return "UNKNOWN"
	default:
		return "-"
	}
}
