package info

import "fmt"

// Version number.
var Version = "0.0.0"

// Release release num
var Release = ""

// BuildTime label of build time.
var BuildTime = ""

// ForPrintFull Returns formated version and build time string for print.
func ForPrintFull() string {
	return fmt.Sprintf("cataloger v%s\nBuild time %s\n", Version, BuildTime)
}

// ForPrint Returns formated version string for print.
func ForPrint() string {
	return fmt.Sprintf("cataloger v%s\n", Version)
}
