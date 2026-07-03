package cmd

type Code int

const (
	ExitOK Code = 0

	ExitSoftware Code = 70 // EX_SOFTWARE (sysexits)
	ExitUsage    Code = 64 // EX_USAGE
	ExitConfig   Code = 78 // EX_CONFIG

	ExitUnavailable Code = 69 // EX_UNAVAILABLE
	ExitTempFail    Code = 75 // EX_TEMPFAIL
)
