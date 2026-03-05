package exitcode

type Code int

const (
	OK Code = 0

	Software Code = 70 // EX_SOFTWARE (sysexits)
	Usage    Code = 64 // EX_USAGE
	Config   Code = 78 // EX_CONFIG

	Unavailable Code = 69 // EX_UNAVAILABLE
	TempFail    Code = 75 // EX_TEMPFAIL
)
