package main

import (
	"os"
)

func IsRedirected() (stdout, stderr bool) {
	o, _ := os.Stdout.Stat()
	stdout = (o.Mode() & os.ModeCharDevice) == 0

	o, _ = os.Stderr.Stat()
	stderr = (o.Mode() & os.ModeCharDevice) == 0

	return
}
