package utils

import "log"

func LogVerbose(verbose bool, format string, v ...interface{}) {
	if verbose {
		log.Printf(format, v...)
	}
}
